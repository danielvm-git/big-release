// story: BUG-git-push-error-swallowed
package git

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

// TestPush_SurfacesGitStderr guards against Push() collapsing every failure
// into a bare "exit status 1". Push runs git with a nil Stderr passed
// straight to cmd.Run(), which Go connects to /dev/null — the actual
// reason a push was rejected must still reach the caller.
func TestPush_SurfacesGitStderr(t *testing.T) {
	// A branch with no upstream configured is the simplest reliable way to
	// make `git push` (no refspec) fail with real, stable stderr text.
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.SetRemoteOrigin(filepath.Join(repo.Dir, "does-not-exist"))

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	pushErr := client.Push("origin")
	if pushErr == nil {
		t.Fatal("expected push with no upstream configured to fail")
	}
	if strings.HasSuffix(pushErr.Error(), "exit status 1") || strings.HasSuffix(pushErr.Error(), "exit status 128") {
		t.Fatalf("Push() error dropped git's stderr, got bare exit status: %q", pushErr.Error())
	}
	if !strings.Contains(pushErr.Error(), "upstream") {
		t.Fatalf("expected git's stderr detail in error, got %q", pushErr.Error())
	}
}

// TestPushTags_SurfacesGitStderr mirrors TestPush_SurfacesGitStderr for the
// tag-push path used by plugins.GitPlugin.Publish.
func TestPushTags_SurfacesGitStderr(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag("v1.0.0")
	repo.SetRemoteOrigin(filepath.Join(repo.Dir, "does-not-exist"))

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	pushErr := client.PushTags("origin")
	if pushErr == nil {
		t.Fatal("expected tag push to a nonexistent remote to fail")
	}
	if pushErr.Error() == "failed to push tags: exit status 1" {
		t.Fatalf("PushTags() error dropped git's stderr, got bare exit status: %q", pushErr.Error())
	}
	if !strings.Contains(pushErr.Error(), "does-not-exist") {
		t.Fatalf("expected git's stderr detail in error, got %q", pushErr.Error())
	}
}
