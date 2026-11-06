# Test Plan: e18 Notes & Changelog Engine

## 1. Overview & Risk Assessment
Epic `e18` makes the changelog engine configurable, introduces custom section mapping, clickable hash/issue links, and revert filtering.
Overall Risk Level: **Medium** (Core algorithm output touched, string formatting and regex matching across commit bodies).

## 2. Risk-Scaled Test Levels

### Level 1: Unit Tests (Low Risk / Component Isolation)
- **Target:** `internal/algorithm/types.go`, `internal/config/config.go`
- **Scope:**
  - `TestCommitTypeConfig_Defaults`: Validate default seeded map matches semantic-release convention.
  - `TestCommitTypeConfig_CustomSections`: Validate unmarshaling of user-defined type arrays and section titles.
- **Coverage Goal:** 100% of new struct serialization and configuration logic.
- **Verify Command:** `go test ./internal/algorithm -run TestCommitType`

### Level 2: Engine Integration Tests (Medium Risk / String Parsing)
- **Target:** `internal/algorithm/generator.go`
- **Scope:**
  - `TestGenerator_DynamicSectionOrdering`: Render changelogs with altered section priorities.
  - `TestGenerator_HideNonReleaseTypes`: Ensure `chore`, `docs`, `ci` commits do not generate sections when `Hidden: true`.
  - `TestGenerator_HiddenWithBreakingChange`: Ensure breaking changes from hidden types still surface in `BREAKING CHANGES`.
  - `TestGenerator_ClickableCommitLinks`: Validate markdown link generation against mock repository URLs.
  - `TestGenerator_RevertCommitPairing`: Validate matching reverts remove both revert and target commits.
- **Coverage Goal:** >95% statement coverage on `generator.go`.
- **Verify Command:** `go test ./internal/algorithm -run TestGenerator_`

### Level 3: Release Orchestration & Contract Tests (Medium-High Risk / Wiring)
- **Target:** `pkg/release/release.go`
- **Scope:**
  - Verify `commitTypes` from loaded configuration are properly passed into the Generator instance.
  - Verify dry-run and actual release cycles produce identical notes.
- **Verify Command:** `go test ./pkg/release -run TestRun_NotesGeneration`

### Level 4: Full Suite & Regression Verification (Zero-Tolerance)
- **Verify Command:** `make preflight`
