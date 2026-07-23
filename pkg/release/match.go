package release

import (
	"regexp"
	"strings"
)

// matchBranchPattern implements a focused subset of extglob matching
// (the kind semantic-release/micromatch uses) sufficient for branch name
// patterns like:
//
//	+([0-9]).x       → 1.x, 2.x, 12.x  (one or more digits, then ".x")
//	release/+([0-9]) → release/1, release/2
//
// It converts the extglob pattern to a Go regular expression and anchors
// it. Unsupported patterns fall back to false (callers should already have
// tried an exact match via matchBranch).
func matchBranchPattern(pattern, branch string) bool {
	re, err := extglobToRegexp(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(branch)
}

// extglobToRegexp converts a small subset of bash extglob patterns into a
// compiled regular expression anchored to the full string. Supported:
//
//   - +(...)  → (?:...)+     (one or more of the group)
//   - *       → [^/]*        (glob star, not crossing /)
//   - ?       → [^/]         (single char, not /)
//   - literal characters are regexp-escaped
//
// Anything else is treated as literal. The group contents inside +(...)
// are passed through as a regexp body (e.g. [0-9]).
func extglobToRegexp(pattern string) (*regexp.Regexp, error) {
	var sb strings.Builder
	sb.WriteString("^")
	i := 0
	for i < len(pattern) {
		c := pattern[i]
		switch {
		case c == '+' && i+1 < len(pattern) && pattern[i+1] == '(':
			// +(...) extglob → (?:...)+
			depth := 1
			j := i + 2
			for j < len(pattern) && depth > 0 {
				switch pattern[j] {
				case '(':
					depth++
				case ')':
					depth--
				}
				if depth > 0 {
					j++
				}
			}
			if depth != 0 {
				// Unbalanced — treat the rest literally.
				sb.WriteString(regexp.QuoteMeta(pattern[i:]))
				return regexp.Compile(sb.String() + "$")
			}
			inner := pattern[i+2 : j] // contents between ( and )
			sb.WriteString("(?:")
			sb.WriteString(inner)
			sb.WriteString(")+")
			i = j + 1
		case c == '*':
			sb.WriteString("[^/]*")
			i++
		case c == '?':
			sb.WriteString("[^/]")
			i++
		default:
			// Escape literal characters, accumulating runs.
			start := i
			for i < len(pattern) {
				ch := pattern[i]
				if ch == '+' || ch == '*' || ch == '?' {
					// Stop if +(...) extglob or glob char follows.
					if ch != '+' || (i+1 < len(pattern) && pattern[i+1] == '(') {
						break
					}
				}
				i++
			}
			sb.WriteString(regexp.QuoteMeta(pattern[start:i]))
		}
	}
	sb.WriteString("$")
	return regexp.Compile(sb.String())
}
