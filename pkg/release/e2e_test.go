// story: e08s03
package release

import (
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func TestE2E_DryRunPipeline(t *testing.T) {
	testrepo.ScrubEnv(t)
	t.Setenv("CI", "true")

	repo := testrepo.Init(t)
	repo.SetRemoteOrigin("https://github.com/example/e2e-test.git")
	repo.Commit("feat: initial feature")
	repo.Tag("v0.1.0")
	repo.Commit("feat: next release")

	gitClient, err := git.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches:       []algorithm.BranchConfig{{Name: "main"}},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"changelog", "git"},
			CommitTypes:    algorithm.DefaultCommitTypes(),
			Publishers:     map[string]algorithm.PublisherConfig{},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	if err := r.Run(); err != nil {
		t.Fatalf("E2E dry-run failed: %v", err)
	}
}

func TestE2E_NoPublishSideEffects(t *testing.T) {
	testrepo.ScrubEnv(t)
	t.Setenv("CI", "true")

	repo := testrepo.Init(t)
	repo.SetRemoteOrigin("https://github.com/example/e2e-test.git")
	repo.Commit("feat: dry run only")

	gitClient, err := git.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches:       []algorithm.BranchConfig{{Name: "main"}},
			TagFormat:      "v${version}",
			InitialVersion: "0.1.0",
			Plugins:        []string{"changelog"},
			CommitTypes:    algorithm.DefaultCommitTypes(),
			Publishers:     map[string]algorithm.PublisherConfig{"npm": {Enabled: false}},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	if err := r.Run(); err != nil {
		t.Fatalf("dry-run should complete: %v", err)
	}

	tags, err := gitClient.GetTags()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) > 0 {
		t.Fatalf("dry-run should not create tags, got %v", tags)
	}
	_ = repo
}
