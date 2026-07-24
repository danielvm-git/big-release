package release

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// Context holds CLI-level inputs for the release orchestrator.
type Context struct {
	Config  *algorithm.Config
	Git     git.GitAPI
	Logger  *zap.Logger
	DryRun  bool
	Verbose bool
}

// Release orchestrates the full release lifecycle: plugins then publishers.
type Release struct {
	ctx *Context
}

// New creates a new Release orchestrator from a CLI-level Context.
func New(ctx *Context) *Release {
	return &Release{ctx: ctx}
}

// Run executes the full release lifecycle:
// 1. Gathers git state (branch, last release, commits, repo URL)
// 2. Runs plugins in order: verify → analyze → generate notes → prepare → publish
// 3. Runs detected publishers: prepare, publish, verify
// 4. Calls Success hooks on all plugins, or Fail hooks on error
// detectCI auto-enables dry-run when no CI environment variables are detected.
func (r *Release) detectCI() {
	if r.ctx.DryRun {
		return
	}
	ciVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI", "TRAVIS"}
	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			return
		}
	}
	r.ctx.Logger.Warn("No CI environment detected, auto-enabling dry-run mode")
	r.ctx.DryRun = true
}

// detectPR reports whether this run was triggered by a pull request.
// Checks GitHub Actions, GitLab CI, and Azure DevOps indicators.
func (r *Release) detectPR() bool {
	if os.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		return true
	}
	if os.Getenv("CI_MERGE_REQUEST_ID") != "" {
		return true
	}
	if os.Getenv("BUILD_REASON") == "PullRequest" {
		return true
	}
	return false
}

// validateBranch checks that the current branch is in the configured
// release branches. Matching supports both exact names and glob patterns
// (e.g. "+([0-9]).x" matches "1.x", "2.x"). This enables maintenance
// branch configurations without enumerating every versioned branch.
func (r *Release) validateBranch(branchName string) error {
	for _, bc := range r.ctx.Config.Branches {
		if matchBranch(bc.Name, branchName) {
			return nil
		}
	}
	return fmt.Errorf("branch %q not in release branches, skipping", branchName)
}

// matchBranch reports whether a configured branch pattern matches the
// actual branch name. It first tries an exact match (the common case and
// the historical behavior), then falls back to extglob-style pattern
// matching via matchBranchPattern.
func matchBranch(pattern, branch string) bool {
	if pattern == branch {
		return true
	}
	return matchBranchPattern(pattern, branch)
}

// mapBranchConfig finds the BranchConfig matching branchName (via exact or
// glob match) and maps it to an algorithm.Branch, propagating Type,
// Channel, and Prerelease fields. Returns a default Branch with only Name
// set if no config matches.
func mapBranchConfig(branchName string, configs []algorithm.BranchConfig) *algorithm.Branch {
	branch := &algorithm.Branch{Name: branchName}
	for _, bc := range configs {
		if matchBranch(bc.Name, branchName) {
			if bc.Type != "" {
				branch.Type = algorithm.BranchType(bc.Type)
			}
			branch.Channel = bc.Channel
			branch.Prerelease = bc.Prerelease
			break
		}
	}
	return branch
}

// buildAlgoContext gathers git state and constructs ReadOnlyContext and MutableState.
func (r *Release) buildAlgoContext() (*algorithm.ReadOnlyContext, *algorithm.MutableState, error) {
	branchName, err := r.ctx.Git.GetCurrentBranch()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	lastRelease, err := r.ctx.Git.GetLastRelease(r.ctx.Config.TagFormat)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get last release: %w", err)
	}

	from := ""
	if lastRelease != nil {
		from = lastRelease.GitTag
	}
	commits, err := r.ctx.Git.GetCommits(from, "HEAD")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get commits: %w", err)
	}

	repoURL, err := r.ctx.Git.GetRepositoryURL()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get repository URL: %w", err)
	}

	readOnly := &algorithm.ReadOnlyContext{
		Config:        r.ctx.Config,
		Branch:        mapBranchConfig(branchName, r.ctx.Config.Branches),
		Commits:       commits,
		Releases:      nil,
		RepositoryURL: repoURL,
		DryRun:        r.ctx.DryRun,
	}
	state := &algorithm.MutableState{
		LastRelease: lastRelease,
		NextRelease: nil,
	}
	return readOnly, state, nil
}

