# ADR 0001 — Explicit release lifecycle state machines

- **Status:** Accepted
- **Date:** 2026-11-04
- **Deciders:** maintainer
- **Supersedes:** none

## Context

The release pipeline encoded its lifecycle as phase comments and mutable
values on `algorithm.MutableState`. Three concrete gaps followed:

- **Non-terminal "no release".** The lifecycle could finish without a
  version, yet `Run` still entered the publisher phase with an empty version.
- **Process-global registries.** A `Release` reused the global plugin
  registry, so a `GitPlugin` bound to one run's `GitAPI` could be reused by
  another. Forge plugins also remembered a release ID/tag across runs, so a
  rollback in a later failing run could delete a release it did not create.
- **Invisible skips.** A pull-request-triggered skip returned `nil` with no
  `$GITHUB_OUTPUT` signal, so CI could not distinguish a skip from a publish.

## Decision

1. **FSM-1 — release lifecycle.** `algorithm.ReleaseStatus` plus
   `legalReleaseTransitions` in `pkg/release/release.go` model
   `pending → analyzing → {no_release | versioned → prepared → published |
   failed}`. Illegal transitions are logged and ignored; they never panic.
2. **FSM-3 — publisher lifecycle.** `algorithm.PublisherRunState` records
   `pending → prepared → published → verified | failed` per publisher on
   `MutableState.Publishers`; failures carry the error text.
3. **Compensation.** `plugins.Rollbacker` is the contract. `Run` calls it on
   failure before `Fail` hooks. Rollback errors are logged and never mask the
   original lifecycle error. `GitPlugin`, `GitHubPlugin` and `GitLabPlugin`
   implement it.
4. **Per-run registries.** Each `Release` takes `Registry.Snapshot()` of
   plugins and publishers (overridable via `Context.Plugins` /
   `Context.Publishers`). `plugins.RunResetter` clears per-run state (forge
   release IDs) before the lifecycle starts.
5. **CI signal.** `Run` emits `published`, `released` and `reason`
   (`published | pull_request | no_release | dry_run`) to `$GITHUB_OUTPUT`.

## Consequences

**Positive**

- Illegal lifecycle states are detectable instead of silently overwriting.
- Skipped runs are distinguishable from publishes in CI.
- A rollback cannot touch a previous run's artifacts.
- Tests inject registries without polluting process-global state.

**Negative / accepted**

- Rollback is best-effort: a failed DELETE leaves an orphan and is logged.
- The git plugin cannot undo a pushed release commit (documented, deliberate).
- `Success`/`Fail` hook errors remain non-fatal by design (NT-8/NT-9).

## Alternatives considered

- **Panic on illegal transitions** — rejected; a release tool must not crash
  on a bookkeeping bug.
- **Factory-per-plugin registries** — rejected for now; it needs constructors
  for every plugin and a larger refactor. Snapshot + `RunResetter` fixes the
  observed leaks with a smaller surface.
- **Dropping `$GITHUB_OUTPUT` writes** — rejected; downstream CI needs the
  skip/publish signal.

## References

- `specs/tech-architecture/HOT-PATHS_AND_STATE-MACHINES_LATEST.md`
- `specs/plans/EXECUTION-PLAN_LATEST.md`
- `specs/bugs/BUG-versionfile-top-level-config.md`
- `specs/bugs/BUG-versionfile-top-level-order.md`
