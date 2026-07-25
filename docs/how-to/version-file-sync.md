# Version File Sync

## Overview

big-release can automatically write a version file after each release. This is useful for applications that need to display their version at runtime (e.g., in a footer, `/about` endpoint, or health-check payload).

## Configuration

Add `versionFile` to your `.big-release.yml`:

```yaml
versionFile:
  path: VERSION
  template: "{{.Version}}"
```

### Options

| Option | Default | Description |
|--------|---------|-------------|
| `path` | `VERSION` | Path to the version file |
| `template` | `{{.Version}}` | Go template for the file content |

### Template Variables

The template has access to these variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Version}}` | Semantic version | `1.2.3` |
| `{{.Type}}` | Release type | `major`, `minor`, `patch` |
| `{{.Channel}}` | Release channel | `latest`, `beta`, `alpha` |
| `{{.GitTag}}` | Git tag name | `v1.2.3` |

## How It Works

1. big-release analyzes commits and determines the next version
2. The git plugin creates and pushes the tag
3. The `versionFile` plugin writes the version file
4. The git plugin commits the version file with `[skip ci]`

## Examples

### Simple VERSION file

```yaml
versionFile:
  path: VERSION
```

Result: `1.2.3`

### Go version constant

```yaml
versionFile:
  path: internal/version/version.go
  template: |
    package version

    const Version = "{{.Version}}"
```

### JSON version file

```yaml
versionFile:
  path: public/version.json
  template: |
    {
      "version": "{{.Version}}",
      "type": "{{.Type}}"
    }
```

## CI Workaround

If you don't want the plugin approach, you can use a CI step after big-release:

```yaml
- name: Write VERSION file
  run: |
    LATEST_TAG="$(git describe --tags --abbrev=0)"
    VERSION="${LATEST_TAG#v}"
    echo "$VERSION" > VERSION
    git add VERSION
    git -c user.name="big-release-bot" -c user.email="bot@example.com" \
      commit -m "chore: update VERSION to ${VERSION} [skip ci]"
    git push
```

## Troubleshooting

### Version file not committed

Make sure the git plugin has `postTagAssets` configured:

```yaml
plugins: [changelog, git, github, versionfile]
pluginConfigs:
  git:
    postTagAssets:
      - VERSION
```

### Template error

Check your template syntax. Common mistakes:
- Missing `{{` or `}}`
- Using `{{version}}` instead of `{{.Version}}`
- Using JavaScript template literals `${version}` instead of Go templates
