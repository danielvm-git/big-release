package npm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/danielvm-git/big-release/internal/publishers"
)

// Publisher publishes to npm
type Publisher struct{}

// NewPublisher creates a new Publisher
func NewPublisher() *Publisher {
	return &Publisher{}
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

// Prepare prepares the package for publishing
func (p *Publisher) Prepare(version string) error {
	// Read package.json
	data, err := os.ReadFile("package.json")
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	// Parse package.json
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Update version
	pkg["version"] = version

	// Write back
	data, err = json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	if err := os.WriteFile("package.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	return nil
}

// Publish publishes the package
func (p *Publisher) Publish(version string) error {
	cmd := exec.Command("npm", "publish")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to publish to npm: %w", err)
	}

	return nil
}

// Verify verifies the publication
func (p *Publisher) Verify(version string) error {
	// Read package.json to get package name
	data, err := os.ReadFile("package.json")
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	name, ok := pkg["name"].(string)
	if !ok {
		return fmt.Errorf("package name not found in package.json")
	}

	// Check if version exists on npm
	cmd := exec.Command("npm", "view", name, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to verify publication: %w", err)
	}

	publishedVersion := filepath.Base(string(output))
	if publishedVersion != version {
		return fmt.Errorf("published version %s does not match expected version %s", publishedVersion, version)
	}

	return nil
}

func init() {
	publishers.Register(NewPublisher())
}
