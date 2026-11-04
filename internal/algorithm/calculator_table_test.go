// story: e01s02
package algorithm

import (
	"strings"
	"testing"
)

func TestCalculateNextVersion_Table(t *testing.T) {
	release := &Branch{Name: "main", Type: BranchTypeRelease}
	maintenance := &Branch{Name: "1.x", Type: BranchTypeMaintenance}

	cases := []struct {
		name    string
		last    *Release
		rel     ReleaseType
		branch  *Branch
		initial string
		want    string
		wantErr string
	}{
		{
			name: "first release default", rel: ReleaseTypeMinor,
			branch: release, want: FirstRelease,
		},
		{
			name: "first release uses initial version", rel: ReleaseTypeMinor,
			branch: release, initial: "0.5.0", want: "0.5.0",
		},
		{
			name: "first release treats empty last version as absent",
			last: &Release{Version: ""}, rel: ReleaseTypeMinor,
			branch: release, want: FirstRelease,
		},
		{
			name: "first release on prerelease branch", rel: ReleaseTypeMinor,
			branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want:   "1.0.0-beta.1",
		},
		{
			name: "first prerelease honors initial version", rel: ReleaseTypeMinor,
			branch:  &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			initial: "2.0.0", want: "2.0.0-beta.1",
		},
		{
			name: "regular major", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypeMajor, branch: release, want: "2.0.0",
		},
		{
			name: "regular minor", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypeMinor, branch: release, want: "1.3.0",
		},
		{
			name: "regular patch", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypePatch, branch: release, want: "1.2.4",
		},
		{
			name: "nil branch defaults to release", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypePatch, branch: nil, want: "1.2.4",
		},
		{
			name: "prerelease type on release branch is unreleasable",
			last: &Release{Version: "1.2.3"}, rel: ReleaseTypePrerelease,
			branch: release, wantErr: "no releasable commits",
		},
		{
			name: "unknown release type is unreleasable",
			last: &Release{Version: "1.2.3"}, rel: ReleaseType("weird"),
			branch: release, wantErr: "no releasable commits",
		},
		{
			name: "invalid last version", last: &Release{Version: "not-a-version"},
			rel: ReleaseTypePatch, branch: release, wantErr: "invalid last version",
		},
		{
			name: "maintenance major", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypeMajor, branch: maintenance, want: "2.0.0",
		},
		{
			name: "maintenance minor", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypeMinor, branch: maintenance, want: "1.3.0",
		},
		{
			name: "maintenance patch", last: &Release{Version: "1.2.3"},
			rel: ReleaseTypePatch, branch: maintenance, want: "1.2.4",
		},
		{
			name: "maintenance empty type defaults to patch",
			last: &Release{Version: "1.2.3"}, rel: ReleaseType(""),
			branch: maintenance, want: "1.2.4",
		},
		{
			name: "prerelease starts minor series", last: &Release{Version: "1.0.0"},
			rel: ReleaseTypeMinor, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.1.0-beta.1",
		},
		{
			name: "prerelease starts patch series", last: &Release{Version: "1.0.0"},
			rel: ReleaseTypePatch, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.0.1-beta.1",
		},
		{
			name: "prerelease starts major series", last: &Release{Version: "1.0.0"},
			rel: ReleaseTypeMajor, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "2.0.0-beta.1",
		},
		{
			name: "prerelease preid falls back to branch name",
			last: &Release{Version: "1.0.0"}, rel: ReleaseTypeMinor,
			branch: &Branch{Name: "next", Type: BranchTypePrerelease},
			want:   "1.1.0-next.1",
		},
		{
			name: "prerelease literal true falls back to branch name",
			last: &Release{Version: "1.0.0"}, rel: ReleaseTypeMinor,
			branch: &Branch{Name: "next", Type: BranchTypePrerelease, Prerelease: "true"},
			want:   "1.1.0-next.1",
		},
		{
			name: "prerelease increments same series", last: &Release{Version: "1.1.0-beta.1"},
			rel: ReleaseTypePatch, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.1.0-beta.2",
		},
		{
			name: "prerelease increment ignores release type", last: &Release{Version: "1.1.0-beta.1"},
			rel: ReleaseTypeMajor, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.1.0-beta.2",
		},
		{
			name: "prerelease increments past nine", last: &Release{Version: "1.1.0-beta.9"},
			rel: ReleaseTypePatch, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.1.0-beta.10",
		},
		{
			name: "prerelease switches series", last: &Release{Version: "1.1.0-alpha.3"},
			rel: ReleaseTypePatch, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.1.1-beta.1",
		},
		{
			name: "prerelease without number stays put", last: &Release{Version: "1.0.0-beta"},
			rel: ReleaseTypePatch, branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want: "1.0.0-beta",
		},
		{
			name: "prerelease extra parts ignored on increment",
			last: &Release{Version: "1.0.0-beta.1.9"}, rel: ReleaseTypePatch,
			branch: &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			want:   "1.0.0-beta.2",
		},
		{
			name: "prerelease unknown release type errors",
			last: &Release{Version: "1.0.0-alpha.1"}, rel: ReleaseType(""),
			branch:  &Branch{Type: BranchTypePrerelease, Prerelease: "beta"},
			wantErr: "cannot start prerelease series",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewCalculator().CalculateNextVersion(tc.last, tc.rel, tc.branch, tc.initial)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error = %v, want substring %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("version = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSplitPrerelease(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", nil},
		{"single part", "beta", []string{"beta"}},
		{"two parts", "beta.1", []string{"beta", "1"}},
		{"empty segments are skipped", "beta..1", []string{"beta", "1"}},
		{"trailing dot is ignored", "beta.1.", []string{"beta", "1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitPrerelease(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("splitPrerelease(%q) = %v, want %v", tc.input, got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("splitPrerelease(%q) = %v, want %v", tc.input, got, tc.want)
				}
			}
		})
	}
}
