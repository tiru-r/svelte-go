package time

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"svelte-go/internal/shared/types"
)

// Handlers handles HTTP requests for the time module
type Handlers struct {
	service *Service
}

// NewHandlers creates new HTTP handlers
func NewHandlers(service *Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// Request/Response types
type StartTimerRequest struct {
	UserID      string `json:"user_id"`
	ProjectID   string `json:"project_id"`
	Description string `json:"description"`
}

type TimerResponse struct {
	Success  bool             `json:"success"`
	Data     *types.TimeEntry `json:"data,omitempty"`
	Error    string           `json:"error,omitempty"`
	Message  string           `json:"message,omitempty"`
	Duration int64            `json:"current_duration,omitempty"`
}

type UpdateTimerRequest struct {
	UserID      string   `json:"user_id"`
	TimeEntryID string   `json:"time_entry_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// HTTP Handlers

func (h *Handlers) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartTimerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" || req.ProjectID == "" {
		h.sendError(w, "user_id and project_id are required", http.StatusBadRequest)
		return
	}

	entry, err := h.service.StartTimer(req.UserID, req.ProjectID, req.Description)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccess(w, entry, "Timer started successfully")
}

func (h *Handlers) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	entry, err := h.service.StopTimer(userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, entry, "Timer stopped successfully")
}

func (h *Handlers) handlePause(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	entry, err := h.service.PauseTimer(userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, entry, "Timer paused successfully")
}

func (h *Handlers) handleResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	entry, err := h.service.ResumeTimer(userID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, entry, "Timer resumed successfully")
}

func (h *Handlers) handleGetCurrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	entry := h.service.GetActiveTimer(userID)
	if entry == nil {
		h.sendResponse(w, TimerResponse{
			Success: true,
			Message: "No active timer",
			Data:    nil,
		})
		return
	}

	// Get current duration
	duration := h.service.GetCurrentDuration(userID)

	response := TimerResponse{
		Success:  true,
		Data:     entry,
		Duration: duration,
		Message:  "Active timer found",
	}

	h.sendResponse(w, response)
}

func (h *Handlers) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateTimerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.TimeEntryID == "" {
		h.sendError(w, "user_id and time_entry_id are required", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateTimeEntry(req.UserID, req.TimeEntryID, req.Description, req.Tags)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, nil, "Timer updated successfully")
}

func (h *Handlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]any{
		"status":    "healthy",
		"service":   "time-api",
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"database":  "badger",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Helper methods

func (h *Handlers) sendSuccess(w http.ResponseWriter, data *types.TimeEntry, message string) {
	response := TimerResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	h.sendResponse(w, response)
}

func (h *Handlers) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := TimerResponse{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(statusCode)
	h.sendResponse(w, response)
}

func (h *Handlers) sendResponse(w http.ResponseWriter, response TimerResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
