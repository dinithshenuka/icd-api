.PHONY: all build dev generate test clean help migrate-up migrate-down migrate-create import

# Configuration
MIGRATIONS_DIR=./database/migrations
DB_PATH=./database/icd11.db

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

## migrate-up: Run all database migrations
migrate-up:
	@echo "Running migrations up..."
	go run cmd/migrate/main.go -action up -path $(MIGRATIONS_DIR) -db $(DB_PATH)

## migrate-down: Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	go run cmd/migrate/main.go -action down -path $(MIGRATIONS_DIR) -db $(DB_PATH)

## migrate-create: Create a new migration file (usage: make migrate-create name=my_migration)
migrate-create:
	@echo "Creating new migration..."
	/opt/homebrew/bin/migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

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
