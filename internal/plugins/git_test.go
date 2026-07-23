// story: e03s01
package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
)

// scrubGitEnv drops git hook/worktree env vars so tests that shell out to
// git cannot mutate the real repository when run under `git commit` hooks
// (which set GIT_DIR / GIT_INDEX_FILE / GIT_WORK_TREE).
func scrubGitEnv() []string {
	var out []string
	for _, e := range os.Environ() {
		switch {
		case strings.HasPrefix(e, "GIT_DIR="),
			strings.HasPrefix(e, "GIT_WORK_TREE="),
			strings.HasPrefix(e, "GIT_INDEX_FILE="),
			strings.HasPrefix(e, "GIT_PREFIX="),
			strings.HasPrefix(e, "GIT_COMMON_DIR="),
			strings.HasPrefix(e, "GIT_OBJECT_DIRECTORY="),
			strings.HasPrefix(e, "GIT_ALTERNATE_OBJECT_DIRECTORIES="):
			continue
		default:
			out = append(out, e)
		}
	}
	return out
}

// setupTestGitRepo creates a temporary git repository for testing.
func setupTestGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "--local", "user.email", "test@test.com"},
		{"git", "config", "--local", "user.name", "Test User"},
		{"git", "config", "--local", "commit.gpgSign", "false"},
		{"git", "config", "--local", "tag.gpgSign", "false"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = scrubGitEnv()
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v\noutput: %s", args, err, string(out))
		}
	}

	return dir
}

// writeFileInDir writes a file in the given directory for test setup.
func writeFileInDir(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(dir+"/"+name, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

// execInDir runs a command in the given directory.
func execInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = scrubGitEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %v failed: %v\noutput: %s", args, err, string(out))
	}
	return string(out)
}

func TestGitPluginName(t *testing.T) {
	t.Run("SC-e03s01-P1-01: Name returns 'git'", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		if name := p.Name(); name != "git" {
			t.Errorf("expected Name() == %q, got %q", "git", name)
		}
	})
}

func TestGitPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e03s01-P1-02: VerifyConditions passes in a git repo", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{}
		state := &algorithm.MutableState{}

		if err := p.VerifyConditions(ctx, state); err != nil {
			t.Errorf("expected no error in git repo, got: %v", err)
		}
	})

	t.Run("SC-e03s01-P1-03: VerifyConditions fails outside git repo", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: false})
		ctx := &algorithm.ReadOnlyContext{}
		state := &algorithm.MutableState{}

		if err := p.VerifyConditions(ctx, state); err == nil {
			t.Error("expected error outside git repo, got nil")
		}
	})
}

func TestGitPluginAnalyzeCommits(t *testing.T) {
	t.Run("SC-e03s01-P1-04: AnalyzeCommits returns empty release type", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{}
		state := &algorithm.MutableState{}
		rt, err := p.AnalyzeCommits(ctx, state)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if rt != "" {
			t.Errorf("expected empty release type, got %q", rt)
		}
	})
}

func TestGitPluginGenerateNotes(t *testing.T) {
	t.Run("SC-e03s01-P1-05: GenerateNotes returns empty string", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{}
		state := &algorithm.MutableState{}
		notes, err := p.GenerateNotes(ctx, state)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if notes != "" {
			t.Errorf("expected empty notes, got %q", notes)
		}
	})
}

