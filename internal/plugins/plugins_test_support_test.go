// Shared test helpers for plugins tests.
package plugins

import (
	"os"
	"testing"
)

// mockCommandRunner is a mock for testing exec plugin.
type mockCommandRunner struct {
	runFunc func(name string, args ...string) (string, error)
}

func (m *mockCommandRunner) Run(name string, args ...string) (string, error) {
	if m.runFunc != nil {
		return m.runFunc(name, args...)
	}
	return "", nil
}

// chdirTempDir changes to a temp directory and returns a cleanup function.
func chdirTempDir(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	return func() { _ = os.Chdir(oldDir) }
}
