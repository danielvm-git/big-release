// story: e08s02
package git

import (
	"fmt"
	"strings"
	"testing"
)

func TestGitAuthURL_HTTPSInject(t *testing.T) {
	const token = "ghp_injecttesttoken12345678901234567890"
	got, err := AuthURL("https://github.com/org/repo.git", token)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "x-access-token:") {
		t.Fatalf("expected x-access-token userinfo, got %q", got)
	}
	if !strings.Contains(got, token) {
		t.Fatalf("expected token in authenticated URL, got %q", got)
	}
}

func TestGitAuthURL_SSHPassthrough(t *testing.T) {
	remote := "git@github.com:org/repo.git"
	got, err := AuthURL(remote, "ghp_shouldnotappear")
	if err != nil {
		t.Fatal(err)
	}
	if got != remote {
		t.Fatalf("expected SSH passthrough %q, got %q", remote, got)
	}
}

func TestGitAuthURL_NoDoubleInject(t *testing.T) {
	existing := "https://user:pass@github.com/org/repo.git"
	got, err := AuthURL(existing, "ghp_newtoken")
	if err != nil {
		t.Fatal(err)
	}
	if got != existing {
		t.Fatalf("expected no double inject, got %q", got)
	}
}

func TestGitAuthURL_RedactedErrors(t *testing.T) {
	const token = "ghp_redactederrortoken123456789012345678"
	err := AuthURLError(fmt.Errorf("push failed for %s", token), token)
	if strings.Contains(err.Error(), token) {
		t.Fatalf("error must redact token, got %q", err.Error())
	}
}
