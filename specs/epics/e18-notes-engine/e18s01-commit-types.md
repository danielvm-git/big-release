# story: e18s01
# Title: Configurable Commit Type Sections & Visibility

## User Story
As a repository maintainer,
I want to configure custom commit types, custom section titles, and hidden status in the changelog,
So that release notes highlight relevant changes cleanly and match project documentation standards.

## Acceptance Criteria
1. `CommitTypeConfig` struct (`Type`, `Section`, `Hidden`) introduced in `internal/algorithm/types.go`.
2. `DefaultConfig` in `internal/config/config.go` seeds standard types (`feat`, `fix`, `perf`, `revert` visible; `docs`, `chore`, `ci`, `refactor`, `style`, `test`, `build` hidden).
3. `Generator` orders sections dynamically according to configured commit types.
4. Breaking changes always render regardless of `Hidden` setting.
