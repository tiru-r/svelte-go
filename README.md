# Svelte + Go Application

A PocketBase-like single binary application using Svelte with SvelteKit on the frontend and Go with BadgerDB on the backend.

## Tech Stack

- **Frontend**: Svelte 5 + SvelteKit + TailwindCSS
- **Backend**: Go 1.25 + net/http + BadgerDB
- **Build Tool**: Bun
- **Documentation**: JSDoc
- **Deployment**: Single binary with SSR frontend

## Quick Start

### Prerequisites

- Go 1.25+
- Bun

### Build and Run (Single Binary)

```bash
# Build everything into a single binary
./scripts/build.sh

# Run the server
./bin/server
```

Your application will be available at `http://localhost:8080` with both SSR frontend and API served from the same binary.

### Development Mode

For development with hot reloading:

```bash
# Install dependencies
make install

# Run frontend and backend separately
make dev
```

Or run them separately:

```bash
# Terminal 1: Frontend (with hot reload)
cd web && bun run dev

# Terminal 2: Backend
go run main.go
```

### Build Options

```bash
# Build everything
make build

# Build only frontend
make frontend

# Build only backend (requires frontend to be built first)
make backend

# Clean build artifacts
make clean
```

### Documentation

Generate JSDoc documentation:

```bash
cd web
bunx jsdoc src/lib/utils.js -d docs/
# or for all JS files:
bunx jsdoc -c jsdoc.config.json
```

## Project Structure

```
svelte-go/
├── main.go             # Go application entry point
├── internal/
│   ├── db/             # Database layer (BadgerDB)
│   ├── handlers/       # HTTP handlers
│   └── models/         # Data models
├── web/                # Svelte frontend
│   ├── src/
│   │   └── routes/     # SvelteKit routes
│   └── build/          # Built frontend (embedded in Go binary)
├── scripts/            # Build scripts
└── data/               # BadgerDB data files
```

## API Endpoints

- `GET /api/health` - Health check endpoint (publishes health_check events)
- `POST /api/events` - Publish events to NATS
- `GET /api/nats/stats` - Get NATS server statistics
- `GET /` - Svelte frontend (SSR) with event publishing

## Features

- ✅ Single binary deployment (like PocketBase)
- ✅ Server-side rendering (SSR) with SvelteKit
- ✅ BadgerDB for fast, embedded database
- ✅ **Embedded NATS for event-driven architecture**
- ✅ JetStream for persistent event storage
- ✅ Event publishing and subscription
- ✅ Real-time event statistics
- ✅ TailwindCSS for styling
- ✅ JavaScript with JSDoc type checking
- ✅ Hot reloading in development
- ✅ JSDoc documentation
- ✅ Go proxy to Bun SSR server