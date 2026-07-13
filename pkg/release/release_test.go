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
)

// versionTestPlugin is a minimal plugin for testing version calculation.
type versionTestPlugin struct {
	name            string
	releaseType     algorithm.ReleaseType
	verifyCalled    bool
	verifyReleaseFn func(*algorithm.Context) error
}

func (p *versionTestPlugin) Name() string                                { return p.name }
func (p *versionTestPlugin) VerifyConditions(_ *algorithm.Context) error { return nil }
func (p *versionTestPlugin) AnalyzeCommits(_ *algorithm.Context) (algorithm.ReleaseType, error) {
	return p.releaseType, nil
}
func (p *versionTestPlugin) VerifyRelease(ctx *algorithm.Context) error {
	p.verifyCalled = true
	if p.verifyReleaseFn != nil {
		return p.verifyReleaseFn(ctx)
	}
	return nil
}
func (p *versionTestPlugin) GenerateNotes(_ *algorithm.Context) (string, error) { return "", nil }
func (p *versionTestPlugin) Prepare(_ *algorithm.Context) error                 { return nil }
func (p *versionTestPlugin) Publish(_ *algorithm.Context) (*algorithm.Release, error) {
	return nil, nil
}
func (p *versionTestPlugin) Success(_ *algorithm.Context) error       { return nil }
func (p *versionTestPlugin) Fail(_ *algorithm.Context, _ error) error { return nil }

var _ plugins.Plugin = (*versionTestPlugin)(nil)

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

	ctx := &Context{
		Config: &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main"},
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
	algoCtx, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if algoCtx.NextRelease == nil {
		t.Fatal("expected NextRelease to be set after runPluginLifecycle")
	}
	if algoCtx.NextRelease.Version == "" {
		t.Error("expected NextRelease.Version to be non-empty")
	}
	if algoCtx.NextRelease.Type != algorithm.ReleaseTypeMinor {
		t.Errorf("expected release type minor, got %q", algoCtx.NextRelease.Type)
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
	algoCtx, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	if err := r.runPluginLifecycle(algoCtx); err != nil {
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
		verifyReleaseFn: func(_ *algorithm.Context) error {
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
	algoCtx, err := r.buildAlgoContext()
	if err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}

	err = r.runPluginLifecycle(algoCtx)
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
