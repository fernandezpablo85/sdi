.PHONY: build build-cli build-api clean

build: build-cli build-api

build-cli:
	go build -o bin/cli ./cmd/cli

build-api:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin/
