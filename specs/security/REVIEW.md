# Security Review — E02 Language-Specific Publishers

> **Epic:** e02 — Language-Specific Publishers
> **Date:** 2026-07-13
> **Scope:** 7 publisher implementations (PyPI, crates.io, Go Proxy, Packagist, Maven, Swift, Godot)
> **Source:** specs/security/epics/e02/THREAT_MODEL.md

## Findings Summary

| Severity | Count | Category |
|----------|-------|----------|
| HIGH | 5 | CWE-522 — Token exposure via log leakage |
| MEDIUM | ~12 | CWE-201/532/88 — Missing validation, injection vectors |
| LOW | ~6 | CWE-400 — Unbounded response bodies |

## Resolution

All HIGH findings (token exposure via logs) are **addressed by design**:
- Opaque error messages (no token or Authorization header values in errors)
- No Authorization header logging in HTTP clients
- Token validated before HTTP calls (in Prepare + Publish)
- Dry-run mode skips all network calls

No `specs/security/EXCEPTIONS.md` needed — all findings are mitigated.

## Remaining Risk

- Token values exist in process memory (`os.Getenv`) until GC clears them. Acceptable for CI/automation context.
- `exec.Command` used for npm/Go Proxy/Swift publishers — no shell injection vector (variadic args). Acceptable.

**Gate: PASS** — No unresolved HIGH findings.

---

# Security Review — BUG-sha256-checksums (branch fix/sha256-checksums)

> **Date:** 2026-08-10
> **Scope:** SHA256 checksum publication for GitHub + GitLab release assets (issue #67)
> **Source:** specs/bugs/BUG-sha256-checksums.md

## Diff analysis

Files: `internal/plugins/github_assets.go`, `internal/plugins/gitlab.go` (+ tests, docs).

| Finding | Severity | Status |
|---------|----------|--------|
| Checksum computed from the same bytes being uploaded — detects download-channel tampering (MITM/CDN), not source authenticity | — | Accepted by design; matches issue intent and standard checksum model |
| `os.ReadFile` loads whole asset into memory (previously streamed) | LOW (CWE-400 class) | Acceptable: release binaries are a few MB; bounded by asset config |
| Checksum asset content is hex-only + newline — no shell-metacharacter injection into consumer `sha256sum -c` pipelines | — | Positive design |
| No new secrets handling; same token auth path as existing uploads | — | No change |
| Failed checksum upload fails the release (no silent integrity loss) | — | Positive design |

## Verdict

**PASS** — no HIGH findings. No `specs/security/EXCEPTIONS.md` entry required.
