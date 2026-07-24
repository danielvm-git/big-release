// story: e24s01

package pnpm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/danielvm-git/big-release/internal/publishers"
)

// Publisher publishes with pnpm.
type Publisher struct {
	// DryRun, when true, skips actual pnpm publish and pnpm view calls.
	DryRun bool
	// channel is the npm-compatible dist-tag; empty means default "latest".
	channel string
	// ExecCommand is the function used to run external commands. Defaults to exec.Command.
	ExecCommand func(name string, arg ...string) *exec.Cmd
}

// NewPublisher creates a new Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		ExecCommand: exec.Command,
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "pnpm"
}

// Detect returns true when pnpm-lock.yaml or pnpm-workspace.yaml is present.
func (p *Publisher) Detect() bool {
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return true
	}
	if _, err := os.Stat("pnpm-workspace.yaml"); err == nil {
		return true
	}
	return false
}

// Prepare prepares the package for publishing.
func (p *Publisher) Prepare(version string) error {
	pkg, err := readPackageJSON()
	if err != nil {
		return err
	}

	pkg["version"] = version

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("pnpm: failed to marshal package.json: %w", err)
	}

	if err := os.WriteFile("package.json", data, 0644); err != nil {
		return fmt.Errorf("pnpm: failed to write package.json: %w", err)
	}

	return nil
}

// Publish publishes the package.
func (p *Publisher) Publish(version string) error {
	if p.DryRun {
		return nil
	}

	args := []string{"publish", "--no-git-checks"}
	if tag := p.distTag(); tag != "" {
		args = append(args, "--tag", tag)
	}
	cmd := p.ExecCommand("pnpm", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pnpm: publish failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

func (p *Publisher) distTag() string {
	tag := strings.TrimSpace(p.channel)
	if tag == "" || tag == "latest" {
		return ""
	}
	return tag
}

// Verify verifies the publication.
func (p *Publisher) Verify(version string) error {
	name, err := readPackageName()
	if err != nil {
		return err
	}

	if p.DryRun {
		return nil
	}

	cmd := p.ExecCommand("pnpm", "view", name, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("pnpm: failed to verify publication: %w", err)
	}

	publishedVersion := strings.TrimSpace(string(output))
	if publishedVersion != version {
		return fmt.Errorf("pnpm: published version %s does not match expected version %s", publishedVersion, version)
	}

	return nil
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

// SetChannel sets the dist-tag from the release channel.
func (p *Publisher) SetChannel(channel string) {
	p.channel = channel
}

func init() {
	publishers.Register(NewPublisher())
}

var npmNamePattern = regexp.MustCompile(`^[a-z0-9][-a-z0-9._]*$`)
var npmScopePattern = regexp.MustCompile(`^@[a-z0-9][-a-z0-9._]*$`)

func isValidPackageName(name string) bool {
	if len(name) == 0 || len(name) > 214 {
		return false
	}
	if name[0] == '@' {
		if !strings.Contains(name, "/") {
			return false
		}
		scopeEnd := strings.Index(name, "/")
		scope := name[:scopeEnd]
		if !npmScopePattern.MatchString(scope) {
			return false
		}
		rest := name[scopeEnd+1:]
		if len(rest) == 0 {
			return false
		}
		return npmNamePattern.MatchString(rest)
	}
	return npmNamePattern.MatchString(name)
}

func readPackageJSON() (map[string]interface{}, error) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("pnpm: failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("pnpm: failed to parse package.json: %w", err)
	}

	return pkg, nil
}

func readPackageName() (string, error) {
	pkg, err := readPackageJSON()
	if err != nil {
		return "", err
	}

	name, ok := pkg["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("pnpm: package name not found or not a string in package.json")
	}
	if !isValidPackageName(name) {
		return "", fmt.Errorf("pnpm: invalid package name %q in package.json", name)
	}

	return name, nil
}
