package release

import (
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

func assertStepOutput(t *testing.T, path string, want ...string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read GITHUB_OUTPUT: %v", err)
	}
	body := string(data)
	for _, w := range want {
		if !strings.Contains(body, w) {
			t.Errorf("GITHUB_OUTPUT missing %q in %q", w, body)
		}
	}
}

func TestRun_PRSkip_WritesReleasedFalse(t *testing.T) {
	out := t.TempDir() + "/gh-out"
	t.Setenv("GITHUB_OUTPUT", out)
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")

	if err := New(&Context{Logger: zap.NewNop()}).Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	assertStepOutput(t, out, "published=false", "released=false", "reason=pull_request")
}

func TestRun_NoRelease_WritesReleasedFalse(t *testing.T) {
	out := t.TempDir() + "/gh-out"
	t.Setenv("GITHUB_OUTPUT", out)
	t.Setenv("CI", "true")
	clearPREnv(t)

	cfg := &algorithm.Config{
		Branches:   []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:  "v${version}",
		Plugins:    []string{},
		Publishers: map[string]algorithm.PublisherConfig{},
	}
	g := &stubGit{
		branch:  "main",
		commits: []*algorithm.Commit{{Hash: "abcdef1234567890", Message: "docs: clarify README"}},
		repoURL: "https://example.com/repo",
	}
	r := New(&Context{
		Config:     cfg,
		Git:        g,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    plugins.NewRegistry(),
		Publishers: publishers.NewRegistry(),
	})

	if err := r.Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	assertStepOutput(t, out, "published=false", "released=false", "reason=no_release")
}

func TestRun_DryRun_WritesReleasedFalseReasonDryRun(t *testing.T) {
	out := t.TempDir() + "/gh-out"
	t.Setenv("GITHUB_OUTPUT", out)
	clearPREnv(t)

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{},
		Publishers:     map[string]algorithm.PublisherConfig{},
	}
	g := &stubGit{
		branch:  "main",
		commits: []*algorithm.Commit{{Hash: "abcdef1234567890", Message: "feat: add thing"}},
		repoURL: "https://example.com/repo",
	}
	r := New(&Context{
		Config:     cfg,
		Git:        g,
		Logger:     zap.NewNop(),
		DryRun:     true,
		Plugins:    plugins.NewRegistry(),
		Publishers: publishers.NewRegistry(),
	})

	if err := r.Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	assertStepOutput(t, out, "published=false", "released=false", "reason=dry_run")
}
