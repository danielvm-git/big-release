// story: e02s07

package npm

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

// Publisher publishes to npm.
type Publisher struct {
	// DryRun, when true, skips actual npm publish and npm view calls.
	DryRun bool
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

// Detect detects if this publisher should be used
func (p *Publisher) Detect() bool {
	_, err := os.Stat("package.json")
	return err == nil
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
		return fmt.Errorf("npm: failed to marshal package.json: %w", err)
	}

	if err := os.WriteFile("package.json", data, 0644); err != nil {
		return fmt.Errorf("npm: failed to write package.json: %w", err)
	}

	return nil
}

// Publish publishes the package.
func (p *Publisher) Publish(version string) error {
	if p.DryRun {
		return nil
	}

	cmd := p.ExecCommand("npm", "publish")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm: publish failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
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

func init() {
	publishers.Register(NewPublisher())
}

// npmNamePattern matches valid npm package names per the npm registry spec.
// Does not include scope prefix; for scoped names use isValidPackageName.
var npmNamePattern = regexp.MustCompile(`^[a-z0-9][-a-z0-9._]*$`)

// npmScopePattern matches valid npm scope names (after the initial @ and before /).
var npmScopePattern = regexp.MustCompile(`^@[a-z0-9][-a-z0-9._]*$`)

// isValidPackageName validates an npm package name.
// Supports both unscoped ("my-package") and scoped ("@scope/my-package") formats.
// Enforces npm's name rules to prevent flag injection via crafted names.
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

// readPackageJSON reads and parses the package.json file in the working directory.
func readPackageJSON() (map[string]interface{}, error) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("npm: failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("npm: failed to parse package.json: %w", err)
	}

	return pkg, nil
}

// readPackageName reads and validates the package name from package.json.
func readPackageName() (string, error) {
	pkg, err := readPackageJSON()
	if err != nil {
		return "", err
	}

	name, ok := pkg["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("npm: package name not found or not a string in package.json")
	}
	if !isValidPackageName(name) {
		return "", fmt.Errorf("npm: invalid package name %q in package.json", name)
	}

	return name, nil
}
