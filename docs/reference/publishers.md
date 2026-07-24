# Publishers Reference

> Reference for built-in language publishers in big-release.

Publishers implement Detect → Prepare → Publish → Verify. Registration is via package `init()`. Shared Node helpers live in `internal/publishers/nodeutil`.

## Detection precedence (npm vs pnpm)

| Markers present | Active publisher |
|-----------------|------------------|
| `pnpm-lock.yaml` and/or `pnpm-workspace.yaml` | **pnpm** (npm Detect returns false) |
| `package.json` only (no pnpm markers) | **npm** |
| Neither | neither Node publisher |

This prevents double-publishing when a pnpm project also has `package.json`.

## npm

| Field | Value |
|-------|-------|
| Name | `npm` |
| Detect | `package.json` present **and** no pnpm markers |
| Prepare | Writes `version` into `package.json` |
| Publish | `npm publish` [+ `--tag <channel>`] |
| Verify | `npm view -- <name> version` |
| Auth | `NPM_TOKEN` / `.npmrc` (CLI-managed) |
| DryRun | Skips publish and view |

Channel `latest` (or empty) omits `--tag`. Other channels pass `--tag <channel>`.

## pnpm

| Field | Value |
|-------|-------|
| Name | `pnpm` |
| Detect | `pnpm-lock.yaml` **OR** `pnpm-workspace.yaml` |
| Prepare | Writes `version` into `package.json` (via nodeutil) |
| Publish | `pnpm publish --no-git-checks` [+ `--tag <channel>`] |
| Verify | `pnpm view <name> version` |
| Auth | Same npm-compatible registry auth (`.npmrc` / env) |
| DryRun | Skips publish and view |

`--no-git-checks` is always passed so CI releases are not blocked by dirty git state after Prepare mutates `package.json`.

### Example config

```yaml
publishers:
  pnpm:
    enabled: true
```

Workspace monorepo orchestration (publishing many packages from one root) is out of scope; use per-package CWD.

## Other publishers

See [Publisher Guide](../how-to/publishers/README.md) for PyPI, crates.io, Go Proxy, Packagist, Maven, Swift, and Godot.
