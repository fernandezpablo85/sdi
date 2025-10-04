.PHONY: build build-cli build-api test clean

build: build-cli build-api

build-cli:
	go build -o bin/cli ./cmd/cli

build-api:
	go build -o bin/api ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/
