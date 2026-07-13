// story: e03s02
package plugins

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	_ = os.Setenv(key, value)
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	_ = os.Unsetenv(key)
}

func TestGitHubPluginName(t *testing.T) {
	t.Run("SC-e03s02-P1-01: Name returns 'github'", func(t *testing.T) {
		p := NewGitHubPlugin()
		if name := p.Name(); name != "github" {
			t.Errorf("expected Name() == %q, got %q", "github", name)
		}
	})
}

func TestGitHubPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e03s02-P1-02: VerifyConditions passes with valid env vars", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "valid-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("SC-e03s02-P1-03: VerifyConditions fails with missing GITHUB_TOKEN", func(t *testing.T) {
		unsetEnv(t, "GITHUB_TOKEN")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.Context{}); err == nil {
			t.Error("expected error with missing GITHUB_TOKEN, got nil")
		}
	})

	t.Run("SC-e03s02-P1-04: VerifyConditions fails with missing GITHUB_REPOSITORY", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "valid-token")
		unsetEnv(t, "GITHUB_REPOSITORY")
		defer unsetEnv(t, "GITHUB_TOKEN")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.Context{}); err == nil {
			t.Error("expected error with missing GITHUB_REPOSITORY, got nil")
		}
	})

	t.Run("SC-e03s02-P1-05: VerifyConditions fails with invalid GITHUB_REPOSITORY format", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "valid-token")
		setEnv(t, "GITHUB_REPOSITORY", "invalid-format")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.Context{}); err == nil {
			t.Error("expected error with invalid repo format, got nil")
		}
	})
}

func TestGitHubPluginAnalyzeCommits(t *testing.T) {
	t.Run("SC-e03s02-P1-06: AnalyzeCommits returns empty release type", func(t *testing.T) {
		p := NewGitHubPlugin()
		rt, err := p.AnalyzeCommits(&algorithm.Context{})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if rt != "" {
			t.Errorf("expected empty release type, got %q", rt)
		}
	})
}

func TestGitHubPluginGenerateNotes(t *testing.T) {
	t.Run("SC-e03s02-P1-07: GenerateNotes returns empty string", func(t *testing.T) {
		p := NewGitHubPlugin()
		notes, err := p.GenerateNotes(&algorithm.Context{})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if notes != "" {
			t.Errorf("expected empty notes, got %q", notes)
		}
	})
}

func TestGitHubPluginPrepare(t *testing.T) {
	t.Run("SC-e03s02-P1-08: Prepare returns nil", func(t *testing.T) {
		p := NewGitHubPlugin()
		if err := p.Prepare(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitHubPluginPublish(t *testing.T) {
	t.Run("SC-e03s02-P1-09: Publish returns nil in dry-run mode", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      true,
		}
		release, err := p.Publish(ctx)
		if err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
		if release != nil {
			t.Errorf("expected nil release in dry-run, got %v", release)
		}
	})

	t.Run("SC-e03s02-P1-10: Publish succeeds on HTTP 201", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("expected Bearer token auth")
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": "12345"}`))
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{
				Version: "1.0.0",
				Type:    algorithm.ReleaseTypePatch,
				Notes:   "Release notes",
			},
			DryRun: false,
		}

		_, err := p.Publish(ctx)
		if err != nil {
			t.Errorf("expected no error on 201, got: %v", err)
		}
	})

	t.Run("SC-e03s02-P1-11: Publish returns auth error on HTTP 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message": "Bad credentials"}`))
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "bad-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      false,
		}

		_, err := p.Publish(ctx)
		if err == nil {
			t.Fatal("expected error on 401, got nil")
		}
	})

	t.Run("SC-e03s02-P1-12: Publish returns error on HTTP 422 (duplicate release)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"message": "Validation Failed"}`))
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL

		ctx := &algorithm.Context{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
			DryRun:      false,
		}

		_, err := p.Publish(ctx)
		if err == nil {
			t.Fatal("expected error on 422, got nil")
		}
	})
}

func TestGitHubPluginSuccess(t *testing.T) {
	t.Run("SC-e03s02-P1-13: Success returns nil", func(t *testing.T) {
		p := NewGitHubPlugin()
		if err := p.Success(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitHubPluginFail(t *testing.T) {
	t.Run("SC-e03s02-P1-14: Fail returns nil", func(t *testing.T) {
		p := NewGitHubPlugin()
		if err := p.Fail(&algorithm.Context{}, nil); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestGitHubPluginAutoRegistration(t *testing.T) {
	t.Run("SC-e03s02-P1-15: GitHubPlugin auto-registered in global registry", func(t *testing.T) {
		found := false
		for _, name := range List() {
			if name == "github" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'github' to be registered in global registry")
		}
	})
}
