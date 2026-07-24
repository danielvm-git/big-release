package main

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

// --- BUG-release-workflow-softprops-and-verbose: verbose logger ---

func TestBuildLogger_VerboseEnablesDebugLevel(t *testing.T) {
	// --verbose must produce a logger that emits Debug-level messages so a
	// successful release is observable in CI.
	logger, err := buildLogger(true)
	if err != nil {
		t.Fatalf("buildLogger(true) failed: %v", err)
	}
	if !logger.Core().Enabled(zapcore.DebugLevel) {
		t.Error("verbose logger should emit Debug-level messages")
	}
	_ = logger.Sync()
}

func TestBuildLogger_NonVerboseHidesDebugLevel(t *testing.T) {
	// The default (non-verbose) logger stays at Info level (production JSON)
	// so Debug noise is suppressed.
	logger, err := buildLogger(false)
	if err != nil {
		t.Fatalf("buildLogger(false) failed: %v", err)
	}
	if logger.Core().Enabled(zapcore.DebugLevel) {
		t.Error("non-verbose logger should NOT emit Debug-level messages")
	}
	if !logger.Core().Enabled(zapcore.InfoLevel) {
		t.Error("non-verbose logger should emit Info-level messages")
	}
	_ = logger.Sync()
}
