// story: e23s01 e23s02 e23s03
package plugins

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func TestGitLabPluginName(t *testing.T) {
	t.Run("SC-e23s01-P1-01: Name returns 'gitlab'", func(t *testing.T) {
		p := NewGitLabPlugin()
		if name := p.Name(); name != "gitlab" {
			t.Errorf("expected Name() == %q, got %q", "gitlab", name)
		}
	})
}

func TestGitLabPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e23s01-P1-02: VerifyConditions passes with GITLAB_TOKEN and CI_PROJECT_ID", func(t *testing.T) {
		setEnv(t, "GITLAB_TOKEN", "valid-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("SC-e23s01-P1-03: VerifyConditions passes with CI_JOB_TOKEN in CI", func(t *testing.T) {
		unsetEnv(t, "GITLAB_TOKEN")
		setEnv(t, "CI_JOB_TOKEN", "job-token")
		setEnv(t, "CI_PROJECT_ID", "99")
		defer unsetEnv(t, "CI_JOB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Errorf("expected no error with CI_JOB_TOKEN, got: %v", err)
		}
	})

	t.Run("SC-e23s01-P1-04: VerifyConditions skips credentials in dry-run", func(t *testing.T) {
		unsetEnv(t, "GITLAB_TOKEN")
		unsetEnv(t, "CI_JOB_TOKEN")
		unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		ctx := &algorithm.ReadOnlyContext{DryRun: true}
		if err := p.VerifyConditions(ctx, &algorithm.MutableState{}); err != nil {
			t.Errorf("dry-run must skip credential checks, got: %v", err)
		}
	})

	t.Run("SC-e23s01-P1-05: VerifyConditions fails with missing token", func(t *testing.T) {
		unsetEnv(t, "GITLAB_TOKEN")
		unsetEnv(t, "CI_JOB_TOKEN")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err == nil {
			t.Error("expected error with missing token, got nil")
		}
	})

	t.Run("SC-e23s01-P1-06: VerifyConditions fails with missing project ID", func(t *testing.T) {
		setEnv(t, "GITLAB_TOKEN", "valid-token")
		unsetEnv(t, "CI_PROJECT_ID")
		unsetEnv(t, "GITLAB_PROJECT_ID")
		defer unsetEnv(t, "GITLAB_TOKEN")

		p := NewGitLabPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err == nil {
			t.Error("expected error with missing project ID, got nil")
		}
	})
}

func TestGitLabPluginPublish(t *testing.T) {
	t.Run("SC-e23s02-P1-01: Publish returns nil in dry-run mode", func(t *testing.T) {
		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		ctx := &algorithm.ReadOnlyContext{DryRun: true}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}
		release, err := p.Publish(ctx, state)
		if err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
		if release != nil {
			t.Errorf("expected nil release in dry-run, got %v", release)
		}
	})

	t.Run("SC-e23s02-P1-02: Publish succeeds on HTTP 201", func(t *testing.T) {
		var authHeader string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if !strings.HasSuffix(r.URL.Path, "/releases") {
				t.Errorf("expected releases path, got %s", r.URL.Path)
			}
			authHeader = r.Header.Get("Private-Token")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{
				Version: "1.0.0",
				Type:    algorithm.ReleaseTypePatch,
				Notes:   "Release notes",
			},
		}

		if _, err := p.Publish(ctx, state); err != nil {
			t.Errorf("expected no error on 201, got: %v", err)
		}
		if authHeader != "test-token" {
			t.Errorf("expected Private-Token auth, got %q", authHeader)
		}
	})

	t.Run("SC-e23s02-P1-03: Publish uses JOB-TOKEN when CI_JOB_TOKEN set", func(t *testing.T) {
		var jobToken string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jobToken = r.Header.Get("Job-Token")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
		}))
		defer server.Close()

		unsetEnv(t, "GITLAB_TOKEN")
		setEnv(t, "CI_JOB_TOKEN", "job-token-value")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "CI_JOB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
		}

		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if jobToken != "job-token-value" {
			t.Errorf("expected Job-Token header, got %q", jobToken)
		}
	})

	t.Run("SC-e23s02-P1-04: Publish returns auth error on HTTP 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"401 Unauthorized"}`))
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "bad-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		_, err := p.Publish(ctx, state)
		if err == nil {
			t.Fatal("expected error on 401, got nil")
		}
		if !strings.Contains(err.Error(), "authentication") {
			t.Errorf("expected auth error, got: %v", err)
		}
	})

	t.Run("SC-e23s02-P1-05: Publish returns error on HTTP 409 duplicate release", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"message":"Already exists"}`))
		}))
		defer server.Close()

		setEnv(t, "GITLAB_TOKEN", "test-token")
		setEnv(t, "CI_PROJECT_ID", "12345")
		defer unsetEnv(t, "GITLAB_TOKEN")
		defer unsetEnv(t, "CI_PROJECT_ID")

		p := NewGitLabPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL + "/api/v4"

		ctx := &algorithm.ReadOnlyContext{DryRun: false}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		_, err := p.Publish(ctx, state)
		if err == nil {
			t.Fatal("expected error on 409, got nil")
		}
	})
}

func TestGitLabPluginPublish_UploadsAssets(t *testing.T) {
	tmp := t.TempDir()
	assetPath := tmp + "/big-release-linux-amd64"
	if err := os.WriteFile(assetPath, []byte("binary-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	var releaseCalled, uploadCalled, linkCalled int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases"):
			releaseCalled++
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/uploads"):
			uploadCalled++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"url":"/uploads/abc/binary","full_path":"/group/project/uploads/abc/binary"}`))
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/assets/links"):
			linkCalled++
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "big-release-linux-amd64") {
				t.Errorf("asset link missing name, got: %s", body)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	setEnv(t, "GITLAB_TOKEN", "test-token")
	setEnv(t, "CI_PROJECT_ID", "12345")
	setEnv(t, "CI_PROJECT_URL", server.URL+"/group/project")
	defer unsetEnv(t, "GITLAB_TOKEN")
	defer unsetEnv(t, "CI_PROJECT_ID")
	defer unsetEnv(t, "CI_PROJECT_URL")

	p := NewGitLabPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL + "/api/v4"
	p.assets = []algorithm.AssetConfig{{Path: assetPath, Label: "big-release-linux-amd64"}}

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if releaseCalled != 1 {
		t.Errorf("expected 1 release POST, got %d", releaseCalled)
	}
	if uploadCalled != 1 {
		t.Errorf("expected 1 upload POST, got %d", uploadCalled)
	}
	if linkCalled != 1 {
		t.Errorf("expected 1 asset link POST, got %d", linkCalled)
	}
}