func (r *Release) Run() error {
	r.detectCI()
	if r.detectPR() {
		r.ctx.Logger.Info("This run was triggered by a pull request and therefore a new version won't be published")
		return nil
	}

	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		return err
	}

	if err := r.validateBranch(algoCtx.Branch.Name); err != nil {
		return err
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		r.callFailHooks(algoCtx, state, err)
		return err
	}

	if err := r.runPublishers(algoCtx, state); err != nil {
		r.callFailHooks(algoCtx, state, err)
		return err
	}

	if !r.ctx.DryRun {
		r.callSuccessHooks(algoCtx, state)
		// Expose the published version to the CI runner (GitHub Actions
		// $GITHUB_OUTPUT) so downstream steps can observe a release
		// (BUG-release-workflow-softprops-and-verbose).
		if err := r.writeStepOutput(state); err != nil {
			r.ctx.Logger.Warn("failed to write step output", zap.Error(err))
		}
	}

	return nil
}

// writeStepOutput exposes the computed version (and a published flag) to the
// enclosing CI runner by appending to the file named by $GITHUB_OUTPUT.
// Outside GitHub Actions the env var is unset and this is a silent no-op.
func (r *Release) writeStepOutput(state *algorithm.MutableState) error {
	if state == nil || state.NextRelease == nil || state.NextRelease.Version == "" {
		return nil
	}
	out := os.Getenv("GITHUB_OUTPUT")
	if out == "" {
		return nil
	}
	line := fmt.Sprintf("version=%s\npublished=true\n", state.NextRelease.Version)
	f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open GITHUB_OUTPUT: %w", err)
	}
	defer func() { _ = f.Close() }()
	if _, err := f.WriteString(line); err != nil {
		return fmt.Errorf("write GITHUB_OUTPUT: %w", err)
	}
	return nil
}

func (r *Release) runPluginLifecycle(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	pluginNames := r.ctx.Config.Plugins

	// Register GitPlugin with the GitAPI from the context
	plugins.Register(plugins.NewGitPlugin(r.ctx.Git))

	// Phase 0: Configure plugins that accept typed config (e.g. github assets).
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if cp, ok := p.(plugins.ConfigurablePlugin); ok {
			raw := ctx.Config.PluginConfigs[name]
			if raw == nil {
				raw = map[string]interface{}{}
			}
			if err := cp.Configure(raw); err != nil {
				return fmt.Errorf("plugin %q configure failed: %w", name, err)
			}
		}
	}

	// Phase 1: VerifyConditions
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if cv, ok := p.(plugins.ConditionVerifier); ok {
			if err := cv.VerifyConditions(ctx, state); err != nil {
				return fmt.Errorf("plugin %q verify conditions failed: %w", name, err)
			}
		}
	}

	// Phase 2: AnalyzeCommits — collect release type from plugins (priority-based)
	var releaseType algorithm.ReleaseType
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if ca, ok := p.(plugins.CommitAnalyzer); ok {
			rt, err := ca.AnalyzeCommits(ctx, state)
			if err != nil {
				return fmt.Errorf("plugin %q analyze commits failed: %w", name, err)
			}
			if rt != "" {
				if releaseType == "" || algorithm.BumpPriority(rt) > algorithm.BumpPriority(releaseType) {
					releaseType = rt
				}
			}
		}
	}

	// Phase 2 fallback: use built-in Analyzer if no plugin returned a release type
	if releaseType == "" {
		analyzer := algorithm.NewAnalyzer()
		releaseType = analyzer.AnalyzeCommits(ctx.Commits)
	}

	// Phase 2.5: Generate notes using the single algorithm Generator
	gen := algorithm.NewGenerator(ctx.Config.CommitTypes)
	gen.SetRepositoryURL(ctx.RepositoryURL)
	notes := gen.GenerateNotes(ctx.Commits, state.LastRelease, state.NextRelease)

	if releaseType != "" || notes != "" {
		if state.NextRelease == nil {
			state.NextRelease = &algorithm.Release{}
		}
		// Propagate the branch channel so publishers can target the right
		// distribution channel (e.g. npm dist-tag, prerelease line).
		state.NextRelease.Channel = ctx.Branch.Channel
		state.NextRelease.Branch = ctx.Branch.Name
		if releaseType != "" {
			state.NextRelease.Type = releaseType
		}
		if notes != "" {
			state.NextRelease.Notes = notes
		}
	}

	// Early exit: no relevant changes detected
	if releaseType == "" && notes == "" {
		r.ctx.Logger.Info("No relevant changes, skipping release")
		return nil
	}

	// Phase 3: Let plugins contribute additional notes (e.g. ChangelogPlugin)
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if ng, ok := p.(plugins.NotesGenerator); ok {
			n, err := ng.GenerateNotes(ctx, state)
			if err != nil {
				return fmt.Errorf("plugin %q generate notes failed: %w", name, err)
			}
			if n != "" && state.NextRelease != nil && n != state.NextRelease.Notes {
				state.NextRelease.Notes += "\n" + n
			}
		}
	}

	// Phase 3.5: CalculateNextVersion
	if state.NextRelease != nil && releaseType != "" {
		calc := algorithm.NewCalculator()
		version, err := calc.CalculateNextVersion(
			state.LastRelease, releaseType, ctx.Branch, ctx.Config.InitialVersion,
		)
		if err != nil {
			return fmt.Errorf("failed to calculate next version: %w", err)
		}
		state.NextRelease.Version = version
		// Narrate the computed version so a successful release is visible in
		// CI output (BUG-release-workflow-softprops-and-verbose).
		r.ctx.Logger.Info("Computed next release",
			zap.String("version", version),
			zap.String("type", string(releaseType)),
		)
	}

	// Phase 4: VerifyRelease
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if rv, ok := p.(plugins.ReleaseVerifier); ok {
			if err := rv.VerifyRelease(ctx, state); err != nil {
				return fmt.Errorf("plugin %q verify release failed: %w", name, err)
			}
		}
	}

	// Phase 4.5: AddChannel — record distribution channel metadata (e22).
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if cm, ok := p.(plugins.ChannelManager); ok {
			if err := cm.AddChannel(ctx, state); err != nil {
				return fmt.Errorf("plugin %q add channel failed: %w", name, err)
			}
		}
	}

	// Phase 5 + 6: Prepare and Publish (skipped in dry-run mode)
	if r.ctx.DryRun {
		r.ctx.Logger.Info("Dry run: skipping plugin prepare and publish")
		return nil
	}

	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if prep, ok := p.(plugins.Preparer); ok {
			if err := prep.Prepare(ctx, state); err != nil {
				return fmt.Errorf("plugin %q prepare failed: %w", name, err)
			}
		}
	}

	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if pub, ok := p.(plugins.Publisher); ok {
			rel, err := pub.Publish(ctx, state)
			if err != nil {
				return fmt.Errorf("plugin %q publish failed: %w", name, err)
			}
			if rel != nil {
				state.NextRelease = rel
			}
			r.ctx.Logger.Info("Plugin published", zap.String("plugin", name))
		}
	}

	return nil
}

