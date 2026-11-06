# story: e18s02
# Title: Hide Non-Release Commit Types From Changelog

## User Story
As a library user reading release notes,
I want non-release commits (docs, chores, ci, tests) omitted from the changelog by default,
So that I only see user-facing changes and bug fixes.

## Acceptance Criteria
1. Default changelog excludes `docs`, `chore`, `ci`, `refactor`, `style`, `test`, `build`.
2. Any commit containing a breaking change note (`BREAKING CHANGE:`) is still surfaced in the breaking changes section even if its type is marked hidden.
