---
bug_id: BUG-tag-ignores-tagformat
status: fixed
severity: high
scope: internal/plugins
title: git plugin's CreateTag ignores tagFormat, breaking every second release with a non-default prefix
---

## Problem

`GetLastRelease(tagFormat)` (`internal/git/client.go:235`) correctly applies `tagFormat` when
parsing existing tags to find the last release. But `internal/plugins/git.go:182`'s `Publish`
step creates the new tag with the bare semver string instead:

```go
return p.Git.CreateTag(version, fmt.Sprintf("release %s", version))
```

`tagFormat` (config key, default `"v${version}"`) is never applied when *writing* the tag —
only when *reading* it back. The first release on any repo appears to work (it creates a tag,
e.g. `0.1.0` instead of the configured `v0.1.0`, and nothing checks it). The second release is
where it breaks: `GetLastRelease("v${version}")` looks for a tag matching `v*`, doesn't find
`0.1.0`, concludes there is no prior release, recomputes the *same* next version, and
`CreateTag` fails with `exit status 128` (tag already exists).

This affects every repo using a non-empty tag prefix — which is the *default* config
(`DefaultTagFormat = "v${version}"`, `internal/config/config.go:20`). big-release's own
`.big-release.yml` uses `tagFormat: "${version}"` (no prefix) — see its comment "Tags in
this repo use format '2.0.0' (no v prefix)" — which happens to make `GetLastRelease` and
`CreateTag` agree by coincidence, masking the bug in this repo's own dogfooding.

## Reproduce

Discovered live in `danielvm-git/bigbase-canary-go` (a canary repo built specifically to
exercise big-release + bigbase-deploy):

1. `.big-release.yml` with `tagFormat: "v${version}"` (the default), `initialVersion: "0.1.0"`.
2. `big-release release` on the first `feat:`/`chore:` push → succeeds, computes `0.1.0`,
   creates git tag `0.1.0` (not `v0.1.0`), pushes changelog commit + tag.
3. A second push with a `fix:` commit → `big-release release` again:
   ```
   INFO  Computed next release  {"version": "0.1.0", "type": "minor"}
   ERROR Release failed, running fail hooks {"error": "plugin \"git\" publish failed: failed to create tag: exit status 128"}
   ```
   Expected: `0.1.1` (or `v0.1.1`). Actual: recomputes `0.1.0` again and fails to tag it
   because `0.1.0` already exists.

## Root Cause

`internal/plugins/git.go`'s `Publish` (~line 182) calls `p.Git.CreateTag(version, ...)` with
the raw semver, never running it through the same template substitution
`GetLastRelease(tagFormat)` uses to parse tags back. The write path and read path disagree
on tag naming whenever `tagFormat` isn't the empty/bare format.

## Fix Approach

Apply `tagFormat` to `version` before calling `CreateTag`, using the same substitution
`GetLastRelease` already implements (`${version}` → semver string) — one shared helper so the
two paths can't drift again. Add a regression test that runs `release` twice in sequence with
a prefixed `tagFormat` (e.g. `"v${version}"`) and asserts the second run computes a bumped
version and both tags (`v1.0.0`, `v1.0.1`) exist.

## Files

- `internal/plugins/git.go`
- `internal/git/client.go` (tag-format substitution helper, if not already shared)