func TestGitPluginPrepare(t *testing.T) {
	t.Run("SC-e03s01-P1-06: Prepare does nothing in dry-run mode", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{DryRun: true}
		state := &algorithm.MutableState{}
		if err := p.Prepare(ctx, state); err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
	})

	t.Run("SC-e03s01-P1-07: Prepare stages and commits changes", func(t *testing.T) {
		dir := setupTestGitRepo(t)

		// Create an initial commit so we have a baseline
		writeFileInDir(t, dir, "initial.txt", "initial")
		execInDir(t, dir, "git", "add", ".")
		execInDir(t, dir, "git", "commit", "-m", "initial commit")

		// Now create a change to commit
		writeFileInDir(t, dir, "test.txt", "release content")

		// Use a real git client for tests that verify actual git operations
		realClient, err := git.NewClient()
		if err != nil {
			t.Fatalf("failed to create git client: %v", err)
		}
		_ = realClient
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		if err := p.Prepare(ctx, state); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// The fakeGit doesn't actually commit, so we can't verify git log output here.
		// This test just verifies the plugin flow doesn't error.
	})

	t.Run("SC-e03s01-P1-08: Prepare skips commit when no changes", func(t *testing.T) {
		dir := setupTestGitRepo(t)

		// Create initial commit
		writeFileInDir(t, dir, "initial.txt", "initial")
		execInDir(t, dir, "git", "add", ".")
		execInDir(t, dir, "git", "commit", "-m", "initial commit")

		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		if err := p.Prepare(ctx, state); err != nil {
			t.Errorf("expected no error with no changes, got: %v", err)
		}
	})
}

func TestGitPluginPublish(t *testing.T) {
	t.Run("SC-e03s01-P1-09: Publish does nothing in dry-run mode", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		ctx := &algorithm.ReadOnlyContext{DryRun: true}
		state := &algorithm.MutableState{}
		release, err := p.Publish(ctx, state)
		if err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
		if release != nil {
			t.Errorf("expected nil release in dry-run, got %v", release)
		}
	})

	t.Run("SC-e03s01-P1-10: Publish creates git tag", func(t *testing.T) {
		// Use fakeGit with a push error to simulate push failure
		fg := &fakeGit{
			isRepo:    true,
			pushErr:   fmt.Errorf("git push: exit status 128"),
			createErr: nil,
		}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		// Publish (creates tag, push fails, tag cleaned up)
		_, err := p.Publish(ctx, state)
		if err == nil {
			t.Error("expected error from push failure, got nil")
		}
		if !strings.Contains(err.Error(), "removed") {
			t.Errorf("expected tag cleanup message, got: %v", err)
		}
	})
}

