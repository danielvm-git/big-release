// story: e08s01
package secure

import (
	"bytes"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMasking_Pattern(t *testing.T) {
	in := "failed with token=ghp_abcdefghijklmnopqrstuvwxyz1234567890"
	out := Redact(in)
	if strings.Contains(out, "ghp_") {
		t.Fatalf("expected token redacted, got %q", out)
	}
	if !strings.Contains(out, Redacted) {
		t.Fatalf("expected %q in output, got %q", Redacted, out)
	}
}

func TestMasking_KnownSecrets(t *testing.T) {
	const secret = "ghp_testknownsecrettokenvalue1234567890"
	t.Setenv("GH_TOKEN", secret)
	in := "auth failed: " + secret
	out := RedactKnownSecrets(in)
	if strings.Contains(out, secret) {
		t.Fatalf("known secret must be redacted, got %q", out)
	}
	if !strings.Contains(out, Redacted) {
		t.Fatalf("expected %q in output", Redacted)
	}
}

func TestMasking_ZapCore(t *testing.T) {
	const secret = "npm_testzaptokensecretvalue1234567890"
	t.Setenv("NPM_TOKEN", secret)

	var buf bytes.Buffer
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := WrapCore(zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.InfoLevel))
	logger := zap.New(core)
	logger.Info("publish failed", zap.String("detail", "bad token "+secret))

	out := buf.String()
	if strings.Contains(out, secret) {
		t.Fatalf("zap output leaked secret: %s", out)
	}
}
