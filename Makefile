# ctx Makefile
# Build, test, and install ctx

# Variables
BINARY_NAME=ctx
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_CLEAN=$(GO_CMD) clean
GO_GET=$(GO_CMD) get
GO_MOD=$(GO_CMD) mod
INSTALL_PATH?=$(HOME)/bin

# Version info
VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "local")
DATE?=$(shell date -u +%Y-%m-%d)

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Detect OS and architecture
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Default target
.DEFAULT_GOAL := build

# Help target
.PHONY: help
help:
	@echo "ctx Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Build ctx binary"
	@echo "  make install    - Build and install ctx to PATH"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make run        - Build and run ctx with test command"
	@echo "  make uninstall  - Remove ctx from PATH"
	@echo "  make release    - Build release binaries for all platforms"
	@echo ""
	@echo "Variables:"
	@echo "  INSTALL_PATH   - Installation directory (default: ~/bin)"
	@echo ""

# Build the binary
.PHONY: build
build:
	@echo "Building ctx..."
	@$(GO_BUILD) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Install ctx to PATH
.PHONY: install
install: build
	@echo "Installing ctx to $(INSTALL_PATH)..."
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed successfully to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo ""
	@if echo $$PATH | grep -q "$(INSTALL_PATH)"; then \
		echo "✓ $(INSTALL_PATH) is in PATH"; \
		echo "You can now use: ctx <command>"; \
	else \
		echo "⚠️  $(INSTALL_PATH) is not in PATH"; \
		echo "Add this to your shell config:"; \
		echo "  export PATH=\"$(INSTALL_PATH):$$PATH\""; \
	fi

# Quick install using ./install.sh
.PHONY: quick-install
quick-install:
	@./install.sh

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO_TEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GO_CLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f dist/*
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO_MOD) download
	@$(GO_MOD) tidy
	@echo "Dependencies ready"

# Build and run with test command
.PHONY: run
run: build
	@echo "Testing ctx..."
	@./$(BINARY_NAME) echo "ctx is working!"

# Uninstall ctx
.PHONY: uninstall
uninstall:
	@echo "Uninstalling ctx..."
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@rm -f $(HOME)/go/bin/$(BINARY_NAME)
	@echo "Uninstalled successfully"

# Build for all platforms
.PHONY: release
release:
	@echo "Building release binaries..."
	@mkdir -p dist
	
	@echo "Building for macOS (Intel)..."
	@GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "Building for macOS (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 $(GO_BUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	
	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 $(GO_BUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "Release binaries built in dist/"
	@ls -lh dist/

# Development build with debug info
.PHONY: dev
dev:
	@echo "Building ctx with debug info..."
	@$(GO_BUILD) -gcflags="all=-N -l" -o $(BINARY_NAME) .
	@echo "Debug build complete: ./$(BINARY_NAME)"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO_CMD) fmt ./...
	@echo "Code formatted"

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Check for updates
.PHONY: update
update:
	@echo "Checking for dependency updates..."
	@$(GO_CMD) list -u -m all

# Benchmark
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	@$(GO_TEST) -bench=. -benchmem ./...

.PHONY: all
all: clean deps test build