package godot_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/godot"
)

// --- Unit tests ---

func TestGodotName(t *testing.T) {
	t.Run("SC-e02s07-P3-01: Name returns 'godot'", func(t *testing.T) {
		p := godot.NewPublisher()
		if name := p.Name(); name != "godot" {
			t.Errorf("expected Name() == %q, got %q", "godot", name)
		}
	})
}

func TestGodotDetect(t *testing.T) {
	t.Run("SC-e02s07-P3-02: Detect true when project.godot exists", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "project.godot", projectGodotContent("My Game", "1.0.0"))
		p := godot.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when project.godot exists")
			}
		})
	})

	t.Run("SC-e02s07-P3-03: Detect false when project.godot absent", func(t *testing.T) {
		dir := t.TempDir()
		p := godot.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when project.godot is absent")
			}
		})
	})
}

func TestGodotPrepare(t *testing.T) {
	t.Run("SC-e02s07-P3-04: Prepare updates config/version in project.godot", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "project.godot", projectGodotContent("My Game", "0.1.0"))

		p := godot.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			data := readFile(t, filepath.Join(dir, "project.godot"))
			if !strings.Contains(string(data), `config/version="2.0.0"`) {
				t.Errorf("expected config/version=\"2.0.0\" in project.godot, got:\n%s", string(data))
			}
		})
	})

	t.Run("SC-e02s07-P3-05: Prepare returns error when project.godot missing", func(t *testing.T) {
		dir := t.TempDir()
		p := godot.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error when project.godot is missing, got nil")
			}
		})
	})
}

func TestGodotEnvValidation(t *testing.T) {
	t.Run("SC-e02s07-P3-11: Publish returns error when GITHUB_TOKEN is empty", func(t *testing.T) {
		_ = os.Unsetenv("GITHUB_TOKEN")
		p := godot.NewPublisher()
		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error when GITHUB_TOKEN is empty, got nil")
		}
		if !strings.Contains(err.Error(), "GITHUB_TOKEN") {
			t.Errorf("expected error message to mention GITHUB_TOKEN, got %q", err.Error())
		}
	})
}

func TestGodotAutoRegistration(t *testing.T) {
	t.Run("SC-e02s07-P3-14: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("godot")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "godot" {
			t.Errorf("expected name %q, got %q", "godot", got.Name())
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

// --- test data helpers ---

func projectGodotContent(name, version string) string {
	return `; Engine configuration

[application]
config/name="` + name + `"
config/version="` + version + `"
run/main_scene="res://main.tscn"

[rendering]
quality/driver/driver_name="OpenGL 3"
`
}

// --- Integration tests ---

func TestGodotPublishHTTP(t *testing.T) {
	t.Run("SC-e02s07-P3-06: Publish creates GitHub Release, returns nil on 201", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if auth := r.Header.Get("Authorization"); auth != "token test-token-123" {
				t.Errorf("expected Authorization header, got %q", auth)
			}
			w.WriteHeader(http.StatusCreated)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("SC-e02s07-P3-07: Publish returns auth error on HTTP 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "401") {
			t.Errorf("expected error to mention HTTP 401, got %q", err.Error())
		}
	})

	t.Run("SC-e02s07-P3-08: Publish returns server error on HTTP 5xx", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error to mention HTTP 500, got %q", err.Error())
		}
	})

	t.Run("SC-e02s07-P3-09: Publish retries with backoff on 429", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected eventual success, got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("SC-e02s07-P3-10: Publish dry-run makes zero HTTP requests", func(t *testing.T) {
		requestCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases", func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusCreated)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")

		p := &godot.Publisher{}
		p.GitHubAPI = srv.URL
		p.DryRun = true

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected nil in dry-run, got %v", err)
		}
		if requestCount != 0 {
			t.Errorf("expected 0 HTTP requests in dry-run, got %d", requestCount)
		}
	})

	t.Run("SC-e02s07-P3-12: Verify returns nil on version match", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases/tags/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Verify("2.0.0")
		if err != nil {
			t.Fatalf("expected nil on version match, got %v", err)
		}
	})

	t.Run("SC-e02s07-P3-13: Verify returns error on HTTP 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/testowner/testrepo/releases/tags/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "GITHUB_TOKEN", "test-token-123")
		setenv(t, "GITHUB_OWNER", "testowner")
		setenv(t, "GITHUB_REPO", "testrepo")

		p := godot.NewPublisher()
		p.GitHubAPI = srv.URL

		err := p.Verify("2.0.0")
		if err == nil {
			t.Fatal("expected error on 404, got nil")
		}
		if !strings.Contains(err.Error(), "404") {
			t.Errorf("expected error to mention HTTP 404, got %q", err.Error())
		}
	})
}
