// story: NT-3 NT-7
package release

import (
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/plugins"
)

func TestIsLegalReleaseTransition(t *testing.T) {
	cases := []struct {
		name string
		from algorithm.ReleaseStatus
		to   algorithm.ReleaseStatus
		want bool
	}{
		{"pending to analyzing", algorithm.ReleaseStatusPending, algorithm.ReleaseStatusAnalyzing, true},
		{"pending to failed", algorithm.ReleaseStatusPending, algorithm.ReleaseStatusFailed, true},
		{"zero value to analyzing", "", algorithm.ReleaseStatusAnalyzing, true},
		{"analyzing to no_release", algorithm.ReleaseStatusAnalyzing, algorithm.ReleaseStatusNoRelease, true},
		{"analyzing to versioned", algorithm.ReleaseStatusAnalyzing, algorithm.ReleaseStatusVersioned, true},
		{"analyzing to failed", algorithm.ReleaseStatusAnalyzing, algorithm.ReleaseStatusFailed, true},
		{"versioned to prepared", algorithm.ReleaseStatusVersioned, algorithm.ReleaseStatusPrepared, true},
		{"versioned to failed", algorithm.ReleaseStatusVersioned, algorithm.ReleaseStatusFailed, true},
		{"prepared to published", algorithm.ReleaseStatusPrepared, algorithm.ReleaseStatusPublished, true},
		{"prepared to failed", algorithm.ReleaseStatusPrepared, algorithm.ReleaseStatusFailed, true},
		{"published to failed", algorithm.ReleaseStatusPublished, algorithm.ReleaseStatusFailed, true},

		{"pending to versioned", algorithm.ReleaseStatusPending, algorithm.ReleaseStatusVersioned, false},
		{"pending to published", algorithm.ReleaseStatusPending, algorithm.ReleaseStatusPublished, false},
		{"analyzing to prepared", algorithm.ReleaseStatusAnalyzing, algorithm.ReleaseStatusPrepared, false},
		{"no_release to analyzing", algorithm.ReleaseStatusNoRelease, algorithm.ReleaseStatusAnalyzing, false},
		{"no_release to versioned", algorithm.ReleaseStatusNoRelease, algorithm.ReleaseStatusVersioned, false},
		{"no_release to failed", algorithm.ReleaseStatusNoRelease, algorithm.ReleaseStatusFailed, false},
		{"versioned to published", algorithm.ReleaseStatusVersioned, algorithm.ReleaseStatusPublished, false},
		{"versioned to analyzing", algorithm.ReleaseStatusVersioned, algorithm.ReleaseStatusAnalyzing, false},
		{"prepared to versioned", algorithm.ReleaseStatusPrepared, algorithm.ReleaseStatusVersioned, false},
		{"published to prepared", algorithm.ReleaseStatusPublished, algorithm.ReleaseStatusPrepared, false},
		{"published to analyzing", algorithm.ReleaseStatusPublished, algorithm.ReleaseStatusAnalyzing, false},
		{"failed to analyzing", algorithm.ReleaseStatusFailed, algorithm.ReleaseStatusAnalyzing, false},
		{"failed to pending", algorithm.ReleaseStatusFailed, algorithm.ReleaseStatusPending, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isLegalReleaseTransition(tc.from, tc.to); got != tc.want {
				t.Errorf("isLegalReleaseTransition(%q, %q) = %v, want %v", tc.from, tc.to, got, tc.want)
			}
		})
	}
}

func TestSetStatus_IllegalTransitionLogsAndDoesNotPanic(t *testing.T) {
	core, observed := observedCore()
	r := New(&Context{Logger: zap.New(core)})
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusNoRelease}

	r.setStatus(state, algorithm.ReleaseStatusAnalyzing)

	if state.Status != algorithm.ReleaseStatusNoRelease {
		t.Errorf("illegal transition changed status to %q, want %q",
			state.Status, algorithm.ReleaseStatusNoRelease)
	}
	found := false
	for _, e := range observed.All() {
		if strings.Contains(e.Message, "Illegal release status transition") {
			found = true
		}
	}
	if !found {
		t.Error("expected an illegal-transition warning to be logged")
	}
}

// statusProgressPlugin records the status it sees at each mutating phase.
type statusProgressPlugin struct {
	name             string
	statusAtPrepare  algorithm.ReleaseStatus
	statusAtPublish  algorithm.ReleaseStatus
	statusAtPostPush algorithm.ReleaseStatus
}

func (p *statusProgressPlugin) Name() string { return p.name }

func (p *statusProgressPlugin) AnalyzeCommits(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return algorithm.ReleaseTypeMinor, nil
}

func (p *statusProgressPlugin) Prepare(_ *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	p.statusAtPrepare = state.Status
	return nil
}

func (p *statusProgressPlugin) Publish(_ *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error) {
	p.statusAtPublish = state.Status
	return nil, nil
}

func (p *statusProgressPlugin) PostPublish(_ *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	p.statusAtPostPush = state.Status
	return nil
}

var (
	_ plugins.Plugin         = (*statusProgressPlugin)(nil)
	_ plugins.CommitAnalyzer = (*statusProgressPlugin)(nil)
	_ plugins.Preparer       = (*statusProgressPlugin)(nil)
	_ plugins.Publisher      = (*statusProgressPlugin)(nil)
	_ plugins.PostPublisher  = (*statusProgressPlugin)(nil)
)

