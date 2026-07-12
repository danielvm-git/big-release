package algorithm

import (
	"regexp"
	"strings"
)

var (
	// Conventional commit pattern: type(scope): description
	// Examples:
	//   feat(auth): add OAuth2 support
	//   fix!: remove deprecated API
	//   perf: optimize database queries
	conventionalCommitPattern = regexp.MustCompile(`^(?P<type>feat|fix|perf|docs|chore|style|refactor|test|build|ci|revert)(?:\((?P<scope>[^)]+)\))?(?P<breaking>!)?: (?P<description>.+)$`)
	
	// Breaking change pattern in commit body
	breakingChangePattern = regexp.MustCompile(`BREAKING CHANGE:\s*(.+)`)
)

// Analyzer analyzes commits to determine release type
type Analyzer struct{}

// NewAnalyzer creates a new Analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeCommits analyzes commits and returns the release type
func (a *Analyzer) AnalyzeCommits(commits []*Commit) ReleaseType {
	var releaseType ReleaseType

	for _, commit := range commits {
		// Skip commits with [skip release] or [release skip]
		if strings.Contains(commit.Message, "[skip release]") ||
			strings.Contains(commit.Message, "[release skip]") {
			continue
		}

		// Parse the commit
		a.parseCommit(commit)

		// Determine bump type
		bump := a.determineBump(commit)

		// Update release type (highest priority wins)
		if bump != "" {
			if releaseType == "" || bumpPriority(bump) > bumpPriority(releaseType) {
				releaseType = bump
			}
		}
	}

	return releaseType
}

// parseCommit parses a conventional commit message
func (a *Analyzer) parseCommit(commit *Commit) {
	// Extract type, scope, breaking, and description
	matches := conventionalCommitPattern.FindStringSubmatch(commit.Message)
	if matches == nil {
		// Not a conventional commit, treat as no release
		commit.Type = ""
		return
	}

	commit.Type = matches[1]
	commit.Scope = matches[2]
	commit.Breaking = matches[3] == "!"
	commit.Subject = matches[4]

	// Check for BREAKING CHANGE in commit body
	if breakingChangePattern.MatchString(commit.Body) {
		commit.Breaking = true
	}
}

// determineBump determines the bump type for a commit
func (a *Analyzer) determineBump(commit *Commit) ReleaseType {
	// Check for breaking changes first (highest priority)
	if commit.Breaking {
		return ReleaseTypeMajor
	}

	// Check commit type
	switch commit.Type {
	case "feat":
		return ReleaseTypeMinor
	case "fix", "perf":
		return ReleaseTypePatch
	default:
		// docs, chore, style, refactor, test, build, ci, revert
		// No release needed
		return ""
	}
}

// bumpPriority returns the priority of a bump type
func bumpPriority(bump ReleaseType) int {
	switch bump {
	case ReleaseTypeMajor:
		return 3
	case ReleaseTypeMinor:
		return 2
	case ReleaseTypePatch:
		return 1
	default:
		return 0
	}
}
