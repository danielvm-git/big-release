# Dogfooding Plan: npm Demo Package

## Overview

This plan walks through creating a demo npm package that uses `big-release` for versioning, changelog generation, and publishing. The goal is to validate the full pipeline and demonstrate version progression from `0.1.0` through `0.2.0` to `1.0.0`.

## Repository Setup

### 1. Create Repository

```bash
# Create directory
mkdir big-release-demo-npm
cd big-release-demo-npm
git init
git checkout -b main
```

### 2. Initialize npm Package

```bash
npm init -y
```

Edit `package.json`:
```json
{
  "name": "@danielvm-git/big-release-demo",
  "version": "0.0.0",
  "description": "Demo package for dogfooding big-release",
  "main": "index.js",
  "scripts": {
    "test": "echo \"OK\""
  },
  "keywords": ["demo", "big-release"],
  "author": "danielvm-git",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/danielvm-git/big-release-demo-npm.git"
  }
}
```

### 3. Add big-release Configuration

Create `.big-release.yml`:
```yaml
branches:
  - name: main

tagFormat: "v${version}"

plugins:
  - changelog
  - git

publishers:
  npm:
    enabled: true
```

### 4. Initial Files

Create `index.js`:
```javascript
function greet(name) {
  return `Hello, ${name}!`;
}

module.exports = { greet };
```

Create `CHANGELOG.md`:
```markdown
# Changelog

All notable changes to this project will be documented in this file.
```

### 5. Initial Commit

```bash
git add -A
git commit -m "feat: initial project setup"
git remote add origin https://github.com/danielvm-git/big-release-demo-npm.git
git push -u origin main
```

## Release Walkthrough

### Release 1: v0.1.0 (First Release)

**Step 1: Add a feature**
```bash
# Edit index.js to add a new function
cat > index.js << 'EOF'
function greet(name) {
  if (!name) return "Hello, World!";
  return `Hello, ${name}!`;
}

function farewell(name) {
  if (!name) return "Goodbye, World!";
  return `Goodbye, ${name}!`;
}

module.exports = { greet, farewell };
EOF

git add -A
git commit -m "feat: add farewell function"
```

**Step 2: Run big-release**
```bash
# Dry-run first
big-release release --dry-run

# If satisfied, run for real
big-release release
```

**Expected output:**
- Version: `0.1.0` (first release, using InitialVersion from config)
- Tag: `v0.1.0`
- Changelog updated with "feat: add farewell function"
- Package published to npm

### Release 2: v0.2.0 (Minor Version)

**Step 1: Add more features**
```bash
cat > index.js << 'EOF'
function greet(name, language = "en") {
  if (!name) name = "World";
  
  const greetings = {
    en: `Hello, ${name}!`,
    es: `Hola, ${name}!`,
    fr: `Bonjour, ${name}!`
  };
  
  return greetings[language] || greetings.en;
}

function farewell(name, language = "en") {
  if (!name) name = "World";
  
  const farewells = {
    en: `Goodbye, ${name}!`,
    es: `Adiós, ${name}!`,
    fr: `Au revoir, ${name}!`
  };
  
  return farewells[language] || farewells.en;
}

module.exports = { greet, farewell };
EOF

git add -A
git commit -m "feat(i18n): add language support for greet and farewell"
```

**Step 2: Run big-release**
```bash
big-release release
```

**Expected output:**
- Version: `0.2.0` (minor bump from feat commit)
- Tag: `v0.2.0`
- Changelog updated with "feat(i18n): add language support"

### Release 3: v0.2.1 (Patch Version)

