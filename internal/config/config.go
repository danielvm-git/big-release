package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

const (
	// DefaultConfigFile is the default configuration file name
	DefaultConfigFile = ".big-release.yml"

	// DefaultTagFormat is the default tag format
	DefaultTagFormat = "v${version}"
)

// DefaultConfig returns the default configuration
func DefaultConfig() *algorithm.Config {
	return &algorithm.Config{
		Branches: []algorithm.BranchConfig{
			{Name: "main"},
			{Name: "next"},
			{Name: "N.x"},
			{Name: "next-major"},
			{Name: "beta", Prerelease: "beta"},
			{Name: "alpha", Prerelease: "alpha"},
		},
		TagFormat:      DefaultTagFormat,
		InitialVersion: "0.1.0",
		Plugins: []string{
			"changelog",
			"git",
			"github",
		},
		Publishers: map[string]algorithm.PublisherConfig{
			"npm":     {Enabled: true},
			"pypi":    {Enabled: true},
			"crates":  {Enabled: true},
			"goproxy": {Enabled: true},
		},
	}
}

// Load loads configuration from file or returns default
func Load(configFile string) (*algorithm.Config, error) {
	if configFile == "" {
		configFile = findConfigFile()
	}

	if configFile == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// findConfigFile searches for configuration file in current directory and parents
func findConfigFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		// Check for config file
		configPath := filepath.Join(dir, DefaultConfigFile)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// Check for alternative config files
		altConfigs := []string{
			".big-release.json",
			".big-release.yaml",
			"big-release.config.js",
		}
		for _, alt := range altConfigs {
			configPath = filepath.Join(dir, alt)
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// ValidateConfig validates the configuration
func ValidateConfig(c *algorithm.Config) error {
	// Validate branches
	if len(c.Branches) == 0 {
		return fmt.Errorf("at least one branch must be configured")
	}

	// Validate tag format
	if c.TagFormat == "" {
		return fmt.Errorf("tag format must not be empty")
	}

	// Validate release branches (max 3)
	releaseCount := 0
	for _, branch := range c.Branches {
		if branch.Type == "" || branch.Type == "release" {
			releaseCount++
		}
	}
	if releaseCount > 3 {
		return fmt.Errorf("maximum 3 release branches allowed, got %d", releaseCount)
	}

	// Validate branch types
	validTypes := map[string]bool{"": true, "release": true, "maintenance": true, "prerelease": true}
	for _, branch := range c.Branches {
		if !validTypes[branch.Type] {
			return fmt.Errorf("invalid branch type %q for branch %q: must be one of release, maintenance, prerelease, or empty", branch.Type, branch.Name)
		}
	}

	// Validate maintenance branches
	for _, branch := range c.Branches {
		if branch.Type == "maintenance" {
			// Maintenance branches should have a range pattern
			// For now, we just validate they have a name
			if branch.Name == "" {
				return fmt.Errorf("maintenance branch must have a name")
			}
		}
	}

	// Validate prerelease branches
	for _, branch := range c.Branches {
		if branch.Prerelease != "" {
			// Prerelease branches should have a valid prerelease identifier
			if branch.Name == "" {
				return fmt.Errorf("prerelease branch must have a name")
			}
		}
	}

	return nil
}
