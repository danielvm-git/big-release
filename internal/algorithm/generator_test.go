package algorithm

import (
	"strings"
	"testing"
)

// --- e18s01: configurable commit type sections & visibility (plumbing) ---

func TestGenerator_WithConfigurableCommitTypes_RespectsSectionTitle(t *testing.T) {
	// Custom section title for feat.
	types := DefaultCommitTypes()
	types[0].Section = "Awesome Features" // feat is first in DefaultCommitTypes

	g := NewGenerator(types)
	commits := []*Commit{
		{Type: "feat", Subject: "add login"},
	}
	notes := g.GenerateNotes(commits, nil, nil)
	if !strings.Contains(notes, "### Awesome Features") {
		t.Errorf("expected custom section title, got:\n%s", notes)
	}
}

func TestGenerator_WithConfigurableCommitTypes_HidesHiddenType(t *testing.T) {
	types := DefaultCommitTypes()
	// Mark chore as hidden.
	for i := range types {
		if types[i].Type == "chore" {
			types[i].Hidden = true
		}
	}

	g := NewGenerator(types)
	commits := []*Commit{
		{Type: "feat", Subject: "add login"},
		{Type: "chore", Subject: "update deps"},
	}
	notes := g.GenerateNotes(commits, nil, nil)
	if strings.Contains(notes, "update deps") {
		t.Errorf("hidden chore commit must not appear in notes:\n%s", notes)
	}
	if !strings.Contains(notes, "add login") {
		t.Errorf("non-hidden feat commit must appear in notes:\n%s", notes)
	}
}

func TestGenerator_WithConfigurableCommitTypes_BreakingStillShownFromHiddenType(t *testing.T) {
	types := DefaultCommitTypes()
	for i := range types {
		if types[i].Type == "refactor" {
			types[i].Hidden = true
		}
	}

	g := NewGenerator(types)
	commits := []*Commit{
		{Type: "refactor", Subject: "rewrite core", Breaking: true, Body: "BREAKING CHANGE: everything"},
	}
	notes := g.GenerateNotes(commits, nil, nil)
	if !strings.Contains(notes, "### BREAKING CHANGES") {
		t.Errorf("breaking change from a hidden type must still appear:\n%s", notes)
	}
	if !strings.Contains(notes, "rewrite core") {
		t.Errorf("breaking change subject must appear in BREAKING CHANGES:\n%s", notes)
	}
}

func TestDefaultCommitTypes_CoversAllCurrentTypes(t *testing.T) {
	types := DefaultCommitTypes()
	seen := map[string]bool{}
	for _, ct := range types {
		seen[ct.Type] = true
	}
	for _, want := range []string{"feat", "fix", "perf", "docs", "refactor", "chore", "style", "test"} {
		if !seen[want] {
			t.Errorf("DefaultCommitTypes missing %q: %+v", want, seen)
		}
	}
}

func TestNewGenerator_DefaultsToStandardCommitTypes(t *testing.T) {
	// NewGenerator() with no args must behave identically to NewGenerator(DefaultCommitTypes())
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "add login"},
		{Type: "fix", Subject: "fix crash"},
	}
	notes := g.GenerateNotes(commits, nil, nil)
	// All standard sections still render with default constructor.
	for _, want := range []string{"### Features", "### Bug Fixes"} {
		if !strings.Contains(notes, want) {
			t.Errorf("default constructor missing %q:\n%s", want, notes)
		}
	}
}

func TestGenerator_GenerateNotes_Empty(t *testing.T) {
	g := NewGenerator()
	notes := g.GenerateNotes(nil, nil, nil)
	if notes != "" {
		t.Errorf("expected empty notes for nil commits, got %q", notes)
	}
}

