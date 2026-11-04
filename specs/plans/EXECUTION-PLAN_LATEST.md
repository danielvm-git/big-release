# Execution Plan — Close Gaps & Open Defects

**Objective:** close every open non-terminal state, fix the latent retry defect,
land the missing e09s05 E2E test, and reconcile spec drift — using parallel
worktrees, one branch per workstream, Orca-managed checkouts, and bigpowers
gates.

**Baseline:** 549 tests pass (`-race`), `go vet` clean, `golangci-lint` 0 issues
on `perf/hot-paths-concurrency`.

---

## 1. Findings that shape the plan

### 1.1 Spec drift (the largest "gap" is bookkeeping)

| Epic | Code reality | `epic.yaml` | `release-plan.yaml` | `execution-status.yaml` |
|------|--------------|-------------|---------------------|-------------------------|
| e05 Publisher/Git Hardening | implemented (RetryClient, GitAPI, unified git path) | `done` | `planning` | `todo` |
| e06 Algorithm/Plugin Unification | implemented (one generator, capability interfaces, priority aggregation) | `planning` | `planning` | `todo` |
| e07 Context Discipline | implemented (`ReadOnlyContext` + `MutableState`) | `planning` | `planning` | `todo` |
| e09 Version File Sync | s01–s04 implemented + docs; **s05 E2E missing** | `planning` | `complete` | `todo` |

Also stale: `specs/state.yaml` (bugs open 1 / deferred 12 — registry says all
fixed), `specs/planning-status.yaml` (discover/not-started), `e09s01-tasks.yaml`
(all tasks `failing` while the code passes).

### 1.2 Open defects

| ID | Defect | Severity | Files |
|----|--------|----------|-------|
| NT-3 | `PreparedNotPublished` — Prepare mutates files, Publish fails | high | release, plugins |
| NT-4 | `TagPushedCommitFailed` — rollback deletes local tag only | high | git plugin, git client |
| NT-5 | `ReleaseCreatedAssetsFailed` — release without assets | high | github/gitlab plugins |
| NT-6 | `DraftCreatedPublishFailed` — orphan draft | medium | github plugin |
| NT-7 | `PostPublishFailed` — tag exists, VERSION not on branch | medium | git plugin |
| NT-8/9 | Success/Fail hook failures swallowed | low | release |
| NT-10 | `RetryClient` cannot retry a body without `GetBody` | medium (latent) | httputil |
| NT-11 | Global registry drift across invocations | low | plugins/publishers |
| NT-12 | PR-skip returns nil (indistinguishable from success) | low | release |
| H8 | GitLab asset uploads serial + no checksum parity | medium | gitlab plugin |
| e09s05 | `TestE2E_VersionFileSync` missing | medium | pkg/release |

Already fixed on `perf/hot-paths-concurrency`: NT-1, NT-2, H1, H6, H7.

---

## 2. Workstreams

Each workstream is one Orca worktree + one branch + one owner. Files are
disjoint inside a wave; that is what makes the wave parallel-safe.

| ID | Branch | Scope | Closes | Files | Depends on |
|----|--------|-------|--------|-------|------------|
| W0 | `perf/hot-paths-concurrency` | land the perf work (squash to main, local) | NT-1/2, H1/6/7 | — | — |
| WA | `chore/spec-reconciliation` | statuses, ledgers, registry dedupe | drift | `specs/**` | W0 |
| WB | `fix/git-ref-atomicity` | remote-tag rollback | NT-4 | `internal/plugins/git.go`, `internal/git/*` | W0 |
| WC | `fix/retry-body-rewind` | rewind body per attempt | NT-10 | `internal/publishers/httputil/retry.go` | W0 |
| WD | `test/e2e-versionfile` | missing E2E | e09s05 | `pkg/release/e2e_test.go` | W0 |
| WE | `fix/release-compensation` | asset failure cleanup + GitLab fan-out/checksums | NT-5, NT-6, H8 | `internal/plugins/github*.go`, `gitlab.go` | W0 |
| WF | `feat/lifecycle-state-machine` | `ReleaseStatus` + `Rollback` capability + per-publisher state | NT-3/7/8/9/11/12, FSM-1/3 | `pkg/release/release.go`, `internal/algorithm/types.go`, `internal/plugins/registry.go` | WB, WE |
| WG | `feat/plugin-rollback` | implement `Rollback` in git/github/gitlab | NT-4/5/6/7 | plugins + git | WF |
| WH | integration | merge, preflight, verify-work, release | all | — | WA–WG |

