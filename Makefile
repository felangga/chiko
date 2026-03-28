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

# --- Version bumping ---
# Usage: make bump-patch | bump-minor | bump-major
# Each target:
#   1. Reads the latest git tag (defaults to v0.0.0 if none exists)
#   2. Increments the appropriate semver component
#   3. Updates internal/entity/app.go so 'go run' also shows the new version
#   4. Creates a new annotated git tag
CURRENT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
CURRENT_VER := $(shell echo "$(CURRENT_TAG)" | sed 's/^v//')
VER_MAJOR   := $(shell echo "$(CURRENT_VER)" | cut -d. -f1)
VER_MINOR   := $(shell echo "$(CURRENT_VER)" | cut -d. -f2)
VER_PATCH   := $(shell echo "$(CURRENT_VER)" | cut -d. -f3)

.PHONY: bump-patch
bump-patch:
	$(eval NEW_VER := $(VER_MAJOR).$(VER_MINOR).$(shell echo $$(($(VER_PATCH)+1))))
	git commit --allow-empty -m "chore: bump version to v$(NEW_VER)"
	git tag -a v$(NEW_VER) -m "Release v$(NEW_VER)"
	@echo "Bumped to v$(NEW_VER) — run 'git push && git push --tags' to publish"

.PHONY: bump-minor
bump-minor:
	$(eval NEW_VER := $(VER_MAJOR).$(shell echo $$(($(VER_MINOR)+1))).0)
	git commit --allow-empty -m "chore: bump version to v$(NEW_VER)"
	git tag -a v$(NEW_VER) -m "Release v$(NEW_VER)"
	@echo "Bumped to v$(NEW_VER) — run 'git push && git push --tags' to publish"

.PHONY: bump-major
bump-major:
	$(eval NEW_VER := $(shell echo $$(($(VER_MAJOR)+1))).0.0)
	git commit --allow-empty -m "chore: bump version to v$(NEW_VER)"
	git tag -a v$(NEW_VER) -m "Release v$(NEW_VER)"
	@echo "Bumped to v$(NEW_VER) — run 'git push && git push --tags' to publish"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Clean, test, and build the application"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for Linux, macOS, and Windows"
	@echo "  test        - Run all tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Run the application"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  lint        - Run linter"
	@echo "  bump-patch  - Tag a new patch release (x.y.Z+1)"
	@echo "  bump-minor  - Tag a new minor release (x.Y+1.0)"
	@echo "  bump-major  - Tag a new major release (X+1.0.0)"
	@echo "  help        - Show this help message"