# ADR 0005: Plan Stress-Test & Grilling Decisions

## Status
Accepted

## Context
Before committing to implementation phases, the release architecture and planning spine underwent a structured assumption-surfacing stress test (`grill-me` protocol) against semantic-release precedents, POSIX signal handling, and multi-forge edge cases.

## Key Grilling Challenges & Resolved Decisions

### 1. Partial Publication & Network Severance
- **Challenge:** If npm publishes successfully but GitHub Releases fails due to network partition or token expiration, the repository is left in an inconsistent state (package released, but no git tag or changelog).
- **Resolution:** Introduced two-phase state machine with explicit `Rollbacker` interface (`pkg/release/`). If a post-publish or forge step fails, compensation hooks rollback newly pushed tags, record `ReleaseStatus: failed`, and emit machine-readable reason to `$GITHUB_OUTPUT`.

### 2. Multi-Branch Prerelease Isolation
- **Challenge:** Semantic release multi-branch specifications can cause version collisions if alpha and beta branches calculate from the same base tag.
- **Resolution:** Explicit `BranchConfig` unmarshaling with strict type validation (`release`, `prerelease`, `maintenance`). Disallow conflicting prerelease tags across branches with identical branch names.

### 3. In-Memory Plugin State Leakage
- **Challenge:** Repeated runs in a long-lived process or testing environment could cross-contaminate plugin registry hooks.
- **Resolution:** `RunResetter` snapshots registry per execution run. State mutations exist strictly within `MutableState` and cannot leak to concurrent or sequential runs.

### 4. Secret Exposure via Third-Party Errors
- **Challenge:** Third-party CLI tools (e.g. `npm publish`, `git push`) often echo URL query params or auth headers on HTTP 401/403 errors into stderr.
- **Resolution:** Mandatory `internal/secure` redaction layer streams through a custom `zapcore.Core` wrapping all subprocess output before logging or displaying.
