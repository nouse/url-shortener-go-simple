# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Start the server (runs on localhost:8080)
go run ./cmd/service

# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run a specific test
go test ./storage -run TestFileStorageShortCode

# Run linter (golangci-lint v2.7.2)
golangci-lint run
```

## Architecture

This is a URL shortener service with a three-layer architecture:

- **cmd/service/main.go**: HTTP server entrypoint with graceful shutdown handling
- **handlers/shortener.go**: HTTP handlers using Go 1.22+ method-based routing on `net/http.ServeMux`
- **storage/storage.go**: Storage interface with file-based implementation

### HTTP Endpoints

| Method | Path         | Description                |
|--------|--------------|----------------------------|
| POST   | /x           | Encode URL to short code   |
| GET    | /x/{code}    | Retrieve URL by short code |
| GET    | /info/{code} | Get URL and visit count    |
| GET    | /ping        | Health check               |

### Key Design Decisions

- **Dependency injection**: Handler receives logger and storage via function parameters
- **Interface-based storage**: `Storage` interface allows swapping implementations
- **Structured logging**: Uses `log/slog` with JSON output
- **Short code generation**: 6-character codes using base32 alphabet, generated with `math/rand/v2`
- **File persistence**: URLs stored in `urls.txt` (created on startup if missing)

## Linting

The project uses extensive linting via `.golangci.toml` with 40+ linters. Key requirements:
- Context required for all slog calls (`sloglint`)
- String keys required for structured logging (`loggercheck`)
- `nolint` comments require explanations and specific linter names
