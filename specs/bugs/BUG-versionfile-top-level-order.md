---
bug_id: BUG-versionfile-top-level-order
status: fixed
severity: high
scope: pkg/release
title: Top-level versionFile auto-adds the plugin after git, so the file is never committed
refs:
  - https://github.com/danielvm-git/big-release/issues/75
  - https://github.com/danielvm-git/big-release/issues/74
---

## Summary

#74 wired the top-level `versionFile` block and auto-added `versionfile`, but
`effectivePlugins()` appends it **after** the configured plugins. Phase 6.5
runs `PostPublish` in declaration order, so `git.PostPublish` stages/commits
`postTagAssets` before `versionfile.PostPublish` writes the file. The tag is
created, but the VERSION commit is empty and the file stays uncommitted.

## Root Cause

`effectivePlugins()` returns `append(names, "versionfile")`. With
`plugins: [changelog, git, github]` and a top-level `versionFile`, the
effective order is `[changelog, git, github, versionfile]` — the opposite of
the documented working order. The docs claim "The top-level `versionFile` form
handles this for you", which is false.

## Fix Approach

Insert `versionfile` immediately **before** the first `git` entry instead of
appending it. If `git` is absent, append it as before. If `versionfile` is
already listed, leave the order untouched.

## Verification

```bash
go test -race ./pkg/release/ -run 'EffectivePlugins|TopLevelVersionFile'
go test -race ./...
go vet ./...
golangci-lint run ./...
```

Expected: the effective list is `[changelog, versionfile, git, github]`, and a
stub git `PostPublish` sees the version file already written.

## Resolution

- `effectivePlugins()` inserts `versionfile` before the first `git` entry.
- Regression tests: `TestEffectivePlugins_InsertsVersionFileBeforeGit`,
  `TestEffectivePlugins_AppendsWhenGitAbsent`,
  `TestEffectivePlugins_ExplicitVersionFileUnchanged`,
  `TestTopLevelVersionFile_WritesBeforeGitPostPublish`.
- Gate: 771 tests (`-race`), vet, lint, build green.
