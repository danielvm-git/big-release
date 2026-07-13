package packagist_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/packagist"
)

func TestPackagistName(t *testing.T) {
	t.Run("SC-e02s04-P3-01: Name returns 'packagist'", func(t *testing.T) {
		p := packagist.NewPublisher()
		if name := p.Name(); name != "packagist" {
			t.Errorf("expected Name() == %q, got %q", "packagist", name)
		}
	})
}

func TestPackagistDetect(t *testing.T) {
	t.Run("SC-e02s04-P3-02: Detect true with composer.json", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "composer.json", composerJSONContent("test-vendor/test-pkg", "0.1.0"))
		p := packagist.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when composer.json exists")
			}
		})
	})

	t.Run("SC-e02s04-P3-03: Detect false without composer.json", func(t *testing.T) {
		dir := t.TempDir()
		p := packagist.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when composer.json is absent")
			}
		})
	})
}

func TestPackagistPrepare(t *testing.T) {
	t.Run("SC-e02s04-P3-04: Prepare updates version in composer.json", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "composer.json", composerJSONContent("test-vendor/test-pkg", "0.1.0"))
		p := packagist.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			data := readFile(t, filepath.Join(dir, "composer.json"))
			var cfg map[string]interface{}
			if err := json.Unmarshal(data, &cfg); err != nil {
				t.Fatalf("failed to parse updated composer.json: %v", err)
			}
			ver, ok := cfg["version"].(string)
			if !ok {
				t.Fatalf("expected version to be a string, got %T", cfg["version"])
			}
			if ver != "2.0.0" {
				t.Errorf("expected version 2.0.0, got %q", ver)
			}
		})
	})

	t.Run("SC-e02s04-P3-05: Prepare returns error when composer.json missing", func(t *testing.T) {
		dir := t.TempDir()
		p := packagist.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error when composer.json is missing, got nil")
			}
			if !strings.Contains(err.Error(), "composer.json") {
				t.Errorf("expected error to mention composer.json, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s04-P3-06: Prepare returns error on malformed JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "composer.json", `{invalid json content`)
		p := packagist.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error for malformed JSON, got nil")
			}
		})
	})
}

func TestPackagistEnvValidation(t *testing.T) {
	t.Run("SC-e02s04-P3-12: Publish returns error when PACKAGIST_TOKEN is empty", func(t *testing.T) {
		_ = os.Unsetenv("PACKAGIST_TOKEN")
		p := packagist.NewPublisher()
		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error when PACKAGIST_TOKEN is empty, got nil")
		}
		if !strings.Contains(err.Error(), "PACKAGIST_TOKEN") {
			t.Errorf("expected error message to mention PACKAGIST_TOKEN, got %q", err.Error())
		}
	})
}

func TestPackagistAutoRegistration(t *testing.T) {
	t.Run("SC-e02s04-P3-15: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("packagist")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "packagist" {
			t.Errorf("expected name %q, got %q", "packagist", got.Name())
		}
	})
}

// --- helpers ---

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return data
}

func withDir(t *testing.T, dir string, fn func()) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore directory to %s: %v", orig, err)
		}
	}()
	fn()
}

func composerJSONContent(name, version string) string {
	return `{
    "name": "` + name + `",
    "description": "A test package",
    "type": "library",
    "version": "` + version + `",
    "require": {
        "php": ">=8.0"
    }
}
`
}

func setenv(t *testing.T, key, value string) {
	t.Helper()
	orig := os.Getenv(key)
	_ = os.Setenv(key, value)
	t.Cleanup(func() {
		if orig == "" {
			_ = os.Unsetenv(key)
		} else {
			_ = os.Setenv(key, orig)
		}
	})
}

// --- Integration tests ---

func TestPackagistPublishHTTP(t *testing.T) {
	t.Run("SC-e02s04-P3-06: Publish returns nil on HTTP 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if auth := r.Header.Get("Authorization"); auth != "token test-token-123" {
				t.Errorf("expected Authorization header, got %q", auth)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	})

	t.Run("SC-e02s04-P3-07: Publish returns auth error on HTTP 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "401") {
				t.Errorf("expected error to mention HTTP 401, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s04-P3-08: Publish returns server error on HTTP 5xx", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "500") {
				t.Errorf("expected error to mention HTTP 500, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s04-P3-09a: Publish retries with backoff on 429, succeeds", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected eventual success, got %v", err)
			}
			if attempts != 3 {
				t.Errorf("expected 3 attempts, got %d", attempts)
			}
		})
	})

	t.Run("SC-e02s04-P3-09b: Publish returns error after exhausting retries on 429", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error after retry exhaustion, got nil")
			}
			if !strings.Contains(err.Error(), "retries") {
				t.Errorf("expected error to mention retries, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s04-P3-10: Publish in dry-run makes zero HTTP requests", func(t *testing.T) {
		requestCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/update-package", func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PACKAGIST_TOKEN", "test-token-123")
		p := &packagist.Publisher{}
		p.APIURL = srv.URL
		p.DryRun = true

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil in dry-run, got %v", err)
			}
			if requestCount != 0 {
				t.Errorf("expected 0 HTTP requests in dry-run, got %d", requestCount)
			}
		})
	})

	t.Run("SC-e02s04-P3-12: Verify returns nil on version match", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/packages/test-vendor/test-pkg.json", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]interface{}{
				"package": map[string]interface{}{
					"versions": map[string]interface{}{
						"2.0.0": map[string]interface{}{},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "composer.json", composerJSONContent("test-vendor/test-pkg", "0.1.0"))
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err != nil {
				t.Fatalf("expected nil on version match, got %v", err)
			}
		})
	})

	t.Run("SC-e02s04-P3-13: Verify returns error on HTTP 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/packages/test-vendor/test-pkg.json", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "composer.json", composerJSONContent("test-vendor/test-pkg", "0.1.0"))
		p := packagist.NewPublisher()
		p.APIURL = srv.URL

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err == nil {
				t.Fatal("expected error on 404, got nil")
			}
			if !strings.Contains(err.Error(), "404") {
				t.Errorf("expected error to mention HTTP 404, got %q", err.Error())
			}
		})
	})
}
