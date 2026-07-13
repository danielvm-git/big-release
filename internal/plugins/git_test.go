// story: e03s01
package plugins

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// setupTestGitRepo creates a temporary git repository for testing.
func setupTestGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "commit.gpgSign", "false"},
		{"git", "config", "tag.gpgSign", "false"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v\noutput: %s", args, err, string(out))
		}
	}

	return dir
}

// writeFileInDir writes a file in the given directory for test setup.
func writeFileInDir(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(dir+"/"+name, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

// execInDir runs a command in the given directory.
func execInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %v failed: %v\noutput: %s", args, err, string(out))
	}
	return string(out)
}

func TestGitPluginName(t *testing.T) {
	t.Run("SC-e03s01-P1-01: Name returns 'git'", func(t *testing.T) {
		p := NewGitPlugin()
		if name := p.Name(); name != "git" {
			t.Errorf("expected Name() == %q, got %q", "git", name)
		}
	})
}

func TestGitPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e03s01-P1-02: VerifyConditions passes in a git repo", func(t *testing.T) {
		dir := setupTestGitRepo(t)
		p := NewGitPlugin()
		p.Dir = dir
		ctx := &algorithm.Context{}

		if err := p.VerifyConditions(ctx); err != nil {
			t.Errorf("expected no error in git repo, got: %v", err)
		}
	})

	t.Run("SC-e03s01-P1-03: VerifyConditions fails outside git repo", func(t *testing.T) {
		dir := t.TempDir()
		p := NewGitPlugin()
		p.Dir = dir
		ctx := &algorithm.Context{}

		if err := p.VerifyConditions(ctx); err == nil {
			t.Error("expected error outside git repo, got nil")
		}
	})
}

func TestGitPluginAnalyzeCommits(t *testing.T) {
	t.Run("SC-e03s01-P1-04: AnalyzeCommits returns empty release type", func(t *testing.T) {
		p := NewGitPlugin()
		ctx := &algorithm.Context{}
		rt, err := p.AnalyzeCommits(ctx)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if rt != "" {
			t.Errorf("expected empty release type, got %q", rt)
		}
	})
}

func TestGitPluginGenerateNotes(t *testing.T) {
	t.Run("SC-e03s01-P1-05: GenerateNotes returns empty string", func(t *testing.T) {
		p := NewGitPlugin()
		ctx := &algorithm.Context{}
		notes, err := p.GenerateNotes(ctx)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if notes != "" {
			t.Errorf("expected empty notes, got %q", notes)
		}
	})
}

func TestGitPluginPrepare(t *testing.T) {
	t.Run("SC-e03s01-P1-06: Prepare does nothing in dry-run mode", func(t *testing.T) {
		p := NewGitPlugin()
		ctx := &algorithm.Context{DryRun: true}
		if err := p.Prepare(ctx); err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
	})

	t.Run("SC-e03s01-P1-07: Prepare stages and commits changes", func(t *testing.T) {
		dir := setupTestGitRepo(t)

		// Create an initial commit so we have a baseline
		writeFileInDir(t, dir, "initial.txt", "initial")
		execInDir(t, dir, "git", "add", ".")
		execInDir(t, dir, "git", "commit", "-m", "initial commit")

		// Now create a change to commit
		writeFileInDir(t, dir, "test.txt", "release content")

		p := NewGitPlugin()
		p.Dir = dir
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      false,
		}

		if err := p.Prepare(ctx); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Verify the commit was made
		cmd := exec.Command("git", "log", "--oneline", "-1")
		cmd.Dir = dir
		out, _ := cmd.Output()
		if !strings.Contains(string(out), "1.0.0") {
			t.Errorf("expected commit message to contain version, got: %s", string(out))
		}
	})

	t.Run("SC-e03s01-P1-08: Prepare skips commit when no changes", func(t *testing.T) {
		dir := setupTestGitRepo(t)

		// Create initial commit
		writeFileInDir(t, dir, "initial.txt", "initial")
		execInDir(t, dir, "git", "add", ".")
		execInDir(t, dir, "git", "commit", "-m", "initial commit")

		p := NewGitPlugin()
		p.Dir = dir
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      false,
		}

		if err := p.Prepare(ctx); err != nil {
			t.Errorf("expected no error with no changes, got: %v", err)
		}
	})
}

func TestGitPluginPublish(t *testing.T) {
	t.Run("SC-e03s01-P1-09: Publish does nothing in dry-run mode", func(t *testing.T) {
		p := NewGitPlugin()
		ctx := &algorithm.Context{DryRun: true}
		release, err := p.Publish(ctx)
		if err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
		if release != nil {
			t.Errorf("expected nil release in dry-run, got %v", release)
		}
	})

	t.Run("SC-e03s01-P1-10: Publish creates git tag", func(t *testing.T) {
		dir := setupTestGitRepo(t)

		// Create initial commit
		writeFileInDir(t, dir, "initial.txt", "initial")
		execInDir(t, dir, "git", "add", ".")
		execInDir(t, dir, "git", "commit", "-m", "initial commit")

		p := NewGitPlugin()
		p.Dir = dir
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      false,
		}

		if err := p.Prepare(ctx); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		// Publish (tag creation, push will fail since no remote, but tag should work)
		_, err := p.Publish(ctx)

		// Check tag exists
		tagCmd := exec.Command("git", "tag", "-l", "1.0.0")
		tagCmd.Dir = dir
		tagOut, _ := tagCmd.Output()
		if strings.TrimSpace(string(tagOut)) != "1.0.0" {
			t.Errorf("expected tag 1.0.0 to exist, got: %s", string(tagOut))
		}
		_ = err
	})
}

func TestGitPluginSuccess(t *testing.T) {
	t.Run("SC-e03s01-P1-11: Success returns nil", func(t *testing.T) {
		p := NewGitPlugin()
		if err := p.Success(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitPluginFail(t *testing.T) {
	t.Run("SC-e03s01-P1-12: Fail returns nil", func(t *testing.T) {
		p := NewGitPlugin()
		if err := p.Fail(&algorithm.Context{}, nil); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitPluginAutoRegistration(t *testing.T) {
	t.Run("SC-e03s01-P1-13: GitPlugin auto-registered in global registry", func(t *testing.T) {
		found := false
		for _, name := range List() {
			if name == "git" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'git' to be registered in global registry")
		}
	})
}
