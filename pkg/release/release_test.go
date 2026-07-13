package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// ── Helpers ──────────────────────────────────────────────────────────────

// runCmdInDir runs a command in the given directory and fails the test on error.
func runCmdInDir(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd %s %v: %v\n%s", name, args, err, out)
	}
}

// writeFile writes content to a file in dir, failing the test on error.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// initReleaseTestRepo creates a minimal git repo in dir with the default
// branch "main", an initial commit, and a .big-release.yml config.
// Returns a git.Client configured for that repo.
func initReleaseTestRepo(t *testing.T, dir string) *git.Client {
	t.Helper()
	runCmdInDir(t, dir, "git", "init", "-b", "main")
	runCmdInDir(t, dir, "git", "config", "user.email", "test@example.com")
	runCmdInDir(t, dir, "git", "config", "user.name", "Test User")
	writeFile(t, dir, "README.md", "# test")
	runCmdInDir(t, dir, "git", "add", ".")
	runCmdInDir(t, dir, "git", "commit", "-m", "chore: initial commit")

	// Write default config with main + beta + maintenance branches
	writeFile(t, dir, ".big-release.yml", `branches:
  - name: main
  - name: beta
    prerelease: beta
  - name: "1.x"
    type: maintenance
tagFormat: v${version}
`)

	client, err := git.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// addCommit creates a commit with the given message in dir.
func addCommit(t *testing.T, dir, message string) {
	t.Helper()
	// Make a change so git has something to commit
	fname := fmt.Sprintf("file-%d.txt", len(message))
	writeFile(t, dir, fname, message)
	runCmdInDir(t, dir, "git", "add", ".")
	runCmdInDir(t, dir, "git", "commit", "-m", message)
}

// addTag creates a lightweight tag in dir.
func addTag(t *testing.T, dir, tag string) {
	t.Helper()
	runCmdInDir(t, dir, "git", "tag", tag)
}

// chdir changes the working directory to dir and returns a function to restore.
func chdir(t *testing.T, dir string) func() {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() { _ = os.Chdir(orig) }
}

// nopLogger returns a no-op zap logger for tests.
func nopLogger() *zap.Logger {
	return zap.NewNop()
}

// ── Mock types for calculator / generator / publisher ────────────────────

type mockCalculator struct {
	nextVersion string
	err         error
}

func (m *mockCalculator) CalculateNextVersion(lastRelease *algorithm.Release, releaseType algorithm.ReleaseType, branch *algorithm.Branch) (string, error) {
	return m.nextVersion, m.err
}

type mockGenerator struct {
	notes string
}

func (m *mockGenerator) GenerateNotes(commits []*algorithm.Commit, lastRelease *algorithm.Release, nextRelease *algorithm.Release) string {
	return m.notes
}

type mockPublisher struct {
	name            string
	detectResult    bool
	prepareErr      error
	publishErr      error
	verifyErr       error
	prepareCalled   bool
	publishCalled   bool
	verifyCalled    bool
	preparedVersion string
}

func (m *mockPublisher) Name() string { return m.name }
func (m *mockPublisher) Detect() bool { return m.detectResult }
func (m *mockPublisher) Prepare(version string) error {
	m.prepareCalled = true
	m.preparedVersion = version
	return m.prepareErr
}
func (m *mockPublisher) Publish(version string) error { m.publishCalled = true; return m.publishErr }
func (m *mockPublisher) Verify(version string) error  { m.verifyCalled = true; return m.verifyErr }

// ── Phase-level tests ────────────────────────────────────────────────────

func TestInitialize(t *testing.T) {
	// SC-e01s10-P0-04: initialize() returns error when not in a git repo
	t.Run("error when not in a git repo", func(t *testing.T) {
		dir := t.TempDir()
		defer chdir(t, dir)()

		client, err := git.NewClient()
		if err != nil {
			t.Fatal(err)
		}

		r := &Releaser{
			ctx: &Context{Git: client},
		}

		branch, err := r.initialize()
		if err == nil {
			t.Error("expected error when not in a git repo")
		}
		if branch != "" {
			t.Errorf("expected empty branch on error, got %q", branch)
		}
		if !strings.Contains(err.Error(), "not a git repository") {
			t.Errorf("expected 'not a git repository' error, got: %v", err)
		}
	})

	t.Run("success in a git repo", func(t *testing.T) {
		dir := t.TempDir()
		gitClient := initReleaseTestRepo(t, dir)
		defer chdir(t, dir)()

		r := &Releaser{
			ctx: &Context{Git: gitClient},
		}

		branch, err := r.initialize()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "main" {
			t.Errorf("expected branch 'main', got %q", branch)
		}
	})
}

func TestAnalyzeBranch(t *testing.T) {
	// SC-e01s10-P0-05: error for unconfigured branch
	t.Run("error for unconfigured branch", func(t *testing.T) {
		cfg := configWithBranches()
		r := &Releaser{
			ctx: &Context{Config: cfg},
		}

		_, err := r.analyzeBranch("feature-x")
		if err == nil {
			t.Error("expected error for unconfigured branch")
		}
		if !strings.Contains(err.Error(), "not configured for release") {
			t.Errorf("expected 'not configured for release' error, got: %v", err)
		}
	})

	// SC-e01s10-P0-07: matches prerelease branch config
	t.Run("matches prerelease branch (beta)", func(t *testing.T) {
		cfg := configWithBranches()
		r := &Releaser{
			ctx: &Context{Config: cfg},
		}

		b, err := r.analyzeBranch("beta")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Type != algorithm.BranchTypePrerelease {
			t.Errorf("expected BranchTypePrerelease, got %q", b.Type)
		}
		if b.Prerelease != "beta" {
			t.Errorf("expected Prerelease='beta', got %q", b.Prerelease)
		}
	})

	// SC-e01s10-P0-08: matches maintenance branch config
	t.Run("matches maintenance branch (1.x)", func(t *testing.T) {
		cfg := configWithBranches()
		r := &Releaser{
			ctx: &Context{Config: cfg},
		}

		// Test that 1.0 matches the 1.x maintenance pattern
		b, err := r.analyzeBranch("1.0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Type != algorithm.BranchTypeMaintenance {
			t.Errorf("expected BranchTypeMaintenance, got %q", b.Type)
		}
	})

	t.Run("matches default main branch", func(t *testing.T) {
		cfg := configWithBranches()
		r := &Releaser{
			ctx: &Context{Config: cfg},
		}

		b, err := r.analyzeBranch("main")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Type != algorithm.BranchTypeRelease {
			t.Errorf("expected BranchTypeRelease, got %q", b.Type)
		}
	})
}

func TestMatchesPattern(t *testing.T) {
	// SC-e01s10-P0-09: matchesPattern exact match
	t.Run("exact match", func(t *testing.T) {
		if !matchesPattern("main", "main") {
			t.Error("expected matchesPattern('main','main') = true")
		}
	})

	t.Run("no match", func(t *testing.T) {
		if matchesPattern("main", "develop") {
			t.Error("expected matchesPattern('main','develop') = false")
		}
	})

	t.Run("N.x maintenance pattern matches", func(t *testing.T) {
		if !matchesPattern("1.x", "1.0") {
			t.Error("expected matchesPattern('1.x','1.0') = true")
		}
		if !matchesPattern("1.x", "1.5") {
			t.Error("expected matchesPattern('1.x','1.5') = true")
		}
	})

	t.Run("N.x pattern no match on different prefix", func(t *testing.T) {
		if matchesPattern("1.x", "2.0") {
			t.Error("expected matchesPattern('1.x','2.0') = false")
		}
	})

	t.Run("N.x pattern no match on short branch", func(t *testing.T) {
		if matchesPattern("1.x", "main") {
			t.Error("expected matchesPattern('1.x','main') = false")
		}
	})
}

func TestFindLastRelease(t *testing.T) {
	// SC-e01s10-P0-10: findLastRelease delegates to gitClient.GetLastRelease
	t.Run("delegates to gitClient with tags", func(t *testing.T) {
		dir := t.TempDir()
		gitClient := initReleaseTestRepo(t, dir)
		addTag(t, dir, "v1.0.0")
		addTag(t, dir, "v1.1.0")
		defer chdir(t, dir)()

		cfg := defaultTestConfig()
		r := &Releaser{
			ctx: &Context{Config: cfg, Git: gitClient},
		}

		release, err := r.findLastRelease()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if release == nil {
			t.Fatal("expected a release, got nil")
		}
		if release.Version != "1.1.0" {
			t.Errorf("expected version 1.1.0, got %q", release.Version)
		}
		if release.GitTag != "v1.1.0" {
			t.Errorf("expected GitTag v1.1.0, got %q", release.GitTag)
		}
	})

	t.Run("returns nil when no tags exist", func(t *testing.T) {
		dir := t.TempDir()
		gitClient := initReleaseTestRepo(t, dir)
		defer chdir(t, dir)()

		cfg := defaultTestConfig()
		r := &Releaser{
			ctx: &Context{Config: cfg, Git: gitClient},
		}

		release, err := r.findLastRelease()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if release != nil {
			t.Errorf("expected nil release when no tags, got %+v", release)
		}
	})
}

func TestCalculateVersion(t *testing.T) {
	// SC-e01s10-P0-11: calculateVersion delegates to Calculator
	t.Run("delegates to mock calculator", func(t *testing.T) {
		dir := t.TempDir()
		gitClient := initReleaseTestRepo(t, dir)
		defer chdir(t, dir)()

		cfg := defaultTestConfig()
		mockCalc := &mockCalculator{nextVersion: "1.1.0"}
		r := &Releaser{
			ctx:        &Context{Config: cfg, Git: gitClient},
			calculator: mockCalc,
		}

		lastRelease := &algorithm.Release{Version: "1.0.0", GitTag: "v1.0.0"}
		branch := &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease}

		nextRelease, err := r.calculateVersion(lastRelease, algorithm.ReleaseTypeMinor, branch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nextRelease.Version != "1.1.0" {
			t.Errorf("expected version 1.1.0, got %q", nextRelease.Version)
		}
		if nextRelease.GitTag != "v1.1.0" {
			t.Errorf("expected tag v1.1.0, got %q", nextRelease.GitTag)
		}
		if nextRelease.Type != algorithm.ReleaseTypeMinor {
			t.Errorf("expected release type minor, got %q", nextRelease.Type)
		}
	})

	t.Run("returns error from calculator", func(t *testing.T) {
		dir := t.TempDir()
		gitClient := initReleaseTestRepo(t, dir)
		defer chdir(t, dir)()

		cfg := defaultTestConfig()
		mockCalc := &mockCalculator{err: fmt.Errorf("calc error")}
		r := &Releaser{
			ctx:        &Context{Config: cfg, Git: gitClient},
			calculator: mockCalc,
		}

		_, err := r.calculateVersion(nil, "", &algorithm.Branch{Name: "main"})
		if err == nil {
			t.Error("expected error from mock calculator")
		}
	})
}

func TestGenerateNotes(t *testing.T) {
	// SC-e01s10-P0-12: generateNotes delegates to Generator
	t.Run("delegates to mock generator", func(t *testing.T) {
		mockGen := &mockGenerator{notes: "## Test Notes\n- feat: some feature"}
		r := &Releaser{
			ctx:       &Context{Config: defaultTestConfig()},
			generator: mockGen,
		}

		notes := r.generateNotes(nil, nil, nil)
		if notes != "## Test Notes\n- feat: some feature" {
			t.Errorf("expected test notes, got %q", notes)
		}
	})
}

func TestPublish(t *testing.T) {
	// SC-e01s10-P0-13: publish delegates to detected publishers
	t.Run("calls publisher methods in order", func(t *testing.T) {
		mockPub := &mockPublisher{name: "mock", detectResult: true}
		r := &Releaser{
			ctx:                &Context{Config: defaultTestConfig(), Logger: nopLogger()},
			publishersOverride: []publishers.Publisher{mockPub},
		}

		rel := &algorithm.Release{Version: "1.0.0"}
		err := r.publish(rel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mockPub.prepareCalled {
			t.Error("expected Prepare to be called")
		}
		if !mockPub.publishCalled {
			t.Error("expected Publish to be called")
		}
		if !mockPub.verifyCalled {
			t.Error("expected Verify to be called")
		}
		if mockPub.preparedVersion != "1.0.0" {
			t.Errorf("expected prepared version 1.0.0, got %q", mockPub.preparedVersion)
		}
	})

	t.Run("dry-run skips publisher methods", func(t *testing.T) {
		mockPub := &mockPublisher{name: "mock", detectResult: true}
		r := &Releaser{
			ctx:                &Context{Config: defaultTestConfig(), DryRun: true, Logger: nopLogger()},
			publishersOverride: []publishers.Publisher{mockPub},
		}

		rel := &algorithm.Release{Version: "1.0.0"}
		err := r.publish(rel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mockPub.prepareCalled {
			t.Error("expected Prepare NOT to be called in dry-run")
		}
		if mockPub.publishCalled {
			t.Error("expected Publish NOT to be called in dry-run")
		}
		if mockPub.verifyCalled {
			t.Error("expected Verify NOT to be called in dry-run")
		}
	})

	t.Run("returns error when publisher has no detect result", func(t *testing.T) {
		r := &Releaser{
			ctx: &Context{Config: defaultTestConfig(), Logger: nopLogger()},
		}
		// No publishers override -> will use publishers.Detect()
		// In a test without package.json, no publishers are detected
		err := r.publish(&algorithm.Release{Version: "1.0.0"})
		if err != nil {
			t.Fatalf("unexpected error when no publishers: %v", err)
		}
	})
}

// ── Integration / E2E tests ──────────────────────────────────────────────

func TestRun_E2E(t *testing.T) {
	// SC-e01s10-P0-01: Full E2E release with feat: commit on main
	dir := t.TempDir()
	gitClient := initReleaseTestRepo(t, dir)

	// Create an initial release tag
	addTag(t, dir, "v1.0.0")

	// Add a feat: commit
	addCommit(t, dir, "feat: add new feature")

	defer chdir(t, dir)()

	cfg := defaultTestConfig()
	r := &Releaser{
		ctx:        &Context{Config: cfg, Git: gitClient, Logger: nopLogger()},
		analyzer:   algorithm.NewAnalyzer(),
		calculator: algorithm.NewCalculator(),
		generator:  algorithm.NewGenerator(),
	}

	if err := r.Run(); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	// Verify tag was created
	tags, err := gitClient.GetTags()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tag := range tags {
		if tag == "v1.1.0" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected tag v1.1.0 to be created, got tags: %v", tags)
	}
}

func TestRun_DryRun(t *testing.T) {
	// SC-e01s10-P0-02: --dry-run skips tag creation
	// SC-e01s10-P0-03: --dry-run skips publisher execution
	dir := t.TempDir()
	gitClient := initReleaseTestRepo(t, dir)

	addTag(t, dir, "v1.0.0")
	addCommit(t, dir, "feat: dry run feature")

	defer chdir(t, dir)()

	cfg := defaultTestConfig()
	r := &Releaser{
		ctx:        &Context{Config: cfg, Git: gitClient, Logger: nopLogger(), DryRun: true},
		analyzer:   algorithm.NewAnalyzer(),
		calculator: algorithm.NewCalculator(),
		generator:  algorithm.NewGenerator(),
	}

	if err := r.Run(); err != nil {
		t.Fatalf("Run() in dry-run mode failed: %v", err)
	}

	// Verify no tags were created beyond the initial v1.0.0
	tags, err := gitClient.GetTags()
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range tags {
		if tag != "v1.0.0" {
			t.Errorf("unexpected tag created in dry-run mode: %s", tag)
		}
	}
}

func TestRun_NoReleasableCommits(t *testing.T) {
	// SC-e01s10-P0-06: No releasable commits → releaseType="", Run exits 0
	dir := t.TempDir()
	gitClient := initReleaseTestRepo(t, dir)

	addTag(t, dir, "v1.0.0")
	// Only chore: commits, no feat/fix/perf
	addCommit(t, dir, "chore: update deps")
	addCommit(t, dir, "docs: update readme")

	defer chdir(t, dir)()

	cfg := defaultTestConfig()
	r := &Releaser{
		ctx:        &Context{Config: cfg, Git: gitClient, Logger: nopLogger()},
		analyzer:   algorithm.NewAnalyzer(),
		calculator: algorithm.NewCalculator(),
		generator:  algorithm.NewGenerator(),
	}

	// Should not error — just logs that no releasable commits were found
	if err := r.Run(); err != nil {
		t.Fatalf("Run() with no releasable commits should not error: %v", err)
	}

	// Verify no new tags were created
	tags, err := gitClient.GetTags()
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range tags {
		if tag != "v1.0.0" {
			t.Errorf("unexpected tag created with no releasable commits: %s", tag)
		}
	}
}

// ── Config helpers ────────────────────────────────────────────────────────

func configWithBranches() *algorithm.Config {
	return &algorithm.Config{
		Branches: []algorithm.BranchConfig{
			{Name: "main"},
			{Name: "beta", Prerelease: "beta"},
			{Name: "1.x", Type: string(algorithm.BranchTypeMaintenance)},
		},
		TagFormat: "v${version}",
	}
}

func defaultTestConfig() *algorithm.Config {
	return configWithBranches()
}
