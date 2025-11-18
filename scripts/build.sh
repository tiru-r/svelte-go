#!/bin/bash
set -e

echo "🔨 Building Svelte + Go Application"
echo "=================================="

# Build frontend
echo "📦 Building frontend..."
cd web
bun install
bun run build
cd ..

# Build backend with SSR frontend (Bun)
echo "🚀 Building backend with SSR frontend (Bun)..."
mkdir -p bin
go build -ldflags="-s -w" -o bin/server main.go

echo "✅ Build complete!"
echo "📁 Binary location: ./bin/server"
echo "🎯 Run with: ./bin/server"
echo "🌐 Access at: http://localhost:8080"