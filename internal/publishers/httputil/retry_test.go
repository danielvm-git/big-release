// story: e05s01

package httputil

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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

func TestRetryClient_BodyRewindOnRetry(t *testing.T) {
	const payload = `{"name":"pkg"}`
	var attempts int
	var bodies []string

	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		bodies = append(bodies, string(body))
		attempts++
		if attempts == 1 {
			return &http.Response{StatusCode: http.StatusTooManyRequests, Body: http.NoBody}, nil
		}
		return &http.Response{StatusCode: http.StatusCreated, Body: http.NoBody}, nil
	})

	client := NewRetryClient(&http.Client{Transport: rt})
	client.MaxRetries = 2
	client.RetryBase = time.Millisecond

	req, err := http.NewRequest(http.MethodPost, "https://example.test", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	if req.GetBody == nil {
		t.Fatal("test setup: expected GetBody to be set")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected successful retry, got %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if len(bodies) != 2 {
		t.Fatalf("expected 2 request bodies, got %d", len(bodies))
	}
	for i, body := range bodies {
		if body != payload {
			t.Fatalf("attempt %d sent body %q, want %q", i+1, body, payload)
		}
	}
}

func TestRetryClient_BodyWithoutGetBodyFails(t *testing.T) {
	var attempts int32
	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt32(&attempts, 1)
		return &http.Response{StatusCode: http.StatusTooManyRequests, Body: http.NoBody}, nil
	})

	client := NewRetryClient(&http.Client{Transport: rt})
	client.MaxRetries = 3
	client.RetryBase = time.Millisecond

	req, err := http.NewRequest(http.MethodPost, "https://example.test", strings.NewReader("payload"))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.GetBody = nil

	_, err = client.Do(req)
	if err == nil {
		t.Fatal("expected error for a request body that cannot be rewound")
	}
	if !strings.Contains(err.Error(), "GetBody") {
		t.Fatalf("expected descriptive error mentioning GetBody, got %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("expected no retry attempt after non-rewindable body, got %d attempts", got)
	}
}

func TestRetryClient_BodylessRetryStillWorks(t *testing.T) {
	tests := []struct {
		name string
		body io.Reader
	}{
		{name: "nil body", body: nil},
		{name: "http.NoBody", body: http.NoBody},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts int32
			rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if atomic.AddInt32(&attempts, 1) == 1 {
					return &http.Response{StatusCode: http.StatusTooManyRequests, Body: http.NoBody}, nil
				}
				return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
			})

			client := NewRetryClient(&http.Client{Transport: rt})
			client.MaxRetries = 2
			client.RetryBase = time.Millisecond

			req, err := http.NewRequest(http.MethodPost, "https://example.test", tt.body)
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("expected bodyless retry to succeed, got %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}
			if got := atomic.LoadInt32(&attempts); got != 2 {
				t.Fatalf("expected 2 attempts, got %d", got)
			}
		})
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
