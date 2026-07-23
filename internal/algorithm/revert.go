package algorithm

import (
	"regexp"
	"strings"
)

// revertBodyPattern matches the "This reverts commit {hash}." body form.
// Captures the commit hash (full 40-char or short).
var revertBodyPattern = regexp.MustCompile(`This reverts commit ([0-9a-f]{7,40})`)

// revertFooterPattern matches the "Reverts: {hash}" footer form.
var revertFooterPattern = regexp.MustCompile(`(?im)^\s*Reverts:\s*([0-9a-f]{7,40})`)

// parseRevertedHash extracts the commit hash targeted by a revert, from
// either the "Reverts:" footer or the "This reverts commit {hash}." body
// sentence. Returns "" if the body is not a revert reference.
func parseRevertedHash(body string) string {
	if body == "" {
		return ""
	}
	if m := revertFooterPattern.FindStringSubmatch(body); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	if m := revertBodyPattern.FindStringSubmatch(body); len(m) == 2 {
		return m[1]
	}
	return ""
}

// filterReverted removes both revert commits and the commits they revert
// from the list, mirroring semantic-release's conventional-commits-filter.
//
// If a revert references a commit that is present in the list, BOTH are
// dropped. An orphaned revert (target not in the list) is kept and still
// appears under the Reverts section. Matching is by full or short hash
// prefix — a revert carries the full hash in its body, while commits in
// the list always carry their full hash, so we prefix-match.
func filterReverted(commits []*Commit) []*Commit {
	// Collect hashes that are reverted by some revert commit in the list.
	revertedHashes := make(map[string]bool)
	for _, c := range commits {
		if c.Type != "revert" {
			continue
		}
		if h := parseRevertedHash(c.Body); h != "" {
			revertedHashes[h] = true
		}
	}

	if len(revertedHashes) == 0 {
		return commits
	}

	out := make([]*Commit, 0, len(commits))
	for _, c := range commits {
		if c.Type == "revert" {
			// Drop the revert only if its target is present (matched).
			if h := parseRevertedHash(c.Body); h != "" && revertedHashes[h] {
				// Confirm the target actually exists in the list; orphaned
				// reverts (target not among commits) stay.
				if commitExists(commits, h) {
					continue
				}
			}
			out = append(out, c)
			continue
		}
		// Drop any commit whose hash is reverted by a revert in the list.
		if isRevertedBy(c.Hash, revertedHashes) {
			continue
		}
		out = append(out, c)
	}
	return out
}

// commitExists reports whether any commit's hash matches or starts with h.
func commitExists(commits []*Commit, h string) bool {
	for _, c := range commits {
		if c.Hash == h || strings.HasPrefix(c.Hash, h) {
			return true
		}
	}
	return false
}

// isRevertedBy reports whether hash is targeted by any entry in
// revertedHashes (full match or prefix).
func isRevertedBy(hash string, revertedHashes map[string]bool) bool {
	for h := range revertedHashes {
		if hash == h || strings.HasPrefix(hash, h) {
			return true
		}
	}
	return false
}
