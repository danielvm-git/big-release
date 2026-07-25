# e09s01: PostPublisher Plugin Interface

## Story
Add a new plugin interface `PostPublisher` that runs after the git tag is created and pushed (Phase 6), but before Success hooks. This allows plugins to write version-derived files that reference the finalized version.

## Requirements
- **ADDED**: `PostPublisher` interface with `PostPublish(ctx, state) error` method
- **ADDED**: Phase 6.5 in `runPluginLifecycle` that calls PostPublish on all plugins implementing the interface
- **ADDED**: PostPublish runs after Publish (tag created) but before Success hooks

## Tasks

### Task 1: Add PostPublisher interface to registry
- Add `PostPublisher` interface to `internal/plugins/registry.go`
- Interface: `PostPublish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error`
- Verify: `go vet ./internal/plugins/`

### Task 2: Wire PostPublisher into release lifecycle
- Add Phase 6.5 in `pkg/release/release.go` `runPluginLifecycle`
- Call PostPublish on all plugins implementing PostPublisher after Phase 6
- Verify: `go test ./pkg/release/ -run TestLifecycle`

### Task 3: Add unit test for PostPublisher phase
- Create test that verifies PostPublish is called after Publish
- Use mock plugin implementing PostPublisher
- Verify: `go test ./internal/plugins/ -run TestPostPublisher`

## Verification
- All existing tests pass
- New PostPublisher interface is callable
- Phase ordering is correct (Publish → PostPublish → Success)
