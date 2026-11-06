# Test Plan: e20 Multi-Branch Release Support

## 1. Overview & Risk Assessment
Epic `e20` enables releases from multiple branch lines (stable, maintenance, prerelease channels) with per-branch tag range scans and channel propagation to GitHub Releases.
Risk Level: **High** (Changes git tag matching and version resolution across non-trunk branches).

## 2. Risk-Scaled Test Levels

### Level 1: Unit Tests (Low Risk)
- **Target:** `internal/git/branch.go`, `internal/git/branch_glob_test.go`
- **Scope:** Verify `MatchBranchPattern` handles exact strings (`main`), glob wildcards (`release/*`), and `N.x` release patterns (`1.x`, `1.0.x`).
- **Verify Command:** `go test ./internal/git -run TestClient_BranchGlobbing`

### Level 2: Tag Scanning & Branch Isolation Tests (Medium Risk)
- **Target:** `internal/git/client.go`
- **Scope:** Ensure `GetLatestTag` filters tags belonging strictly to the active branch line rather than global commit histories.
- **Verify Command:** `go test ./internal/git -run TestClient_GetLatestTag`

### Level 3: Release Orchestration Integration Tests (High Risk)
- **Target:** `pkg/release/release.go`
- **Scope:** Verify `BranchConfig.Channel` propagates to `state.NextRelease.Channel` and GitHub release flags (`prerelease: true`).
- **Verify Command:** `go test ./pkg/release -run TestRelease_Channel`

### Level 4: Full Verification Stack
- **Verify Command:** `make preflight`
