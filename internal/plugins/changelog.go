// story: e03s04
package plugins

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
func (p *ChangelogPlugin) VerifyConditions(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// AnalyzeCommits analyzes commits and returns release type.
func (p *ChangelogPlugin) AnalyzeCommits(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return "", nil
}

// VerifyRelease is not applicable for the changelog plugin.
func (p *ChangelogPlugin) VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// GenerateNotes returns release notes already computed by the algorithm
// Generator and stored in state.NextRelease.Notes by the orchestrator.
func (p *ChangelogPlugin) GenerateNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	if state.NextRelease == nil {
		return "", nil
	}
	return state.NextRelease.Notes, nil
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
	firstLine := ""
	if len(lines) > 0 {
		firstLine = lines[0]
	}
	if strings.HasPrefix(firstLine, "# ") {
		sb.WriteString(existingContent)
	} else {
		sb.WriteString(newFileChangelog())
		if firstLine != "" {
			sb.WriteString(existingContent)
		}
	}
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

func (p *ChangelogPlugin) resolveNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	if state.NextRelease.Notes != "" {
		return state.NextRelease.Notes, nil
	}
	return p.GenerateNotes(ctx, state)
}

func (p *ChangelogPlugin) readChangelogFile() (string, error) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err == nil {
		return string(data), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	return "", fmt.Errorf("failed to read CHANGELOG.md: %w", err)
}

// Prepare writes the release notes to CHANGELOG.md.
func (p *ChangelogPlugin) Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun {
		return nil
	}
	notes, err := p.resolveNotes(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to generate release notes: %w", err)
	}
	existing, err := p.readChangelogFile()
	if err != nil {
		return err
	}
	merged := mergeChangelogContent(notes, existing)
	trimmed := strings.TrimSpace(merged) + "\n"
	if err := os.WriteFile("CHANGELOG.md", []byte(trimmed), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}
	return nil
}

// Publish publishes the release.
func (p *ChangelogPlugin) Publish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error) {
	return nil, nil
}

// Success is called after successful release.
func (p *ChangelogPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// Fail is called on release failure.
func (p *ChangelogPlugin) Fail(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, err error) error {
	return nil
}

func init() {
	Register(NewChangelogPlugin())
}
