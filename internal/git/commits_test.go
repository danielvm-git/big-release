// story: e08s04
package git

import (
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func TestCommitTraversal_LastTag(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag("v1.0.0")
	repo.Commit("feat: second release")
	repo.Tag("v1.1.0")
	repo.Commit("fix: after tag")

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	last, err := client.GetLastRelease("v${version}")
	if err != nil {
		t.Fatal(err)
	}
	if last == nil {
		t.Fatal("expected last release")
	}
	if last.GitTag != "v1.1.0" {
		t.Fatalf("expected v1.1.0, got %q", last.GitTag)
	}
}

func TestCommitTraversal_Range(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag("v1.0.0")
	repo.Commit("feat: new feature")
	repo.Commit("fix: bugfix")

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	commits, err := client.GetCommits("v1.0.0", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits after tag, got %d", len(commits))
	}
}

func TestCommitTraversal_EmptyRange(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag("v1.0.0")

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	head, err := client.GetHead()
	if err != nil {
		t.Fatal(err)
	}
	_ = head

	commits, err := client.GetCommits("v1.0.0", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 0 {
		t.Fatalf("expected empty range at tag, got %d commits", len(commits))
	}
}
