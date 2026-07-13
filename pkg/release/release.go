package release

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// Context carries dependencies into the release process.
type Context struct {
	Config  *algorithm.Config
	Git     *git.Client
	Logger  *zap.Logger
	DryRun  bool
	Verbose bool
}

// commitAnalyzer is the interface for analyzing commits.
type commitAnalyzer interface {
	AnalyzeCommits(commits []*algorithm.Commit) algorithm.ReleaseType
}

// versionCalculator is the interface for calculating next version.
type versionCalculator interface {
	CalculateNextVersion(lastRelease *algorithm.Release, releaseType algorithm.ReleaseType, branch *algorithm.Branch) (string, error)
}

// notesGenerator is the interface for generating release notes.
type notesGenerator interface {
	GenerateNotes(commits []*algorithm.Commit, lastRelease *algorithm.Release, nextRelease *algorithm.Release) string
}

// Releaser orchestrates the 11-phase release pipeline.
type Releaser struct {
	ctx        *Context
	analyzer   commitAnalyzer
	calculator versionCalculator
	generator  notesGenerator

	// publishersOverride, if set, is used instead of publishers.Detect() (test injection).
	publishersOverride []publishers.Publisher
}

// New creates a new Releaser with default implementations.
func New(ctx *Context) *Releaser {
	return &Releaser{
		ctx:        ctx,
		analyzer:   algorithm.NewAnalyzer(),
		calculator: algorithm.NewCalculator(),
		generator:  algorithm.NewGenerator(),
	}
}

// Run executes the full 11-phase release pipeline.
func (r *Releaser) Run() error {
	// Phase 1: Initialize
	branch, err := r.initialize()
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Phase 2: Analyze branch
	branchCfg, err := r.analyzeBranch(branch)
	if err != nil {
		return fmt.Errorf("branch analysis failed: %w", err)
	}

	// Phase 3: Verify auth (skip if no remote configured)
	if !r.ctx.DryRun {
		if err := r.verifyAuth(branch); err != nil {
			r.ctx.Logger.Warn("auth verification skipped", zap.Error(err))
		}
	}

	// Phase 4: Find last release
	lastRelease, err := r.findLastRelease()
	if err != nil {
		return fmt.Errorf("find last release failed: %w", err)
	}

	// Phase 5: Get commits
	commits, err := r.getCommits(lastRelease)
	if err != nil {
		return fmt.Errorf("get commits failed: %w", err)
	}

	// Phase 6: Analyze commits
	releaseType := r.analyzeCommits(commits)
	if releaseType == "" {
		if r.ctx.Verbose {
			r.ctx.Logger.Info("No releasable commits found")
		}
		return nil
	}

	// Phase 7: Calculate version
	nextRelease, err := r.calculateVersion(lastRelease, releaseType, branchCfg)
	if err != nil {
		return fmt.Errorf("version calculation failed: %w", err)
	}

	// Phase 8: Generate notes
	notes := r.generateNotes(commits, lastRelease, nextRelease)
	nextRelease.Notes = notes

	// Phase 9: Create tag
	if !r.ctx.DryRun {
		tagRef := nextRelease.GitHead
		if tagRef == "" {
			tagRef = "HEAD"
		}
		if err := r.ctx.Git.CreateTag(nextRelease.GitTag, tagRef); err != nil {
			return fmt.Errorf("create tag failed: %w", err)
		}
	} else {
		r.ctx.Logger.Info("Dry run: would create tag", zap.String("tag", nextRelease.GitTag))
	}

	// Phase 10: Publish
	if err := r.publish(nextRelease); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	// Phase 11: Notify success
	if err := r.notifySuccess(nextRelease); err != nil {
		return fmt.Errorf("notify success failed: %w", err)
	}

	return nil
}

// initialize verifies we're in a git repo and returns the current branch name.
func (r *Releaser) initialize() (string, error) {
	if !r.ctx.Git.IsGitRepo() {
		return "", fmt.Errorf("not a git repository")
	}

	branch, err := r.ctx.Git.GetCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return branch, nil
}

