# Threat Model — E02 Language-Specific Publishers

> **Epic:** e02 — Language-Specific Publishers  
> **Date:** 2026-07-13  
> **Reviewer:** security-review skill (build-epic Step 0)  
> **Scope:** 7 publisher implementations (PyPI, crates.io, Go Proxy, Packagist, Maven, Swift, Godot)  
> **Confidence threshold:** 8/10 (per skill mandate)

---

## Executive Summary

| Severity | Count | CWE Coverage |
|----------|-------|--------------|
<<<<<<< Updated upstream
| HIGH     | 3     | CWE-522     |
| MEDIUM   | 8     | CWE-201, CWE-532, CWE-88 |
| LOW      | 5     | CWE-400     |
=======
| HIGH     | 4     | CWE-522     |
| MEDIUM   | 10    | CWE-201, CWE-532, CWE-88 |
| LOW      | 6     | CWE-400     |
>>>>>>> Stashed changes

---

## e02s01 — PyPI Publisher

> **Risk Tier:** P1  
> **Token:** `PYPI_TOKEN` (env var, Bearer token in HTTP header)  
> **API Endpoints:**
> - Publish: `POST https://upload.pypi.org/legacy/`
> - Verify: `GET https://pypi.org/pypi/<package>/json`

### Data Flow

```
PYPI_TOKEN (env var)
    │
    ▼
┌─────────────────────────┐
│  PyPI Publisher          │
│                         │
│  Publish():              │
│    POST /legacy/         │
│    Authorization: token  │
│    <PYPI_TOKEN>          │
│                         │
│  Verify():               │
│    GET /pypi/<pkg>/json  │
└──────┬──────────────────┘
       │
       ▼
   PyPI Registry (upload.pypi.org / pypi.org)
```

**Trust boundary:** The `PYPI_TOKEN` environment variable is set outside the Go process (CI secrets or shell env). The publisher reads it at publish time.

---

## FINDING E02-PYPI-01 — HIGH (Confidence 9/10)

### Token exposure via environment dump or log leakage

**File:** `internal/publishers/pypi/pypi.go` (Publish)  
**Severity:** HIGH  
**Category:** Insufficiently Protected Credentials  
**CWE:** CWE-522

#### Description

The `PYPI_TOKEN` is read from the environment using `os.Getenv("PYPI_TOKEN")`. If the token value is ever logged (e.g., in a debug log, error message, or panic recovery), it will be exposed. The token is also present in the process's `/proc/self/environ` and in any crash dump.

Additionally, if the HTTP client is configured to log request headers (common during debugging), the `Authorization: token <PYPI_TOKEN>` header would be written to logs verbatim.

#### Exploit Scenario

1. A developer runs `big-release --verbose release` and a transient network error occurs.
2. The error handler logs the full HTTP request including headers: `Authorization: token pypi-xxxxx`.
3. The log is uploaded to a log aggregation service (e.g., Datadog, CloudWatch) with lesser access controls than the CI secrets manager.
4. An attacker with read access to logs extracts the PyPI token.

#### Recommendation

1. **Never log the token value.** If logging the Authorization header is necessary, replace the token value with `[REDACTED]`:
   ```go
   func redactAuthHeader(req *http.Request) string {
       h := req.Header.Get("Authorization")
       if h != "" {
           req.Header.Set("Authorization", "token [REDACTED]")
       }
       return fmt.Sprintf("%s %s", req.Method, req.URL)
   }
   ```

2. **Return opaque errors** from Publish — don't include `os.Getenv` output or raw `Authorization` header values:
   ```go
   // Good:
   return fmt.Errorf("pypi: publish failed (HTTP %d)", resp.StatusCode)
   
   // Bad:
   return fmt.Errorf("pypi: publish failed with token %s: HTTP %d", token, resp.StatusCode)
   ```

3. **Consider using `os.Getenv` with immediate zeroing** — since Go strings are immutable, the token will remain in memory until garbage collected. For high-security deployments, consider using `memguard` or similar for sensitive strings.

---

## FINDING E02-PYPI-02 — MEDIUM (Confidence 8/10)

### Missing token validation before HTTP call

**File:** `internal/publishers/pypi/pypi.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Sent Data  
**CWE:** CWE-201

#### Description

`Publish(version)` checks for an empty `PYPI_TOKEN` and returns an error before making HTTP calls. However, the check happens at the `Publish` call boundary only. If `Prepare()` is called first and succeeds (updating version in config files), and then `Publish()` fails with an empty token error, the config files have already been modified — the system is left in a dirty state.

#### Exploit Scenario

1. `PYPI_TOKEN` is accidentally unset in CI.
2. `Prepare("2.0.0")` runs successfully, updating `setup.cfg` version.
3. `Publish("2.0.0")` returns error: "PYPI_TOKEN is empty".
4. CI pipeline shows "release failed" but the working tree has already been modified.
5. Next `git status` shows dirty tree; re-running may double-commit the version bump.

#### Recommendation

1. **Pre-validate the token at construction time** — call `os.Getenv("PYPI_TOKEN")` in `NewPublisher()` and store the result. Validate it in `Prepare()` (before mutating files) and in `Publish()`.

2. **Or, validate in `Prepare()`** so that version bump doesn't happen without a valid token:

   ```go
   func (p *Publisher) Prepare(version string) error {
       if p.token == "" {
           return fmt.Errorf("pypi: PYPI_TOKEN is empty, cannot prepare for publish")
       }
       // ... version bump logic ...
   }
   ```

---

## FINDING E02-PYPI-03 — MEDIUM (Confidence 7/10)

### Package name and version injection in Verify URL

**File:** `internal/publishers/pypi/pypi.go` (Verify)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Query Strings in Request / Injection  
**CWE:** CWE-532 / CWE-88

#### Description

`Verify(version)` constructs a URL from the package name and version:

```go
url := fmt.Sprintf("https://pypi.org/pypi/%s/json", pkgName)
```

If the package name contains special characters (e.g., `../../`, newlines, URL-unsafe characters), the URL could resolve to an unexpected endpoint. While PyPI enforces valid package name rules on their side, a malicious `pyproject.toml` or `setup.cfg` could contain a crafted `name` field.

#### Exploit Scenario

1. A repository has a `pyproject.toml` with `name = "../../config/special"`.
2. When `Verify` runs, it calls `GET https://pypi.org/pypi/../../config/special/json`.
3. Depending on the HTTP client and server, this might resolve unexpectedly or leak information about the request.