func TestRunPluginLifecycle_StatusProgression(t *testing.T) {
	probe := &statusProgressPlugin{name: "status-progress-stub"}
	plugins.Register(probe)

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{probe.name},
	}
	ctx := &Context{Config: cfg, Logger: zap.NewNop(), DryRun: false}
	algoCtx := &algorithm.ReadOnlyContext{
		Config:  cfg,
		Branch:  &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease},
		Commits: []*algorithm.Commit{{Type: "feat", Message: "feat: add thing", Subject: "add thing"}},
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(ctx).runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	if probe.statusAtPrepare != algorithm.ReleaseStatusVersioned {
		t.Errorf("Prepare saw status %q, want %q", probe.statusAtPrepare, algorithm.ReleaseStatusVersioned)
	}
	if probe.statusAtPublish != algorithm.ReleaseStatusPrepared {
		t.Errorf("Publish saw status %q, want %q", probe.statusAtPublish, algorithm.ReleaseStatusPrepared)
	}
	if probe.statusAtPostPush != algorithm.ReleaseStatusPrepared {
		t.Errorf("PostPublish saw status %q, want %q", probe.statusAtPostPush, algorithm.ReleaseStatusPrepared)
	}
	if state.Status != algorithm.ReleaseStatusPublished {
		t.Errorf("final status = %q, want %q", state.Status, algorithm.ReleaseStatusPublished)
	}
}

func TestRunPluginLifecycle_NoRelevantChangesSetsNoRelease(t *testing.T) {
	noop := &priorityTestPlugin{name: "no-release-status-stub", releaseType: ""}
	plugins.Register(noop)

	cfg := &algorithm.Config{
		Branches:  []algorithm.BranchConfig{{Name: "main"}},
		TagFormat: "v${version}",
		Plugins:   []string{noop.name},
	}
	ctx := &Context{Config: cfg, Logger: zap.NewNop(), DryRun: true}
	algoCtx := &algorithm.ReadOnlyContext{
		Config: cfg,
		Branch: &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease},
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(ctx).runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}
	if state.Status != algorithm.ReleaseStatusNoRelease {
		t.Errorf("status = %q, want %q", state.Status, algorithm.ReleaseStatusNoRelease)
	}
}

// rollbackStubPlugin fails verification and records rollback/fail ordering.
type rollbackStubPlugin struct {
	name        string
	events      *[]string
	rollbackErr error
}

func (p *rollbackStubPlugin) Name() string { return p.name }

func (p *rollbackStubPlugin) AnalyzeCommits(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return algorithm.ReleaseTypeMinor, nil
}

func (p *rollbackStubPlugin) VerifyRelease(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return errors.New("verify release exploded")
}

func (p *rollbackStubPlugin) Rollback(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	*p.events = append(*p.events, "rollback")
	return p.rollbackErr
}

func (p *rollbackStubPlugin) Success(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return nil
}

func (p *rollbackStubPlugin) Fail(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState, _ error) error {
	*p.events = append(*p.events, "fail")
	return nil
}

var (
	_ plugins.Plugin          = (*rollbackStubPlugin)(nil)
	_ plugins.CommitAnalyzer  = (*rollbackStubPlugin)(nil)
	_ plugins.ReleaseVerifier = (*rollbackStubPlugin)(nil)
	_ plugins.Rollbacker      = (*rollbackStubPlugin)(nil)
	_ plugins.LifecycleHook   = (*rollbackStubPlugin)(nil)
)

// fakeGit implements only the reads used before the plugin lifecycle.
type fakeGit struct {
	git.GitAPI
}

func (f *fakeGit) GetCurrentBranch() (string, error)                 { return "main", nil }
func (f *fakeGit) GetLastRelease(string) (*algorithm.Release, error) { return nil, nil }
func (f *fakeGit) GetRepositoryURL() (string, error)                 { return "https://example.test/repo.git", nil }
func (f *fakeGit) GetCommits(string, string) ([]*algorithm.Commit, error) {
	return []*algorithm.Commit{{Type: "feat", Message: "feat: add thing", Subject: "add thing"}}, nil
}

func TestRun_OnFailure_CallsRollbackBeforeFailHooks(t *testing.T) {
	t.Setenv("CI", "true")
	t.Setenv("GITHUB_EVENT_NAME", "push")
	t.Setenv("CI_MERGE_REQUEST_ID", "")
	t.Setenv("BUILD_REASON", "")

	events := []string{}
	stub := &rollbackStubPlugin{
		name:        "rollback-on-failure-stub",
		events:      &events,
		rollbackErr: errors.New("rollback boom"),
	}
	plugins.Register(stub)

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{stub.name},
	}
	ctx := &Context{
		Config: cfg,
		Git:    &fakeGit{},
		Logger: zap.NewNop(),
		DryRun: false,
	}

	err := New(ctx).Run()
	if err == nil {
		t.Fatal("expected Run to fail")
	}
	// A rollback error must never replace the original lifecycle error.
	if !strings.Contains(err.Error(), "verify release exploded") {
		t.Errorf("Run error = %q, want the original verify error", err)
	}
	if strings.Contains(err.Error(), "rollback boom") {
		t.Errorf("Run error = %q, rollback error must not mask the original", err)
	}
	if len(events) != 2 || events[0] != "rollback" || events[1] != "fail" {
		t.Errorf("events = %v, want [rollback fail]", events)
	}
}
