# Contributing to big-release

Thank you for your interest in contributing to big-release! This document provides guidelines and information for contributors.

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- Go 1.26 or later
- Git
- Make
- golangci-lint (for linting)

### Development Setup

1. Clone the repository:
   ```bash
   gh repo clone danielvm-git/big-release
   cd big-release
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feat/my-feature
```

### 2. Make Changes

- Follow Go coding standards
- Write tests for new functionality
- Update documentation if needed

### 3. Run Preflight

```bash
make preflight
```

This runs:
- Linting (`golangci-lint`)
- Vet (`go vet`)
- Tests (`go test`)

### 4. Commit Changes

Use Conventional Commits:

```bash
git commit -m "feat: add new publisher for Dart"
```

### 5. Create a Pull Request

```bash
gh pr create
```

## Commit Message Format

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `perf` | Performance improvement | Patch |
| `docs` | Documentation | None |
| `chore` | Maintenance | None |
| `style` | Code style | None |
| `refactor` | Code refactoring | None |
| `test` | Tests | None |
| `build` | Build system | None |
| `ci` | CI/CD | None |
| `revert` | Revert commit | None |

### Breaking Changes

Use `BREAKING CHANGE:` in the footer or `!` after the type:

```
feat!: remove deprecated API

BREAKING CHANGE: The old API has been removed
```

## Code Style

### Go

- Follow [Effective Go](https://go.dev/doc/effective-go) guidelines
- Use `gofmt` and `goimports` for formatting
- Run `golangci-lint` before committing

### Documentation

- Use clear, concise language
- Include code examples
- Follow Diátaxis framework (tutorials, how-to, reference, explanation)

## Testing

### Unit Tests

```bash
make test-unit
```

### Integration Tests

```bash
make test-integration
```

### All Tests

```bash
make test
```

## Pull Request Guidelines

1. **One feature per PR** - Keep PRs focused
2. **Write tests** - All new features must have tests
3. **Update documentation** - If changing public API
4. **Pass preflight** - All checks must pass
5. **Conventional Commits** - Follow commit message format

## Git Attribution

**Important:** All commits must appear as if authored solely by the human user. Never include:

- `Co-authored-by:` footer
- `Co-Authored-By:` footer
- Any spelling or casing variant

This is enforced by pre-tool-use hooks and CI.

## Release Process

Releases are automated by `.github/workflows/test-build-release.yml` using big-release itself:

1. Push to `main` (or merge a PR) — trunk-based; no separate `deploy.yml`
2. CI runs **lint → test → build**; build uploads artifact `big-release-<sha>`
3. `release` job downloads that artifact (no rebuild) and runs `big-release release`
4. big-release analyzes commits, determines the next version, tags, and creates the GitHub release

### Branch protection

`main` requires status checks **lint**, **test**, and **build**. Solo owner may push directly to `main` (`enforce_admins: false`); prefer PRs so the checks run before merge.

## Getting Help

- Open an issue for bugs or feature requests
- Start a discussion for questions
- Join our community chat (if available)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
