# story: e19s01
# Title: Upload Binary Assets to GitHub Releases

## User Story
As a maintainer publishing compiled binaries,
I want big-release to attach compiled artifacts to the created GitHub release,
So that users can download prebuilt distributions directly from GitHub.

## Acceptance Criteria
1. Introduce `GitHubConfig.Assets` in `algorithm.Config`.
2. Capture release ID from 201 response during release creation.
3. Upload assets via GitHub upload endpoint with correct MIME types and glob expansion.
