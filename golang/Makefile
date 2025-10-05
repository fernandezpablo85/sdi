.PHONY: build build-cli build-api test clean

build: build-cli build-api

build-cli:
	go build -o bin/cli ./cmd/cli

build-api:
	go build -o bin/api ./cmd/api

test:
	gotestsum --format testname

test-no-cache:
	gotestsum --format testname -- -count=1 ./...

clean:
	rm -rf bin/
