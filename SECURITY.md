# Security Policy

## Supported Versions

big-release is pre-1.0 software. Only the latest released version receives
security updates.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| older   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. Thank you for improving the security
of this project.

**Please do not report security vulnerabilities through public GitHub issues.**

big-release consumes registry credentials (`NPM_TOKEN`, `PYPI_TOKEN`,
`CARGO_TOKEN`, etc.) during release and runs in CI environments, so reports
involving secret handling, command injection, or arbitrary code execution during
publishing are especially welcome.

Report vulnerabilities privately using one of these methods:

- **GitHub private vulnerability reporting** (preferred) — use the
  *"Report a vulnerability"* button under the **Security** tab of this
  repository:
  https://github.com/danielvm-git/big-release/security/advisories/new
- **Email** — send details to **INSERT SECURITY EMAIL**.

Please include:

1. A description of the vulnerability and its impact.
2. Steps to reproduce, including any proof-of-concept.
3. Affected versions (if known).
4. Any suggested mitigations.

### Response timeline

- We will acknowledge your report within **48 hours**.
- We will provide an initial assessment and a planned fix date within **7 days**.
- We will notify you when the vulnerability is fixed and credit you in the
  advisory (unless you prefer to remain anonymous).

### Scope

In scope: vulnerabilities in this project's own code, its published binaries,
and the release/publishing workflow.
Out of scope: vulnerabilities in upstream dependencies (report those upstream),
social engineering, or denial-of-service against the project's infrastructure.

## Disclosure policy

We follow coordinated disclosure: details are published only after a fix is
available and affected users have had reasonable time to upgrade.

## Related security artifacts

This file is the public vulnerability-reporting policy. The project's internal
security *review* artifacts live under `specs/security/` (e.g.
`specs/security/REVIEW.md`) and are not part of the public reporting channel.
