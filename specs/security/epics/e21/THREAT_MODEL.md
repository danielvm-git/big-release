# Threat Model — e21 Advanced GitHub

## Surface area

- GitHub REST API: release POST/PATCH (discussion_category_name, make_latest)
- Git plugin: commit message templates, asset path staging

## Risk level

**Low** — extends existing authenticated API calls and local git staging; no new secrets or trust boundaries.

## Mitigations

- Templates parsed with Go text/template (no code execution)
- Asset globs limited to modified files already in the working tree
- GitHub tokens remain env-scoped; discussion creation uses same token as release publish
