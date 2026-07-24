# Threat Model — E08 Semantic Release Test Parity

> **Epic:** e08 — Semantic Release Test Parity  
> **Date:** 2026-07-24  
> **Reviewer:** security-review skill (build-epic Step 0)  
> **Scope:** Secret masking/redaction; git AuthURL token injection; temp-git E2E fixtures; commit traversal; AggregateError; config file/CLI load  
> **Confidence threshold:** 8/10 (per skill mandate)

---

## Executive Summary

| Severity | Count | CWE Coverage |
|----------|-------|--------------|
| HIGH     | 0     | — |
| MEDIUM   | 2     | CWE-532, CWE-598 / CWE-200 |
| LOW      | 2     | CWE-377 / CWE-459, CWE-532 |

Overall residual risk after mitigations: **LOW**. Epic primarily closes test gaps and lands minimal parity helpers; the security-relevant surface is credential handling in logs/errors and in-memory URL construction — not new network publish paths.

---

## Surface Area

| Component | Trust boundary | Sensitive data |
|-----------|----------------|----------------|
| `internal/secure` Redact / RedactKnownSecrets | Process logs ↔ stdout/stderr / notes | `GH_TOKEN`, `GITHUB_TOKEN`, `NPM_TOKEN`, registry tokens in env |
| zap core wrapper (`cmd/big-release`) | Logger → operator / CI logs | Any field values that may embed secrets |
| Notes generator (`internal/algorithm/generator.go`) | Release notes output | Commit subjects/bodies that might echo tokens (defense in depth) |
| `internal/git.AuthURL` | In-memory URL string for git remotes | HTTPS credentials embedded as `https://x-access-token:<token>@host/...` |
| `internal/git/testrepo` temp repos | Local filesystem (test only) | Test tokens / scrubbed env; no production remotes |
| Config load (`internal/config`) | Filesystem YAML/JSON + CLI flags | Paths and dry-run flags only (no JS execution) |
| AggregateError (`pkg/release`) | Error return values | Plugin error strings that may include remote URLs |

### Data Flow

```
env tokens (GH_TOKEN / NPM_TOKEN / …)
        │
        ├──────────────────────────────┐
        ▼                              ▼
┌───────────────────┐         ┌────────────────────┐
│ zap + secure      │         │ AuthURL(remote,tok)│
│ RedactKnownSecrets│         │ HTTPS inject only  │
│ notes Redact      │         │ SSH passthrough    │
└─────────┬─────────┘         └─────────┬──────────┘
          │                             │
          ▼                             ▼
   CI logs / notes              git fetch/push argv
                                (never log raw URL)
```

---

## FINDING E08-MASK-01 — MEDIUM (Confidence 9/10)

### Credential leakage via unstructured logs and release notes

**File:** `internal/secure/*`, `cmd/big-release/main.go`, `internal/algorithm/generator.go`  
**Severity:** MEDIUM  
**Category:** Insertion of Sensitive Information into Log File  
**CWE:** CWE-532  

#### Description

e08s01 introduces/extends secret masking so token values never appear in zap logs or generated notes. Without redaction, env-backed tokens can leak into CI artifacts, GitHub Actions logs, and published release notes when errors or debug fields include URLs, headers, or env dumps.

#### Mitigation (required in implementation)

- Centralize `Redact` / `RedactKnownSecrets` in `internal/secure`.
- Wrap zap core so all log fields pass through redaction.
- Generator notes path delegates string scrubbing to `secure`.
- Tests assert known token fixtures never appear in captured log/note output (`TestMasking`).

#### Status

**Accepted with mitigation** — primary deliverable of e08s01; gate via SC-e08s01-P0 scenarios.

---

## FINDING E08-AUTH-01 — MEDIUM (Confidence 9/10)

### Token embedded in HTTPS remote URL (in-memory / error paths)

**File:** `internal/git/auth_url.go`  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Query String / Credential in URI  
**CWE:** CWE-598 / CWE-200  

#### Description

e08s02 `AuthURL(remoteURL, token)` injects credentials into HTTPS remotes for git operations (semantic-release parity). The constructed URL contains the raw token. If that string is logged, wrapped into `fmt.Errorf("%v", url)`, or written to process listings, the token leaks. SSH remotes must remain untouched. Double-injection must not occur.

#### Mitigation (required in implementation)

- Pure function; do not mutate `remote.origin.url` on disk (out of scope / forbidden).
- HTTPS only inject; SSH and non-HTTP schemes passthrough.
- Idempotent: skip inject if credentials already present.
- All error messages that may include the URL MUST pass through `secure` redaction before return/log.
- Unit tests cover inject / passthrough / no-double-inject / redacted errors (`TestGitAuthURL`).

#### Status

**Accepted with mitigation** — primary deliverable of e08s02; residual risk is operator misuse of debug logging outside redacted paths.

---

## FINDING E08-TEMP-01 — LOW (Confidence 8/10)

### Temp-git fixture directories and env leakage in tests

**File:** `internal/git/testrepo/` (and consumers e08s03 / e08s04)  
**Severity:** LOW  
**Category:** Insecure Temporary File / Incomplete Cleanup  
**CWE:** CWE-377 / CWE-459  

#### Description

E2E and traversal tests create temporary git repositories. Residual risk is (a) leftover dirs containing crafted remotes or token-like fixtures, and (b) leaking host `GH_TOKEN` into child git processes if env is not scrubbed.

#### Mitigation

- Use `t.TempDir()` (auto-cleanup) for all fixtures.
- Provide `ScrubEnv` helper that clears/overrides credential env vars for subprocesses.
- Never point test remotes at production hosts; dry-run / `CI=true` for pipeline E2E.
- Do not commit token fixtures into the repo; generate ephemeral strings in-test.

#### Status

**Accepted with mitigation** — LOW residual under CI trust model.

---

## FINDING E08-ERR-01 — LOW (Confidence 8/10)

### AggregateError may concatenate sensitive plugin messages

**File:** `pkg/release/errors.go`, `pkg/release/release.go`  
**Severity:** LOW  
**Category:** Insertion of Sensitive Information into Log File  
**CWE:** CWE-532  

#### Description

e08s05 collects multiple VerifyConditions failures. If individual plugin errors embed AuthURL or registry responses, joining them can amplify leakage in a single returned error.

#### Mitigation

- Prefer redacted plugin errors at source (e08s01/e08s02).
- AggregateError tests use stub messages without real secrets.
- Do not expand aggregation beyond VerifyConditions (out of scope).

#### Status

**Accepted with mitigation** — residual LOW if MASK/AUTH mitigations land first (Wave 1 before Wave 2/3).

---

## Residual Risk

| Story | Risk after mitigations |
|-------|------------------------|
| e08s01 | LOW — masking is the control; tests prove absence of tokens |
| e08s02 | LOW–MEDIUM — token-in-URL by design; must never log raw |
| e08s03 | LOW — temp-git + scrubbed env + dry-run |
| e08s04 | LOW — local git traversal only |
| e08s05 | LOW — AggregateError; depends on redacted sources |
| e08s06 | NONE–LOW — YAML/JSON/CLI only; no JS execution |

**Gate:** Proceed. No HIGH findings ≥ 8 confidence unresolved. Threat tags feed plan-tests risk tiers and plan-work `security:` fields.