**Step 1: Fix a bug**
```bash
cat > index.js << 'EOF'
function greet(name, language = "en") {
  if (!name) name = "World";
  
  const greetings = {
    en: `Hello, ${name}!`,
    es: `Hola, ${name}!`,
    fr: `Bonjour, ${name}!`
  };
  
  return greetings[language] || greetings.en;
}

function farewell(name, language = "en") {
  if (!name) name = "World";
  
  const farewells = {
    en: `Goodbye, ${name}!`,
    es: `Adiós, ${name}!`,
    fr: `Au revoir, ${name}!`
  };
  
  return farewells[language] || farewells.en;
}

module.exports = { greet, farewell };
EOF

git add -A
git commit -m "fix(i18n): handle unsupported language gracefully"
```

**Step 2: Run big-release**
```bash
big-release release
```

**Expected output:**
- Version: `0.2.1` (patch bump from fix commit)
- Tag: `v0.2.1`
- Changelog updated with "fix(i18n): handle unsupported language"

### Release 4: v1.0.0 (Major Version)

**Step 1: Breaking change**
```bash
cat > index.js << 'EOF'
// v1.0.0 - Breaking change: API now returns objects instead of strings

function greet(name, options = {}) {
  const { language = "en", formal = false } = options;
  
  if (!name) name = "World";
  
  const greetings = {
    en: formal ? `Good day, ${name}.` : `Hello, ${name}!`,
    es: formal ? `Buenos días, ${name}.` : `Hola, ${name}!`,
    fr: formal ? `Bonjour, ${name}.` : `Bonjour, ${name}!`
  };
  
  return {
    message: greetings[language] || greetings.en,
    name,
    language,
    formal
  };
}

function farewell(name, options = {}) {
  const { language = "en", formal = false } = options;
  
  if (!name) name = "World";
  
  const farewells = {
    en: formal ? `Farewell, ${name}.` : `Goodbye, ${name}!`,
    es: formal ? `Adiós, ${name}.` : `Adiós, ${name}!`,
    fr: formal ? `Au revoir, ${name}.` : `Au revoir, ${name}!`
  };
  
  return {
    message: farewells[language] || farewells.en,
    name,
    language,
    formal
  };
}

module.exports = { greet, farewell };
EOF

git add -A
git commit -m "feat(api)!: change return type from string to object

BREAKING CHANGE: greet() and farewell() now return objects instead of strings.
Update your code to use .message property to get the string."
```

**Step 2: Run big-release**
```bash
big-release release
```

**Expected output:**
- Version: `1.0.0` (major bump from breaking change)
- Tag: `v1.0.0`
- Changelog updated with breaking change section

## Version Progression Summary

| Release | Version | Commit Type | Description |
|---------|---------|-------------|-------------|
| 1 | 0.1.0 | feat | First release |
| 2 | 0.2.0 | feat | Minor feature addition |
| 3 | 0.2.1 | fix | Bug fix |
| 4 | 1.0.0 | feat! | Breaking change |

## Configuration Reference

### .big-release.yml

```yaml
# Branch configuration
branches:
  - name: main

# Tag format (supports ${version} placeholder)
tagFormat: "v${version}"

# Initial version for first release
initialVersion: "0.1.0"

# Plugins to use
plugins:
  - changelog    # Generates CHANGELOG.md
  - git          # Creates git tags and commits
  - github       # Creates GitHub releases (optional)

# Publishers to use
publishers:
  npm:
    enabled: true
```

## Environment Variables

For publishing, set these tokens:
```bash
export NPM_TOKEN="your-npm-token"
export GITHUB_TOKEN="your-github-token"
```

## Verification Commands

```bash
# Check version
big-release version

# Validate configuration
big-release validate

# Check health
big-release health

# Dry-run release
big-release release --dry-run

# Verbose output
big-release release --dry-run --verbose
```

## Troubleshooting

### Common Issues

1. **"failed to get repository URL"**
   - Ensure git remote is configured: `git remote add origin <url>`

2. **"NPM_TOKEN not set"**
   - Set token: `export NPM_TOKEN="your-token"`

3. **"no releasable commits found"**
   - Ensure commits follow Conventional Commits format

4. **"branch not in release branches"**
   - Add branch to `.big-release.yml` config
