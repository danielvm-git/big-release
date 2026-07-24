# e08s06 — Configuration Loading (File + CLI) Tests

**BCP:** 1 | **Status:** todo

## Acceptance (Gherkin)

```gherkin
# SC-e08s06-P0-01
Scenario: YAML config loads
  Given .big-release.yml with custom tagFormat
  When Load runs
  Then tagFormat is merged over defaults

# SC-e08s06-P0-02
Scenario: JSON config loads
  Given .big-release.json
  When Load runs
  Then values parsed correctly

# SC-e08s06-P1-01
Scenario: Parent discovery
  Given config in parent directory
  When Load with empty path from child dir
  Then parent config found

# SC-e08s06-P1-02
Scenario: Explicit CLI path
  Given --config flag path
  When Load uses explicit file
  Then that file wins over discovery
```

## Files

- `internal/config/config.go`, `config_test.go`
