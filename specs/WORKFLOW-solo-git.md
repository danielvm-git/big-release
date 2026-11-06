# Workflow: Solo Git Mode

## Overview

This repository operates in **solo-git** mode (`workflow_mode: solo-git` in `specs/state.yaml`).
In solo-git mode, development follows strict trunk-based principles with local feature branches and scripted landing to `main`.

## Principles

1. **No Direct Work on Main:** All changes must originate on an isolated branch (`kickoff-branch`).
2. **Always Green:** Local Preflight (`make preflight`) and CI checks must be green before landing.
3. **P1 Git Attribution:** No AI agent attribution (`Co-authored-by:` or variants) is permitted in commit history.
4. **Local Squash Landing:** Branches are landed cleanly using `scripts/land-branch.sh`.

## Branch Lifecycle

```
main ─────────── (branch off) ───────────> feature branch
  │                                              │
  │                                        (TDD & Red-Green)
  │                                              │
  │                                        (Preflight PASS)
  │                                              │
  <──────── bash scripts/land-branch.sh ─────────┘
  │ (squash + land to main + push)
```

## Landing Command

```bash
bash scripts/land-branch.sh <branch-name> "<conventional commit message>"
```

### Landing Verification Sequence

1. Working tree cleanliness check
2. `make preflight` execution
3. Rebase onto latest `origin/main`
4. Conventional Commits verification
5. Attribution footers block check
6. Squash merge into `main`
7. Push to `origin/main`
8. Local feature branch cleanup
