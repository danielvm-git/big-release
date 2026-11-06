# story: e20s01
# Title: Per-Branch Tag Range Scan and Glob Branch Matching

## User Story
As a maintainer maintaining multiple release streams (e.g. 1.x maintenance and 2.x stable),
I want big-release to match branch globs and filter tags per branch,
So that maintenance branches calculate version bumps relative to their own version line.

## Acceptance Criteria
1. Support branch glob patterns (e.g. `1.x`, `*.x`) in `validateBranch`.
2. Filter git tags to scan only tags matching the active branch line.
