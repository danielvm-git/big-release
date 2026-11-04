---
bug_id: BUG-versionfile-top-level-config
status: fixed
severity: medium
scope: pkg/release, docs
title: Top-level versionFile config is ignored and docs recommend the wrong plugin order
refs:
  - https://github.com/danielvm-git/big-release/issues/74
  - https://github.com/danielvm-git/big-release/issues/62
---

## Summary

`algorithm.Config.VersionFile` is defined and documented, but nothing reads
it. Only `pluginConfigs.versionfile` reaches `VersionFilePlugin.Configure`.
The how-to also puts `versionfile` after `git` in the `plugins` list, so the
git plugin commits `postTagAssets` before the version file is written.

## Root Cause

1. **Config not wired.** `pkg/release/release.go` Phase 0 reads
   `ctx.Config.PluginConfigs[name]` only. `Config.VersionFile` (declared in
   `internal/algorithm/types.go`) is never converted into the raw map that
   `VersionFilePlugin.Configure` accepts.
2. **Plugin not auto-added.** The effective plugin list is exactly
   `ctx.Config.Plugins`. A user who sets only top-level `versionFile` never
   gets the plugin in the lifecycle.
3. **Docs order.** `docs/how-to/version-file-sync.md` recommends
   `plugins: [changelog, git, github, versionfile]`. Phase 6.5 iterates in
   declaration order, so `git.PostPublish` (stage/commit) runs before
   `versionfile.PostPublish` (write). The commit is empty and the file is
   left uncommitted. The E2E test already uses the working order:
   `[]string{"versionfile", "git"}`.

## Fix Approach

1. Build the effective plugin list: if `Config.VersionFile != nil` and
   `"versionfile"` is absent, append it.
2. For `versionfile`, merge the top-level config into the plugin config
   (`path`, `template`), with explicit `pluginConfigs.versionfile` winning.
3. Update `docs/how-to/version-file-sync.md`: put `versionfile` before `git`
   in every example and document the ordering requirement; keep top-level
   `versionFile` as the supported primary form.
4. Tests:
   - top-level `versionFile` auto-adds the plugin and writes the configured
     path/template;
   - `pluginConfigs.versionfile` overrides the top-level values;
   - PostPublisher plugins run in declaration order (writer before committer).

## Verification

```bash
go test -race ./pkg/release/ ./internal/plugins/ -run 'VersionFile|PostPublish'
go test -race ./...
go vet ./...
golangci-lint run ./...
```

Expected: top-level `versionFile` writes the configured path/template; the
ordering test shows the writer runs before the committer.

## Resolution

- `Release.effectivePlugins` appends `versionfile` when `Config.VersionFile`
  is set and the plugin is not already listed.
- Phase 0 merges the top-level block into the raw `versionfile` config;
  explicit `pluginConfigs.versionfile` values override it.
- `docs/how-to/version-file-sync.md` puts `versionfile` before `git` in every
  example and documents the ordering requirement.
- Regression tests: `TestRunPluginLifecycle_TopLevelVersionFileConfig`,
  `TestRunPluginLifecycle_PluginConfigOverridesTopLevelVersionFile`,
  `TestPostPublish_RunsInPluginDeclarationOrder`.
- Gate: 764 tests (`-race`), vet, lint, build green.
