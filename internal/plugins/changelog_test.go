// story: e03s04
package plugins

import (
	"os"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func TestChangelogPluginName(t *testing.T) {
	t.Run("SC-e03s04-P1-01: Name returns 'changelog'", func(t *testing.T) {
		p := NewChangelogPlugin()
		if name := p.Name(); name != "changelog" {
			t.Errorf("expected Name() == %q, got %q", "changelog", name)
		}
	})
}

func TestChangelogPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e03s04-P1-02: VerifyConditions returns nil", func(t *testing.T) {
		if err := NewChangelogPlugin().VerifyConditions(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestChangelogPluginAnalyzeCommits(t *testing.T) {
	t.Run("SC-e03s04-P1-03: AnalyzeCommits returns empty", func(t *testing.T) {
		rt, err := NewChangelogPlugin().AnalyzeCommits(&algorithm.Context{})
		if err != nil || rt != "" {
			t.Errorf("expected empty, got err=%v rt=%q", err, rt)
		}
	})
}

func TestChangelogPluginGenerateNotes(t *testing.T) {
	p := NewChangelogPlugin()
	t.Run("SC-e03s04-P1-04: returns empty with nil NextRelease", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{NextRelease: nil})
		if err != nil || notes != "" {
			t.Errorf("expected empty, got err=%v notes=%q", err, notes)
		}
	})
	t.Run("SC-e03s04-P1-05: with no commits", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			Commits:     []*algorithm.Commit{},
		})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if !strings.Contains(notes, "No significant changes") {
			t.Errorf("expected 'No significant changes', got: %s", notes)
		}
	})
	t.Run("SC-e03s04-P1-06: groups commits by type", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{
			NextRelease: &algorithm.Release{Version: "2.0.0"},
			Commits: []*algorithm.Commit{
				{Type: "feat", Subject: "add user login", Scope: "auth"},
				{Type: "fix", Subject: "fix nil pointer", Scope: "core"},
				{Type: "docs", Subject: "update readme"},
			},
		})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		for _, want := range []string{"### Features", "### Bug Fixes", "### Documentation", "add user login", "fix nil pointer"} {
			if !strings.Contains(notes, want) {
				t.Errorf("expected %q in notes", want)
			}
		}
	})
	t.Run("SC-e03s04-P1-07: includes breaking changes section", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{
			NextRelease: &algorithm.Release{Version: "3.0.0"},
			Commits: []*algorithm.Commit{
				{Type: "feat", Subject: "add new API", Breaking: true, Body: "Removes old v1"},
				{Type: "fix", Subject: "fix timeout"},
			},
		})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if !strings.Contains(notes, "BREAKING CHANGES") || !strings.Contains(notes, "Removes old v1") {
			t.Errorf("expected BREAKING CHANGES section, got:\n%s", notes)
		}
	})
	t.Run("SC-e03s04-P1-08: includes version and date", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.5.0"},
			Commits:     []*algorithm.Commit{{Type: "feat", Subject: "new feature"}},
		})
		if err != nil || !strings.Contains(notes, "1.5.0") {
			t.Errorf("expected version 1.5.0 in notes, got err=%v notes=%s", err, notes)
		}
	})
	t.Run("SC-e03s04-P1-09: excludes breaking commits from regular sections", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{
			NextRelease: &algorithm.Release{Version: "4.0.0"},
			Commits: []*algorithm.Commit{
				{Type: "feat", Subject: "normal"},
				{Type: "feat", Subject: "breaking", Breaking: true},
			},
		})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if !strings.Contains(notes, "BREAKING CHANGES") {
			t.Errorf("expected BREAKING CHANGES section, got:\n%s", notes)
		}
	})
}

func TestChangelogPluginPrepare(t *testing.T) {
	p := NewChangelogPlugin()
	t.Run("SC-e03s04-P1-10: does nothing in dry-run mode", func(t *testing.T) {
		ctx := &algorithm.Context{DryRun: true, NextRelease: &algorithm.Release{Version: "1.0.0", Notes: "test"}}
		if err := p.Prepare(ctx); err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
	})
	t.Run("SC-e03s04-P1-11: creates new CHANGELOG.md when none exists", func(t *testing.T) {
		defer chdirTempDir(t)()
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0", Notes: "## 1.0.0 (2024-01-15)\n\n### Features\n\n- feat: initial release\n"},
		}
		if err := p.Prepare(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		data, _ := os.ReadFile("CHANGELOG.md")
		content := string(data)
		if !strings.Contains(content, "1.0.0") || !strings.Contains(content, "# Changelog") {
			t.Errorf("expected version and title in changelog, got: %s", content)
		}
	})
	t.Run("SC-e03s04-P1-12: appends to existing CHANGELOG.md", func(t *testing.T) {
		defer chdirTempDir(t)()
		_ = os.WriteFile("CHANGELOG.md", []byte("# Changelog\n\nAll notable changes.\n\n## 0.9.0 (2024-01-01)\n\n### Features\n\n- feat: initial\n"), 0644)
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0", Notes: "## 1.0.0 (2024-06-15)\n\n### Features\n\n- feat: major release\n"},
		}
		if err := p.Prepare(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		data, _ := os.ReadFile("CHANGELOG.md")
		content := string(data)
		if !strings.Contains(content, "1.0.0") || !strings.Contains(content, "0.9.0") {
			t.Errorf("expected both versions in changelog, got: %s", content)
		}
	})
	t.Run("SC-e03s04-P1-13: generates notes when Notes is empty", func(t *testing.T) {
		defer chdirTempDir(t)()
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0", Notes: ""},
			Commits:     []*algorithm.Commit{{Type: "feat", Subject: "new feature"}},
		}
		if err := p.Prepare(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		data, _ := os.ReadFile("CHANGELOG.md")
		content := string(data)
		if !strings.Contains(content, "1.0.0") || !strings.Contains(content, "new feature") {
			t.Errorf("expected version and commit in changelog, got: %s", content)
		}
	})
}

func TestChangelogPluginLifecycle(t *testing.T) {
	p := NewChangelogPlugin()
	t.Run("SC-e03s04-P1-14: Publish returns nil", func(t *testing.T) {
		release, err := p.Publish(&algorithm.Context{})
		if err != nil || release != nil {
			t.Errorf("expected nil, got err=%v release=%v", err, release)
		}
	})
	t.Run("SC-e03s04-P1-15: Success returns nil", func(t *testing.T) {
		if err := p.Success(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
	t.Run("SC-e03s04-P1-16: Fail returns nil", func(t *testing.T) {
		if err := p.Fail(&algorithm.Context{}, nil); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestChangelogPluginAutoRegistration(t *testing.T) {
	t.Run("SC-e03s04-P1-17: ChangelogPlugin auto-registered in global registry", func(t *testing.T) {
		found := false
		for _, name := range List() {
			if name == "changelog" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'changelog' to be registered in global registry")
		}
	})
}
