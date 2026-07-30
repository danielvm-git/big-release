package crates_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/crates"
)

// --- Unit tests ---

func TestCratesName(t *testing.T) {
	t.Run("SC-e02s02-P1-01: Name returns 'crates'", func(t *testing.T) {
		p := crates.NewPublisher()
		if name := p.Name(); name != "crates" {
			t.Errorf("expected Name() == %q, got %q", "crates", name)
		}
	})
}

func TestCratesDetect(t *testing.T) {
	t.Run("SC-e02s02-P1-02: Detect true with Cargo.toml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))

		p := crates.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when Cargo.toml exists")
			}
		})
	})

	t.Run("SC-e02s02-P1-03: Detect false without Cargo.toml", func(t *testing.T) {
		dir := t.TempDir()
		p := crates.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when Cargo.toml is absent")
			}
		})
	})
}

func TestCratesPrepare(t *testing.T) {
	t.Run("SC-e02s02-P1-04: Prepare updates version in Cargo.toml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))

		p := crates.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			var cfg struct {
				Package struct {
					Name    string `toml:"name"`
					Version string `toml:"version"`
				} `toml:"package"`
			}
			data := readFile(t, filepath.Join(dir, "Cargo.toml"))
			if err := toml.Unmarshal(data, &cfg); err != nil {
				t.Fatalf("failed to parse updated Cargo.toml: %v", err)
			}
			if cfg.Package.Version != "2.0.0" {
				t.Errorf("expected version 2.0.0, got %q", cfg.Package.Version)
			}
		})
	})

	t.Run("SC-e02s02-P1-05: Prepare returns error when Cargo.toml missing", func(t *testing.T) {
		dir := t.TempDir()
		p := crates.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error when Cargo.toml is missing, got nil")
			}
		})
	})

	t.Run("SC-e02s02-P1-06: Prepare returns error on malformed TOML", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", `[package\ninvalid`)

		p := crates.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error for malformed TOML, got nil")
			}
		})
	})
}

func TestCratesEnvValidation(t *testing.T) {
	t.Run("SC-e02s02-P1-13: Publish returns error when CARGO_TOKEN is empty", func(t *testing.T) {
		_ = os.Unsetenv("CARGO_TOKEN")
		p := crates.NewPublisher()
		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error when CARGO_TOKEN is empty, got nil")
		}
		if !strings.Contains(err.Error(), "CARGO_TOKEN") {
			t.Errorf("expected error message to mention CARGO_TOKEN, got %q", err.Error())
		}
	})
}

func TestCratesAutoRegistration(t *testing.T) {
	t.Run("SC-e02s02-P1-16: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("crates")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "crates" {
			t.Errorf("expected name %q, got %q", "crates", got.Name())
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

func cargoTomlContent(name, version string) string {
	return `[package]
name = "` + name + `"
version = "` + version + `"
edition = "2021"

[dependencies]
`
}

// --- Integration tests ---

func TestCratesPublishHTTP(t *testing.T) {
	t.Run("SC-e02s02-P1-07: Publish returns nil on HTTP 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			if auth := r.Header.Get("Authorization"); auth != "test-token-123" {
				t.Errorf("expected Authorization header, got %q", auth)
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	})

	t.Run("SC-e02s02-P1-08: Publish returns auth error on HTTP 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

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

	t.Run("SC-e02s02-P1-09: Publish returns forbidden error on HTTP 403", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "403") {
				t.Errorf("expected error to mention HTTP 403, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s02-P1-10: Publish returns server error on HTTP 5xx", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

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

	t.Run("SC-e02s02-P1-11a: Publish retries with backoff on 429, succeeds", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
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
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

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

	t.Run("SC-e02s02-P1-11b: Publish returns error after exhausting retries on 429", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"

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

	t.Run("SC-e02s02-P1-12: Publish in dry-run makes zero HTTP requests", func(t *testing.T) {
		requestCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/new", func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := &crates.Publisher{}
		p.RegistryURL = srv.URL + "/api/v1/crates/new"
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

	t.Run("SC-e02s02-P1-14: Verify returns nil on version match", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/test-crate/versions", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]interface{}{
				"versions": []map[string]string{
					{"num": "1.0.0"},
					{"num": "2.0.0"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"
		p.VerifyURL = srv.URL + "/api/v1/crates"

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err != nil {
				t.Fatalf("expected nil on version match, got %v", err)
			}
		})
	})

	t.Run("SC-e02s02-P1-15a: Verify returns error on HTTP 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/test-crate/versions", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"
		p.VerifyURL = srv.URL + "/api/v1/crates"

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

	t.Run("SC-e02s02-P1-15b: Verify returns error when version not found", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/crates/test-crate/versions", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]interface{}{
				"versions": []map[string]string{
					{"num": "1.0.0"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", cargoTomlContent("test-crate", "0.1.0"))
		setenv(t, "CARGO_TOKEN", "test-token-123")

		p := crates.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/crates/new"
		p.VerifyURL = srv.URL + "/api/v1/crates"

		withDir(t, dir, func() {
			err := p.Verify("9.9.9")
			if err == nil {
				t.Fatal("expected error when version not found, got nil")
			}
		})
	})
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
