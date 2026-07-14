package httputil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryClient_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestRetryClient_RetryOn429(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	client.MaxRetries = 3
	client.RetryBase = time.Millisecond // No delay for tests

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
}

func TestRetryClient_ExhaustRetries(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	client.MaxRetries = 2
	client.RetryBase = time.Millisecond

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}

	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts (1 initial + 2 retries), got %d", got)
	}
}

func TestRetryClient_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestRetryClient_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestRetryClient_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrServerError) {
		t.Fatalf("expected ErrServerError, got %v", err)
	}
}

func TestRetryClient_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	client := NewRetryClient(server.Client())
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryClient_Defaults(t *testing.T) {
	client := NewRetryClient(nil)
	if client.HTTPClient != http.DefaultClient {
		t.Fatal("expected default HTTP client")
	}
	if client.MaxRetries != 3 {
		t.Fatalf("expected MaxRetries=3, got %d", client.MaxRetries)
	}
	if client.RetryBase != 1*time.Second {
		t.Fatalf("expected RetryBase=1s, got %v", client.RetryBase)
	}
}
