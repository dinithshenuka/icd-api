.PHONY: all build dev generate test clean help migrate-up migrate-down migrate-create import

# Configuration
MIGRATIONS_DIR=./database/migrations
DB_PATH=./database/icd11.db
SERVER_MAIN=./cmd/server/main.go
OPENAPI_SPEC=./openapi/openapi.yaml
GEN_DIR=./internal/api/handler

# Default target
all: help

## build: Compile the binary for production
build:
	@echo "Building production binary..."
	go build -o ./bin/api $(SERVER_MAIN)

## dev: Start the development server with hot-reloading (air)
dev:
	@echo "Starting development server with live reload..."
	$(shell go env GOPATH)/bin/air

## generate: Regenerate Go code from OpenAPI specification
generate:
	@echo "Generating API code from $(OPENAPI_SPEC)..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -package handler -generate types,gin $(OPENAPI_SPEC) > $(GEN_DIR)/api.gen.go

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## migrate-up: Run all database migrations
migrate-up:
	@echo "Running migrations up..."
	go run cmd/migrate/main.go -action up -path $(MIGRATIONS_DIR) -db $(DB_PATH)

## migrate-down: Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	go run cmd/migrate/main.go -action down -path $(MIGRATIONS_DIR) -db $(DB_PATH)

## import: Import data from Excel to SQLite
import:
	@echo "Importing data..."
	go run cmd/importer/main.go

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
