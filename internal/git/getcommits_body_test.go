// story: e01s04
package git

import (
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func TestGetCommits_PreservesMultiLineBody(t *testing.T) {
	repo := testrepo.Init(t)
	repo.RunGit("commit", "--allow-empty", "-m", "feat: multi", "-m", "line one\nline two")

	commits, err := newTestClient(t).GetCommits("", "HEAD")
	if err != nil {
		t.Fatalf("GetCommits: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(commits))
	}
	if got, want := commits[0].Body, "line one\nline two"; got != want {
		t.Fatalf("Body = %q, want %q", got, want)
	}
	if commits[0].Message != "feat: multi" {
		t.Fatalf("Message = %q", commits[0].Message)
	}
}

func TestGetCommits_MultiLineBodiesAcrossCommits(t *testing.T) {
	repo := testrepo.Init(t)
	repo.RunGit("commit", "--allow-empty", "-m", "feat: first", "-m", "a1\na2")
	repo.RunGit("commit", "--allow-empty", "-m", "fix: second", "-m", "b1\nb2")

	commits, err := newTestClient(t).GetCommits("", "HEAD")
	if err != nil {
		t.Fatalf("GetCommits: %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("got %d commits, want 2", len(commits))
	}
	// git log returns newest first.
	if got, want := commits[0].Body, "b1\nb2"; got != want {
		t.Fatalf("newest Body = %q, want %q", got, want)
	}
	if got, want := commits[1].Body, "a1\na2"; got != want {
		t.Fatalf("oldest Body = %q, want %q", got, want)
	}
	if commits[0].Message != "fix: second" || commits[1].Message != "feat: first" {
		t.Fatalf("messages = %q, %q", commits[0].Message, commits[1].Message)
	}
}

func TestGetCommits_EmptyBodyStillParses(t *testing.T) {
	repo := testrepo.Init(t)
	repo.RunGit("commit", "--allow-empty", "-m", "chore: no body")

	commits, err := newTestClient(t).GetCommits("", "HEAD")
	if err != nil {
		t.Fatalf("GetCommits: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(commits))
	}
	if commits[0].Body != "" {
		t.Fatalf("Body = %q, want empty", commits[0].Body)
	}
}
