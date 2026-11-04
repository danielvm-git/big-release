# UAT — Execution Plan (gaps + open defects)

Run from the repo root on `main`. Each step must exit 0.

## 1. Spec consistency

```bash
grep -E "status: (complete|done)" specs/release-plan.yaml specs/epics/e05-*/epic.yaml specs/epics/e06-*/epic.yaml specs/epics/e07-*/epic.yaml specs/epics/e09-*/epic.yaml
grep -E "e09s05: done" specs/execution-status.yaml
grep -E "open: 0|deferred: 0" specs/state.yaml
python3 -c "import yaml,glob; [yaml.safe_load(open(f)) for f in glob.glob('specs/**/*.yaml', recursive=True)]"
```

Expected: e05/e06/e07/e09 all complete/done, e09s05 done, bug counts 0 open, YAML parses.

## 2. Hot-path concurrency

```bash
go test -race ./pkg/release/ -run 'NoRelevantChanges|RunsIndependentGitReadsConcurrently' -count=1
go test -race ./internal/plugins/ -run 'PostsCommentsConcurrently|UploadsAssetsConcurrently' -count=1
```

Expected: both pass. Against pre-change code the barrier tests fail because the reads/uploads are serialized.

## 3. NT-4 — remote tag rollback

```bash
go test -race ./internal/git/ -run 'TestDeleteRemoteTag' -count=1
go test -race ./internal/plugins/ -run 'TestGitPluginRollback' -count=1
```

Expected: local + remote tag deleted; no-op when the tag was never pushed.

## 4. NT-10 — retry body rewind

```bash
go test -race ./internal/publishers/httputil/ -run 'TestRetryClient_Body' -count=1
```

Expected: 429→201 resends the full body; body without `GetBody` fails with a descriptive error; bodyless retry still works.

## 5. e09s05 — version file E2E

```bash
go test -race ./pkg/release/ -run 'TestE2E_VersionFileSync' -count=1 -v
```

Expected: `VERSION` contains the computed version; tag `v<version>` exists; the post-tag commit is one commit after the tag.

## 6. FSM-1 — explicit lifecycle

```bash
go test -race ./pkg/release/ -run 'TestIsLegalReleaseTransition|TestSetStatus_IllegalTransitionLogsAndDoesNotPanic|TestRunPluginLifecycle_StatusProgression|TestRunPluginLifecycle_NoRelevantChangesSetsNoRelease' -count=1
```

Expected: legal transitions accepted; illegal ones logged and ignored; status progresses pending→analyzing→versioned→prepared→published; no-release is terminal.

## 7. Rollback compensation

```bash
go test -race ./pkg/release/ -run 'TestRun_OnFailure_CallsRollbackBeforeFailHooks' -count=1
go test -race ./internal/plugins/ -run 'TestGitHubPluginRollback|TestGitLabPluginRollback' -count=1
```

Expected: rollback runs before fail hooks; GitHub/GitLab releases created by a failed run are deleted; rollback is idempotent.

## 8. Full gate

```bash
go test -race ./... && go vet ./... && golangci-lint run ./... && go build ./...
```

Expected: **612 tests pass, 0 lint issues, build OK.**

## Known non-blocking gaps

- NT-8/NT-9: Success/Fail hook errors are logged, not fatal (intentional).
- NT-11: global plugin/publisher registries persist across invocations.
- NT-12: a PR-triggered run returns `nil` (indistinguishable from success to the caller).
- Git rollback gates on local tag presence; a GitHub 201 with an unparseable body leaves `releaseID == 0`.
