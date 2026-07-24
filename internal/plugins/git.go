// story: e03s01 e21s03 e21s04 e22s01
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
)

// GitPlugin commits changes and manages git tags for releases.
type GitPlugin struct {
	Git            git.GitAPI
	commitMessage  string
	commitAssets   []string
	assetsExplicit bool
	channelNoteRef string
}

// NewGitPlugin creates a new GitPlugin.
func NewGitPlugin(gitAPI git.GitAPI) *GitPlugin {
	return &GitPlugin{Git: gitAPI}
}

// Configure decodes git plugin config from PluginConfigs["git"].
func (p *GitPlugin) Configure(raw map[string]interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("git plugin: failed to re-marshal config: %w", err)
	}
	var cfg GitConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("git plugin: invalid config: %w", err)
	}
	p.commitMessage = cfg.Message
	if _, ok := raw["assets"]; ok {
		p.assetsExplicit = true
		p.commitAssets = cfg.Assets
	}
	return nil
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

func (p *GitPlugin) stageChanges(ctx *algorithm.ReadOnlyContext) error {
	if !p.assetsExplicit {
		return p.Git.StageChanges()
	}

	modified, err := p.Git.GetModifiedFiles()
	if err != nil {
		return err
	}
	paths := matchModifiedAssets(modified, p.commitAssets)
	if len(paths) == 0 {
		return nil
	}
	return p.Git.StagePaths(paths)
}

func (p *GitPlugin) hasChangesToCommit() (bool, error) {
	return p.Git.HasChangesToCommit()
}

type gitTemplateContext struct {
	Version     string
	Branch      string
	Notes       string
	NextRelease *algorithm.Release
	LastRelease *algorithm.Release
}

func (p *GitPlugin) buildCommitMessage(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	version := state.NextRelease.Version
	fallback := fmt.Sprintf("chore(release): %s [skip ci]\n\nRelease version %s", version, version)
	if p.commitMessage == "" {
		return fallback, nil
	}

	tctx := gitTemplateContext{
		Version:     version,
		Notes:       state.NextRelease.Notes,
		NextRelease: state.NextRelease,
		LastRelease: state.LastRelease,
	}
	if ctx.Branch != nil {
		tctx.Branch = ctx.Branch.Name
	}

	t, err := template.New("commit message").Parse(p.commitMessage)
	if err != nil {
		return "", fmt.Errorf("invalid commit message template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, tctx); err != nil {
		return "", fmt.Errorf("failed to render commit message template: %w", err)
	}
	return buf.String(), nil
}

func (p *GitPlugin) commitRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	msg, err := p.buildCommitMessage(ctx, state)
	if err != nil {
		return err
	}
	return p.Git.Commit(msg)
}

// Prepare stages configured changes and commits them with the release version.
func (p *GitPlugin) Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun {
		return nil
	}
	if err := p.stageChanges(ctx); err != nil {
		return err
	}
	hasChanges, err := p.hasChangesToCommit()
	if err != nil {
		return err
	}
	if !hasChanges {
		return nil
	}
	return p.commitRelease(ctx, state)
}

// AddChannel records the release channel in a git note on HEAD (e22s01).
func (p *GitPlugin) AddChannel(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun || state.NextRelease == nil || state.NextRelease.Channel == "" {
		return nil
	}
	head, err := p.Git.GetHead()
	if err != nil {
		return fmt.Errorf("git plugin: failed to get HEAD for channel note: %w", err)
	}
	if err := p.Git.AddNote(state.NextRelease.Channel, head); err != nil {
		return fmt.Errorf("git plugin: failed to add channel note: %w", err)
	}
	p.channelNoteRef = head
	return nil
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
	if p.channelNoteRef != "" {
		if err := p.Git.PushNotes("origin", p.channelNoteRef); err != nil {
			return nil, fmt.Errorf("git plugin: failed to push channel notes: %w", err)
		}
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

// matchModifiedAssets returns modified paths matching any configured asset glob.
func matchModifiedAssets(modified, patterns []string) []string {
	if len(patterns) == 0 || len(modified) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		candidates := matches
		if len(matches) == 0 {
			candidates = []string{pattern}
		}
		for _, mod := range modified {
			for _, candidate := range candidates {
				if pathMatchesAsset(mod, candidate, pattern) && !seen[mod] {
					seen[mod] = true
					out = append(out, mod)
				}
			}
		}
	}
	return out
}

func pathMatchesAsset(modified, candidate, pattern string) bool {
	if modified == candidate || modified == pattern {
		return true
	}
	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, modified)
		if matched {
			return true
		}
		base := filepath.Base(pattern)
		if strings.Contains(base, "*") {
			matched, _ := filepath.Match(base, filepath.Base(modified))
			return matched
		}
	}
	return filepath.Base(modified) == filepath.Base(candidate)
}
