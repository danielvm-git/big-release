# e08s02 — Git Authentication & URL Mutation Tests

**BCP:** 2 | **Status:** todo

## Acceptance (Gherkin)

```gherkin
# SC-e08s02-P0-01
Scenario: HTTPS token injection
  Given remote URL https://github.com/org/repo.git
  When AuthURL is called with a token
  Then result embeds token as HTTPS userinfo

# SC-e08s02-P0-02
Scenario: SSH passthrough
  Given remote URL git@github.com:org/repo.git
  When AuthURL is called with a token
  Then result is unchanged

# SC-e08s02-P0-03
Scenario: No double inject
  Given URL already contains credentials
  When AuthURL is called again
  Then token is not injected twice

# SC-e08s02-P1-01
Scenario: Redacted errors
  Given AuthURL fails
  Then error message does not contain raw token
```

## Files

- `internal/git/auth_url.go`, `auth_url_test.go`
