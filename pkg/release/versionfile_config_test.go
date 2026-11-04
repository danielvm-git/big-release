package release

import (
	"os"
	"slices"
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

func releaseWithPlugins(names []string, vf *algorithm.VersionFileConfig) *Release {
	return New(&Context{
		Config:  &algorithm.Config{Plugins: names, VersionFile: vf},
		Logger:  zap.NewNop(),
		Plugins: plugins.NewRegistry(),
	})
}

func TestEffectivePlugins_InsertsVersionFileBeforeGit(t *testing.T) {
	r := releaseWithPlugins(
		[]string{"changelog", "git", "github"},
		&algorithm.VersionFileConfig{Path: "VERSION"},
	)
	want := []string{"changelog", "versionfile", "git", "github"}
	if got := r.effectivePlugins(); !slices.Equal(got, want) {
		t.Fatalf("effectivePlugins = %v, want %v", got, want)
	}
}

func TestEffectivePlugins_AppendsWhenGitAbsent(t *testing.T) {
	r := releaseWithPlugins(
		[]string{"changelog", "github"},
		&algorithm.VersionFileConfig{Path: "VERSION"},
	)
	want := []string{"changelog", "github", "versionfile"}
	if got := r.effectivePlugins(); !slices.Equal(got, want) {
		t.Fatalf("effectivePlugins = %v, want %v", got, want)
	}
}

func TestEffectivePlugins_ExplicitVersionFileUnchanged(t *testing.T) {
	r := releaseWithPlugins(
		[]string{"versionfile", "git"},
		&algorithm.VersionFileConfig{Path: "VERSION"},
	)
	want := []string{"versionfile", "git"}
	if got := r.effectivePlugins(); !slices.Equal(got, want) {
		t.Fatalf("effectivePlugins = %v, want %v", got, want)
	}
}

// fileCheckGitStub stands in for the git plugin and records whether the
// version file existed when its PostPublish ran.
type fileCheckGitStub struct {
	path    string
	sawFile bool
}

func (p *fileCheckGitStub) Name() string { return "git" }
func (p *fileCheckGitStub) PostPublish(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	_, err := os.Stat(p.path)
	p.sawFile = err == nil
	return nil
}

func TestTopLevelVersionFile_WritesBeforeGitPostPublish(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg := &algorithm.Config{
		Branches:       []algorithm.BranchConfig{{Name: "main"}},
		TagFormat:      "v${version}",
		InitialVersion: "0.1.0",
		Plugins:        []string{"git"}, // versionfile is auto-added, not listed
		VersionFile:    &algorithm.VersionFileConfig{Path: "VERSION", Template: "{{.Version}}"},
	}
	gitStub := &fileCheckGitStub{path: "VERSION"}
	reg := plugins.Snapshot()
	reg.Register(gitStub) // shadow the real git plugin for this run

	r := New(&Context{
		Config:     cfg,
		Logger:     zap.NewNop(),
		DryRun:     false,
		Plugins:    reg,
		Publishers: publishers.NewRegistry(),
	})
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusPending}

	if err := r.runPluginLifecycle(featLifecycleContext(cfg), state); err != nil {
		t.Fatalf("runPluginLifecycle failed: %v", err)
	}
	if !gitStub.sawFile {
		t.Fatal("git PostPublish ran before versionfile wrote VERSION")
	}
}
