package httputil

import "errors"

// Sentinel errors for common HTTP status codes.
var (
	// ErrUnauthorized indicates a 401 response.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates a 403 response.
	ErrForbidden = errors.New("forbidden")
	// ErrServerError indicates a 5xx response.
	ErrServerError = errors.New("server error")
)
