package release

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/publishers"
)

// stubGit controls the four git reads buildAlgoContext performs. Embedding the
// interface (left nil) makes any unexpected method call panic loudly instead of
// silently returning a zero value.
type stubGit struct {
	git.GitAPI
	branch      string
	lastRelease *algorithm.Release
	commits     []*algorithm.Commit
	repoURL     string
}

func (g *stubGit) GetCurrentBranch() (string, error)                 { return g.branch, nil }
func (g *stubGit) GetLastRelease(string) (*algorithm.Release, error) { return g.lastRelease, nil }
func (g *stubGit) GetCommits(string, string) ([]*algorithm.Commit, error) {
	return g.commits, nil
}
func (g *stubGit) GetRepositoryURL() (string, error) { return g.repoURL, nil }
func (g *stubGit) IsGitRepo() bool                   { return true }

// recordingPublisher records whether the orchestrator invoked it.
type recordingPublisher struct {
	name string
	mu   sync.Mutex
	ran  bool
}

func (p *recordingPublisher) Name() string { return p.name }
func (p *recordingPublisher) Detect() bool { return true }
func (p *recordingPublisher) Prepare(string) error {
	p.mu.Lock()
	p.ran = true
	p.mu.Unlock()
	return nil
}
func (p *recordingPublisher) Publish(string) error { return nil }
func (p *recordingPublisher) Verify(string) error  { return nil }
func (p *recordingPublisher) SetDryRun(bool)       {}

func (p *recordingPublisher) didRun() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ran
}

// disabledExcept builds a publisher config that disables every registered
// publisher except keep, so a test cannot accidentally run a real publisher.
func disabledExcept(names []string, keep string) map[string]algorithm.PublisherConfig {
	out := make(map[string]algorithm.PublisherConfig, len(names))
	for _, n := range names {
		out[n] = algorithm.PublisherConfig{Enabled: n == keep}
	}
	return out
}

// clearPREnv removes every environment marker detectPR looks at, so Run is not
// short-circuited by an inherited CI environment.
func clearPREnv(t *testing.T) {
	t.Helper()
	t.Setenv("GITHUB_EVENT_NAME", "")
	t.Setenv("CI_MERGE_REQUEST_ID", "")
	t.Setenv("BUILD_REASON", "")
	t.Setenv("CIRCLE_PULL_REQUEST", "")
	t.Setenv("BITBUCKET_PR_ID", "")
	t.Setenv("CHANGE_ID", "")
}

func TestRun_NoRelevantChanges_SkipsPublishers(t *testing.T) {
	// A docs-only commit yields neither a release type nor notes. The run must
	// stop before publishers: otherwise they execute with an empty version,
	// leaving the release in a non-terminal, half-published state.
	t.Setenv("CI", "true")
	clearPREnv(t)

	rec := &recordingPublisher{name: "test-recorder-no-release"}
	publishers.Register(rec)

	cfg := &algorithm.Config{
		Branches:   []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:  "v${version}",
		Plugins:    []string{},
		Publishers: disabledExcept(publishers.List(), rec.Name()),
	}
	g := &stubGit{
		branch:  "main",
		commits: []*algorithm.Commit{{Hash: "abcdef1234567890", Message: "docs: clarify README"}},
		repoURL: "https://example.com/repo",
	}

	r := New(&Context{Config: cfg, Git: g, Logger: zap.NewNop(), DryRun: false})
	if err := r.Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if rec.didRun() {
		t.Fatal("publishers must not run when no release version was computed")
	}
}

func TestHasPublishableRelease(t *testing.T) {
	cases := []struct {
		name  string
		state *algorithm.MutableState
		want  bool
	}{
		{"nil state", nil, false},
		{"nil next release", &algorithm.MutableState{}, false},
		{"notes without version", &algorithm.MutableState{NextRelease: &algorithm.Release{Notes: "n"}}, false},
		{"versioned release", &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.2.3"}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasPublishableRelease(tc.state); got != tc.want {
				t.Errorf("hasPublishableRelease() = %v, want %v", got, tc.want)
			}
		})
	}
}

// barrierGit blocks each independent read until all three have started. A
// sequential implementation never gets past the first read, so the test times
// out instead of passing.
type barrierGit struct {
	git.GitAPI
	started chan struct{}
	proceed chan struct{}
}

func (g *barrierGit) wait() { g.started <- struct{}{}; <-g.proceed }

func (g *barrierGit) GetCurrentBranch() (string, error) { g.wait(); return "main", nil }
func (g *barrierGit) GetLastRelease(string) (*algorithm.Release, error) {
	g.wait()
	return nil, nil
}
func (g *barrierGit) GetRepositoryURL() (string, error) {
	g.wait()
	return "https://example.com/repo", nil
}
func (g *barrierGit) GetCommits(string, string) ([]*algorithm.Commit, error) { return nil, nil }

func TestBuildAlgoContext_RunsIndependentGitReadsConcurrently(t *testing.T) {
	g := &barrierGit{started: make(chan struct{}, 3), proceed: make(chan struct{})}
	var once sync.Once
	release := func() { once.Do(func() { close(g.proceed) }) }
	defer release()

	r := New(&Context{
		Config: &algorithm.Config{
			Branches:  []algorithm.BranchConfig{{Name: "main"}},
			TagFormat: "v${version}",
		},
		Git:    g,
		Logger: zap.NewNop(),
		DryRun: true,
	})

	done := make(chan error, 1)
	go func() {
		_, _, err := r.buildAlgoContext()
		done <- err
	}()

	for i := 0; i < 3; i++ {
		select {
		case <-g.started:
		case <-time.After(2 * time.Second):
			t.Fatal("independent git reads did not start concurrently")
		}
	}
	release()
	if err := <-done; err != nil {
		t.Fatalf("buildAlgoContext failed: %v", err)
	}
}
