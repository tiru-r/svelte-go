# 🚀 Getting Started with Svelte + Go Starter Kit

Welcome to your PocketBase-like starter template! This guide will help you get up and running quickly.

## 📋 Prerequisites

Before you start, make sure you have:

- **Go 1.25+** - [Download here](https://golang.org/dl/)
- **Bun** - [Install here](https://bun.sh/docs/installation)

### Quick Installation Check

```bash
go version    # Should show 1.25 or higher
bun --version # Should show any recent version
```

## 🏃 Quick Start (30 seconds)

### Option 1: Production Build (Recommended)

```bash
# 1. Build everything into a single binary
./scripts/build.sh

# 2. Run the server
./bin/server

# 3. Open http://localhost:8080
```

That's it! You have a fully functional application with:
- ✅ Frontend (SSR)
- ✅ API backend
- ✅ Embedded database
- ✅ Event system

### Option 2: Development Mode

```bash
# 1. Install dependencies
make install

# 2. Start development servers
make dev
```

This will start:
- Frontend dev server: `http://localhost:5173` (with hot reload)
- Backend API server: `http://localhost:8080`

## 🎯 What You Get

Your starter includes everything you need:

### Frontend (Svelte 5)
- **Modern UI**: TailwindCSS v4 styling
- **SSR**: Server-side rendering built-in
- **Type Safety**: JavaScript + JSDoc (no TypeScript complexity)
- **Hot Reload**: Instant development feedback

### Backend (Go)
- **Fast**: Native Go performance
- **Embedded Database**: BadgerDB v4 (no external DB needed)
- **Event System**: NATS with JetStream persistence
- **Single Binary**: Deploy anywhere

## 🛠️ Available Commands

```bash
make help           # Show all available commands

# Development
make install        # Install dependencies
make dev           # Start dev servers (frontend + backend)
make clean         # Clean build artifacts

# Building
make build         # Build everything for production
make frontend      # Build only frontend
make backend       # Build only backend

# Testing & Quality
make test          # Run tests
make lint          # Run linters and type checking
make docs          # Generate JSDoc documentation

# Production
./scripts/build.sh # Build optimized single binary
make start         # Build and run production server
```

## 🏗️ Project Structure Overview

```
your-project/
├── 📄 main.go              # Application entry point
├── 🔧 Makefile             # Development commands
├── 📜 scripts/build.sh     # Production build script
├── 
├── 📁 internal/            # Go backend code
│   ├── db/                 # BadgerDB database layer
│   ├── events/             # NATS event system  
│   ├── handlers/           # HTTP API handlers
│   └── models/             # Data models
├── 
├── 📁 web/                 # Svelte frontend
│   ├── src/
│   │   ├── lib/            # Utilities (with JSDoc types)
│   │   └── routes/         # SvelteKit pages/routes
│   ├── build/              # Built frontend (auto-generated)
│   └── package.json        # Frontend dependencies
├── 
├── 📁 bin/                 # Built binaries (auto-generated)
│   └── server              # Your single executable
└── 
└── 📁 data/                # Runtime data (auto-generated)
    ├── badger/             # Database files
    └── jetstream/          # Event stream storage
```

## 🔥 Development Workflow

### 1. Start Development
```bash
make dev
```

### 2. Make Changes
- Edit frontend: `web/src/` (auto-reloads)
- Edit backend: `internal/` (restart `make dev`)

### 3. Test Your Changes
```bash
make test   # Run tests
make lint   # Check types and formatting
```

### 4. Build for Production
```bash
./scripts/build.sh
```

## 🚀 Deployment

Your application builds to a **single binary** with everything embedded:

```bash
# After building
./bin/server
```

### Deploy Anywhere
```bash
# Copy to server
scp ./bin/server user@yourserver:/opt/

# Run on server (no dependencies needed!)
./server

# Or with custom port
PORT=3000 ./server
```

## 🔌 Using the Event System

### Frontend (JavaScript)
```javascript
import { publishEvent } from '$lib/utils.js';

// Publish an event
await publishEvent('user_action', {
  action: 'button_click',
  timestamp: Date.now(),
  userId: 123
});
```

### Backend (Go)
```go
// Subscribe to events
natsService.Subscribe("user.>", func(event *Event) error {
  log.Printf("User event: %s", event.Type)
  // Process event here
  return nil
})

// Publish events
natsService.PublishEvent("events.notification", Event{
  Type: "email_sent",
  Source: "backend",
  Data: map[string]interface{}{
    "to": "user@example.com",
    "subject": "Welcome!",
  },
})
```

## 📊 Monitoring

Visit these endpoints to monitor your application:

- `http://localhost:8080/` - Your frontend
- `http://localhost:8080/api/health` - Health check
- `http://localhost:8080/api/nats/stats` - Event system statistics

## 🎨 Customization

### Frontend Styling
- Edit `web/src/routes/+page.svelte` for the main page
- Modify `web/tailwind.config.js` for styling
- TailwindCSS v4 classes are available everywhere

### Backend API
- Add new endpoints in `internal/handlers/`
- Create data models in `internal/models/`
- Database operations in `internal/db/`

### Build Process
- Modify `scripts/build.sh` for custom build steps
- Update `Makefile` for new development commands

## 🆘 Troubleshooting

### Port Already in Use
```bash
# Check what's using the port
netstat -tlnp | grep :8080

# Kill processes
make restart
```

### Build Issues
```bash
# Clean and rebuild
make clean
./scripts/build.sh
```

### Database Issues
```bash
# Clean database files
make db-clean
```

## 📚 Next Steps

1. **Explore the Code**: Start with `web/src/routes/+page.svelte` and `main.go`
2. **Add Features**: Create new API endpoints and frontend pages
3. **Deploy**: Use the single binary for easy deployment
4. **Scale**: Add more event subscribers and database models

## 💡 Pro Tips

- Use `make dev` for development with hot reload
- JSDoc provides type safety without TypeScript complexity
- The event system supports pub/sub patterns for scalability
- Single binary deployment means no dependency management on servers
- BadgerDB is embedded - no external database setup needed

---

**Happy coding! 🎉**

Need help? Check the main [README.md](./README.md) for more details.