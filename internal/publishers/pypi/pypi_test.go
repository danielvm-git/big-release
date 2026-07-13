package pypi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/pypi"
)

// --- Unit tests ---

func TestPyPIName(t *testing.T) {
	t.Run("SC-e02s01-P1-01: Name returns 'pypi'", func(t *testing.T) {
		p := pypi.NewPublisher()
		if name := p.Name(); name != "pypi" {
			t.Errorf("expected Name() == %q, got %q", "pypi", name)
		}
	})
}

func TestPyPIDetect(t *testing.T) {
	t.Run("SC-e02s01-P1-02a: Detect true with pyproject.toml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pyproject.toml", `[project]\nname = "test-pkg"\nversion = "0.1.0"\n`)
		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when pyproject.toml exists")
			}
		})
	})

	t.Run("SC-e02s01-P1-02b: Detect true with setup.py", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "setup.py", `from setuptools import setup\nsetup(name="test-pkg", version="0.1.0")\n`)
		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when setup.py exists")
			}
		})
	})

	t.Run("SC-e02s01-P1-03: Detect false when neither file exists", func(t *testing.T) {
		dir := t.TempDir()
		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when neither config file exists")
			}
		})
	})
}

func TestPyPIPrepare(t *testing.T) {
	t.Run("SC-e02s01-P1-04a: Prepare updates version in pyproject.toml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "0.1.0"))

		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			data := readFile(t, filepath.Join(dir, "pyproject.toml"))
			if !strings.Contains(string(data), `version = "2.0.0"`) {
				t.Errorf("expected version 2.0.0 in pyproject.toml, got:\n%s", string(data))
			}
		})
	})

	t.Run("SC-e02s01-P1-04b: Prepare updates version in setup.cfg", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "setup.cfg", setupCfgContent("test-pkg", "0.1.0"))

		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("3.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			data := readFile(t, filepath.Join(dir, "setup.cfg"))
			if !strings.Contains(string(data), "version = 3.0.0") {
				t.Errorf("expected version 3.0.0 in setup.cfg, got:\n%s", string(data))
			}
		})
	})

	t.Run("SC-e02s01-P1-05: Prepare returns error when both files absent", func(t *testing.T) {
		dir := t.TempDir()
		p := pypi.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error when no config file exists, got nil")
			}
		})
	})
}

func TestPyPIEnvValidation(t *testing.T) {
	t.Run("SC-e02s01-P1-13: Publish returns error when PYPI_TOKEN is empty", func(t *testing.T) {
		// Ensure PYPI_TOKEN is unset
		_ = os.Unsetenv("PYPI_TOKEN")
		p := pypi.NewPublisher()
		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error when PYPI_TOKEN is empty, got nil")
		}
		if !strings.Contains(err.Error(), "PYPI_TOKEN") {
			t.Errorf("expected error message to mention PYPI_TOKEN, got %q", err.Error())
		}
	})
}

func TestPyPIAutoRegistration(t *testing.T) {
	t.Run("SC-e02s01-P1-16: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("pypi")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "pypi" {
			t.Errorf("expected name %q, got %q", "pypi", got.Name())
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

// --- test data helpers ---

func pyprojectTomlContent(name, version string) string {
	return `[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "` + name + `"
version = "` + version + `"
description = "A test package"
`
}

func setupCfgContent(name, version string) string {
	return `[metadata]
name = ` + name + `
version = ` + version + `
description = A test package

[options]
packages = find:
`
}

// --- Integration tests ---

func TestPyPIPublishHTTP(t *testing.T) {
	t.Run("SC-e02s01-P1-06: Publish returns nil on HTTP 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if auth := r.Header.Get("Authorization"); auth != "token test-token-123" {
				t.Errorf("expected Authorization header, got %q", auth)
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	})

	t.Run("SC-e02s01-P1-07: Publish returns auth error on HTTP 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

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

	t.Run("SC-e02s01-P1-08: Publish returns forbidden error on HTTP 403", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

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

	t.Run("SC-e02s01-P1-09: Publish returns server error on HTTP 5xx", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

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

	t.Run("SC-e02s01-P1-10: Publish retries with backoff on 429, succeeds", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
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
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

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

	t.Run("SC-e02s01-P1-11: Publish returns error after exhausting retries on 429", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setupDistDir(t, dir)
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "1.0.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.RegistryURL = srv.URL + "/legacy/"

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

	t.Run("SC-e02s01-P1-12: Publish in dry-run makes zero HTTP requests", func(t *testing.T) {
		requestCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/legacy/", func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := &pypi.Publisher{}
		p.RegistryURL = srv.URL + "/legacy/"
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

	t.Run("SC-e02s01-P1-14: Verify returns nil on version match", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/pypi/test-pkg/json", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]interface{}{
				"info": map[string]string{
					"version": "2.0.0",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "0.1.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.VerifyURL = srv.URL + "/pypi"

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err != nil {
				t.Fatalf("expected nil on version match, got %v", err)
			}
		})
	})

	t.Run("SC-e02s01-P1-15: Verify returns error on HTTP 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/pypi/test-pkg/json", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "pyproject.toml", pyprojectTomlContent("test-pkg", "0.1.0"))
		setenv(t, "PYPI_TOKEN", "test-token-123")

		p := pypi.NewPublisher()
		p.VerifyURL = srv.URL + "/pypi"

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

// --- integration test helpers ---

func setupDistDir(t *testing.T, dir string) {
	t.Helper()
	distDir := filepath.Join(dir, "dist")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		t.Fatalf("failed to create dist dir: %v", err)
	}
	writeFile(t, distDir, "test_pkg-1.0.0.tar.gz", "fake-tarball-content")
	writeFile(t, distDir, "test_pkg-1.0.0-py3-none-any.whl", "fake-wheel-content")
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
