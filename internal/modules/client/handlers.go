package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"svelte-go/internal/shared/types"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/client/create", h.handleCreateClient)
	mux.HandleFunc("GET /api/client/list", h.handleGetClients)
	mux.HandleFunc("PUT /api/client/update", h.handleUpdateClient)
	mux.HandleFunc("POST /api/project/create", h.handleCreateProject)
	mux.HandleFunc("GET /api/project/list", h.handleGetProjects)
	mux.HandleFunc("GET /api/project/client", h.handleGetClientProjects)
	mux.HandleFunc("PUT /api/project/update", h.handleUpdateProject)
	mux.HandleFunc("GET /api/client/health", h.handleHealth)

	log.Println("Client API routes configured")
}

func (h *Handlers) handleCreateClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"user_id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Company string `json:"company"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Name == "" || req.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	client, err := h.service.CreateClient(req.UserID, req.Name, req.Email, req.Company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    client,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetClients(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	clients, err := h.service.GetClients(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    clients,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleUpdateClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string         `json:"client_id"`
		Updates  map[string]any `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ClientID == "" {
		http.Error(w, "client_id required", http.StatusBadRequest)
		return
	}

	client, err := h.service.UpdateClient(req.ClientID, req.Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    client,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID    string  `json:"client_id"`
		UserID      string  `json:"user_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		HourlyRate  float64 `json:"hourly_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ClientID == "" || req.UserID == "" || req.Name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	project, err := h.service.CreateProject(req.ClientID, req.UserID, req.Name, req.Description, req.HourlyRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    project,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	projects, err := h.service.GetProjects(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    projects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetClientProjects(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "client_id required", http.StatusBadRequest)
		return
	}

	projects, err := h.service.GetProjectsByClient(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    projects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID string         `json:"project_id"`
		Updates   map[string]any `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		http.Error(w, "project_id required", http.StatusBadRequest)
		return
	}

	project, err := h.service.UpdateProject(req.ProjectID, req.Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    project,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"status":    "healthy",
		"module":    "client",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
