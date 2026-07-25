# Branch protection — big-release

## Required status checks

`main` requires these jobs from **Test Build Release** to pass:

| Job | Purpose |
|-----|---------|
| `lint` | golangci-lint (pinned v2.5.0) |
| `test` | vet, unit tests, conventional commits, preflight |
| `build` | host binary compile + artifact upload |

Configured via GitHub API (2026-07-25):

```json
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["lint", "test", "build"]
  },
  "enforce_admins": false
}
```

Verify live settings: `gh api repos/danielvm-git/big-release/branches/main/protection`

After the first TBR run, confirm check names match: `gh api repos/danielvm-git/big-release/commits/main/check-runs --jq '.check_runs[].name'`

## Solo owner direct-push policy

- **`enforce_admins: false`** — owner may push directly to `main`; required checks still run on each push.
- **No required PR reviews** — solo maintainer workflow.
- **Prefer PRs** for feature work; direct push reserved for hotfixes and release-bot commits.

## Deploy checklist

**N/A** — CLI publishes via the `release` job. No `deploy.yml`, BigBase host, or post-deploy smoke.
