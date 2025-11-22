#!/bin/bash

# Quick development script
# This is a simplified version of run.sh for rapid development

set -e

echo "🚀 Quick Development Start"

# Kill any existing process on the port
PORT=${PORT:-3003}
if lsof -ti:$PORT > /dev/null 2>&1; then
    echo "🔄 Stopping existing process on port $PORT..."
    kill -9 $(lsof -ti:$PORT) 2>/dev/null || true
    sleep 1
fi

# Generate templates if needed
if find . -name "*.templ" -type f | head -1 | grep -q . && command -v templ &> /dev/null; then
    echo "📝 Generating templates..."
    templ generate
fi

# Quick build and run
echo "🔨 Building..."
go build -o svelte-go-dev ./

echo "🎯 Starting on port $PORT..."
echo "   http://localhost:$PORT"
echo "   http://localhost:$PORT/login"
echo ""

PORT=$PORT ./svelte-go-dev