// story: e24s03

package nodeutil_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers/nodeutil"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func withDir(t *testing.T, dir string, fn func()) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(orig)
	}()
	fn()
}

func TestIsValidPackageName(t *testing.T) {
	valid := []string{"my-package", "@scope/my-package", "my.pkg-name_1"}
	for _, name := range valid {
		if !nodeutil.IsValidPackageName(name) {
			t.Errorf("expected valid: %q", name)
		}
	}
	invalid := []string{"", ".hidden", "MyPackage", "@scope", "@scope/", "my package"}
	for _, name := range invalid {
		if nodeutil.IsValidPackageName(name) {
			t.Errorf("expected invalid: %q", name)
		}
	}
}

func TestReadWritePackageJSON(t *testing.T) {
	t.Run("SC-e24s03-U01: round-trip version update", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"x","version":"1.0.0"}`)
		withDir(t, dir, func() {
			pkg, err := nodeutil.ReadPackageJSON("npm")
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			pkg["version"] = "2.0.0"
			if err := nodeutil.WritePackageJSON("npm", pkg); err != nil {
				t.Fatalf("write: %v", err)
			}
			data, _ := os.ReadFile("package.json")
			if !strings.Contains(string(data), `"version": "2.0.0"`) {
				t.Errorf("got %s", data)
			}
		})
	})

	t.Run("SC-e24s03-U02: read missing file", func(t *testing.T) {
		withDir(t, t.TempDir(), func() {
			_, err := nodeutil.ReadPackageJSON("pnpm")
			if err == nil || !strings.Contains(err.Error(), "pnpm: failed to read package.json") {
				t.Fatalf("got %v", err)
			}
		})
	})
}

func TestReadPackageName(t *testing.T) {
	t.Run("SC-e24s03-U03: valid name", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name":"test-pkg","version":"1.0.0"}`)
		withDir(t, dir, func() {
			name, err := nodeutil.ReadPackageName("npm")
			if err != nil || name != "test-pkg" {
				t.Fatalf("got %q %v", name, err)
			}
		})
	})

	t.Run("SC-e24s03-U04: missing name", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"version":"1.0.0"}`)
		withDir(t, dir, func() {
			_, err := nodeutil.ReadPackageName("npm")
			if err == nil || !strings.Contains(err.Error(), "package name not found") {
				t.Fatalf("got %v", err)
			}
		})
	})
}
