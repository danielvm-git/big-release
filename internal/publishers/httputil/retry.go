package httputil

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// RetryClient wraps an HTTP client with exponential backoff retry on 429.
type RetryClient struct {
	// HTTPClient is the underlying HTTP client.
	HTTPClient *http.Client
	// MaxRetries is the maximum number of retry attempts (default 3).
	MaxRetries int
	// RetryBase is the base backoff duration (default 1s).
	RetryBase time.Duration
}

// NewRetryClient creates a RetryClient with sensible defaults.
func NewRetryClient(client *http.Client) *RetryClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &RetryClient{
		HTTPClient: client,
		MaxRetries: 3,
		RetryBase:  1 * time.Second,
	}
}

// Do executes an HTTP request with retry logic for 429 responses.
func (rc *RetryClient) Do(req *http.Request) (*http.Response, error) {
	if rc.MaxRetries == 0 {
		rc.MaxRetries = 3
	}
	if rc.RetryBase == 0 {
		rc.RetryBase = 1 * time.Second
	}

	backoff := rc.RetryBase
	var lastErr error

	for attempt := 0; attempt <= rc.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(
				float64(backoff)*2,
				float64(rc.RetryBase*time.Duration(math.Pow(2, float64(rc.MaxRetries)))),
			))
		}

		resp, err := rc.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		switch {
		case resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated:
			return resp, nil
		case resp.StatusCode == http.StatusTooManyRequests:
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("rate limited (HTTP %d)", resp.StatusCode)
			continue
		case resp.StatusCode == http.StatusUnauthorized:
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%w (HTTP %d)", ErrUnauthorized, resp.StatusCode)
		case resp.StatusCode == http.StatusForbidden:
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%w (HTTP %d)", ErrForbidden, resp.StatusCode)
		case resp.StatusCode >= 500:
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("%w (HTTP %d)", ErrServerError, resp.StatusCode)
		default:
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("unexpected status (HTTP %d)", resp.StatusCode)
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", rc.MaxRetries, lastErr)
}
