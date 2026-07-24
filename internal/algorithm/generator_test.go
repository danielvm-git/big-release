package algorithm

import (
	"strings"
	"testing"
)

// --- e18s04 (#9): filter revert commits from release notes ---

func TestGenerator_GenerateNotes_RevertRemovesBothRevertAndTarget(t *testing.T) {
	// If commit A is reverted by commit B, neither appears in notes.
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Scope: "auth", Subject: "add OAuth", Hash: "aaaaaaa111111111"},
		{Type: "revert", Subject: "Revert \"add OAuth\"", Hash: "bbbbbbb222222222", Body: "This reverts commit aaaaaaa111111111."},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if strings.Contains(notes, "add OAuth") {
		t.Errorf("reverted target must not appear in notes:\n%s", notes)
	}
	if strings.Contains(notes, "Revert") {
		t.Errorf("revert commit must not appear in notes:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_RevertFromRevertsFooter(t *testing.T) {
	// The Reverts: footer form (alternative to the body sentence).
	g := NewGenerator()
	commits := []*Commit{
		{Type: "fix", Subject: "fix bug", Hash: "ccccccc333333333"},
		{Type: "revert", Subject: "undo fix", Hash: "ddddddd444444444", Body: "Reverts: ccccccc333333333"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if strings.Contains(notes, "fix bug") {
		t.Errorf("reverted target must not appear:\n%s", notes)
	}
	if strings.Contains(notes, "undo fix") {
		t.Errorf("revert commit must not appear:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_OrphanedRevertStillShown(t *testing.T) {
	// A revert whose target is not in the commit list still appears.
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "unrelated feature", Hash: "eeeeeee555555555"},
		{Type: "revert", Subject: "Revert something old", Hash: "fffffff666666666", Body: "This reverts commit 9999999999999999."},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "### Removed") {
		t.Errorf("orphaned revert must render under Removed section:\n%s", notes)
	}
	if !strings.Contains(notes, "Revert something old") {
		t.Errorf("orphaned revert subject must appear:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_NonRevertedFeatureStillShown(t *testing.T) {
	// Unrelated feat commits are unaffected by revert filtering.
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "keep this", Hash: "1111111111111111"},
		{Type: "revert", Subject: "Revert something", Hash: "2222222222222222", Body: "This reverts commit 9999999999999999."},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "keep this") {
		t.Errorf("non-reverted feat must still appear:\n%s", notes)
	}
}

func TestParseRevertedHash_BodySentence(t *testing.T) {
	body := "This reverts commit abcdef1234567890abcdef1234567890abcdef12.\n\nReason: broken"
	got := parseRevertedHash(body)
	want := "abcdef1234567890abcdef1234567890abcdef12"
	if got != want {
		t.Errorf("parseRevertedHash body form: got %q want %q", got, want)
	}
}

func TestParseRevertedHash_RevertsFooter(t *testing.T) {
	body := "Some description\n\nReverts: 1234567890abcdef1234567890abcdef12345678"
	got := parseRevertedHash(body)
	want := "1234567890abcdef1234567890abcdef12345678"
	if got != want {
		t.Errorf("parseRevertedHash footer form: got %q want %q", got, want)
	}
}

func TestParseRevertedHash_None(t *testing.T) {
	if got := parseRevertedHash("just a normal body"); got != "" {
		t.Errorf("expected empty hash for non-revert body, got %q", got)
	}
}

func TestParseRevertedHash_ShortHashFromFooter(t *testing.T) {
	// Footers may carry a short hash; we still return it for matching.
	body := "Reverts: abc1234"
	got := parseRevertedHash(body)
	if got != "abc1234" {
		t.Errorf("expected short hash abc1234, got %q", got)
	}
}

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
	if !strings.Contains(notes, "### Changed") {
		t.Errorf("breaking change from a hidden type must still appear:\n%s", notes)
	}
	if !strings.Contains(notes, "rewrite core") {
		t.Errorf("breaking change subject must appear in Changed:\n%s", notes)
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
	for _, want := range []string{"### Added", "### Fixed"} {
		if !strings.Contains(notes, want) {
			t.Errorf("default constructor missing %q:\n%s", want, notes)
		}
	}
}

// BUG-changelog-format: defaults map to Keep-a-Changelog 1.1.0 categories.
func TestDefaultCommitTypes_KeepAChangelogSections(t *testing.T) {
	want := map[string]string{
		"feat":   "Added",
		"fix":    "Fixed",
		"perf":   "Changed",
		"revert": "Removed",
	}
	for _, ct := range DefaultCommitTypes() {
		if section, ok := want[ct.Type]; ok && ct.Section != section {
			t.Errorf("type %q: section = %q, want %q", ct.Type, ct.Section, section)
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

func TestGenerator_GenerateNotes_DefaultHidesNonReleaseTypes(t *testing.T) {
	// #7: by default, only release-relevant types appear in the changelog.
	// Hidden: docs, chore, refactor, style, test, build, ci.
	// Visible: feat, fix, perf, revert.
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

	// Visible sections + subjects (Keep-a-Changelog categories).
	for _, want := range []string{
		"### Added",
		"### Fixed",
		"### Changed",
		"add login",
		"fix crash",
		"speed up queries",
	} {
		if !strings.Contains(notes, want) {
			t.Errorf("notes missing visible %q:\n%s", want, notes)
		}
	}

	// Hidden sections + subjects must not appear.
	for _, unwanted := range []string{
		"### Documentation",
		"### Refactoring",
		"### Chores",
		"### Style",
		"### Tests",
		"update readme",
		"extract util",
		"update deps",
		"format code",
		"add unit tests",
	} {
		if strings.Contains(notes, unwanted) {
			t.Errorf("hidden %q must not appear in notes:\n%s", unwanted, notes)
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

	if !strings.Contains(notes, "### Changed") {
		t.Errorf("missing Changed section:\n%s", notes)
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

	addedIdx := strings.Index(notes, "### Added")
	changedIdx := strings.Index(notes, "### Changed")
	if addedIdx == -1 || changedIdx == -1 {
		t.Fatalf("expected Added and Changed sections:\n%s", notes)
	}
	addedSection := notes[addedIdx:changedIdx]
	if strings.Contains(addedSection, "breaking feature") {
		t.Errorf("breaking commit should not appear in Added section:\n%s", addedSection)
	}
	changedSection := notes[changedIdx:]
	if !strings.Contains(changedSection, "breaking feature") {
		t.Errorf("breaking commit should appear in Changed section:\n%s", changedSection)
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
		{Type: "fix", Subject: "add fix"},
		{Type: "feat", Subject: "add feature"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	addedIdx := strings.Index(notes, "### Added")
	fixedIdx := strings.Index(notes, "### Fixed")
	if addedIdx == -1 || fixedIdx == -1 {
		t.Fatalf("expected both sections:\n%s", notes)
	}
	if addedIdx > fixedIdx {
		t.Errorf("Added should come before Fixed:\n%s", notes)
	}
}

// --- e18s02 (#7): breaking changes from hidden types still appear ---

func TestGenerator_GenerateNotes_BreakingFromHiddenTypeStillShown(t *testing.T) {
	// refactor is hidden by default, but a breaking refactor must surface
	// in Changed per Keep-a-Changelog.
	g := NewGenerator()
	commits := []*Commit{
		{Type: "refactor", Subject: "rewrite core", Breaking: true, Body: "BREAKING CHANGE: everything"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "### Changed") {
		t.Errorf("breaking change from hidden type must still appear:\n%s", notes)
	}
	if !strings.Contains(notes, "rewrite core") {
		t.Errorf("breaking change subject must appear in Changed:\n%s", notes)
	}
	// The refactor section itself must NOT render (it is hidden).
	if strings.Contains(notes, "### Refactoring") {
		t.Errorf("hidden Refactoring section must not render:\n%s", notes)
	}
}

// --- e18s03 (#8): clickable commit hash and issue links ---

func TestGenerator_GenerateNotes_ClickableCommitHash_Scoped(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	commits := []*Commit{
		{Type: "feat", Scope: "auth", Subject: "add OAuth", Hash: "abc1234567890def"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	want := "[abc1234](https://github.com/owner/repo/commit/abc1234567890def)"
	if !strings.Contains(notes, want) {
		t.Errorf("expected clickable commit hash %q:\n%s", want, notes)
	}
}

func TestGenerator_GenerateNotes_ClickableCommitHash_NoRepoURL(t *testing.T) {
	// Without a repo URL, fall back to the plain-text hash (current behavior).
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Scope: "auth", Subject: "add OAuth", Hash: "abc1234567890def"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if !strings.Contains(notes, "(abc1234)") {
		t.Errorf("expected plain-text hash fallback when no repo URL:\n%s", notes)
	}
	if strings.Contains(notes, "](http") {
		t.Errorf("should not emit a link without repo URL:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_HashRenderedEvenWithoutScope(t *testing.T) {
	// #8: the no-scope case previously omitted the hash entirely.
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	commits := []*Commit{
		{Type: "fix", Subject: "fix bug", Hash: "def567890abcdef1"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	want := "[def5678](https://github.com/owner/repo/commit/def567890abcdef1)"
	if !strings.Contains(notes, want) {
		t.Errorf("expected hash link in no-scope case %q:\n%s", want, notes)
	}
}

func TestGenerator_GenerateNotes_ClickableCompareLink(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	commits := []*Commit{
		{Type: "feat", Subject: "new feature"},
	}
	lastRelease := &Release{GitTag: "v1.0.0"}
	nextRelease := &Release{GitTag: "v1.1.0"}

	notes := g.GenerateNotes(commits, lastRelease, nextRelease)

	want := "[Full Changelog](https://github.com/owner/repo/compare/v1.0.0...v1.1.0)"
	if !strings.Contains(notes, want) {
		t.Errorf("expected clickable compare link %q:\n%s", want, notes)
	}
}

func TestGenerator_GenerateNotes_CompareLinkWithoutRepoURL_OmitsLink(t *testing.T) {
	// No repo URL: keep the historical prose form (do not emit a broken link).
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "new feature"},
	}
	lastRelease := &Release{GitTag: "v1.0.0"}
	nextRelease := &Release{GitTag: "v1.1.0"}

	notes := g.GenerateNotes(commits, lastRelease, nextRelease)

	if strings.Contains(notes, "](http") {
		t.Errorf("should not emit compare link without repo URL:\n%s", notes)
	}
	if !strings.Contains(notes, "v1.0.0") || !strings.Contains(notes, "v1.1.0") {
		t.Errorf("should still mention tags in prose form:\n%s", notes)
	}
}

func TestGenerator_GenerateNotes_IssueReferenceLinkified(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	commits := []*Commit{
		{Type: "fix", Subject: "resolve crash, fixes #456", Hash: "abcdef1234567890"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	want := "[#456](https://github.com/owner/repo/issues/456)"
	if !strings.Contains(notes, want) {
		t.Errorf("expected issue reference linkified %q:\n%s", want, notes)
	}
}

func TestGenerator_GenerateNotes_ClosesResolvesLinkified(t *testing.T) {
	g := NewGenerator()
	g.SetRepositoryURL("https://github.com/owner/repo")
	commits := []*Commit{
		{Type: "feat", Subject: "add thing, closes #100 and resolves #200", Hash: "abcdef1234567890"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	for _, want := range []string{
		"[#100](https://github.com/owner/repo/issues/100)",
		"[#200](https://github.com/owner/repo/issues/200)",
	} {
		if !strings.Contains(notes, want) {
			t.Errorf("expected %q:\n%s", want, notes)
		}
	}
}

func TestGenerator_GenerateNotes_IssueRefWithoutRepoURL_NotLinkified(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "fix", Subject: "fix crash, fixes #456", Hash: "abcdef1234567890"},
	}

	notes := g.GenerateNotes(commits, nil, nil)

	if strings.Contains(notes, "](http") {
		t.Errorf("should not linkify without repo URL:\n%s", notes)
	}
	if !strings.Contains(notes, "#456") {
		t.Errorf("should keep literal #456:\n%s", notes)
	}
}

// --- e08s01: secret masking in release notes ---

func TestMasking_GeneratorNotes(t *testing.T) {
	g := NewGenerator()
	commits := []*Commit{
		{Type: "feat", Subject: "add auth token=ghp_notleaked1234567890123456789012", Hash: "abc1234567890"},
	}
	notes := g.GenerateNotes(commits, nil, nil)
	if strings.Contains(notes, "ghp_notleaked") {
		t.Fatalf("release notes leaked token:\n%s", notes)
	}
	if !strings.Contains(notes, "[secure]") {
		t.Fatalf("expected redacted placeholder in notes:\n%s", notes)
	}
}
