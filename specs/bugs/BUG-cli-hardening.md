---
bug_id: BUG-cli-hardening
status: fixed
severity: high
scope: cmd
title: Three CLI hardening bugs (compile blocker, swallowed errors, dry-run dead)
---

## Bugs Fixed

### Bug 5: pkg/release package missing
- **Root cause:** The import existed but the package directory was never created
- **Fix applied:** Created pkg/release/release.go with orchestrator, pkg/release/release_test.go with 4 tests
- **Hardening:** Tests verify orchestrator instantiation and dry-run behavior

### Bug 1: Error output swallowed
- **Root cause:** os.Exit(1) called without printing the error
- **Fix applied:** Added fmt.Fprintln(os.Stderr, err) before os.Exit(1)
- **Hardening:** Code review standard — all main.go error handlers should print before exit

### Bug 2: --dry-run dead flag
- **Root cause:** Publishers had DryRun fields but nothing set them; GitHub plugin VerifyConditions ran before dry-run gate
- **Fix applied:** Added SetDryRun to Publisher interface, implemented on 8 publishers, wired in orchestrator; GitHub VerifyConditions skips credential check when dry-run
- **Hardening:** Publisher interface contract now includes SetDryRun; new publishers must implement it

## Verification

**Tests:** go test ./... passes (279 passed, 15 packages)
**Build:** go build ./... compiles cleanly
**Lint:** make lint — 0 issues
**Behavioral proof:** `./bin/big-release --dry-run release` exits 0 with dry-run info messages, no side effects
