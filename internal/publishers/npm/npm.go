// story: e02s07 e22s02 e24s01 e24s03

package npm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/nodeutil"
)

// Publisher publishes to npm.
type Publisher struct {
	// DryRun, when true, skips actual npm publish and npm view calls.
	DryRun bool
	// channel is the npm dist-tag (e22); empty means default "latest".
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

// Name returns the publisher name
func (p *Publisher) Name() string {
	return "npm"
}

// Detect detects if this publisher should be used.
// Returns false when pnpm markers are present so pnpm owns the project (e24s01).
func (p *Publisher) Detect() bool {
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return false
	}
	if _, err := os.Stat("pnpm-workspace.yaml"); err == nil {
		return false
	}
	_, err := os.Stat("package.json")
	return err == nil
}

// Prepare prepares the package for publishing.
func (p *Publisher) Prepare(version string) error {
	pkg, err := nodeutil.ReadPackageJSON("npm")
	if err != nil {
		return err
	}

	pkg["version"] = version
	return nodeutil.WritePackageJSON("npm", pkg)
}

// Publish publishes the package.
func (p *Publisher) Publish(version string) error {
	if p.DryRun {
		return nil
	}

	args := []string{"publish"}
	if tag := p.distTag(); tag != "" {
		args = append(args, "--tag", tag)
	}
	cmd := p.ExecCommand("npm", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm: publish failed: %w: %s", err, strings.TrimSpace(stderr.String()))
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
	name, err := nodeutil.ReadPackageName("npm")
	if err != nil {
		return err
	}

	if p.DryRun {
		return nil
	}

	cmd := p.ExecCommand("npm", "view", "--", name, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("npm: failed to verify publication: %w", err)
	}

	publishedVersion := strings.TrimSpace(string(output))
	if publishedVersion != version {
		return fmt.Errorf("npm: published version %s does not match expected version %s", publishedVersion, version)
	}

	return nil
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

// SetChannel sets the npm dist-tag from the release channel (e22s02).
func (p *Publisher) SetChannel(channel string) {
	p.channel = channel
}

func init() {
	publishers.Register(NewPublisher())
}
