# big-release

> 🚀 Unified, multi-language release automation — one tool, all languages.

[![CI](https://github.com/danielvm-git/big-release/actions/workflows/ci.yml/badge.svg)](https://github.com/danielvm-git/big-release/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielvm-git/big-release)](https://goreportcard.com/report/github.com/danielvm-git/big-release)
[![Release](https://img.shields.io/github/v/release/danielvm-git/big-release)](https://github.com/danielvm-git/big-release/releases)

## What is big-release?

**big-release** is a unified release tool that automatically:

- 📊 **Analyzes commits** using Conventional Commits
- 🔢 **Determines the next version** (patch, minor, major)
- 📝 **Generates changelogs** from commit history
- 🏷️ **Creates git tags** with proper formatting
- 📦 **Publishes packages** to any registry (npm, PyPI, crates.io, etc.)
- 🎯 **Creates GitHub releases** with assets

## Why big-release?

| Problem | Solution |
|---------|----------|
| Different tools per language | One tool for all |
| Inconsistent workflows | Unified behavior |
| Complex setup | Single binary, zero config |
| Language-specific learning curve | Same CLI everywhere |

## Supported Languages

| Language | Publisher | Registry |
|----------|-----------|----------|
| JavaScript/TypeScript | `npm` | npmjs.com |
| Python | `pypi` | pypi.org |
| Rust | `crates` | crates.io |
| Go | `goproxy` | proxy.golang.org |
| PHP | `packagist` | packagist.org |
| Java | `maven` | maven central |
| Swift | `swift` | swiftpackageindex.com |
| Godot/GDScript | `godot` | GitHub Releases |

## Quick Start

### Install

```bash
# macOS
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-darwin-arm64 -o big-release
chmod +x big-release
sudo mv big-release /usr/local/bin/

# Linux
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-linux-amd64 -o big-release
chmod +x big-release
sudo mv big-release /usr/local/bin/
```

### Usage

```bash
# Basic release
big-release

# Dry run (see what would happen)
big-release --dry-run

# Verbose output
big-release --verbose

# Validate configuration
big-release validate

# Show current version
big-release version
```

## Configuration

Create `.big-release.yml` in your project root:

```yaml
# Branch configuration
branches:
  - main                    # release branch
  - next                    # release branch
  - "N.x"                   # maintenance branch
  - name: beta
    prerelease: true        # prerelease branch

# Tag format
tagFormat: "v${version}"

# Publishers (auto-detected, but can be configured)
publishers:
  npm:
    enabled: true
  pypi:
    enabled: true
  crates:
    enabled: true

# Plugins
plugins:
  - changelog
  - git
  - github
```

## GitHub Action

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    branches: [main, next, 'N.x', beta, alpha]

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Install big-release
        run: |
          curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-linux-amd64 -o big-release
          chmod +x big-release
          sudo mv big-release /usr/local/bin/
      
      - name: Run big-release
        run: big-release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Documentation

- [Getting Started](docs/tutorials/getting-started.md)
- [Configuration Reference](docs/reference/configuration.md)
- [Plugin Development](docs/how-to/develop-plugins.md)
- [Publisher Guide](docs/how-to/publishers/README.md)
- [Algorithm Deep Dive](docs/explanation/algorithm.md)

## Architecture

```
big-release/
├── cmd/big-release/          # CLI entry point
├── internal/
│   ├── algorithm/            # Core release algorithm
│   ├── git/                  # Git operations
│   ├── config/               # Configuration loading
│   ├── plugins/              # Plugin system
│   └── publishers/           # Language-specific publishers
├── pkg/release/              # Public API
├── docs/                     # Documentation (Diátaxis)
├── specs/                    # Planning & specs
└── tests/                    # Test suite
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - see [LICENSE](LICENSE)

---

Built with ❤️ by [danielvm-git](https://github.com/danielvm-git)
