# Svelte + Go Starter Kit

🚀 **PocketBase-like starter template** - A modern full-stack application with single binary deployment, embedded database, event-driven architecture, and SSR frontend.

## 🏗️ Tech Stack

- **Frontend**: Svelte 5 + SvelteKit + TailwindCSS v4
- **Backend**: Go 1.25 + net/http + BadgerDB v4
- **Events**: Embedded NATS server + JetStream
- **Runtime**: Bun (instead of Node.js)
- **Types**: JavaScript + JSDoc (no TypeScript needed)
- **Deployment**: Single binary with embedded frontend

## ⚡ Quick Start

### Prerequisites

- Go 1.25+
- Bun

### 🚀 Production (Single Binary)

```bash
# Build everything into a single binary
./scripts/build.sh

# Run the server (like PocketBase)
./bin/server

# Server runs on http://localhost:8080
# Frontend (SSR) + API + Database + Events - all in one binary!
```

### 🔥 Development Mode

```bash
# Install dependencies and start development
make install && make dev
```

Or manually:

```bash
# Terminal 1: Frontend (hot reload)
cd web && bun run dev

# Terminal 2: Backend (file watching)
go run main.go
```

## 📁 Project Structure

```
svelte-go/                    # PocketBase-like starter
├── main.go                   # Single binary entry point
├── internal/
│   ├── db/                   # BadgerDB embedded database
│   ├── events/               # NATS embedded event system
│   ├── handlers/             # HTTP API handlers
│   └── models/               # Data models
├── web/                      # Svelte frontend
│   ├── src/
│   │   ├── lib/              # Utilities with JSDoc types
│   │   └── routes/           # SvelteKit SSR routes
│   └── build/                # Built frontend (embedded)
├── scripts/
│   └── build.sh             # Production build script
├── data/                    # Runtime data (BadgerDB + NATS)
└── bin/                     # Built binary
    └── server               # Single executable
```

## 🌐 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Svelte frontend (SSR) |
| `GET` | `/api/health` | Health check + event publishing |
| `POST` | `/api/events` | Publish events to NATS |
| `GET` | `/api/nats/stats` | Real-time NATS statistics |

## ✨ Features

### Core Features
- ✅ **Single binary deployment** (like PocketBase)
- ✅ **Server-side rendering (SSR)** with SvelteKit  
- ✅ **Embedded BadgerDB** - No external database needed
- ✅ **Embedded NATS server** - Event-driven architecture built-in
- ✅ **Hot reloading** in development
- ✅ **Type safety** with JSDoc (no TypeScript complexity)

### Developer Experience
- ✅ **Bun runtime** - Faster than Node.js
- ✅ **TailwindCSS v4** - Latest styling framework
- ✅ **Event system** - Publish/subscribe with JetStream persistence
- ✅ **Real-time stats** - Monitor your application
- ✅ **JSDoc documentation** - Generate docs from code
- ✅ **One-command build** - `./scripts/build.sh`

## 🔧 Development Commands

```bash
# Development
make install          # Install dependencies
make dev             # Start dev servers (frontend + backend)
make clean           # Clean build artifacts

# Building
make build           # Build everything
make frontend        # Build only frontend
make backend         # Build only backend

# Documentation
make docs            # Generate JSDoc documentation

# Production
./scripts/build.sh   # Build single binary
./bin/server         # Run production server
```

## 🎯 Use Cases

Perfect starter for:
- **Admin dashboards** 
- **Internal tools**
- **Prototype applications**
- **Event-driven systems**
- **Real-time applications**
- **Single binary deployments**

## 🔥 Getting Started

1. **Clone this starter:**
   ```bash
   git clone <your-repo> my-app
   cd my-app
   ```

2. **Start developing:**
   ```bash
   make install && make dev
   ```

3. **Build for production:**
   ```bash
   ./scripts/build.sh
   ./bin/server
   ```

4. **Deploy:** Just copy the `./bin/server` binary anywhere!

## 🚀 Production Deployment

The built binary is completely self-contained:
- Embedded frontend (SSR)
- Embedded database (BadgerDB)
- Embedded event system (NATS)
- Zero external dependencies

```bash
# Copy binary to server
scp ./bin/server user@server:/opt/myapp/

# Run on server
./server
# or with custom port:
PORT=8080 ./server
```

## 🔌 Event System

Built-in event-driven architecture:

```javascript
// Frontend: Publish events
await publishEvent('user_action', {
  action: 'button_click',
  timestamp: Date.now()
});
```

```go
// Backend: Subscribe to events
natsService.Subscribe("user.>", func(event *Event) error {
  log.Printf("User action: %s", event.Type)
  return nil
})
```

## 📊 Monitoring

Real-time statistics available at `/api/nats/stats`:
- Active connections
- Message throughput  
- Byte transfer rates
- Server uptime

---

**Built with ❤️ as a modern alternative to traditional full-stack setups**