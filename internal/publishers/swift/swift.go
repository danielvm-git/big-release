package swift

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
func (p *Publisher) Publish(version string) error {
	if !isValidVersion(version) {
		return fmt.Errorf("swift: invalid version %q", version)
	}

	// Create git tag.
	cmd := p.ExecCommand("git", "tag", version)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to create git tag: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	if p.DryRun {
		return nil
	}

	// Push git tag.
	cmd = p.ExecCommand("git", "push", "origin", version)
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to push git tag: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

// Verify checks that the versioned git tag exists locally.
func (p *Publisher) Verify(version string) error {
	cmd := p.ExecCommand("git", "tag", "-l", version)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swift: failed to list git tags: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	tag := strings.TrimSpace(stdout.String())
	if tag != version {
		return fmt.Errorf("swift: tag %q not found", version)
	}

	return nil
}

func init() {
	publishers.Register(NewPublisher())
}

// --- internal helpers ---

// isValidVersion validates that the version string is safe for use in git commands.
func isValidVersion(v string) bool {
	return len(v) > 0 && len(v) <= 128 && validVersionRegex.MatchString(v)
}
