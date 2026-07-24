---
bug_id: BUG-release-workflow-softprops-and-verbose
status: open
severity: high
scope: ci
title: "Release workflow: softprops upload fails on missing tag, and --verbose emits no output"
---

## Summary

Two compounding defects in the `Release` GitHub Actions workflow
(`.github/workflows/release.yml`) cause the `release` job to fail on every
push to `main`, and leave `big-release release --verbose` effectively silent
in CI.

- **Bug A** — the `Upload release artifacts` step
  (`softprops/action-gh-release@v3`) fails with
  `⚠️ GitHub Releases requires a tag` because the workflow runs on a
  branch push (`refs/heads/main`), not a tag push, and no `tag_name` input
  is supplied.
- **Bug B** — `big-release release --verbose` produces zero useful output:
  the logger writes JSON to stderr and `--verbose` lowers the level to
  Debug, but there are **no `.Debug()` calls** anywhere in the release path
  and success-path `.Info()` calls are sparse.

## Evidence (CI run 30116517165, job 89558642923)

The `Run big-release` step ran for ~7s with **no stdout/stderr**:

```
18:22:12 ##[group]Run big-release release --verbose
18:22:19 ##[group]Run softprops/action-gh-release@v3
```

The tag `2.1.0` and release `v2.1.0` were created (big-release's plugins
worked), but the upload step immediately aborted:

```
Upload release artifacts  ##[error]⚠️ GitHub Releases requires a tag
```

## Root Cause Analysis (4-phase)

### 1. Reproduce

Reproduced directly from the run log (`gh run view --job 89558642923 --log`).
The job shows: `release.yml:150` runs `big-release release --verbose` on
`refs/heads/main`; `release.yml:162-171` then invokes
`softprops/action-gh-release@v3` with only a `files:` input — no `tag_name`.

### 2. Isolate

- **Tag is created correctly** by the `git` plugin:
  `internal/plugins/git.go:181-210` (`createTag` → `pushRefs`) calls
  `internal/git/client.go:110-117` (`git tag -a`) and
  `internal/git/client.go:355-361` (`git push origin --tags`). So the tag
  is a genuine git ref pushed mid-run.
- **The action cannot see it.** `softprops/action-gh-release` derives
  `tag_name` from `github.ref_name` when no input is given. GitHub Actions
  fixes `GITHUB_REF` at workflow start from the **triggering** event
  (`refs/heads/main`); a tag pushed during the run never updates it. The
  action sees `main`, which is not a tag, and refuses.
- **The github plugin already does release + assets via API.**
  `internal/plugins/github.go:274-302` (`Publish`) POSTs to
  `/repos/{repo}/releases` with `tag_name`, then `uploadAssets`
  (`internal/plugins/github_assets.go`) uploads configured assets. This
  mirrors `@semantic-release/github` exactly — making the softprops step
  redundant. The only reason binaries were not uploaded by the plugin is
  that `.big-release.yml` configures **no `pluginConfigs.github.assets`**.
- **`--verbose` is a no-op.** `cmd/big-release/main.go:93-97` builds the
  logger from `zap.NewProductionConfig()` (JSON → stderr); `--verbose` only
  flips the level to Debug. A grep of `internal/` and `pkg/` for `.Debug(`
  returns **zero** matches in the release path. Success-path `.Info()`
  calls (`pkg/release/release.go:156,275,336,379,401,432`) cover only
  PR-skip, no-changes, dry-run, disabled-publisher, and failure cases — a
  successful release logs almost nothing.

### 3. Hypothesize

- **A:** The fix is to let big-release's own github plugin own release
  creation **and** asset upload atomically (the semantic-release pattern),
  by configuring assets and removing the softprops step entirely. Passing
  `tag_name` to softprops would keep two release mechanisms in play and is
  more fragile.
- **B:** (1) When `--verbose`, build the logger from
  `zap.NewDevelopmentConfig()` (human-readable console encoder at Debug
  level) so verbose output is legible. (2) Add the missing success-path
  `.Info()` logs in `release.go` (version computed, plugin published) so a
  successful release narrates itself regardless of `--verbose`.

### 4. Verify

- **semantic-release (opensrc)** creates the GitHub release and uploads
  assets entirely via the REST API inside `@semantic-release/github`'s
  `publish.js` — there is **no** `softprops/action-gh-release` step in its
  workflow. Confirmed via the opensrc cache. big-release's github plugin
  already mirrors this; only the config (assets list) and the redundant
  workflow step differ.
- `softprops/action-gh-release` README documents that `tag_name` defaults
  to `github.ref_name` and recommends guarding with
  `if: github.ref_type == 'tag'` — i.e. it is designed for tag-push runs,
  not branch-push runs. This confirms removing it (rather than feeding it a
  tag) is the cleaner fix.

## Fix Approach

1. **`.big-release.yml`** — add `pluginConfigs.github.assets` listing the
   four cross-compiled binaries so the github plugin uploads them.
2. **`.github/workflows/release.yml`** — delete the `Upload release
   artifacts` step (lines 162-171) and its `env:` block.
3. **`cmd/big-release/main.go`** — extract `buildLogger(verbose bool)`;
   use `zap.NewDevelopmentConfig()` when verbose, else production JSON at
   Info. Keep `zap.WrapCore(secure.WrapCore)`.
4. **`pkg/release/release.go`** — add success-path Info logs (computed
   version, plugin published) and `writeStepOutput` (writes
   `version=`/`published=true` to `$GITHUB_OUTPUT` when set) so downstream
   CI steps can observe the release.

### Files to modify

- `.big-release.yml`
- `.github/workflows/release.yml`
- `cmd/big-release/main.go` (+ `main_test.go`)
- `pkg/release/release.go` (+ `release_test.go`)

## Verify

1. `make preflight` is green (test + lint + build).
2. New tests pass:
   - `TestBuildLogger_VerboseUsesDebugLevel` — verbose emits Debug;
     non-verbose does not.
   - `TestRelease_WritesGitHubOutput` — `$GITHUB_OUTPUT` gains
     `version=`/`published=true` after a successful run.
   - `TestRelease_LogsComputedVersion` — an Info log carrying the version
     is recorded.
3. PR opens; `gh pr checks` is green.
4. After merge to `main`, the next release run creates the tag, the GitHub
   release, and uploads all four binaries via the github plugin — with no
   softprops step and no "requires a tag" error.
