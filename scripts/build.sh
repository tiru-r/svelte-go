#!/bin/bash
set -e

echo "🚀 PocketBase-like Starter Kit Build"
echo "===================================="
echo ""

# Check prerequisites
echo "🔍 Checking prerequisites..."
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Aborting." >&2; exit 1; }
command -v bun >/dev/null 2>&1 || { echo "❌ Bun is required but not installed. Aborting." >&2; exit 1; }
echo "✅ Prerequisites check passed"
echo ""

# Clean previous build
echo "🧹 Cleaning previous build..."
rm -rf bin/
rm -rf web/build/
echo "✅ Clean complete"
echo ""

# Build frontend
echo "📦 Building frontend..."
cd web
echo "   📥 Installing dependencies..."
bun install --silent
echo "   🏗️  Building Svelte app for SSR..."
bun run build
echo "   ✅ Frontend built to web/build/"
cd ..
echo ""

# Build backend with embedded frontend
echo "🔧 Building Go backend..."
echo "   📁 Creating bin directory..."
mkdir -p bin
echo "   ⚡ Compiling Go binary with optimizations..."
go build -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo 'dev')" -o bin/server main.go
echo "   ✅ Backend built to bin/server"
echo ""

# Get binary size
BINARY_SIZE=$(du -h bin/server | cut -f1)
echo "📊 Build Summary"
echo "================"
echo "   📁 Binary size: $BINARY_SIZE"
echo "   📁 Binary location: ./bin/server"
echo "   🔌 Embedded: Frontend (SSR) + Database + Event System"
echo ""

echo "🎉 Build complete!"
echo ""
echo "🚀 Next steps:"
echo "   1. Run: ./bin/server"
echo "   2. Open: http://localhost:8080"
echo "   3. Deploy: Copy ./bin/server anywhere!"
echo ""
echo "💡 Tips:"
echo "   • Set PORT environment variable for custom port"
echo "   • Binary includes everything - no external dependencies needed"
echo "   • Perfect for single binary deployment like PocketBase"