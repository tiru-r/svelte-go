# 🎯 PocketBase-like Starter Template

This is a complete starter template for building modern full-stack applications with a single binary deployment model.

## 🏆 What Makes This Special

Unlike traditional full-stack setups, this starter gives you:

### ✅ **Zero External Dependencies**
- No Docker required
- No external database setup
- No message queue installation
- No complex deployment pipelines

### ✅ **Modern Stack**
- **Frontend**: Svelte 5 + SvelteKit + TailwindCSS v4
- **Backend**: Go 1.25 + BadgerDB + NATS
- **Runtime**: Bun (faster than Node.js)
- **Types**: JSDoc (simpler than TypeScript)

### ✅ **Developer Experience**
- Hot reload in development
- Type safety without TypeScript complexity  
- One command build and deploy
- Real-time event system built-in
- Comprehensive documentation

## 🚀 Perfect For

- **Prototypes**: Get ideas running fast
- **Internal Tools**: Admin dashboards, monitoring tools
- **Side Projects**: Personal applications
- **MVPs**: Minimum viable products
- **Learning**: Full-stack development
- **Production Apps**: Scale when ready

## 📦 What's Included

### Frontend (Svelte)
```
web/
├── src/
│   ├── routes/+page.svelte      # Main page with event demo
│   └── lib/utils.js             # API utilities with JSDoc
├── package.json                 # Frontend dependencies
├── tailwind.config.js           # TailwindCSS v4 config
└── svelte.config.js            # SvelteKit SSR config
```

### Backend (Go)
```
internal/
├── db/badger.go                # BadgerDB embedded database
├── events/nats.go              # NATS embedded event system
├── handlers/handlers.go        # HTTP API endpoints
└── models/                     # Data models (add your own)
```

### Build System
```
scripts/build.sh               # Production build script
Makefile                       # Development commands
.env.example                   # Environment configuration
```

## 🎨 Customization Examples

### Add a New API Endpoint

1. **Add handler** (`internal/handlers/handlers.go`):
```go
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
    // Your logic here
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Users API"})
}
```

2. **Register route** (`main.go`):
```go
mux.HandleFunc("/api/users", h.Users)
```

3. **Use in frontend** (`web/src/lib/utils.js`):
```javascript
export async function getUsers() {
    return apiRequest('/api/users');
}
```

### Add a New Database Model

1. **Create model** (`internal/models/user.go`):
```go
type User struct {
    ID       string    `json:"id"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}
```

2. **Add database methods** (`internal/db/users.go`):
```go
func (d *DB) CreateUser(user *User) error {
    // BadgerDB operations
}
```

### Add Event Subscribers

```go
// Subscribe to user events
natsService.Subscribe("user.>", func(event *Event) error {
    log.Printf("User event: %s", event.Type)
    // Send email, update analytics, etc.
    return nil
})
```

## 🔧 Advanced Configuration

### Environment Variables
Copy `.env.example` to `.env` and customize:
```bash
PORT=8080
DB_PATH=./data/badger
NATS_PORT=4223
LOG_LEVEL=INFO
```

### Custom Build Flags
Edit `scripts/build.sh`:
```bash
go build -ldflags="-s -w -X main.version=v1.0.0" -o bin/server main.go
```

### TailwindCSS Customization
Edit `web/tailwind.config.js`:
```javascript
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        brand: '#your-color'
      }
    }
  }
}
```

## 🌟 Deployment Strategies

### 1. Simple VPS Deployment
```bash
# Build locally
./scripts/build.sh

# Copy to server
scp ./bin/server user@yourserver:/opt/myapp/

# Run on server
ssh user@yourserver "cd /opt/myapp && ./server"
```

### 2. Systemd Service
Create `/etc/systemd/system/myapp.service`:
```ini
[Unit]
Description=My App
After=network.target

[Service]
Type=simple
User=myapp
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/server
Environment=PORT=8080
Restart=always

[Install]
WantedBy=multi-user.target
```

### 3. Docker (Optional)
```dockerfile
FROM scratch
COPY bin/server /server
EXPOSE 8080
CMD ["/server"]
```

## 🎓 Learning Path

1. **Start Simple**: Run the default app, explore the code
2. **Add Features**: Create new endpoints and pages
3. **Use Events**: Implement pub/sub patterns
4. **Scale Database**: Add more models and relationships
5. **Deploy**: Use single binary deployment
6. **Monitor**: Add metrics and logging
7. **Optimize**: Profile and improve performance

## 🔄 Upgrade Path

This starter grows with your needs:
- **Small**: Single binary deployment
- **Medium**: Add external services (Redis, PostgreSQL)
- **Large**: Microservices architecture
- **Enterprise**: Kubernetes deployment

## 💡 Pro Tips

- **Development**: Use `make dev` for hot reload
- **Type Safety**: JSDoc provides IntelliSense without TypeScript
- **Events**: Use for real-time features and decoupling
- **Database**: BadgerDB is fast for read-heavy workloads
- **Deployment**: Single binary means no dependency hell
- **Monitoring**: Built-in NATS stats for real-time insights

## 🤝 Community

This template is designed to be:
- **Forkable**: Easy to customize for your needs
- **Educational**: Learn modern full-stack patterns
- **Production-ready**: Scale from prototype to production
- **Maintainable**: Clear structure and documentation

---

**Ready to build something amazing? Start with `make install && make dev`** 🚀