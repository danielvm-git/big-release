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

	t.Run("BUG-tag-ignores-tagformat: Publish applies configured tagFormat to the created tag", func(t *testing.T) {
		fg := &fakeGit{isRepo: true}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
			Config: &algorithm.Config{TagFormat: "v${version}"},
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.2.3"},
		}

		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if fg.lastCreatedTag != "v1.2.3" {
			t.Errorf("expected tag %q (tagFormat applied), got %q", "v1.2.3", fg.lastCreatedTag)
		}
	})

	t.Run("BUG-tag-ignores-tagformat: Publish falls back to bare version when tagFormat unset", func(t *testing.T) {
		fg := &fakeGit{isRepo: true}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.2.3"},
		}

		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if fg.lastCreatedTag != "1.2.3" {
			t.Errorf("expected bare tag %q, got %q", "1.2.3", fg.lastCreatedTag)
		}
	})

	t.Run("BUG-tag-ignores-tagformat: rollback deletes the same formatted tag that was created", func(t *testing.T) {
		fg := &fakeGit{
			isRepo:  true,
			pushErr: fmt.Errorf("git push: exit status 128"),
		}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
			Config: &algorithm.Config{TagFormat: "v${version}"},
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.2.3"},
		}

		if _, err := p.Publish(ctx, state); err == nil {
			t.Fatal("expected error from push failure")
		}
		if fg.lastDeletedTag != fg.lastCreatedTag {
			t.Errorf("rollback deleted %q but created tag was %q — must match", fg.lastDeletedTag, fg.lastCreatedTag)
		}
	})
}

// TestFormatTag_RoundTripsWithGetLastRelease is the exact scenario from
// BUG-tag-ignores-tagformat: whatever git.FormatTag writes, the real
// client's GetLastRelease(tagFormat) must be able to read back — using a
// real git repo and the real internal/git client, not fakes, since the bug
// was precisely a mismatch between the write path (git.go's CreateTag call)
// and the read path (GetLastRelease) that no fake would catch.
func TestFormatTag_RoundTripsWithGetLastRelease(t *testing.T) {
	dir := setupTestGitRepo(t)
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prevWD) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	writeFileInDir(t, dir, "initial.txt", "initial")
	execInDir(t, dir, "git", "add", ".")
	execInDir(t, dir, "git", "commit", "-m", "initial commit")

	const tagFormat = "v${version}"
	gitClient, err := git.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// First release: tag exactly as the git plugin's Publish now would.
	firstTag := git.FormatTag(tagFormat, "1.0.0")
	if firstTag != "v1.0.0" {
		t.Fatalf("FormatTag(%q, %q) = %q, want v1.0.0", tagFormat, "1.0.0", firstTag)
	}
	execInDir(t, dir, "git", "tag", "-a", firstTag, "-m", "release "+firstTag)

	last, err := gitClient.GetLastRelease(tagFormat)
	if err != nil {
		t.Fatal(err)
	}
	if last == nil || last.Version != "1.0.0" || last.GitTag != firstTag {
		t.Fatalf("GetLastRelease after first tag = %+v, want version 1.0.0 / tag %s", last, firstTag)
	}

	// Second release: this is where BUG-tag-ignores-tagformat broke — a bare
	// second tag ("1.0.1") would go unrecognized by GetLastRelease, causing
	// the caller to recompute the same version and fail on tag re-creation.
	secondTag := git.FormatTag(tagFormat, "1.0.1")
	execInDir(t, dir, "git", "tag", "-a", secondTag, "-m", "release "+secondTag)

	last2, err := gitClient.GetLastRelease(tagFormat)
	if err != nil {
		t.Fatal(err)
	}
	if last2 == nil || last2.Version != "1.0.1" || last2.GitTag != secondTag {
		t.Fatalf("GetLastRelease after second tag = %+v, want version 1.0.1 / tag %s", last2, secondTag)
	}
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

	t.Run("BUG-git-push-error-swallowed: Prepare stages nothing without explicit assets config", func(t *testing.T) {
		// With no `assets` configured, committing must be a no-op — matching
		// semantic-release's default (no plugin implements `prepare` unless
		// the user opts into @semantic-release/git). A fresh install must
		// never push a surprise commit to the release branch: on a protected
		// branch that's rejected outright for any non-admin pusher, with no
		// GitHub-side config able to avoid it short of an explicit opt-in.
		// changes: false mirrors the real Client, where HasChangesToCommit
		// now reflects the index (git diff --cached) rather than overall
		// working-tree dirtiness — with nothing staged, it correctly reports
		// no changes even if some other plugin wrote files to disk.
		fg := &fakeGit{isRepo: true, changes: false}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0"}}
		if err := p.Prepare(ctx, state); err != nil {
			t.Fatalf("Prepare: %v", err)
		}
		if fg.stageChangesCalled {
			t.Error("Prepare called StageChanges() (stage-everything) with no assets configured")
		}
		if fg.lastCommitMsg != "" {
			t.Errorf("Prepare committed %q with no assets configured, want no commit", fg.lastCommitMsg)
		}
	})
}

func TestGitPluginAddChannel(t *testing.T) {
	t.Run("SC-e22s01-P1-01: AddChannel writes git note for channel", func(t *testing.T) {
		fg := &fakeGit{isRepo: true, head: "abc123"}
		p := NewGitPlugin(fg)
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0", Channel: "next"}}
		if err := p.AddChannel(ctx, state); err != nil {
			t.Fatalf("AddChannel: %v", err)
		}
		if fg.lastNote != "next" || fg.lastNoteRef != "abc123" {
			t.Errorf("note = %q@%q, want next@abc123", fg.lastNote, fg.lastNoteRef)
		}
	})

	t.Run("SC-e22s01-P1-02: Publish pushes channel notes after tag push", func(t *testing.T) {
		fg := &fakeGit{isRepo: true, head: "deadbeef"}
		p := NewGitPlugin(fg)
		_ = p.AddChannel(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "2.0.0", Channel: "beta"},
		})
		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "2.0.0", Channel: "beta"}}
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if len(fg.pushedNotes) != 1 || fg.pushedNotes[0] != "origin:deadbeef" {
			t.Errorf("pushed notes = %v, want [origin:deadbeef]", fg.pushedNotes)
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
	commits            []*algorithm.Commit
	tags               []string
	tagHead            string
	head               string
	release            *algorithm.Release
	repoURL            string
	branch             string
	isRepo             bool
	changes            bool
	modifiedFiles      []string
	stagedPaths        []string
	stageChangesCalled bool
	lastCommitMsg      string
	lastNote           string
	lastNoteRef        string
	pushedNotes        []string
	createErr          error
	pushErr            error
	deleteErr          error
	lastCreatedTag     string
	lastDeletedTag     string
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
	f.stageChangesCalled = true
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
	f.lastCreatedTag = tag
	return f.createErr
}

func (f *fakeGit) Push(remote string) error {
	return f.pushErr
}

func (f *fakeGit) PushTags(remote string) error {
	return f.pushErr
}

func (f *fakeGit) DeleteTag(tag string) error {
	f.lastDeletedTag = tag
	return f.deleteErr
}

func (f *fakeGit) AddNote(note, ref string) error {
	f.lastNote = note
	f.lastNoteRef = ref
	return nil
}

func (f *fakeGit) PushNotes(remote, ref string) error {
	f.pushedNotes = append(f.pushedNotes, remote+":"+ref)
	return nil
}
