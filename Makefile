.PHONY: all build test lint fmt clean preflight help

# Variables
BINARY_NAME := big-release
BUILD_DIR := ./bin
CMD_DIR := ./cmd/big-release
GO_FILES := $(shell find . -name '*.go' -type f -not -path './vendor/*')

# Default target
all: preflight build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Built $(BUILD_DIR)/$(BINARY_NAME)"

# Run all tests
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	@go test -v -race ./internal/...

# Run integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	@go test -v -race -tags=integration ./tests/integration/...

# Lint the code
lint:
	@echo "🔍 Linting..."
	@golangci-lint run ./...

# Format the code
fmt:
	@echo "✨ Formatting..."
	@gofmt -s -w $(GO_FILES)
	@goimports -w $(GO_FILES)

# Vet the code
vet:
	@echo "🔍 Vetting..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Preflight checks (must pass before any forward work)
preflight: lint vet test
	@echo "✅ Preflight PASS"

# Run the binary
run:
	@go run $(CMD_DIR)

# Install the binary
install:
	@echo "📦 Installing $(BINARY_NAME)..."
	@go install $(CMD_DIR)

# Cross-compile for all platforms
release:
	@echo "🚀 Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "✅ Release binaries built in $(BUILD_DIR)/release/"

# Validate configuration
validate:
	@echo "🔍 Validating configuration..."
	@$(BUILD_DIR)/$(BINARY_NAME) validate

# Show help
help:
	@echo "big-release - Unified release automation"
	@echo ""
	@echo "Targets:"
	@echo "  make build          Build the binary"
	@echo "  make test           Run all tests"
	@echo "  make test-unit      Run unit tests only"
	@echo "  make test-integration Run integration tests"
	@echo "  make lint           Lint the code"
	@echo "  make fmt            Format the code"
	@echo "  make vet            Vet the code"
	@echo "  make clean          Clean build artifacts"
	@echo "  make preflight      Run all checks (must pass before forward work)"
	@echo "  make run            Run the binary"
	@echo "  make install        Install the binary"
	@echo "  make release        Cross-compile for all platforms"
	@echo "  make validate       Validate configuration"
	@echo "  make help           Show this help"
