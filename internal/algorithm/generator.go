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

	// IssueRefPattern matches bare or keyword-prefixed issue references:
	//   #123, fixes #123, closes #123, resolves #123
	// Captures the keyword prefix (group 1) and the issue number (group 2).
	// Exported so plugins (e.g. github) can reuse without duplication.
	IssueRefPattern = regexp.MustCompile(`(?i)(fixes|closes|resolves)?\s*#(\d+)`)
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

// changedSection is the Keep-a-Changelog bucket for breaking changes and
// performance improvements. It renders after Fixed and before Removed.
const changedSection = "Changed"

// sections derives the ordered, renderable section list from the
// generator's commit types. Breaking changes and perf commits share the
// synthetic "changed" section per Keep-a-Changelog 1.1.0.
func (g *Generator) sections() []sectionEntry {
	out := make([]sectionEntry, 0, len(g.commitTypes)+1)
	insertedChanged := false
	for _, ct := range g.commitTypes {
		if ct.Type == "perf" {
			continue
		}
		out = append(out, sectionEntry{
			key:    ct.Type,
			title:  ct.Section,
			hidden: ct.Hidden || ct.Section == "",
		})
		if ct.Type == "fix" {
			out = append(out, sectionEntry{key: "changed", title: changedSection, hidden: false})
			insertedChanged = true
		}
	}
	if !insertedChanged {
		// Ensure breaking/perf commits always have a renderable home.
		out = append([]sectionEntry{{key: "changed", title: changedSection, hidden: false}}, out...)
	}
	return out
}

// GenerateNotes generates release notes from commits.
func (g *Generator) GenerateNotes(commits []*Commit, lastRelease *Release, nextRelease *Release) string {
	if len(commits) == 0 {
		return ""
	}

	// Drop revert pairs (a revert and its matched target) from the notes.
	// Orphaned reverts are retained and render under the Removed section.
	commits = filterReverted(commits)
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

// groupCommits groups commits by type. Breaking commits and perf commits go
// into the "changed" bucket regardless of hidden status; non-breaking
// commits are grouped by their Type and skipped when their type is hidden.
func (g *Generator) groupCommits(commits []*Commit) map[string][]*Commit {
	groups := make(map[string][]*Commit)
	hidden := make(map[string]bool)
	for _, ct := range g.commitTypes {
		if ct.Hidden || ct.Section == "" {
			hidden[ct.Type] = true
		}
	}

	for _, commit := range commits {
		if commit.Breaking {
			groups["changed"] = append(groups["changed"], commit)
			continue
		}
		if commit.Type == "perf" && !hidden["perf"] {
			groups["changed"] = append(groups["changed"], commit)
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
	return IssueRefPattern.ReplaceAllStringFunc(text, func(match string) string {
		sub := IssueRefPattern.FindStringSubmatch(match)
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
// Per Keep-a-Changelog 1.1.0, only release-relevant types are visible by
// default: feat→Added, fix→Fixed, perf→Changed, revert→Removed. The
// remaining types (docs, refactor, chore, style, test, build, ci) are hidden
// but still parsed — their commits are excluded from the changelog unless
// they carry a breaking change, which surfaces in the Changed section.
func DefaultCommitTypes() []CommitTypeConfig {
	return []CommitTypeConfig{
		{Type: "feat", Section: "Added"},
		{Type: "fix", Section: "Fixed"},
		{Type: "perf", Section: "Changed"},
		{Type: "revert", Section: "Removed"},
		{Type: "docs", Hidden: true},
		{Type: "refactor", Hidden: true},
		{Type: "build", Hidden: true},
		{Type: "ci", Hidden: true},
		{Type: "chore", Hidden: true},
		{Type: "style", Hidden: true},
		{Type: "test", Hidden: true},
	}
}
