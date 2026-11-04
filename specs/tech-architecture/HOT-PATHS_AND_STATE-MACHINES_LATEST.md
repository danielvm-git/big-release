# Hot Paths & Implicit State Machines — Latest

> **Status 2026-11-04:** NT-1..NT-7, NT-10 and NT-11 are fixed. FSM-1/FSM-3 are
> explicit (`ReleaseStatus`, `PublisherRunState`), `plugins.Rollbacker`
> compensates git tags and forge releases, and each run takes a registry
> snapshot with a `RunResetter` pass so plugin state cannot leak across runs.
> Remaining: NT-8/NT-9 (hook failures are logged, not fatal — intentional),
> NT-12 (PR skip returns nil — intentional).

**Reason for existence:** name every hot path in the release pipeline, record the
sequential I/O that was removed, inventory the implicit finite-state machines
(FSMs) the service already runs, and list the non-terminal states where a run can
stop without a clean terminal outcome. This is the input for making the machines
explicit.

Scope: `pkg/release/` (orchestrator), `internal/plugins/`, `internal/publishers/`,
`internal/git/`, `cmd/big-release/`.

---

## 1. Hot-path sweep

Hot path = code on every `big-release release` invocation.

| # | Hot path | I/O per call | Before | After |
|---|----------|--------------|--------|-------|
| H1 | `Release.buildAlgoContext` → `readGitState` | 3 `git` subprocesses (`rev-parse`, `tag --list`+`rev-list`, `config`) | sequential | **concurrent** (`sync.WaitGroup`) |
| H2 | `Release.buildAlgoContext` → `GetCommits` | 1 `git log` | sequential | sequential (depends on H1's last tag) |
| H3 | `Release.runPluginLifecycle` phases 0–6.5 | plugin calls | sequential | sequential **by design** — each phase consumes the previous phase's output |
| H4 | `publishers.Detect` | `os.Stat` per publisher | sequential | sequential (cheap, filesystem cache) |
| H5 | `Release.runPublishers` Prepare→Publish→Verify | external commands / HTTP | sequential across publishers | sequential **by design** — two-phase barrier + fail-fast ordering |
| H6 | `GitHubPlugin.Success` issue comments | 1 HTTP POST per issue | sequential | **concurrent**, bounded to 4 |
| H7 | `GitHubPlugin.uploadAssets` | 2 HTTP POSTs per asset (file + `.sha256`) | sequential | **concurrent**, bounded to 4 |
| H8 | `GitLabPlugin.uploadAssets` | 2 HTTP POSTs per asset | sequential | sequential — **not fixed** (same pattern as H7, candidate next) |
| H9 | `ExecPlugin.Prepare` commands | user shell commands | sequential | sequential **by design** — the user's list order is the contract |

### Removed sequential I/O (this change)

- **H1** — three independent git reads were awaited one at a time. Each is a
  process spawn; the wall-clock cost was `3 × spawn`. They now run together;
  `GetCommits` still runs after `GetLastRelease` because its range depends on the
  last tag.
- **H6** — N referenced issues meant N serial HTTP round-trips in a success hook
  that runs after the release is already published.
- **H7** — N release assets (binaries + checksums) were uploaded serially.

Bounded concurrency (`maxCommentWorkers = 4`, `maxUploadWorkers = 4`) keeps
memory and connection use predictable and stays below GitHub secondary limits.

### Discovered defect (logged, not fixed)

`internal/publishers/httputil/retry.go` — `RetryClient.Do` reuses the same
`*http.Request` across attempts. `http.Client.Do` closes/consumes the body on the
first attempt, so a retry of a body-bearing request sends an empty body. This is
a non-terminal retry state: the request is neither correctly retried nor cleanly
failed. Fix by rebuilding the request from `req.GetBody` (or a `[]byte` payload)
per attempt.

---

## 2. Implicit state machines

None of these are modelled as an explicit type. Each is a de-facto FSM encoded in
control flow and field mutation.

| ID | Machine | Where | States | Terminal states | Conflict |
|----|---------|-------|--------|-----------------|----------|
| FSM-1 | Release lifecycle | `pkg/release/release.go` | `Idle → PRSkip → ContextBuilt → BranchRejected → Verifying → Analyzing → NoRelease → Versioned → Verified → Channeled → Prepared → Published → PostPublished → Success` | `PRSkip`, `BranchRejected`, `NoRelease`, `Success`, `Failed` | `NoRelease` used to fall through to FSM-3 (fixed) |
| FSM-2 | Plugin capability set | `internal/plugins/registry.go` | each plugin independently implements any subset of 10 capability interfaces | n/a | phase order assumes a plugin that `Publish`es also `Prepare`d |
| FSM-3 | Publisher lifecycle | `internal/publishers/base.go`, `runPublishers` | `Detected → ConfigFiltered → Configured → Prepared → Published → Verified` | `Verified`, `Failed` | no per-publisher state; a failed `Publish` leaves earlier `Prepare` mutations in place |
| FSM-4 | GitHub/GitLab release | `github.go`, `github_assets.go`, `gitlab.go` | `Created → AssetsUploaded → DraftPublished` | `Published` | `Created` and `AssetsUploaded` are non-terminal and unrecoverable |
| FSM-5 | Git ref lifecycle | `internal/plugins/git.go` | `TagCreated → TagsPushed → CommitPushed → NotesPushed` | `NotesPushed` | rollback deletes only the **local** tag |
| FSM-6 | Push outcome | `internal/git/push_outcome.go` | `Unknown / RejectedPolicy / RejectedNonFastForward / NoUpstream / AuthFailed` | classification only | none — this is the model to copy |
| FSM-7 | Branch classification | `pkg/release/release.go` `mapBranchConfig` | `release / maintenance / prerelease` | n/a | empty `type` silently defaults to `release` |
| FSM-8 | Release type | `internal/algorithm/types.go`, `analyzer.go` | `"" / patch / minor / major / prerelease` | `""` = no release | `prerelease` is a type but is also derived from branch type |
| FSM-9 | CI mode | `detectCI`, `detectPR` | `Normal / AutoDryRun / PRSkip` | `AutoDryRun`, `PRSkip` | auto-dry-run silently overrides an explicit `--dry-run=false` |
| FSM-10 | Version file sync | `versionfile.go`, `git.go` `PostPublish` | `Tagged → VersionFileWritten → Committed → Pushed` | `Pushed` | `tagOnly` skips the write entirely |

### FSM-1 detail (the orchestrator)

Phases are comments, not states:

```
0 Configure → 1 VerifyConditions → 2 AnalyzeCommits → 2.5 GenerateNotes
→ 3 Plugin notes → 3.5 CalculateNextVersion → 4 VerifyRelease
→ 4.5 AddChannel → 5 Prepare → 6 Publish → 6.5 PostPublish
→ Success | Fail
```

Transitions are enforced only by `return err` and by `if r.ctx.DryRun`. There is
no transition table, no illegal-transition check, and no state on
`algorithm.MutableState` that says which phase produced the current values.

### FSM-2 detail (capability set)

`Plugin` is one method (`Name`). Nine optional interfaces add capabilities.
A plugin can be a `Publisher` without being a `Preparer`, or a `LifecycleHook`
without ever participating in the pipeline. The orchestrator type-asserts at each
phase, so "not implemented" and "implemented as no-op" are indistinguishable.

### FSM-5 detail (git refs)

```
createTag → PushTags → tagOnly? ──yes──▶ terminal
                │no
                ▼
             Push(commit) ──fail──▶ DeleteTag (local only) → Failed
```

The remote tag from `PushTags` survives the rollback.

---

## 3. Non-terminal states in the release service

A non-terminal state is reachable but has no transition to a clean
`Success`/`Failed` outcome, or leaves external state half-written.

| ID | State | Reached when | Why it is non-terminal | Status |
|----|-------|--------------|------------------------|--------|
| NT-1 | `NoRelease → Publishers` | docs/chore-only commits | publishers ran with `version == ""` | **FIXED** — `hasPublishableRelease` guard in `Run` |
| NT-2 | `NotesWithoutType` | notes generated, no release type | `NextRelease` non-nil, `Version == ""`; nil-check alone missed it | **FIXED** — same guard |
| NT-3 | `PreparedNotPublished` | plugin `Prepare` committed, `Publish` failed | local commit / `CHANGELOG.md` mutated, no tag | open |
| NT-4 | `TagPushedCommitFailed` | `PushTags` ok, `Push` fails | remote tag orphaned; rollback is local-only | open |
| NT-5 | `ReleaseCreatedAssetsFailed` | release POST ok, asset upload fails | release exists without assets/checksums | open |
| NT-6 | `DraftCreatedPublishFailed` | draft release ok, `PATCH draft=false` fails | invisible draft left behind | open |
| NT-7 | `PostPublishFailed` | tag ok, `PostPublish` write/push fails | tag exists, `VERSION` not on branch | open |
| NT-8 | `SuccessHookFailed` | release ok, comment/notification fails | logged and swallowed; release reports success | open (intentional) |
| NT-9 | `FailHookFailed` | failure handler itself fails | logged and swallowed; original error still returned | open (intentional) |
| NT-10 | `RetryBodyConsumed` | 429 retry on a body request | retry sends empty body; not a real retry | open (see §1) |
| NT-11 | `RegistryDrift` | repeated orchestrator invocations | global `plugins`/`publishers` registries keep the first `GitAPI`; mitigated by the `plugins.Get("git")` guard | open |
| NT-12 | `PRSkip` | PR-triggered run | returns `nil`; indistinguishable from a successful release to the caller | open (intentional) |

### Cross-machine conflicts

1. **FSM-1 terminal `NoRelease` vs FSM-3 entry** — the lifecycle declared itself
   done, then the orchestrator entered the publisher machine anyway. Fixed by
   making "has a concrete version" the precondition for FSM-3.
2. **FSM-5 rollback vs remote** — `DeleteTag` only touches the local repo.
3. **FSM-4 draft vs asset failure** — assets are uploaded while the release is
   draft, but a failure leaves the draft and no compensation.
4. **FSM-2 capability subset vs FSM-1 phase assumptions** — no-op and absent are
   the same at the call site.
5. **FSM-9 auto-dry-run vs explicit intent** — `--dry-run=false` is overridden
   when no CI variable is present.
6. **FSM-10 `tagOnly` vs `PostPublish`** — handled by an explicit early return,
   but the skip is silent.

---

## 4. Path to explicit machines

Ordered by leverage. Do not model all ten at once.

1. **`ReleaseStatus` on `MutableState`** — add an enum
   (`pending / analyzing / no_release / versioned / prepared / published / failed`)
   and assert legal transitions in `Run`. This turns FSM-1 from comments into a
   type and makes NT-1/NT-2 impossible by construction.
2. **Per-publisher state** — record `prepared / published / verified` per
   publisher so a retry or failure hook can see half-applied mutations (FSM-3,
   NT-3).
3. **`Rollback` capability** — plugins that create external state (tag, release,
   draft) implement `Rollback(ctx, state)`; the orchestrator calls it on failure
   (FSM-4, FSM-5, NT-4/5/6/7).
4. **Capability manifest** — declare each plugin's capabilities in one table
   instead of scattered type assertions (FSM-2).
5. **Fix `RetryClient.Do`** — rebuild the request body per attempt (NT-10).

`FSM-6` (`PushOutcome`) is the template: a small string enum with a classifier
and a `Hint()` method. Copy that shape.

---

## 5. Verification

```bash
# Full suite with the race detector
go test -race ./...

# The two behaviors added by this change
go test -race ./pkg/release/ -run 'NoRelevantChanges|RunsIndependentGitReadsConcurrently'
go test -race ./internal/plugins/ -run 'PostsCommentsConcurrently|UploadsAssetsConcurrently'

# Lint the touched packages
golangci-lint run ./pkg/release/... ./internal/plugins/...
```

Expected: all pass; the concurrency tests fail against the pre-change code
because the barrier never fills.

`verify: go test -race ./pkg/release/ ./internal/plugins/ && golangci-lint run ./pkg/release/... ./internal/plugins/...`
