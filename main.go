package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"svelte-go/internal/events"
	"svelte-go/internal/modules/client"
	"svelte-go/internal/modules/expense"
	"svelte-go/internal/modules/invoice"
	timemodule "svelte-go/internal/modules/time"
	"svelte-go/internal/shared/database"
	"svelte-go/internal/shared/types"
)

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
	
	// Time tracking module
	timeService := timemodule.NewService(eventBus, db)
	timeHandlers := timemodule.NewHandlers(timeService)
	
	// Expense tracking module
	expenseService := expense.NewService(eventBus, db)
	expenseHandlers := expense.NewHandlers(expenseService)
	
	// Client & project management module
	clientService := client.NewService(eventBus, db)
	clientHandlers := client.NewHandlers(clientService)
	
	// Invoice generation module
	invoiceService := invoice.NewService(eventBus, db)
	invoiceHandlers := invoice.NewHandlers(invoiceService)
	
	// Set up HTTP routes
	mux := http.NewServeMux()
	
	// Module routes
	timeHandlers.SetupRoutes(mux)
	expenseHandlers.SetupRoutes(mux)
	clientHandlers.SetupRoutes(mux)
	invoiceHandlers.SetupRoutes(mux)
	
	// System routes
	mux.HandleFunc("/api/health", handleOverallHealth(db, eventBus))
	
	// Serve static frontend files
	staticDir := "./web/build"
	fileServer := http.FileServer(http.Dir(staticDir))
	
	// Handle SPA routing by serving index.html for non-API routes
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve API routes normally
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}
		
		// Try to serve the requested file
		filePath := filepath.Join(staticDir, r.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}
		
		// If file doesn't exist, serve index.html for SPA routing
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}))
	
	// Publish system startup event
	startupEvent := types.NewEvent("system_started", "main", map[string]interface{}{
		"modules": []string{"time", "expense", "client", "invoice"},
		"architecture": "modular_monolith",
		"database": "badger",
		"events": "nats_embedded",
		"start_time": time.Now(),
	})
	eventBus.Publish("system.startup", startupEvent)

	// HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
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
		log.Printf("   GET  /api/health          - Overall health")
		log.Printf("   POST /api/time/start      - Start timer")
		log.Printf("   POST /api/time/stop       - Stop timer")
		log.Printf("   GET  /api/time/current    - Current timer")
		log.Printf("   POST /api/time/pause      - Pause timer")
		log.Printf("   POST /api/time/resume     - Resume timer")
		log.Printf("   PUT  /api/time/update     - Update timer")
		log.Printf("   GET  /health              - Time module health")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("📤 Shutting down...")

	// Publish shutdown event
	shutdownEvent := types.NewEvent("system_shutdown", "main", map[string]interface{}{
		"reason": "signal",
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
		w.WriteHeader(http.StatusOK)
		
		response := `{
			"status": "healthy",
			"architecture": "modular_monolith",
			"database": "badger", 
			"events": "nats_embedded",
			"modules": ["time"],
			"timestamp": "` + time.Now().Format(time.RFC3339) + `"
		}`
		
		w.Write([]byte(response))
	}
}