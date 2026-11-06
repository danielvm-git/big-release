# story: e20s02
# Title: Channel Propagation to GitHub Release

## User Story
As a project maintainer releasing prereleases or distribution channels (e.g. beta, next),
I want the branch distribution channel propagated to the GitHub release prerelease flag and tag metadata,
So that prereleases are not marked as the latest production release.

## Acceptance Criteria
1. Propagate `BranchConfig.Channel` to `NextRelease.Channel`.
2. Configure GitHub release payload with `prerelease: true` when a channel is set.
