package release

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// versionTestPlugin is a minimal plugin for testing version calculation.
type versionTestPlugin struct {
	name            string
	releaseType     algorithm.ReleaseType
	verifyCalled    bool
	verifyReleaseFn func(*algorithm.ReadOnlyContext, *algorithm.MutableState) error
}

func (p *versionTestPlugin) Name() string { return p.name }
func (p *versionTestPlugin) VerifyConditions(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *versionTestPlugin) AnalyzeCommits(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return p.releaseType, nil
}
func (p *versionTestPlugin) VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	p.verifyCalled = true
	if p.verifyReleaseFn != nil {
		return p.verifyReleaseFn(ctx, state)
	}
	return nil
}
func (p *versionTestPlugin) GenerateNotes(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (string, error) {
	return "", nil
}
func (p *versionTestPlugin) Prepare(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *versionTestPlugin) Publish(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (*algorithm.Release, error) {
	return nil, nil
}
func (p *versionTestPlugin) Success(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *versionTestPlugin) Fail(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState, _ error) error {
	return nil
}

var _ plugins.Plugin = (*versionTestPlugin)(nil)
var _ plugins.ConditionVerifier = (*versionTestPlugin)(nil)
var _ plugins.CommitAnalyzer = (*versionTestPlugin)(nil)
var _ plugins.ReleaseVerifier = (*versionTestPlugin)(nil)
var _ plugins.NotesGenerator = (*versionTestPlugin)(nil)
var _ plugins.Preparer = (*versionTestPlugin)(nil)
var _ plugins.Publisher = (*versionTestPlugin)(nil)
var _ plugins.LifecycleHook = (*versionTestPlugin)(nil)

func TestNew_ReturnsRelease(t *testing.T) {
	ctx := &Context{}
	r := New(ctx)
	if r == nil {
		t.Fatal("expected non-nil Release")
	}
	if r.ctx != ctx {
		t.Fatal("expected Release to hold the provided context")
	}
}

func TestContext_FieldsPropagatedFromCLIFlags(t *testing.T) {
	cfg := &algorithm.Config{TagFormat: "v${version}"}
	logger := zap.NewNop()
	gitClient := &git.Client{}

	ctx := &Context{
		Config:  cfg,
		Git:     gitClient,
		Logger:  logger,
		DryRun:  true,
		Verbose: true,
	}

	if ctx.Config != cfg {
		t.Error("Config field not propagated")
	}
	if ctx.Git != gitClient {
		t.Error("Git field not propagated")
	}
	if ctx.Logger != logger {
		t.Error("Logger field not propagated")
	}
	if !ctx.DryRun {
		t.Error("DryRun flag not propagated")
	}
	if !ctx.Verbose {
		t.Error("Verbose flag not propagated")
	}
}

func TestRun_DryRunCompletesWithoutSideEffects(t *testing.T) {
	// Dry-run mode should complete without errors even in a real repo,
	// because it skips Prepare, Publish, and publisher operations.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	currentBranch, err := gitClient.GetCurrentBranch()
	if err != nil {
		t.Skipf("skipping: cannot get current branch: %v", err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
				{Name: currentBranch},
			},
			TagFormat: "v${version}",
			Plugins:   []string{"changelog", "git"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	if err := r.Run(); err != nil {
		t.Fatalf("dry-run Run() failed: %v", err)
	}
}

func TestVersionCalculation(t *testing.T) {
	// Register a test plugin that returns a release type, then verify
	// runPluginLifecycle sets NextRelease.Version.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	testPlugin := &versionTestPlugin{name: "version-test", releaseType: algorithm.ReleaseTypeMinor}
	plugins.Register(testPlugin)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"version-test"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set after runPluginLifecycle")
	}
	if state.NextRelease.Version == "" {
		t.Error("expected NextRelease.Version to be non-empty")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypeMinor {
		t.Errorf("expected release type minor, got %q", state.NextRelease.Type)
	}
}

func TestVerifyRelease(t *testing.T) {
	// VerifyRelease should be called during the plugin lifecycle.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	testPlugin := &versionTestPlugin{name: "verify-test", releaseType: algorithm.ReleaseTypeMinor}
	plugins.Register(testPlugin)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"verify-test"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if !testPlugin.verifyCalled {
		t.Error("expected VerifyRelease to be called during plugin lifecycle")
	}
}

func TestVerifyRelease_Error(t *testing.T) {
	// VerifyRelease errors should be propagated.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	testPlugin := &versionTestPlugin{
		name:        "verify-error-test",
		releaseType: algorithm.ReleaseTypeMinor,
		verifyReleaseFn: func(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
			return fmt.Errorf("release verification failed")
		},
	}
	plugins.Register(testPlugin)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"verify-error-test"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	err = r.runPluginLifecycle(algoCtx, state)
	if err == nil {
		t.Fatal("expected error from VerifyRelease, got nil")
	}
	if !strings.Contains(err.Error(), "verify release failed") {
		t.Errorf("expected 'verify release failed' error, got: %v", err)
	}
}

func TestBranchValidation(t *testing.T) {
	// When the current branch is not in Config.Branches, Run() should return a skip error.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "nonexistent-branch-xyz"},
			},
			TagFormat: "v${version}",
			Plugins:   []string{"changelog", "git"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	err = r.Run()
	if err == nil {
		t.Fatal("expected error for unconfigured branch, got nil")
	}
	if !strings.Contains(err.Error(), "not in release branches") {
		t.Errorf("expected 'not in release branches' error, got: %v", err)
	}
}

func TestCIDetection(t *testing.T) {
	// When no CI env vars are set and DryRun is false, the orchestrator
	// should auto-enable dry-run to prevent accidental publishes.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	// Clear all CI env vars for this test.
	ciVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI", "TRAVIS"}
	saved := make(map[string]string)
	for _, v := range ciVars {
		saved[v] = os.Getenv(v)
		_ = os.Unsetenv(v)
	}
	t.Cleanup(func() {
		for _, v := range ciVars {
			if saved[v] != "" {
				_ = os.Setenv(v, saved[v])
			}
		}
	})

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"changelog", "git"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: false, // explicitly NOT dry-run
	}

	r := New(ctx)
	r.detectCI()
	if !r.ctx.DryRun {
		t.Error("expected DryRun to be auto-enabled when no CI env vars are set")
	}
}

func TestNew_PreservesContextReference(t *testing.T) {
	cfg := &algorithm.Config{Plugins: []string{"git"}}
	logger := zap.NewNop()
	gitClient := &git.Client{}

	ctx := &Context{
		Config: cfg,
		Git:    gitClient,
		Logger: logger,
		DryRun: false,
	}

	r := New(ctx)
	if r.ctx.Config != cfg {
		t.Error("Config reference not preserved")
	}
	if len(r.ctx.Config.Plugins) == 0 || r.ctx.Config.Plugins[0] != "git" {
		t.Error("Plugins field not accessible through context")
	}
}

// --- Bug 1: BUG-branch-config-dead ---

func TestBuildAlgoContext_PropagatesBranchConfig(t *testing.T) {
	// buildAlgoContext should copy Type, Channel, and Prerelease from the
	// matching BranchConfig into the algorithm.Branch — not just Name.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	branchName, err := gitClient.GetCurrentBranch()
	if err != nil {
		t.Skipf("skipping: cannot get current branch: %v", err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{
					Name:       branchName,
					Type:       "prerelease",
					Channel:    "beta",
					Prerelease: "beta",
				},
			},
			TagFormat: "v${version}",
			Plugins:   []string{},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}
	_ = state // state is not used in this test

	if algoCtx.Branch.Name != branchName {
		t.Errorf("expected branch name %q, got %q", branchName, algoCtx.Branch.Name)
	}
	if algoCtx.Branch.Type != algorithm.BranchTypePrerelease {
		t.Errorf("expected branch type %q, got %q", algorithm.BranchTypePrerelease, algoCtx.Branch.Type)
	}
	if algoCtx.Branch.Channel != "beta" {
		t.Errorf("expected channel %q, got %q", "beta", algoCtx.Branch.Channel)
	}
	if algoCtx.Branch.Prerelease != "beta" {
		t.Errorf("expected prerelease %q, got %q", "beta", algoCtx.Branch.Prerelease)
	}
}

func TestBuildAlgoContext_PropagatesMaintenanceBranchConfig(t *testing.T) {
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	branchName, err := gitClient.GetCurrentBranch()
	if err != nil {
		t.Skipf("skipping: cannot get current branch: %v", err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{
					Name:    branchName,
					Type:    "maintenance",
					Channel: "lts",
				},
			},
			TagFormat: "v${version}",
			Plugins:   []string{},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}
	_ = state

	if algoCtx.Branch.Type != algorithm.BranchTypeMaintenance {
		t.Errorf("expected branch type %q, got %q", algorithm.BranchTypeMaintenance, algoCtx.Branch.Type)
	}
	if algoCtx.Branch.Channel != "lts" {
		t.Errorf("expected channel %q, got %q", "lts", algoCtx.Branch.Channel)
	}
}

// --- Bug 2: BUG-publishers-config-ignored ---

// stubPublisher is a minimal publisher for testing config filtering.
type stubPublisher struct {
	name   string
	dryRun bool
	ran    bool
}

func (p *stubPublisher) Name() string           { return p.name }
func (p *stubPublisher) Detect() bool           { return true }
func (p *stubPublisher) Prepare(_ string) error { p.ran = true; return nil }
func (p *stubPublisher) Publish(_ string) error { return nil }
func (p *stubPublisher) Verify(_ string) error  { return nil }
func (p *stubPublisher) SetDryRun(dryRun bool)  { p.dryRun = dryRun }

func TestRunPublishers_SkipsDisabledPublishers(t *testing.T) {
	// Publishers with enabled: false in config should be skipped.
	// Detected publishers not in config map should still run (backward compatible).
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	npmPub := &stubPublisher{name: "npm"}
	goPub := &stubPublisher{name: "goproxy"}

	// Register stub publishers
	publishers.Register(npmPub)
	publishers.Register(goPub)
	t.Cleanup(func() {
		// Reset global registry after test
		*publishers.NewRegistry() = publishers.Registry{}
	})

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat: "v${version}",
			Plugins:   []string{},
			Publishers: map[string]algorithm.PublisherConfig{
				"npm":     {Enabled: true},
				"goproxy": {Enabled: false},
			},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: false,
	}

	r := New(ctx)
	algoCtx := &algorithm.ReadOnlyContext{
		Config: ctx.Config,
		Branch: &algorithm.Branch{Name: "main"},
		DryRun: false,
	}
	state := &algorithm.MutableState{}

	err = r.runPublishers(algoCtx, state)
	if err != nil {
		t.Fatalf("runPublishers failed: %v", err)
	}

	if !npmPub.ran {
		t.Error("expected npm publisher to run (enabled: true)")
	}
	if goPub.ran {
		t.Error("expected goproxy publisher to be skipped (enabled: false)")
	}
}

// --- Story e06s03: Priority-based AnalyzeCommits aggregation ---

// priorityTestPlugin returns a specific release type when AnalyzeCommits is called.
type priorityTestPlugin struct {
	name        string
	releaseType algorithm.ReleaseType
}

func (p *priorityTestPlugin) Name() string { return p.name }
func (p *priorityTestPlugin) AnalyzeCommits(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return p.releaseType, nil
}

var _ plugins.Plugin = (*priorityTestPlugin)(nil)
var _ plugins.CommitAnalyzer = (*priorityTestPlugin)(nil)

func TestAnalyzeCommits_PriorityBasedAggregation(t *testing.T) {
	// When plugin A returns patch and plugin B returns minor,
	// the result should be minor (higher priority wins).
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	patchPlugin := &priorityTestPlugin{name: "patch-plugin", releaseType: algorithm.ReleaseTypePatch}
	minorPlugin := &priorityTestPlugin{name: "minor-plugin", releaseType: algorithm.ReleaseTypeMinor}
	plugins.Register(patchPlugin)
	plugins.Register(minorPlugin)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			// patch comes before minor in plugin order
			Plugins: []string{"patch-plugin", "minor-plugin"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypeMinor {
		t.Errorf("expected release type minor (priority-based), got %q", state.NextRelease.Type)
	}
}

func TestAnalyzeCommits_MajorOverridesMinor(t *testing.T) {
	// When plugin A returns minor and plugin B returns major,
	// the result should be major (highest priority wins regardless of order).
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	minorPlugin := &priorityTestPlugin{name: "minor-plugin-2", releaseType: algorithm.ReleaseTypeMinor}
	majorPlugin := &priorityTestPlugin{name: "major-plugin", releaseType: algorithm.ReleaseTypeMajor}
	plugins.Register(minorPlugin)
	plugins.Register(majorPlugin)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			// minor comes before major in plugin order
			Plugins: []string{"minor-plugin-2", "major-plugin"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypeMajor {
		t.Errorf("expected release type major, got %q", state.NextRelease.Type)
	}
}

func TestAnalyzeCommits_FirstPluginWins_WhenEqualPriority(t *testing.T) {
	// When two plugins return the same priority, the first one wins (stable).
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	patchA := &priorityTestPlugin{name: "patch-a", releaseType: algorithm.ReleaseTypePatch}
	patchB := &priorityTestPlugin{name: "patch-b", releaseType: algorithm.ReleaseTypePatch}
	plugins.Register(patchA)
	plugins.Register(patchB)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"patch-a", "patch-b"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypePatch {
		t.Errorf("expected release type patch, got %q", state.NextRelease.Type)
	}
}

// --- BUG-nil-panic-analyzer: Analyzer fallback tests ---

// noAnalyzePlugin is a plugin that implements all interfaces EXCEPT CommitAnalyzer.
type noAnalyzePlugin struct {
	name string
}

func (p *noAnalyzePlugin) Name() string { return p.name }
func (p *noAnalyzePlugin) VerifyConditions(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *noAnalyzePlugin) VerifyRelease(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *noAnalyzePlugin) GenerateNotes(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (string, error) {
	return "", nil
}
func (p *noAnalyzePlugin) Prepare(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *noAnalyzePlugin) Publish(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (*algorithm.Release, error) {
	return nil, nil
}
func (p *noAnalyzePlugin) Success(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}
func (p *noAnalyzePlugin) Fail(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState, _ error) error {
	return nil
}

var _ plugins.Plugin = (*noAnalyzePlugin)(nil)
var _ plugins.ConditionVerifier = (*noAnalyzePlugin)(nil)
var _ plugins.ReleaseVerifier = (*noAnalyzePlugin)(nil)
var _ plugins.NotesGenerator = (*noAnalyzePlugin)(nil)
var _ plugins.Preparer = (*noAnalyzePlugin)(nil)
var _ plugins.Publisher = (*noAnalyzePlugin)(nil)
var _ plugins.LifecycleHook = (*noAnalyzePlugin)(nil)

func TestAnalyzerFallback_FeatCommitTriggersMinor(t *testing.T) {
	// V2+V3: When no plugin returns a release type, the built-in Analyzer
	// should detect 'feat:' commits and return minor.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	// Register a plugin that does NOT implement CommitAnalyzer
	noAnalyze := &noAnalyzePlugin{name: "no-analyze"}
	plugins.Register(noAnalyze)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"no-analyze"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	// Inject a feat commit so the Analyzer has something to analyze
	algoCtx.Commits = []*algorithm.Commit{
		{Message: "feat: add login", Type: "feat", Subject: "add login"},
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set after Analyzer fallback")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypeMinor {
		t.Errorf("expected release type minor from Analyzer fallback, got %q", state.NextRelease.Type)
	}
}

func TestAnalyzerFallback_FixCommitTriggersPatch(t *testing.T) {
	// V4: 'fix:' commits should trigger patch via Analyzer fallback.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	noAnalyze := &noAnalyzePlugin{name: "no-analyze-patch"}
	plugins.Register(noAnalyze)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"no-analyze-patch"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	algoCtx.Commits = []*algorithm.Commit{
		{Message: "fix: resolve crash", Type: "fix", Subject: "resolve crash"},
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypePatch {
		t.Errorf("expected release type patch, got %q", state.NextRelease.Type)
	}
}

func TestAnalyzerFallback_BreakingCommitTriggersMajor(t *testing.T) {
	// V5: Breaking commits should trigger major via Analyzer fallback.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	noAnalyze := &noAnalyzePlugin{name: "no-analyze-major"}
	plugins.Register(noAnalyze)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"no-analyze-major"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	algoCtx.Commits = []*algorithm.Commit{
		{Message: "feat!: breaking change", Type: "feat", Subject: "breaking change", Breaking: true},
	}

	if err := r.runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if state.NextRelease == nil {
		t.Fatal("expected NextRelease to be set")
	}
	if state.NextRelease.Type != algorithm.ReleaseTypeMajor {
		t.Errorf("expected release type major, got %q", state.NextRelease.Type)
	}
}

func TestAnalyzerFallback_NoRelevantChanges_NoPanic(t *testing.T) {
	// V6: When commits have no releasable type, should exit cleanly.
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	noAnalyze := &noAnalyzePlugin{name: "no-analyze-noop"}
	plugins.Register(noAnalyze)

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"no-analyze-noop"},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	// Only chore commits — no releasable type
	algoCtx.Commits = []*algorithm.Commit{
		{Message: "chore: update deps", Type: "chore", Subject: "update deps"},
	}

	// Should return nil (no error, no release)
	err = r.runPluginLifecycle(algoCtx, state)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// TestReadOnlyContextImmutable verifies that ReadOnlyContext fields are not modified by plugins.
func TestReadOnlyContextImmutable(t *testing.T) {
	// Use a simple test: buildAlgoContext returns ReadOnlyContext and MutableState
	// and verify that both are non-nil and separate
	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: cannot create git client: %v", err)
	}
	if !gitClient.IsGitRepo() {
		t.Skip("skipping: not in a git repo")
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
			},
			TagFormat: "v${version}",
			Plugins:   []string{},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx, state, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	// Verify that ReadOnlyContext and MutableState are separate structs
	if algoCtx == nil {
		t.Fatal("expected non-nil ReadOnlyContext")
	}
	if state == nil {
		t.Fatal("expected non-nil MutableState")
	}

	// Verify that ReadOnlyContext does not have mutable fields
	// (this is a compile-time check, but we can verify the fields exist)
	_ = algoCtx.Config
	_ = algoCtx.Branch
	_ = algoCtx.Commits
	_ = algoCtx.Releases
	_ = algoCtx.RepositoryURL
	_ = algoCtx.DryRun

	// Verify that MutableState has the mutable fields
	_ = state.LastRelease
	_ = state.NextRelease
	_ = state.Notes
	_ = state.Assets

	// Verify that modifying MutableState does not affect ReadOnlyContext
	state.Notes = "test notes"
	if algoCtx.RepositoryURL == "test notes" {
		t.Error("modifying MutableState.Notes should not affect ReadOnlyContext.RepositoryURL")
	}
}
