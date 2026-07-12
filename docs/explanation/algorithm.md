# Algorithm Deep Dive

This document explains the core algorithm behind big-release in detail.

## Overview

big-release follows a 12-phase algorithm to automate the release process:

```
Initialize → Analyze Branch → Verify Auth → Find Last Release → Get Commits → Analyze Commits → Calculate Version → Generate Notes → Create Tag → Publish → GitHub Release → Notify
```

## Phase 1: Initialize

**Purpose:** Set up the release environment

**Steps:**
1. Check if running in CI environment
2. Load configuration file
3. Set up git author (for automated commits)
4. Initialize logging

**Key Decisions:**
- If not in CI and not dry-run, force dry-run mode
- If triggered by PR, skip release

## Phase 2: Analyze Branch

**Purpose:** Determine the branch type and configuration

**Steps:**
1. Get current branch name from CI environment
2. Match branch against configured branches
3. Classify branch type:
   - **Release:** `main`, `next`
   - **Maintenance:** `1.x`, `2.x.x`
   - **Prerelease:** `beta`, `alpha`

**Example:**
```yaml
branches:
  - main                    # Release branch
  - "1.x"                   # Maintenance branch
  - name: beta
    prerelease: true        # Prerelease branch
```

## Phase 3: Verify Auth

**Purpose:** Ensure we can push to the repository

**Steps:**
1. Get repository URL from git remote
2. Attempt dry-run push
3. If failed, check if branch is up-to-date
4. If behind remote, skip release

**Error Handling:**
- `EGITNOPERMISSION`: Cannot push to repository
- `EBRANCHBEHIND`: Branch is behind remote

## Phase 4: Find Last Release

**Purpose:** Find the most recent release

**Steps:**
1. Get all tags for current branch
2. Filter tags matching tag format
3. Parse version numbers
4. Sort by semver precedence
5. Return highest version

**Example:**
```
Tags: v1.0.0, v1.1.0, v1.2.0-beta.1
Last Release: v1.2.0-beta.1
```

## Phase 5: Get Commits

**Purpose:** Retrieve commits since last release

**Steps:**
1. Determine commit range (from last release to HEAD)
2. Run `git log` with formatting
3. Parse commit messages
4. Filter out skipped commits

**Example:**
```
git log v1.0.0..HEAD --pretty=format:"%H|%s|%an|%ae|%ai|%b"
```

## Phase 6: Analyze Commits

**Purpose:** Determine the release type from commits

**Steps:**
1. Parse each commit using Conventional Commits
2. Extract type, scope, breaking changes
3. Determine bump type:
   - `feat` → minor
   - `fix`, `perf` → patch
   - Breaking change → major
4. Return highest priority bump

**Conventional Commits Format:**
```
<type>(<scope>): <description>

BREAKING CHANGE: <description>
```

**Examples:**
```
feat(auth): add OAuth2 support          → minor
fix!: remove deprecated API             → major
perf: optimize database queries         → patch
docs: update README                     → no release
```

## Phase 7: Calculate Version

**Purpose:** Calculate the next semantic version

**Algorithm:**
```
IF no last release:
  IF prerelease branch:
    version = "1.0.0-{preid}.1"
  ELSE:
    version = "1.0.0"

ELSE:
  PARSE last version (major.minor.patch)
  
  IF prerelease branch:
    IF current version is prerelease with same preid:
      version = increment prerelease number
    ELSE:
      version = increment based on type + "-{preid}.1"
  
  ELSE:
    version = increment based on type
```

**Increment Rules:**
- Major: `X.y.z` → `(X+1).0.0`
- Minor: `x.Y.z` → `x.(Y+1).0`
- Patch: `x.y.Z` → `x.y.(Z+1)`
- Prerelease: `x.y.z-pre.N` → `x.y.z-pre.(N+1)`

## Phase 8: Generate Notes

**Purpose:** Create release notes from commits

**Steps:**
1. Group commits by type:
   - Breaking changes
   - Features
   - Bug fixes
   - Performance improvements
2. Format each commit
3. Add comparison link
4. Hide sensitive information

**Example Output:**
```markdown
## Features

- **auth:** add OAuth2 support (abc1234)

## Bug Fixes

- **api:** fix rate limiting (def5678)

---

Full Changelog: comparing changes from v1.0.0 to v1.1.0
```

## Phase 9: Create Tag

**Purpose:** Create the git tag

**Steps:**
1. Format tag name (e.g., `v1.1.0`)
2. Validate tag name
3. Create annotated tag
4. Add channel notes
5. Push tag to remote
6. Push notes to remote

## Phase 10: Publish

**Purpose:** Publish the package to registry

**Steps:**
1. Detect package type (npm, PyPI, crates, etc.)
2. Update version in manifest
3. Execute publish command
4. Verify publication

**Language-Specific:**
- **npm:** `npm publish`
- **PyPI:** `twine upload dist/*`
- **crates.io:** `cargo publish`

## Phase 11: GitHub Release

**Purpose:** Create GitHub release with notes

**Steps:**
1. Create release via GitHub API
2. Attach release notes
3. Upload assets (if any)

## Phase 12: Notify

**Purpose:** Notify stakeholders of release

**Steps:**
1. Log success message
2. Output release notes (if dry-run)
3. Return release information

## Error Handling

### Error Types

| Error | Description | Action |
|-------|-------------|--------|
| `EGITNOPERMISSION` | Cannot push to repository | Check credentials |
| `EINVALIDTAGFORMAT` | Invalid tag format | Update configuration |
| `EINVALIDNEXTVERSION` | Version out of range | Check branch configuration |
| `EPLUGIN` | Plugin error | Check plugin configuration |
| `EPUBLISH` | Publishing failed | Check registry credentials |

### Recovery Strategies

1. **Retry:** For transient network errors
2. **Rollback:** For failed publications
3. **Skip:** For non-critical failures
4. **Abort:** For critical failures

## Performance Considerations

- **Git operations:** Use shallow clones when possible
- **Commit parsing:** Cache parsed commits
- **Network calls:** Batch API requests
- **Parallel execution:** Run independent plugins concurrently

## Security Considerations

- **Credentials:** Never log sensitive information
- **Tokens:** Use environment variables
- **Verification:** Verify publication success
- **Rollback:** Support rollback on failure

## Testing Strategy

- **Unit tests:** Test each phase independently
- **Integration tests:** Test complete workflow
- **Edge cases:** Test error conditions
- **Performance tests:** Test with large repositories

## Future Improvements

1. **Monorepo support:** Multiple packages per repository
2. **Custom plugins:** User-defined plugins
3. **Web UI:** Release history visualization
4. **Analytics:** Release metrics and insights
