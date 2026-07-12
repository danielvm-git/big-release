# danielvm-git/big-release — AI Agents

> **Multi-agent context** — This file is the canonical project context for **Cline**, **Aider**, **OpenCode**, and other AGENTS.md-native tools. Claude Code and Cursor read it via the `CLAUDE.md` symlink.

Read CONVENTIONS.md before any GitHub or git operation.

<!-- BEGIN bigpowers:context-routing -->
## Context Routing

| Glob / trigger | Load first |
|----------------|------------|
| `internal/**/*.go` | `docs/reference/` for API contracts |
| `publishers/**` | Language-specific docs under `docs/how-to/publishers/` |
| `specs/**` | Matching spec file |
| Default / session start | This file → `CONVENTIONS.md` → `specs/state.yaml` |
<!-- END bigpowers:context-routing -->

<!-- BEGIN bigpowers:learned-preferences -->
## Learned User Preferences

- Follow bigpowers fix-or-log — never dismiss gate failures.
- Use `rtk`-prefixed shell commands for git, test, build output.
- Prefer `gh` CLI over raw `git push` / `curl` for GitHub operations.

## Workspace Facts

- This repo is **big-release** — a unified, multi-language release tool.
- Written in Go for single-binary distribution.
- Core algorithm ported from semantic-release (Node.js).
- Plugin system for language-specific publishers.
- Stack: Go / YAML / Bash
- 35+ portfolio repos will consume this tool.
<!-- END bigpowers:learned-preferences -->

<!-- BEGIN bigpowers:project -->
## Project

Unified release tool for danielvm-git — determines version, generates changelog, publishes packages across all languages.
Stack: Go, YAML, Bash

## Commands

| Action | Command |
|--------|---------|
| Preflight | `make preflight` |
| Lint | `make lint` |
| Test | `make test` |
| Build | `make build` |
| CI | `gh pr checks` (when a PR is open) |
| Health | `./bin/big-release health` |
| Setup | `bash scripts/setup.sh` |

## Observability

| What | Command |
|------|---------|
| Health check (JSON) | `./bin/big-release health` |
| View logs | `./bin/big-release --verbose release` (structured JSON via zap) |
| Validate config | `./bin/big-release validate` |
| Show version | `./bin/big-release version` |
| CI status | `gh run list --limit 5` |

## Architecture

Three layers: **algorithm** (internal/algorithm/) → **git operations** (internal/git/) → **publishers** (internal/publishers/) → **CLI** (cmd/big-release/).

## Conventions

- Conventional Commits on all changes; `feat:` = minor, `fix:` = patch
- No direct work on `main` — feature branches only, squash-merge via `scripts/land-branch.sh`
- `guard-git` hooks block dangerous git operations
- `bigpowers` owns methodology; this repo owns release automation
- Branch protection: require PR, require CI pass, no force-push
- Docs follow Diátaxis structure (tutorials / how-to / reference / explanation)

## Never

- Never dismiss reproducible gate failures as pre-existing or out of scope
- Never proceed on red Preflight or red CI — invoke quick-fix or fix-bug first
- Never edit another repo's files from this repo
- Never duplicate bigpowers methodology — this repo extends, not replaces
- **Never emit `Co-authored-by:` footers** in any commit (`git commit -m`). All commits must appear as if authored solely by the human user. Enforced by pre-tool-use hook, `land-branch.sh`, and CI. No exceptions.

## Agent Rules

- **Workflow Mandate:** Use bigpowers skills (e.g. `plan-work`, `develop-tdd`) for structured work.
- **P1 Git Attribution Rule:** Never include `Co-authored-by:` / `Co-Authored-By:` in any commit. See CONVENTIONS.md § Git Attribution. The hook blocks them at `git commit`; CI blocks them at PR. You will be blocked.
- **Always Green:** Preflight and CI must be green before forward work.
- Read specs/ and CONVENTIONS.md before writing code.
- Write the minimum code that solves the stated problem.
- All planning output goes in specs/.
- New workflow templates must include `concurrency` group and `timeout-minutes`.
<!-- END bigpowers:project -->
