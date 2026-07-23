---
bug_id: BUG-nil-panic-analyzer
status: fixed
severity: critical
scope: release-pipeline
title: Nil pointer panic in ChangelogPlugin + unused Analyzer
source: https://github.com/danielvm-git/big-release/issues/6
---

## Summary

Two critical bugs make big-release non-functional:
1. `ChangelogPlugin.resolveNotes()` panics on nil `NextRelease`
2. `algorithm.Analyzer` exists but is never called — no version bump ever happens

## Environment

- OS: macOS arm64
- Version: latest main (commit 2f5dbf3)
- Config: single branch (main), plugins: changelog + git + github

## Bug 1: Nil pointer dereference in ChangelogPlugin.Prepare

### Reproduction
```bash
git tag v0.1.0
git commit --allow-empty -m "feat(ci): test"
CI=1 GITHUB_TOKEN=xxx GITHUB_REPOSITORY=owner/repo big-release release
```

### Root Cause
1. `buildAlgoContext()` initializes `NextRelease: nil`
2. `runPluginLifecycle()` only creates `NextRelease` if `releaseType != ""` or `notes != ""`
3. No plugin implements `AnalyzeCommits()` — all return `("", nil)`
4. So `NextRelease` stays `nil`
5. `ChangelogPlugin.Prepare()` calls `resolveNotes()` which accesses `state.NextRelease.Notes` → panic

### Location
`internal/plugins/changelog.go:99` — `resolveNotes()` method

## Bug 2: Analyzer exists but is never used

### Root Cause
- `internal/algorithm/analyzer.go` has a full `Analyzer` implementation that parses Conventional Commits
- `runPluginLifecycle()` Phase 2 loops over plugins asking for `AnalyzeCommits`
- All 3 built-in plugins (git, github, changelog) return `("", nil)`
- The `Analyzer` is never called as a fallback

### Cascade
1. No plugin returns a release type → `releaseType` stays `""`
2. `NextRelease` is never created
3. Phase 3 `GenerateNotes` runs but `NextRelease` is nil → notes empty
4. Phase 5 `Prepare` calls `resolveNotes` → nil pointer panic

## Verify Steps

- [x] V1: `resolveNotes` handles nil `NextRelease` without panic
- [x] V2: Built-in Analyzer is called as fallback when no plugin returns a type
- [x] V3: `feat:` commit triggers minor release type detection
- [x] V4: `fix:` commit triggers patch release type detection
- [x] V5: Breaking commit triggers major release type detection
- [x] V6: Early exit when no relevant changes (no panic, clean log)
- [x] V7: Full release flow with `feat:` commit completes without panic
- [x] V8: All existing tests still pass

## Fix Approach

1. Add nil guard in `resolveNotes()` (defensive)
2. Wire `algorithm.Analyzer` as fallback in Phase 2 of `runPluginLifecycle()`
3. Add early exit when no release type detected
4. TDD: write failing tests first, then implement