#### Recommendation

Validate the package name and version before using them in URL construction:

```go
var validPyPIPackageName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

func isValidPyPIPackageName(name string) bool {
    return len(name) >= 1 && len(name) <= 200 && validPyPIPackageName.MatchString(name)
}
```

---

## FINDING E02-PYPI-04 — LOW (Confidence 5/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/pypi/pypi.go` (Publish, Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Publish` and `Verify` methods read HTTP response bodies without a size limit. If PyPI returns an unexpectedly large response (e.g., a 5xx error page with megabytes of HTML), the publisher will allocate memory to hold the full response body. While unlikely in practice, this could be triggered by a misconfigured proxy or CDN returning a large error page.

#### Recommendation

Use `http.MaxBytesReader` or `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 10 * 1024 * 1024 // 10 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("pypi: response body too large")
}
```

---

## Risk Matrix

| Finding | Severity | Confidence | CWE | Fix Effort |
|---------|----------|------------|-----|------------|
| Token exposure via log leakage (E02-PYPI-01) | HIGH | 9/10 | CWE-522 | Small (redact logs, opaque errors) |
| Missing token validation before HTTP call (E02-PYPI-02) | MEDIUM | 8/10 | CWE-201 | Small (validate in Prepare) |
| Package name injection in Verify URL (E02-PYPI-03) | MEDIUM | 7/10 | CWE-532 | Small (validate name format) |
| Unbounded response body (E02-PYPI-04) | LOW | 5/10 | CWE-400 | Trivial (add LimitReader) |
| Token exposure via env/logs (E02-CRATES-01) | HIGH | 9/10 | CWE-522 | Small (redact logs, opaque errors) |
| Missing token validation before HTTP call (E02-CRATES-02) | MEDIUM | 8/10 | CWE-201 | Small (validate in Prepare) |
| Package name injection in Verify URL (E02-CRATES-03) | MEDIUM | 7/10 | CWE-532 | Small (validate name format) |
| Unbounded response body (E02-CRATES-04) | LOW | 5/10 | CWE-400 | Trivial (add LimitReader) |
|| Git tag injection via version string (E02-GOPROXY-01) | MEDIUM | 7/10 | CWE-88 | Small (version validation) |
<<<<<<< Updated upstream
|| Unbounded response body (E02-GOPROXY-02) | LOW | 6/10 | CWE-400 | Trivial (add LimitReader) |
|| Module name injection in Verify URL (E02-GOPROXY-03) | LOW | 5/10 | CWE-88 | Small (validate module path) |
||| Token exposure via env/logs (E02-MAVEN-01) | HIGH | 9/10 | CWE-522 | Small (redact logs, opaque errors) |
||| Missing token validation before HTTP call (E02-MAVEN-02) | MEDIUM | 8/10 | CWE-201 | Small (validate in Prepare) |
||| POM group/artifact injection in Verify URL (E02-MAVEN-03) | MEDIUM | 7/10 | CWE-88 | Small (validate group/artifact/version) |
||| Unbounded response body (E02-MAVEN-04) | LOW | 5/10 | CWE-400 | Trivial (add LimitReader) |
=======
|| Token exposure via env/logs (E02-PACKAGIST-01) | HIGH | 9/10 | CWE-522 | Small (redact logs, opaque errors) |
|| Missing token validation before HTTP call (E02-PACKAGIST-02) | MEDIUM | 8/10 | CWE-201 | Small (validate in Prepare) |
|| Package name injection in Verify URL (E02-PACKAGIST-03) | MEDIUM | 7/10 | CWE-532 | Small (validate name format) |
|| Unbounded response body (E02-PACKAGIST-04) | LOW | 5/10 | CWE-400 | Trivial (add LimitReader) |
|| Git tag injection via version string (E02-SWIFT-01) | MEDIUM | 7/10 | CWE-88 | Small (version validation) |
>>>>>>> Stashed changes
|| Token exposure via env/logs (E02-GODOT-01) | HIGH | 9/10 | CWE-522 | Small (redact logs, opaque errors) |
|| Missing token validation before HTTP call (E02-GODOT-02) | MEDIUM | 8/10 | CWE-201 | Small (validate in Prepare) |
|| Owner/repo injection in GitHub API URL (E02-GODOT-03) | MEDIUM | 7/10 | CWE-88 | Small (validate owner/repo) |
|| Unbounded response body (E02-GODOT-04) | LOW | 5/10 | CWE-400 | Trivial (add LimitReader) |
<<<<<<< Updated upstream
=======

>>>>>>> Stashed changes

---

## e02s02 — crates.io Publisher

