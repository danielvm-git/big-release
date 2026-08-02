package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerReturns200(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerSetsContentType(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Fatalf("expected Content-Type text/html; charset=utf-8, got %s", ct)
	}
}

func TestHandlerContainsBigRelease(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "big-release") {
		t.Fatalf("expected body to contain 'big-release'")
	}
}

func TestListenAddrUsesPort(t *testing.T) {
	t.Setenv("PORT", "9999")
	if got := listenAddr(); got != ":9999" {
		t.Fatalf("expected :9999, got %s", got)
	}
}

func TestListenAddrDefaultsTo8080(t *testing.T) {
	t.Setenv("PORT", "")
	if got := listenAddr(); got != ":8080" {
		t.Fatalf("expected :8080, got %s", got)
	}
}
