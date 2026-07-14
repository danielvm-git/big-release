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
	VerifyConditions(ctx *algorithm.Context) error
}

// CommitAnalyzer analyzes commits and returns a release type.
type CommitAnalyzer interface {
	AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error)
}

// ReleaseVerifier verifies the calculated release before proceeding.
type ReleaseVerifier interface {
	VerifyRelease(ctx *algorithm.Context) error
}

// NotesGenerator generates release notes.
type NotesGenerator interface {
	GenerateNotes(ctx *algorithm.Context) (string, error)
}

// Preparer prepares the release (stages changes, runs hooks, etc.).
type Preparer interface {
	Prepare(ctx *algorithm.Context) error
}

// Publisher publishes the release (creates tags, GitHub releases, etc.).
type Publisher interface {
	Publish(ctx *algorithm.Context) (*algorithm.Release, error)
}

// LifecycleHook is called after success or on failure.
type LifecycleHook interface {
	Success(ctx *algorithm.Context) error
	Fail(ctx *algorithm.Context, err error) error
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
