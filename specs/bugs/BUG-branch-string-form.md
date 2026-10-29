---
bug_id: BUG-branch-string-form
status: fixed
severity: medium
scope: internal/algorithm, internal/config
title: Documented `branches: - main` string form fails to parse
---

## Problem

`README.md` (§ Configuration) and `docs/reference/configuration.md` show branches as bare
names (`- main`, `- "N.x"`). Every such config fails before any release logic runs.

## Reproduce

1. `.big-release.yml`:
   ```yaml
   branches:
     - main
   tagFormat: "v${version}"
   ```
2. `big-release validate`
3. Expected: `✅ Configuration is valid`
4. Actual: `cannot unmarshal !!str `main` into algorithm.BranchConfig` (same for JSON `["main"]`)

Found by a consumer repo (danielvm-git/big-remote-control) during a release dry-run on 2.5.0; reproduced on `main`.

## Root Cause

`BranchConfig` is a struct with no custom decoder, so YAML/JSON accept only the object form.

## Fix Approach

`BranchConfig.UnmarshalYAML` / `UnmarshalJSON` (`internal/algorithm/branch_config.go`):
a scalar becomes `BranchConfig{Name: value}`; an object decodes as before; anything else
fails with a hint. Docs need no change — they now describe working behavior.

## Resolution

**Fixed:** 2026-10-29
**Evidence:** `go test ./internal/config/ -run 'BranchStringForm|BranchRejects'`

## Not fixed here

`"N.x"` given as a bare name defaults to `type: release`, so the README example
(`main`, `next`, `"N.x"`, `beta`) still fails validation with "maximum 3 release branches".
Inferring `maintenance` from the `N.x` pattern is a separate behavior change.
