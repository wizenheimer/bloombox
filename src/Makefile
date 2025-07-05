.PHONY: build run test clean deps fmt lint vet install dev

# Build variables
BINARY_NAME=bloombox
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin

# Colors
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
all: build

# Build the application
build:
	@echo "$(BLUE)"
	@echo "██████  ██      ██████   ██████  ███    ███ ██████   ██████  ██   ██ "
	@echo "██   ██ ██     ██    ██ ██    ██ ████  ████ ██   ██ ██    ██  ██ ██  "
	@echo "██████  ██     ██    ██ ██    ██ ██ ████ ██ ██████  ██    ██   ███   "
	@echo "██   ██ ██     ██    ██ ██    ██ ██  ██  ██ ██   ██ ██    ██  ██ ██  "
	@echo "██████  ███████ ██████   ██████  ██      ██ ██████   ██████  ██   ██ "
	@echo "$(NC)"
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Run the application
run:
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@go run ./cmd/server

# Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installed$(NC)"

# Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)Tests completed$(NC)"

# Run tests with coverage
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

# Run linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)Linting completed$(NC)"

# Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Vet completed$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed$(NC)"

# Install the binary
install: build
	@echo "$(BLUE)Installing binary...$(NC)"
	@mv $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)Binary installed to $(GOPATH)/bin/$(NC)"

# Development setup
dev: deps fmt vet test
	@echo "$(GREEN)Development setup completed$(NC)"

# Quick development cycle
dev-run: fmt vet run

# Build and run in one command
build-run: build
	@echo "$(BLUE)Running built binary...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Install golangci-lint if not present
install-lint:
	@echo "$(BLUE)Installing golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)golangci-lint installed$(NC)"

# Run with hot reload (requires air)
dev-hot:
	@echo "$(BLUE)Starting development server with hot reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air not found. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(BLUE)Falling back to regular run...$(NC)"; \
		go run ./cmd/server; \
	fi

# Build for multiple platforms
build-all:
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	@echo "$(GREEN)Built linux-amd64$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	@echo "$(GREEN)Built darwin-amd64$(NC)"
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	@echo "$(GREEN)Built windows-amd64$(NC)"
	@echo "$(GREEN)All builds completed$(NC)"

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)

help:
	@echo "$(BLUE)"
	@echo "██████  ██      ██████   ██████  ███    ███ ██████   ██████  ██   ██ "
	@echo "██   ██ ██     ██    ██ ██    ██ ████  ████ ██   ██ ██    ██  ██ ██  "
	@echo "██████  ██     ██    ██ ██    ██ ██ ████ ██ ██████  ██    ██   ███   "
	@echo "██   ██ ██     ██    ██ ██    ██ ██  ██  ██ ██   ██ ██    ██  ██ ██  "
	@echo "██████  ███████ ██████   ██████  ██      ██ ██████   ██████  ██   ██ "
	@echo "$(NC)"
	@echo "$(YELLOW)Available targets:$(NC)"
	@echo "$(GREEN)  help           $(NC)Show this help message"
	@echo "$(GREEN)  build          $(NC)Build the application"
	@echo "$(GREEN)  run            $(NC)Run the application"
	@echo "$(GREEN)  deps           $(NC)Install dependencies"
	@echo "$(GREEN)  test           $(NC)Run tests"
	@echo "$(GREEN)  test-coverage  $(NC)Run tests with coverage report"
	@echo "$(GREEN)  fmt            $(NC)Format code"
	@echo "$(GREEN)  lint           $(NC)Run linter"
	@echo "$(GREEN)  vet            $(NC)Run go vet"
	@echo "$(GREEN)  clean          $(NC)Clean build artifacts"
	@echo "$(GREEN)  install        $(NC)Install the binary"
	@echo "$(GREEN)  dev            $(NC)Development setup (deps, fmt, vet, test)"
	@echo "$(GREEN)  dev-run        $(NC)Quick dev cycle (fmt, vet, run)"
	@echo "$(GREEN)  build-run      $(NC)Build and run in one command"
	@echo "$(GREEN)  install-lint   $(NC)Install golangci-lint if not present"
	@echo "$(GREEN)  dev-hot        $(NC)Run with hot reload (requires air)"
	@echo "$(GREEN)  build-all      $(NC)Build for multiple platforms"