// story: e03s04
package plugins

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// ChangelogPlugin generates CHANGELOG.md.
type ChangelogPlugin struct{}

// NewChangelogPlugin creates a new ChangelogPlugin.
func NewChangelogPlugin() *ChangelogPlugin {
	return &ChangelogPlugin{}
}

// Name returns the plugin name.
func (p *ChangelogPlugin) Name() string {
	return "changelog"
}

// VerifyConditions verifies pre-release conditions.
func (p *ChangelogPlugin) VerifyConditions(ctx *algorithm.Context) error {
	return nil
}

// AnalyzeCommits analyzes commits and returns release type.
func (p *ChangelogPlugin) AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error) {
	return "", nil
}

// filterCommits returns non-breaking commits matching the given type.
func filterCommits(commits []*algorithm.Commit, commitType string) []*algorithm.Commit {
	var result []*algorithm.Commit
	for _, c := range commits {
		if c.Type == commitType && !c.Breaking {
			result = append(result, c)
		}
	}
	return result
}

type commitCategory struct {
	title   string
	commits []*algorithm.Commit
}

func (p *ChangelogPlugin) writeCategorySection(sb *strings.Builder, cat commitCategory) {
	if len(cat.commits) == 0 {
		return
	}
	fmt.Fprintf(sb, "### %s\n\n", cat.title)
	for _, c := range cat.commits {
		scope := ""
		if c.Scope != "" {
			scope = fmt.Sprintf(" **%s**", c.Scope)
		}
		fmt.Fprintf(sb, "- %s%s: %s\n", c.Type, scope, c.Subject)
	}
	sb.WriteString("\n")
}

func (p *ChangelogPlugin) writeBreakingChanges(sb *strings.Builder, commits []*algorithm.Commit) {
	var breaking []*algorithm.Commit
	for _, c := range commits {
		if c.Breaking {
			breaking = append(breaking, c)
		}
	}
	if len(breaking) == 0 {
		return
	}
	sb.WriteString("### BREAKING CHANGES\n\n")
	for _, c := range breaking {
		fmt.Fprintf(sb, "- %s: %s\n", c.Subject, c.Body)
	}
	sb.WriteString("\n")
}

func defaultCategories(commits []*algorithm.Commit) []commitCategory {
	return []commitCategory{
		{title: "Features", commits: filterCommits(commits, "feat")},
		{title: "Bug Fixes", commits: filterCommits(commits, "fix")},
		{title: "Performance Improvements", commits: filterCommits(commits, "perf")},
		{title: "Documentation", commits: filterCommits(commits, "docs")},
		{title: "Refactoring", commits: filterCommits(commits, "refactor")},
		{title: "Chores", commits: filterCommits(commits, "chore")},
		{title: "Style", commits: filterCommits(commits, "style")},
		{title: "Tests", commits: filterCommits(commits, "test")},
	}
}

// GenerateNotes generates release notes from commits grouped by type.
func (p *ChangelogPlugin) GenerateNotes(ctx *algorithm.Context) (string, error) {
	if ctx.NextRelease == nil {
		return "", nil
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "## %s (%s)\n\n", ctx.NextRelease.Version, time.Now().Format("2006-01-02"))

	if len(ctx.Commits) == 0 {
		sb.WriteString("No significant changes.\n")
		return sb.String(), nil
	}
	for _, cat := range defaultCategories(ctx.Commits) {
		p.writeCategorySection(&sb, cat)
	}
	p.writeBreakingChanges(&sb, ctx.Commits)
	return strings.TrimSpace(sb.String()), nil
}

// findContentStartIdx scans existing changelog lines for the first "## " header after line 0.
func findContentStartIdx(lines []string) int {
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") && i > 0 {
			return i
		}
	}
	return -1
}

func newFileChangelog() string {
	return "# Changelog\n\nAll notable changes to this project will be documented in this file.\n\n"
}

func mergeIntoExisting(sb *strings.Builder, existingContent string) {
	lines := strings.Split(existingContent, "\n")
	startIdx := findContentStartIdx(lines)
	if startIdx > 0 {
		sb.WriteString(strings.Join(lines[startIdx:], "\n"))
		return
	}
	if !strings.HasPrefix(lines[0], "# ") {
		sb.WriteString(newFileChangelog())
	}
	sb.WriteString(existingContent)
}

func mergeChangelogContent(lastRelease string, existingContent string) string {
	var sb strings.Builder
	sb.WriteString(lastRelease)
	sb.WriteString("\n\n")

	if existingContent == "" {
		sb.WriteString(newFileChangelog())
		return sb.String()
	}
	mergeIntoExisting(&sb, existingContent)
	return sb.String()
}

func (p *ChangelogPlugin) resolveNotes(ctx *algorithm.Context) (string, error) {
	if ctx.NextRelease.Notes != "" {
		return ctx.NextRelease.Notes, nil
	}
	return p.GenerateNotes(ctx)
}

func (p *ChangelogPlugin) readChangelogFile() string {
	if data, err := os.ReadFile("CHANGELOG.md"); err == nil {
		return string(data)
	}
	return ""
}

// Prepare writes the release notes to CHANGELOG.md.
func (p *ChangelogPlugin) Prepare(ctx *algorithm.Context) error {
	if ctx.DryRun {
		return nil
	}
	notes, err := p.resolveNotes(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate release notes: %w", err)
	}
	merged := mergeChangelogContent(notes, p.readChangelogFile())
	trimmed := strings.TrimSpace(merged) + "\n"
	if err := os.WriteFile("CHANGELOG.md", []byte(trimmed), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}
	return nil
}

// Publish publishes the release.
func (p *ChangelogPlugin) Publish(ctx *algorithm.Context) (*algorithm.Release, error) {
	return nil, nil
}

// Success is called after successful release.
func (p *ChangelogPlugin) Success(ctx *algorithm.Context) error {
	return nil
}

// Fail is called on release failure.
func (p *ChangelogPlugin) Fail(ctx *algorithm.Context, err error) error {
	return nil
}

func init() {
	Register(NewChangelogPlugin())
}
