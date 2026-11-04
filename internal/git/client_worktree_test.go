// story: e01s04
package git

import (
	"os"
	"testing"
	"time"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func writeFile(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestGetModifiedFiles_RenameReportsDestination(t *testing.T) {
	repo := testrepo.Init(t)
	writeFile(t, "old.txt", "content\n")
	repo.RunGit("add", "old.txt")
	repo.RunGit("commit", "-m", "chore: add old")

	repo.RunGit("mv", "old.txt", "new.txt")

	files, err := newTestClient(t).GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles: %v", err)
	}
	if len(files) != 1 || files[0] != "new.txt" {
		t.Fatalf("GetModifiedFiles() = %v, want [new.txt]", files)
	}
}

func TestGetModifiedFiles_StagedAndUntracked(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	files, err := client.GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("clean repo returned %v, want no files", files)
	}

	writeFile(t, "staged.txt", "x\n")
	if err := client.StagePaths([]string{"staged.txt"}); err != nil {
		t.Fatalf("StagePaths: %v", err)
	}
	writeFile(t, "untracked.txt", "y\n")

	files, err = client.GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles: %v", err)
	}
	seen := map[string]bool{}
	for _, f := range files {
		seen[f] = true
	}
	for _, want := range []string{"staged.txt", "untracked.txt"} {
		if !seen[want] {
			t.Fatalf("GetModifiedFiles() = %v, missing %q", files, want)
		}
	}
}

func TestHasChangesToCommit_LooksAtIndexOnly(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	if dirty, err := client.HasChangesToCommit(); err != nil || dirty {
		t.Fatalf("clean repo: dirty=%v err=%v", dirty, err)
	}

	writeFile(t, "untracked.txt", "x\n")
	if dirty, err := client.HasChangesToCommit(); err != nil || dirty {
		t.Fatalf("untracked file must not count as staged: dirty=%v err=%v", dirty, err)
	}

	if err := client.StagePaths([]string{"untracked.txt"}); err != nil {
		t.Fatalf("StagePaths: %v", err)
	}
	if dirty, err := client.HasChangesToCommit(); err != nil || !dirty {
		t.Fatalf("staged file: dirty=%v err=%v", dirty, err)
	}
}

func TestStageChanges_StagesWorkingTree(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	writeFile(t, "all.txt", "x\n")

	client := newTestClient(t)
	if err := client.StageChanges(); err != nil {
		t.Fatalf("StageChanges: %v", err)
	}
	if dirty, err := client.HasChangesToCommit(); err != nil || !dirty {
		t.Fatalf("after StageChanges: dirty=%v err=%v", dirty, err)
	}
	files, err := client.GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles: %v", err)
	}
	if len(files) != 1 || files[0] != "all.txt" {
		t.Fatalf("GetModifiedFiles() = %v, want [all.txt]", files)
	}
}

func TestStagePaths_EmptyAndMissing(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	if err := client.StagePaths(nil); err != nil {
		t.Fatalf("StagePaths(nil) = %v, want nil", err)
	}
	if err := client.StagePaths([]string{"missing.txt"}); err == nil {
		t.Fatal("expected staging a missing path to fail")
	}
}

func TestCommit_CommitsStagedOnly(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	if err := client.Commit("chore: nothing staged"); err == nil {
		t.Fatal("expected commit with nothing staged to fail")
	}

	writeFile(t, "feature.txt", "x\n")
	if err := client.StagePaths([]string{"feature.txt"}); err != nil {
		t.Fatalf("StagePaths: %v", err)
	}
	if err := client.Commit("feat: add feature file"); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if dirty, err := client.HasChangesToCommit(); err != nil || dirty {
		t.Fatalf("after commit: dirty=%v err=%v", dirty, err)
	}
	if got := runGitTest(t, repo.Dir, "log", "-1", "--format=%s"); got != "feat: add feature file\n" {
		t.Fatalf("HEAD subject = %q", got)
	}
	if got := runGitTest(t, repo.Dir, "log", "-1", "--format=%an"); got != "big-release[bot]\n" {
		t.Fatalf("author = %q, want big-release[bot]", got)
	}
}

func TestGetRepositoryURL(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	if _, err := client.GetRepositoryURL(); err == nil {
		t.Fatal("expected an error without a remote")
	}

	repo.SetRemoteOrigin("https://github.com/owner/repo.git")
	got, err := client.GetRepositoryURL()
	if err != nil {
		t.Fatalf("GetRepositoryURL: %v", err)
	}
	if got != "https://github.com/owner/repo.git" {
		t.Fatalf("GetRepositoryURL() = %q", got)
	}
}

func TestGetCurrentBranch_AndDetached(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	got, err := client.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch: %v", err)
	}
	if got != "main" {
		t.Fatalf("GetCurrentBranch() = %q, want main", got)
	}

	repo.RunGit("checkout", "--detach")
	got, err = client.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch detached: %v", err)
	}
	if got != "HEAD" {
		t.Fatalf("GetCurrentBranch() detached = %q, want HEAD", got)
	}
}

func TestIsGitRepo_OutsideRepository(t *testing.T) {
	testrepo.Init(t)
	client := newTestClient(t)
	if !client.IsGitRepo() {
		t.Fatal("IsGitRepo() = false inside a repository")
	}

	t.Chdir(t.TempDir())
	if client.IsGitRepo() {
		t.Fatal("IsGitRepo() = true outside a repository")
	}
}

func TestGetHead_OutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, err := newTestClient(t).GetHead(); err == nil {
		t.Fatal("expected GetHead to fail outside a repository")
	}
}

func TestGetCurrentTime_ParsesWithGitLayout(t *testing.T) {
	got := GetCurrentTime()
	if _, err := time.Parse("2006-01-02 15:04:05 -0700", got); err != nil {
		t.Fatalf("GetCurrentTime() = %q, not in git layout: %v", got, err)
	}
}
