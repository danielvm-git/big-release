package swift_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/swift"
)

// --- Unit tests ---

func TestSwiftName(t *testing.T) {
	t.Run("SC-e02s06-P3-01: Name returns 'swift'", func(t *testing.T) {
		p := swift.NewPublisher()
		if name := p.Name(); name != "swift" {
			t.Errorf("expected Name() == %q, got %q", "swift", name)
		}
	})
}

func TestSwiftDetect(t *testing.T) {
	t.Run("SC-e02s06-P3-02: Detect true with Package.swift", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "Package.swift", "// swift-tools-version:5.5\n")
		p := swift.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when Package.swift exists")
			}
		})
	})

	t.Run("SC-e02s06-P3-03: Detect false without Package.swift", func(t *testing.T) {
		dir := t.TempDir()
		p := swift.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when Package.swift is absent")
			}
		})
	})
}

func TestSwiftPrepare(t *testing.T) {
	t.Run("SC-e02s06-P3-04: Prepare(version) is no-op, returns nil", func(t *testing.T) {
		p := swift.NewPublisher()
		if err := p.Prepare("1.0.0"); err != nil {
			t.Fatalf("expected nil error from Prepare, got %v", err)
		}
	})
}

func TestSwiftAutoRegistration(t *testing.T) {
	t.Run("SC-e02s06-P3-11: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("swift")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "swift" {
			t.Errorf("expected name %q, got %q", "swift", got.Name())
		}
	})
}

// --- Integration tests with exec mocking ---

func TestSwiftPublish(t *testing.T) {
	t.Run("SC-e02s06-P3-05: Publish(version) pushes tag, returns nil on success", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)

		var commands []string
		mockExec := func(name string, args ...string) *exec.Cmd {
			commands = append(commands, name+" "+strings.Join(args, " "))
			switch name {
			case "git":
				return exec.Command("true")
			default:
				return exec.Command("true")
			}
		}

		p := swift.NewPublisher()
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})

		if len(commands) < 2 {
			t.Fatalf("expected at least 2 exec commands (git tag, git push), got %d: %v", len(commands), commands)
		}
		if !strings.Contains(commands[0], "git tag 1.0.0") {
			t.Errorf("expected git tag command, got %q", commands[0])
		}
		if !strings.Contains(commands[1], "git push origin 1.0.0") {
			t.Errorf("expected git push command, got %q", commands[1])
		}
	})

	t.Run("SC-e02s06-P3-06: Publish(version) returns error on git push failure", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)

		pushAttempted := false
		mockExec := func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "tag" {
				return exec.Command("true")
			}
			if name == "git" && len(args) > 0 && args[0] == "push" {
				pushAttempted = true
				return exec.Command("false")
			}
			return exec.Command("true")
		}

		p := swift.NewPublisher()
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error on git push failure, got nil")
			}
			if !strings.Contains(err.Error(), "failed to push") {
				t.Errorf("expected error to mention 'failed to push', got %q", err.Error())
			}
		})

		if !pushAttempted {
			t.Error("expected git push to be attempted")
		}
	})

	t.Run("SC-e02s06-P3-07: Publish(version) in dry-run sets tag locally, no push", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)

		var commands []string
		mockExec := func(name string, args ...string) *exec.Cmd {
			commands = append(commands, name+" "+strings.Join(args, " "))
			return exec.Command("true")
		}

		p := swift.NewPublisher()
		p.DryRun = true
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil in dry-run, got %v", err)
			}
		})

		if len(commands) != 1 {
			t.Fatalf("expected 1 exec command (git tag only) in dry-run, got %d: %v", len(commands), commands)
		}
		if !strings.Contains(commands[0], "git tag 1.0.0") {
			t.Errorf("expected git tag command, got %q", commands[0])
		}
	})

	t.Run("SC-e02s06-P3-08: Publish(version) returns error for invalid version", func(t *testing.T) {
		dir := t.TempDir()
		mockExec := func(name string, args ...string) *exec.Cmd {
			t.Error("unexpected exec call for invalid version")
			return exec.Command("true")
		}

		p := swift.NewPublisher()
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Publish("--delete")
			if err == nil {
				t.Fatal("expected error for invalid version, got nil")
			}
			if !strings.Contains(err.Error(), "invalid version") {
				t.Errorf("expected error to mention 'invalid version', got %q", err.Error())
			}
		})
	})
}

func TestSwiftVerify(t *testing.T) {
	t.Run("SC-e02s06-P3-09: Verify(version) returns nil when tag exists", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)

		mockExec := func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) >= 3 && args[0] == "tag" && args[1] == "-l" {
				// Return a command that outputs the tag name
				cmd := exec.Command("echo", args[2])
				return cmd
			}
			return exec.Command("true")
		}

		p := swift.NewPublisher()
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Verify("1.0.0")
			if err != nil {
				t.Fatalf("expected nil when tag exists, got %v", err)
			}
		})
	})

	t.Run("SC-e02s06-P3-10: Verify(version) returns error when tag missing", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)

		mockExec := func(name string, args ...string) *exec.Cmd {
			// Return a command that outputs nothing (git tag -l returns empty for missing tags)
			cmd := exec.Command("true")
			return cmd
		}

		p := swift.NewPublisher()
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Verify("9.9.9")
			if err == nil {
				t.Fatal("expected error when tag missing, got nil")
			}
			if !strings.Contains(err.Error(), "not found") {
				t.Errorf("expected error to mention 'not found', got %q", err.Error())
			}
		})
	})
}

// --- helpers ---

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func withDir(t *testing.T, dir string, fn func()) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore directory to %s: %v", orig, err)
		}
	}()
	fn()
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v: %s", args, err, string(out))
		}
	}
}