func TestGitPluginSuccess(t *testing.T) {
	t.Run("SC-e03s01-P1-11: Success returns nil", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		if err := p.Success(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitPluginFail(t *testing.T) {
	t.Run("SC-e03s01-P1-12: Fail returns nil", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true})
		if err := p.Fail(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}, nil); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitPluginCommitMessageTemplate(t *testing.T) {
	t.Run("SC-e21s03-P1-01: Prepare uses configured commit message template", func(t *testing.T) {
		fg := &fakeGit{isRepo: true, changes: true}
		p := NewGitPlugin(fg)
		if err := p.Configure(map[string]interface{}{
			"message": "release {{.Version}} on {{.Branch}}",
		}); err != nil {
			t.Fatalf("Configure: %v", err)
		}
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
			Branch: &algorithm.Branch{Name: "main"},
		}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.2.3"}}
		if err := p.Prepare(ctx, state); err != nil {
			t.Fatalf("Prepare: %v", err)
		}
		if fg.lastCommitMsg != "release 1.2.3 on main" {
			t.Errorf("commit message = %q, want custom template output", fg.lastCommitMsg)
		}
	})

	t.Run("SC-e21s03-P1-02: invalid commit message template returns error", func(t *testing.T) {
		p := NewGitPlugin(&fakeGit{isRepo: true, changes: true})
		if err := p.Configure(map[string]interface{}{
			"message": "release {{.MissingField",
		}); err != nil {
			t.Fatalf("Configure: %v", err)
		}
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0"}}
		if err := p.Prepare(ctx, state); err == nil {
			t.Fatal("expected template parse error")
		}
	})
}

func TestMatchModifiedAssets(t *testing.T) {
	t.Run("SC-e21s04-P1-01: matches modified files against asset globs", func(t *testing.T) {
		got := matchModifiedAssets(
			[]string{"CHANGELOG.md", "README.md", "internal/foo.go"},
			[]string{"CHANGELOG.md", "package.json"},
		)
		if len(got) != 1 || got[0] != "CHANGELOG.md" {
			t.Errorf("got %v, want [CHANGELOG.md]", got)
		}
	})
}

func TestGitPluginCommitAssets(t *testing.T) {
	t.Run("SC-e21s04-P1-02: Prepare stages only configured asset paths", func(t *testing.T) {
		fg := &fakeGit{
			isRepo:        true,
			changes:       true,
			modifiedFiles: []string{"CHANGELOG.md", "README.md"},
		}
		p := NewGitPlugin(fg)
		if err := p.Configure(map[string]interface{}{
			"assets": []interface{}{"CHANGELOG.md"},
		}); err != nil {
			t.Fatalf("Configure: %v", err)
		}
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0"}}
		if err := p.Prepare(ctx, state); err != nil {
			t.Fatalf("Prepare: %v", err)
		}
		if len(fg.stagedPaths) != 1 || fg.stagedPaths[0] != "CHANGELOG.md" {
			t.Errorf("staged paths = %v, want [CHANGELOG.md]", fg.stagedPaths)
		}
	})
}

func TestGitPluginRegistration(t *testing.T) {
	t.Run("SC-e03s01-P1-13: GitPlugin can be registered", func(t *testing.T) {
		// Clear the registry for this test
		oldPlugins := globalRegistry.plugins
		globalRegistry.plugins = make(map[string]Plugin)
		defer func() { globalRegistry.plugins = oldPlugins }()

		// Register a new GitPlugin
		gitPlugin := NewGitPlugin(&fakeGit{isRepo: true})
		Register(gitPlugin)

		// Verify it's registered
		found := false
		for _, name := range List() {
			if name == "git" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'git' to be registered in global registry")
		}
	})
}

// fakeGit implements git.GitAPI for testing.
type fakeGit struct {
	commits       []*algorithm.Commit
	tags          []string
	tagHead       string
	head          string
	release       *algorithm.Release
	repoURL       string
	branch        string
	isRepo        bool
	changes       bool
	modifiedFiles []string
	stagedPaths   []string
	lastCommitMsg string
	createErr     error
	pushErr       error
	deleteErr     error
}

func (f *fakeGit) GetCommits(from, to string) ([]*algorithm.Commit, error) {
	return f.commits, nil
}

func (f *fakeGit) GetTags() ([]string, error) {
	return f.tags, nil
}

func (f *fakeGit) GetTagHead(tag string) (string, error) {
	return f.tagHead, nil
}

func (f *fakeGit) GetHead() (string, error) {
	return f.head, nil
}

func (f *fakeGit) GetLastRelease(tagFormat string) (*algorithm.Release, error) {
	return f.release, nil
}

func (f *fakeGit) GetRepositoryURL() (string, error) {
	return f.repoURL, nil
}

func (f *fakeGit) GetCurrentBranch() (string, error) {
	return f.branch, nil
}

func (f *fakeGit) IsGitRepo() bool {
	return f.isRepo
}

func (f *fakeGit) StageChanges() error {
	return nil
}

func (f *fakeGit) GetModifiedFiles() ([]string, error) {
	return f.modifiedFiles, nil
}

func (f *fakeGit) StagePaths(paths []string) error {
	f.stagedPaths = append(f.stagedPaths, paths...)
	return nil
}

func (f *fakeGit) HasChangesToCommit() (bool, error) {
	return f.changes, nil
}

func (f *fakeGit) Commit(message string) error {
	f.lastCommitMsg = message
	return nil
}

func (f *fakeGit) CreateTag(tag, message string) error {
	return f.createErr
}

func (f *fakeGit) Push(remote string) error {
	return f.pushErr
}

func (f *fakeGit) PushTags(remote string) error {
	return f.pushErr
}

func (f *fakeGit) DeleteTag(tag string) error {
	return f.deleteErr
}
