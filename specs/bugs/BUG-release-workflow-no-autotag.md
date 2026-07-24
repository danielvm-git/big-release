---
bug_id: BUG-release-workflow-no-autotag
status: fixed
severity: high
scope: ci
title: Release workflow skips release because no tag exists (chicken-and-egg)
---

## Summary

The `.github/workflows/release.yml` workflow has a "Check tag exists" gate that
prevents `big-release` from running unless a tag already exists on the commit.
Since `big-release` is the tool that should create tags, no releases are ever
created — a chicken-and-egg problem.

## Root Cause

The workflow checks for an existing tag before running `big-release`:

```yaml
- name: Check tag exists
  id: check-tag
  run: |
    if git describe --exact-match --tags HEAD >/dev/null 2>&1; then
      echo "tag_exists=true"
    else
      echo "tag_exists=false"
      echo "No tag on this commit, skipping release."
    fi
```

All subsequent steps are conditional on `tag_exists == 'true'`.

## Evidence

- Only 1 release exists: `2.0.0` (created manually or by earlier workflow)
- 15+ commits since that release (features, fixes, deps) — none triggered a release
- CI logs confirm: `No tag on this commit, skipping release.`

## How semantic-release solved this

semantic-release's workflow has **no tag check** — it runs directly on push to main:

```yaml
on:
  push:
    branches: [master, next, beta, alpha, "*.x"]

steps:
  - uses: actions/checkout@v6
  - run: npx semantic-release
```

Internally, semantic-release:
1. Finds existing tags via `git tag --merged branch`
2. Gets commits since last tag
3. Analyzes Conventional Commits to determine version bump
4. Creates the tag itself: `git tag v2.1.0 HEAD`
5. Pushes tag: `git push origin --follow-tags`

## Fix Approach

Remove the tag check gate and let `big-release` handle version detection and
tag creation — the same pattern semantic-release uses.

### Files to modify

- `.github/workflows/release.yml` — remove "Check tag exists" step and all
  `if: steps.check-tag.outputs.tag_exists == 'true'` conditions

## Verify

1. Push a `feat:` commit to main
2. Release workflow should run `big-release` without tag gate
3. `big-release` should detect last tag, analyze commits, create new tag
4. GitHub release should be created with the new version
