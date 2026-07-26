# Configuration Reference

big-release uses a YAML configuration file (`.big-release.yml`) to define release behavior.

## Configuration File Locations

big-release looks for configuration in the following locations (in order):

1. `.big-release.yml`
2. `.big-release.json`
3. `.big-release.yaml`
4. `big-release.config.js`

## Configuration Options

### branches

Defines which branches trigger releases.

```yaml
branches:
  - main                    # Release branch
  - next                    # Release branch
  - "N.x"                   # Maintenance branch (e.g., 1.x, 2.x)
  - name: beta
    prerelease: true        # Prerelease branch
    channel: beta           # Distribution channel
```

**Branch Types:**

| Type | Description | Example |
|------|-------------|---------|
| Release | Default branch type | `main`, `next` |
| Maintenance | Patches older versions | `1.x`, `2.x.x` |
| Prerelease | Pre-release versions | `beta`, `alpha`, `canary` |

### tagFormat

The format for git tags.

```yaml
tagFormat: "v${version}"    # Default: v1.2.3
tagFormat: "${version}"     # Alternative: 1.2.3
tagFormat: "release/${version}"  # Alternative: release/1.2.3
```

### plugins

List of plugins to execute.

```yaml
plugins:
  - changelog              # Generate CHANGELOG.md
  - git                    # Commit changes back to repo
  - github                 # Create GitHub releases
```

#### git plugin options

| Option | Type | Default | Description |
|---|---|---|---|
| `message` | string | `chore(release): <version> [skip ci]` | Release commit message template |
| `assets` | []string | none | Globs to stage. With none set, no release commit is made |
| `postTagAssets` | []string | none | Globs committed after the tag is created |
| `tagOnly` | bool | `false` | Push tags without pushing the release commit |

##### tagOnly — releasing onto a protected branch

By default a failed push **fails the release**, matching semantic-release. If
the remote refuses the commit, big-release reports it with git's stderr and a
hint, and removes the local tag.

Some setups cannot grant the release identity push access to the default
branch at all — GitHub branch protection or a ruleset with no bypass available
(bypass actors need an organization), GitLab protected branches, or any
`pre-receive` policy. Set `tagOnly` to publish the tag and skip the commit:

```yaml
plugins:
  - git:
      tagOnly: true
```

This is deliberately explicit. The version bump and CHANGELOG will **not** land
on the branch, so keep them out of `assets` or regenerate them elsewhere.
Prefer granting push access when you can; reach for `tagOnly` when you can't.

### publishers

Language-specific package publishers.

```yaml
publishers:
  npm:
    enabled: true
    registry: "https://registry.npmjs.org"
  pypi:
    enabled: true
    registry: "https://pypi.org"
  crates:
    enabled: true
  goproxy:
    enabled: true
  packagist:
    enabled: true
  maven:
    enabled: true
```

## Environment Variables

big-release uses the following environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `GITHUB_TOKEN` | GitHub API token | For GitHub releases |
| `NPM_TOKEN` | npm authentication token | For npm publishing |
| `PYPI_TOKEN` | PyPI authentication token | For PyPI publishing |
| `CARGO_TOKEN` | crates.io authentication token | For crates.io publishing |

## Example Configurations

### JavaScript/TypeScript Project

```yaml
branches:
  - main
  - next
  - "N.x"
  - name: beta
    prerelease: true

tagFormat: "v${version}"

plugins:
  - changelog
  - git
  - github

publishers:
  npm:
    enabled: true
```

### Python Project

```yaml
branches:
  - main
  - develop
  - "N.x"

tagFormat: "v${version}"

plugins:
  - changelog
  - git
  - github

publishers:
  pypi:
    enabled: true
```

### Rust Project

```yaml
branches:
  - main
  - next
  - "N.x"

tagFormat: "v${version}"

plugins:
  - changelog
  - git
  - github

publishers:
  crates:
    enabled: true
```

### Godot Project

```yaml
branches:
  - main

tagFormat: "v${version}"

plugins:
  - changelog
  - git
  - github
  - exec

exec:
  prepareCmd: |
    if [ -f export_presets.cfg ]; then
      sed -i 's/version\\/name=".*"/version\\/name="${nextRelease.version}"/g' export_presets.cfg
    fi

publishers:
  godot:
    enabled: true
```

## Validation Rules

1. **At least one branch** must be configured
2. **Maximum 3 release branches** allowed
3. **Tag format** must contain `${version}` placeholder
4. **Maintenance branches** must follow `N.x` or `N.x.x` pattern
5. **Prerelease branches** must have valid prerelease identifier

## Advanced Configuration

### Custom Exec Commands

Execute custom commands during the release process:

```yaml
plugins:
  - exec

exec:
  prepareCmd: |
    echo "Preparing release ${nextRelease.version}"
    npm run build
  publishCmd: |
    echo "Publishing ${nextRelease.version}"
    npm publish
```

### Multiple Channels

Publish to different distribution channels:

```yaml
branches:
  - main                    # Default channel
  - name: beta
    channel: beta           # Beta channel
  - name: alpha
    channel: alpha          # Alpha channel
```

### Monorepo Support

Configure for monorepo projects:

```yaml
branches:
  - main

publishers:
  npm:
    enabled: true
    workspaces: true
```
