// story: e08s03
package release

import (
	"os"
	"slices"
	"strings"
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

// story: e09s05
func TestE2E_VersionFileSync(t *testing.T) {
	testrepo.ScrubEnv(t)
	t.Setenv("CI", "true")

	repo := testrepo.Init(t)
	// A local bare remote keeps Publish/PostPublish pushes off the network.
	remote := t.TempDir()
	repo.RunGit("init", "--bare", remote)
	repo.SetRemoteOrigin(remote)
	repo.RunGit("config", "push.default", "current")

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
			Plugins:        []string{"versionfile", "git"},
			CommitTypes:    algorithm.DefaultCommitTypes(),
			Publishers:     map[string]algorithm.PublisherConfig{},
			PluginConfigs: map[string]map[string]interface{}{
				"versionfile": {"path": "VERSION", "template": "{{.Version}}"},
				"git":         {"postTagAssets": []string{"VERSION"}},
			},
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: false,
	}

	if err := New(ctx).Run(); err != nil {
		t.Fatalf("E2E release failed: %v", err)
	}

	const wantVersion = "0.2.0"

	// The versionfile plugin writes the computed version after the tag exists.
	got, err := os.ReadFile("VERSION")
	if err != nil {
		t.Fatalf("VERSION file not written: %v", err)
	}
	if v := strings.TrimSpace(string(got)); v != wantVersion {
		t.Fatalf("VERSION = %q, want %q", v, wantVersion)
	}

	tags, err := gitClient.GetTags()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(tags, "v"+wantVersion) {
		t.Fatalf("expected tag v%s, got %v", wantVersion, tags)
	}

	// Tag/commit relationship: the tag is created at HEAD during Publish, and
	// the version-file commit lands on top of it in PostPublish.
	tagHead, err := gitClient.GetTagHead("v" + wantVersion)
	if err != nil {
		t.Fatal(err)
	}
	head, err := gitClient.GetHead()
	if err != nil {
		t.Fatal(err)
	}
	if tagHead == head {
		t.Fatalf("expected version-file commit after tag; tag and HEAD are both %s", head)
	}

	afterTag, err := gitClient.GetCommits(tagHead, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(afterTag) != 1 {
		t.Fatalf("expected exactly one post-tag commit, got %d", len(afterTag))
	}
	wantMsg := "chore(release): update version files to " + wantVersion + " [skip ci]"
	if afterTag[0].Message != wantMsg {
		t.Fatalf("post-tag commit message = %q, want %q", afterTag[0].Message, wantMsg)
	}
}