> **Risk Tier:** P1  
> **Token:** `CARGO_TOKEN` (env var, Bearer token in HTTP header)  
> **API Endpoints:**
> - Publish: `PUT https://crates.io/api/v1/crates/new`
> - Verify: `GET https://crates.io/api/v1/crates/<name>/versions`

### Data Flow

```
CARGO_TOKEN (env var)
    │
    ▼
┌──────────────────────────┐
│  crates.io Publisher      │
│                          │
│  Publish():               │
│    PUT /api/v1/crates/new │
│    Authorization: <token> │
│                          │
│  Verify():                │
│    GET /api/v1/crates/    │
│      <name>/versions      │
└──────┬───────────────────┘
       │
       ▼
   crates.io Registry (crates.io)
```

**Trust boundary:** The `CARGO_TOKEN` environment variable is set outside the Go process (CI secrets or shell env). The publisher reads it at publish time.

---

## FINDING E02-CRATES-01 — HIGH (Confidence 9/10)

### Token exposure via environment dump or log leakage

**File:** `internal/publishers/crates/crates.go` (Publish)  
**Severity:** HIGH  
**Category:** Insufficiently Protected Credentials  
**CWE:** CWE-522

#### Description

The `CARGO_TOKEN` is read from the environment using `os.Getenv("CARGO_TOKEN")`. If the token value is ever logged (e.g., in a debug log, error message, or panic recovery), it will be exposed. The token is also present in the process's `/proc/self/environ` and in any crash dump.

If the HTTP client is configured to log request headers (common during debugging), the `Authorization: <CARGO_TOKEN>` header would be written to logs verbatim.

#### Exploit Scenario

1. A developer runs `big-release --verbose release` and a transient network error occurs.
2. The error handler logs the full HTTP request including headers: `Authorization: cargo-xxxxx`.
3. The log is uploaded to a log aggregation service (e.g., Datadog, CloudWatch) with lesser access controls than the CI secrets manager.
4. An attacker with read access to logs extracts the crates.io token.

#### Recommendation

1. **Never log the token value.** If logging the Authorization header is necessary, replace the token value with `[REDACTED]`:
   ```go
   func redactAuthHeader(req *http.Request) string {
       h := req.Header.Get("Authorization")
       if h != "" {
           req.Header.Set("Authorization", "Bearer [REDACTED]")
       }
       return fmt.Sprintf("%s %s", req.Method, req.URL)
   }
   ```

2. **Return opaque errors** from Publish — don't include `os.Getenv` output or raw `Authorization` header values:
   ```go
   // Good:
   return fmt.Errorf("crates: publish failed (HTTP %d)", resp.StatusCode)
   
   // Bad:
   return fmt.Errorf("crates: publish failed with token %s: HTTP %d", token, resp.StatusCode)
   ```

3. **Consider using `os.Getenv` with immediate zeroing** — since Go strings are immutable, the token will remain in memory until garbage collected. For high-security deployments, consider using `memguard` or similar for sensitive strings.

---

## FINDING E02-CRATES-02 — MEDIUM (Confidence 8/10)

### Missing token validation before HTTP call

**File:** `internal/publishers/crates/crates.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Sent Data  
**CWE:** CWE-201

#### Description

`Publish(version)` checks for an empty `CARGO_TOKEN` and returns an error before making HTTP calls. However, the check happens at the `Publish` call boundary only. If `Prepare()` is called first and succeeds (updating version in Cargo.toml), and then `Publish()` fails with an empty token error, the config files have already been modified — the system is left in a dirty state.

#### Exploit Scenario

1. `CARGO_TOKEN` is accidentally unset in CI.
2. `Prepare("2.0.0")` runs successfully, updating `Cargo.toml` version.
3. `Publish("2.0.0")` returns error: "CARGO_TOKEN is empty".
4. CI pipeline shows "release failed" but the working tree has already been modified.
5. Next `git status` shows dirty tree; re-running may double-commit the version bump.

#### Recommendation

1. **Pre-validate the token at construction time** — call `os.Getenv("CARGO_TOKEN")` in `NewPublisher()` and store the result. Validate it in `Prepare()` (before mutating files) and in `Publish()`.

2. **Or, validate in `Prepare()`** so that version bump doesn't happen without a valid token:

   ```go
   func (p *Publisher) Prepare(version string) error {
       if os.Getenv("CARGO_TOKEN") == "" {
           return fmt.Errorf("crates: CARGO_TOKEN is empty, cannot prepare for publish")
       }
       // ... version bump logic ...
   }
   ```

---

## FINDING E02-CRATES-03 — MEDIUM (Confidence 7/10)

### Package name and version injection in Verify URL

**File:** `internal/publishers/crates/crates.go` (Verify)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Query Strings in Request / Injection  
**CWE:** CWE-532 / CWE-88

#### Description

`Verify(version)` constructs a URL from the package name:

```go
url := fmt.Sprintf("https://crates.io/api/v1/crates/%s/versions", pkgName)
```

If the package name contains special characters (e.g., `../../`, newlines, URL-unsafe characters), the URL could resolve to an unexpected endpoint. While crates.io enforces valid package name rules on their side, a malicious `Cargo.toml` could contain a crafted `name` field.

#### Exploit Scenario

1. A repository has a `Cargo.toml` with `name = "../../config/special"`.
2. When `Verify` runs, it calls `GET https://crates.io/api/v1/crates/../../config/special/versions`.
3. Depending on the HTTP client and server, this might resolve unexpectedly or leak information about the request.

#### Recommendation

Validate the package name before using it in URL construction:

```go
var validCrateName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func isValidCrateName(name string) bool {
    return len(name) >= 1 && len(name) <= 128 && validCrateName.MatchString(name)
}
```

