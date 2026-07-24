// story: e03s02
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
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("SC-e03s02-P1-03: VerifyConditions fails with missing GITHUB_TOKEN", func(t *testing.T) {
		unsetEnv(t, "GITHUB_TOKEN")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err == nil {
			t.Error("expected error with missing GITHUB_TOKEN, got nil")
		}
	})

	t.Run("SC-e03s02-P1-04: VerifyConditions fails with missing GITHUB_REPOSITORY", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "valid-token")
		unsetEnv(t, "GITHUB_REPOSITORY")
		defer unsetEnv(t, "GITHUB_TOKEN")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err == nil {
			t.Error("expected error with missing GITHUB_REPOSITORY, got nil")
		}
	})

	t.Run("SC-e03s02-P1-05: VerifyConditions fails with invalid GITHUB_REPOSITORY format", func(t *testing.T) {
		setEnv(t, "GITHUB_TOKEN", "valid-token")
		setEnv(t, "GITHUB_REPOSITORY", "invalid-format")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		if err := p.VerifyConditions(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err == nil {
			t.Error("expected error with invalid repo format, got nil")
		}
	})
}

func TestGitHubPluginAnalyzeCommits(t *testing.T) {
	t.Run("SC-e03s02-P1-06: AnalyzeCommits returns empty release type", func(t *testing.T) {
		p := NewGitHubPlugin()
		rt, err := p.AnalyzeCommits(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{})
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
		notes, err := p.GenerateNotes(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{})
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
		if err := p.Prepare(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
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
		ctx := &algorithm.ReadOnlyContext{
			DryRun: true,
		}
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

		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{
				Version: "1.0.0",
				Type:    algorithm.ReleaseTypePatch,
				Notes:   "Release notes",
			},
		}

		_, err := p.Publish(ctx, state)
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

		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		_, err := p.Publish(ctx, state)
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

		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
		}
		state := &algorithm.MutableState{
			NextRelease: &algorithm.Release{Version: "1.0.0"},
		}

		_, err := p.Publish(ctx, state)
		if err == nil {
			t.Fatal("expected error on 422, got nil")
		}
	})
}

func TestGitHubPluginSuccess(t *testing.T) {
	t.Run("SC-e03s02-P1-13: Success returns nil", func(t *testing.T) {
		p := NewGitHubPlugin()
		if err := p.Success(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

// --- e19s01 (#10): upload binary assets to GitHub releases ---

func TestGitHubPluginPublish_UploadsAssets(t *testing.T) {
	// The release POST must be followed by asset upload POSTs to the
	// uploads host, one per configured asset.
	tmp := t.TempDir()
	assetPath := tmp + "/big-release-linux-amd64"
	if err := os.WriteFile(assetPath, []byte("binary-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	var releaseCalled, assetCalled int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/releases") && r.Method == http.MethodPost {
			releaseCalled++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 42}`))
			return
		}
		// Asset upload endpoint: /repos/{repo}/releases/{id}/assets
		if strings.Contains(r.URL.Path, "/assets") && r.Method == http.MethodPost {
			assetCalled++
			if name := r.URL.Query().Get("name"); name == "" {
				t.Errorf("asset upload missing ?name= query")
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.uploadBaseURL = server.URL // route uploads to the same test server
	p.assets = []algorithm.AssetConfig{{Path: assetPath, Label: "big-release (linux-amd64)"}}

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
	if assetCalled != 1 {
		t.Errorf("expected 1 asset upload, got %d", assetCalled)
	}
}

func TestGitHubPluginPublish_MissingAssetLogsWarningNotFailure(t *testing.T) {
	// A configured asset that does not exist on disk must NOT fail the release.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 42}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.uploadBaseURL = server.URL
	p.assets = []algorithm.AssetConfig{{Path: "/nonexistent/missing-binary"}}

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Errorf("missing asset must not fail the release, got: %v", err)
	}
}

func TestGitHubPluginPublish_DryRunSkipsAssets(t *testing.T) {
	tmp := t.TempDir()
	assetPath := tmp + "/binary"
	if err := os.WriteFile(assetPath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	var anyCall bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anyCall = true
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.uploadBaseURL = server.URL
	p.assets = []algorithm.AssetConfig{{Path: assetPath}}

	ctx := &algorithm.ReadOnlyContext{DryRun: true}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("dry-run must not error, got: %v", err)
	}
	if anyCall {
		t.Error("dry-run must not make any HTTP calls")
	}
}

func TestGitHubPluginPublish_NoAssetsSkipsUpload(t *testing.T) {
	// Default plugin (no assets configured) behaves exactly as before:
	// single release POST, no upload calls.
	var releaseCalled, assetCalled int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			releaseCalled++
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 42}`))
			return
		}
		if strings.Contains(r.URL.Path, "/assets") {
			assetCalled++
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

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
	if assetCalled != 0 {
		t.Errorf("expected 0 asset uploads when no assets configured, got %d", assetCalled)
	}
}

// --- e19s02 (#13): draft GitHub releases ---

func TestGitHubPluginPublish_DraftReleaseCreatesAsDraftThenPublishes(t *testing.T) {
	// draftRelease:true → POST release with draft:true, then after assets
	// are uploaded, PATCH the release to draft:false.
	tmp := t.TempDir()
	assetPath := tmp + "/binary"
	if err := os.WriteFile(assetPath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	var createBody []byte
	var patchCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases"):
			createBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 99}`))
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/assets"):
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodPatch && strings.Contains(r.URL.Path, "/releases/99"):
			patchCalled = true
			patchBody, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(patchBody), `"draft":false`) {
				t.Errorf("PATCH must set draft:false, got: %s", patchBody)
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
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
	p.uploadBaseURL = server.URL
	p.draftRelease = true
	p.assets = []algorithm.AssetConfig{{Path: assetPath}}

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(string(createBody), `"draft":true`) {
		t.Errorf("create request must include draft:true, got: %s", createBody)
	}
	if !patchCalled {
		t.Error("expected a PATCH to publish the draft after asset upload")
	}
}

func TestGitHubPluginPublish_DraftReleaseNoAssets_PublishesImmediately(t *testing.T) {
	// draftRelease:true but no assets → still PATCH to draft:false right after create.
	var patchCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases"):
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 99}`))
		case r.Method == http.MethodPatch && strings.Contains(r.URL.Path, "/releases/99"):
			patchCalled = true
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
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
	p.draftRelease = true

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !patchCalled {
		t.Error("draft with no assets must still PATCH to publish")
	}
}

func TestGitHubPluginPublish_NonDraftNeverPatches(t *testing.T) {
	// Default (draftRelease:false) must NEVER send a PATCH.
	var patchCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			patchCalled = true
		}
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 99}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if patchCalled {
		t.Error("non-draft release must never send a PATCH")
	}
}

