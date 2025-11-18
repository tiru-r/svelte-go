package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"svelte-go/internal/db"
	"svelte-go/internal/events"
)

type Handler struct {
	db   *db.DB
	nats *events.NATSService
}

func New(database *db.DB, natsService *events.NATSService) *Handler {
	h := &Handler{
		db:   database,
		nats: natsService,
	}
	
	// Setup event subscriptions
	h.setupEventHandlers()
	
	return h
}

func (h *Handler) setupEventHandlers() {
	// Subscribe to health check events
	h.nats.Subscribe("events.health.*", func(event *events.Event) error {
		// Handle health check events
		if event.Type == "health_check" {
			// Could log, update metrics, etc.
			return nil
		}
		return nil
	})
	
	// Subscribe to user action events
	h.nats.Subscribe("user.action.*", func(event *events.Event) error {
		// Handle user actions - could update analytics, trigger workflows, etc.
		return nil
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Publish health check event
	h.nats.PublishEvent("events.health.check", events.Event{
		Type:   "health_check",
		Source: "api",
		Data: map[string]interface{}{
			"endpoint": "/api/health",
			"timestamp": time.Now(),
		},
	})

	response := map[string]interface{}{
		"status":    "healthy",
		"database":  "connected",
		"nats":      "running",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var eventData struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&eventData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Publish user action event
	subject := "user.action." + eventData.Type
	event := events.Event{
		Type:   eventData.Type,
		Source: "frontend",
		Data:   eventData.Data,
	}

	if err := h.nats.PublishEvent(subject, event); err != nil {
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Event published successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) NATSStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := h.nats.GetStats()
	if stats == nil {
		http.Error(w, "NATS stats not available", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"connections": stats["connections"],
		"subscriptions": stats["subscriptions"],
		"messages": map[string]interface{}{
			"in":  stats["in_msgs"],
			"out": stats["out_msgs"],
		},
		"bytes": map[string]interface{}{
			"in":  stats["in_bytes"],
			"out": stats["out_bytes"],
		},
		"uptime": stats["uptime"].(time.Duration).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}