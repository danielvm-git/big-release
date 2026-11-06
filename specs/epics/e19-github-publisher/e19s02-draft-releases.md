# story: e19s02
# Title: Support Draft GitHub Releases

## User Story
As a maintainer releasing large assets,
I want releases created as drafts until all assets are uploaded,
So that users do not observe incomplete releases with missing downloads.

## Acceptance Criteria
1. Support `draftRelease: true` in configuration.
2. Publish release as draft initially, attach assets, and finalize with `draft: false`.
