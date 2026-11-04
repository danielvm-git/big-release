// story: e01s03
package algorithm

import "testing"

func TestParseRevertedHash_Table(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{"empty body", "", ""},
		{"unrelated body", "just a normal body", ""},
		{"body sentence", "This reverts commit abcdef1234567890.", "abcdef1234567890"},
		{"footer", "description\n\nReverts: 1234567890abcdef", "1234567890abcdef"},
		{"footer is case insensitive", "Reverts: ABC1234", "ABC1234"},
		{"short hash from footer", "Reverts: abc1234", "abc1234"},
		{"footer wins over body sentence",
			"This reverts commit aaaaaaa1111.\n\nReverts: bbbbbbb2222", "bbbbbbb2222"},
		{"non-hex target is ignored", "This reverts commit zzzzzzz.", ""},
		{"too-short target is ignored", "This reverts commit abc.", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseRevertedHash(tc.body); got != tc.want {
				t.Fatalf("parseRevertedHash(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}

func TestFilterReverted_Table(t *testing.T) {
	target := &Commit{Type: "feat", Subject: "feature", Hash: "aaaaaaa111111111"}
	revert := &Commit{Type: "revert", Subject: "revert feature", Hash: "bbbbbbb222222222",
		Body: "This reverts commit aaaaaaa111111111."}
	orphan := &Commit{Type: "revert", Subject: "revert old", Hash: "ccccccc333333333",
		Body: "This reverts commit 9999999999999999."}

	cases := []struct {
		name    string
		commits []*Commit
		want    []*Commit
	}{
		{"nil list", nil, nil},
		{"no reverts", []*Commit{target}, []*Commit{target}},
		{"revert without body is kept", []*Commit{
			{Type: "revert", Subject: "empty", Hash: "ddddddd444444444"},
		}, []*Commit{{Type: "revert", Subject: "empty", Hash: "ddddddd444444444"}}},
		{"matched revert drops both", []*Commit{target, revert}, nil},
		{"orphan revert is kept", []*Commit{target, orphan}, []*Commit{target, orphan}},
		{"short target hash matches full commit hash", []*Commit{
			{Type: "feat", Subject: "feature", Hash: "abc1234ffffffff"},
			{Type: "revert", Subject: "revert", Hash: "fffffff000000000",
				Body: "Reverts: abc1234"},
		}, nil},
		{"matched and orphan reverts mixed", []*Commit{
			target, revert, orphan,
		}, []*Commit{orphan}},
		{"target hash longer than commit hash is not a match", []*Commit{
			{Type: "feat", Subject: "short", Hash: "abc1234"},
			{Type: "revert", Subject: "revert", Hash: "fffffff000000000",
				Body: "Reverts: abc1234ffff"},
		}, []*Commit{
			{Type: "feat", Subject: "short", Hash: "abc1234"},
			{Type: "revert", Subject: "revert", Hash: "fffffff000000000",
				Body: "Reverts: abc1234ffff"},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterReverted(tc.commits)
			if len(got) != len(tc.want) {
				t.Fatalf("filterReverted() returned %d commits, want %d", len(got), len(tc.want))
			}
			for i := range tc.want {
				if got[i].Hash != tc.want[i].Hash {
					t.Fatalf("filterReverted()[%d] = %q, want %q", i, got[i].Hash, tc.want[i].Hash)
				}
			}
		})
	}
}

func TestCommitExists_PrefixMatching(t *testing.T) {
	commits := []*Commit{{Hash: "abcdef1234567890"}}
	if !commitExists(commits, "abcdef1") {
		t.Error("commitExists() = false for short prefix, want true")
	}
	if commitExists(commits, "1234567") {
		t.Error("commitExists() = true for unrelated prefix, want false")
	}
	if commitExists(nil, "abcdef1") {
		t.Error("commitExists() = true for empty list, want false")
	}
}

func TestIsRevertedBy_PrefixMatching(t *testing.T) {
	reverted := map[string]bool{"abcdef1234567890": true}
	if !isRevertedBy("abcdef1234567890", reverted) {
		t.Error("isRevertedBy() = false for exact hash, want true")
	}
	if !isRevertedBy("abcdef1234567890ffff", reverted) {
		t.Error("isRevertedBy() = false for prefixed hash, want true")
	}
	if isRevertedBy("abcdef1", reverted) {
		t.Error("isRevertedBy() = true when commit hash is shorter than target, want false")
	}
	if isRevertedBy("abcdef1234567890", nil) {
		t.Error("isRevertedBy() = true for empty map, want false")
	}
}
