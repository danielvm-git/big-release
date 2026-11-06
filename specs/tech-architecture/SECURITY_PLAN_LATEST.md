# Security Plan & Threat Model

## 1. Security Architecture & Threat Model

### Assets to Protect
1. **Forge and Registry Authentication Tokens:** `GITHUB_TOKEN`, `NPM_TOKEN`, `PYPI_TOKEN`, `CARGO_TOKEN`, `GITLAB_TOKEN`.
2. **Repository Integrity:** Git tag creation, release asset binaries, and commit histories.
3. **Distribution Supply Chain:** The `big-release` binary distribution artifacts and their cryptographic checksums.

### Trust Boundaries
```
[External CI Environment / Secrets] 
        │ (Environment Variables)
        ▼
┌───────────────────────────────────────────────┐
│              big-release Runtime              │
│  ┌─────────────────────────────────────────┐  │
│  │   Secure In-Memory Boundary             │  │
│  │   - Token Redaction Engine              │  │
│  │   - Authenticated Git Remote Transport   │  │
│  └─────────────────────────────────────────┘  │
└───────────────────────────────────────────────┘
        │ (Subprocess & HTTPS)
        ▼
[Public Package Registries & Forge APIs]
```

### Threat Modeling Matrix (STRIDE)

| Threat Type | Target Vector | Risk | Mitigation Strategy |
| :--- | :--- | :--- | :--- |
| **Spoofing** | Git author / committer identity spoofing | High | Enforce mandatory bot committer config (`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`); require verified CI forge context. |
| **Tampering** | Asset tampering during build or publication | High | Binaries built in ephemeral runners; checksum generation (`sha256sum`); atomic git tag push. |
| **Repudiation** | Untraceable release actions | Medium | Structured JSON logging with Zap; audit records in git history with conventional commit provenance. |
| **Information Disclosure** | Token leakage in stdout, stderr, or CI logs | Critical | Dynamic redaction core (`internal/secure/redact.go`); scrub URLs before process invocation; block `.env*` in `.gitignore`. |
| **Denial of Service** | Process hang or infinite loop during network I/O | Medium | Mandatory `context.WithTimeout` on all HTTP and subprocess executions; rate-limit backoff. |
| **Elevation of Privilege** | Command injection via commit message or branch name | High | Strictly parameterized CLI exec (`exec.CommandContext`), zero shell evaluation (`sh -c` prohibited for external inputs). |

## 2. Token Masking and Secure Execution Invariants

1. **Subprocess Isolation:** Never pass secrets as command-line arguments (which appear in `ps` and process tables). Pass credentials via environment variables or stdin.
2. **Log Redaction:** All output channels (stdout, stderr, log files) route through `secure.ZapCore` to redact tokens matching known environment keys and high-entropy patterns.
3. **No Persistent Token Storage:** No credentials may ever be saved to `.git/config`, disk caches, or local repository files.
4. **Git Remote Cleaning:** Dynamic remote URLs containing credentials (`https://token@github.com/...`) are scrubbed immediately following fetch or push operations.
