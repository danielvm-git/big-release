// story: e20s01
package git

import "testing"

func TestClient_BranchGlobbing(t *testing.T) {
	cases := []struct {
		pattern string
		branch  string
		want    bool
	}{
		{"main", "main", true},
		{"main", "master", false},
		{"1.x", "1.x", true},
		{"1.x", "1.0.x", true},
		{"release/*", "release/v1.0.0", true},
		{"release/*", "main", false},
	}

	for _, tc := range cases {
		got := MatchBranchPattern(tc.pattern, tc.branch)
		if got != tc.want {
			t.Errorf("MatchBranchPattern(%q, %q) = %v; want %v", tc.pattern, tc.branch, got, tc.want)
		}
	}
}
