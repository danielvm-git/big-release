# Dogfooding Plan: Todo API with big-release

## Overview

Build a Todo API (Node.js/Express) and use `big-release` to version it through 4 milestones. Demonstrates real-world version progression from `0.1.0` to `1.0.0`.

## Project Structure

```
big-release-demo-todo/
  .big-release.yml
  package.json
  index.js
  CHANGELOG.md
  README.md
  src/
    db.js
    auth.js
    routes/
      tasks.js
      users.js
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
  "name": "@danielvm-git/todo-api",
  "version": "0.0.0",
  "description": "A Todo API with auth, CRUD, and collaboration",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "test": "node test.js"
  },
  "dependencies": {
    "express": "^4.18.0",
    "uuid": "^9.0.0"
  },
  "keywords": ["todo", "api", "task-management"],
  "author": "danielvm-git",
  "license": "MIT"
}
```

## Release Walkthrough

### Release 1: v0.1.0 — Basic CRUD

**Commits:**
```bash
git commit -m "feat: initial project setup with Express server"
git commit -m "feat(db): add in-memory task storage"
git commit -m "feat(tasks): add CRUD endpoints for tasks"
git commit -m "test: add task CRUD tests"
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
  - feat: initial project setup with Express server
  - feat(db): add in-memory task storage
  - feat(tasks): add CRUD endpoints for tasks

  ### Tests
  - test: add task CRUD tests
  ```

---

### Release 2: v0.2.0 — Authentication

**Commits:**
```bash
git commit -m "feat(auth): add user registration and login"
git commit -m "feat(auth): add JWT token validation middleware"
git commit -m "feat(users): add user profile endpoint"
git commit -m "fix(auth): handle invalid token gracefully"
git commit -m "test: add authentication tests"
```

**Run release:**
```bash
big-release release
```

**Expected output:**
- Version: `0.2.0` (minor bump from `feat:` commits)
- Tag: `v0.2.0`
- Changelog appends new section above previous

---

### Release 3: v0.3.0 — Task Sharing

**Commits:**
```bash
git commit -m "feat(sharing): add task assignment to users"
git commit -m "feat(sharing): add shared task view endpoint"
git commit -m "fix(sharing): prevent duplicate assignments"
git commit -m "docs: update API documentation"
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
git commit -m "feat(api)!: change task response format

BREAKING CHANGE: tasks now return nested objects instead of flat structures.
Update client code to use task.title instead of taskTitle."
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
  - feat(api): change task response format
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
mkdir big-release-demo-todo
cd big-release-demo-todo
git init && git checkout -b main

# Initialize
npm init -y
# Edit package.json (see above)
# Create .big-release.yml (see above)
# Create src/ structure

# First commit
git add -A
git commit -m "feat: initial project setup with Express server"

# Add remote
git remote add origin https://github.com/danielvm-git/big-release-demo-todo.git
git push -u origin main

# Release
big-release release
```

## Environment Variables

```bash
export NPM_TOKEN="your-npm-token"
export GITHUB_TOKEN="your-github-token"
```
