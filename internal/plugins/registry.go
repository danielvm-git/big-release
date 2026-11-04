package plugins

import (
	"fmt"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// Plugin is the minimal interface every plugin must implement.
// Capability checks are done via type assertions to the interfaces below.
type Plugin interface {
	Name() string
}

// ConditionVerifier verifies pre-release conditions.
type ConditionVerifier interface {
	VerifyConditions(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// CommitAnalyzer analyzes commits and returns a release type.
type CommitAnalyzer interface {
	AnalyzeCommits(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (algorithm.ReleaseType, error)
}

// ReleaseVerifier verifies the calculated release before proceeding.
type ReleaseVerifier interface {
	VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// NotesGenerator generates release notes.
type NotesGenerator interface {
	GenerateNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error)
}

// Preparer prepares the release (stages changes, runs hooks, etc.).
type Preparer interface {
	Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// ChannelManager records distribution-channel metadata before publish (e22).
type ChannelManager interface {
	AddChannel(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// Publisher publishes the release (creates tags, GitHub releases, etc.).
type Publisher interface {
	Publish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error)
}

// PostPublisher runs after the tag is created and pushed (Phase 6.5).
// Use this to write version-derived files that reference the finalized version.
type PostPublisher interface {
	PostPublish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// LifecycleHook is called after success or on failure.
type LifecycleHook interface {
	Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
	Fail(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, err error) error
}

// Rollbacker compensates external state created by a plugin when the release
// lifecycle fails after it has begun (FSM-1 failure edge). Rollback errors are
// logged by the orchestrator and never mask the original lifecycle error.
type Rollbacker interface {
	Rollback(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error
}

// ConfigurablePlugin receives its typed configuration from
// Config.PluginConfigs before the pipeline runs. Plugins that need config
// implement this; the orchestrator type-asserts and calls Configure.
type ConfigurablePlugin interface {
	Configure(raw map[string]interface{}) error
}

// Registry manages plugins
type Registry struct {
	plugins map[string]Plugin
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register registers a plugin
func (r *Registry) Register(plugin Plugin) {
	r.plugins[plugin.Name()] = plugin
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, error) {
	plugin, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %q not found", name)
	}
	return plugin, nil
}

// List returns all registered plugin names
func (r *Registry) List() []string {
	var names []string
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

// Snapshot returns a copy of the registry. A release run takes a snapshot so
// it can register per-run plugins (e.g. a GitPlugin bound to that run's
// GitAPI) without mutating the process-wide registry (NT-11).
func (r *Registry) Snapshot() *Registry {
	out := NewRegistry()
	for name, p := range r.plugins {
		out.plugins[name] = p
	}
	return out
}

// Snapshot returns a copy of the global registry.
func Snapshot() *Registry { return globalRegistry.Snapshot() }

// RunResetter resets state a plugin instance remembers between runs. The
// registry reuses instances, so a run must clear stale state (e.g. a release
// ID from a previous run) before it can roll anything back.
type RunResetter interface {
	ResetRunState()
}

// Global registry instance
var globalRegistry = NewRegistry()

// Register registers a plugin in the global registry
func Register(plugin Plugin) {
	globalRegistry.Register(plugin)
}

// Get retrieves a plugin from the global registry
func Get(name string) (Plugin, error) {
	return globalRegistry.Get(name)
}

// List returns all registered plugins from the global registry
func List() []string {
	return globalRegistry.List()
}
