# ADR 0004: Distribution and Self-Release Pipeline

## Status
Accepted

## Context
As a release automation tool, `big-release` must itself be released, versioned, and distributed reliably. Other repositories need a fast, trusted method to fetch and run `big-release` during CI workflows.

## Decision
1. **GitHub Releases as Distribution Host:** Compiled platform binaries are attached as assets to GitHub Releases (`big-release-<os>-<arch>.tar.gz`) alongside SHA-256 checksums.
2. **Self-Releasing Dogfooding:** `big-release` releases itself in GitHub Actions using its own compiled host binary (`make release` followed by `big-release release`).
3. **Download Helper:** A lightweight installer script (`curl -fsSL https://... | sh`) allows consumers to easily pin and install specific releases or `latest` in containerized runners.

## Consequences
- **Positive:**
  - True dogfooding proves tool correctness on every main merge.
  - Zero third-party package registry dependency for distributing the tool itself.
  - Transparent release assets and cryptographic provenance.
- **Negative:**
  - Self-release requires careful staging: host binary must be built and validated before invoking release commands.
