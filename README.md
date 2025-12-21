# A simplistic URL shortener service

## Components

- cmd/service: Entrypoint of HTTP service.
- handlers: HTTP handlers based on `net/http.ServeMux`.
- storage: Storage URL and code, in plaintext.

## Workflow

- Install Go 1.25+
- Start server with `go run ./cmd/service`
- Run tests with `go test ./...`
- TODO: print test coverage
- Run linters with `golangci-lint run`