// analyzeBranch matches the current branch against the config.
func (r *Releaser) analyzeBranch(branch string) (*algorithm.Branch, error) {
	for _, bc := range r.ctx.Config.Branches {
		if !matchesPattern(bc.Name, branch) {
			continue
		}

		b := &algorithm.Branch{
			Name:       bc.Name,
			Prerelease: bc.Prerelease,
			Channel:    bc.Channel,
		}

		switch {
		case bc.Type == string(algorithm.BranchTypeMaintenance):
			b.Type = algorithm.BranchTypeMaintenance
		case bc.Type == string(algorithm.BranchTypePrerelease) || bc.Prerelease != "":
			b.Type = algorithm.BranchTypePrerelease
		default:
			b.Type = algorithm.BranchTypeRelease
		}

		return b, nil
	}

	return nil, fmt.Errorf("branch %q is not configured for release", branch)
}

// matchesPattern checks if a branch name matches a configured pattern.
func matchesPattern(pattern, branch string) bool {
	if pattern == branch {
		return true
	}
	// Support N.x patterns for maintenance branches (e.g., "1.x" matches "1.0", "1.1")
	if strings.HasSuffix(pattern, ".x") {
		prefix := strings.TrimSuffix(pattern, ".x")
		if strings.HasPrefix(branch, prefix+".") {
			return true
		}
	}
	return false
}

// verifyAuth verifies push authentication.
func (r *Releaser) verifyAuth(branch string) error {
	return r.ctx.Git.VerifyAuth("origin", branch)
}

// findLastRelease finds the most recent release tag.
func (r *Releaser) findLastRelease() (*algorithm.Release, error) {
	return r.ctx.Git.GetLastRelease(r.ctx.Config.TagFormat)
}

// getCommits retrieves commits since the last release.
func (r *Releaser) getCommits(lastRelease *algorithm.Release) ([]*algorithm.Commit, error) {
	var from string
	if lastRelease != nil {
		from = lastRelease.GitTag
	}

	head, err := r.ctx.Git.GetHead()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	return r.ctx.Git.GetCommits(from, head)
}

// analyzeCommits determines the release type from commits.
func (r *Releaser) analyzeCommits(commits []*algorithm.Commit) algorithm.ReleaseType {
	return r.analyzer.AnalyzeCommits(commits)
}

// calculateVersion computes the next version and builds a Release object.
func (r *Releaser) calculateVersion(lastRelease *algorithm.Release, releaseType algorithm.ReleaseType, branch *algorithm.Branch) (*algorithm.Release, error) {
	nextVersion, err := r.calculator.CalculateNextVersion(lastRelease, releaseType, branch)
	if err != nil {
		return nil, err
	}

	tagFormat := r.ctx.Config.TagFormat
	if tagFormat == "" {
		tagFormat = "v${version}"
	}

	gitTag := strings.Replace(tagFormat, "${version}", nextVersion, 1)

	head, _ := r.ctx.Git.GetHead()
	if head == "" {
		head = "HEAD"
	}

	return &algorithm.Release{
		Version: nextVersion,
		GitTag:  gitTag,
		GitHead: head,
		Type:    releaseType,
		Branch:  branch.Name,
	}, nil
}

// generateNotes generates release notes from commits.
func (r *Releaser) generateNotes(commits []*algorithm.Commit, lastRelease *algorithm.Release, nextRelease *algorithm.Release) string {
	return r.generator.GenerateNotes(commits, lastRelease, nextRelease)
}

// publish runs all detected publishers.
func (r *Releaser) publish(release *algorithm.Release) error {
	detected := r.publishersOverride
	if detected == nil {
		detected = publishers.Detect()
	}

	if len(detected) == 0 {
		return nil
	}

	for _, pub := range detected {
		if r.ctx.DryRun {
			r.ctx.Logger.Info("Dry run: would publish package", zap.String("publisher", pub.Name()))
			continue
		}

		if err := pub.Prepare(release.Version); err != nil {
			return fmt.Errorf("publisher %q prepare failed: %w", pub.Name(), err)
		}
		if err := pub.Publish(release.Version); err != nil {
			return fmt.Errorf("publisher %q publish failed: %w", pub.Name(), err)
		}
		if err := pub.Verify(release.Version); err != nil {
			return fmt.Errorf("publisher %q verify failed: %w", pub.Name(), err)
		}
	}

	return nil
}

// notifySuccess notifies plugins of a successful release.
func (r *Releaser) notifySuccess(release *algorithm.Release) error {
	return nil
}
