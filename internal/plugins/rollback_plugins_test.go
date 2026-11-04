// story: e21s03 e21s04
package plugins

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func rollbackState(version string) *algorithm.MutableState {
	return &algorithm.MutableState{NextRelease: &algorithm.Release{Version: version}}
}

func rollbackTagConfig() *algorithm.Config {
	return &algorithm.Config{TagFormat: "v${version}"}
}

func TestGitPluginRollback(t *testing.T) {
	t.Run("deletes the local and remote release tag", func(t *testing.T) {
		f := &fakeGit{tags: []string{"v1.0.0"}}
		p := NewGitPlugin(f)
		ctx := &algorithm.ReadOnlyContext{Config: rollbackTagConfig()}

		if err := p.Rollback(ctx, rollbackState("1.0.0")); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if f.lastDeletedTag != "v1.0.0" {
			t.Errorf("local tag deleted = %q, want v1.0.0", f.lastDeletedTag)
		}
		if !f.remoteDeleteCalled {
			t.Error("expected remote tag to be deleted")
		}
		if f.lastRemoteTagRemote != "origin" {
			t.Errorf("remote = %q, want origin", f.lastRemoteTagRemote)
		}
		if f.lastRemoteTagDeleted != "v1.0.0" {
			t.Errorf("remote tag deleted = %q, want v1.0.0", f.lastRemoteTagDeleted)
		}
	})

	t.Run("no-op when there is no next release", func(t *testing.T) {
		f := &fakeGit{}
		p := NewGitPlugin(f)

		if err := p.Rollback(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if f.deleteTagCalls != 0 || f.deleteRemoteTagCalls != 0 {
			t.Errorf("expected no tag deletes, got local=%d remote=%d", f.deleteTagCalls, f.deleteRemoteTagCalls)
		}
	})

	t.Run("no-op when the version is empty", func(t *testing.T) {
		f := &fakeGit{tags: []string{"v"}}
		p := NewGitPlugin(f)

		if err := p.Rollback(&algorithm.ReadOnlyContext{}, rollbackState("")); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if f.deleteTagCalls != 0 || f.deleteRemoteTagCalls != 0 {
			t.Errorf("expected no tag deletes, got local=%d remote=%d", f.deleteTagCalls, f.deleteRemoteTagCalls)
		}
	})

	t.Run("no-op when the release tag does not exist", func(t *testing.T) {
		f := &fakeGit{tags: []string{"v0.9.0"}}
		p := NewGitPlugin(f)
		ctx := &algorithm.ReadOnlyContext{Config: rollbackTagConfig()}

		if err := p.Rollback(ctx, rollbackState("1.0.0")); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if f.deleteTagCalls != 0 || f.deleteRemoteTagCalls != 0 {
			t.Errorf("expected no tag deletes, got local=%d remote=%d", f.deleteTagCalls, f.deleteRemoteTagCalls)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		f := &fakeGit{tags: []string{"v1.0.0"}}
		p := NewGitPlugin(f)
		ctx := &algorithm.ReadOnlyContext{Config: rollbackTagConfig()}

		if err := p.Rollback(ctx, rollbackState("1.0.0")); err != nil {
			t.Fatalf("first Rollback() error = %v", err)
		}
		f.tags = nil // the successful delete removed the local tag

		if err := p.Rollback(ctx, rollbackState("1.0.0")); err != nil {
			t.Fatalf("second Rollback() error = %v", err)
		}
		if f.deleteTagCalls != 1 {
			t.Errorf("local deletes = %d, want 1", f.deleteTagCalls)
		}
		if f.deleteRemoteTagCalls != 1 {
			t.Errorf("remote deletes = %d, want 1", f.deleteRemoteTagCalls)
		}
	})

	t.Run("returns the cleanup error when tag deletion fails", func(t *testing.T) {
		f := &fakeGit{tags: []string{"v1.0.0"}, deleteErr: errors.New("delete denied")}
		p := NewGitPlugin(f)
		ctx := &algorithm.ReadOnlyContext{Config: rollbackTagConfig()}

		err := p.Rollback(ctx, rollbackState("1.0.0"))
		if err == nil {
			t.Fatal("expected an error when local tag deletion fails")
		}
	})
}

func TestGitHubPluginRollback(t *testing.T) {
	t.Run("deletes the release created by Publish", func(t *testing.T) {
		var deletePath string
		var deleteCalls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": 4242}`))
			case http.MethodDelete:
				atomic.AddInt32(&deleteCalls, 1)
				deletePath = r.URL.Path
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected method %s", r.Method)
			}
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&deleteCalls); got != 1 {
			t.Errorf("delete calls = %d, want 1", got)
		}
		if want := "/repos/owner/repo/releases/4242"; deletePath != want {
			t.Errorf("delete path = %q, want %q", deletePath, want)
		}
	})

	t.Run("no-op when no release was created", func(t *testing.T) {
		var calls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&calls, 1)
		}))
		defer server.Close()

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		if err := p.Rollback(&algorithm.ReadOnlyContext{}, rollbackState("1.0.0")); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != 0 {
			t.Errorf("requests = %d, want 0", got)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		var deleteCalls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": 7}`))
			case http.MethodDelete:
				atomic.AddInt32(&deleteCalls, 1)
				w.WriteHeader(http.StatusNoContent)
			}
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("first Rollback() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("second Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&deleteCalls); got != 1 {
			t.Errorf("delete calls = %d, want 1", got)
		}
	})

	t.Run("returns the cleanup error when the delete fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": 9}`))
			case http.MethodDelete:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`boom`))
			}
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err == nil {
			t.Fatal("expected an error when the release delete fails")
		}
	})
}

func TestGitLabPluginRollback(t *testing.T) {
	t.Run("deletes the release created by Publish", func(t *testing.T) {
		var deletePath string
		var deleteCalls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
			case http.MethodDelete:
				atomic.AddInt32(&deleteCalls, 1)
				deletePath = r.URL.Path
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected method %s", r.Method)
			}
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&deleteCalls); got != 1 {
			t.Errorf("delete calls = %d, want 1", got)
		}
		if want := "/api/v4/projects/12345/releases/1.0.0"; deletePath != want {
			t.Errorf("delete path = %q, want %q", deletePath, want)
		}
	})

	t.Run("no-op when no release was created", func(t *testing.T) {
		var calls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&calls, 1)
		}))
		defer server.Close()

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		if err := p.Rollback(&algorithm.ReadOnlyContext{}, rollbackState("1.0.0")); err != nil {
			t.Fatalf("Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != 0 {
			t.Errorf("requests = %d, want 0", got)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		var deleteCalls int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
			case http.MethodDelete:
				atomic.AddInt32(&deleteCalls, 1)
				w.WriteHeader(http.StatusNoContent)
			}
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("first Rollback() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err != nil {
			t.Fatalf("second Rollback() error = %v", err)
		}
		if got := atomic.LoadInt32(&deleteCalls); got != 1 {
			t.Errorf("delete calls = %d, want 1", got)
		}
	})

	t.Run("returns the cleanup error when the delete fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
			case http.MethodDelete:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`boom`))
			}
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{}
		state := rollbackState("1.0.0")
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if err := p.Rollback(ctx, state); err == nil {
			t.Fatal("expected an error when the release delete fails")
		}
	})
}
