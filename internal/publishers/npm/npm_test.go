// story: e02s07

package npm_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	npm "github.com/danielvm-git/big-release/internal/publishers/npm"
)

func TestNpmName(t *testing.T) {
	t.Run("SC-e02-npm-U01: Name returns 'npm'", func(t *testing.T) {
		p := npm.NewPublisher()
		if name := p.Name(); name != "npm" {
			t.Errorf("expected Name() == %q, got %q", "npm", name)
		}
	})
}

func TestNpmDetect(t *testing.T) {
	t.Run("SC-e02-npm-U02: Detect true with package.json", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name": "test-pkg", "version": "1.0.0"}`)
		withDir(t, dir, func() {
			if !npm.NewPublisher().Detect() {
				t.Errorf("expected Detect() == true")
			}
		})
	})

	t.Run("SC-e02-npm-U03: Detect false without package.json", func(t *testing.T) {
		withDir(t, t.TempDir(), func() {
			if npm.NewPublisher().Detect() {
				t.Errorf("expected Detect() == false")
			}
		})
	})
}

func TestNpmPrepare(t *testing.T) {
	t.Run("SC-e02-npm-U04: Prepare updates version", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name": "test-pkg", "version": "1.0.0"}`)
		withDir(t, dir, func() {
			if err := npm.NewPublisher().Prepare("2.0.0"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			data, _ := os.ReadFile("package.json")
			if !strings.Contains(string(data), `"version": "2.0.0"`) {
				t.Errorf("expected version 2.0.0, got %s", string(data))
			}
		})
	})

	t.Run("SC-e02-npm-U05: Prepare error on missing package.json", func(t *testing.T) {
		err := withDirWrap(t.TempDir(), func() error { return npm.NewPublisher().Prepare("2.0.0") })
		if err == nil || !strings.Contains(err.Error(), "failed to read package.json") {
			t.Fatalf("expected error about reading, got %v", err)
		}
	})

	t.Run("SC-e02-npm-U06: Prepare error on malformed JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{invalid}`)
		err := withDirWrap(dir, func() error { return npm.NewPublisher().Prepare("2.0.0") })
		if err == nil || !strings.Contains(err.Error(), "failed to parse package.json") {
			t.Fatalf("expected error about parsing, got %v", err)
		}
	})
}

func TestNpmAutoRegistration(t *testing.T) {
	t.Run("SC-e02-npm-U07: Auto-registered via init()", func(t *testing.T) {
		got, err := publishers.Get("npm")
		if err != nil {
			t.Fatalf("expected publisher registered, got error: %v", err)
		}
		if got.Name() != "npm" {
			t.Errorf("expected name %q, got %q", "npm", got.Name())
		}
	})
}

func TestIsValidPackageName(t *testing.T) {
	valid := []struct{ name, value string }{
		{"SC-e02-npm-V01: valid unscoped", "my-package"},
		{"SC-e02-npm-V02: valid scoped", "@scope/my-package"},
		{"SC-e02-npm-V03: valid dots hyphens underscores", "my.pkg-name_1"},
	}
	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "package.json", `{"name":"`+tc.value+`","version":"1.0.0"}`)
			p := npm.NewPublisher()
			p.DryRun = true
			withDir(t, dir, func() {
				if err := p.Verify("1.0.0"); err != nil {
					t.Errorf("expected valid name %q to pass, got %v", tc.value, err)
				}
			})
		})
	}

	invalid := []struct{ name, value string }{
		{"SC-e02-npm-V04: empty name", ""},
		{"SC-e02-npm-V05: starts with dot", ".hidden-pkg"},
		{"SC-e02-npm-V06: uppercase", "MyPackage"},
		{"SC-e02-npm-V07: scoped missing slash", "@scope"},
		{"SC-e02-npm-V08: scoped empty package", "@scope/"},
		{"SC-e02-npm-V09: spaces", "my package"},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "package.json", `{"name":"`+tc.value+`","version":"1.0.0"}`)
			p := npm.NewPublisher()
			p.DryRun = true
			withDir(t, dir, func() {
				if err := p.Verify("1.0.0"); err == nil {
					t.Errorf("expected error for invalid name %q, got nil", tc.value)
				}
			})
		})
	}
}

func TestNpmPublish(t *testing.T) {
	t.Run("SC-e02-npm-I01: Publish success", func(t *testing.T) {
		called := false
		p := npm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd {
			called = true
			if name != "npm" || strings.Join(args, " ") != "publish" {
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

	t.Run("SC-e02-npm-I02: Publish failure", func(t *testing.T) {
		p := npm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("false") }
		if err := p.Publish("1.0.0"); err == nil || !strings.Contains(err.Error(), "npm: publish failed") {
			t.Fatalf("expected 'publish failed' error, got %v", err)
		}
	})

	t.Run("SC-e02-npm-I03: Publish dry-run skips exec", func(t *testing.T) {
		called := false
		p := npm.NewPublisher()
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
}

func TestNpmVerify(t *testing.T) {
	t.Run("SC-e02-npm-I04: Verify version match", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := npm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("echo", "1.0.0") }
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	})

	t.Run("SC-e02-npm-I05: Verify version mismatch", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := npm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("echo", "2.0.0") }
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err == nil || !strings.Contains(err.Error(), "does not match") {
				t.Fatalf("expected 'does not match' error, got %v", err)
			}
		})
	})

	t.Run("SC-e02-npm-I06: Verify npm view failure", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		p := npm.NewPublisher()
		p.ExecCommand = func(name string, args ...string) *exec.Cmd { return exec.Command("false") }
		withDir(t, dir, func() {
			if err := p.Verify("1.0.0"); err == nil || !strings.Contains(err.Error(), "npm: failed to verify") {
				t.Fatalf("expected 'failed to verify' error, got %v", err)
			}
		})
	})

	t.Run("SC-e02-npm-I07: Verify missing package name", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"version":"1.0.0"}`)
		called := false
		p := npm.NewPublisher()
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

	t.Run("SC-e02-npm-I08: Verify dry-run skips exec", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		called := false
		p := npm.NewPublisher()
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

// withDirWrap runs fn in dir and returns any error.
func withDirWrap(dir string, fn func() error) error {
	orig, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}
	defer os.Chdir(orig) //nolint:errcheck
	return fn()
}
