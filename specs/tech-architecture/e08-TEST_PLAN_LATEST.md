# Test Design: e08-semantic-release-test-parity

> **Epic:** e08 — Semantic Release Test Parity  
> **Date:** 2026-07-24  
> **Threat model:** [`specs/security/epics/e08/THREAT_MODEL.md`](../security/epics/e08/THREAT_MODEL.md)  
> **Skill:** plan-tests (epic-scoped, once)

Risk-scaled parity tests against semantic-release Node.js behaviors. Push tests to the lowest level: unit for pure helpers; integration via `internal/git/testrepo` for git client; E2E dry-run for pipeline.

Scenario ID format: `SC-e08s0Y-P{0|1|2|3}-NN`.

---

## 1. Risk Matrix & Scenarios

### e08s01 — Secret Masking & Redaction (Unit)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s01-P0-01 | Known env token (`GH_TOKEN`/`NPM_TOKEN`) never appears in `RedactKnownSecrets` output | P0 | Unit | `internal/secure` |
| SC-e08s01-P0-02 | Pattern forms (`token=ghp_…`, bearer headers) redacted by `Redact` | P0 | Unit | `internal/secure` |
| SC-e08s01-P1-01 | Generator notes hide token patterns (defense in depth) | P1 | Unit | `internal/algorithm/generator` |
| SC-e08s01-P1-02 | Zap core wrapper redacts message/field values | P1 | Unit | `internal/secure` + `cmd/big-release` |
| SC-e08s01-P2-01 | Empty / no-secret input passes through unchanged | P2 | Unit | `internal/secure` |

**Package verify:** `go test ./internal/... -run TestMasking`

