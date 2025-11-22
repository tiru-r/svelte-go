#!/bin/bash

# Svelte-Go Application Runner Script
# This script handles building and running the application with proper configuration

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="svelte-go"
DEFAULT_PORT="8080"
DATA_DIR="./data"
BUILD_DIR="./bin"

# Parse command line arguments
CLEAN=false
PORT=${PORT:-$DEFAULT_PORT}
PROD=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        --prod)
            PROD=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -p, --port PORT    Set the port (default: 8080)"
            echo "  -c, --clean        Clean build and data directories before starting"
            echo "  --prod             Run in production mode"
            echo "  -h, --help         Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  PORT              Server port (default: 8080)"
            echo "  JWT_SECRET        JWT secret key (auto-generated if not set)"
            echo ""
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

# Clean function
clean() {
    log_info "Cleaning build and data directories..."
    
    # Clean build directory
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        log_success "Cleaned build directory"
    fi
    
    # Clean data directory (with confirmation)
    if [ -d "$DATA_DIR" ] && [ "$CLEAN" = true ]; then
        log_warning "This will delete all application data!"
        read -p "Are you sure you want to delete the data directory? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$DATA_DIR"
            log_success "Cleaned data directory"
        else
            log_info "Skipping data directory cleanup"
        fi
    fi
    
    # Clean compiled binary
    if [ -f "$APP_NAME" ]; then
        rm "$APP_NAME"
        log_success "Cleaned compiled binary"
    fi
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.19 or later."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    log_success "Go version: $GO_VERSION"
    
    # Check templ (if .templ files exist)
    if find . -name "*.templ" -type f | head -1 | grep -q .; then
        if ! command -v templ &> /dev/null; then
            log_warning "templ command not found. Installing..."
            go install github.com/a-h/templ/cmd/templ@latest
            if ! command -v templ &> /dev/null; then
                log_error "Failed to install templ. Please install manually: go install github.com/a-h/templ/cmd/templ@latest"
                exit 1
            fi
        fi
        log_success "templ is available"
    fi
}

# Generate templates
generate_templates() {
    if find . -name "*.templ" -type f | head -1 | grep -q .; then
        log_info "Generating templates..."
        if templ generate; then
            log_success "Templates generated"
        else
            log_error "Failed to generate templates"
            exit 1
        fi
    fi
}

# Build application
build() {
    log_info "Building application..."
    
    # Create build directory
    mkdir -p "$BUILD_DIR"
    
    # Set build flags
    BUILD_FLAGS="-trimpath"
    if [ "$PROD" = true ]; then
        BUILD_FLAGS="$BUILD_FLAGS -ldflags='-w -s'"
        log_info "Building in production mode (optimized)"
    fi
    
    # Build
    if go build $BUILD_FLAGS -o "$BUILD_DIR/$APP_NAME" ./; then
        log_success "Build completed: $BUILD_DIR/$APP_NAME"
    else
        log_error "Build failed"
        exit 1
    fi
}

# Setup directories
setup_directories() {
    log_info "Setting up directories..."
    
    # Create data directory
    mkdir -p "$DATA_DIR"
    log_success "Data directory ready: $DATA_DIR"
    
    # Create logs directory
    mkdir -p "./logs"
    log_success "Logs directory ready: ./logs"
}

# Set environment variables
setup_environment() {
    log_info "Setting up environment..."
    
    export PORT="$PORT"
    log_success "Port set to: $PORT"
    
    # Generate JWT secret if not provided
    if [ -z "$JWT_SECRET" ]; then
        log_warning "JWT_SECRET not set. A secret will be auto-generated."
        log_warning "Set JWT_SECRET environment variable for production use."
    else
        log_success "JWT_SECRET is configured"
    fi
    
    if [ "$PROD" = true ]; then
        export GIN_MODE="release"
        log_success "Production mode enabled"
    else
        log_info "Development mode (use --prod for production)"
    fi
}

# Start application
start_app() {
    log_info "Starting $APP_NAME on port $PORT..."
    echo ""
    echo -e "${GREEN}🚀 Application starting...${NC}"
    echo -e "${BLUE}📡 Server will be available at: http://localhost:$PORT${NC}"
    echo -e "${BLUE}📊 Health check: http://localhost:$PORT/api/health${NC}"
    echo -e "${BLUE}🔐 Login page: http://localhost:$PORT/login${NC}"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
    echo ""
    
    # Start the application
    exec "$BUILD_DIR/$APP_NAME"
}

# Cleanup function for graceful shutdown
cleanup() {
    log_info "Shutting down..."
    exit 0
}

# Set trap for graceful shutdown
trap cleanup SIGINT SIGTERM

# Main execution
main() {
    echo -e "${GREEN}🏗️  Svelte-Go Application Runner${NC}"
    echo ""
    
    # Clean if requested
    if [ "$CLEAN" = true ]; then
        clean
    fi
    
    # Check dependencies
    check_dependencies
    
    # Generate templates
    generate_templates
    
    # Build application
    build
    
    # Setup directories
    setup_directories
    
    # Setup environment
    setup_environment
    
    # Start application
    start_app
}

# Run main function
main "$@"