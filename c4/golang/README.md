# Rate Limiter - Go Implementation

Go implementation for concepts from Chapter 4 of "System Design Interview" by Alex Xu.

## Prerequisites

- Go 1.21+
- `gotestsum` (for testing): `go install gotest.tools/gotestsum@latest`

## Building

```bash
make build          # Build both CLI and API
make build-cli      # Build CLI only
make build-api      # Build API only
```

Binaries are output to `bin/` directory.

## Testing

```bash
make test           # Run tests with gotestsum
make test-no-cache  # Run tests without cache
make coverage       # Generate coverage report (coverage.html)
```

## Running

### API Server

```bash
./bin/api
```

### CLI Tool

```bash
./bin/cli [command]
```

## Design Philosophy

This implementation follows principles from John Ousterhout's "A Philosophy of Software Design":

- **Deep modules**: The rate limiter provides a simple interface while handling complex token bucket logic internally
- **Clarity over cleverness**: Code prioritizes readability and maintainability
- **Good abstractions**: Clean separation between rate limiting logic, storage, and application layers
