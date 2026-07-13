package plugins

import (
	"fmt"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// Plugin defines the interface for plugins
type Plugin interface {
	// Name returns the plugin name
	Name() string

	// VerifyConditions verifies pre-release conditions
	VerifyConditions(ctx *algorithm.Context) error

	// AnalyzeCommits analyzes commits and returns release type
	AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error)

	// VerifyRelease verifies the calculated release before proceeding
	VerifyRelease(ctx *algorithm.Context) error

	// GenerateNotes generates release notes
	GenerateNotes(ctx *algorithm.Context) (string, error)

	// Prepare prepares the release
	Prepare(ctx *algorithm.Context) error

	// Publish publishes the release
	Publish(ctx *algorithm.Context) (*algorithm.Release, error)

	// Success is called after successful release
	Success(ctx *algorithm.Context) error

	// Fail is called on release failure
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

// RunPlugins runs all plugins in order
func (r *Registry) RunPlugins(ctx *algorithm.Context, pluginNames []string, fn func(Plugin) error) error {
	for _, name := range pluginNames {
		plugin, err := r.Get(name)
		if err != nil {
			return fmt.Errorf("failed to get plugin %q: %w", name, err)
		}

		if err := fn(plugin); err != nil {
			return fmt.Errorf("plugin %q failed: %w", name, err)
		}
	}
	return nil
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
