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

## Follow-up 1: confirmed root cause (resolved)

The fix above landed and the next CI run (30161410418) showed the real error:

```
remote: error: GH006: Protected branch update failed for refs/heads/main.
remote: - 3 of 3 required status checks are expected.
! [remote rejected] main -> main (protected branch hook declined)
```

`main`'s branch protection requires status checks before accepting *any* push to
the branch ref (not just PR merges) — GitHub can't have those checks already
reported for a commit SHA that never existed before, so it rejects direct pushes
outright for any non-admin actor. `github-actions[bot]` (the release job's
`GITHUB_TOKEN` identity) is not an admin, so this is a hard, unconditional block —
no code-level retry or client change can route around a server-side rejection.

Compared against `semantic-release`'s own default behavior (verified via its
source): semantic-release's core release flow never commits anything back to the
branch by default — no plugin implements `prepare` unless the user explicitly adds
`@semantic-release/git`. It only ever creates a tag + pushes that tag, then
creates the GitHub Release via the REST API, which auto-creates the tag server-side
without ever touching `refs/heads/main` — sidestepping branch protection entirely.
Committing version bumps/changelogs to the branch is opt-in, and semantic-release's
own docs warn that opting in requires the user to separately handle branch
protection themselves. There is no way to push new commits to a protected branch
with zero GitHub-side accommodation — that's inherent to git+GitHub, not a
big-release limitation.

### Fix 2

`internal/plugins/git.go`'s `stageChanges()` defaulted to `git add .` (stage
*everything*) whenever no `assets` glob was configured — the opposite of
semantic-release's safe default. Worse, `internal/git/client.go`'s `Commit()`
unconditionally ran `git add -A` internally regardless of what a caller had
already staged, which silently defeated the *existing* explicit `assets:` config
feature too (any file dirtied by another plugin, e.g. `changelog` writing
`CHANGELOG.md` to disk, got swept into the commit regardless of `assets:` scoping).
And `HasChangesToCommit()` checked overall working-tree dirtiness
(`git status --porcelain`) rather than what was staged, so an unrelated dirty file
alone was enough to trigger a commit attempt even when nothing had been staged.

Changed all three to match semantic-release's contract:
- `stageChanges()`: with no `assets` configured, stage nothing (no-op), instead of
  staging everything.
- `Client.Commit()`: commit whatever the caller already staged; stop re-staging
  everything itself.
- `Client.HasChangesToCommit()`: check the index (`git diff --cached --name-only`),
  not overall working-tree state.

Net effect: a fresh `big-release` install with a vanilla config (like this repo's
own `.big-release.yml`, which never set `git.assets`) never attempts to commit or
push anything to the release branch — only the tag, via `git push`, plus the
GitHub Release via the API — matching semantic-release's zero-GitHub-config
default exactly.

### Verify

- `internal/plugins/git_test.go`: new case under `TestGitPluginPrepare` proves
  `Prepare()` stages nothing and commits nothing with no `assets` configured
  (fake git client).
- `internal/plugins/git_prepare_realclient_test.go` (new): end-to-end check
  against a real `git.Client` and real repo — writes an unstaged file (simulating
  `changelog`'s `CHANGELOG.md` write), confirms `Prepare()` leaves `HEAD` and the
  working tree untouched. Confirmed red against the pre-fix code, green after.
- Existing explicit `assets:` config test (`SC-e03s01-P1-...` staged-paths case)
  still passes unchanged.
- `go build`, `go vet`, `go test ./...` (513 tests), all green.

### Files

- `internal/git/client.go` (`Commit`, `HasChangesToCommit`)
- `internal/plugins/git.go` (`stageChanges`)
- `internal/plugins/git_test.go`, `internal/plugins/git_prepare_realclient_test.go` (new)

## Follow-up 2: this repo's own `.big-release.yml`

With the safe default in place, `big-release`'s own repo needs an explicit decision
on whether it still wants `CHANGELOG.md` committed to `main` (which would require
opting back in via `git.assets` *and* separately accommodating branch protection —
unavoidable for that specific feature). Decision: drop the in-repo `CHANGELOG.md`
file; release notes continue to show up on GitHub Releases via the `github`
plugin's API-created release body, matching semantic-release's own default. This
keeps `big-release` installable by any downstream repo with zero GitHub
configuration required.
