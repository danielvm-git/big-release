# big-release — Conventions

## Conventional Commits & Semantic Versioning

All changes to this repository MUST follow the [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification. Versioning MUST strictly adhere to [Semantic Versioning 2.0.0](https://semver.org/).

### Commit Message Format
`<type>(<scope>): <description>` (Space after colon is MANDATORY)

### Types & Version Bumps
- `feat`: Minor (x.Y.z) - New feature
- `fix`: Patch (x.y.Z) - Bug fix
- `perf`: Patch (x.y.Z) - Performance improvement
- `docs`, `chore`, `style`, `refactor`, `test`: No bump (unless breaking)
- `BREAKING CHANGE:` (or `!` after type): Major (X.y.z)

## Git Attribution

**P1 rule** — Commits must appear as if authored solely by the human user. Never include:

- `Co-authored-by:` footer
- `Co-Authored-By:` footer
- Any spelling or casing variant of the above

Enforced by: pre-tool-use hook (blocks at `git commit`), `scripts/land-branch.sh` (blocks at merge), CI (verify job in `ci.yml`).

Rationale: AI agents are tools, not collaborators. Attribution footers distort contribution metrics, confuse `git blame`, and violate the single-author model of solo-local workflow.

## GitHub & Git Operations

- No direct work on `main` or `master`. Every task MUST start with a feature branch or worktree via `kickoff-branch`.
- **Integrate (solo profile):** Ship with `bash scripts/land-branch.sh <branch> "<conventional message>"` after `release-branch` gates — local squash to `main`, then push. PR is optional (remote CI / branch protection only).
- `git push origin <feature-branch>` is allowed for backup or CI; never push directly to `main`/`master` except via `land-branch.sh`.
- Use `gh repo clone` not `git clone` for GitHub repos
- Use `gh run view` / `gh run watch` for CI status
- Verify auth with `gh auth status` before operations
- **Git Attribution:** NEVER include `Co-authored-by`, `Co-Authored-By`, or any other footer that attributes code to an AI agent. All commits must appear as if they were authored solely by the human user.
- Never call GitHub REST API directly (curl, fetch, etc.)

## Agent Workflow Mandates

**AGENTS MUST NEVER BYPASS THE BIGPOWERS WORKFLOW.**

- **No Direct Coding:** When a user issues a directive like "build feature X", you MUST NOT execute the request by writing code directly.
- **Required Skills:** You MUST route all work through the appropriate bigpowers skills.
  - Start with `survey-context` if you lack context.
  - Use `plan-work` to flesh out tasks in `specs/epics/eNN-*.yaml` before writing any feature code.
  - Use `develop-tdd` or `execute-plan` to implement the plan.
  - Use `investigate-bug` for bug reports before writing a fix.
- **Verification Mandate:** Every story implementation MUST end with a step-by-step manual verification script provided to the user.
- **Traceability Mandate:** Every story MUST have at least one `story: eNNsNN` tag in its implementing code or test file.
- **Stream Continuity:** When writing large files or long documents, you MUST output continuously in chunks of ~200 lines.

## Always Green / Shift Left

**Always Green** means Preflight and CI are green before any forward work.

**Preflight** — the project's full local verification stack (lint, test, build). Preflight MUST pass before kickoff, develop, or verify phases advance.

**CI green** — when a PR exists or remote CI applies, `gh pr checks` MUST show passing before merge or land.

## Discovered Defects

Any **reproducible gate failure** encountered during unrelated work is a discovered defect.

**fix-or-log ladder (mandatory):**

1. **quick-fix** — trivial, data-only, or single-file fixes within guardrails.
2. **fix-bug** — when quick-fix guardrails abort, or the failure needs investigation.
3. **Log** — only when reproduction is blocked after good-faith attempt.

**Hard block:** Red Preflight or red CI blocks kickoff-branch, develop-tdd, and verify-work forward progress until fix-or-log produces green.

## Never

- Never dismiss reproducible gate failures as pre-existing or out of scope
- Never proceed on red Preflight or red CI — invoke quick-fix or fix-bug first
- Never edit another repo's files from this repo — use `gh repo clone` and work there
