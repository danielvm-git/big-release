// story: e01s01
package algorithm

import "testing"

func TestAnalyzer_AnalyzeCommits_Table(t *testing.T) {
	cases := []struct {
		name    string
		commits []*Commit
		want    ReleaseType
	}{
		{
			name:    "empty list",
			commits: nil,
			want:    "",
		},
		{
			name:    "feature is minor",
			commits: []*Commit{{Message: "feat: add login"}},
			want:    ReleaseTypeMinor,
		},
		{
			name:    "fix is patch",
			commits: []*Commit{{Message: "fix: repair crash"}},
			want:    ReleaseTypePatch,
		},
		{
			name:    "performance is patch",
			commits: []*Commit{{Message: "perf: speed up queries"}},
			want:    ReleaseTypePatch,
		},
		{
			name: "non-release types produce no bump",
			commits: []*Commit{
				{Message: "docs: update readme"},
				{Message: "chore: bump deps"},
				{Message: "style: format"},
				{Message: "refactor: extract"},
				{Message: "test: add cases"},
				{Message: "build: wire target"},
				{Message: "ci: add job"},
				{Message: "revert: undo something"},
			},
			want: "",
		},
		{
			name:    "bang marks breaking",
			commits: []*Commit{{Message: "feat!: remove old API"}},
			want:    ReleaseTypeMajor,
		},
		{
			name:    "breaking body marks breaking",
			commits: []*Commit{{Message: "fix: tweak", Body: "BREAKING CHANGE: config shape changed"}},
			want:    ReleaseTypeMajor,
		},
		{
			name:    "non-conventional commit is ignored",
			commits: []*Commit{{Message: "just a message"}},
			want:    "",
		},
		{
			name:    "skip release marker is ignored",
			commits: []*Commit{{Message: "feat: hidden [skip release]"}},
			want:    "",
		},
		{
			name:    "release skip marker is ignored",
			commits: []*Commit{{Message: "feat: hidden [release skip]"}},
			want:    "",
		},
		{
			name: "highest priority wins regardless of order",
			commits: []*Commit{
				{Message: "fix: first"},
				{Message: "feat: second"},
			},
			want: ReleaseTypeMinor,
		},
		{
			name: "major beats earlier minor",
			commits: []*Commit{
				{Message: "feat: first"},
				{Message: "refactor!: second"},
			},
			want: ReleaseTypeMajor,
		},
		{
			name: "later patch does not downgrade minor",
			commits: []*Commit{
				{Message: "feat: first"},
				{Message: "fix: second"},
			},
			want: ReleaseTypeMinor,
		},
		{
			name: "skip marker does not hide other commits",
			commits: []*Commit{
				{Message: "feat!: skipped [skip release]"},
				{Message: "fix: counted"},
			},
			want: ReleaseTypePatch,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NewAnalyzer().AnalyzeCommits(tc.commits); got != tc.want {
				t.Fatalf("AnalyzeCommits() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAnalyzer_ParsesCommitFields(t *testing.T) {
	commit := &Commit{Message: "feat(api): add pagination", Body: "details"}
	NewAnalyzer().AnalyzeCommits([]*Commit{commit})

	if commit.Type != "feat" {
		t.Errorf("Type = %q, want feat", commit.Type)
	}
	if commit.Scope != "api" {
		t.Errorf("Scope = %q, want api", commit.Scope)
	}
	if commit.Subject != "add pagination" {
		t.Errorf("Subject = %q, want add pagination", commit.Subject)
	}
	if commit.Breaking {
		t.Error("Breaking = true, want false")
	}
}

func TestAnalyzer_NonConventionalResetsType(t *testing.T) {
	// A commit carrying a stale Type must not be treated as releasable when
	// its message does not parse.
	commit := &Commit{Message: "not conventional", Type: "feat"}
	if got := NewAnalyzer().AnalyzeCommits([]*Commit{commit}); got != "" {
		t.Fatalf("AnalyzeCommits() = %q, want no release", got)
	}
	if commit.Type != "" {
		t.Fatalf("Type = %q, want empty", commit.Type)
	}
}

func TestAnalyzer_BreakingBodyIgnoredOnNonConventional(t *testing.T) {
	// parseCommit returns before the body check, so a non-conventional
	// message with a breaking body cannot force a major bump.
	commit := &Commit{Message: "not conventional", Body: "BREAKING CHANGE: nope"}
	if got := NewAnalyzer().AnalyzeCommits([]*Commit{commit}); got != "" {
		t.Fatalf("AnalyzeCommits() = %q, want no release", got)
	}
}

func TestBumpPriority(t *testing.T) {
	cases := []struct {
		bump ReleaseType
		want int
	}{
		{ReleaseTypeMajor, 3},
		{ReleaseTypeMinor, 2},
		{ReleaseTypePatch, 1},
		{ReleaseTypePrerelease, 0},
		{ReleaseType(""), 0},
	}
	for _, tc := range cases {
		if got := BumpPriority(tc.bump); got != tc.want {
			t.Errorf("BumpPriority(%q) = %d, want %d", tc.bump, got, tc.want)
		}
	}
}
