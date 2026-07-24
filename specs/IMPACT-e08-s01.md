# Impact Report — e08s01 Secret Masking

> **Mode:** lightweight | **Date:** 2026-07-24 | **Story:** e08s01

## Change summary

New `internal/secure` package; wire generator + zap logger to redact tokens.

## Dependents

| Area | Impact |
|------|--------|
| `internal/algorithm/generator.go` | Delegates hideSensitive to secure |
| `cmd/big-release/main.go` | Logger uses redacting zap core |

## Risk score: **4 / 10**

New package, narrow surface. **Gate:** proceed.
