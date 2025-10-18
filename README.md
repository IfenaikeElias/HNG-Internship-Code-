# HNG_stage_0 — Go / Gin demo

Simple HTTP service exposing a single endpoint:

GET /me — returns a JSON profile and a cat fact fetched from https://catfact.ninja/fact.

## Repository contents

- `main.go` — service implementation
- `main_test.go` — unit test for the `/me` endpoint
- `go.mod` — module dependencies

## Prerequisites

- Linux (commands assume a POSIX shell)
- Go 1.18 or newer installed and on PATH

Verify Go is installed:

```bash
go version
```

## Setup

1. Open a terminal and change to the project root (where `go.mod` lives):

```bash
cd /home/elias/Documents/HNG/HNG_stage_0
```

2. Download and verify dependencies:

```bash
go mod tidy
```

This will fetch the dependencies declared in `go.mod`.

## Dependencies

Primary libraries used:

- github.com/gin-gonic/gin — HTTP router
- github.com/go-resty/resty/v2 — HTTP client used to call the cat-fact API
- github.com/stretchr/testify — testing helpers

You do not need to install these manually if you run `go mod tidy`. If you prefer manual install:

```bash
go get github.com/gin-gonic/gin
go get github.com/go-resty/resty/v2
go get github.com/stretchr/testify
```

## Environment variables

None required by the code as provided. The service binds to `localhost:8080` and uses the hard-coded external API URL in `main.go`. If you want configuration via environment variables, update `main.go` accordingly.

## Run locally

Run the service directly:

```bash
go run main.go
```

Or build and run:

```bash
go build -o hng_stage_0 main.go
./hng_stage_0
```

The server logs:

```
Server running on http://localhost:8080
```

Test the endpoint manually:

```bash
curl http://localhost:8080/me
```

Example successful response:

```json
{
  "status": "success",
  "user": {
    "email": "eifenaike@gmail.com",
    "name": "Ifenaike Elias Ayooluwa",
    "stack": "Go/Gin"
  },
  "timestamp": "2025-10-18T12:34:56.789Z",
  "fact": "Cats sleep 70% of their lives."
}
```

## Tests

Unit tests are provided in `main_test.go`. Note: the test calls the running handler which in turn calls the external cat-fact API. Tests therefore require network access and may fail if the external service is unreachable.

Run all tests:

```bash
go test ./... -v
```

Run tests for the current package:

```bash
go test -v
```

Run a single test by name:

```bash
go test -run TestGetProfile -v
```

If you need deterministic tests without external HTTP calls, consider refactoring `fetchCatFact` to accept an injected HTTP client or a function parameter and add a mock in tests.

## Troubleshooting

- If module download issues occur, try:

```bash
GOPROXY=https://proxy.golang.org go mod tidy
```

- If the `/me` endpoint returns an error, check network connectivity to `https://catfact.ninja/fact`.

## Notes

- No environment variables required for the current implementation.
- To change the listen address or external API URL, edit `main.go`.

## License

Add a `LICENSE` file if you plan to publish this repository.