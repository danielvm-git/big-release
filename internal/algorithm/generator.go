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
)

// Generator generates release notes.
type Generator struct {
	commitTypes []CommitTypeConfig
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
		fmt.Fprintf(&sb, "Full Changelog: comparing changes from %s to %s\n", lastRelease.GitTag, nextRelease.GitTag)
	}

	return g.hideSensitive(sb.String())
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

// formatCommit formats a commit for release notes.
func (g *Generator) formatCommit(commit *Commit) string {
	if commit.Scope != "" {
		return fmt.Sprintf("**%s:** %s (%s)", commit.Scope, commit.Subject, commit.Hash[:7])
	}
	return fmt.Sprintf("**%s:** %s", commit.Type, commit.Subject)
}

// hideSensitive hides sensitive information in release notes.
func (g *Generator) hideSensitive(text string) string {
	for _, pattern := range sensitivePatterns {
		text = pattern.ReplaceAllString(text, "[secure]")
	}
	return text
}

// DefaultCommitTypes returns the seed list of commit type configurations
// used when the user has not configured commitTypes. e18s01 keeps the
// historical defaults (all sections visible); e18s02 flips the
// non-release types to hidden.
func DefaultCommitTypes() []CommitTypeConfig {
	return []CommitTypeConfig{
		{Type: "feat", Section: "Features"},
		{Type: "fix", Section: "Bug Fixes"},
		{Type: "perf", Section: "Performance Improvements"},
		{Type: "docs", Section: "Documentation"},
		{Type: "refactor", Section: "Refactoring"},
		{Type: "chore", Section: "Chores"},
		{Type: "style", Section: "Style"},
		{Type: "test", Section: "Tests"},
	}
}
