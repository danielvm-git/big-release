---
bug_id: BUG-sha256-checksums
status: open
severity: medium
scope: internal/plugins
title: Release binaries published without SHA256 checksums (supply chain risk)
---

# BUG-2026-08-10T131934: Release binaries published without SHA256 checksums

## Problem

- **What happens**: GitHub releases publish only the raw binaries
  (`big-release-linux-amd64`, `big-release-darwin-amd64`, ...). There is no
  `.sha256` checksum asset anywhere on the release. Consumers installing via
  `curl -sL ... -o big-release && chmod +x big-release` have no way to verify
  binary integrity.
- **What should happen**: Each release publishes a `<binary>.sha256` asset
  containing the hex SHA-256 digest so consumers can verify, per issue #67:
  `echo "$(cat big-release.sha256)  big-release" | sha256sum -c -`
  (note: this pattern requires the `.sha256` file to contain ONLY the hex
  digest, not the standard `sha256sum` two-column output).
- **How to reproduce**:
  1. `gh release view v2.4.3 --repo danielvm-git/big-release --json assets -q '.assets[].name'`
  2. Observed: only `big-release-{linux,darwin}-*` and `...-windows-amd64.exe`.
     No `.sha256` files. (Verified 2026-08-10.)
  3. Found during QA audit of `bigbase-canary-go`'s `test-build-release.yml`,
     which installs big-release with no integrity verification.

**Security impact: MEDIUM** — supply-chain hardening gap.
**Exploit path**: A compromised CDN or MITM attacker between the consumer's CI
pipeline and `github.com/.../releases/latest/download/` can substitute a
tampered binary; without a checksum the pipeline cannot detect the
substitution. Adding published SHA-256 digests lets consumers detect any
tampering.

## Root Cause Analysis

**Reproduce** — See above; latest release assets confirmed to lack checksums.

**Isolate** — The release publish path is:
`make release` (cross-compiles 4 binaries into `bin/`) →
`big-release release` → `GitHubPlugin.Publish` →
`uploadAssets` (`internal/plugins/github_assets.go`).
- The uploader iterates configured assets and uploads each file verbatim; it
  computes no digests and creates no derived assets.
- The `Makefile` `release` target only builds binaries — no `sha256sum` step.
- `.big-release.yml` lists only the four binaries as assets.

**Hypothesize** (ranked):
1. H1 — No checksum logic exists anywhere in the publish path.
   *Falsification*: `grep -rn "sha256|checksum" internal/ pkg/` → zero hits. **Confirmed.**
2. H2 — Checksums are expected to be produced by the Makefile.
   *Falsification*: `make release` target contains no `sha256sum`/`shasum` step. **Confirmed.**
3. H3 — The CI workflow generates checksums separately.
   *Falsification*: `.github/workflows/test-build-release.yml` has no checksum step. **Confirmed.**

**Verify** — Root cause confirmed with evidence: **checksum generation is
entirely absent from the release pipeline; the GitHub plugin uploads exactly
the configured assets and computes no digests.** Single root cause, no
contributing defects. Risk level: Low (purely additive feature; no existing
behavior changes).

This is a **missing feature** (not a regression): the fix is to make checksum
publishing a built-in capability of the GitHub plugin so every release —
this repo's and all consumer repos' — gains verifiable integrity. Fixing it in
the plugin (rather than per-repo Makefiles) is the durable fix: 35+ portfolio
repos consume this tool, and the plugin is the single shared publish path.

## TDD Fix Plan

Implemented against `internal/plugins/github_assets.go`. The checksum asset
is `<uploaded-name>.sha256` with hex-only content + trailing newline, matching
the consumer pattern in issue #67.

1. **RED**: `TestGitHubPluginPublish_UploadsChecksums` — configure one binary
   asset; httptest asserts that after the binary upload, a second upload named
   `<name>.sha256` arrives whose body equals the hex SHA-256 of the file
   content (assert via `sha256.Sum256` on the original bytes).
   **GREEN**: In `uploadOneAsset`, after a successful upload, compute the
   SHA-256 digest and upload `<name>.sha256` (hex-only).
   **verify**: `go test ./internal/plugins/ -run UploadsChecksums -count=1`

2. **RED**: `TestGitHubPluginUploadChecksum_UploadFailureFails` — server
   returns 500 for the `.sha256` upload; `Publish` must return an error (a
   release that silently loses checksum integrity must not pass).
   **GREEN**: checksum upload failures return the error from the HTTP status
   check, propagating out of `uploadAssets`.
   **verify**: `go test ./internal/plugins/ -run UploadChecksum -count=1`

3. **GREEN (guard)**: existing `TestGitHubPluginPublish_MissingAssetLogsWarningNotFailure`
   must still pass — a missing asset warns and skips, and must NOT attempt a
   checksum upload or crash.
   **verify**: `go test ./internal/plugins/ -run MissingAsset -count=1`

**REFACTOR**: none anticipated — the change is additive within `uploadOneAsset`;
no public interface changes.

## Acceptance Criteria

- [ ] Each GitHub release includes a `<binary>.sha256` asset per uploaded binary
- [ ] Checksum file contains only the lowercase hex digest + newline (works with `echo "$(cat f.sha256)  f" | sha256sum -c -`)
- [ ] A failed checksum upload fails the release (no silent integrity loss)
- [ ] Missing asset files still warn-and-skip (no checksum attempt, no crash)
- [ ] `make preflight` green; full `go test ./...` green
- [ ] Update `docs/` install instructions with verification snippet (issue #67 pattern)

## Resolution

<!-- filled in by validate-fix -->
