# story: e18s03
# Title: Clickable Commit Hash and Issue Links

## User Story
As a developer reading the generated release notes on GitHub,
I want commit hashes and issue numbers formatted as clickable markdown links,
So that I can easily navigate directly to the underlying diffs and pull requests.

## Acceptance Criteria
1. Render commit hash as markdown link `[{short}]({repoURL}/commit/{full})` in scoped and unscoped formats.
2. Linkify issue references `#123` or `fixes #123` to `{repoURL}/issues/123`.
3. Provide fallback formatting when `RepositoryURL` is omitted.
