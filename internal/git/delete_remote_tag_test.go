// story: NT-4
package git

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

// TestDeleteRemoteTag_RemovesTagFromRemote is the NT-4 counterpart to
// DeleteTag: after PushTags succeeds and the release commit push fails, the
// rollback must be able to remove the tag that already landed on the remote.
func TestDeleteRemoteTag_RemovesTagFromRemote(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag("v1.0.0")
	initBareRemote(t, repo)
	runGitTest(t, repo.Dir, "push", "origin", "refs/tags/v1.0.0")

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteRemoteTag("origin", "v1.0.0"); err != nil {
		t.Fatalf("DeleteRemoteTag() error = %v", err)
	}

	if out := runGitTest(t, repo.Dir, "ls-remote", "--tags", "origin"); strings.Contains(out, "v1.0.0") {
		t.Fatalf("remote tag still present after DeleteRemoteTag:\n%s", out)
	}
}

// TestDeleteRemoteTag_AbsentTagIsIdempotent documents the boundary git
// actually gives us: deleting a tag the remote never had exits 0 with a
// warning, so a rollback retried after a partial cleanup still succeeds.
func TestDeleteRemoteTag_AbsentTagIsIdempotent(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	initBareRemote(t, repo)

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteRemoteTag("origin", "v9.9.9"); err != nil {
		t.Fatalf("expected deleting an absent remote tag to succeed, got %v", err)
	}
}

// TestDeleteRemoteTag_SurfacesGitStderr guards the error path: an unreachable
// remote must produce git's reason, not a bare exit status.
func TestDeleteRemoteTag_SurfacesGitStderr(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.SetRemoteOrigin(filepath.Join(repo.Dir, "does-not-exist"))

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	delErr := client.DeleteRemoteTag("origin", "v1.0.0")
	if delErr == nil {
		t.Fatal("expected deleting a remote tag from a nonexistent remote to fail")
	}
	if strings.HasSuffix(delErr.Error(), "exit status 1") {
		t.Fatalf("DeleteRemoteTag() error dropped git's stderr: %q", delErr.Error())
	}
	if !strings.Contains(delErr.Error(), "does-not-exist") {
		t.Fatalf("expected git's stderr detail in error, got %q", delErr.Error())
	}
}

func initBareRemote(t *testing.T, repo *testrepo.Repo) string {
	t.Helper()
	remoteDir := filepath.Join(t.TempDir(), "remote.git")
	runGitTest(t, repo.Dir, "init", "--bare", remoteDir)
	repo.SetRemoteOrigin(remoteDir)
	return remoteDir
}

// runGitTest runs a git command in dir with the same environment scrubbing as
// the production client, so an ambient git hook cannot redirect the test.
func runGitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = scrubGitEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return string(out)
}
