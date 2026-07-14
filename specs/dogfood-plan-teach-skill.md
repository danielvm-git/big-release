# Dogfooding Plan: Teach Skill + Demo Repos

## Overview

Use the `teach` skill structure to create an HTML teaching site for big-release, then create demo repositories for each publisher module that reproduce the site with ecosystem-specific content.

## Architecture

```
big-release-teach/                    # Main teaching site
  MISSION.md                          # Why learn big-release
  RESOURCES.md                        # SemVer, Conventional Commits, semantic-release links
  NOTES.md                            # Teaching preferences
  GLOSSARY.md                         # big-release terminology
  lessons/
    0001-what-is-big-release.html     # Overview of the tool
    0002-conventional-commits.html    # Commit message format
    0003-versioning.html              # SemVer and version calculation
    0004-greenfield-setup.html        # Starting a new project
    0005-brownfield-migration.html    # Adding to existing project
    0006-publishing.html              # Publishing to registries
  reference/
    cheat-sheet.html                 # Quick reference card
    config-reference.html            # .big-release.yml options
  assets/
    styles.css                       # Shared styles
    quiz.js                          # Quiz widget

big-release-demo-npm/                # npm module demo
  lessons/
    0001-npm-setup.html              # npm-specific setup
    0002-npm-publishing.html         # npm publishing flow
    0003-npm-versioning.html         # npm version conventions
  reference/
    npm-cheat-sheet.html             # npm-specific reference

big-release-demo-pypi/               # PyPI module demo
  lessons/
    0001-pypi-setup.html
    0002-pypi-publishing.html
  reference/
    pypi-cheat-sheet.html

big-release-demo-crates/             # crates.io module demo
  lessons/
    0001-crates-setup.html
    0002-crates-publishing.html
  reference/
    crates-cheat-sheet.html

big-release-demo-go/                 # Go Proxy module demo
  lessons/
    0001-go-setup.html
    0002-go-publishing.html
  reference/
    go-cheat-sheet.html
```

## Lesson Content Plan

### Main Site Lessons

#### 0001-what-is-big-release.html
- What big-release does (version, changelog, publish)
- Why it exists (unified tool for multiple ecosystems)
- How it compares to semantic-release
- Interactive: Show a release being calculated

#### 0002-conventional-commits.html
- Commit message format: `type(scope): description`
- Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Breaking changes: `feat!:` or `BREAKING CHANGE:` footer
- Interactive: Quiz on commit message format

#### 0003-versioning.html
- SemVer: MAJOR.MINOR.PATCH
- How commit types map to version bumps
- Initial version configuration
- Interactive: Calculate version from commits

#### 0004-greenfield-setup.html
- Starting a new project with big-release
- Configuration file setup
- First release workflow
- Step-by-step walkthrough

#### 0005-brownfield-migration.html
- Adding big-release to existing project
- Handling existing tags
- Migrating from other tools
- Step-by-step walkthrough

#### 0006-publishing.html
- Publishing to npm, PyPI, crates.io, Go Proxy
- Token setup
- Dry-run mode
- Verification

### Module-Specific Lessons

Each demo repo would have lessons specific to that ecosystem:

**npm demo:**
- package.json configuration
- Scoped packages
- npm token setup
- Unpublishing (if needed)

**PyPI demo:**
- setup.py / pyproject.toml configuration
- PyPI token setup
- TestPyPI for testing

**crates.io demo:**
- Cargo.toml configuration
- crates.io token setup
- Documentation publishing

**Go Proxy demo:**
- go.mod configuration
- Tag-based versioning
- Proxy caching

## Execution Plan

### Phase 1: Create Main Teaching Site
1. Create `big-release-teach/` directory structure
2. Write MISSION.md, RESOURCES.md, GLOSSARY.md
3. Create shared assets (styles.css, quiz.js)
4. Write 6 HTML lessons
5. Create reference documents

### Phase 2: Create Demo Repos
1. Create `big-release-demo-npm/` with npm-specific lessons
2. Create `big-release-demo-pypi/` with PyPI-specific lessons
3. Create `big-release-demo-crates/` with crates-specific lessons
4. Create `big-release-demo-go/` with Go-specific lessons

### Phase 3: Test and Validate
1. Open each lesson in browser
2. Verify interactive elements work
3. Test big-release dry-run on each demo repo
4. Document any issues

## Teaching Philosophy

Following the teach skill philosophy:
- **Knowledge**: SemVer, Conventional Commits, registry-specific rules
- **Skills**: Setting up big-release, creating releases, troubleshooting
- **Wisdom**: When to use which publisher, handling edge cases

Each lesson should:
- Be self-contained HTML
- Have beautiful typography (Tufte-inspired)
- Include interactive elements
- Link to other lessons
- Recommend primary sources
- Be completable quickly

## Success Criteria

- [ ] Main teaching site has 6 lessons covering all key topics
- [ ] Each demo repo has 2-3 ecosystem-specific lessons
- [ ] All lessons are beautiful HTML with good typography
- [ ] Interactive elements work (quizzes, calculators)
- [ ] big-release dry-run passes on each demo repo
- [ ] Lessons reference real documentation
