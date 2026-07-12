#!/usr/bin/env bash
# setup.sh — Idempotent environment setup for big-release.
# Safe to run multiple times. Produces the same result on every run.

set -euo pipefail

echo "big-release setup"

# Check Go
if command -v go >/dev/null 2>&1; then
  echo "✅ Go: $(go version)"
else
  echo "❌ Go not found. Install: https://go.dev/dl/"
  exit 1
fi

# Check git
if command -v git >/dev/null 2>&1; then
  echo "✅ Git: $(git --version)"
else
  echo "❌ Git not found."
  exit 1
fi

# Check golangci-lint
if command -v golangci-lint >/dev/null 2>&1; then
  echo "✅ golangci-lint: $(golangci-lint --version)"
else
  echo "⚠  golangci-lint not found. Install: https://golangci-lint.run/welcome/install/"
fi

# Check jq (required by guard-git hooks)
if command -v jq >/dev/null 2>&1; then
  echo "✅ jq: $(jq --version)"
else
  echo "⚠  jq not found. Required for guard-git hooks. Install: brew install jq"
fi

# Download Go dependencies (idempotent — no-op if already cached)
echo ""
echo "→ go mod download..."
go mod download
echo "✅ Dependencies ready"

# Build binary
echo ""
echo "→ make build..."
make build
echo "✅ Binary built"

# Configure git hooks path
if [ -d .githooks ]; then
  git config core.hooksPath .githooks
  echo "✅ Git hooks configured"
fi

echo ""
echo "✅ Setup complete"