// --- e19s03 (#11): configurable release name and body templates ---

func TestGitHubPluginPublish_CustomReleaseNameTemplate(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			capturedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 1}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.releaseNameTemplate = "Release {{.Version}} - {{.Date}}"

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "2.0.0", Type: algorithm.ReleaseTypeMinor, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.HasPrefix(string(capturedBody), `{"tag_name":"2.0.0","name":"Release 2.0.0`) {
		t.Errorf("expected templated release name, got: %s", capturedBody)
	}
}

func TestGitHubPluginPublish_CustomBodyTemplate(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			capturedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 1}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.releaseBodyTemplate = "Shipped {{.Version}} on branch {{.Branch}}\n\n{{.Notes}}"

	ctx := &algorithm.ReadOnlyContext{
		DryRun: false,
		Branch: &algorithm.Branch{Name: "main"},
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.5.0", Type: algorithm.ReleaseTypeMinor, Notes: "### Added\n- new thing"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	bodyStr := string(capturedBody)
	if !strings.Contains(bodyStr, `Shipped 1.5.0 on branch main`) {
		t.Errorf("expected templated body with version+branch, got: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, `### Added`) {
		t.Errorf("expected notes interpolated into body, got: %s", bodyStr)
	}
}

func TestGitHubPluginPublish_DefaultNameTemplatePreservesVPrefix(t *testing.T) {
	// With no template configured, name must still be "v{version}".
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases") {
			capturedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 1}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "3.1.4", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	if _, err := p.Publish(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(string(capturedBody), `"name":"v3.1.4"`) {
		t.Errorf("default name must be v{version}, got: %s", capturedBody)
	}
}

func TestGitHubPluginPublish_InvalidNameTemplateReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.releaseNameTemplate = "{{ .Version" // malformed

	ctx := &algorithm.ReadOnlyContext{DryRun: false}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0", Type: algorithm.ReleaseTypePatch, Notes: "n"},
	}

	_, err := p.Publish(ctx, state)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
	if !strings.Contains(err.Error(), "template") {
		t.Errorf("expected template-related error, got: %v", err)
	}
}

// --- e19s04 (#12): comment on resolved issues/PRs after release ---

