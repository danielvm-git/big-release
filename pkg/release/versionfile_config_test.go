package release

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/plugins"
	"github.com/danielvm-git/big-release/internal/publishers"
)

func versionFileContext(cfg *algorithm.Config) *Context {
	return &Context{
		Config:     cfg,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    plugins.Snapshot(),
		Publishers: publishers.NewRegistry(),
	}
}

func featLifecycleContext(cfg *algorithm.Config) *algorithm.ReadOnlyContext {
	return &algorithm.ReadOnlyContext{
		Config:  cfg,
		Branch:  &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease},
		Commits: []*algorithm.Commit{{Type: "feat", Message: "feat: add thing", Subject: "add thing"}},
	}
}

func TestRunPluginLifecycle_TopLevelVersionFileConfig(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{}, // versionfile is not listed explicitly
		VersionFile:    &algorithm.VersionFileConfig{Path: "BUILD_VERSION", Template: "v{{.Version}}"},
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(versionFileContext(cfg)).runPluginLifecycle(featLifecycleContext(cfg), state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	got, err := os.ReadFile("BUILD_VERSION")
	if err != nil {
		t.Fatalf("top-level versionFile did not write BUILD_VERSION: %v", err)
	}
	if string(got) != "v0.1.0" {
		t.Fatalf("BUILD_VERSION = %q, want %q", got, "v0.1.0")
	}
}

func TestRunPluginLifecycle_PluginConfigOverridesTopLevelVersionFile(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{"versionfile"},
		VersionFile:    &algorithm.VersionFileConfig{Path: "BUILD_VERSION", Template: "v{{.Version}}"},
		PluginConfigs: map[string]map[string]interface{}{
			"versionfile": {"path": "OVERRIDE_VERSION", "template": "{{.Version}}"},
		},
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(versionFileContext(cfg)).runPluginLifecycle(featLifecycleContext(cfg), state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}

	got, err := os.ReadFile("OVERRIDE_VERSION")
	if err != nil {
		t.Fatalf("pluginConfigs override did not win: %v", err)
	}
	if string(got) != "0.1.0" {
		t.Fatalf("OVERRIDE_VERSION = %q, want %q", got, "0.1.0")
	}
	if _, err := os.Stat("BUILD_VERSION"); !os.IsNotExist(err) {
		t.Fatalf("top-level path should not be written when overridden, stat err: %v", err)
	}
}

// postPublishRecorder records the order in which PostPublish hooks run.
type postPublishRecorder struct {
	name   string
	events *[]string
}

func (p *postPublishRecorder) Name() string { return p.name }
func (p *postPublishRecorder) PostPublish(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	*p.events = append(*p.events, p.name)
	return nil
}

func TestPostPublish_RunsInPluginDeclarationOrder(t *testing.T) {
	var events []string

	reg := plugins.NewRegistry()
	reg.Register(&postPublishRecorder{name: "versionfile", events: &events})
	reg.Register(&postPublishRecorder{name: "git", events: &events})

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		// The documented working order: versionfile writes before git commits.
		Plugins: []string{"versionfile", "git"},
	}
	ctx := &Context{
		Config:     cfg,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    reg,
		Publishers: publishers.NewRegistry(),
	}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := New(ctx).runPluginLifecycle(featLifecycleContext(cfg), state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}
	if len(events) != 2 || events[0] != "versionfile" || events[1] != "git" {
		t.Fatalf("PostPublish order = %v, want [versionfile git]", events)
	}
}
