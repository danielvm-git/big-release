// story: BUG-push-fails-silently
package plugins

import (
	"errors"
	"strings"
	"testing"
)

// errPolicyRejection is git's own output when any server refuses a ref. The same
// markers appear on GitHub, GitLab, Gitea and Bitbucket, which is why the
// plugin must not special-case a forge here.
var errPolicyRejection = errors.New(`remote: GitLab: You are not allowed to push code to protected branches on this project.
 ! [remote rejected] main -> main (pre-receive hook declined)
error: failed to push some refs to 'origin'`)

// THE regression. 4e77fed made pushRefs `return nil` when the commit push
// failed, so a release reported success while the release commit never landed
// — re-creating the silent failure the bug was named for, and diverging from
// semantic-release, where a failed push fails the release.
func TestPushRefs_PolicyRejectionOnCommitPushIsFatal(t *testing.T) {
	f := &fakeGit{pushCommitErr: errPolicyRejection}
	p := NewGitPlugin(f)

	err := p.pushRefs()

	if err == nil {
		t.Fatal("commit push rejected by remote policy must fail the release, got nil")
	}
	if !f.pushTagsCalled {
		t.Error("tags should still be pushed first")
	}
	if !strings.Contains(err.Error(), "protected branches") {
		t.Errorf("error dropped git's stderr: %q", err.Error())
	}
}

// The one thing that must still be swallowed is nothing at all: git exits 0
// with "Everything up-to-date" when there is nothing to push, so a successful
// push must not be turned into an error.
func TestPushRefs_SucceedsWhenPushSucceeds(t *testing.T) {
	f := &fakeGit{}
	p := NewGitPlugin(f)

	if err := p.pushRefs(); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !f.pushCalled || !f.pushTagsCalled {
		t.Errorf("both tags and commits should be pushed: tags=%v commits=%v", f.pushTagsCalled, f.pushCalled)
	}
}

// tagOnly is the declared opt-out for remotes the release identity cannot push
// commits to. It must skip the commit push entirely rather than push-and-ignore.
func TestPushRefs_TagOnlySkipsCommitPush(t *testing.T) {
	f := &fakeGit{pushCommitErr: errPolicyRejection}
	p := NewGitPlugin(f)
	if err := p.Configure(map[string]interface{}{"tagOnly": true}); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	if err := p.pushRefs(); err != nil {
		t.Fatalf("tagOnly should not fail on an unpushable commit, got %v", err)
	}
	if !f.pushTagsCalled {
		t.Error("tagOnly must still push tags")
	}
	if f.pushCalled {
		t.Error("tagOnly must not attempt the commit push at all")
	}
}

// Default must remain semantic-release parity: commits are pushed unless the
// user opts out. Guards against tagOnly silently becoming the default again.
func TestPushRefs_DefaultPushesCommits(t *testing.T) {
	f := &fakeGit{}
	p := NewGitPlugin(f)
	if err := p.Configure(map[string]interface{}{"message": "chore: x"}); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	if err := p.pushRefs(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.pushCalled {
		t.Fatal("commit push must happen by default (semantic-release parity)")
	}
}

// A failed tag push was already fatal; keep it that way.
func TestPushRefs_TagPushFailureIsFatal(t *testing.T) {
	f := &fakeGit{pushTagsErr: errors.New("boom")}
	p := NewGitPlugin(f)

	if err := p.pushRefs(); err == nil {
		t.Fatal("tag push failure must fail the release")
	}
	if f.pushCalled {
		t.Error("commit push should not be attempted after a failed tag push")
	}
}
