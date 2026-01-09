.PHONY: build run test clean install deps examples

# Binary name
BINARY_NAME=synq-pipe
BINARY_PATH=./bin/$(BINARY_NAME)

# Build variables
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Default target
all: deps build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Build binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_PATH) cmd/synq-pipe/main.go
	@echo "✓ Build complete: $(BINARY_PATH)"

# Build for production (optimized)
build-prod:
	@echo "🚀 Building for production..."
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build -trimpath $(LDFLAGS) -o $(BINARY_PATH) cmd/synq-pipe/main.go
	@echo "✓ Production build complete: $(BINARY_PATH)"

# Build for multiple platforms
build-all:
	@echo "🌍 Building for multiple platforms..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 cmd/synq-pipe/main.go
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 cmd/synq-pipe/main.go
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 cmd/synq-pipe/main.go
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe cmd/synq-pipe/main.go
	@echo "✓ Multi-platform build complete"

# Install binary to system
install: build
	@echo "📥 Installing $(BINARY_NAME)..."
	@mkdir -p $(GOPATH)/bin
	@cp $(BINARY_PATH) $(GOPATH)/bin/
	@echo "✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Run with example workflow
run: build
	@echo "▶️  Running example workflow..."
	$(BINARY_PATH) run examples/simple.yaml --tui

# Run Ghost mode example
run-ghost: build
	@echo "👻 Running Ghost mode example..."
	$(BINARY_PATH) run examples/ghost.yaml --tui

# Run all examples
examples: build
	@echo "📚 Running all examples..."
	@echo "\n=== Simple Pipeline ==="
	$(BINARY_PATH) run examples/simple.yaml
	@echo "\n=== Ghost Mode ==="
	$(BINARY_PATH) run examples/ghost.yaml
	@echo "\n=== Advanced Pipeline ==="
	$(BINARY_PATH) run examples/advanced.yaml

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
fmt:
	@echo "✨ Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Code formatted"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Development mode (rebuild on changes)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

# Show help
help:
	@echo "Synqly Multiverse CLI - Makefile Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make deps           Install dependencies"
	@echo "  make build          Build binary"
	@echo "  make install        Install binary to system"
	@echo ""
	@echo "Development:"
	@echo "  make dev            Run in development mode (auto-reload)"
	@echo "  make fmt            Format code"
	@echo "  make lint           Lint code"
	@echo "  make test           Run tests"
	@echo "  make test-coverage  Run tests with coverage"
	@echo ""
	@echo "Running:"
	@echo "  make run            Run simple example"
	@echo "  make run-ghost      Run Ghost mode example"
	@echo "  make examples       Run all examples"
	@echo ""
	@echo "Building:"
	@echo "  make build-prod     Production build (optimized)"
	@echo "  make build-all      Build for all platforms"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          Remove build artifacts"