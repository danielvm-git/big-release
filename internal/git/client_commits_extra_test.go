// story: e01s04 e08s04
package git

import (
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/git/testrepo"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestGetCommits_NULParsing_PipeInSubject(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("feat: support a | b in the subject")

	commits, err := newTestClient(t).GetCommits("", "HEAD")
	if err != nil {
		t.Fatalf("GetCommits: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(commits))
	}
	c := commits[0]
	if c.Message != "feat: support a | b in the subject" {
		t.Fatalf("Message = %q", c.Message)
	}
	if c.Body != "" {
		t.Fatalf("Body = %q, want empty", c.Body)
	}
	if c.Author != "Test User" || c.Email != "test@example.com" {
		t.Fatalf("author = %q <%q>", c.Author, c.Email)
	}
	if c.Date == "" {
		t.Fatal("Date is empty")
	}
}

func TestGetCommits_ParsesBody(t *testing.T) {
	repo := testrepo.Init(t)
	repo.RunGit("commit", "--allow-empty", "-m", "feat: with body", "-m", "body line")

	commits, err := newTestClient(t).GetCommits("", "HEAD")
	if err != nil {
		t.Fatalf("GetCommits: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(commits))
	}
	if commits[0].Body != "body line" {
		t.Fatalf("Body = %q, want %q", commits[0].Body, "body line")
	}
}

func TestGetCommits_InvalidRef(t *testing.T) {
	testrepo.Init(t)
	if _, err := newTestClient(t).GetCommits("does-not-exist", "HEAD"); err == nil {
		t.Fatal("expected an error for an unknown ref")
	}
}

func TestGetTags(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	tags, err := client.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("got %v, want no tags", tags)
	}

	repo.Tag("v1.0.0")
	repo.Tag("v1.1.0")
	tags, err = client.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	seen := map[string]bool{}
	for _, tag := range tags {
		seen[tag] = true
	}
	for _, want := range []string{"v1.0.0", "v1.1.0"} {
		if !seen[want] {
			t.Fatalf("GetTags() = %v, missing %q", tags, want)
		}
	}
}

func TestGetTagHead(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	client := newTestClient(t)
	head, err := client.GetHead()
	if err != nil {
		t.Fatalf("GetHead: %v", err)
	}
	repo.Tag("v1.0.0")

	got, err := client.GetTagHead("v1.0.0")
	if err != nil {
		t.Fatalf("GetTagHead: %v", err)
	}
	if got != head {
		t.Fatalf("GetTagHead() = %q, want %q", got, head)
	}
	if _, err := client.GetTagHead("v9.9.9"); err == nil {
		t.Fatal("expected an error for an unknown tag")
	}
}

func TestGetLastRelease_SemverOrdering(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	for _, tag := range []string{"v1.0.0", "v1.2.0", "v1.10.0"} {
		repo.Tag(tag)
	}

	last, err := newTestClient(t).GetLastRelease("v${version}")
	if err != nil {
		t.Fatalf("GetLastRelease: %v", err)
	}
	if last == nil {
		t.Fatal("expected a release")
	}
	if last.GitTag != "v1.10.0" {
		t.Fatalf("GitTag = %q, want v1.10.0 (semver, not lexical)", last.GitTag)
	}
	if last.Version != "1.10.0" {
		t.Fatalf("Version = %q, want 1.10.0", last.Version)
	}
}

func TestGetLastRelease_PrefersStableOverPrerelease(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	for _, tag := range []string{"v1.0.0", "v1.1.0-beta.1", "v1.1.0"} {
		repo.Tag(tag)
	}

	last, err := newTestClient(t).GetLastRelease("v${version}")
	if err != nil {
		t.Fatalf("GetLastRelease: %v", err)
	}
	if last == nil || last.GitTag != "v1.1.0" {
		t.Fatalf("last = %+v, want v1.1.0", last)
	}
}

func TestGetLastRelease_PrereleaseOnly(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	for _, tag := range []string{"v1.1.0-beta.1", "v1.1.0-beta.2"} {
		repo.Tag(tag)
	}

	last, err := newTestClient(t).GetLastRelease("v${version}")
	if err != nil {
		t.Fatalf("GetLastRelease: %v", err)
	}
	if last == nil || last.GitTag != "v1.1.0-beta.2" {
		t.Fatalf("last = %+v, want v1.1.0-beta.2", last)
	}
}

func TestGetLastRelease_TagFormats(t *testing.T) {
	cases := []struct {
		name      string
		tags      []string
		tagFormat string
		wantTag   string
		wantVer   string
	}{
		{"empty format uses bare version", []string{"1.0.0", "2.0.0"}, "", "2.0.0", "2.0.0"},
		{"custom prefix", []string{"release-1.0.0", "release-1.1.0"},
			"release-${version}", "release-1.1.0", "1.1.0"},
		{"no version token", []string{"1.0.0"}, "release", "1.0.0", "1.0.0"},
		{"only invalid tags", []string{"nightly", "v1", "vnot.semver"}, "v${version}", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := testrepo.Init(t)
			repo.Commit("chore: baseline")
			for _, tag := range tc.tags {
				repo.Tag(tag)
			}

			last, err := newTestClient(t).GetLastRelease(tc.tagFormat)
			if err != nil {
				t.Fatalf("GetLastRelease: %v", err)
			}
			if tc.wantTag == "" {
				if last != nil {
					t.Fatalf("last = %+v, want nil", last)
				}
				return
			}
			if last == nil {
				t.Fatal("expected a release")
			}
			if last.GitTag != tc.wantTag || last.Version != tc.wantVer {
				t.Fatalf("last = %+v, want tag %q version %q", last, tc.wantTag, tc.wantVer)
			}
		})
	}
}

func TestGetLastRelease_NoTags(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")

	last, err := newTestClient(t).GetLastRelease("v${version}")
	if err != nil {
		t.Fatalf("GetLastRelease: %v", err)
	}
	if last != nil {
		t.Fatalf("last = %+v, want nil", last)
	}
}

func TestFormatTag_Table(t *testing.T) {
	cases := []struct {
		name      string
		tagFormat string
		version   string
		want      string
	}{
		{"v prefix", "v${version}", "1.2.3", "v1.2.3"},
		{"empty format", "", "1.2.3", "1.2.3"},
		{"bare token", "${version}", "1.2.3", "1.2.3"},
		{"custom prefix", "release-${version}", "1.2.3", "release-1.2.3"},
		{"suffix after token is ignored", "prefix-${version}-suffix", "1.2.3", "prefix-1.2.3"},
		{"no token", "noversion", "1.2.3", "1.2.3"},
		{"first token wins", "x${version}y${version}", "1.0.0", "x1.0.0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatTag(tc.tagFormat, tc.version); got != tc.want {
				t.Fatalf("FormatTag(%q, %q) = %q, want %q",
					tc.tagFormat, tc.version, got, tc.want)
			}
		})
	}
}

func TestFormatTag_RoundTripsThroughGetLastRelease(t *testing.T) {
	repo := testrepo.Init(t)
	repo.Commit("chore: baseline")
	repo.Tag(FormatTag("v${version}", "3.4.5"))

	last, err := newTestClient(t).GetLastRelease("v${version}")
	if err != nil {
		t.Fatalf("GetLastRelease: %v", err)
	}
	if last == nil || last.GitTag != "v3.4.5" {
		t.Fatalf("last = %+v, want v3.4.5", last)
	}
	if strings.TrimSpace(last.Version) != "3.4.5" {
		t.Fatalf("Version = %q, want 3.4.5", last.Version)
	}
}
