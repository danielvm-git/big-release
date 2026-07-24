// story: e24s01

package pnpm_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	pnpm "github.com/danielvm-git/big-release/internal/publishers/pnpm"
)

func TestPnpmName(t *testing.T) {
	t.Run("SC-e24s01-U01: Name returns pnpm", func(t *testing.T) {
		p := pnpm.NewPublisher()
		if name := p.Name(); name != "pnpm" {
			t.Errorf("expected Name() == %q, got %q", "pnpm", name)
		}
	})
}

func TestPnpmDetect(t *testing.T) {
	t.Run("SC-e24s01-U02: Detect true with pnpm-lock.yaml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
		withDir(t, dir, func() {
			if !pnpm.NewPublisher().Detect() {
				t.Errorf("expected Detect() == true")
			}
		})
	})

	t.Run("SC-e24s01-U03: Detect true with pnpm-workspace.yaml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pnpm-workspace.yaml", "packages:\n  - 'packages/*'\n")
		withDir(t, dir, func() {
			if !pnpm.NewPublisher().Detect() {
				t.Errorf("expected Detect() == true")
			}
		})
	})

	t.Run("SC-e24s01-U04: Detect false without pnpm markers", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"x","version":"1.0.0"}`)
		withDir(t, dir, func() {
			if pnpm.NewPublisher().Detect() {
				t.Errorf("expected Detect() == false")
			}
		})
	})
}

func TestPnpmPrepare(t *testing.T) {
	t.Run("SC-e24s01-U05: Prepare updates version", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name": "test-pkg", "version": "1.0.0"}`)
		withDir(t, dir, func() {
			if err := pnpm.NewPublisher().Prepare("2.0.0"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			data, _ := os.ReadFile("package.json")
			if !strings.Contains(string(data), `"version": "2.0.0"`) {
				t.Errorf("expected version 2.0.0, got %s", string(data))
			}
		})
	})

	t.Run("SC-e24s01-U06: Prepare error on missing package.json", func(t *testing.T) {
		err := withDirWrap(t.TempDir(), func() error { return pnpm.NewPublisher().Prepare("2.0.0") })
		if err == nil || !strings.Contains(err.Error(), "failed to read package.json") {
			t.Fatalf("expected error about reading, got %v", err)
		}
	})
}

func TestPnpmPublish(t *testing.T) {
	t.Run("SC-e24s01-I01: Publish success with --no-git-checks", func(t *testing.T) {
		called := false
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			called = true
			if name != "pnpm" || strings.Join(args, " ") != "publish --no-git-checks" {
				t.Errorf("unexpected: %s %v", name, args)
			}
			return exec.Command("true")
		}
		if err := p.Publish("1.0.0"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected exec call")
		}
	})

	t.Run("SC-e24s01-I02: Publish dry-run skips exec", func(t *testing.T) {
		called := false
		p := pnpm.NewPublisher()
		p.DryRun = true
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			called = true
			return exec.Command("true")
		}
		if err := p.Publish("1.0.0"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Error("expected no exec call in dry-run")
		}
	})

	t.Run("SC-e24s01-I03: Publish uses --tag from channel", func(t *testing.T) {
		var args []string
		p := pnpm.NewPublisher()
		p.SetChannel("next")
		p.ExecCommand = func(name string, a ...string) *exec.Cmd {
			args = a
			return exec.Command("true")
		}
		if err := p.Publish("1.0.0"); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if strings.Join(args, " ") != "publish --no-git-checks --tag next" {
			t.Errorf("args = %v, want publish --no-git-checks --tag next", args)
		}
	})
}

func TestPnpmVerify(t *testing.T) {
	t.Run("SC-e24s01-I04: Verify version match", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			if name != "pnpm" || strings.Join(args, " ") != "view test-pkg version" {
				t.Errorf("unexpected: %s %v", name, args)
			}
			return exec.Command("echo", "1.0.0")
		}
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	})

	t.Run("SC-e24s01-I05: Verify dry-run skips exec", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		called := false
		p := pnpm.NewPublisher()
		p.DryRun = true
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			called = true
			return exec.Command("echo", "1.0.0")
		}
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
		if called {
			t.Error("expected no exec call in dry-run")
		}
	})
}
