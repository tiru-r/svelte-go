package time

import (
	"log"
	"net/http"
)

// SetupRoutes sets up the HTTP routes for the time module
func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	// Time tracking endpoints
	mux.HandleFunc("/api/time/start", h.handleStart)
	mux.HandleFunc("/api/time/stop", h.handleStop)
	mux.HandleFunc("/api/time/pause", h.handlePause)
	mux.HandleFunc("/api/time/resume", h.handleResume)
	mux.HandleFunc("/api/time/current", h.handleGetCurrent)
	mux.HandleFunc("/api/time/update", h.handleUpdate)

	// Health check
	mux.HandleFunc("/health", h.handleHealth)

	log.Println("Time API routes configured")
}
