// story: e01s03
package algorithm

import (
	"strings"
	"testing"
)

func TestGenerator_FormatCommit_ScopedWithoutHash(t *testing.T) {
	notes := NewGenerator().GenerateNotes([]*Commit{
		{Type: "feat", Scope: "auth", Subject: "add OAuth"},
	}, nil, nil)

	if !strings.Contains(notes, "**auth:** add OAuth") {
		t.Fatalf("expected scoped formatting without hash:\n%s", notes)
	}
	if strings.Contains(notes, "(abc") {
		t.Fatalf("unexpected hash suffix:\n%s", notes)
	}
}

func TestGenerator_PerfRendersUnderChanged(t *testing.T) {
	notes := NewGenerator().GenerateNotes([]*Commit{
		{Type: "perf", Subject: "speed up queries"},
	}, nil, nil)

	if !strings.Contains(notes, "### Changed") {
		t.Fatalf("expected Changed section for perf:\n%s", notes)
	}
	if !strings.Contains(notes, "speed up queries") {
		t.Fatalf("expected perf subject in notes:\n%s", notes)
	}
	if strings.Contains(notes, "### Performance") {
		t.Fatalf("perf must not get its own section:\n%s", notes)
	}
}

func TestGenerator_HiddenPerfIsExcluded(t *testing.T) {
	types := []CommitTypeConfig{
		{Type: "perf", Hidden: true},
		{Type: "feat", Section: "Added"},
	}
	notes := NewGenerator(types).GenerateNotes([]*Commit{
		{Type: "perf", Subject: "speed up queries"},
		{Type: "feat", Subject: "add login"},
	}, nil, nil)

	if strings.Contains(notes, "speed up queries") {
		t.Fatalf("hidden perf commit must not render:\n%s", notes)
	}
	if !strings.Contains(notes, "add login") {
		t.Fatalf("visible feat commit must render:\n%s", notes)
	}
}

func TestGenerator_NoFixTypeStillRendersChanged(t *testing.T) {
	// Without a fix type the Changed bucket is prepended, so a breaking
	// commit always has a home.
	types := []CommitTypeConfig{{Type: "feat", Section: "Added"}}
	notes := NewGenerator(types).GenerateNotes([]*Commit{
		{Type: "feat", Subject: "rewrite core", Breaking: true},
	}, nil, nil)

	if !strings.Contains(notes, "### Changed") {
		t.Fatalf("expected prepended Changed section:\n%s", notes)
	}
	if !strings.Contains(notes, "rewrite core") {
		t.Fatalf("expected breaking subject:\n%s", notes)
	}
}

func TestGenerator_EmptyCommitTypesFallsBackToDefaults(t *testing.T) {
	notes := NewGenerator([]CommitTypeConfig{}).GenerateNotes([]*Commit{
		{Type: "feat", Subject: "add login"},
	}, nil, nil)

	if !strings.Contains(notes, "### Added") {
		t.Fatalf("empty commit types should fall back to defaults:\n%s", notes)
	}
}

func TestGenerator_EmptySectionTitleHidesType(t *testing.T) {
	types := []CommitTypeConfig{
		{Type: "feat"}, // no Section, not explicitly Hidden
		{Type: "fix", Section: "Fixed"},
	}
	notes := NewGenerator(types).GenerateNotes([]*Commit{
		{Type: "feat", Subject: "add login"},
		{Type: "fix", Subject: "fix crash"},
	}, nil, nil)

	if strings.Contains(notes, "add login") {
		t.Fatalf("type with empty section title must be hidden:\n%s", notes)
	}
	if !strings.Contains(notes, "fix crash") {
		t.Fatalf("visible fix commit must render:\n%s", notes)
	}
}

func TestGenerator_SetRepositoryURLTrimsTrailingSlash(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo/")
	notes := g.GenerateNotes([]*Commit{
		{Type: "feat", Subject: "add OAuth", Hash: "abc1234567890def"},
	}, nil, nil)

	if !strings.Contains(notes, "https://github.com/owner/repo/commit/") {
		t.Fatalf("expected single-slash commit link:\n%s", notes)
	}
	if strings.Contains(notes, "repo//commit") {
		t.Fatalf("trailing slash was not trimmed:\n%s", notes)
	}
}

func TestGenerator_LinkifiesBareAndUppercaseIssues(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	notes := g.GenerateNotes([]*Commit{
		{Type: "fix", Subject: "Fixes #7 and also #42"},
	}, nil, nil)

	for _, want := range []string{
		"fixes [#7](https://github.com/owner/repo/issues/7)",
		"[#42](https://github.com/owner/repo/issues/42)",
	} {
		if !strings.Contains(notes, want) {
			t.Fatalf("expected %q in notes:\n%s", want, notes)
		}
	}
}

func TestGenerator_MultipleCommitsPerSection(t *testing.T) {
	notes := NewGenerator().GenerateNotes([]*Commit{
		{Type: "feat", Subject: "first feature"},
		{Type: "feat", Subject: "second feature"},
	}, nil, nil)

	if strings.Count(notes, "### Added") != 1 {
		t.Fatalf("expected exactly one Added section:\n%s", notes)
	}
	for _, want := range []string{"first feature", "second feature"} {
		if !strings.Contains(notes, want) {
			t.Fatalf("expected %q in notes:\n%s", want, notes)
		}
	}
}
