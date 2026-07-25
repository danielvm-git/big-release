// story: e08s03 e08s04
package testrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Repo is a temporary git repository for integration tests.
type Repo struct {
	Dir string
	t   *testing.T
}

// Init creates a new temp git repo and chdirs into it for the test duration.
func Init(t *testing.T) *Repo {
	t.Helper()
	dir := t.TempDir()
	r := &Repo{Dir: dir, t: t}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
	})

	runGit(t, "init")
	runGit(t, "config", "user.email", "test@example.com")
	runGit(t, "config", "user.name", "Test User")
	runGit(t, "checkout", "-b", "main")
	return r
}

func runGit(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = "."
	// Strip git worktree/hook context vars so this always operates on the
	// tempdir at cmd.Dir, never an ambient repo leaked from a calling git
	// hook (e.g. running `go test` from inside a pre-commit hook).
	cmd.Env = scrubGitEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func scrubGitEnv() []string {
	env := os.Environ()
	out := make([]string, 0, len(env))
	for _, e := range env {
		switch {
		case strings.HasPrefix(e, "GIT_DIR="),
			strings.HasPrefix(e, "GIT_WORK_TREE="),
			strings.HasPrefix(e, "GIT_INDEX_FILE="),
			strings.HasPrefix(e, "GIT_PREFIX="),
			strings.HasPrefix(e, "GIT_COMMON_DIR="),
			strings.HasPrefix(e, "GIT_OBJECT_DIRECTORY="),
			strings.HasPrefix(e, "GIT_ALTERNATE_OBJECT_DIRECTORIES="):
			continue
		default:
			out = append(out, e)
		}
	}
	return out
}

// Commit writes a file and creates a commit with the given message.
func (r *Repo) Commit(message string) {
	r.t.Helper()
	path := filepath.Join(r.Dir, "file.txt")
	prev := ""
	if b, err := os.ReadFile(path); err == nil {
		prev = string(b)
	}
	content := prev + message + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
	runGit(r.t, "add", "-A")
	runGit(r.t, "commit", "-m", message)
}

// Tag creates an annotated tag at HEAD.
func (r *Repo) Tag(name string) {
	r.t.Helper()
	runGit(r.t, "tag", "-a", name, "-m", name)
}

// SetRemoteOrigin sets remote.origin.url without auth injection side effects.
func (r *Repo) SetRemoteOrigin(url string) {
	r.t.Helper()
	runGit(r.t, "remote", "add", "origin", url)
}

// ScrubEnv clears sensitive env vars for the test process.
func ScrubEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{"GH_TOKEN", "GITHUB_TOKEN", "NPM_TOKEN"} {
		t.Setenv(name, "")
	}
}
