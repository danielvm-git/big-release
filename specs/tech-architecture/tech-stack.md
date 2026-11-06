# Tech Stack & Architecture Decisions

## Language & Runtime
- **Go 1.26+** (compiled static binary, cross-platform distribution without runtime dependencies)
- **YAML / JSON** (declarative configuration parsing via `gopkg.in/yaml.v3` and `encoding/json`)
- **Bash** (operational wrappers, preflight, and local git hooks)

## Core Dependencies
- **spf13/cobra** — CLI routing, flag parsing, and help generation
- **uber-go/zap** — Structured JSON logging and secret redaction
- **Masterminds/semver/v3** — Strict Semantic Versioning 2.0.0 parsing and comparison
- **gopkg.in/yaml.v3** — Configuration loading and serialization

## Architecture & Modules
Four distinct architectural layers:
1. **Algorithm (`internal/algorithm/`)** — Pure business logic; commit message parsing, SemVer calculation, release notes rendering. Zero external I/O or network dependency.
2. **Git Operations (`internal/git/`)** — Repository inspection, commit log retrieval, authenticated remote URLs, tag creation, atomic push operations.
3. **Plugins & Lifecycle (`internal/plugins/`, `pkg/release/`)** — Finite-state release lifecycle engine (`ReleaseStatus`, `PublisherRunState`), rollback compensation, changelog generation, and forge integrations.
4. **Publishers (`internal/publishers/`)** — Modular ecosystem adapters (`npm`, `pnpm`, `pypi`, `cargo`, `go`) implementing detection, version bumping, and publishing interfaces.
5. **CLI (`cmd/big-release/`)** — Root command execution, subcommands (`release`, `validate`, `version`, `health`), flag parsing, and JSON error reporting.

## Domain Language (Ubiquitous Language)
- **Release Run:** A single execution cycle evaluating commit history against the latest tag to determine and execute release actions.
- **Publisher:** A language-specific adapter responsible for manifest inspection, version preparation, and artifact publication.
- **Plugin:** Lifecycle extension point hooked into pre-publish, post-publish, or rollback events.
- **Preflight:** The comprehensive local gate (lint + vet + race-tested test suite) that must pass before any forward progress.
- **Rollback:** Automated compensation routine reversing git tags, remote forge releases, or staging state upon mid-pipeline failure.
- **Dry Run:** Idempotent execution phase calculating bumps and changelog notes without side effects.

## Error-Handling Decisions
- **Aggregate Error Pattern (`pkg/release/errors.go`):** Pipeline steps that encounter multiple non-fatal or composite validation failures accumulate errors into `AggregateError`, reporting all issues together rather than failing on the first.
- **Explicit Error Wrapping:** All internal boundary returns wrap errors with `fmt.Errorf("<context>: %w", err)` to preserve causal chains and unwrapping capability (`errors.Is`, `errors.As`).
- **Graceful Non-Panic Invariants:** State machine illegal transitions log diagnostic context via Zap and fail gracefully rather than triggering panics in production pipelines.
- **Clean CLI Exit Codes:** Controlled exit codes (0 for success/no-release, 1 for fatal failure) paired with clean error messages on `os.Stderr` and structured JSON logs under `--verbose`.

## API-Shape Decisions
- **Publisher Contract:**
  ```go
  type Publisher interface {
      Name() string
      Detect(ctx context.Context, dir string) (bool, error)
      Prepare(ctx context.Context, dir string, nextVersion string) error
      Publish(ctx context.Context, dir string, nextVersion string, opts PublishOptions) error
      Verify(ctx context.Context, dir string, nextVersion string) error
  }
  ```
- **Plugin Lifecycle Hooks:**
  Plugins implement optional lifecycle hooks (`VerifyConditions`, `AnalyzeCommits`, `GenerateNotes`, `Prepare`, `Publish`, `Success`, `Fail`) evaluated in strict sequential order.
- **CLI Subcommand Structure:**
  `big-release [command] [flags]`
  Commands: `release` (execute or dry-run release), `validate` (configuration check), `version` (print tool version), `health` (JSON health status).
