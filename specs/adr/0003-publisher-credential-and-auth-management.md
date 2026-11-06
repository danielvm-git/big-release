# ADR 0003: Publisher Credential and Authentication Management

## Status
Accepted

## Context
Publishing packages across disparate registries (npm, PyPI, crates.io, GitHub Packages, GitLab Package Registry) requires handling sensitive tokens (`NPM_TOKEN`, `PYPI_TOKEN`, `CARGO_TOKEN`, `GITHUB_TOKEN`). Leaking these secrets into CI logs, stdout, or Git commits is a catastrophic security failure.

## Decision
1. **Environment Variable Injection:** Authentication tokens are passed strictly through environment variables or native toolchain credential stores (e.g. system git credential helpers, `~/.netrc`), never through CLI flags or persisted configuration files.
2. **Redaction Layer (`internal/secure`):** A dedicated masking wrapper redacts known environment tokens and high-entropy patterns from all log outputs, stdout, stderr, and error messages before emitting to streams.
3. **Authenticated Git Remote URLs:** Git remote operations requiring authentication dynamically embed tokens in memory for specific commands and scrub them immediately (`git/auth_url.go`), preventing tokens from remaining in local `.git/config`.

## Consequences
- **Positive:**
  - Zero token persistence on disk in target workspaces.
  - Defense-in-depth against accidental secret leakage in build logs.
  - Consistent token naming conventions across all portfolio repositories.
- **Negative:**
  - Requires CI pipelines to explicitly map repository secrets to standard environment variable names.
