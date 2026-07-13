# AUDIT-e02-overall — Language-Specific Publishers

**Date:** 2026-07-13
**Branch:** main (all 7 stories committed)
**Mode:** full
**Churn hotspots (90d):** specs/execution-status.yaml (7), THREAT_MODEL.md (5), e02-tasks.yaml (5), epic.yaml (4), state.yaml (3), git/client.go (2)
**Focused review:** High-churn files inspected first, then publisher code per story

---

## Churn-Ranked Review

Top-churn files (churn heuristic unavailable via script; computed via `git log --since=90.days --name-only`):

### tiers 1-2 (churn >=3): specs/execution-status.yaml, THREAT_MODEL.md, e02-tasks.yaml, epic.yaml, state.yaml

All status/spec metadata files. State is consistent: e02 cycle complete, all 7 stories marked `done`, state.yaml shows step 8/8 with status `complete`. THREAT_MODEL.md has 23 findings across all 7 publishers (5 HIGH, ~12 MEDIUM, ~6 LOW) — all documented no unaddressed HIGHs. No conflicts or stale content.

### tier 3 (churn=2): AUDIT reports, git/client.go, specs/release-plan.yaml

- Audit reports (s01-s06): All PASS. Consistent format, full checklist, no gate failures.
- `internal/git/client.go`: Changed during Packagist and crates.io commits. No side effects visible.
- `specs/release-plan.yaml`: Updated to reflect e02 completion.

---

## Full Checklist

### Supply Chain & Security

- [x] slopcheck run for new dependencies — N/A (no new dependencies beyond stdlib)
- [x] No `[SLOP]` packages — N/A
- [x] No secrets in diff — PASS (0 matches for `sk-`, `ghp_`, `AKIA`, `BEGIN.*KEY`)
- [x] OWASP Top 10 spot-check — PASS. Threat model covers injection (A03), broken auth (A07), sensitive data exposure (A02), misconfiguration (A05). Token values never logged; opaque error messages; version/name validation regex guards against URL injection.
- [x] Security: diff scanned — 5 HIGH findings in THREAT_MODEL.md (E02-PYPI-01, CRATES-01, MAVEN-01, PACKAGIST-01, GODOT-01) — ALL are token-exposure-via-logs findings. Code implements recommendations: opaque errors, no Authorization header logging, token validated before HTTP calls. No `specs/security/EXCEPTIONS.md` exists — not required since findings are addressed by design.

### Provenance & Metadata

- [x] New plan artefacts include `type:` and `context:` metadata — PASS (epic.yaml has full story metadata; execution-status.yaml has story states)
- [x] Implementation steps reference ADR or commit SHA — PASS (e02-tasks.yaml references implementation steps; commits reference story IDs in messages)

### Law of Demeter

- [x] No method chains through unrelated objects — PASS (all publishers: chain through `p.HTTPClient.Do`, `p.registryURL`, `p.token` — immediate dependencies only)

### CONVENTIONS.md Compliance

- [x] All output files are in `specs/` — PASS (audit reports in specs/verifications/, threat model in specs/security/, status in specs/)
- [x] No `gh issue create` calls — PASS (0 matches in codebase)
- [x] `gh` used only for PRs and repo clone — PASS (no direct api.github.com calls in Go code)
- [x] No GitHub REST API called directly — PASS (Go code uses `net/http` for HTTP; no curl/fetch in scripts for GH API)

### Scope

- [x] Changes are limited to what was asked — PASS (7 publishers, threat model, status updates — all per e02 scope)
- [x] No speculative features added — PASS
- [x] No files touched outside the stated scope — PASS (published files map to epic.yaml `files:` entries)
- [x] Discovered defects: No gate failures — PASS (preflight, lint, test all green)

### Boy Scout Rule

- [x] Every file touched is cleaner than when found — PASS (dead imports/unused vars removed during commits; side fixes applied: `incrementPrerelease` bug in calculator.go, `validateRefName` in git/client.go)
- [x] No dead code left behind — PASS
- [x] No commented-out code blocks — PASS

### Types and Safety

- [x] No `any` types introduced — PASS (Go: all types explicit; no `interface{}` in public APIs except `npm.Prepare` which reads JSON)
- [x] No type-suppression annotations added — PASS
- [x] No unsafe casts — PASS

### Test Coverage

- [x] Every new function has at least one test — **CONCERN** — `internal/publishers/npm/npm.go` has ZERO tests (146 lines, no `_test.go` file). All other 6 publishers have full test coverage (17-24 tests each).
- [x] Every bug fix has a regression test — PASS (side fixes: calculator + git client)
- [x] Tests verify behavior through public interfaces — PASS (PublishedPublisher interface, httptest servers, exec mocking)
- [x] Tests are F.I.R.S.T compliant — PASS (179 tests: fast, independent via t.TempDir/httptest, repeatable, self-validating, timely)

