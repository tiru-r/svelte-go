.PHONY: build dev clean frontend backend

# Build everything for production
build: frontend backend

# Development mode
dev:
	@echo "Starting development servers..."
	@make -j2 dev-frontend dev-backend

dev-frontend:
	cd web && bun run dev

dev-backend:
	go run main.go

# Build frontend
frontend:
	@echo "Building frontend..."
	cd web && bun run build

# Build backend with SSR frontend (Bun)
backend: frontend
	@echo "Building backend with SSR frontend (Bun)..."
	go build -o bin/server main.go

# Clean build artifacts
clean:
	rm -rf web/build
	rm -rf bin
	rm -rf data

# Install dependencies
install:
	go mod download
	cd web && bun install

# Run the built server
start: build
	./bin/server