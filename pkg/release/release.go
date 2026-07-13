package release

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// Context holds CLI-level inputs for the release orchestrator.
type Context struct {
	Config  *algorithm.Config
	Git     *git.Client
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
func (r *Release) Run() error {
	branchName, err := r.ctx.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	lastRelease, err := r.ctx.Git.GetLastRelease(r.ctx.Config.TagFormat)
	if err != nil {
		return fmt.Errorf("failed to get last release: %w", err)
	}

	from := ""
	if lastRelease != nil {
		from = lastRelease.GitTag
	}
	commits, err := r.ctx.Git.GetCommits(from, "HEAD")
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	repoURL, err := r.ctx.Git.GetRepositoryURL()
	if err != nil {
		return fmt.Errorf("failed to get repository URL: %w", err)
	}

	algoCtx := &algorithm.Context{
		Config:        r.ctx.Config,
		Branch:        &algorithm.Branch{Name: branchName},
		LastRelease:   lastRelease,
		NextRelease:   nil,
		Commits:       commits,
		Releases:      nil,
		RepositoryURL: repoURL,
		DryRun:        r.ctx.DryRun,
	}

	if err := r.runPluginLifecycle(algoCtx); err != nil {
		r.callFailHooks(algoCtx, err)
		return err
	}

	if err := r.runPublishers(algoCtx); err != nil {
		r.callFailHooks(algoCtx, err)
		return err
	}

	if !r.ctx.DryRun {
		r.callSuccessHooks(algoCtx)
	}

	return nil
}

func (r *Release) runPluginLifecycle(ctx *algorithm.Context) error {
	pluginNames := r.ctx.Config.Plugins

	// Phase 1: VerifyConditions
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if err := p.VerifyConditions(ctx); err != nil {
			return fmt.Errorf("plugin %q verify conditions failed: %w", name, err)
		}
	}

	// Phase 2: AnalyzeCommits — collect release type from plugins
	var releaseType algorithm.ReleaseType
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		rt, err := p.AnalyzeCommits(ctx)
		if err != nil {
			return fmt.Errorf("plugin %q analyze commits failed: %w", name, err)
		}
		if rt != "" {
			releaseType = rt
		}
	}

	// Phase 3: GenerateNotes — collect release notes from plugins
	var notes string
	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		n, err := p.GenerateNotes(ctx)
		if err != nil {
			return fmt.Errorf("plugin %q generate notes failed: %w", name, err)
		}
		if n != "" {
			if notes != "" {
				notes += "\n"
			}
			notes += n
		}
	}

	if releaseType != "" || notes != "" {
		if ctx.NextRelease == nil {
			ctx.NextRelease = &algorithm.Release{}
		}
		if releaseType != "" {
			ctx.NextRelease.Type = releaseType
		}
		if notes != "" {
			ctx.NextRelease.Notes = notes
		}
	}

	// Phase 4 + 5: Prepare and Publish (skipped in dry-run mode)
	if r.ctx.DryRun {
		r.ctx.Logger.Info("Dry run: skipping plugin prepare and publish")
		return nil
	}

	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		if err := p.Prepare(ctx); err != nil {
			return fmt.Errorf("plugin %q prepare failed: %w", name, err)
		}
	}

	for _, name := range pluginNames {
		p, err := plugins.Get(name)
		if err != nil {
			return fmt.Errorf("plugin %q not found: %w", name, err)
		}
		rel, err := p.Publish(ctx)
		if err != nil {
			return fmt.Errorf("plugin %q publish failed: %w", name, err)
		}
		if rel != nil {
			ctx.NextRelease = rel
		}
	}

	return nil
}

func (r *Release) runPublishers(ctx *algorithm.Context) error {
	detected := publishers.Detect()

	for _, p := range detected {
		if setter, ok := p.(interface{ SetDryRun(bool) }); ok {
			setter.SetDryRun(r.ctx.DryRun)
		}
	}

	if r.ctx.DryRun {
		r.ctx.Logger.Info("Dry run: skipping publisher operations")
		return nil
	}

	version := ""
	if ctx.NextRelease != nil {
		version = ctx.NextRelease.Version
	}

	for _, pub := range detected {
		if err := pub.Prepare(version); err != nil {
			return fmt.Errorf("publisher %q prepare failed: %w", pub.Name(), err)
		}
	}

	for _, pub := range detected {
		if err := pub.Publish(version); err != nil {
			return fmt.Errorf("publisher %q publish failed: %w", pub.Name(), err)
		}
	}

	for _, pub := range detected {
		if err := pub.Verify(version); err != nil {
			return fmt.Errorf("publisher %q verify failed: %w", pub.Name(), err)
		}
	}

	return nil
}

func (r *Release) callFailHooks(ctx *algorithm.Context, originalErr error) {
	r.ctx.Logger.Error("Release failed, running fail hooks", zap.Error(originalErr))
	for _, name := range plugins.List() {
		p, err := plugins.Get(name)
		if err != nil {
			continue
		}
		if err := p.Fail(ctx, originalErr); err != nil {
			r.ctx.Logger.Error("Fail hook failed", zap.String("plugin", name), zap.Error(err))
		}
	}
}

func (r *Release) callSuccessHooks(ctx *algorithm.Context) {
	for _, name := range plugins.List() {
		p, err := plugins.Get(name)
		if err != nil {
			continue
		}
		if err := p.Success(ctx); err != nil {
			r.ctx.Logger.Error("Success hook failed", zap.String("plugin", name), zap.Error(err))
		}
	}
}
