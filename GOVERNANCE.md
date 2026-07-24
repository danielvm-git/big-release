# Governance

## Project status

**big-release** is an open-source release-automation tool maintained by
`danielvm-git`. It determines a project's next version from its commit history,
generates a changelog, and publishes packages across languages from a single
binary.

## Decision-making structure

| Role | Who | Decides |
| --- | --- | --- |
| **Maintainer** | `@danielvm-git` (repository owner) | Final say on architecture, the release algorithm's behavior, and what gets merged. |
| **Contributors** | Anyone who opens a merged PR | Propose changes via issues and pull requests. |
| **Community** | Anyone using big-release | Raises questions in Discussions; influences direction through feedback and bug reports. |

## How decisions are made

1. **Small fixes** (typos, docs, a bug fix with an obvious cause) — a
   contributor opens a PR; the maintainer reviews and merges. No ceremony.

2. **New publishers or configuration surface** — open an issue first. The change
   must fit the plugin/publisher model documented under
   [`docs/how-to/publishers/`](./docs/how-to/publishers) and
   [`docs/explanation/algorithm.md`](./docs/explanation/algorithm.md). Cite the
   registry's publishing contract in the PR.

3. **Changes to the release algorithm** (version-determination rules, branch
   matching, changelog generation, SemVer handling) — these are the
   highest-stakes changes because 35+ consuming repos depend on consistent
   behavior. They require:
   - An issue describing the gap or inconsistency.
   - A proposed change that preserves compatibility with the
     [Conventional Commits](https://www.conventionalcommits.org) and
     [SemVer](https://semver.org) contracts big-release implements, or a clear
     migration story if it breaks them.
   - Tests updated in the same PR; the algorithm is ported from
     semantic-release, so divergence must be deliberate and documented.

4. **Changes to repository conventions** (`CONVENTIONS.md`, the git attribution
   rule, the branch-protection model) — the highest bar. These exist to keep
   history clean and CI honest. Removing or relaxing one requires evidence that
   the problem it prevents no longer applies.

## The team

Maintainers and regular contributors are listed in
[`CODEOWNERS`](./CODEOWNERS). Code ownership guides which review is requested
for a given path, not who is permitted to contribute.

## Amendments to this document

Propose changes by opening a PR that updates this file. Explain the governance
gap the change addresses.
