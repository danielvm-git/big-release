# Test Plan — Epic e08 Semantic Release Test Parity

> **Epic:** e08  
> **Date:** 2026-07-24  
> **Threat model:** `specs/security/epics/e08/THREAT_MODEL.md`

## Strategy

Risk-scaled parity tests against semantic-release Node.js behaviors. Unit tests for pure helpers; integration via `internal/git/testrepo` for git client; E2E dry-run for pipeline.

## Fixture Plan

| Fixture | Location | Purpose |
|---------|----------|---------|
| `testrepo.Repo` | `internal/git/testrepo/` | Temp git init, commit, tag, chdir, env scrub |
| Stub verify plugins | `pkg/release/aggregate_error_test.go` | Inject ≥2 VerifyConditions failures |
| Config temp files | `internal/config/config_test.go` | YAML/JSON/parent discovery |

---

## e08s01 — Secret Masking (Unit)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s01-P0-01 | P0 | Known env token value never appears in RedactKnownSecrets output | `go test ./internal/secure -run TestMasking_KnownSecrets` |
| SC-e08s01-P0-02 | P0 | Pattern `token=ghp_x` redacted in Redact | `go test ./internal/secure -run TestMasking_Pattern` |
| SC-e08s01-P1-01 | P1 | Generator notes hide token patterns | `go test ./internal/... -run TestMasking_GeneratorNotes` |
| SC-e08s01-P1-02 | P1 | Zap wrapper redacts log message fields | `go test ./internal/secure -run TestMasking_ZapCore` |

---

## e08s02 — Auth URL (Unit)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s02-P0-01 | P0 | HTTPS URL gets token injected as userinfo | `go test ./internal/git -run TestGitAuthURL_HTTPSInject` |
| SC-e08s02-P0-02 | P0 | SSH URL unchanged (passthrough) | `go test ./internal/git -run TestGitAuthURL_SSHPassthrough` |
| SC-e08s02-P0-03 | P0 | URL with existing credentials not double-injected | `go test ./internal/git -run TestGitAuthURL_NoDoubleInject` |
| SC-e08s02-P1-01 | P1 | Error messages redact token values | `go test ./internal/git -run TestGitAuthURL_RedactedErrors` |

---

## e08s03 — E2E Pipeline (Integration)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s03-P0-01 | P0 | Temp repo init→commit→tag→Run dry-run completes | `go test ./pkg/release -run TestE2E_DryRunPipeline` |
| SC-e08s03-P1-01 | P1 | CI env set; no accidental publish side effects | `go test ./pkg/release -run TestE2E_NoPublishSideEffects` |

---

## e08s04 — Commit Traversal (Integration)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s04-P0-01 | P0 | GetLastRelease picks highest semver tag | `go test ./internal/git -run TestCommitTraversal_LastTag` |
| SC-e08s04-P0-02 | P0 | GetCommits from tag..HEAD returns expected range | `go test ./internal/git -run TestCommitTraversal_Range` |
| SC-e08s04-P1-01 | P1 | Empty range when HEAD equals last tag | `go test ./internal/git -run TestCommitTraversal_EmptyRange` |

---

## e08s05 — AggregateError (Unit)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s05-P0-01 | P0 | Two VerifyConditions failures collected in one error | `go test ./pkg/release -run TestAggregateErrors_TwoFailures` |
| SC-e08s05-P1-01 | P1 | AggregateError unwraps all constituent errors | `go test ./pkg/release -run TestAggregateErrors_Unwrap` |

---

## e08s06 — Config Loading (Unit)

| ID | Priority | Scenario | Verify |
|----|----------|----------|--------|
| SC-e08s06-P0-01 | P0 | YAML config file loads and merges defaults | `go test ./internal/config -run TestFileLoading_YAML` |
| SC-e08s06-P0-02 | P0 | JSON config file loads | `go test ./internal/config -run TestFileLoading_JSON` |
| SC-e08s06-P1-01 | P1 | Parent directory discovery finds config | `go test ./internal/config -run TestFileLoading_ParentDiscovery` |
| SC-e08s06-P1-02 | P1 | Explicit CLI config path overrides discovery | `go test ./internal/config -run TestFileLoading_ExplicitPath` |

---

## Coverage Targets

| Package | Target |
|---------|--------|
| internal/secure | 90%+ |
| internal/git (auth + traversal) | 85%+ |
| pkg/release (aggregate + e2e) | 80%+ |
| internal/config | 85%+ |

## Commands

```bash
make test
go test ./internal/secure -run TestMasking
go test ./internal/git -run 'TestGitAuthURL|TestCommitTraversal'
go test ./pkg/release -run 'TestE2E|TestAggregateErrors'
go test ./internal/config -run TestFileLoading
```
