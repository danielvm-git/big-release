package release

import (
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

func emptyPublishers() *publishers.Registry { return publishers.NewRegistry() }

func TestNew_IsolatesPluginRegistries(t *testing.T) {
	regA := plugins.NewRegistry()
	regB := plugins.NewRegistry()

	a := New(&Context{Logger: zap.NewNop(), Plugins: regA, Publishers: emptyPublishers()})
	b := New(&Context{Logger: zap.NewNop(), Plugins: regB, Publishers: emptyPublishers()})

	onlyA := &versionTestPlugin{name: "only-in-registry-a"}
	regA.Register(onlyA)

	if _, err := a.plugins.Get("only-in-registry-a"); err != nil {
		t.Fatalf("registry A should contain the plugin: %v", err)
	}
	if _, err := b.plugins.Get("only-in-registry-a"); err == nil {
		t.Fatal("plugin registered in registry A leaked into registry B")
	}
}

func TestRunPluginLifecycle_UsesInjectedPluginRegistry(t *testing.T) {
	reg := plugins.NewRegistry()
	rec := &recorderPlugin{}
	rec.name = "injected-recorder"
	rec.releaseType = algorithm.ReleaseTypePatch
	reg.Register(rec)

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{rec.name},
	}
	ctx := &Context{
		Config:     cfg,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    reg,
		Publishers: emptyPublishers(),
	}
	algoCtx := &algorithm.ReadOnlyContext{
		Config:  cfg,
		Branch:  &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease},
		Commits: []*algorithm.Commit{{Type: "fix", Message: "fix: bug", Subject: "bug"}},
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(ctx).runPluginLifecycle(algoCtx, state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}
	if !rec.published {
		t.Fatal("plugin from the injected registry was not published")
	}
}

func TestRunPublishers_UsesInjectedPublisherRegistry(t *testing.T) {
	reg := publishers.NewRegistry()
	rec := &recordingPublisher{name: "injected-publisher"}
	reg.Register(rec)

	cfg := &algorithm.Config{
		Publishers: map[string]algorithm.PublisherConfig{
			rec.Name(): {Enabled: true},
		},
	}
	ctx := &Context{
		Config:     cfg,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    plugins.NewRegistry(),
		Publishers: reg,
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0"},
	}
	algoCtx := &algorithm.ReadOnlyContext{Config: cfg, Branch: &algorithm.Branch{Name: "main"}}

	if err := New(ctx).runPublishers(algoCtx, state); err != nil {
		t.Fatalf("runPublishers failed: %v", err)
	}
	if !rec.didRun() {
		t.Fatal("publisher from the injected registry was not prepared")
	}
	if len(state.Publishers) != 1 || state.Publishers[0].Status != algorithm.PublisherStatusVerified {
		t.Fatalf("ledger = %+v, want one verified entry", state.Publishers)
	}
}

type resetStubPlugin struct {
	versionTestPlugin
	resets int
}

func (p *resetStubPlugin) ResetRunState() { p.resets++ }

func TestRun_ResetsPluginRunStateBeforeLifecycle(t *testing.T) {
	t.Setenv("CI", "true")
	clearPREnv(t)

	reg := plugins.NewRegistry()
	stub := &resetStubPlugin{}
	stub.name = "reset-stub"
	reg.Register(stub)

	cfg := &algorithm.Config{
		Branches:   []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:  "v${version}",
		Plugins:    []string{},
		Publishers: map[string]algorithm.PublisherConfig{},
	}
	g := &stubGit{
		branch:  "main",
		commits: []*algorithm.Commit{{Hash: "abcdef1234567890", Message: "docs: clarify README"}},
		repoURL: "https://example.com/repo",
	}
	r := New(&Context{
		Config:     cfg,
		Git:        g,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    reg,
		Publishers: emptyPublishers(),
	})

	if err := r.Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if stub.resets != 1 {
		t.Fatalf("ResetRunState called %d times, want 1", stub.resets)
	}
}
