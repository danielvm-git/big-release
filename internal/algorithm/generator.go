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

// Generator generates release notes
type Generator struct{}

// NewGenerator creates a new Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// sectionOrder defines the display order for commit type sections.
var sectionOrder = []struct {
	key   string
	title string
}{
	{"breaking", "BREAKING CHANGES"},
	{"feat", "Features"},
	{"fix", "Bug Fixes"},
	{"perf", "Performance Improvements"},
	{"docs", "Documentation"},
	{"refactor", "Refactoring"},
	{"chore", "Chores"},
	{"style", "Style"},
	{"test", "Tests"},
}

// GenerateNotes generates release notes from commits
func (g *Generator) GenerateNotes(commits []*Commit, lastRelease *Release, nextRelease *Release) string {
	if len(commits) == 0 {
		return ""
	}

	var sb strings.Builder

	// Group commits by type
	groups := g.groupCommits(commits)

	// Generate sections in deterministic order
	for _, sec := range sectionOrder {
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

// groupCommits groups commits by type. Breaking commits go into
// the "breaking" bucket; non-breaking commits are grouped by their Type.
func (g *Generator) groupCommits(commits []*Commit) map[string][]*Commit {
	groups := make(map[string][]*Commit)

	for _, commit := range commits {
		if commit.Breaking {
			groups["breaking"] = append(groups["breaking"], commit)
			continue
		}
		if commit.Type != "" {
			groups[commit.Type] = append(groups[commit.Type], commit)
		}
	}

	return groups
}

// formatCommit formats a commit for release notes
func (g *Generator) formatCommit(commit *Commit) string {
	if commit.Scope != "" {
		return fmt.Sprintf("**%s:** %s (%s)", commit.Scope, commit.Subject, commit.Hash[:7])
	}
	return fmt.Sprintf("**%s:** %s", commit.Type, commit.Subject)
}

// hideSensitive hides sensitive information in release notes
func (g *Generator) hideSensitive(text string) string {
	for _, pattern := range sensitivePatterns {
		text = pattern.ReplaceAllString(text, "[secure]")
	}
	return text
}
