package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"datastar-go/internal/events"
	"datastar-go/internal/modules/auth"
	"datastar-go/internal/modules/client"
	"datastar-go/internal/modules/expense"
	"datastar-go/internal/modules/invoice"
	timemodule "datastar-go/internal/modules/time"
	"datastar-go/internal/shared/database"
	"datastar-go/internal/shared/types"
	"datastar-go/internal/web"
)

// EventBusAdapter adapts the existing EventBus to auth module interface
type EventBusAdapter struct {
	bus *events.EventBus
}

func (e *EventBusAdapter) Publish(eventType string, data any) error {
	eventData, ok := data.(map[string]any)
	if !ok {
		eventData = map[string]any{"data": data}
	}
	event := types.NewEvent(eventType, "auth", eventData)
	return e.bus.Publish("AUTH_EVENTS", event)
}

func main() {
	// Initialize shared database
	db, err := database.NewBadgerDB("./data/freelancer_db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize event bus
	eventBus, err := events.NewEventBus()
	if err != nil {
		log.Fatal("Failed to initialize event bus:", err)
	}
	defer eventBus.Close()

	// Initialize modules
	log.Println("🏗️  Initializing modular monolith...")

	// Authentication module
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = auth.GenerateSecretKey()
		log.Println("⚠️  Using generated JWT secret. Set JWT_SECRET env var for production")
	}
	authRepo := auth.NewRepository(db.DB())
	// Create event adapter
	eventAdapter := &EventBusAdapter{bus: eventBus}
	authService := auth.NewService(authRepo, jwtSecret, eventAdapter)
	authHandler := auth.NewHandler(authService)
	authMiddleware := auth.NewMiddleware(authService)

	// Time tracking module
	timeService := timemodule.NewService(eventBus, db.DB())
	timeHandlers := timemodule.NewHandlers(timeService)

	// Expense tracking module
	expenseService := expense.NewService(eventBus, db.DB())
	expenseHandlers := expense.NewHandlers(expenseService)

	// Client & project management module
	clientService := client.NewService(eventBus, db.DB())
	clientHandlers := client.NewHandlers(clientService)

	// Invoice generation module
	invoiceService := invoice.NewService(eventBus, db.DB())
	invoiceHandlers := invoice.NewHandlers(invoiceService)

	// Web handlers for Templ/Datastar frontend
	webHandlers := web.NewHandlers(clientService, expenseService, invoiceService, timeService)

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Auth routes (no middleware)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/verify", authHandler.VerifyToken)
	mux.HandleFunc("/api/auth/profile", authHandler.GetProfile)
	mux.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)

	// Module routes (protected) - wrap existing mux with auth middleware
	protectedMux := http.NewServeMux()
	timeHandlers.SetupRoutes(protectedMux)
	expenseHandlers.SetupRoutes(protectedMux)
	clientHandlers.SetupRoutes(protectedMux)
	invoiceHandlers.SetupRoutes(protectedMux)

	// Wrap all non-auth API routes with authentication
	mux.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for specific endpoints
		if r.URL.Path == "/health" {
			protectedMux.ServeHTTP(w, r)
			return
		}
		authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			protectedMux.ServeHTTP(w, r)
		})(w, r)
	})))

	// System routes
	mux.HandleFunc("/api/health", handleOverallHealth(db, eventBus))

	// Set up web routes 
	webMux := http.NewServeMux()
	webHandlers.SetupRoutes(webMux)
	
	// Add login/register pages (no auth required)
	mux.HandleFunc("/login", webHandlers.LoginPage)
	mux.HandleFunc("/register", webHandlers.RegisterPage)
	
	// Wrap protected web routes with auth middleware
	mux.Handle("/", authMiddleware.RequireWebAuth(func(w http.ResponseWriter, r *http.Request) {
		// Serve API routes normally
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		// Skip auth for login/register pages  
		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			webMux.ServeHTTP(w, r)
			return
		}

		// Serve protected web pages
		webMux.ServeHTTP(w, r)
	}))

	// Publish system startup event
	startupEvent := types.NewEvent("system_started", "main", map[string]any{
		"modules":      []string{"auth", "time", "expense", "client", "invoice"},
		"architecture": "modular_monolith",
		"database":     "badger",
		"events":       "nats_embedded",
		"start_time":   time.Now(),
	})
	eventBus.Publish("system.startup", startupEvent)

	// HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Freelancer app starting on port %s", port)
		log.Printf("📡 Event-driven modular monolith")
		log.Printf("🗄️  Database: Badger (%s)", "./data/freelancer_db")
		log.Printf("📊 Available endpoints:")
		log.Printf("   POST /api/auth/register   - Register user")
		log.Printf("   POST /api/auth/login      - Login user")
		log.Printf("   POST /api/auth/logout     - Logout user")
		log.Printf("   GET  /api/health          - Overall health")
		log.Printf("   POST /api/time/start      - Start timer (protected)")
		log.Printf("   POST /api/time/stop       - Stop timer (protected)")
		log.Printf("   GET  /api/time/current    - Current timer (protected)")
		log.Printf("   And more protected endpoints...")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("📤 Shutting down...")

	// Publish shutdown event
	shutdownEvent := types.NewEvent("system_shutdown", "main", map[string]any{
		"reason":         "signal",
		"uptime_seconds": time.Now().Unix(),
	})
	eventBus.Publish("system.shutdown", shutdownEvent)

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Freelancer app stopped gracefully")
}

// handleOverallHealth returns overall system health
func handleOverallHealth(db *database.BadgerDB, eventBus *events.EventBus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database health
		dbStatus := "healthy"
		if db == nil {
			dbStatus = "unhealthy"
		}

		// Check event bus health
		eventsStatus := "healthy"
		if eventBus == nil {
			eventsStatus = "unhealthy"
		}

		status := "healthy"
		if dbStatus != "healthy" || eventsStatus != "healthy" {
			status = "unhealthy"
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		response := `{
			"status": "` + status + `",
			"architecture": "modular_monolith",
			"database": {"status": "` + dbStatus + `", "type": "badger"}, 
			"events": {"status": "` + eventsStatus + `", "type": "nats_embedded"},
			"modules": ["auth", "time", "expense", "client", "invoice"],
			"timestamp": "` + time.Now().Format(time.RFC3339) + `"
		}`

		w.Write([]byte(response))
	}
}
