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
| HIGH     | 1     | CWE-522     |
| MEDIUM   | 3     | CWE-201, CWE-532, CWE-88 |
| LOW      | 3     | CWE-400     |

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
|| Unbounded response body (E02-GOPROXY-02) | LOW | 6/10 | CWE-400 | Trivial (add LimitReader) |
|| Module name injection in Verify URL (E02-GOPROXY-03) | LOW | 5/10 | CWE-88 | Small (validate module path) |

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
