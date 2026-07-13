package goproxy_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/goproxy"
)

// --- Unit tests ---

func TestGoProxyName(t *testing.T) {
	t.Run("SC-e02s03-P2-01: Name returns 'goproxy'", func(t *testing.T) {
		p := goproxy.NewPublisher()
		if name := p.Name(); name != "goproxy" {
			t.Errorf("expected Name() == %q, got %q", "goproxy", name)
		}
	})
}

func TestGoProxyDetect(t *testing.T) {
	t.Run("SC-e02s03-P2-02: Detect true with go.mod", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")
		p := goproxy.NewPublisher()
		withDir(t, dir, func() {
			if !p.Detect() {
				t.Errorf("expected Detect() == true when go.mod exists")
			}
		})
	})

	t.Run("SC-e02s03-P2-03: Detect false without go.mod", func(t *testing.T) {
		dir := t.TempDir()
		p := goproxy.NewPublisher()
		withDir(t, dir, func() {
			if p.Detect() {
				t.Errorf("expected Detect() == false when go.mod is absent")
			}
		})
	})
}

func TestGoProxyPrepare(t *testing.T) {
	t.Run("SC-e02s03-P2-04: Prepare(version) is no-op, returns nil", func(t *testing.T) {
		p := goproxy.NewPublisher()
		if err := p.Prepare("1.0.0"); err != nil {
			t.Fatalf("expected nil error from Prepare, got %v", err)
		}
	})
}

func TestGoProxyAutoRegistration(t *testing.T) {
	t.Run("SC-e02s03-P2-12: Auto-registered in global registry via init()", func(t *testing.T) {
		got, err := publishers.Get("goproxy")
		if err != nil {
			t.Fatalf("expected publisher to be registered, got error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil publisher from global registry")
		}
		if got.Name() != "goproxy" {
			t.Errorf("expected name %q, got %q", "goproxy", got.Name())
		}
	})
}

// --- Integration tests ---

func TestGoProxyPublishExec(t *testing.T) {
	t.Run("SC-e02s03-P2-05: Publish(version) pushes tag and runs go list -m", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		var commands []string
		mockExec := func(name string, args ...string) *exec.Cmd {
			commands = append(commands, name+" "+strings.Join(args, " "))
			// Simulate success by returning a no-op command that just exits 0
			switch name {
			case "git", "go":
				return exec.Command("true")
			default:
				return exec.Command("true")
			}
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/1.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		p := goproxy.NewPublisher()
		p.ExecCommand = mockExec
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})

		if len(commands) < 2 {
			t.Fatalf("expected at least 2 exec commands (git tag, git push), got %d: %v", len(commands), commands)
		}
		if !strings.Contains(commands[0], "git tag v1.0.0") {
			t.Errorf("expected git tag command, got %q", commands[0])
		}
		if !strings.Contains(commands[1], "git push origin v1.0.0") {
			t.Errorf("expected git push command, got %q", commands[1])
		}
	})

	t.Run("SC-e02s03-P2-08: Publish in dry-run skips tag push and proxy call", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		execCalls := 0
		mockExec := func(name string, args ...string) *exec.Cmd {
			execCalls++
			return exec.Command("true")
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			t.Error("unexpected HTTP call in dry-run mode")
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		p := &goproxy.Publisher{}
		p.ProxyURL = srv.URL
		p.DryRun = true
		p.ExecCommand = mockExec

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected nil in dry-run, got %v", err)
			}
		})

		if execCalls != 0 {
			t.Errorf("expected 0 exec calls in dry-run, got %d", execCalls)
		}
	})
}

func TestGoProxyPublishHTTP(t *testing.T) {
	t.Run("SC-e02s03-P2-06: Publish returns error on HTTP 5xx from proxy", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		mockExec := func(name string, args ...string) *exec.Cmd {
			return exec.Command("true")
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/1.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		p := goproxy.NewPublisher()
		p.ExecCommand = mockExec
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err == nil {
				t.Fatal("expected error on 5xx, got nil")
			}
			if !strings.Contains(err.Error(), "500") {
				t.Errorf("expected error to mention HTTP 500, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s03-P2-07: Publish retries with backoff on 429, succeeds", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		mockExec := func(name string, args ...string) *exec.Cmd {
			return exec.Command("true")
		}

		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/1.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		p := goproxy.NewPublisher()
		p.ExecCommand = mockExec
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Publish("1.0.0")
			if err != nil {
				t.Fatalf("expected eventual success, got %v", err)
			}
			if attempts != 3 {
				t.Errorf("expected 3 HTTP attempts, got %d", attempts)
			}
		})
	})

	t.Run("SC-e02s03-P2-07b: Publish returns error after exhausting retries on 429", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		mockExec := func(name string, args ...string) *exec.Cmd {
			return exec.Command("true")
		}

		attempts := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/1.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		p := goproxy.NewPublisher()
		p.ExecCommand = mockExec
		p.ProxyURL = srv.URL

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
}

func TestGoProxyVerify(t *testing.T) {
	t.Run("SC-e02s03-P2-09: Verify returns nil on success", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/2.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]string{"Version": "2.0.0"}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		p := goproxy.NewPublisher()
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err != nil {
				t.Fatalf("expected nil on success, got %v", err)
			}
		})
	})

	t.Run("SC-e02s03-P2-10: Verify returns error on HTTP 404", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/9.9.9.info", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		p := goproxy.NewPublisher()
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Verify("9.9.9")
			if err == nil {
				t.Fatal("expected error on 404, got nil")
			}
			if !strings.Contains(err.Error(), "404") {
				t.Errorf("expected error to mention HTTP 404, got %q", err.Error())
			}
		})
	})

	t.Run("SC-e02s03-P2-11: Verify returns error on version mismatch", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/example.com/test/@v/2.0.0.info", func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]string{"Version": "3.0.0"}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/test\n")

		p := goproxy.NewPublisher()
		p.ProxyURL = srv.URL

		withDir(t, dir, func() {
			err := p.Verify("2.0.0")
			if err == nil {
				t.Fatal("expected error on version mismatch, got nil")
			}
			if !strings.Contains(err.Error(), "version mismatch") {
				t.Errorf("expected error to mention version mismatch, got %q", err.Error())
			}
		})
	})
}

// --- helpers ---

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
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
