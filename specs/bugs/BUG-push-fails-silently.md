# BUG-2026-07-25T060000: Git push fails silently in CI with exit status 1

## Problem

big-release's release job in CI fails with:
```
plugin "git" publish failed: push failed, local tag 2.3.0 removed: failed to push: exit status 1
```

The error message shows "exit status 1" but no git stderr output, making it impossible to diagnose the actual push failure reason.

**Expected behavior:** The push should succeed (creating the tag on the remote) OR the error message should include the actual git stderr output (e.g., "rejected by branch protection", "permission denied", etc.).

**Actual behavior:** The push fails with a generic "exit status 1" and no stderr output.

## Root Cause Analysis

### Phase 1: Reproduce
The bug reproduces consistently in the CI release job when big-release tries to push the newly created tag.

### Phase 2: Isolate
The issue is in `internal/git/client.go`'s `Push` function:
```go
func (c *Client) Push(remote string) error {
    cmd := gitCmd("push", remote)
    if err := runGit(cmd); err != nil {
        return fmt.Errorf("failed to push: %w", err)
    }
    return nil
}
```

The `runGit` function captures stderr, but the error message in CI shows no stderr output. This could mean:
1. git is not producing stderr output for this specific failure
2. The stderr is being swallowed somewhere
3. The push is failing for a reason that doesn't produce stderr

### Phase 3: Hypothesize
The most likely hypothesis is that the push is failing because:
1. **Branch protection**: The main branch is protected and the GITHUB_TOKEN doesn't have permission to push
2. **No new commits**: The push is trying to push commits that don't exist (the checkout action fetches code but doesn't create new commits)
3. **Credential issue**: The git credentials configured by the checkout action are not being used correctly

### Phase 4: Verify

Looking at the CI workflow:
- The release job has `contents: write` permission ✅
- The checkout action uses `persist-credentials: true` ✅
- The push is called after creating a tag locally

**Root cause identified:** The `pushRefs()` function pushes commits AND tags sequentially:

```go
func (p *GitPlugin) pushRefs() error {
    if err := p.Git.Push("origin"); err != nil {
        return err  // ← If this fails, tags are never pushed
    }
    return p.Git.PushTags("origin")
}
```

The first `git push origin` is trying to push the current branch. In CI, the checkout action fetches the code but doesn't create new commits. If the push fails for ANY reason (branch protection, no new commits, credential issue), the tag push is never attempted.

**The real issue:** Tags should be pushed independently of commits. The tag is the important artifact, not the commit push.

**Security impact: NONE** — This is a CI workflow issue, not a security vulnerability.

## TDD Fix Plan

### 1. **RED**: Write a test that verifies tags are pushed even if commit push fails
   **GREEN**: Modify `pushRefs()` to push tags independently of commits
   **verify**: `go test ./internal/plugins/ -run TestGitPlugin_PushTagsIndependent`

### 2. **RED**: Write a test that verifies Push includes stderr in error message
   **GREEN**: Already implemented in `runGit()` — verify it works
   **verify**: `go test ./internal/git/ -run TestPush`

### 3. **RED**: Write a test that verifies pushRefs handles "Everything up-to-date" gracefully
   **GREEN**: Modify `pushRefs()` to not fail when there are no new commits to push
   **verify**: `go test ./internal/plugins/ -run TestGitPlugin_PushRefs`

### 4. **RED**: Write a test that verifies the actual git error is surfaced in CI
   **GREEN**: Add verbose logging to the Push function to log the actual git command and output
   **verify**: `go test ./internal/git/ -run TestPush_VerboseLogging`

## Acceptance Criteria

- [ ] Tags are pushed to remote even if commit push has no changes
- [ ] Push failures include the actual git stderr output in the error message
- [ ] pushRefs handles "Everything up-to-date" gracefully (no failure when no new commits)
- [ ] CI logs show the actual git command and output when push fails
- [ ] All existing tests still pass

## Resolution

<!-- filled in by validate-fix -->
