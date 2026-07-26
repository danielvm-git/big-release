// story: BUG-push-fails-silently
package git

import (
	"fmt"
	"strings"
)

// PushOutcome classifies why a `git push` failed.
//
// Classification keys on git's own client-side output, never on a forge's
// wording. `! [remote rejected] ... (pre-receive hook declined)` is emitted by
// git itself when any server refuses a ref, so the same markers work on
// GitHub, GitLab, Gitea, Bitbucket and self-hosted remotes. Forge-specific
// strings arrive on `remote:` lines and are used only to enrich Hint().
type PushOutcome string

const (
	PushOutcomeUnknown         PushOutcome = "unknown"
	PushRejectedPolicy         PushOutcome = "rejected-policy"
	PushRejectedNonFastForward PushOutcome = "rejected-non-fast-forward"
	PushNoUpstream             PushOutcome = "no-upstream"
	PushAuthFailed             PushOutcome = "auth-failed"
)

// PushError wraps a failed push with its classification and git's raw output.
//
// Note there is deliberately no "nothing to push" outcome: `git push` with no
// new commits exits 0 printing "Everything up-to-date", so it never reaches
// here. Every non-zero push is a real failure and must be reported as one.
type PushError struct {
	Outcome PushOutcome
	Raw     string
	Err     error
}

func (e *PushError) Error() string {
	msg := fmt.Sprintf("push failed (%s): %v", e.Outcome, e.Err)
	if hint := e.Hint(); hint != "" {
		msg += "\nhint: " + hint
	}
	return msg
}

func (e *PushError) Unwrap() error { return e.Err }

// newPushError classifies err and wraps it. Returns nil for a nil error so
// callers can use it inline.
func newPushError(err error) error {
	if err == nil {
		return nil
	}
	raw := err.Error()
	return &PushError{Outcome: classifyPush(raw), Raw: raw, Err: err}
}

// classifyPush maps git's output to an outcome using git's portable markers.
func classifyPush(output string) PushOutcome {
	low := strings.ToLower(output)

	switch {
	case containsAny(low, "authentication failed", "could not read username",
		"permission denied (publickey)", "invalid username or password"):
		return PushAuthFailed

	case strings.Contains(low, "has no upstream branch"),
		strings.Contains(low, "no upstream configured"):
		return PushNoUpstream

	// Server-side refusal. git prints `! [remote rejected]` for these; the
	// parenthetical reason is git's own, not the forge's.
	case containsAny(low, "pre-receive hook declined", "protected branch hook declined",
		"push declined", "remote rejected"):
		return PushRejectedPolicy

	// Local ref is behind. `(fetch first)` and `(non-fast-forward)` are both
	// emitted by git depending on whether the remote ref is known locally.
	case containsAny(low, "non-fast-forward", "fetch first"),
		strings.Contains(low, "! [rejected]"):
		return PushRejectedNonFastForward
	}

	return PushOutcomeUnknown
}

// forgeHints maps a substring of the remote's message to actionable advice.
// Adding a forge is a new row, not a new code path.
var forgeHints = []struct {
	match string
	hint  string
}{
	{"not allowed to push code to protected branches",
		"GitLab protected branch: grant the pushing user/token 'Allowed to push' on the branch, or use a project access token with Maintainer role."},
	{"protected branch hook declined",
		"GitHub branch protection or ruleset rejected the push: use a token permitted to bypass it, or stop pushing commits (see the git plugin's tagOnly option)."},
	{"changes must be made through a pull request",
		"GitHub ruleset requires a pull request: releases cannot push straight to this branch. Use a bypass actor, or the git plugin's tagOnly option."},
}

// Hint returns actionable advice. Forge-specific text only selects the
// wording — the outcome was already decided from git's portable markers.
func (e *PushError) Hint() string {
	low := strings.ToLower(e.Raw)

	switch e.Outcome {
	case PushRejectedPolicy:
		for _, fh := range forgeHints {
			if strings.Contains(low, fh.match) {
				return fh.hint + " Alternatively set tagOnly on the git plugin to publish tags without pushing commits."
			}
		}
		return "The remote refused the push by policy (protected branch or pre-receive hook). " +
			"Allow the release identity to push, or set tagOnly on the git plugin to publish tags without pushing commits."

	case PushRejectedNonFastForward:
		return "The remote has commits the local branch does not. Fetch and rebase before releasing; " +
			"in CI this usually means the checkout was shallow or another job pushed first."

	case PushNoUpstream:
		return "The current branch has no upstream. In CI, check out a branch rather than a detached HEAD, or configure the remote tracking ref."

	case PushAuthFailed:
		return "The remote rejected the credentials. Check the token or SSH key available to the release job and that it has write access."
	}

	return ""
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
