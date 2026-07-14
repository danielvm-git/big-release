// story: e03s01
package plugins

import (
	"fmt"
	"os/exec"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
)

// GitPlugin commits changes and manages git tags for releases.
type GitPlugin struct {
	// Git is the git API implementation.
	Git git.GitAPI
}

// NewGitPlugin creates a new GitPlugin.
func NewGitPlugin(gitAPI git.GitAPI) *GitPlugin {
	return &GitPlugin{Git: gitAPI}
}

// Name returns the plugin name.
func (p *GitPlugin) Name() string {
	return "git"
}

// VerifyConditions checks that git is installed and the working directory is a git repository.
func (p *GitPlugin) VerifyConditions(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH: %w", err)
	}
	if !p.Git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}
	return nil
}

// AnalyzeCommits is not applicable for the git plugin.
func (p *GitPlugin) AnalyzeCommits(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return "", nil
}

// VerifyRelease is not applicable for the git plugin.
func (p *GitPlugin) VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// GenerateNotes is not applicable for the git plugin.
func (p *GitPlugin) GenerateNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	return "", nil
}

func (p *GitPlugin) stageChanges() error {
	return p.Git.StageChanges()
}

func (p *GitPlugin) hasChangesToCommit() (bool, error) {
	return p.Git.HasChangesToCommit()
}

func (p *GitPlugin) commitRelease(version string) error {
	msg := fmt.Sprintf("chore(release): %s [skip ci]\n\nRelease version %s", version, version)
	return p.Git.Commit(msg)
}

// Prepare stages all changes and commits them with the release version.
func (p *GitPlugin) Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
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
	return p.commitRelease(state.NextRelease.Version)
}

func (p *GitPlugin) createTag(version string) error {
	return p.Git.CreateTag(version, fmt.Sprintf("release %s", version))
}

func (p *GitPlugin) pushRefs() error {
	if err := p.Git.Push("origin"); err != nil {
		return err
	}
	return p.Git.PushTags("origin")
}

// Publish creates a git tag and pushes changes and tags to the remote.
func (p *GitPlugin) Publish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error) {
	if ctx.DryRun {
		return nil, nil
	}
	if err := p.createTag(state.NextRelease.Version); err != nil {
		return nil, err
	}
	if err := p.pushRefs(); err != nil {
		_ = p.deleteTag(state.NextRelease.Version)
		return nil, fmt.Errorf("push failed, local tag %s removed: %w", state.NextRelease.Version, err)
	}
	return nil, nil
}

func (p *GitPlugin) deleteTag(version string) error {
	return p.Git.DeleteTag(version)
}

// Success is called after a successful release.
func (p *GitPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// Fail is called on release failure.
func (p *GitPlugin) Fail(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, err error) error {
	return nil
}
