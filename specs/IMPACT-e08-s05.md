# Impact Report — e08s05 AggregateError

> **Mode:** lightweight | **Date:** 2026-07-24 | **Story:** e08s05

## Change summary

Collect all VerifyConditions errors via `AggregateError` + `errors.Join` in `release.go`.

## Dependents

| Area | Impact |
|------|--------|
| Plugin verify phase | Behavior change: all failures reported |

## Risk score: **4 / 10**

Semantic-release parity; scoped to VerifyConditions only. **Gate:** proceed.
