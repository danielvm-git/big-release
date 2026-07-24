# big-release

> 🚀 Unified, multi-language release automation — one tool, all languages.

[![CI](https://github.com/danielvm-git/big-release/actions/workflows/ci.yml/badge.svg)](https://github.com/danielvm-git/big-release/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielvm-git/big-release)](https://goreportcard.com/report/github.com/danielvm-git/big-release)
[![Release](https://img.shields.io/github/v/release/danielvm-git/big-release)](https://github.com/danielvm-git/big-release/releases)

## What is big-release?

**big-release** is a unified release tool that automatically:

- 📊 **Analyzes commits** using Conventional Commits
- 🔢 **Determines the next version** (patch, minor, major)
- 📝 **Generates changelogs** in [Keep a Changelog](https://keepachangelog.com/) format
- 🏷️ **Creates git tags** with proper formatting
- 📦 **Publishes packages** to any registry (npm, pnpm, PyPI, crates.io, Maven, Go, Swift, Packagist, Godot)
- 🎯 **Creates GitHub releases** with assets, templates, and issue commenting
- 🔀 **Supports GitLab releases** with assets and issue commenting
- 📢 **Multi-channel releases** via git notes and dist-tags

## Why big-release?

| Problem | Solution |
|---------|----------|
| Different tools per language | One tool for all |
| Inconsistent workflows | Unified behavior |
| Complex setup | Single binary, zero config |
| Language-specific learning curve | Same CLI everywhere |
| semantic-release requires Node.js | Go binary, no runtime |

## Supported Languages

| Language | Publisher | Registry |
|----------|-----------|----------|
| JavaScript/TypeScript | `npm` | npmjs.com |
| JavaScript/TypeScript (pnpm) | `pnpm` | npm-compatible registries |
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
# macOS (Apple Silicon)
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-darwin-arm64 -o big-release
chmod +x big-release
sudo mv big-release /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-darwin-amd64 -o big-release
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

# GitHub plugin configuration (optional)
pluginConfigs:
  github:
    assets:
      - path: "dist/*.tar.gz"
      - path: "dist/*.zip"
        label: "Source code"
    draftRelease: false
    releaseName: "v${version}"
    successComment: "🎉 Released in version ${version}"
    releasedLabels:
      - released

# Commit type visibility in changelog (optional)
commitTypes:
  - type: feat
    section: Features
  - type: fix
    section: Bug Fixes
  - type: perf
    section: Performance
  - type: revert
    section: Reverts
  - type: docs
    hidden: true
  - type: chore
    hidden: true
  - type: refactor
    hidden: true
```

## GitHub Action

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    branches: [main, next, 'N.x', beta, alpha]
    tags: ['*']

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v7
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      
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

## GitLab Action

```yaml
# .gitlab-ci.yml
release:
  stage: release
  image: golang:1.26
  rules:
    - if: $CI_COMMIT_TAG
  script:
    - curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-linux-amd64 -o big-release
    - chmod +x big-release
    - ./big-release
  variables:
    GITLAB_TOKEN: $GITLAB_TOKEN
```

## Documentation

- [Getting Started](docs/tutorials/getting-started.md)
- [Configuration Reference](docs/reference/configuration.md)
- [Publishers Reference](docs/reference/publishers.md)
- [Plugin Development](docs/how-to/develop-plugins.md)
- [Publisher Guide](docs/how-to/publishers/README.md)
- [Algorithm Deep Dive](docs/explanation/algorithm.md)

## Architecture

```
big-release/
├── cmd/big-release/          # CLI entry point
├── internal/
│   ├── algorithm/            # Core release algorithm
│   │   ├── analyzer.go       # Commit analysis
│   │   ├── calculator.go     # Version calculation
│   │   ├── generator.go      # Changelog generation
│   │   ├── revert.go         # Revert commit filtering
│   │   └── types.go          # Data types
│   ├── git/                  # Git operations
│   ├── config/               # Configuration loading
│   ├── plugins/              # Plugin system
│   │   ├── changelog.go      # Changelog plugin
│   │   ├── git.go            # Git plugin
│   │   ├── github.go         # GitHub releases
│   │   ├── github_assets.go  # GitHub asset uploads
│   │   ├── github_success.go # GitHub issue commenting
│   │   └── registry.go       # Plugin registry
│   └── publishers/           # Language-specific publishers
│       ├── npm/              # npm publishing
│       ├── pnpm/             # pnpm publishing
│       ├── nodeutil/         # shared package.json helpers
│       ├── pypi/             # PyPI publishing
│       ├── crates/           # crates.io publishing
│       ├── goproxy/          # Go proxy publishing
│       ├── maven/            # Maven publishing
│       ├── swift/            # Swift Package Manager
│       ├── packagist/        # Packagist publishing
│       └── godot/            # Godot Asset Library
├── pkg/release/              # Public API
├── docs/                     # Documentation (Diátaxis)
├── specs/                    # Planning & specs
└── tests/                    # Test suite
```

## Features

### Conventional Commits

big-release analyzes commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` → Minor version bump
- `fix:` → Patch version bump
- `perf:` → Patch version bump
- `BREAKING CHANGE:` → Major version bump
- `revert:` → Handled and filtered from changelog

### Changelog Generation

Generates changelogs in [Keep a Changelog](https://keepachangelog.com/) format:

- `### Added` — New features
- `### Fixed` — Bug fixes
- `### Changed` — Performance improvements
- `### Removed` — Reverted changes

### Multi-Branch Support

Supports multiple release branches with different behaviors:

- **Release branches** (main, next) — Standard releases
- **Maintenance branches** (1.x, 2.x) — Patch releases for older versions
- **Prerelease branches** (beta, alpha) — Pre-release versions

### GitHub Integration

- Creates GitHub releases with customizable templates
- Uploads binary assets to releases
- Comments on issues/PRs referenced in commits
- Adds labels to released issues
- Supports draft releases
- Links to GitHub Discussions

### GitLab Integration

- Creates GitLab releases with assets
- Comments on issues and merge requests
- Supports GitLab CI/CD

### Multi-Channel Releases

Supports different release channels:

- `main` → Stable channel
- `next` → Next channel
- `beta` → Beta channel

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - see [LICENSE](LICENSE)

---

Built with ❤️ by [danielvm-git](https://github.com/danielvm-git)
