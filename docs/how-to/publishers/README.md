# Publisher Guide

This guide explains how to configure and use language-specific publishers.

## Overview

Publishers handle the actual package publishing to registries. big-release auto-detects the package type based on manifest files.

## Supported Publishers

| Publisher | Manifest File | Registry | Environment Variable |
|-----------|---------------|----------|---------------------|
| npm | `package.json` | npmjs.com | `NPM_TOKEN` |
| PyPI | `pyproject.toml`, `setup.py` | pypi.org | `PYPI_TOKEN` |
| crates.io | `Cargo.toml` | crates.io | `CARGO_TOKEN` |
| Go Proxy | `go.mod` | proxy.golang.org | - |
| Packagist | `composer.json` | packagist.org | `PACKAGIST_TOKEN` |
| Maven | `pom.xml` | maven central | `MAVEN_TOKEN` |
| Swift | `Package.swift` | swiftpackageindex.com | - |
| Godot | `export_presets.cfg` | GitHub Releases | - |

## Configuration

### npm

```yaml
publishers:
  npm:
    enabled: true
    options:
      registry: "https://registry.npmjs.org"
      access: "public"  # or "private"
```

**Environment Variables:**
- `NPM_TOKEN`: npm authentication token

**Steps:**
1. Update `package.json` version
2. Run `npm publish`
3. Verify publication

### PyPI

```yaml
publishers:
  pypi:
    enabled: true
    options:
      repository: "https://upload.pypi.org/legacy/"
```

**Environment Variables:**
- `PYPI_TOKEN`: PyPI authentication token

**Steps:**
1. Update `pyproject.toml` or `setup.py` version
2. Build package: `python -m build`
3. Upload: `twine upload dist/*`
4. Verify publication

### crates.io

```yaml
publishers:
  crates:
    enabled: true
```

**Environment Variables:**
- `CARGO_TOKEN`: crates.io authentication token

**Steps:**
1. Update `Cargo.toml` version
2. Run `cargo publish`
3. Verify publication

### Go Proxy

```yaml
publishers:
  goproxy:
    enabled: true
```

**Steps:**
1. Create git tag with module path
2. Run `go mod download`
3. Verify publication

### Packagist

```yaml
publishers:
  packagist:
    enabled: true
    options:
      webhook_url: "https://packagist.org/webhook"
```

**Environment Variables:**
- `PACKAGIST_TOKEN`: Packagist authentication token

**Steps:**
1. Update `composer.json` version
2. Trigger Packagist webhook
3. Verify publication

### Maven

```yaml
publishers:
  maven:
    enabled: true
    options:
      repository: "https://oss.sonatype.org/staging/"
```

**Environment Variables:**
- `MAVEN_TOKEN`: Maven Central authentication token

**Steps:**
1. Update `pom.xml` version
2. Run `mvn release:perform`
3. Verify publication

### Swift

```yaml
publishers:
  swift:
    enabled: true
    options:
      registry: "https://swiftpackageindex.com"
```

**Steps:**
1. Update `Package.swift` version
2. Create git tag
3. Verify publication

### Godot

```yaml
publishers:
  godot:
    enabled: true
    options:
      export_presets: "export_presets.cfg"
```

**Steps:**
1. Update `export_presets.cfg` version
2. Update version code (timestamp)
3. Commit changes
4. Create GitHub release

## Custom Publishers

### Creating a Custom Publisher

1. Create a new directory in `internal/publishers/`
2. Implement the `Publisher` interface:

```go
type Publisher interface {
    Name() string
    Detect() bool
    Prepare(version string) error
    Publish(version string) error
    Verify(version string) error
}
```

3. Register the publisher in `init()`:

```go
func init() {
    publishers.Register(NewPublisher())
}
```

### Example Custom Publisher

```go
package custom

import (
    "github.com/danielvm-git/big-release/internal/publishers"
)

type Publisher struct{}

func NewPublisher() *Publisher {
    return &Publisher{}
}

func (p *Publisher) Name() string {
    return "custom"
}

func (p *Publisher) Detect() bool {
    // Detect if this publisher should be used
    _, err := os.Stat("custom-manifest.json")
    return err == nil
}

func (p *Publisher) Prepare(version string) error {
    // Update manifest with new version
    return nil
}

func (p *Publisher) Publish(version string) error {
    // Publish to custom registry
    return nil
}

func (p *Publisher) Verify(version string) error {
    // Verify publication
    return nil
}

func init() {
    publishers.Register(NewPublisher())
}
```

## Troubleshooting

### "Publisher not found"

Ensure the publisher is registered in `init()`:

```go
func init() {
    publishers.Register(NewPublisher())
}
```

### "Detection failed"

Check if the manifest file exists:

```go
func (p *Publisher) Detect() bool {
    _, err := os.Stat("package.json")
    return err == nil
}
```

### "Publish failed"

1. Check environment variables are set
2. Verify credentials are valid
3. Check network connectivity
4. Review registry-specific error messages

### "Verification failed"

1. Check if package exists on registry
2. Verify version number matches
3. Wait for registry propagation (may take a few minutes)
