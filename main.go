package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"svelte-go/internal/db"
	"svelte-go/internal/events"
	"svelte-go/internal/handlers"
)

func main() {
	database, err := db.Init()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	// Initialize NATS embedded server
	natsService, err := events.NewNATSService()
	if err != nil {
		log.Fatal("Failed to initialize NATS:", err)
	}
	defer natsService.Close()

	h := handlers.New(database, natsService)

	// Find available port for SvelteKit server
	sveltePort := findAvailablePort(3001)
	log.Printf("Using port %d for SvelteKit server", sveltePort)

	// Start SvelteKit SSR server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svelteCmd := exec.CommandContext(ctx, "bun", "run", "web/build/index.js")
	svelteCmd.Env = append(os.Environ(),
		"PORT="+fmt.Sprintf("%d", sveltePort),
		"HOST=127.0.0.1")

	// Start SvelteKit server
	if err := svelteCmd.Start(); err != nil {
		log.Fatal("Failed to start SvelteKit server:", err)
	}
	defer func() {
		if svelteCmd.Process != nil {
			svelteCmd.Process.Kill()
		}
	}()

	// Wait for SvelteKit server to be ready with proper health check
	svelteURL := fmt.Sprintf("http://127.0.0.1:%d", sveltePort)
	if !waitForServer(svelteURL, 10*time.Second) {
		log.Fatal("SvelteKit server failed to become ready")
	}

	// Create reverse proxy for SvelteKit
	parsedURL, _ := url.Parse(svelteURL)
	svelteProxy := httputil.NewSingleHostReverseProxy(parsedURL)

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/events", h.Events)
	mux.HandleFunc("/api/nats/stats", h.NATSStats)

	// Proxy everything else to SvelteKit SSR server
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Proxy to SvelteKit
		svelteProxy.ServeHTTP(w, r)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Frontend (SSR) and API available at http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("Shutting down...")

	// Shutdown server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

// findAvailablePort finds an available port starting from the given port
func findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port
		}
	}
	log.Fatal("No available ports found")
	return startPort
}

// waitForServer waits for a server to become ready by making HTTP requests
func waitForServer(url string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return true
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}
