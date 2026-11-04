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

---

# Security Review — v2.6.0 wave (W3)

> **Date:** 2026-11-04
> **Reviewer:** Worker W3 — `security-review` skill (5-phase scan)
> **Scope:** `git diff 795129b..HEAD` — 13 commits (perf + 12)
> **HEAD:** `9f87c90215a1827549fdb3a58464b6a7ce46d9a9`
> **Diff size:** 37 files, +2792 / −375 (22 Go files; remainder specs/docs)

## Scope

Commits under review:

| Commit | Subject | Security surface |
|--------|---------|------------------|
| `f267814` | perf(release): parallelize independent git reads and github fan-out | concurrency |
| `aaacda5` | fix(httputil): rewind request body on retry | HTTP retry |
| `851d493` | fix(git): roll back remote tag when commit push fails | rollback, shell exec |
| `0177b33` | fix(plugins): clean up half-published releases and parallelize GitLab… | rollback, concurrency |
| `6bb7cdc` | feat(release): make lifecycle status explicit and add rollback hook | rollback orchestration |
| `28592f8` | feat(plugins): implement rollback for git tags and forge releases | rollback |
| `bec2d17`, `9f108aa`, `6d3be4f`, `54da1f8`, `7436412`, `ee009da`, `9f87c90` | chore/spec/test/docs | none |

Security-relevant files:

| Area | Files |
|------|-------|
| Concurrency | `pkg/release/release.go`, `internal/plugins/github_assets.go`, `internal/plugins/github_success.go`, `internal/plugins/gitlab.go` |
| Rollback / cleanup | `internal/plugins/git.go`, `internal/plugins/github.go`, `internal/plugins/gitlab.go`, `internal/git/client.go`, `internal/plugins/registry.go` |
| HTTP retry | `internal/publishers/httputil/retry.go` |
| Token handling | `internal/plugins/github.go`, `internal/plugins/gitlab.go`, `internal/git/client.go` |
| Shell exec | `internal/git/client.go` (`gitCmd` / `runGit`) |
| Path handling | `internal/plugins/github_assets.go`, `internal/plugins/gitlab.go` |

## Method

Five-phase `security-review` scan: scope resolution → context research → vulnerability
assessment (input → sink trace) → false-positive filtering → report. Each candidate was
scored on the 0–10 confidence rubric; anything below 8 is suppressed. CWE mapping applied
where relevant.

Evidence commands and results:

| Check | Command | Result |
|-------|---------|--------|
| Diff scope | `git diff --stat 795129b..HEAD` | 37 files, +2792/−375 |
| Race + tests | `go test -race ./...` | **615 passed in 20 packages** (go1.26.3 darwin/arm64) |
| Token/secret logging | `git diff 795129b..HEAD \| grep -Ei '(token\|secret\|password\|authorization\|bearer\|api[_-]?key).*(log\|printf\|warn\|info\|error)'` | no matches in diff |
| Token/secret logging (repo) | grep of `Logger.*` / `fmt.*` carrying secret terms | only `internal/secure/redact_test.go:45` (synthetic test token) |
| Dangerous sinks | grep for `exec.Command`, `sh -c`, `os.ReadFile`, `InsecureSkipVerify`, `http://`, `PathEscape`, `Sprintf("%s/repos` | reviewed individually; see findings |

## Findings

**No reportable findings at or above the confidence-8 threshold. No HIGH or CRITICAL
findings.** The change set is net-positive for security: it removes half-published
artifacts, makes HTTP retries fail closed, and scopes destructive cleanup to state created
by the current run.

### Suppressed observations (confidence < 8 — not gate-blocking)

| # | File:Line | Category / CWE | Sev | Conf | Observation |
|---|-----------|----------------|-----|------|-------------|
| O1 | `internal/plugins/git.go:301` | Destructive action / integrity | LOW | 6 | `GitPlugin.Rollback` deletes the local and remote tag for `NextRelease.Version` whenever a local tag of that name exists. A pre-existing tag could in theory be removed, but `GetLastRelease` always yields a strictly higher next version, so an existing tag cannot be targeted. Not arbitrary tag deletion. |
| O2 | `internal/plugins/github.go:356,375` | Secrets exposure (defense-in-depth) / CWE-532 | LOW | 5 | `Rollback`/`deleteReleaseByID` trust `GITHUB_REPOSITORY` and print upstream response bodies to stderr via `fmt.Fprintf`. Env is a trusted input (precedent #3); GitHub/GitLab error bodies do not echo `Authorization`. The stderr path bypasses `secure.WrapCore`, but no concrete secret reaches it. |
| O3 | `internal/plugins/gitlab.go:340` | Path handling / CWE-22 | LOW | 4 | Release tag is `neturl.PathEscape`d before URL insertion (positive control). `p.projectID()` is not escaped but is env/config-trusted and pre-existing behavior. |
| O4 | `internal/publishers/httputil/retry.go:37-42` | Latent data race / CWE-362 | LOW | 5 | `RetryClient.Do` lazily writes `rc.MaxRetries`/`rc.RetryBase` on a shared receiver. `NewRetryClient` sets non-zero defaults and publishers run sequentially, so this write path is not reached concurrently today. Re-check if a `RetryClient` is ever shared across the new fan-out. |
| O5 | `internal/publishers/httputil/retry.go:98` | Resource management (excluded #4) | — | — | `rewindRequestBody` replaces `req.Body` without closing the consumed body (fd/GC only, not security). Retrying a non-idempotent POST after a transport error can duplicate a publish — integrity, not confidentiality. |
| O6 | `internal/plugins/github_assets.go:28` | Concurrency correctness | — | — | Bounded fan-out (`maxUploadWorkers = 4`) with `sync.WaitGroup` and a mutex-guarded first error; `go test -race` is clean. `uploadOneAsset` touches no shared mutable plugin state. |

## Positive controls verified

- **Shell exec:** `DeleteRemoteTag` uses `gitCmd(...)` → `exec.Command("git", args...)` variadic. No `sh -c`, no shell interpolation. The tag is passed as one refspec argument prefixed `:refs/tags/`, so it cannot be parsed as an option (`--...`).
- **Token handling:** tokens are read from env only (`GITHUB_TOKEN`, `GITLAB_TOKEN`/`CI_JOB_TOKEN`) and never logged. Zap output passes through `secure.WrapCore` (`cmd/big-release/main.go:100`), which redacts known secret env values and `ghp_`/`gho_`/`npm_` patterns.
- **HTTP retry:** non-rewindable bodies now fail closed (`rewindRequestBody` returns an error) instead of silently resending an empty payload. 429 bodies are drained and closed before retry.
- **Rollback scoping:** GitHub/GitLab rollback deletes only the release remembered from this run (`createdReleaseID` / `createdReleaseTag`, mutex-guarded) and treats 404 as already-deleted (idempotent). Git rollback deletes the remote tag only when `PushTags` actually succeeded (`tagsPushedBy`), so it never removes a tag this run did not push.
- **Path handling:** GitHub asset names are `url.QueryEscape`d; GitLab uses `multipart.CreateFormFile` (Go escapes quotes in the filename header). Asset paths and labels originate from trusted config.
- **Concurrency:** `readGitState` writes three distinct named return values from three goroutines and reads them only after `wg.Wait()` (happens-before); `firstErr` is mutex-guarded. No race reported under `-race`.

## Verdict

**PASS** — 0 HIGH/CRITICAL findings at confidence ≥ 8. No `specs/security/EXCEPTIONS.md`
entry required. Observations O1–O4 are recorded for defense-in-depth and are not
gate-blocking.
