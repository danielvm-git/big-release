// story: e03s01
package plugins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// GitPlugin commits changes and manages git tags for releases.
type GitPlugin struct {
	// Dir is the working directory for git commands; empty means current dir.
	Dir string
}

// NewGitPlugin creates a new GitPlugin.
func NewGitPlugin() *GitPlugin {
	return &GitPlugin{}
}

// Name returns the plugin name.
func (p *GitPlugin) Name() string {
	return "git"
}

// gitCommand creates an exec.Cmd with the plugin's working directory.
func (p *GitPlugin) gitCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	if p.Dir != "" {
		cmd.Dir = p.Dir
	}
	return cmd
}

// VerifyConditions checks that git is installed and the working directory is a git repository.
func (p *GitPlugin) VerifyConditions(ctx *algorithm.Context) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH: %w", err)
	}
	cmd := p.gitCommand("rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}
	if strings.TrimSpace(string(out)) == "" {
		return fmt.Errorf("not a git repository")
	}
	return nil
}

// AnalyzeCommits is not applicable for the git plugin.
func (p *GitPlugin) AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error) {
	return "", nil
}

// GenerateNotes is not applicable for the git plugin.
func (p *GitPlugin) GenerateNotes(ctx *algorithm.Context) (string, error) {
	return "", nil
}

func (p *GitPlugin) stageChanges() error {
	cmd := p.gitCommand("add", ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %w\noutput: %s", err, string(out))
	}
	return nil
}

func (p *GitPlugin) hasChangesToCommit() (bool, error) {
	cmd := p.gitCommand("status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git status failed: %w", err)
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func (p *GitPlugin) commitRelease(version string) error {
	msg := fmt.Sprintf("chore(release): %s [skip ci]\n\nRelease version %s", version, version)
	cmd := p.gitCommand("commit", "-m", msg)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %w\noutput: %s", err, string(out))
	}
	return nil
}

// Prepare stages all changes and commits them with the release version.
func (p *GitPlugin) Prepare(ctx *algorithm.Context) error {
	if ctx.DryRun {
		return nil
	}
	if err := p.stageChanges(); err != nil {
		return err
	}
	hasChanges, err := p.hasChangesToCommit()
	if err != nil {
		return err
	}
	if !hasChanges {
		return nil
	}
	return p.commitRelease(ctx.NextRelease.Version)
}

func (p *GitPlugin) createTag(version string) error {
	cmd := p.gitCommand("tag", "-a", version, "-m", fmt.Sprintf("release %s", version))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git tag failed: %w\noutput: %s", err, string(out))
	}
	return nil
}

func (p *GitPlugin) pushRefs() error {
	pushCmd := p.gitCommand("push")
	if out, err := pushCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %w\noutput: %s", err, string(out))
	}
	pushTagsCmd := p.gitCommand("push", "--tags")
	if out, err := pushTagsCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push --tags failed: %w\noutput: %s", err, string(out))
	}
	return nil
}

// Publish creates a git tag and pushes changes and tags to the remote.
func (p *GitPlugin) Publish(ctx *algorithm.Context) (*algorithm.Release, error) {
	if ctx.DryRun {
		return nil, nil
	}
	if err := p.createTag(ctx.NextRelease.Version); err != nil {
		return nil, err
	}
	if err := p.pushRefs(); err != nil {
		_ = p.deleteTag(ctx.NextRelease.Version)
		return nil, fmt.Errorf("push failed, local tag %s removed: %w", ctx.NextRelease.Version, err)
	}
	return nil, nil
}

func (p *GitPlugin) deleteTag(version string) error {
	cmd := p.gitCommand("tag", "-d", version)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git tag -d failed: %w\noutput: %s", err, string(out))
	}
	return nil
}

// Success is called after a successful release.
func (p *GitPlugin) Success(ctx *algorithm.Context) error {
	return nil
}

// Fail is called on release failure.
func (p *GitPlugin) Fail(ctx *algorithm.Context, err error) error {
	return nil
}

func init() {
	Register(NewGitPlugin())
}