`WF`/`WG` are split so the interface lands before three plugins implement it;
that avoids a three-way conflict on `registry.go`.

---

## 3. Waves

### Wave 0 — Baseline (serial, blocks everything)

1. Commit `perf/hot-paths-concurrency` (Conventional Commit, no push).
2. Local squash-merge to `main`; rerun `go test -race ./...`, `go vet ./...`,
   `golangci-lint run ./...`.
3. Push only after the user approves (branch protection + attribution rules).

### Wave 1 — Parallel (5 worktrees)

```
WA specs      ─┐
WB git refs   ─┤
WC retry      ─┼─▶ gate: each worktree green (test+vet+lint)
WD e2e        ─┤
WE compensation┘
```

### Wave 2 — Parallel (after WB+WE merge)

```
WF lifecycle state machine ──▶ gate ──▶ WG plugin rollback ──▶ gate
```

### Wave 3 — Integration (serial)

1. Merge order: WA → WC → WD → WB → WE → WF → WG.
2. Full preflight + `verify-work` UAT script.
3. Update `state.yaml`, `execution-status.yaml`, `release-plan.yaml`.
4. Release decision (`release-branch`): PR or `land-branch.sh`.

---

## 4. Orca topology

Repo id: `bc843c02-bae5-4093-a8e7-8051c74edbcb` (big-release).

```text
# one independent checkout per workstream
orca worktree create --name wa-specs --no-parent --base-branch main --json
orca worktree create --name wb-git-refs --no-parent --base-branch main --json
orca worktree create --name wc-retry --no-parent --base-branch main --json
orca worktree create --name wd-e2e --no-parent --base-branch main --json
orca worktree create --name we-compensation --no-parent --base-branch main --json

# one agent per checkout (opencode is the session's agent)
orca agent launch --agent opencode --worktree path:<abs> --prompt "<brief>" --json

# supervision
orca worktree ps --json
orca terminal read --terminal <handle>
orca terminal wait --terminal <handle> --condition exit
```

Worker brief (same for every workstream):

> Read AGENTS.md and CONVENTIONS.md. Use develop-tdd. Touch only the files
> listed for this workstream. Run `go test -race ./...`, `go vet ./...`,
> `golangci-lint run ./...`. Do not commit, push, or edit specs outside your
> scope. Report: files changed, tests added, verify output, open risks.

---

## 5. Gates and merge policy

| Gate | Command | Blocks |
|------|---------|--------|
| G1 per-workstream | `go test -race ./... && go vet ./... && golangci-lint run ./...` | merge |
| G2 no drift | `git diff --check` + YAML parse of touched specs | WA merge |
| G3 integration | `make preflight` | WH |
| G4 attribution | no `Co-authored-by` (hook + CI) | any commit |
| G5 UAT | `verify-work` manual script | release |

Merge policy: squash-merge per workstream, Conventional Commits, no direct
push to main until Wave 3. Rebase each workstream on main immediately before
merge; if a file overlaps outside its declared scope, abort and re-plan.

---

## 6. Rollback

- Per worktree: `git worktree remove --force` + `git branch -D` (work is
  isolated, main untouched).
- Per wave: `git reset --hard <wave-base>` on the integration branch.
- No workstream writes to `main` directly.

---

## 7. Todo ledger

Tracked live with the todo list. Definition of done per item: G1 green +
behavior test + no scope creep. The plan is complete when WH's G3/G5 pass and
the spec drift table shows one consistent status per epic.

`verify: go test -race ./... && go vet ./... && golangci-lint run ./...`
