package swift

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// validVersionRegex matches semver-like version strings: start with alphanumeric,
// followed by alphanumeric, dots, underscores, or hyphens.
var validVersionRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

// Publisher publishes Swift packages via git tags.
type Publisher struct {
	// DryRun, when true, creates git tags locally but skips git push.
	DryRun bool
	// ExecCommand is the function used to run external commands. Defaults to exec.Command.
	ExecCommand func(name string, arg ...string) *exec.Cmd
}

// NewPublisher creates a new Swift Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		ExecCommand: exec.Command,
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "swift"
}

// Detect returns true when Package.swift exists in the working directory.
func (p *Publisher) Detect() bool {
	_, err := os.Stat("Package.swift")
	return err == nil
}

// Prepare returns nil — Swift uses git tag-based versioning, no file mutation needed.
func (p *Publisher) Prepare(version string) error {
	return nil
}

// Publish creates a versioned git tag and pushes it to origin.
// In dry-run mode, the tag is created locally but not pushed.
// Uses git.FormatTag with the configured tagFormat instead of bare version
// (BUG-swift-bypasses-tagformat).
func (p *Publisher) Publish(version string) error {
	if !isValidVersion(version) {
		return fmt.Errorf("swift: invalid version %q", version)
	}

	tag := p.resolveTag(version)

	// Create git tag.
	cmd := p.ExecCommand("git", "tag", "-a", tag, "-m", "release "+version)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to create git tag: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	if p.DryRun {
		return nil
	}

	// Push git tag.
	cmd = p.ExecCommand("git", "push", "origin", tag)
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to push git tag: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

// resolveTag returns the git tag for the given version using the configured
// tagFormat from .big-release.yml, falling back to bare version (BUG-swift-bypasses-tagformat).
func (p *Publisher) resolveTag(version string) string {
	cfgPath := ".big-release.yml"
	if _, err := os.Stat(cfgPath); err == nil {
		data, err := os.ReadFile(cfgPath)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "tagFormat:") {
					parts := strings.SplitN(trimmed, ":", 2)
					if len(parts) == 2 {
						fmt := strings.TrimSpace(parts[1])
						fmt = strings.Trim(fmt, `"'`)
						return git.FormatTag(fmt, version)
					}
				}
			}
		}
	}
	return version
}

// Verify checks that the versioned git tag exists locally.
// Uses resolveTag to match the tag created by Publish (BUG-swift-bypasses-tagformat).
func (p *Publisher) Verify(version string) error {
	tag := p.resolveTag(version)
	cmd := p.ExecCommand("git", "tag", "-l", tag)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to list git tags: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	found := strings.TrimSpace(stdout.String())
	if found != tag {
		return fmt.Errorf("swift: tag %q not found", tag)
	}

	return nil
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

func init() {
	publishers.Register(NewPublisher())
}

// --- internal helpers ---

// isValidVersion validates that the version string is safe for use in git commands.
func isValidVersion(v string) bool {
	return len(v) > 0 && len(v) <= 128 && validVersionRegex.MatchString(v)
}
