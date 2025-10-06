.PHONY: build build-cli build-api test coverage clean

build: build-cli build-api

build-cli:
	go build -o bin/cli ./cmd/cli

build-api:
	go build -o bin/api ./cmd/api

test:
	gotestsum --format testname

test-no-cache:
	gotestsum --format testname -- -count=1 ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	rm -rf bin/ coverage.out coverage.html
