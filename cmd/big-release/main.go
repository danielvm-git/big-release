package main

import (
	"fmt"
	"os"

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

	rootCmd.AddCommand(releaseCmd, validateCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRelease(dryRun, verbose bool, configFile string) error {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

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

	if err := cfg.Validate(); err != nil {
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
