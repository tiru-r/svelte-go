# Svelte-Go Makefile
# Simple shortcuts for common development tasks

.PHONY: help dev run build clean test lint fmt deps

# Default target
help: ## Show this help message
	@echo "🏗️  Svelte-Go Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "📚 Examples:"
	@echo "  make dev          # Quick development start"
	@echo "  make run PORT=3000 # Production mode on port 3000"
	@echo "  make clean        # Clean all build artifacts"

dev: ## Start development server (quick)
	@echo "🚀 Starting development server..."
	./dev.sh

run: ## Start production server
	@echo "🏭 Starting production server..."
	./run.sh --prod

build: ## Build production binary
	@echo "🔨 Building production binary..."
	./run.sh --prod --port 0 2>/dev/null || true
	@echo "✓ Build completed: bin/svelte-go"

clean: ## Clean build artifacts and data
	@echo "🧹 Cleaning build artifacts..."
	./run.sh --clean --port 0 2>/dev/null || true
	rm -f svelte-go-dev svelte-go
	@echo "✓ Clean completed"

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...
	@if command -v templ >/dev/null 2>&1; then \
		templ fmt .; \
	else \
		echo "⚠️  templ not found for template formatting"; \
	fi

deps: ## Install/update dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@if ! command -v templ >/dev/null 2>&1; then \
		echo "📥 Installing templ..."; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	fi
	@echo "✓ Dependencies ready"

install: deps ## Install dependencies and build
	@echo "⚙️  Setting up project..."
	$(MAKE) deps
	$(MAKE) build
	@echo ""
	@echo "🎉 Setup complete! Run 'make dev' to start developing"

# Environment-specific targets
prod: ## Start production server
	PORT=${PORT:-8080} ./run.sh --prod --port ${PORT}

# Port-specific shortcuts  
dev-3000: ## Start dev server on port 3000
	PORT=3000 ./dev.sh

dev-3003: ## Start dev server on port 3003
	PORT=3003 ./dev.sh

dev-8080: ## Start dev server on port 8080
	PORT=8080 ./dev.sh