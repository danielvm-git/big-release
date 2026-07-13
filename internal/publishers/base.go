package publishers

import (
	"fmt"
)

// Publisher defines the interface for package publishers
type Publisher interface {
	// Name returns the publisher name
	Name() string

	// Detect detects if this publisher should be used
	Detect() bool

	// Prepare prepares the package for publishing
	Prepare(version string) error

	// Publish publishes the package
	Publish(version string) error

	// Verify verifies the publication
	Verify(version string) error

	// SetDryRun sets the dry-run mode on the publisher
	SetDryRun(dryRun bool)
}

// Registry manages publishers
type Registry struct {
	publishers map[string]Publisher
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		publishers: make(map[string]Publisher),
	}
}

// Register registers a publisher
func (r *Registry) Register(publisher Publisher) {
	r.publishers[publisher.Name()] = publisher
}

// Get retrieves a publisher by name
func (r *Registry) Get(name string) (Publisher, error) {
	publisher, ok := r.publishers[name]
	if !ok {
		return nil, fmt.Errorf("publisher %q not found", name)
	}
	return publisher, nil
}

// Detect returns all publishers that should be used
func (r *Registry) Detect() []Publisher {
	var detected []Publisher
	for _, publisher := range r.publishers {
		if publisher.Detect() {
			detected = append(detected, publisher)
		}
	}
	return detected
}

// List returns all registered publisher names
func (r *Registry) List() []string {
	var names []string
	for name := range r.publishers {
		names = append(names, name)
	}
	return names
}

// Global registry instance
var globalRegistry = NewRegistry()

// Register registers a publisher in the global registry
func Register(publisher Publisher) {
	globalRegistry.Register(publisher)
}

// Get retrieves a publisher from the global registry
func Get(name string) (Publisher, error) {
	return globalRegistry.Get(name)
}

// Detect returns all detected publishers from the global registry
func Detect() []Publisher {
	return globalRegistry.Detect()
}

// List returns all registered publishers from the global registry
func List() []string {
	return globalRegistry.List()
}
