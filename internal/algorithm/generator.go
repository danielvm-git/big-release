package algorithm

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Sensitive patterns to hide
	sensitivePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(token|password|credential|secret|private).*[=:]\s*\S+`),
		regexp.MustCompile(`(?i)(token|password|credential|secret|private)=\S+`),
	}

	// issueRefPattern matches bare or keyword-prefixed issue references:
	//   #123, fixes #123, closes #123, resolves #123
	// Captures the keyword prefix (group 1) and the issue number (group 2).
	issueRefPattern = regexp.MustCompile(`(?i)(fixes|closes|resolves)?\s*#(\d+)`)
)

// Generator generates release notes.
type Generator struct {
	commitTypes []CommitTypeConfig
	repoURL     string
}

// NewGenerator creates a new Generator. When commitTypes is nil or empty,
// the generator falls back to DefaultCommitTypes(), preserving the
// historical default behavior.
func NewGenerator(commitTypes ...[]CommitTypeConfig) *Generator {
	g := &Generator{commitTypes: DefaultCommitTypes()}
	if len(commitTypes) > 0 && len(commitTypes[0]) > 0 {
		g.commitTypes = commitTypes[0]
	}
	return g
}

// SetRepositoryURL configures the repository URL used to render
// clickable commit, compare, and issue links in the generated notes.
// When unset (empty), the generator falls back to plain-text rendering.
func (g *Generator) SetRepositoryURL(repoURL string) {
	g.repoURL = strings.TrimRight(repoURL, "/")
}

// sectionEntry defines a single renderable section.
type sectionEntry struct {
	key    string // commit type ("breaking" for the synthetic breaking bucket)
	title  string // display title
	hidden bool
}

// breakingSection is the synthetic bucket for breaking changes from any
// commit type. It always renders first and is never hidden.
const breakingSection = "BREAKING CHANGES"

// sections derives the ordered, renderable section list from the
// generator's commit types. The synthetic "breaking" section is always
// prepended.
func (g *Generator) sections() []sectionEntry {
	out := make([]sectionEntry, 0, len(g.commitTypes)+1)
	out = append(out, sectionEntry{key: "breaking", title: breakingSection, hidden: false})
	for _, ct := range g.commitTypes {
		out = append(out, sectionEntry{
			key:    ct.Type,
			title:  ct.Section,
			hidden: ct.Hidden || ct.Section == "",
		})
	}
	return out
}

// GenerateNotes generates release notes from commits.
func (g *Generator) GenerateNotes(commits []*Commit, lastRelease *Release, nextRelease *Release) string {
	if len(commits) == 0 {
		return ""
	}

	var sb strings.Builder

	// Group commits by type
	groups := g.groupCommits(commits)

	// Generate sections in deterministic order
	for _, sec := range g.sections() {
		if sec.hidden {
			continue
		}
		if len(groups[sec.key]) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "### %s\n\n", sec.title)
		for _, commit := range groups[sec.key] {
			fmt.Fprintf(&sb, "- %s\n", g.formatCommit(commit))
		}
		sb.WriteString("\n")
	}

	// Add comparison link if there's a last release
	if lastRelease != nil && lastRelease.GitTag != "" && nextRelease != nil {
		sb.WriteString("\n---\n\n")
		g.writeComparison(&sb, lastRelease.GitTag, nextRelease.GitTag)
	}

	return g.hideSensitive(sb.String())
}

// writeComparison emits the "Full Changelog" line. With a repo URL it is
// a clickable markdown compare link; without one it falls back to the
// historical prose form.
func (g *Generator) writeComparison(sb *strings.Builder, lastTag, nextTag string) {
	if g.repoURL != "" {
		fmt.Fprintf(sb, "[Full Changelog](%s/compare/%s...%s)\n", g.repoURL, lastTag, nextTag)
		return
	}
	fmt.Fprintf(sb, "Full Changelog: comparing changes from %s to %s\n", lastTag, nextTag)
}

