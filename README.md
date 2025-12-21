# A simplistic URL shortener service

## Components

- cmd/service: Entrypoint of HTTP service.
- handlers: HTTP handlers based on `net/http.ServeMux`.
- storage: Storage URL and code, in JSON.

## Workflow

- Install Go 1.25+
- Start server with `go run ./cmd/service`
- Run tests with `go test -v -coverprofile=cov ./...`
- View test coverage with `go tool cover -html=coverage.out`
- Run linters with `golangci-lint run`
