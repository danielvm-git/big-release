# QA Audit Report — big-release

**Date:** 2026-07-30
**Auditor:** MiMo-v2.5 (automated QA audit)
**Scope:** Full end-to-end — algorithm, git operations, publishers, plugins, CLI, config, CI/CD

## Run Config

| Parameter | Value | Source |
|-----------|-------|--------|
| `<N>` (ceiling) | 40 | ~16.6k LOC (5-25k range) |
| `<FROZEN>` | Publisher interface, Plugin interface, `git.GitAPI`, `algorithm.Config` schema | AGENTS.md, interfaces |
| Hotspot files | `pkg/release/release.go` (18 churn), `internal/plugins/git.go` (12), `internal/plugins/github.go` (12), `internal/git/client.go` (11) | git log 12mo |
| Preflight | `make preflight` (lint + vet + test) | AGENTS.md |
| Open issues | 0 | gh issue list |
| CI status | All green | gh run list |

## Summary

| Metric | Value |
|--------|-------|
| Total bugs found | 28 |
| Previously fixed | 14 |
| New bugs fixed this audit | 8 |
| New bugs deferred | 12 |
| Critical fixed | 2 (calculator nil-deref, nil-baseVersion) |
| High fixed | 5 (godot/swift Verify, crates/pypi retry body, exec panic) |
| Files audited | 41 source + 36 test |
| Test packages | 18/18 pass |
| CI | Green |

## Bugs Fixed This Audit (New)

### CRITICAL

| ID | Title | Fix |
|----|-------|-----|
| BUG-calculator-nil-branch | `CalculateNextVersion` panics on nil `branch` when `lastRelease` is non-nil | Added nil guard: `if branch == nil { branch = &Branch{Type: BranchTypeRelease} }` |
| BUG-calculator-nil-baseversion | `calculatePrerelease` panics on nil `baseVersion` when `releaseType` unrecognized | Added `default: return nil, fmt.Errorf(...)` |

### HIGH

| ID | Title | Fix |
|----|-------|-----|
| BUG-godot-verify-hardcoded-v | Godot Verify hardcodes `v` prefix, ignores tagFormat | Verify now calls `p.resolveTag(version)` |
| BUG-swift-verify-bare-version | Swift Verify uses bare version, ignores tagFormat | Verify now calls `p.resolveTag(version)` |
| BUG-crates-retry-body-empty | Crates `*bytes.Buffer` body not rewindable on 429 retry | Changed to `bytes.NewReader(body.Bytes())` |
| BUG-pypi-retry-body-empty | PyPI `*bytes.Buffer` body not rewindable on 429 retry | Changed to `bytes.NewReader(buf.Bytes())` |
| BUG-exec-shellsplit-panic | `shellSplit` returns empty slice → panic | Added `len(parts) == 0` guard |

## Previously Fixed (14)

All verified solid: BUG-tag-ignores-tagformat, BUG-cli-hardening, BUG-version-calc, BUG-branch-validation, BUG-nil-panic-analyzer, BUG-ci-binary-missing, BUG-pr-detection, BUG-changelog-format, BUG-changelog-title, BUG-branch-validation-gaps, BUG-push-fails-silently, BUG-branch-config-dead, BUG-publishers-config-ignored, BUG-release-workflow-softprops-and-verbose

## Deferred (12)

| ID | Severity | Title |
|----|----------|-------|
| BUG-git-postpublish-no-tag-check | LOW | PostPublish pushes without tag check |
| BUG-config-explicit-silent-fallback | HIGH | `--config nonexistent.yml` silently uses defaults |
| BUG-ci-cancel-in-progress | HIGH | `cancel-in-progress: true` risks partial releases |
| BUG-ci-missing-arm64 | HIGH | Missing linux/arm64 release binary |
| BUG-getcommits-pipe-corruption | MEDIUM | `GetCommits` pipe-delimiter parsing corrupts on `\|` |
| BUG-git-readpath-stderr | MEDIUM | 9 read-path git functions discard stderr |
| BUG-detectpr-missing-ci | MEDIUM | `detectPR()` missing CircleCI, Bitbucket, Jenkins |
| BUG-health-always-nil | MEDIUM | Health returns nil even when broken |
| BUG-config-js-unparseable | MEDIUM | `big-release.config.js` discovered but YAML-parsed |
| BUG-postpublish-ignores-tagonly | MEDIUM | PostPublish ignores `tagOnly` flag |
| BUG-attribution-push-bypass | MEDIUM | Attribution check skips pushes to main |
| BUG-analyzer-nil-commit | MEDIUM | nil `*Commit` in slice will panic |

## Verification

```
$ go vet ./...         — Clean
$ go test ./...        — 18/18 packages PASS
$ go build ./...       — Compiles clean
$ golangci-lint run    — 2 pre-existing errcheck issues (main_test.go), 0 new
```

## Audit Methodology

6 parallel subagents audited each module boundary:
1. **Algorithm** — calculator nil safety, prerelease edge cases, analyzer coverage
2. **Git operations** — pipe-delimiter parsing, stderr surfacing, dead code
3. **Publishers** — all 9 publishers: tagFormat symmetry, retry body rewinding
4. **Plugins** — shellSplit safety, PostPublish tagOnly, GitHub/GitLab tagFormat
5. **CLI/Config/Release** — config fallback, PR detection, health exit codes
6. **CI/CD** — concurrency groups, attribution checks, release binaries

## Hotspot Analysis (12-month churn)

| File | Churn | Risk |
|------|-------|------|
| `pkg/release/release.go` | 18 | HIGH (orchestrator core) |
| `internal/plugins/github.go` | 12 | HIGH (release API) |
| `internal/plugins/git.go` | 12 | HIGH (tag/push lifecycle) |
| `internal/git/client.go` | 11 | MEDIUM (git operations) |
| `internal/config/config.go` | 9 | MEDIUM (config loading) |