---

## FINDING E02-CRATES-04 — LOW (Confidence 5/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/crates/crates.go` (Publish, Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Publish` and `Verify` methods read HTTP response bodies without a size limit. If crates.io returns an unexpectedly large response (e.g., a 5xx error page with megabytes of HTML), the publisher will allocate memory to hold the full response body. While unlikely in practice, this could be triggered by a misconfigured proxy or CDN returning a large error page.

#### Recommendation

Use `http.MaxBytesReader` or `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 10 * 1024 * 1024 // 10 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("crates: response body too large")
}
```

---

## e02s03 — Go Proxy Publisher

> **Risk Tier:** P2  
> **Token:** None (token-less — tag-based versioning only)  
> **API Endpoints:**
> - Publish: `git tag <version> && git push origin <tag>`, then `go list -m <module>@<version>` (poll Go proxy)
> - Verify: `GET https://proxy.golang.org/<module>/@v/<version>.info`

### Data Flow

```
┌──────────────────────────────┐
│  Go Proxy Publisher           │
│                              │
│  Publish():                   │
│    git tag v<version>        │
│    git push origin v<version>│
│    go list -m <mod>@<version>│
│      (GOPROXY env override)  │
│                              │
│  Verify():                    │
│    GET /<module>/@v/<ver>.info│
└──────┬───────────────────────┘
       │
       ├── os/exec: git tag, git push, go list -m
       │
       ▼
   Go Module Mirror (proxy.golang.org)
```

**Trust boundary:** No credentials are handled. The publisher relies on the user's local Git credentials (`git push`) and the Go module mirror's public API (`go list -m`, `proxy.golang.org`). The GOPROXY environment variable may be overridden but contains only a URL, not a secret.

---

## FINDING E02-GOPROXY-01 — MEDIUM (Confidence 7/10)

### Git tag injection via version string

**File:** `internal/publishers/goproxy/goproxy.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Argument Injection  
**CWE:** CWE-88

#### Description

`Publish(version)` constructs a git tag as `v<version>` and passes it to `git tag` and `git push` via `os/exec`. If the version string contains shell metacharacters (e.g., `; rm -rf /`, `$(malicious)`, `` `backtick` ``), and the command is executed via a shell, arbitrary command injection is possible.

Even with direct `os/exec.Command` (no shell), characters like newlines or flag injection (`--delete`) in the version string could result in unexpected tag operations.

#### Exploit Scenario

1. An attacker-controlled input (e.g., from commit message parsing or changelog generation) produces a version string like `1.0.0; rm -rf /`.
2. If passed to `exec.Command("sh", "-c", ...)`, the shell interprets the semicolon and executes the destructive command.
3. If passed to `exec.Command("git", "tag", tagName)`, the `exec.Command` variant is safe from shell injection, but a version like `--delete` or `-d` could be interpreted as a git flag.

#### Recommendation

1. **Always use `exec.Command` (variadic args) instead of `exec.Command("sh", "-c", ...)`.** Go's `os/exec.Command` with separate arguments does not invoke a shell, preventing shell injection.
2. **Validate version string** before constructing the tag — reject characters beyond `[a-zA-Z0-9._-]`:
   ```go
   var validVersion = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

   func isValidVersion(v string) bool {
       return len(v) > 0 && len(v) <= 128 && validVersion.MatchString(v)
   }
   ```
3. **Validate version in Prepare()** as a defense-in-depth measure — all publishers accept `version` as input, so a shared validation helper would catch injection at the earliest point.
4. **Do not echo or log the version string verbatim** in error messages that could appear in a pipeline log; return opaque errors.

---

## FINDING E02-GOPROXY-02 — LOW (Confidence 6/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/goproxy/goproxy.go` (Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Verify` method reads the Go proxy's `.info` endpoint response body without a size limit. While Go proxy responses are typically small (a few KB of JSON), a misconfigured proxy or CDN could return a large response, causing excessive memory allocation.

#### Recommendation

Use `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 1 * 1024 * 1024 // 1 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("goproxy: response body too large")
}
```

---

## FINDING E02-GOPROXY-03 — LOW (Confidence 5/10)

### Module name injection in Verify URL

**File:** `internal/publishers/goproxy/goproxy.go` (Verify)  
**Severity:** LOW  
**Category:** Injection  
**CWE:** CWE-88

#### Description

`Verify(version)` constructs a URL from the module name:
```
/<module>/@v/<version>.info
```
If the module name (read from `go.mod`) contains special characters, the URL could resolve to an unexpected endpoint. Go module paths are validated by the Go toolchain, but a malicious or malformed `go.mod` from a fork could contain a crafted module path.

#### Recommendation

Validate the module name against Go module path rules before using it in URL construction:

```go
var validGoModPath = regexp.MustCompile(`^[a-zA-Z0-9./_-]+$`)

func isValidModulePath(path string) bool {
    return len(path) > 0 && len(path) <= 256 && validGoModPath.MatchString(path)
}
```

---

<<<<<<< Updated upstream
=======
## e02s04 — Packagist Publisher

> **Risk Tier:** P3  
> **Token:** `PACKAGIST_TOKEN` (env var, Bearer token in HTTP header)  
> **API Endpoints:**
> - Publish: `POST https://packagist.org/api/update-package`
> - Verify: `GET https://packagist.org/packages/<vendor>/<package>.json`

### Data Flow

