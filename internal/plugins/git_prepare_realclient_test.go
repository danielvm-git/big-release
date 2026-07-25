// story: BUG-git-push-error-swallowed
package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

// TestGitPluginPrepare_RealClient_NoAssetsConfigSkipsCommit is an end-to-end
// check (real git.Client, not fakeGit) that Prepare() leaves an unconfigured
// GitPlugin's release branch untouched even when another plugin (e.g.
// changelog) has written an unstaged file to disk. Before this fix,
// Client.Commit() unconditionally ran `git add -A`, so any file present in
// the working tree — staged or not — got swept into a real commit and
// pushed to the release branch, which a protected branch's required status
// checks reject outright for a non-admin pusher (GH006).
func TestGitPluginPrepare_RealClient_NoAssetsConfigSkipsCommit(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client, err := git.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	headBefore, err := client.GetHead()
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the changelog plugin writing CHANGELOG.md without staging it.
	if err := os.WriteFile(filepath.Join(repo.Dir, "CHANGELOG.md"), []byte("# changelog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := NewGitPlugin(client) // no Configure() call: assetsExplicit stays false
	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0"}}
	if err := p.Prepare(ctx, state); err != nil {
		t.Fatalf("Prepare: %v", err)
	}

	headAfter, err := client.GetHead()
	if err != nil {
		t.Fatal(err)
	}
	if headAfter != headBefore {
		t.Fatalf("Prepare created a commit with no assets configured: HEAD moved %s -> %s", headBefore, headAfter)
	}
}
