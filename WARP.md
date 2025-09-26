# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

**Chiko** is a beautiful terminal user interface (TUI) client for gRPC built in Go. It provides an intuitive interface for testing gRPC services, built on top of `grpcurl` and the `tview` library. The application features server reflection, authentication support, metadata management, payload generation, and bookmark functionality.

## Architecture

### Core Structure
- **`cmd/chiko/main.go`**: Application entry point that parses flags and initializes the UI
- **`internal/controller/`**: Business logic layer containing flag parsing, gRPC operations, bookmarks, and storage
- **`internal/entity/`**: Data models and structures (Session, Auth, Bookmark, Theme, etc.)
- **`internal/ui/`**: TUI components using `tview` library for all interface elements
- **`internal/logger/`**: Logging system with channels for real-time log display

### Key Components
- **Session Management**: Core entity that maintains connection state, authentication, and request configuration
- **gRPC Controller**: Handles gRPC connections, reflection, and method invocation via `grpcurl`
- **Bookmark System**: Persistent storage of favorite requests and configurations
- **UI Layout**: Split-pane interface with sidebar menu, bookmark tree, output panel, and log viewer

## Development Commands

### Build & Run
```bash
# Run from source (development)
go run ./cmd/chiko/main.go

# Build binary
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Testing
```bash
# Run all tests
make test
# OR
go test ./...

# Run tests with verbose output
go test -v ./...

# Test specific package
go test ./internal/entity/
```

### Code Quality
```bash
# Run linter (installs golangci-lint if needed)
make lint

# Format code
go fmt ./...

# Vet code
go vet ./...

# Download and tidy dependencies
make deps
```

### Development Workflow
```bash
# Clean build artifacts
make clean

# Complete build pipeline (clean, test, build)
make all

# Run with specific flags (examples)
go run ./cmd/chiko/main.go -plaintext localhost:9090
go run ./cmd/chiko/main.go -plaintext -d '{"name": "test"}' localhost:9090 my.service/Method
```

## Testing Patterns

The project uses Go's standard testing framework with the following patterns:

### Test Structure
- Test files follow `*_test.go` naming convention
- Tests use table-driven test patterns where appropriate
- Focus areas: flag validation, session parsing, bookmark operations, and storage functions

### Running Specific Tests
```bash
# Test a single function
go test -run TestFunctionName ./internal/package/

# Test with coverage
go test -cover ./...

# Run tests for modified packages only
go test ./internal/entity/ ./internal/controller/
```

## Key Dependencies

- **`github.com/rivo/tview`**: Terminal UI framework for the interface
- **`github.com/fullstorydev/grpcurl`**: Core gRPC functionality and reflection
- **`github.com/epiclabs-io/winman`**: Window management for modal dialogs
- **`github.com/gdamore/tcell/v2`**: Terminal cell handling and keyboard events
- **`github.com/google/uuid`**: UUID generation for sessions
- **`github.com/atotto/clipboard`**: Clipboard operations for copy functionality

## Configuration & State

### Session Configuration
The `entity.Session` struct is the central configuration object containing:
- Connection details (server URL, TLS settings)
- Authentication (Bearer tokens)
- Request metadata and payloads
- gRPC call options (timeouts, message sizes)

### Command Line Flags
Chiko accepts grpcurl-compatible flags:
- `-plaintext`: Use HTTP/2 without TLS (default: true)
- `-insecure`: Skip certificate verification
- `-cert`, `-key`, `-cacert`: SSL certificate files
- `-d`: Request data payload
- `-connect-timeout`, `-max-time`: Timeout configurations

### Bookmark Storage
Bookmarks are stored in `.bookmarks` files with JSON format, organized by categories containing session configurations.

## UI Components

### Layout Architecture
- **Main Layout**: Flex container with title bar and split content
- **Sidebar**: Menu list (top) and bookmark tree (bottom)
- **Main Panel**: Output panel (top) and log viewer (bottom)
- **Modal System**: Overlay dialogs for configuration and user input

### Key UI Files
- `init_sidebar_menu.go`: Main navigation menu
- `init_bookmark_menu.go`: Bookmark tree management
- `init_output_panel.go`: Response display and interaction commands
- `init_log_list.go`: Real-time log display

## Release Process

The project uses GitHub Actions with GoReleaser for automated releases:

### Release Workflow
1. Tag version: `git tag v1.0.0 && git push --tags`
2. GitHub Action triggers GoReleaser
3. Builds for Linux, macOS, Windows (amd64, arm64)
4. Creates GitHub release with binaries
5. Updates Homebrew tap (felangga/homebrew-chiko)

### Manual Release Testing
```bash
# Test goreleaser locally
goreleaser release --snapshot --clean
```

## Troubleshooting Development

### Common Issues
- **Build failures**: Ensure Go 1.22.5+ and run `go mod download`
- **TUI rendering issues**: Verify terminal supports 256 colors
- **gRPC connection problems**: Check server reflection support and TLS configuration
- **Test failures**: Ensure no external dependencies or file system permissions issues

### Debug Mode
The application logs to both the UI log panel and can export logs to files. Use the built-in log viewer for debugging TUI interactions and gRPC operations.