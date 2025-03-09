# Makefile for Chiko project

# Project variables
BINARY_NAME=chiko
VERSION=$(shell git describe --tags --always --dirty)
BUILD_DIR=build
CMD_DIR=cmd/chiko
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: clean test build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	
# Install the application
.PHONY: install
install: 
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)

# Build for all platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)/main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)/main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)/main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	$(GOCMD) clean -cache

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(CMD_DIR)/main.go

# Install dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Lint the code
.PHONY: lint
lint:
	@which golangci-lint >/dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       - Clean, test, and build the application"
	@echo "  build     - Build the application"
	@echo "  build-all - Build for Linux, macOS, and Windows"
	@echo "  test      - Run all tests"
	@echo "  clean     - Remove build artifacts"
	@echo "  run       - Run the application"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  lint      - Run linter"
	@echo "  help      - Show this help message"