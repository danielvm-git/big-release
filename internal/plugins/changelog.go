package plugins

import (
	"fmt"
	"os"
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// ChangelogPlugin generates CHANGELOG.md
type ChangelogPlugin struct{}

// NewChangelogPlugin creates a new ChangelogPlugin
func NewChangelogPlugin() *ChangelogPlugin {
	return &ChangelogPlugin{}
}

// Name returns the plugin name
func (p *ChangelogPlugin) Name() string {
	return "changelog"
}

// VerifyConditions verifies pre-release conditions
func (p *ChangelogPlugin) VerifyConditions(ctx *algorithm.Context) error {
	return nil
}

// AnalyzeCommits analyzes commits and returns release type
func (p *ChangelogPlugin) AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error) {
	return "", nil
}

// GenerateNotes generates release notes
func (p *ChangelogPlugin) GenerateNotes(ctx *algorithm.Context) (string, error) {
	return "", nil
}

// Prepare prepares the release
func (p *ChangelogPlugin) Prepare(ctx *algorithm.Context) error {
	if ctx.DryRun {
		return nil
	}

	// Read existing CHANGELOG.md
	changelogPath := "CHANGELOG.md"
	existingContent := ""
	if data, err := os.ReadFile(changelogPath); err == nil {
		existingContent = string(data)
	}

	// Generate new content
	var sb strings.Builder

	// Add new release section
	sb.WriteString(fmt.Sprintf("## %s (%s)\n\n", ctx.NextRelease.Version, "2024-01-01"))
	sb.WriteString(ctx.NextRelease.Notes)
	sb.WriteString("\n\n")

	// Add existing content (skip header if exists)
	if existingContent != "" {
		// Find the first ## or # after the title
		lines := strings.Split(existingContent, "\n")
		startIdx := 0
		for i, line := range lines {
			if strings.HasPrefix(line, "## ") && i > 0 {
				startIdx = i
				break
			}
		}
		if startIdx > 0 {
			sb.WriteString(strings.Join(lines[startIdx:], "\n"))
		}
	} else {
		sb.WriteString("# Changelog\n\n")
		sb.WriteString("All notable changes to this project will be documented in this file.\n\n")
	}

	// Write to file
	if err := os.WriteFile(changelogPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}

	return nil
}

// Publish publishes the release
func (p *ChangelogPlugin) Publish(ctx *algorithm.Context) (*algorithm.Release, error) {
	return nil, nil
}

// Success is called after successful release
func (p *ChangelogPlugin) Success(ctx *algorithm.Context) error {
	return nil
}

// Fail is called on release failure
func (p *ChangelogPlugin) Fail(ctx *algorithm.Context, err error) error {
	return nil
}

func init() {
	Register(NewChangelogPlugin())
}
