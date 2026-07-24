<!--
  Thanks for contributing to big-release! Fill in the sections below.
  Changes to the release algorithm or a publisher's publishing contract
  require a linked issue.
-->

## Summary

<!-- One or two sentences: what does this PR change, and why? -->

## Linked issue

<!-- "Closes #N" or "Refs #N". Algorithm/publisher changes require a linked issue. -->

## What kind of change is this?

<!-- Check one. -->

- [ ] Bug fix — fixes incorrect release behavior
- [ ] New feature — adds a publisher, plugin, or config option
- [ ] Breaking change — changes SemVer/Conventional-Commits handling or the CLI surface
- [ ] Refactor — no behavior change
- [ ] Documentation
- [ ] Tooling / CI
- [ ] Other

## Testing done

<!-- How did you verify this? Which publisher(s) or path does it affect? -->

## Checklist

- [ ] I have read [`CONTRIBUTING.md`](../CONTRIBUTING.md) and [`CONVENTIONS.md`](../CONVENTIONS.md).
- [ ] `make preflight` passes (lint + vet + test).
- [ ] New behavior has tests; existing tests still pass.
- [ ] Documentation updated if the public CLI, config, or a publisher contract changed.
- [ ] Commit message(s) follow Conventional Commits (`feat:`, `fix:`, etc.).
- [ ] No `Co-authored-by:` / `Co-Authored-By:` footer in any commit (enforced by hook and CI).

## Notes for review

<!-- Anything reviewers should pay attention to, or migration notes for consumers. -->
