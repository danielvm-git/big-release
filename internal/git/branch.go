// story: e20s01
package git

import (
	"path"
	"regexp"
	"strings"
)

// MatchBranchPattern matches branch names against exact strings, wildcards, or version globs (e.g. 1.x).
func MatchBranchPattern(pattern, branch string) bool {
	if pattern == branch {
		return true
	}

	// 1. Try standard glob match
	if matched, err := path.Match(pattern, branch); err == nil && matched {
		return true
	}

	// 2. Handle N.x / N.N.x patterns common in release channels (e.g. 1.x -> 1\..*)
	if strings.HasSuffix(pattern, ".x") {
		prefix := strings.TrimSuffix(pattern, ".x")
		escaped := regexp.QuoteMeta(prefix)
		re, err := regexp.Compile("^" + escaped + `(\..+)?(\.x)?$`)
		if err == nil && re.MatchString(branch) {
			return true
		}
	}

	return false
}
