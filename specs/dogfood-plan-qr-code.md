# Dogfooding Plan: QR Code Generator with big-release

## Overview

Build a QR Code Generator npm package and use `big-release` to version it through 4 milestones. Demonstrates real-world version progression from `0.1.0` to `1.0.0`.

## Why QR Code Generator

- Pure computation (no auth/DB)
- Visual output (fun to demo)
- Natural feature progression
- Real-world utility
- Small enough to build quickly

## Project Structure

```
big-release-demo-qr/
  .big-release.yml
  package.json
  index.js
  CHANGELOG.md
  README.md
  src/
    generator.js
    options.js
    renderer.js
```

## Configuration

### .big-release.yml

```yaml
branches:
  - name: main

tagFormat: "v${version}"
initialVersion: "0.1.0"

plugins:
  - changelog
  - git

publishers:
  npm:
    enabled: true
```

### package.json

```json
{
  "name": "@danielvm-git/qr-code-generator",
  "version": "0.0.0",
  "description": "Generate QR codes from text with customizable options",
  "main": "index.js",
  "scripts": {
    "test": "node test.js"
  },
  "dependencies": {
    "qrcode": "^1.5.3"
  },
  "keywords": ["qr", "qrcode", "generator", "barcode"],
  "author": "danielvm-git",
  "license": "MIT"
}
```

## Release Walkthrough

### Release 1: v0.1.0 — Basic QR Generation

**Commits:**
```bash
git commit -m "feat: initial project setup"
git commit -m "feat(generator): add basic QR code generation from text"
git commit -m "feat(renderer): add PNG output support"
git commit -m "test: add basic generation tests"
```

**Run release:**
```bash
big-release release --dry-run   # preview
big-release release              # publish v0.1.0
```

**Expected output:**
- Version: `0.1.0`
- Tag: `v0.1.0`
- Changelog:
  ```
  ## 0.1.0 (2026-07-13)

  ### Features
  - feat: initial project setup
  - feat(generator): add basic QR code generation from text
  - feat(renderer): add PNG output support

  ### Tests
  - test: add basic generation tests
  ```

---

### Release 2: v0.2.0 — Options and Formats

**Commits:**
```bash
git commit -m "feat(options): add size customization"
git commit -m "feat(options): add error correction level (L/M/Q/H)"
git commit -m "feat(options): add foreground and background colors"
git commit -m "feat(renderer): add SVG output support"
git commit -m "fix(options): validate size constraints"
git commit -m "test: add option validation tests"
```

**Run release:**
```bash
big-release release
```

**Expected output:**
- Version: `0.2.0` (minor bump from `feat:` commits)
- Tag: `v0.2.0`

---

### Release 3: v0.3.0 — Logo Overlay

**Commits:**
```bash
git commit -m "feat(logo): add center logo overlay"
git commit -m "feat(logo): add logo size and margin options"
git commit -m "fix(logo): handle transparent PNG logos"
git commit -m "docs: add API documentation with examples"
```

**Run release:**
```bash
big-release release
```

**Expected output:**
- Version: `0.3.0` (minor bump)
- Tag: `v0.3.0`

---

### Release 4: v1.0.0 — Breaking Change

**Commits:**
```bash
git commit -m "feat(api)!: return structured result object

BREAKING CHANGE: generate() now returns { buffer, format, width, height }
instead of raw Buffer. Access the buffer via result.buffer."
```

**Run release:**
```bash
big-release release
```

**Expected output:**
- Version: `1.0.0` (major bump from `feat!:` + `BREAKING CHANGE:`)
- Tag: `v1.0.0`
- Changelog includes:
  ```
  ### BREAKING CHANGES
  - feat(api): return structured result object
  ```

---

## Version Progression Summary

| Release | Version | Commits | Bump Reason |
|---------|---------|---------|-------------|
| 1 | `0.1.0` | `feat:` x3, `test:` x1 | First release (initialVersion) |
| 2 | `0.2.0` | `feat:` x3, `fix:` x1, `test:` x1 | Minor (feat commits) |
| 3 | `0.3.0` | `feat:` x2, `fix:` x1, `docs:` x1 | Minor (feat commits) |
| 4 | `1.0.0` | `feat!:` x1 | Major (BREAKING CHANGE) |

## What This Demonstrates

1. **Conventional Commits drive versioning** — `feat:` = minor, `fix:` = patch, `feat!:` = major
2. **Changelog auto-generation** — Grouped by type, cumulative across releases
3. **Initial version config** — Starts at 0.1.0, not 1.0.0
4. **Git tagging** — Clean `v${version}` tags
5. **npm publishing** — Version bump in package.json + publish
6. **Real-world progression** — Shows how a project grows from MVP to v1.0

## Setup Commands

```bash
# Create repo
mkdir big-release-demo-qr
cd big-release-demo-qr
git init && git checkout -b main

# Initialize
npm init -y
# Edit package.json (see above)
# Create .big-release.yml (see above)
# Create src/ structure

# First commit
git add -A
git commit -m "feat: initial project setup"

# Add remote
git remote add origin https://github.com/danielvm-git/big-release-demo-qr.git
git push -u origin main

# Release
big-release release
```

## Environment Variables

```bash
export NPM_TOKEN="your-npm-token"
export GITHUB_TOKEN="your-github-token"
```
