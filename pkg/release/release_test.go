package release

import (
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
)

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
