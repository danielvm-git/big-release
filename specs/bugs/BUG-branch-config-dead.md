---
bug_id: BUG-branch-config-dead
status: fixed
severity: high
scope: pkg/release
title: Branch config Type/Channel/Prerelease never propagated to algorithm layer
---

## Problem

`pkg/release/release.go:92` creates `Branch{Name: branchName}` but never sets `Type`, `Channel`, or `Prerelease` from `BranchConfig`. The Calculator's branch-type dispatch (`calculator.go:48`) always hits the `default` case. Prerelease and maintenance branch support is configured but functionally dead.

## Reproduce

1. Set `branches[].type: prerelease` in `.big-release.yml`
2. Run `big-release release` on that branch
3. Expected: prerelease version (e.g., `1.0.0-alpha.1`)
4. Actual: regular version (e.g., `1.0.0`)

## Root Cause

`release.go:92` only maps `Name`:

```go
Branch: &algorithm.Branch{Name: branchName}
```

`BranchConfig.Type`, `BranchConfig.Channel`, `BranchConfig.Prerelease` are never copied.

## Fix Approach

1. Add `mapBranchConfig()` mapper in `release.go` that copies all fields
2. Add validation in `config.ValidateConfig` that `BranchConfig.Type` is one of `release|maintenance|prerelease` or empty

## Resolution

**Fixed:** 2026-07-25
**Root cause confirmed:** mapBranchConfig() now copies Type, Channel, and Prerelease from BranchConfig
**Fix applied:** Added mapBranchConfig() function in release.go
**Evidence:** `grep -A15 "func mapBranchConfig" pkg/release/release.go`
**Commit:** Integrated into release pipeline

## Files

- `pkg/release/release.go`
- `internal/config/config.go`
