package maven_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/maven"
)

func TestMavenName(t *testing.T) {
	t.Run("SC-e02s05-P2-01: Name returns 'maven'", func(t *testing.T) {
		p := maven.NewPublisher()
		if name := p.Name(); name != "maven" {
			t.Errorf("expected Name() == %q, got %q", "maven", name)
		}
	})
}

func TestMavenDetect(t *testing.T) {
	t.Run("SC-e02s05-P2-02: Detect true when pom.xml exists", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", pomXMLContent("com.example", "test-artifact", "1.0.0"))
		p := maven.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when pom.xml exists")
			}
		})
	})

	t.Run("SC-e02s05-P2-03: Detect false when pom.xml absent", func(t *testing.T) {
		dir := t.TempDir()
		p := maven.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when pom.xml absent")
			}
		})
	})
}

func TestMavenPrepare(t *testing.T) {
	t.Run("SC-e02s05-P2-04: Prepare updates version in pom.xml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", pomXMLContent("com.example", "test-artifact", "1.0.0"))
		p := maven.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			data := readFile(t, filepath.Join(dir, "pom.xml"))
			content := string(data)
			if !strings.Contains(content, "<version>2.0.0</version>") {
				t.Errorf("expected version 2.0.0 in pom.xml, got:\n%s", content)
			}
		})
	})

	t.Run("SC-e02s05-P2-05: Prepare returns error when file missing", func(t *testing.T) {
		dir := t.TempDir()
		p := maven.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error when pom.xml missing, got nil")
			}
			if !strings.Contains(err.Error(), "pom.xml not found") {
				t.Errorf("expected error to mention pom.xml not found, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s05-P2-06: Prepare returns error on malformed XML", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", `<project><version>1.0.0</version></unclosed>`)
		p := maven.NewPublisher()
		withDir(t, dir, func() {
			err := p.Prepare("2.0.0")
			if err == nil {
				t.Fatal("expected error on malformed XML, got nil")
			}
		})
	})
}

func TestMavenEnvValidation(t *testing.T) {
	t.Run("SC-e02s05-P2-13: Publish returns error when MAVEN_TOKEN is empty", func(t *testing.T) {
		_ = os.Unsetenv("MAVEN_TOKEN")
		p := maven.NewPublisher()
		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error when MAVEN_TOKEN is empty, got nil")
		}
		if !strings.Contains(err.Error(), "MAVEN_TOKEN") {
			t.Errorf("expected error message to mention MAVEN_TOKEN, got %q", err.Error())
		}
	})
}

func TestMavenAutoRegistration(t *testing.T) {
	t.Run("Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("maven")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "maven" {
			t.Errorf("expected name %q, got %q", "maven", got.Name())
		}
	})
}

func TestMavenPublishHTTP(t *testing.T) {
	t.Run("SC-e02s05-P2-07: Publish sends POST, returns nil on 200", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if auth := r.Header.Get("Authorization"); auth != "Bearer test-token-123" {
				t.Errorf("expected Authorization: Bearer, got %q", auth)
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")
		p := maven.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("SC-e02s05-P2-08: Publish auth error on HTTP 401", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")
		p := maven.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"

		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "401") {
			t.Errorf("expected error to mention HTTP 401, got %q", err.Error())
		}
	})

	t.Run("SC-e02s05-P2-09: Publish forbidden on HTTP 403", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")
		p := maven.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"

		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "403") {
			t.Errorf("expected error to mention HTTP 403, got %q", err.Error())
		}
	})

	t.Run("SC-e02s05-P2-10: Publish server error on 5xx", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")
		p := maven.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"

		err := p.Publish("1.0.0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error to mention HTTP 500, got %q", err.Error())
		}
	})

	t.Run("SC-e02s05-P2-11: Publish retries on 429 then succeeds", func(t *testing.T) {
		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")
		p := maven.NewPublisher()
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected eventual success, got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("SC-e02s05-P2-12: Publish dry-run makes zero HTTP requests", func(t *testing.T) {
		requestCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/publisher/upload", func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		setenv(t, "MAVEN_TOKEN", "test-token-123")

		p := &maven.Publisher{}
		p.RegistryURL = srv.URL + "/api/v1/publisher/upload"
		p.DryRun = true

		err := p.Publish("1.0.0")
		if err != nil {
			t.Fatalf("expected nil in dry-run, got %v", err)
		}
		if requestCount != 0 {
			t.Errorf("expected 0 HTTP requests in dry-run, got %d", requestCount)
		}
	})
}

func TestMavenVerifyHTTP(t *testing.T) {
	t.Run("SC-e02s05-P2-14: Verify returns nil on match", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/solrsearch/select", func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query().Get("q")
			if !strings.Contains(q, "com.example") || !strings.Contains(q, "test-artifact") || !strings.Contains(q, "2.0.0") {
				t.Errorf("unexpected query: %s", q)
			}
			resp := map[string]interface{}{
				"response": map[string]interface{}{
					"numFound": 1,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", pomXMLContent("com.example", "test-artifact", "1.0.0"))

		p := maven.NewPublisher()
		withDir(t, dir, func() {
			p.VerifyURL = srv.URL + "/solrsearch/select"
			err := p.Verify("2.0.0")
			if err != nil {
				t.Fatalf("expected nil on match, got %v", err)
			}
		})
	})

	t.Run("SC-e02s05-P2-14b: Verify returns error on 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/solrsearch/select", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", pomXMLContent("com.example", "test-artifact", "1.0.0"))

		p := maven.NewPublisher()
		withDir(t, dir, func() {
			p.VerifyURL = srv.URL + "/solrsearch/select"
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

func pomXMLContent(groupID, artifactID, version string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>` + groupID + `</groupId>
    <artifactId>` + artifactID + `</artifactId>
    <version>` + version + `</version>
    <packaging>jar</packaging>
</project>`
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
