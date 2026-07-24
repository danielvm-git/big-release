# Impact Report — e24s03 Extract shared package.json utilities

> **Mode:** lightweight (required — npm dependents)  
> **Date:** 2026-07-24  
> **Story:** e24s03

## Change summary

Extract `readPackageJSON`, `writePackageJSON`, `readPackageName`, `isValidPackageName` into `internal/publishers/nodeutil/`. Rewire npm + pnpm to import nodeutil. npm Detect exclusion already landed in e24s01.

## Dependents

| Area | Impact |
|------|--------|
| `internal/publishers/npm` | Refactor imports; behavior must stay identical |
| `internal/publishers/pnpm` | Same |
| npm + pnpm unit tests | Must stay green without assertion rewrites where possible |

## Blast radius

- Shared helpers touch Prepare/Verify for both Node publishers
- No CLI/registry protocol changes
- Error message prefix preserved via `prefix` argument

## Risk score: **5 / 10**

Medium refactor of shared I/O with strong existing tests. Below grill-me threshold (7).

**Gate:** proceed.
