all: test build golangci-lint
lint: golangci-lint super-linter
pipeline: all lint

generate:
	go generate ./...

build: generate
	go build -v ./...

test: generate
	go test -v ./...

run:
	go run .

golangci-lint:
	golangci-lint run

super-linter:
	docker run \
		--rm \
		--volume "$(shell pwd):/work:z" \
		--env DEFAULT_WORKSPACE=/work \
		--env RUN_LOCAL=true \
		--env VALIDATE_GO=false \
		ghcr.io/super-linter/super-linter:slim-v5.7.2 bash

.PHONY: all build generate golangci-lint lint pipeline run super-linter test
