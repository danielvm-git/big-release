---
bug_id: BUG-git-push-error-swallowed
status: fixed
severity: high
scope: internal/git
title: git.Client mutating commands discard stderr, collapsing every push failure to "exit status 1"
---

## Problem

CI run [danielvm-git/big-release#30146351433](https://github.com/danielvm-git/big-release/actions/runs/30146351433/job/89649099099)
("release" job, step "Run big-release") failed with:

```
ERROR release/release.go:484 Release failed, running fail hooks
{"error": "plugin \"git\" publish failed: push failed, local tag 2.3.0 removed: failed to push: exit status 1"}
```

Even re-running with `ACTIONS_STEP_DEBUG=true` (`gh run rerun --debug`) produced no
additional detail — GitHub Actions' step debug only instruments commands run through
its own toolkit exec wrapper, not subprocesses spawned from inside the `big-release`
Go binary. The real `git push` stderr was never captured anywhere, making the CI
failure undiagnosable from logs alone.

## Root Cause

`internal/git/client.go`'s `gitCmd()` builds `*exec.Cmd` values with `Stdout`/`Stderr`
left `nil`. Every mutating operation (`Push`, `PushTags`, `PushNotes`, `CreateTag`,
`Commit`, `AddNote`, `VerifyAuth`, `StageChanges`, `StagePaths`, `DeleteTag`) called
`cmd.Run()` directly. Per Go's `os/exec` semantics, a nil `Stderr` on `Run()` is
connected to `/dev/null` — there is no `ExitError.Stderr` fallback the way there is
for `Output()`. The wrapping `fmt.Errorf("failed to push: %w", err)` therefore only
ever had `*exec.ExitError.Error()` to work with, i.e. the bare `"exit status N"`.

This made the actual push-rejection reason (auth, non-fast-forward, missing
upstream, protected ref, etc.) permanently invisible for this whole command
family, in CI and locally.

## Fix

Added `runGit(cmd *exec.Cmd) error` in `internal/git/client.go`: assigns a
`bytes.Buffer` to `cmd.Stderr` before `Run()` and appends its trimmed contents to
the returned error when non-empty. All prior `cmd.Run()` call sites in the mutating
command family now route through it. `IsGitRepo()` (`cmd.Run() == nil`, exit-code-only
check) is intentionally left untouched.

## Verify

- `internal/git/push_error_test.go`: `TestPush_SurfacesGitStderr` and
  `TestPushTags_SurfacesGitStderr` assert the returned error is no longer the bare
  `"... exit status N"` form and contains real git stderr content.
- `go build ./...`, `go vet ./...`, `go test ./... -count=1` (511 tests), and
  `golangci-lint run ./internal/git/...` all pass.

## Files

- `internal/git/client.go`
- `internal/git/push_error_test.go` (new)

## Follow-up

This fix restores visibility only. The next CI push-triggered release will now show
the actual git rejection reason in its logs if the underlying push still fails —
follow up on that specific cause once observed, rather than guessing further here.
Branch protection, repo `default_workflow_permissions` (currently `read`, overridden
by the `release` job's explicit `permissions: contents: write`), and tag protection
were all checked directly against the live repo and ruled out as blockers.
