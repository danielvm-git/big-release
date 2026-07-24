package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/config"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/pkg/release"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "big-release",
		Short: "Unified release automation for all languages",
		Long: `big-release is a unified release tool that automatically:
- Analyzes commits using Conventional Commits
- Determines the next version (patch, minor, major)
- Generates changelogs from commit history
- Creates git tags with proper formatting
- Publishes packages to any registry
- Creates GitHub releases with assets`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildDate),
	}

	// Global flags
	var dryRun bool
	var verbose bool
	var configFile string

	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Run in dry-run mode (no writes)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default: .big-release.yml)")

	// Release command
	releaseCmd := &cobra.Command{
		Use:   "release",
		Short: "Run the release process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRelease(dryRun, verbose, configFile)
		},
	}

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(configFile)
		},
	}

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(configFile)
		},
	}

	// Health command
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check system health",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHealth()
		},
	}

	rootCmd.AddCommand(releaseCmd, validateCmd, versionCmd, healthCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// buildLogger constructs the zap logger used for the release run.
//
// Without --verbose it uses the production JSON encoder at Info level
// (machine-parseable, for log shippers). With --verbose it switches to the
// development console encoder at Debug level so a successful release is
// legible in CI output — the production logger writes JSON to stderr and,
// combined with the sparse success-path logs, previously made a release
// appear completely silent (BUG-release-workflow-softprops-and-verbose).
func buildLogger(verbose bool) (*zap.Logger, error) {
	if verbose {
		devCfg := zap.NewDevelopmentConfig()
		devCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		return devCfg.Build()
	}
	prodCfg := zap.NewProductionConfig()
	prodCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	return prodCfg.Build()
}

func runRelease(dryRun, verbose bool, configFile string) error {
	// Initialize logger
	logger, err := buildLogger(verbose)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() { _ = logger.Sync() }()

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize git client
	gitClient, err := git.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize git client: %w", err)
	}

	// Create release context
	ctx := &release.Context{
		Config:  cfg,
		Git:     gitClient,
		Logger:  logger,
		DryRun:  dryRun,
		Verbose: verbose,
	}

	// Run release
	releaser := release.New(ctx)
	return releaser.Run()
}

func runValidate(configFile string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("configuration invalid: %w", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("configuration invalid: %w", err)
	}

	fmt.Println("✅ Configuration is valid")
	return nil
}

func runVersion(configFile string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	gitClient, err := git.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize git client: %w", err)
	}

	lastRelease, err := gitClient.GetLastRelease(cfg.TagFormat)
	if err != nil {
		return fmt.Errorf("failed to get last release: %w", err)
	}

	if lastRelease == nil {
		fmt.Println("No releases found")
		return nil
	}

	fmt.Printf("Current version: %s\n", lastRelease.Version)
	fmt.Printf("Git tag: %s\n", lastRelease.GitTag)
	return nil
}

type healthResult struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func runHealth() error {
	checks := make(map[string]string)
	healthy := true

	// Check git
	if out, err := exec.Command("git", "version").Output(); err == nil {
		checks["git"] = string(out[:len(out)-1])
	} else {
		checks["git"] = "NOT FOUND"
		healthy = false
	}

	// Check git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err == nil {
		checks["git_repo"] = "ok"
	} else {
		checks["git_repo"] = "NOT A GIT REPOSITORY"
		healthy = false
	}

	// Check go
	if out, err := exec.Command("go", "version").Output(); err == nil {
		checks["go"] = string(out[:len(out)-1])
	} else {
		checks["go"] = "NOT FOUND"
		healthy = false
	}

	// Check config
	if _, err := config.Load(""); err == nil {
		checks["config"] = "valid"
	} else {
		checks["config"] = fmt.Sprintf("invalid: %v", err)
		healthy = false
	}

	result := healthResult{
		Status: "ok",
		Checks: checks,
	}
	if !healthy {
		result.Status = "degraded"
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("failed to encode health result: %w", err)
	}

	if !healthy {
		return fmt.Errorf("health check failed")
	}
	return nil
}