func (r *Release) runPublishers(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	detected := publishers.Detect()

	// Filter detected publishers against config. Skip publishers with enabled: false.
	// Detected publishers not in the config map still run (backward compatible).
	filtered := make([]publishers.Publisher, 0, len(detected))
	for _, p := range detected {
		if pc, exists := r.ctx.Config.Publishers[p.Name()]; exists && !pc.Enabled {
			r.ctx.Logger.Info("Publisher disabled by config, skipping", zap.String("publisher", p.Name()))
			continue
		}
		filtered = append(filtered, p)
	}

	for _, p := range filtered {
		if setter, ok := p.(interface{ SetDryRun(bool) }); ok {
			setter.SetDryRun(r.ctx.DryRun)
		}
	}
	channel := ""
	if state.NextRelease != nil {
		channel = state.NextRelease.Channel
	}
	for _, p := range filtered {
		if setter, ok := p.(interface{ SetChannel(string) }); ok {
			setter.SetChannel(channel)
		}
	}

	if r.ctx.DryRun {
		r.ctx.Logger.Info("Dry run: skipping publisher operations")
		return nil
	}

	version := ""
	if state.NextRelease != nil {
		version = state.NextRelease.Version
	}

	for _, pub := range filtered {
		if err := pub.Prepare(version); err != nil {
			return fmt.Errorf("publisher %q prepare failed: %w", pub.Name(), err)
		}
	}

	for _, pub := range filtered {
		if err := pub.Publish(version); err != nil {
			return fmt.Errorf("publisher %q publish failed: %w", pub.Name(), err)
		}
	}

	for _, pub := range filtered {
		if err := pub.Verify(version); err != nil {
			return fmt.Errorf("publisher %q verify failed: %w", pub.Name(), err)
		}
	}

	return nil
}

func (r *Release) callFailHooks(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, originalErr error) {
	r.ctx.Logger.Error("Release failed, running fail hooks", zap.Error(originalErr))
	for _, name := range plugins.List() {
		p, err := plugins.Get(name)
		if err != nil {
			continue
		}
		if lh, ok := p.(plugins.LifecycleHook); ok {
			if err := lh.Fail(ctx, state, originalErr); err != nil {
				r.ctx.Logger.Error("Fail hook failed", zap.String("plugin", name), zap.Error(err))
			}
		}
	}
}

func (r *Release) callSuccessHooks(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) {
	for _, name := range plugins.List() {
		p, err := plugins.Get(name)
		if err != nil {
			continue
		}
		if lh, ok := p.(plugins.LifecycleHook); ok {
			if err := lh.Success(ctx, state); err != nil {
				r.ctx.Logger.Error("Success hook failed", zap.String("plugin", name), zap.Error(err))
			}
		}
	}
}