```
PACKAGIST_TOKEN (env var)
    │
    ▼
┌───────────────────────────────┐
│  Packagist Publisher           │
│                               │
│  Publish():                    │
│    POST /api/update-package   │
│    Authorization: token        │
│    <PACKAGIST_TOKEN>          │
│                               │
│  Verify():                     │
│    GET /packages/              │
│      <vendor>/<package>.json   │
└──────┬────────────────────────┘
       │
       ▼
   Packagist Registry (packagist.org)
```

**Trust boundary:** The `PACKAGIST_TOKEN` environment variable is set outside the Go process (CI secrets or shell env). The publisher reads it at publish time.

---

## FINDING E02-PACKAGIST-01 — HIGH (Confidence 9/10)

### Token exposure via environment dump or log leakage

**File:** `internal/publishers/packagist/packagist.go` (Publish)  
**Severity:** HIGH  
**Category:** Insufficiently Protected Credentials  
**CWE:** CWE-522

#### Description

The `PACKAGIST_TOKEN` is read from the environment using `os.Getenv("PACKAGIST_TOKEN")`. If the token value is ever logged (e.g., in a debug log, error message, or panic recovery), it will be exposed. The token is also present in the process's `/proc/self/environ` and in any crash dump.

Additionally, if the HTTP client is configured to log request headers (common during debugging), the `Authorization: token <PACKAGIST_TOKEN>` header would be written to logs verbatim.

#### Exploit Scenario

1. A developer runs `big-release --verbose release` and a transient network error occurs.
2. The error handler logs the full HTTP request including headers: `Authorization: token packagist-xxxxx`.
3. The log is uploaded to a log aggregation service (e.g., Datadog, CloudWatch) with lesser access controls than the CI secrets manager.
4. An attacker with read access to logs extracts the Packagist token.

#### Recommendation

1. **Never log the token value.** If logging the Authorization header is necessary, replace the token value with `[REDACTED]`:
   ```go
   func redactAuthHeader(req *http.Request) string {
       h := req.Header.Get("Authorization")
       if h != "" {
           req.Header.Set("Authorization", "token [REDACTED]")
       }
       return fmt.Sprintf("%s %s", req.Method, req.URL)
   }
   ```

2. **Return opaque errors** from Publish — don't include `os.Getenv` output or raw `Authorization` header values:
   ```go
   // Good:
   return fmt.Errorf("packagist: publish failed (HTTP %d)", resp.StatusCode)

   // Bad:
   return fmt.Errorf("packagist: publish failed with token %s: HTTP %d", token, resp.StatusCode)
   ```

3. **Consider using `os.Getenv` with immediate zeroing** — since Go strings are immutable, the token will remain in memory until garbage collected. For high-security deployments, consider using `memguard` or similar for sensitive strings.

---

## FINDING E02-PACKAGIST-02 — MEDIUM (Confidence 8/10)

### Missing token validation before HTTP call

**File:** `internal/publishers/packagist/packagist.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Sent Data  
**CWE:** CWE-201

#### Description

`Publish(version)` checks for an empty `PACKAGIST_TOKEN` and returns an error before making HTTP calls. However, the check happens at the `Publish` call boundary only. If `Prepare()` is called first and succeeds (updating version in composer.json), and then `Publish()` fails with an empty token error, the config files have already been modified — the system is left in a dirty state.

#### Exploit Scenario

1. `PACKAGIST_TOKEN` is accidentally unset in CI.
2. `Prepare("2.0.0")` runs successfully, updating `composer.json` version.
3. `Publish("2.0.0")` returns error: "PACKAGIST_TOKEN is empty".
4. CI pipeline shows "release failed" but the working tree has already been modified.
5. Next `git status` shows dirty tree; re-running may double-commit the version bump.

#### Recommendation

1. **Pre-validate the token at construction time** — call `os.Getenv("PACKAGIST_TOKEN")` in `NewPublisher()` and store the result. Validate it in `Prepare()` (before mutating files) and in `Publish()`.

2. **Or, validate in `Prepare()`** so that version bump doesn't happen without a valid token:

   ```go
   func (p *Publisher) Prepare(version string) error {
       if os.Getenv("PACKAGIST_TOKEN") == "" {
           return fmt.Errorf("packagist: PACKAGIST_TOKEN is empty, cannot prepare for publish")
       }
       // ... version bump logic ...
   }
   ```

---

## FINDING E02-PACKAGIST-03 — MEDIUM (Confidence 7/10)

### Package name injection in Verify URL

**File:** `internal/publishers/packagist/packagist.go` (Verify)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Query Strings in Request / Injection  
**CWE:** CWE-532 / CWE-88

#### Description

`Verify(version)` constructs a URL from the vendor/package name:
```go
url := fmt.Sprintf("https://packagist.org/packages/%s/%s.json", vendor, pkg)
```
If the vendor or package name (read from `composer.json`) contains special characters (e.g., `../../`, newlines, URL-unsafe characters), the URL could resolve to an unexpected endpoint. While Packagist enforces valid package name rules on their side, a malicious `composer.json` could contain a crafted `name` field.

#### Exploit Scenario

1. A repository has a `composer.json` with `"name": "../../config/special"`.
2. When `Verify` runs, it calls `GET https://packagist.org/packages/../../config/special.json`.
3. Depending on the HTTP client and server, this might resolve unexpectedly or leak information about the request.

#### Recommendation

Validate the vendor and package name before using them in URL construction:

```go
var validPackagistName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

func isValidPackagistName(name string) bool {
    return len(name) >= 1 && len(name) <= 200 && validPackagistName.MatchString(name)
}
```

---

