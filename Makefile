.PHONY: all build dev generate test clean help

# Default target
all: help

## build: Compile the binary for production
build:
	@echo "Building production binary..."
	go build -o ./bin/api ./cmd/api

## dev: Start the development server with hot-reloading (air)
dev:
	@echo "Starting development server with live reload..."
	$(shell go env GOPATH)/bin/air

## generate: Regenerate Go code from OpenAPI specification
generate:
	@echo "Generating API code from openapi.yaml..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -package icdcodes -generate types,gin api/v1/openapi.yaml > internal/icdcodes/api.gen.go

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## clean: Remove build artifacts and temporary files
clean:
	@echo "Cleaning up..."
	rm -rf ./bin ./tmp

## help: Show available commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