// groupCommits groups commits by type. Breaking commits go into the
// "breaking" bucket regardless of hidden status; non-breaking commits
// are grouped by their Type and skipped when their type is hidden.
func (g *Generator) groupCommits(commits []*Commit) map[string][]*Commit {
	groups := make(map[string][]*Commit)
	hidden := make(map[string]bool)
	for _, sec := range g.sections() {
		if sec.hidden {
			hidden[sec.key] = true
		}
	}

	for _, commit := range commits {
		if commit.Breaking {
			groups["breaking"] = append(groups["breaking"], commit)
			continue
		}
		if commit.Type != "" && !hidden[commit.Type] {
			groups[commit.Type] = append(groups[commit.Type], commit)
		}
	}

	return groups
}

// formatCommit formats a commit for release notes. When a repository URL
// is configured, the commit hash is rendered as a clickable markdown link
// and any issue references (#123, fixes #N, closes #N, resolves #N) are
// linkified. The hash is always rendered, even without a scope.
func (g *Generator) formatCommit(commit *Commit) string {
	subject := g.linkifyIssues(commit.Subject)
	scope := commit.Scope

	switch {
	case scope != "" && g.hasHash(commit):
		return fmt.Sprintf("**%s:** %s %s", scope, subject, g.commitHashLink(commit))
	case scope != "":
		return fmt.Sprintf("**%s:** %s", scope, subject)
	case g.hasHash(commit):
		return fmt.Sprintf("**%s:** %s %s", commit.Type, subject, g.commitHashLink(commit))
	default:
		return fmt.Sprintf("**%s:** %s", commit.Type, subject)
	}
}

// hasHash reports whether the commit carries a usable hash.
func (g *Generator) hasHash(commit *Commit) bool {
	return len(commit.Hash) >= 7
}

// commitHashLink renders the short hash. With a repo URL it is a clickable
// markdown link to the full commit; without one it is the historical
// parenthesized plain-text form.
func (g *Generator) commitHashLink(commit *Commit) string {
	short := commit.Hash[:7]
	if g.repoURL == "" {
		return fmt.Sprintf("(%s)", short)
	}
	return fmt.Sprintf("([%s](%s/commit/%s))", short, g.repoURL, commit.Hash)
}

// linkifyIssues turns "#123", "fixes #123", "closes #123", "resolves #123"
// references into clickable markdown links. Without a repo URL it returns
// the input unchanged.
func (g *Generator) linkifyIssues(text string) string {
	if g.repoURL == "" {
		return text
	}
	return issueRefPattern.ReplaceAllStringFunc(text, func(match string) string {
		sub := issueRefPattern.FindStringSubmatch(match)
		keyword := strings.TrimSpace(sub[1])
		num := sub[2]
		if keyword != "" {
			lowered := strings.ToLower(sub[1])
			return fmt.Sprintf("%s [#%s](%s/issues/%s)", lowered, num, g.repoURL, num)
		}
		return fmt.Sprintf("[#%s](%s/issues/%s)", num, g.repoURL, num)
	})
}

// hideSensitive hides sensitive information in release notes.
func (g *Generator) hideSensitive(text string) string {
	for _, pattern := range sensitivePatterns {
		text = pattern.ReplaceAllString(text, "[secure]")
	}
	return text
}

// DefaultCommitTypes returns the seed list of commit type configurations
// used when the user has not configured commitTypes.
//
// Per #7 (parity with @semantic-release/conventional-changelog-conventionalcommits),
// only release-relevant types are visible by default: feat, fix, perf,
// revert. The remaining types (docs, refactor, chore, style, test, build,
// ci) are hidden but still parsed — their commits are excluded from the
// changelog unless they carry a breaking change, which always surfaces
// in the BREAKING CHANGES section.
func DefaultCommitTypes() []CommitTypeConfig {
	return []CommitTypeConfig{
		{Type: "feat", Section: "Features"},
		{Type: "fix", Section: "Bug Fixes"},
		{Type: "perf", Section: "Performance Improvements"},
		{Type: "revert", Section: "Reverts"},
		{Type: "docs", Hidden: true},
		{Type: "refactor", Hidden: true},
		{Type: "build", Hidden: true},
		{Type: "ci", Hidden: true},
		{Type: "chore", Hidden: true},
		{Type: "style", Hidden: true},
		{Type: "test", Hidden: true},
	}
}