func TestParseReferencedIssues(t *testing.T) {
	cases := []struct {
		name    string
		message string
		want    []int
	}{
		{"bare ref", "fix: crash, #123", []int{123}},
		{"fixes keyword", "fix: resolve, fixes #456", []int{456}},
		{"closes keyword", "feat: add, closes #789", []int{789}},
		{"resolves keyword", "feat: add, resolves #101", []int{101}},
		{"multiple", "feat: add, closes #1 and resolves #2, fixes #3", []int{1, 2, 3}},
		{"none", "chore: update deps", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseReferencedIssues(tc.message)
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("index %d: got %d, want %d", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestGitHubPluginSuccess_CommentsOnReferencedIssues(t *testing.T) {
	var commented []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// POST /repos/{repo}/issues/{n}/comments
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/issues/") && strings.Contains(r.URL.Path, "/comments") {
			body, _ := io.ReadAll(r.Body)
			commented = append(commented, r.URL.Path)
			if !strings.Contains(string(body), "2.0.0") {
				t.Errorf("comment body should include version, got: %s", body)
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL
	p.successComment = "Released in {{.Version}}"

	ctx := &algorithm.ReadOnlyContext{
		DryRun:  false,
		Commits: []*algorithm.Commit{{Message: "fix: crash, closes #42 and fixes #99"}},
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "2.0.0", Type: algorithm.ReleaseTypeMajor},
	}

	if err := p.Success(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(commented) != 2 {
		t.Errorf("expected 2 issue comments, got %d: %v", len(commented), commented)
	}
}

func TestGitHubPluginSuccess_DryRunSkipsCommenting(t *testing.T) {
	var anyCall bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anyCall = true
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

	ctx := &algorithm.ReadOnlyContext{
		DryRun:  true,
		Commits: []*algorithm.Commit{{Message: "fix: crash, closes #42"}},
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0"},
	}

	if err := p.Success(ctx, state); err != nil {
		t.Fatalf("dry-run Success must not error, got: %v", err)
	}
	if anyCall {
		t.Error("dry-run must not make any HTTP calls")
	}
}

func TestGitHubPluginSuccess_403IsNonFatal(t *testing.T) {
	// A 403/404 on commenting must NOT fail the release.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

	ctx := &algorithm.ReadOnlyContext{
		DryRun:  false,
		Commits: []*algorithm.Commit{{Message: "fix: crash, closes #42"}},
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0"},
	}

	if err := p.Success(ctx, state); err != nil {
		t.Errorf("commenting 403 must be non-fatal, got: %v", err)
	}
}

func TestGitHubPluginSuccess_NoIssuesNoCalls(t *testing.T) {
	var anyCall bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anyCall = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	setEnv(t, "GITHUB_TOKEN", "test-token")
	setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
	defer unsetEnv(t, "GITHUB_TOKEN")
	defer unsetEnv(t, "GITHUB_REPOSITORY")

	p := NewGitHubPlugin()
	p.client = server.Client()
	p.apiBaseURL = server.URL

	ctx := &algorithm.ReadOnlyContext{
		DryRun:  false,
		Commits: []*algorithm.Commit{{Message: "chore: no issue ref here"}},
	}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{Version: "1.0.0"},
	}

	if err := p.Success(ctx, state); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if anyCall {
		t.Error("no issue refs → no HTTP calls expected")
	}
}

func TestExpandAssetGlobs(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.tar.gz", "b.tar.gz", "c.zip"} {
		if err := os.WriteFile(tmp+"/"+name, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, errs := expandAssetGlobs([]algorithm.AssetConfig{{Path: tmp + "/*.tar.gz"}})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 expanded assets, got %d: %+v", len(got), got)
	}
}

func TestMimeTypeForAsset(t *testing.T) {
	cases := map[string]string{
		"app.tar.gz": "application/gzip",
		"app.zip":    "application/zip",
		"app.exe":    "application/vnd.microsoft.portable-executable",
		"app.bin":    "application/octet-stream",
	}
	for name, want := range cases {
		if got := mimeTypeForAsset(name); got != want {
			t.Errorf("mimeTypeForAsset(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestGitHubPluginFail(t *testing.T) {
	t.Run("SC-e03s02-P1-14: Fail returns nil", func(t *testing.T) {
		p := NewGitHubPlugin()
		if err := p.Fail(&algorithm.ReadOnlyContext{}, &algorithm.MutableState{}, nil); err != nil {
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

func TestGitHubPluginPublish_DiscussionCategory(t *testing.T) {
	t.Run("SC-e21s01-P1-01: Publish sends discussion_category_name when configured", func(t *testing.T) {
		var body string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				b, _ := io.ReadAll(r.Body)
				body = string(b)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 99}`))
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL
		if err := p.Configure(map[string]interface{}{
			"discussionCategoryName": "Announcements",
		}); err != nil {
			t.Fatalf("Configure: %v", err)
		}

		ctx := &algorithm.ReadOnlyContext{DryRun: false, Branch: &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease}}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "2.0.0", Notes: "notes"}}
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if !strings.Contains(body, `"discussion_category_name":"Announcements"`) {
			t.Errorf("expected discussion_category_name in payload, got %s", body)
		}
	})
}

func TestGitHubPlugin_resolveMakeLatest(t *testing.T) {
	t.Run("SC-e21s02-P1-01: main release branch defaults make_latest true", func(t *testing.T) {
		p := NewGitHubPlugin()
		ctx := &algorithm.ReadOnlyContext{Branch: &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease}}
		if got := p.resolveMakeLatest(ctx); got != "true" {
			t.Errorf("got %q, want true", got)
		}
	})

	t.Run("SC-e21s02-P1-02: prerelease branch defaults make_latest false", func(t *testing.T) {
		p := NewGitHubPlugin()
		ctx := &algorithm.ReadOnlyContext{Branch: &algorithm.Branch{Name: "beta", Type: algorithm.BranchTypePrerelease}}
		if got := p.resolveMakeLatest(ctx); got != "false" {
			t.Errorf("got %q, want false", got)
		}
	})

	t.Run("SC-e22s03-P1-01: non-latest channel defaults make_latest false", func(t *testing.T) {
		p := NewGitHubPlugin()
		ctx := &algorithm.ReadOnlyContext{Branch: &algorithm.Branch{Name: "main", Type: algorithm.BranchTypeRelease, Channel: "next"}}
		if got := p.resolveMakeLatest(ctx); got != "false" {
			t.Errorf("got %q, want false for channel next", got)
		}
	})

	t.Run("SC-e21s02-P1-03: explicit makeLatest config overrides default", func(t *testing.T) {
		p := NewGitHubPlugin()
		trueVal := true
		p.makeLatest = &trueVal
		ctx := &algorithm.ReadOnlyContext{Branch: &algorithm.Branch{Name: "beta", Type: algorithm.BranchTypePrerelease}}
		if got := p.resolveMakeLatest(ctx); got != "true" {
			t.Errorf("got %q, want true", got)
		}
	})
}

func TestGitHubPluginPublish_MakeLatest(t *testing.T) {
	t.Run("SC-e21s02-P1-04: Publish sends make_latest false for beta branch", func(t *testing.T) {
		var body string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				b, _ := io.ReadAll(r.Body)
				body = string(b)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 1}`))
		}))
		defer server.Close()

		setEnv(t, "GITHUB_TOKEN", "test-token")
		setEnv(t, "GITHUB_REPOSITORY", "owner/repo")
		defer unsetEnv(t, "GITHUB_TOKEN")
		defer unsetEnv(t, "GITHUB_REPOSITORY")

		p := NewGitHubPlugin()
		p.client = server.Client()
		p.apiBaseURL = server.URL
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
			Branch: &algorithm.Branch{Name: "beta", Type: algorithm.BranchTypePrerelease},
		}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0-beta.1", Type: algorithm.ReleaseTypePrerelease, Notes: "notes"}}
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if !strings.Contains(body, `"make_latest":"false"`) {
			t.Errorf("expected make_latest false in payload, got %s", body)
		}
	})
}

func TestGitHubPluginPublish_DraftIncludesMakeLatest(t *testing.T) {
	t.Run("SC-e21s02-P1-05: publishDraft PATCH includes make_latest", func(t *testing.T) {
		var patchBody string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/releases"):
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": 42}`))
			case r.Method == http.MethodPatch:
				b, _ := io.ReadAll(r.Body)
				patchBody = string(b)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id": 42}`))
			default:
				w.WriteHeader(http.StatusNotFound)
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
		p.draftRelease = true
		ctx := &algorithm.ReadOnlyContext{
			DryRun: false,
			Branch: &algorithm.Branch{Name: "beta", Type: algorithm.BranchTypePrerelease},
		}
		state := &algorithm.MutableState{NextRelease: &algorithm.Release{Version: "1.0.0-beta.1", Type: algorithm.ReleaseTypePrerelease, Notes: "notes"}}
		if _, err := p.Publish(ctx, state); err != nil {
			t.Fatalf("Publish: %v", err)
		}
		if !strings.Contains(patchBody, `"make_latest":"false"`) {
			t.Errorf("expected make_latest in draft publish PATCH, got %s", patchBody)
		}
	})
}
