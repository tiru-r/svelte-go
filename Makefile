# Svelte + Go Starter Kit Makefile
# PocketBase-like development commands

.PHONY: install dev build frontend backend clean docs help test lint

# Default target
help: ## Show this help message
	@echo "🚀 Svelte + Go Starter Kit Commands"
	@echo "==================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "📦 Installing dependencies..."
	@cd web && bun install
	@go mod tidy
	@echo "✅ Dependencies installed"

dev: ## Start development servers (frontend + backend)
	@echo "🔥 Starting development servers..."
	@echo "Frontend: http://localhost:5173"
	@echo "Backend:  http://localhost:8080"
	@echo "Press Ctrl+C to stop both servers"
	@trap 'kill %1 %2' INT; \
	cd web && bun run dev & \
	go run main.go & \
	wait

build: frontend backend ## Build everything for production

frontend: ## Build frontend for production
	@echo "🏗️  Building frontend..."
	@cd web && bun run build
	@echo "✅ Frontend built to web/build/"

backend: ## Build backend binary (requires frontend to be built first)
	@echo "🏗️  Building backend..."
	@mkdir -p bin
	@go build -o bin/server main.go
	@echo "✅ Backend built to bin/server"

clean: ## Clean all build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf web/build/
	@rm -rf web/.svelte-kit/
	@rm -rf bin/
	@rm -rf data/
	@echo "✅ Clean complete"

docs: ## Generate JSDoc documentation
	@echo "📚 Generating documentation..."
	@cd web && bunx jsdoc -c jsdoc.config.json
	@echo "✅ Documentation generated to web/docs/"

# Development helpers
dev-frontend: ## Start only frontend dev server
	@cd web && bun run dev

dev-backend: ## Start only backend dev server
	@go run main.go

# Testing
test: ## Run tests
	@echo "🧪 Running tests..."
	@cd web && bun test
	@go test ./...

# Production helpers
prod-build: ## Build production binary using build script
	@./scripts/build.sh

prod-run: ## Run production server
	@./bin/server

start: build ## Build and run production server
	@./bin/server

# Database
db-clean: ## Clean database files
	@rm -rf data/badger/
	@rm -rf data/jetstream/
	@echo "✅ Database files cleaned"

# Linting and formatting
lint: ## Run linters and type checking
	@echo "🔍 Running linters..."
	@cd web && bunx tsc --noEmit
	@go fmt ./...
	@go vet ./...

# Dependencies
deps-update: ## Update dependencies
	@echo "📦 Updating dependencies..."
	@cd web && bun update
	@go get -u ./...
	@go mod tidy

# Quick commands
quick-start: install build start ## Install, build, and run production server

restart: ## Kill processes and restart development
	@pkill -f "bun.*dev" || true
	@pkill -f "go run main.go" || true
	@sleep 1
	@$(MAKE) dev