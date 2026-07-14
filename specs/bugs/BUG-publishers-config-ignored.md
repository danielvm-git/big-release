---
bug_id: BUG-publishers-config-ignored
status: open
severity: medium
scope: pkg/release
title: Config publishers.enabled flag silently ignored — detection is file-based only
---

## Problem

`runPublishers` calls `publishers.Detect()` which checks file presence (`package.json`, `Cargo.toml`, etc.). The config's `publishers` map with `enabled: true/false` is never consulted. A repo with both `package.json` and `go.mod` always publishes to both, regardless of config.

## Reproduce

1. Create project with `package.json` and `go.mod`
2. Set `publishers: { npm: { enabled: true }, goproxy: { enabled: false } }` in `.big-release.yml`
3. Run `big-release release`
4. Expected: only npm publisher runs
5. Actual: both npm and goproxy publishers run

## Root Cause

`pkg/release/release.go` `runPublishers` method calls `publishers.Detect()` without filtering against `ctx.Config.Publishers`.

## Fix Approach

Filter detected publishers against config map. Skip publishers with `enabled: false`. Backward compatible: detected publishers not in config still run.

## Files

- `pkg/release/release.go`