## FINDING E02-PACKAGIST-04 — LOW (Confidence 5/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/packagist/packagist.go` (Publish, Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Publish` and `Verify` methods read HTTP response bodies without a size limit. If Packagist returns an unexpectedly large response (e.g., a 5xx error page with megabytes of HTML), the publisher will allocate memory to hold the full response body. While unlikely in practice, this could be triggered by a misconfigured proxy or CDN returning a large error page.

#### Recommendation

Use `http.MaxBytesReader` or `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 10 * 1024 * 1024 // 10 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("packagist: response body too large")
}
```

---

## e02s06 — Swift Publisher

> **Risk Tier:** P3  
> **Token:** None (token-less — git tag-based versioning only)  
> **API Endpoints:**
> - Publish: `git tag <version> && git push origin <tag>`
> - Verify: `git tag -l <tag>` (local check only, no HTTP)
>
> ### Data Flow
>
> ```
> ┌──────────────────────────────┐
> │  Swift Publisher               │
> │                              │
> │  Publish():                   │
> │    git tag <version>         │
> │    git push origin <version> │
> │                              │
> │  Prepare():                   │
> │    no-op (tag-based)         │
> │                              │
> │  Verify():                    │
> │    git tag -l <version>      │
> └──────┬───────────────────────┘
>        │
>        └── os/exec: git tag, git push, git tag -l
> ```
>
> **Trust boundary:** No credentials are handled. The publisher relies on the user's local Git credentials (`git push`). All operations use `os/exec.Command` with variadic args — no shell invocation.
>
> ---
>
> ## FINDING E02-SWIFT-01 — MEDIUM (Confidence 7/10)
>
> ### Git tag injection via version string
>
> **File:** `internal/publishers/swift/swift.go` (Publish)  
> **Severity:** MEDIUM  
> **Category:** Argument Injection  
> **CWE:** CWE-88
>
> #### Description
>
> `Publish(version)` constructs a git tag from the version string and passes it to `git tag` and `git push` via `os/exec`. If the version string contains characters interpreted as git flags (e.g., `-d`, `--delete`) or newlines, the resulting git command may perform unintended operations.
>
> Because `os/exec.Command` is used with variadic args (no shell), shell metacharacters (`;`, `|`, `` ` ``, `$()`) are not a vector. However, flag injection is possible — a version string like `--delete` passed as `git tag --delete` would delete an existing tag instead of creating one.
>
> #### Exploit Scenario
>
> 1. An attacker-controlled input (e.g., from commit message parsing or changelog generation) produces a version string like `--delete`.
> 2. `git tag --delete` is executed, deleting the `--delete` tag (or if another positional arg follows, deleting that tag).
> 3. If a crafted tag name includes newlines, multiple git commands could be injected (though `os/exec.Command` passes each argument as a single OS-level argv element, limiting this risk).
>
> #### Recommendation
>
> 1. **Use `exec.Command` (variadic args) exclusively** — already using this pattern with `ExecCommand` field.
> 2. **Validate version string** before constructing the tag — reject characters beyond `[a-zA-Z0-9._-]` and reject leading dashes:
>    ```go
>    var validVersion = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
>
>    func isValidVersion(v string) bool {
>        return len(v) > 0 && len(v) <= 128 && validVersion.MatchString(v)
>    }
>    ```
> 3. **Return opaque errors** — don't include the raw version string in error messages that could appear in pipeline logs:
>    ```go
>    return fmt.Errorf("swift: publish failed (git tag)")
>    ```
>
> ---

>>>>>>> Stashed changes
## e02s07 — Godot Publisher

> **Risk Tier:** P3  
> **Token:** `GITHUB_TOKEN` (env var, Bearer token in HTTP header)  
> **API Endpoints:**
> - Publish: `POST https://api.github.com/repos/<owner>/<repo>/releases`
> - Verify: `GET https://api.github.com/repos/<owner>/<repo>/releases/tags/<version>`

### Data Flow

```
GITHUB_TOKEN (env var)
    │
    ▼
┌───────────────────────────────────┐
│  Godot Publisher                    │
│                                   │
│  Publish():                        │
│    POST /repos/<owner>/<repo>/     │
│      releases                      │
│    Authorization: token             │
│    <GITHUB_TOKEN>                  │
│                                   │
│  Verify():                         │
│    GET /repos/<owner>/<repo>/      │
│      releases/tags/<version>      │
└──────┬────────────────────────────┘
       │
       ▼
   GitHub API (api.github.com)
```

**Trust boundary:** The `GITHUB_TOKEN` environment variable is set outside the Go process (CI secrets or shell env). The publisher reads it at publish time. Owner/repo are configuration parameters embedded in the tool config, not user-controlled at runtime.

---

## FINDING E02-GODOT-01 — HIGH (Confidence 9/10)

### Token exposure via environment dump or log leakage

**File:** `internal/publishers/godot/godot.go` (Publish)  
**Severity:** HIGH  
**Category:** Insufficiently Protected Credentials  
**CWE:** CWE-522

#### Description

The `GITHUB_TOKEN` is read from the environment using `os.Getenv("GITHUB_TOKEN")`. If the token value is ever logged (e.g., in a debug log, error message, or panic recovery), it will be exposed. The token is also present in the process's `/proc/self/environ` and in any crash dump.

Additionally, if the HTTP client is configured to log request headers (common during debugging), the `Authorization: token <GITHUB_TOKEN>` header would be written to logs verbatim.

#### Exploit Scenario

1. A developer runs `big-release --verbose release` and a transient network error occurs.
2. The error handler logs the full HTTP request including headers: `Authorization: token ghp_xxxxx`.
3. The log is uploaded to a log aggregation service (e.g., Datadog, CloudWatch) with lesser access controls than the CI secrets manager.
4. An attacker with read access to logs extracts the GitHub token.

#### Recommendation

1. **Never log the token value.** If logging the Authorization header is necessary, replace the token value with `[REDACTED]`:
   ```go
   func redactAuthHeader(req *http.Request) string {
       h := req.Header.Get("Authorization")
       if h != "" {
           req.Header.Set("Authorization", "token [REDACTED]")
       }
       return fmt.Sprintf("%s %s", req.Method, req.URL)
   }
   ```

2. **Return opaque errors** from Publish — don't include `os.Getenv` output or raw `Authorization` header values:
   ```go
   // Good:
   return fmt.Errorf("godot: publish failed (HTTP %d)", resp.StatusCode)

   // Bad:
   return fmt.Errorf("godot: publish failed with token %s: HTTP %d", token, resp.StatusCode)
   ```

3. **Consider using `os.Getenv` with immediate zeroing** — since Go strings are immutable, the token will remain in memory until garbage collected. For high-security deployments, consider using `memguard` or similar for sensitive strings.

---

## FINDING E02-GODOT-02 — MEDIUM (Confidence 8/10)

### Missing token validation before HTTP call

**File:** `internal/publishers/godot/godot.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Sent Data  
**CWE:** CWE-201

#### Description

`Publish(version)` checks for an empty `GITHUB_TOKEN` and returns an error before making HTTP calls. However, the check happens at the `Publish` call boundary only. If `Prepare()` is called first and succeeds (updating version in `project.godot`), and then `Publish()` fails with an empty token error, the config files have already been modified — the system is left in a dirty state.

#### Exploit Scenario

1. `GITHUB_TOKEN` is accidentally unset in CI.
2. `Prepare("2.0.0")` runs successfully, updating `project.godot` version.
3. `Publish("2.0.0")` returns error: "GITHUB_TOKEN is empty".
4. CI pipeline shows "release failed" but the working tree has already been modified.
5. Next `git status` shows dirty tree; re-running may double-commit the version bump.

#### Recommendation

1. **Pre-validate the token at construction time** — call `os.Getenv("GITHUB_TOKEN")` in `NewPublisher()` and store the result. Validate it in `Prepare()` (before mutating files) and in `Publish()`.

2. **Or, validate in `Prepare()`** so that version bump doesn't happen without a valid token:

   ```go
   func (p *Publisher) Prepare(version string) error {
       if os.Getenv("GITHUB_TOKEN") == "" {
           return fmt.Errorf("godot: GITHUB_TOKEN is empty, cannot prepare for publish")
       }
       // ... version bump logic ...
   }
   ```

---

## FINDING E02-GODOT-03 — MEDIUM (Confidence 7/10)

### Owner/repo injection in GitHub API URL

**File:** `internal/publishers/godot/godot.go` (Publish, Verify)  
**Severity:** MEDIUM  
**Category:** Injection  
**CWE:** CWE-88

#### Description

The `Publish` and `Verify` methods construct GitHub API URLs from the owner and repo:
```go
url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
```
If the owner or repo name (read from configuration) contains special characters (e.g., `../../`, newlines, URL-unsafe characters), the URL could resolve to an unexpected endpoint. While GitHub enforces valid repository name rules, a malicious configuration could contain crafted owner/repo values.

#### Exploit Scenario

1. A repository has a configuration with `owner = "../../config/special"`.
2. When `Publish` runs, it calls `POST https://api.github.com/repos/../../config/special/releases`.
3. Depending on the HTTP client and server, this might resolve unexpectedly or leak information about the request.

#### Recommendation

Validate the owner and repo before using them in URL construction:

```go
var validGitHubName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

func isValidGitHubName(name string) bool {
    return len(name) >= 1 && len(name) <= 100 && validGitHubName.MatchString(name)
}
```

---

## FINDING E02-GODOT-04 — LOW (Confidence 5/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/godot/godot.go` (Publish, Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Publish` and `Verify` methods read HTTP response bodies without a size limit. If GitHub API returns an unexpectedly large response (e.g., a 5xx error page with megabytes of HTML), the publisher will allocate memory to hold the full response body. While unlikely in practice, this could be triggered by a misconfigured proxy or CDN returning a large error page.

#### Recommendation

Use `http.MaxBytesReader` or `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 10 * 1024 * 1024 // 10 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("godot: response body too large")
}
```
<<<<<<< Updated upstream
=======

---
---

## e02s05 — Maven Publisher

> **Risk Tier:** P2  
> **Token:** `MAVEN_TOKEN` (env var, Bearer token in HTTP header)  
> **API Endpoints:**
> - Publish: `POST https://central.sonatype.com/api/v1/publisher/upload`
> - Verify: Maven Central search/query API

### Data Flow

```
MAVEN_TOKEN (env var)
    │
    ▼
┌──────────────────────────────────┐
│  Maven Publisher                  │
│                                  │
│  Publish():                       │
│    POST /api/v1/publisher/upload  │
│    Authorization: Bearer          │
│      <MAVEN_TOKEN>               │
│                                  │
│  Verify():                        │
│    GET search.maven.org/         │
│      solrsearch/select           │
└──────┬───────────────────────────┘
       │
       ▼
   Maven Central (Sonatype)
```

**Trust boundary:** The `MAVEN_TOKEN` environment variable is set outside the Go process (CI secrets or shell env). The publisher reads it at publish time.

---

## FINDING E02-MAVEN-01 — HIGH (Confidence 9/10)

### Token exposure via environment dump or log leakage

**File:** `internal/publishers/maven/maven.go` (Publish)  
**Severity:** HIGH  
**Category:** Insufficiently Protected Credentials  
**CWE:** CWE-522

#### Description

The `MAVEN_TOKEN` is read from the environment using `os.Getenv("MAVEN_TOKEN")`. If the token value is ever logged (e.g., in a debug log, error message, or panic recovery), it will be exposed. The token is also present in the process's `/proc/self/environ` and in any crash dump.

Additionally, if the HTTP client is configured to log request headers (common during debugging), the `Authorization: Bearer <MAVEN_TOKEN>` header would be written to logs verbatim.

#### Exploit Scenario

1. A developer runs `big-release --verbose release` and a transient network error occurs.
2. The error handler logs the full HTTP request including headers: `Authorization: Bearer maven-xxxxx`.
3. The log is uploaded to a log aggregation service (e.g., Datadog, CloudWatch) with lesser access controls than the CI secrets manager.
4. An attacker with read access to logs extracts the Maven Central token.

#### Recommendation

1. **Never log the token value.** If logging the Authorization header is necessary, replace the token value with `[REDACTED]`:
   ```go
   func redactAuthHeader(req *http.Request) string {
       h := req.Header.Get("Authorization")
       if h != "" {
           req.Header.Set("Authorization", "Bearer [REDACTED]")
       }
       return fmt.Sprintf("%s %s", req.Method, req.URL)
   }
   ```

2. **Return opaque errors** from Publish — don't include `os.Getenv` output or raw `Authorization` header values:
   ```go
   // Good:
   return fmt.Errorf("maven: publish failed (HTTP %d)", resp.StatusCode)

   // Bad:
   return fmt.Errorf("maven: publish failed with token %s: HTTP %d", token, resp.StatusCode)
   ```

3. **Consider using `os.Getenv` with immediate zeroing** — since Go strings are immutable, the token will remain in memory until garbage collected. For high-security deployments, consider using `memguard` or similar for sensitive strings.

---

## FINDING E02-MAVEN-02 — MEDIUM (Confidence 8/10)

### Missing token validation before HTTP call

**File:** `internal/publishers/maven/maven.go` (Publish)  
**Severity:** MEDIUM  
**Category:** Information Exposure Through Sent Data  
**CWE:** CWE-201

#### Description

`Publish(version)` checks for an empty `MAVEN_TOKEN` and returns an error before making HTTP calls. However, the check happens at the `Publish` call boundary only. If `Prepare()` is called first and succeeds (updating version in pom.xml), and then `Publish()` fails with an empty token error, the config files have already been modified — the system is left in a dirty state.

#### Exploit Scenario

1. `MAVEN_TOKEN` is accidentally unset in CI.
2. `Prepare("2.0.0")` runs successfully, updating `pom.xml` version.
3. `Publish("2.0.0")` returns error: "MAVEN_TOKEN is empty".
4. CI pipeline shows "release failed" but the working tree has already been modified.
5. Next `git status` shows dirty tree; re-running may double-commit the version bump.

#### Recommendation

1. **Pre-validate the token at construction time** — call `os.Getenv("MAVEN_TOKEN")` in `NewPublisher()` and store the result. Validate it in `Prepare()` (before mutating files) and in `Publish()`.

2. **Or, validate in `Prepare()`** so that version bump doesn't happen without a valid token:

   ```go
   func (p *Publisher) Prepare(version string) error {
       if os.Getenv("MAVEN_TOKEN") == "" {
           return fmt.Errorf("maven: MAVEN_TOKEN is empty, cannot prepare for publish")
       }
       // ... version bump logic ...
   }
   ```

---

## FINDING E02-MAVEN-03 — MEDIUM (Confidence 7/10)

### POM group/artifact/version injection in Verify URL

**File:** `internal/publishers/maven/maven.go` (Verify)  
**Severity:** MEDIUM  
**Category:** Injection  
**CWE:** CWE-88

#### Description

`Verify(version)` constructs a URL from the group ID, artifact ID, and version read from `pom.xml`. If these values contain special characters (e.g., `../../`, newlines, URL-unsafe characters), the URL could resolve to an unexpected endpoint.

#### Exploit Scenario

1. A repository has a `pom.xml` with `<groupId>../../config/special</groupId>`.
2. When `Verify` runs, it calls `GET ...?q=g:../../config/special+AND+a:artifact+v:version`.
3. Depending on the HTTP client and server, this might resolve unexpectedly or leak information.

#### Recommendation

Validate the group ID, artifact ID, and version before using them in URL construction:

```go
var validMavenIdentifier = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func isValidMavenIdentifier(id string) bool {
    return len(id) > 0 && len(id) <= 256 && validMavenIdentifier.MatchString(id)
}
```

---

## FINDING E02-MAVEN-04 — LOW (Confidence 5/10)

### Unbounded response body in HTTP client

**File:** `internal/publishers/maven/maven.go` (Publish, Verify)  
**Severity:** LOW  
**Category:** Uncontrolled Resource Consumption  
**CWE:** CWE-400

#### Description

The `Publish` and `Verify` methods read HTTP response bodies without a size limit. If Maven Central returns an unexpectedly large response (e.g., a 5xx error page with megabytes of HTML), the publisher will allocate memory to hold the full response body.

#### Recommendation

Use `http.MaxBytesReader` or `io.LimitReader` when reading response bodies:

```go
const maxResponseSize = 10 * 1024 * 1024 // 10 MB
body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
if len(body) > maxResponseSize {
    return fmt.Errorf("maven: response body too large")
}
>>>>>>> Stashed changes
```