### e08s02 — Git Authentication & URL Mutation (Unit)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s02-P0-01 | HTTPS remote gets token injected as userinfo | P0 | Unit | `internal/git.AuthURL` |
| SC-e08s02-P0-02 | SSH remote unchanged (passthrough) | P0 | Unit | `internal/git.AuthURL` |
| SC-e08s02-P0-03 | URL with existing credentials not double-injected | P0 | Unit | `internal/git.AuthURL` |
| SC-e08s02-P1-01 | Error paths that may include URL redact token values | P1 | Unit | `internal/git` + `secure` |
| SC-e08s02-P2-01 | Non-HTTP schemes (git://, file://) passthrough | P2 | Unit | `internal/git.AuthURL` |

**Package verify:** `go test ./internal/git -run TestGitAuthURL`

### e08s03 — E2E Git Repository Simulation (E2E / Integration)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s03-P0-01 | Temp repo `Init`→`Commit`→`Tag`→`Run` dry-run completes without error | P0 | E2E | `pkg/release` + `testrepo` |
| SC-e08s03-P1-01 | With `CI=true` / dry-run, no publish or remote side effects | P1 | E2E | `pkg/release` |
| SC-e08s03-P2-01 | ScrubEnv clears host `GH_TOKEN` from child git env | P2 | Integration | `internal/git/testrepo` |

**Package verify:** `go test ./pkg/release -run TestE2E`

### e08s04 — Git Commit Traversal & Filtering (Integration)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s04-P0-01 | `GetLastRelease` selects highest semver tag (sort correct) | P0 | Integration | `internal/git` + `testrepo` |
| SC-e08s04-P0-02 | Commits for `from..HEAD` return expected range | P0 | Integration | `internal/git` + `testrepo` |
| SC-e08s04-P1-01 | Empty range when HEAD equals last release tag | P1 | Integration | `internal/git` + `testrepo` |
| SC-e08s04-P2-01 | Non-semver / unmatched tags ignored by last-release selection | P2 | Integration | `internal/git` |

**Package verify:** `go test ./internal/git -run TestCommitTraversal`

### e08s05 — Aggregate Pipeline Error Handling (Unit)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s05-P0-01 | ≥2 VerifyConditions failures → single AggregateError containing all | P0 | Unit | `pkg/release` |
| SC-e08s05-P1-01 | AggregateError unwraps / joins all constituent errors (`errors.Is`/`Unwrap`) | P1 | Unit | `pkg/release` |
| SC-e08s05-P2-01 | Single VerifyConditions failure still returns that error (no empty aggregate) | P2 | Unit | `pkg/release` |

**Package verify:** `go test ./pkg/release -run TestAggregateErrors`

### e08s06 — Configuration Loading File + CLI (Unit)

| Scenario ID | Behavior Description | Risk | Test Level | Target |
|-------------|----------------------|------|------------|--------|
| SC-e08s06-P0-01 | YAML config loads and merges with defaults | P0 | Unit | `internal/config` |
| SC-e08s06-P0-02 | JSON config loads | P0 | Unit | `internal/config` |
| SC-e08s06-P1-01 | Parent directory discovery finds nearest config | P1 | Unit | `internal/config` |
| SC-e08s06-P1-02 | Explicit CLI config path overrides discovery | P1 | Unit | `internal/config` |
| SC-e08s06-P1-03 | CLI dry-run flag merges onto loaded config | P1 | Unit | `internal/config` |
| SC-e08s06-P3-01 | Missing / invalid config returns actionable error (no panic) | P3 | Unit | `internal/config` |

**Package verify:** `go test ./internal/config -run TestFileLoading`

---

## 2. Fixture Architecture & Isolation

### `internal/git/testrepo` (shared — land before Wave 2/3)

| Helper | Contract |
|--------|----------|
| `Init(t)` | `t.TempDir()` + `git init` + initial branch; returns `*Repo` rooted at temp path |
| `Commit(msg)` | Creates a conventional-commit style commit (configurable subject) |
| `Tag(name)` | Annotated or lightweight tag at HEAD |
| `ScrubEnv(t)` | Clears/overrides `GH_TOKEN`, `GITHUB_TOKEN`, `NPM_TOKEN`, and related credential env for subprocesses |
| (optional) `Chdir` | `t.Chdir` into repo for CWD-sensitive callers |

**Isolation rules:**
- Always `t.TempDir()` — no leftover fixtures under `/tmp` ownership
- Never point remotes at production hosts
- Token fixtures are ephemeral in-test strings, never committed
- E2E runs with dry-run / `CI=true`

### Other fixtures

| Fixture | Location | Purpose |
|---------|----------|---------|
| Stub VerifyConditions plugins | `pkg/release/aggregate_error_test.go` | Inject ≥2 failures without network |
| Config temp trees | `internal/config/config_test.go` | YAML/JSON files + nested parent dirs |
| Log capture buffer | `internal/secure` / zap tests | Assert redaction on written bytes |

No network intercepts required for e08 (AuthURL is pure; E2E is local git only).

---

## 3. NFR Verification

| NFR Type | Requirement | Verification Command |
|----------|-------------|----------------------|
| Security | Token values absent from logs/notes/errors (CWE-532 / AuthURL) | `go test ./internal/... -run 'TestMasking\|TestGitAuthURL_Redacted'` |
| Isolation | Temp repos cleaned; env scrubbed | `go test ./internal/git/testrepo ./pkg/release -run 'TestE2E\|TestCommitTraversal'` |
| Perf (lite) | Unit suites for secure/auth/config complete in < 5s wall | `go test ./internal/secure ./internal/git ./internal/config -count=1` |
| Reliability | Full suite green after each wave | `make test` / `make preflight` |

---

## 4. Out of Scope

- Executing `big-release.config.js` (JS runtime)
- Aggregating errors outside VerifyConditions
- Mutating `remote.origin.url` on disk
- Full Node semantic-release suite port
- Live registry / GitHub API publish in E2E

---

## 5. Coverage Targets

| Package | Target |
|---------|--------|
| `internal/secure` | 90%+ |
| `internal/git` (auth + traversal + testrepo) | 85%+ |
| `pkg/release` (aggregate + e2e) | 80%+ |
| `internal/config` | 85%+ |

## 6. Commands

```bash
make test
go test ./internal/secure -run TestMasking
go test ./internal/git -run 'TestGitAuthURL|TestCommitTraversal'
go test ./pkg/release -run 'TestE2E|TestAggregateErrors'
go test ./internal/config -run TestFileLoading
```

## Story coverage checklist

| Story | Scenario IDs present |
|-------|----------------------|
| e08s01 | SC-e08s01-P0-01, P0-02, P1-01, P1-02, P2-01 |
| e08s02 | SC-e08s02-P0-01, P0-02, P0-03, P1-01, P2-01 |
| e08s03 | SC-e08s03-P0-01, P1-01, P2-01 |
| e08s04 | SC-e08s04-P0-01, P0-02, P1-01, P2-01 |
| e08s05 | SC-e08s05-P0-01, P1-01, P2-01 |
| e08s06 | SC-e08s06-P0-01, P0-02, P1-01, P1-02, P1-03, P3-01 |