func TestGitLabPluginPublish_MissingAssetLogsWarningNotFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"tag_name":"1.0.0"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITLAB_TOKEN", "test-token")
	setEnv(t, "CI_PROJECT_ID", "12345")
	defer unsetEnv(t, "GITLAB_TOKEN")
	defer unsetEnv(t, "CI_PROJECT_ID")

	p := NewGitLabPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL + "/api/v4"
	p.assets = []algorithm.AssetConfig{{Path: "/nonexistent/missing-binary"}}

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Errorf("missing asset must not fail the release, got: %v", err)
	}
}

func TestGitLabPluginPublish_CustomReleaseNameTemplate(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			capturedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"tag_name":"2.0.0"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITLAB_TOKEN", "test-token")
	setEnv(t, "CI_PROJECT_ID", "12345")
	defer unsetEnv(t, "GITLAB_TOKEN")
	defer unsetEnv(t, "CI_PROJECT_ID")

	p := NewGitLabPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL + "/api/v4"
	p.releaseNameTemplate = "Release {{.Version}}"

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "2.0.0", Type: algorithm.ReleaseTypeMinor, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(string(capturedBody), `"name":"Release 2.0.0"`) {
		t.Errorf("expected templated release name, got: %s", capturedBody)
	}
}

func TestGitLabPluginAutoRegistration(t *testing.T) {
	t.Run("SC-e23s03-P1-01: GitLabPlugin auto-registered in global registry", func(t *testing.T) {
		found := false
		for _, name := range List() {
			if name == "gitlab" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'gitlab' to be registered in global registry")
		}
	})
}

func TestGitLabPluginConfigure(t *testing.T) {
	t.Run("SC-e23s01-P1-07: Configure loads assets from plugin config", func(t *testing.T) {
		p := NewGitLabPlugin()
		raw := map[string]interface{}{
			"assets": []interface{}{
				map[string]interface{}{"path": "dist/*.tar.gz", "label": "archives"},
			},
			"releaseName": "v{{.Version}}",
		}
		if err := p.Configure(raw); err != nil {
			t.Fatalf("Configure failed: %v", err)
		}
		if len(p.assets) != 1 || p.assets[0].Path != "dist/*.tar.gz" {
			t.Errorf("expected asset config, got %+v", p.assets)
		}
		if p.releaseNameTemplate != "v{{.Version}}" {
			t.Errorf("expected releaseName template, got %q", p.releaseNameTemplate)
		}
	})
}