func TestGenerator_GenerateNotes_AllCommitTypes(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "add login"},
		{Type: "fix", Subject: "fix crash"},
		{Type: "perf", Subject: "speed up queries"},
		{Type: "docs", Subject: "update readme"},
		{Type: "refactor", Subject: "extract util"},
		{Type: "chore", Subject: "update deps"},
		{Type: "style", Subject: "format code"},
		{Type: "test", Subject: "add unit tests"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	for _, want := range []string{
		"### Features",
		"### Bug Fixes",
		"### Performance Improvements",
		"### Documentation",
		"### Refactoring",
		"### Chores",
		"### Style",
		"### Tests",
		"add login",
		"fix crash",
		"speed up queries",
		"update readme",
		"extract util",
		"update deps",
		"format code",
		"add unit tests",
	} {
		if !strings.Contains(notes, want) {
			t.Errorf("notes missing %q:\n%s", want, notes)
		}
	}
}

func TestGenerator_GenerateNotes_BreakingChanges(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "new API", Breaking: true, Body: "Removes v1"},
		{Type: "fix", Subject: "bug fix"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "### BREAKING CHANGES") {
		t.Errorf("missing BREAKING CHANGES section:\n%s", notes)
	}
	if !strings.Contains(notes, "new API") {
		t.Errorf("missing breaking commit subject:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_BreakingCommitNotInTypeSection(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "normal feature"},
		{Type: "feat", Subject: "breaking feature", Breaking: true},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	// BREAKING CHANGES comes before Features in sectionOrder.
	// The breaking commit should appear in BREAKING CHANGES, not Features.
	featuresIdx := strings.Index(notes, "### Features")
	if featuresIdx == -1 {
		t.Fatalf("expected Features section:\n%s", notes)
	}
	// Everything after Features header is the Features section (until end or next section)
	featuresSection := notes[featuresIdx:]
	if strings.Contains(featuresSection, "breaking feature") {
		t.Errorf("breaking commit should not appear in Features section:\n%s", featuresSection)
	}
	// Verify it IS in the BREAKING CHANGES section
	breakingIdx := strings.Index(notes, "### BREAKING CHANGES")
	if breakingIdx == -1 {
		t.Fatalf("expected BREAKING CHANGES section:\n%s", notes)
	}
	breakingSection := notes[breakingIdx:featuresIdx]
	if !strings.Contains(breakingSection, "breaking feature") {
		t.Errorf("breaking commit should appear in BREAKING CHANGES section:\n%s", breakingSection)
	}
}

func TestGenerator_GenerateNotes_ComparisonLink(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "new feature"},
	}
	lastRelease := &Release{GitTag: "v1.0.0"}
	nextRelease := &Release{GitTag: "v1.1.0"}

	notes := g.GenerateNotes(commits, lastRelease, nextRelease)

	if !strings.Contains(notes, "v1.0.0") || !strings.Contains(notes, "v1.1.0") {
		t.Errorf("missing comparison link:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_NoComparisonLinkWithoutLastRelease(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "new feature"},
	}

	notes := g.GenerateNotes(commits, nil, &Release{GitTag: "v1.0.0"})

	if strings.Contains(notes, "Full Changelog") {
		t.Errorf("should not have comparison link without last release:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_ScopeFormatting(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Scope: "auth", Subject: "add OAuth", Hash: "abc123456789"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "**auth:** add OAuth (abc1234)") {
		t.Errorf("expected scoped formatting:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_NoScopeFormatting(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "fix", Subject: "fix bug"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "**fix:** fix bug") {
		t.Errorf("expected type-based formatting:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_HidesSensitive(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "fix", Subject: "fix token=secret123 in config"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if strings.Contains(notes, "secret123") {
		t.Errorf("sensitive data should be hidden:\n%s", notes)
	}
}

func TestGenerator_SectionOrder(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "test", Subject: "add test"},
		{Type: "feat", Subject: "add feature"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	featIdx := strings.Index(notes, "### Features")
	testIdx := strings.Index(notes, "### Tests")
	if featIdx == -1 || testIdx == -1 {
		t.Fatalf("expected both sections:\n%s", notes)
	}
	if featIdx > testIdx {
		t.Errorf("Features should come before Tests:\n%s", notes)
	}
}
