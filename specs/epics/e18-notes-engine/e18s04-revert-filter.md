# story: e18s04
# Title: Filter Revert Commits From Release Notes

## User Story
As a repository maintainer,
I want revert commits and their reverted target commits excluded from the changelog,
So that clean release notes do not describe temporary features that were retracted before the release.

## Acceptance Criteria
1. Revert parser detects `revert: ...` or `"This reverts commit <hash>"`.
2. Both the revert commit and the target commit it reverted are omitted from notes.
3. Orphaned reverts (where target commit is outside current release range) still display normally.
