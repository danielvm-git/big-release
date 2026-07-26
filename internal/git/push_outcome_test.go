// story: BUG-push-fails-silently
package git

import (
	"errors"
	"strings"
	"testing"
)

// The fixtures below are verbatim git output captured from a real local
// bare-repo lab (pre-receive hook rejection, diverged history, missing
// upstream). They are git's own client-side rendering of the remote's
// refusal, which is why classification keys on them rather than on any
// forge's wording: GitHub, GitLab, Gitea and Bitbucket all produce them.
const (
	fixturePolicyGitHub = `remote: protected branch hook declined
To ../remote.git
 ! [remote rejected] main -> main (pre-receive hook declined)
error: failed to push some refs to '../remote.git'`

	fixturePolicyGitLab = `remote: GitLab: You are not allowed to push code to protected branches on this project.
To gitlab.com:acme/app.git
 ! [remote rejected] main -> main (pre-receive hook declined)
error: failed to push some refs to 'gitlab.com:acme/app.git'`

	// GitHub rulesets (2023+) word the refusal differently from classic branch
	// protection. Both must classify as policy.
	fixtureRulesetGitHub = `remote: error: GH013: Repository rule violations found for refs/heads/main.
remote: error: Changes must be made through a pull request.
 ! [remote rejected] main -> main (push declined due to repository rule violations)`

	fixtureNonFastForward = ` ! [rejected]        main -> main (fetch first)
error: failed to push some refs to '../remote.git'
hint: Updates were rejected because the remote contains work that you do not`

	fixtureNoUpstream = `fatal: The current branch main has no upstream branch.
To push the current branch and set the remote as upstream, use`

	fixtureAuth = `remote: Invalid username or password.
fatal: Authentication failed for 'https://github.com/acme/app.git/'`
)

func TestClassifyPush(t *testing.T) {
	cases := []struct {
		name string
		out  string
		want PushOutcome
	}{
		{"github protected branch", fixturePolicyGitHub, PushRejectedPolicy},
		{"gitlab protected branch", fixturePolicyGitLab, PushRejectedPolicy},
		{"github ruleset", fixtureRulesetGitHub, PushRejectedPolicy},
		{"non-fast-forward", fixtureNonFastForward, PushRejectedNonFastForward},
		{"no upstream", fixtureNoUpstream, PushNoUpstream},
		{"auth failure", fixtureAuth, PushAuthFailed},
		{"unrecognised", "some other git explosion", PushOutcomeUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyPush(tc.out); got != tc.want {
				t.Fatalf("classifyPush() = %q, want %q", got, tc.want)
			}
		})
	}
}

// The whole point of the fix: a policy rejection must never be reported as
// success. This is the behaviour that regressed in 4e77fed.
func TestPushError_PolicyRejectionIsAnError(t *testing.T) {
	err := newPushError(errors.New(fixturePolicyGitHub))
	if err == nil {
		t.Fatal("policy rejection must produce an error, got nil")
	}
	var pe *PushError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PushError, got %T", err)
	}
	if pe.Outcome != PushRejectedPolicy {
		t.Fatalf("Outcome = %q, want %q", pe.Outcome, PushRejectedPolicy)
	}
}

// Forge-specific text must reach the operator as an actionable hint, but only
// as enrichment — the decision above is made from git's portable markers.
func TestPushError_HintIsForgeAware(t *testing.T) {
	cases := []struct {
		name    string
		out     string
		wantSub string
	}{
		{"github hint", fixturePolicyGitHub, "branch protection"},
		{"github ruleset hint", fixtureRulesetGitHub, "pull request"},
		{"gitlab hint", fixturePolicyGitLab, "protected branch"},
		{"non-fast-forward hint", fixtureNonFastForward, "fetch"},
		{"auth hint", fixtureAuth, "credential"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pe := &PushError{Outcome: classifyPush(tc.out), Raw: tc.out}
			hint := pe.Hint()
			if hint == "" {
				t.Fatal("expected a hint, got empty string")
			}
			if !strings.Contains(strings.ToLower(hint), tc.wantSub) {
				t.Fatalf("hint %q does not mention %q", hint, tc.wantSub)
			}
		})
	}
}

// A policy rejection on a protected branch is the case operators hit most, so
// the tagOnly escape hatch must be discoverable from the error itself.
func TestPushError_PolicyHintMentionsTagOnly(t *testing.T) {
	pe := &PushError{Outcome: PushRejectedPolicy, Raw: fixturePolicyGitHub}
	if !strings.Contains(pe.Hint(), "tagOnly") {
		t.Fatalf("policy hint should point at the tagOnly option, got %q", pe.Hint())
	}
}

// Error() must keep git's stderr — the guarantee BUG-git-push-error-swallowed
// established in the layer below. Classification must not replace it.
func TestPushError_PreservesGitStderr(t *testing.T) {
	err := newPushError(errors.New(fixturePolicyGitLab))
	if !strings.Contains(err.Error(), "not allowed to push code to protected branches") {
		t.Fatalf("PushError dropped git's stderr: %q", err.Error())
	}
	if strings.HasSuffix(err.Error(), "exit status 1") {
		t.Fatalf("PushError collapsed to a bare exit status: %q", err.Error())
	}
}

func TestNewPushError_NilStaysNil(t *testing.T) {
	if got := newPushError(nil); got != nil {
		t.Fatalf("newPushError(nil) = %v, want nil", got)
	}
}
