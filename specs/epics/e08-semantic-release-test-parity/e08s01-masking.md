# e08s01 — Secret Masking & Redaction Tests

**BCP:** 2 | **Status:** todo

## Acceptance (Gherkin)

```gherkin
# SC-e08s01-P0-01
Scenario: Known env token redacted
  Given GH_TOKEN is set in environment
  When RedactKnownSecrets processes text containing the token value
  Then output must not contain the raw token

# SC-e08s01-P0-02
Scenario: Pattern-based redaction
  Given text contains "token=ghp_secret123"
  When Redact is applied
  Then output contains "[secure]" not the secret

# SC-e08s01-P1-01
Scenario: Release notes sanitized
  Given commits mention token= values
  When GenerateNotes runs
  Then notes redact sensitive patterns

# SC-e08s01-P1-02
Scenario: Zap logs redacted
  Given logger uses secure zap core
  When logging a message with token value
  Then emitted log redacts the value
```

## Files

- `internal/secure/redact.go`, `redact_test.go`, `zap.go`
- `internal/algorithm/generator.go`
- `cmd/big-release/main.go`
