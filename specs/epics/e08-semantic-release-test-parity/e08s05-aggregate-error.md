# e08s05 — Aggregate Pipeline Error Handling Tests

**BCP:** 2 | **Status:** todo

## Acceptance (Gherkin)

```gherkin
# SC-e08s05-P0-01
Scenario: Two verify failures aggregated
  Given two plugins fail VerifyConditions
  When runPluginLifecycle runs
  Then error contains both failure messages

# SC-e08s05-P1-01
Scenario: Unwrap all errors
  Given AggregateError returned
  When errors.Unwrap is used
  Then all constituent errors accessible
```

## Files

- `pkg/release/errors.go`, `release.go`, `aggregate_error_test.go`
