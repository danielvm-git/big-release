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

// GenerateNotes generates release notes from commits
func (g *Generator) GenerateNotes(commits []*Commit, lastRelease *Release, nextRelease *Release) string {
	if len(commits) == 0 {
		return ""
	}

	var sb strings.Builder

	// Group commits by type
	groups := g.groupCommits(commits)

	// Generate sections
	if len(groups["breaking"]) > 0 {
		sb.WriteString("### ⚠ BREAKING CHANGES\n\n")
		for _, commit := range groups["breaking"] {
			sb.WriteString(fmt.Sprintf("- %s\n", g.formatCommit(commit)))
		}
		sb.WriteString("\n")
	}

	if len(groups["feat"]) > 0 {
		sb.WriteString("### Features\n\n")
		for _, commit := range groups["feat"] {
			sb.WriteString(fmt.Sprintf("- %s\n", g.formatCommit(commit)))
		}
		sb.WriteString("\n")
	}

	if len(groups["fix"]) > 0 {
		sb.WriteString("### Bug Fixes\n\n")
		for _, commit := range groups["fix"] {
			sb.WriteString(fmt.Sprintf("- %s\n", g.formatCommit(commit)))
		}
		sb.WriteString("\n")
	}

	if len(groups["perf"]) > 0 {
		sb.WriteString("### Performance Improvements\n\n")
		for _, commit := range groups["perf"] {
			sb.WriteString(fmt.Sprintf("- %s\n", g.formatCommit(commit)))
		}
		sb.WriteString("\n")
	}

	// Add comparison link if there's a last release
	if lastRelease != nil && lastRelease.GitTag != "" && nextRelease != nil {
		sb.WriteString(fmt.Sprintf("\n---\n\n"))
		sb.WriteString(fmt.Sprintf("Full Changelog: comparing changes from %s to %s\n", lastRelease.GitTag, nextRelease.GitTag))
	}

	return g.hideSensitive(sb.String())
}

// groupCommits groups commits by type
func (g *Generator) groupCommits(commits []*Commit) map[string][]*Commit {
	groups := make(map[string][]*Commit)

	for _, commit := range commits {
		// Check for breaking changes first
		if commit.Breaking {
			groups["breaking"] = append(groups["breaking"], commit)
			continue
		}

		// Group by type
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
