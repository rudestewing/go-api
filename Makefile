# Makefile for Go API Project

.PHONY: help run build test clean migrate-create migrate-up migrate-down migrate-status dev

# Default target
help:
	@echo "Available commands:"
	@echo "  make run              - Run the application with air (hot reload)"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Migration commands:"
	@echo "  make migrate-create name=\"migration_name\"  - Create a new migration"
	@echo "  make migrate-up                            - Run all pending migrations"
	@echo "  make migrate-down                          - Rollback last batch of migrations"
	@echo "  make migrate-status                        - Show migration status"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Start development server with hot reload"

# Application commands
run:
	go run cmd/api/main.go

build:
	go build -o tmp/main cmd/api/main.go

test:
	go test ./...

clean:
	rm -rf tmp/
	go clean

# Development with hot reload
dev:
	air

# Migration commands
migrate-create:
ifndef name
	@echo "Error: name parameter is required"
	@echo "Usage: make migrate-create name=\"migration_name\""
	@exit 1
endif
	go run cmd/migrate/main.go create "$(name)"

migrate-up:
	go run cmd/migrate/main.go migrate

migrate-down:
	go run cmd/migrate/main.go rollback

migrate-status:
	go run cmd/migrate/main.go status

migrate-help:
	go run cmd/migrate/main.go help

# Clean up temporary files
clean-migrations:
	@echo "Cleaning up temporary migration files..."
	@find database/migrations/ -name "*.tmp" -delete 2>/dev/null || true
	@echo "Done."