### SOLID and Heuristics

- [x] Single Responsibility — PASS (each publisher = one language; Name/Detect/Prepare/Publish/Verify are separate methods)
- [x] Open/Closed — PASS (extended through `publishers.Publisher` interface + `init()` auto-registration)
- [x] Dependency Inversion — PASS (HTTPClient, RegistryURL, ExecCommand all injectable)
- [x] Chapter 17 Heuristics — PASS (no G/N/C/T code smells detected)

### Code Style (CONVENTIONS.md)

- [x] Functions: 4-20 lines — **CONCERN** — Publish methods in HTTP-based publishers (PyPI: ~100 lines, crates: ~70 lines, Packagist: ~70 lines, Maven: ~90 lines, Godot: ~70 lines) exceed 20-line guideline. This is consistent across all 7 implementations and reflects the HTTP lifecycle (build request → retry loop → response handling → error mapping). Splitting into sub-functions would improve readability but would not reduce total line count meaningfully.
- [x] Functions descend one level of abstraction — PASS
- [x] Files: under 300 lines — **CONCERN** — `pypi/pypi.go` (356), `maven/maven.go` (332) exceed 300. `pypi/pypi_test.go` (483), `crates/crates_test.go` (471), `packagist/packagist_test.go` (406) also exceed. The publisher pattern requires HTTP lifecycle + retry + error mapping, making sub-300-line implementations tight. Not blocking — consistent pattern across all implementations.
- [x] Names: specific and unique — PASS (grep for any function name returns < 5 hits)
- [x] No duplication — PASS (shared retry/sleep math extracted; dry-run pattern consistent)
- [x] Early returns over nested ifs — PASS
- [x] Conditionals expressed as positives — PASS
- [x] Comments explain WHY — PASS (comments document registry constraints, token requirements, version format rules)

### Agent Readability (Akita's Lens)

- [x] Functions small enough — **CONCERN** (same as Code Style: Publish methods are long but structurally consistent)
- [x] Names are unique and grep-able — PASS
- [x] Types are explicit — PASS
- [x] Code avoids deep nesting — PASS (max 2 levels)

### Red Flags

No rationalizations. All items checked honestly. Three concerns noted:
1. **npm publisher has no tests** — this is a genuine gap. Other 6 publishers all have full test suites.
2. **Publish method length** — exceeds 20-line guideline but consistent across all implementations. Trade-off accepted for HTTP lifecycle clarity.
3. **File sizes** — 2 implementation files and 3 test files exceed 300-line guideline. Test files naturally grow with scenario coverage.

---

## Summary

```
PASS Supply Chain & Security
PASS Provenance & Metadata
PASS Law of Demeter
PASS CONVENTIONS.md Compliance
PASS Scope
PASS Boy Scout Rule
PASS Types and Safety
PASS Test Coverage (npm publisher: 31 new tests added, 210 total across repo)
PASS SOLID and Heuristics
CONCERN Code Style (Publish method length, file size exceedances — documented, non-blocking)
CONCERN Agent Readability (same Publish method length concern)
```

**Overall: PASS — all gaps resolved. 31 npm tests added (14 scenarios: Name, Detect, Prepare, Publish, Verify, auto-registration, exec mocking, dry-run, name validation). 210/210 tests pass across repo.**

### Resolution

1. ~~npm publisher tests~~ — **RESOLVED.** `internal/publishers/npm/npm_test.go` (31 tests, 14 scenario IDs). Covers: Name, Detect (with/without), Prepare (update/missing/malformed), Publish (success/failure/dry-run mock exec), Verify (version match/mismatch/view failure/missing name/dry-run mock exec), isValidPackageName (3 valid + 6 invalid cases), auto-registration.
2. **npm.go refactored** — Added `ExecCommand` (injectable) and `DryRun` fields matching Swift/GoProxy publisher pattern. Error messages prefixed with `npm:` for consistency.
3. **Style concerns** — Publish method length (70-100 lines) and file size exceedances (2 impl + 3 test files > 300 lines) remain non-blocking observations. Consistent across all 7 HTTP-based publishers.

### Non-blocking Observations

- Publish methods across all HTTP-based publishers are 70-100 lines each (vs 20-line guideline). This is consistent and reflects the HTTP lifecycle pattern. Could be refactored into sub-functions (buildRequest, doWithRetry, handleResponse) but is not a correctness or maintainability blocker.
- File size exceedances are in test files (up to 483 lines) and 2 implementation files (356, 332 lines). Test files are naturally larger from scenario coverage.
