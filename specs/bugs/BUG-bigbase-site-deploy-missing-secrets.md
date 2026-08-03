---
bug_id: BUG-bigbase-site-deploy-missing-secrets
status: open
severity: high
scope: ci
title: Site Deploy action failed due to missing BIGBASE_DEPLOY_TOKEN secret
---

## Summary

The `Site Deploy` workflow (run [30769185930](https://github.com/danielvm-git/big-release/actions/runs/30769185930)) failed during the `Deploy to BigBase` step with the error:
`deploy_token not set. Provision a scoped token via provision_ci_credentials — do not use account email/password.`

## Root Cause

1. The repository `danielvm-git/big-release` had not been registered or provisioned as a site on BigBase (`https://release.bigbase.click`).
2. The GitHub repository secrets `BIGBASE_DEPLOY_TOKEN` and `BIGBASE_SITE_ID` were not populated in the `danielvm-git/big-release` repository.

## Fix Approach

Following the instructions from [BigBase MCP Discussion #192](https://github.com/danielvm-git/bigbase/discussions/192):
1. Authenticate with BigBase API and obtain an MCP API key with `mcp:provision` scope.
2. Call BigBase MCP `create_repo` to register `danielvm-git/big-release` (`repo_id: 91d1271af913386a8135aa3c0141816e`).
3. Call BigBase MCP `create_site` to provision the `release` site (`site_id: ef9fb02a15c299f744ca672e69333827`).
4. Call BigBase MCP `provision_ci_credentials` to generate a site-scoped deploy token (`bb_dep_*`).
5. Configure `BIGBASE_DEPLOY_TOKEN` and `BIGBASE_SITE_ID` as GitHub repository secrets using `gh secret set`.
6. Re-run the failed GitHub Action workflow or trigger a new run to verify deployment success.

## Verification

1. `gh secret list --repo danielvm-git/big-release` confirms `BIGBASE_DEPLOY_TOKEN` and `BIGBASE_SITE_ID` are present.
2. Trigger / re-run failed workflow `30769185930` using `gh run rerun 30769185930 --repo danielvm-git/big-release`.
3. Confirm deployment step passes and site `https://release.bigbase.click` passes health check.
