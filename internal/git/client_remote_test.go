// story: e01s04 e08s04
package git

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func TestCreateTag_AndDeleteTag(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	if err := client.CreateTag("v1.0.0", "release 1.0.0"); err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	tags, err := client.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(tags) != 1 || tags[0] != "v1.0.0" {
		t.Fatalf("GetTags() = %v, want [v1.0.0]", tags)
	}
	if err := client.CreateTag("v1.0.0", "duplicate"); err == nil {
		t.Fatal("expected duplicate tag creation to fail")
	}
	if err := client.DeleteTag("v1.0.0"); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
	if err := client.DeleteTag("v9.9.9"); err == nil {
		t.Fatal("expected deleting a missing tag to fail")
	}
}

func TestAddNote_AndPushNotes(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	initBareRemote(t, repo)

	client := newTestClient(t)
	if err := client.AddNote("release note", "HEAD"); err != nil {
		t.Fatalf("AddNote: %v", err)
	}
	if err := client.PushNotes("origin", "HEAD"); err != nil {
		t.Fatalf("PushNotes: %v", err)
	}

	out := runGitTest(t, repo.Dir, "ls-remote", "origin", "refs/notes/big-release-HEAD")
	if !strings.Contains(out, "refs/notes/big-release-HEAD") {
		t.Fatalf("notes ref missing on remote:\n%s", out)
	}
}

func TestAddNote_UnknownRefFails(t *testing.T) {
	testrepo.Init(t)
	if err := newTestClient(t).AddNote("note", "does-not-exist"); err == nil {
		t.Fatal("expected AddNote against an unknown ref to fail")
	}
}

func TestPushNotes_BadRemoteFails(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	client := newTestClient(t)
	if err := client.AddNote("note", "HEAD"); err != nil {
		t.Fatalf("AddNote: %v", err)
	}
	repo.SetRemoteOrigin(filepath.Join(repo.Dir, "does-not-exist"))
	if err := client.PushNotes("origin", "HEAD"); err == nil {
		t.Fatal("expected PushNotes to a missing remote to fail")
	}
}

func TestVerifyAuth(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	initBareRemote(t, repo)
	runGitTest(t, repo.Dir, "push", "origin", "main")

	client := newTestClient(t)
	if err := client.VerifyAuth("origin", "main"); err != nil {
		t.Fatalf("VerifyAuth: %v", err)
	}

	repo.RunGit("remote", "set-url", "origin", filepath.Join(repo.Dir, "does-not-exist"))
	if err := client.VerifyAuth("origin", "main"); err == nil {
		t.Fatal("expected VerifyAuth against a missing remote to fail")
	}
}

func TestIsBranchUpToDate(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	initBareRemote(t, repo)
	runGitTest(t, repo.Dir, "push", "origin", "main")

	client := newTestClient(t)
	upToDate, err := client.IsBranchUpToDate("origin", "main")
	if err != nil {
		t.Fatalf("IsBranchUpToDate: %v", err)
	}
	if !upToDate {
		t.Fatal("IsBranchUpToDate() = false after pushing main")
	}

	repo.Commit("feat: ahead of remote")
	upToDate, err = client.IsBranchUpToDate("origin", "main")
	if err != nil {
		t.Fatalf("IsBranchUpToDate: %v", err)
	}
	if upToDate {
		t.Fatal("IsBranchUpToDate() = true with a local commit ahead of remote")
	}

	upToDate, err = client.IsBranchUpToDate("origin", "missing-branch")
	if err != nil {
		t.Fatalf("IsBranchUpToDate missing branch: %v", err)
	}
	if upToDate {
		t.Fatal("IsBranchUpToDate() = true for a branch absent from the remote")
	}
}

func TestIsBranchUpToDate_BadRemoteFails(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.SetRemoteOrigin(filepath.Join(repo.Dir, "does-not-exist"))

	if _, err := newTestClient(t).IsBranchUpToDate("origin", "main"); err == nil {
		t.Fatal("expected IsBranchUpToDate against a missing remote to fail")
	}
}

func TestPushError_Unwrap(t *testing.T) {
	sentinel := errors.New("boom")
	pe := &PushError{Outcome: PushOutcomeUnknown, Raw: "boom", Err: sentinel}

	if !errors.Is(pe, sentinel) {
		t.Fatal("errors.Is() = false, want true")
	}
	if errors.Unwrap(pe) != sentinel {
		t.Fatal("errors.Unwrap() did not return the wrapped error")
	}
}

func TestScrubGitEnv_RemovesHookContext(t *testing.T) {
	t.Setenv("GIT_DIR", "/tmp/leaked")
	t.Setenv("GIT_WORK_TREE", "/tmp/leaked")
	t.Setenv("GIT_INDEX_FILE", "/tmp/leaked")

	env := scrubGitEnv()
	for _, e := range env {
		for _, prefix := range []string{"GIT_DIR=", "GIT_WORK_TREE=", "GIT_INDEX_FILE="} {
			if strings.HasPrefix(e, prefix) {
				t.Fatalf("scrubGitEnv() leaked %q", e)
			}
		}
	}

	foundPath := false
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			foundPath = true
		}
	}
	if !foundPath {
		t.Fatal("scrubGitEnv() dropped PATH")
	}
}
