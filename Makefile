all: test build golangci-lint
lint: golangci-lint super-linter
pipeline: all lint

build:
	go build -v ./...

test:
	go test ./...

golangci-lint:
	golangci-lint run

super-linter:
	docker run \
		--rm \
		--volume "$(shell pwd):/work:z" \
		--env DEFAULT_WORKSPACE=/work \
		--env RUN_LOCAL=true \
		--env VALIDATE_GO=false \
		ghcr.io/super-linter/super-linter:slim-v5.6.1 bash

.PHONY: all build golangci-lint lint pipeline super-linter test
