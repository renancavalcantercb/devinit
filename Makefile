.PHONY: help build build-all test test-unit test-integration test-e2e clean install validate-templates lint fmt dev

# Variables
BINARY_NAME=devinit
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
help:
	@echo "devinit - Development Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build for current platform"
	@echo "  build-all          Build for all platforms (Linux, macOS)"
	@echo "  test               Run all tests"
	@echo "  test-unit          Run unit tests only"
	@echo "  test-integration   Run integration tests only"
	@echo "  test-e2e           Run end-to-end tests"
	@echo "  clean              Remove build artifacts"
	@echo "  install            Install to \$$GOPATH/bin"
	@echo "  validate-templates Validate all templates"
	@echo "  lint               Run linters"
	@echo "  fmt                Format code"
	@echo "  dev                Build and run (development mode)"

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/devinit
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "  Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/devinit

	@echo "  Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/devinit

	@echo "  Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/devinit

	@echo "  Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/devinit

	@echo "✓ All builds complete"
	@ls -lh $(BUILD_DIR)

# Run all tests
test:
	@echo "Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Tests complete. Coverage report: coverage.html"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -race ./internal/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -race ./test/integration/...

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	@go test -v -race ./test/e2e/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Install to $GOPATH/bin
install:
	@echo "Installing to \$$GOPATH/bin..."
	@go install $(LDFLAGS) ./cmd/devinit
	@echo "✓ Installed: $(shell which $(BINARY_NAME))"

# Validate all templates
validate-templates: build
	@echo "Validating templates..."
	@./$(BUILD_DIR)/$(BINARY_NAME) templates validate

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
		go vet ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

# Development mode (build and show usage)
dev: build
	@echo ""
	@./$(BUILD_DIR)/$(BINARY_NAME) --help

# Quick test - generate a sample project
demo: build
	@echo "Creating demo project..."
	@rm -rf /tmp/demo-api
	@./$(BUILD_DIR)/$(BINARY_NAME) new demo-api \
		--lang python \
		--framework fastapi \
		--database postgres
	@echo ""
	@echo "Demo project created at /tmp/demo-api"
	@echo "To run:"
	@echo "  cd /tmp/demo-api"
	@echo "  docker compose up"
