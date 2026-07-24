# e08s04 — Git Commit Traversal & Filtering Tests

**BCP:** 1 | **Status:** todo

## Acceptance (Gherkin)

```gherkin
# SC-e08s04-P0-01
Scenario: Last semver tag selected
  Given tags v1.0.0 and v1.1.0
  When GetLastRelease runs
  Then v1.1.0 is selected

# SC-e08s04-P0-02
Scenario: Commit range from tag to HEAD
  Given tag v1.0.0 and commits after it
  When GetCommits(from, HEAD) runs
  Then only post-tag commits returned

# SC-e08s04-P1-01
Scenario: Empty range at tag
  Given HEAD equals last tag
  When GetCommits runs
  Then empty commit list
```

## Files

- `internal/git/commits_test.go`, `client.go`, `testrepo/`
