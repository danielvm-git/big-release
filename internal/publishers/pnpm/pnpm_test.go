// story: e24s01 e24s02

package pnpm_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
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

	t.Run("SC-e24s02-U01: Detect false in empty dir", func(t *testing.T) {
		withDir(t, t.TempDir(), func() {
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

	t.Run("SC-e24s02-U02: Prepare error on malformed JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{invalid}`)
		err := withDirWrap(dir, func() error { return pnpm.NewPublisher().Prepare("2.0.0") })
		if err == nil || !strings.Contains(err.Error(), "failed to parse package.json") {
			t.Fatalf("expected error about parsing, got %v", err)
		}
	})
}

func TestPnpmAutoRegistration(t *testing.T) {
	t.Run("SC-e24s02-U03: Auto-registered via init()", func(t *testing.T) {
		got, err := publishers.Get("pnpm")
		if err != nil {
			t.Fatalf("expected publisher registered, got error: %v", err)
		}
		if got.Name() != "pnpm" {
			t.Errorf("expected name %q, got %q", "pnpm", got.Name())
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

	t.Run("SC-e24s02-I01: Publish failure", func(t *testing.T) {
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("false") }
		if err := p.Publish("1.0.0"); err == nil || !strings.Contains(err.Error(), "pnpm: publish failed") {
			t.Fatalf("expected 'publish failed' error, got %v", err)
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

	t.Run("SC-e24s02-I02: latest channel omits dist-tag flag", func(t *testing.T) {
		var args []string
		p := pnpm.NewPublisher()
		p.SetChannel("latest")
		p.ExecCommand = func(name string, a ...string) *exec.Cmd {
			args = a
			return exec.Command("true")
		}
		if err := p.Publish("1.0.0"); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if strings.Join(args, " ") != "publish --no-git-checks" {
			t.Errorf("args = %v, want publish --no-git-checks", args)
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

	t.Run("SC-e24s02-I03: Verify version mismatch", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("echo", "2.0.0") }
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err == nil || !strings.Contains(err.Error(), "does not match") {
				t.Fatalf("expected 'does not match' error, got %v", err)
			}
		})
	})

	t.Run("SC-e24s02-I04: Verify pnpm view failure", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("false") }
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err == nil || !strings.Contains(err.Error(), "pnpm: failed to verify") {
				t.Fatalf("expected 'failed to verify' error, got %v", err)
			}
		})
	})

	t.Run("SC-e24s02-I05: Verify missing package name", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"version":"1.0.0"}`)
		called := false
		p := pnpm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			called = true
			return exec.Command("true")
		}
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err == nil || !strings.Contains(err.Error(), "package name not found") {
				t.Fatalf("expected 'name not found' error, got %v", err)
			}
		})
		if called {
			t.Error("expected no exec call when name missing")
		}
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
