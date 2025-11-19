package invoice

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
	mux.HandleFunc("POST /api/invoice/create", h.handleCreateInvoice)
	mux.HandleFunc("POST /api/invoice/generate", h.handleGenerateFromTimeEntries)
	mux.HandleFunc("GET /api/invoice/list", h.handleGetInvoices)
	mux.HandleFunc("GET /api/invoice/client", h.handleGetClientInvoices)
	mux.HandleFunc("PUT /api/invoice/status", h.handleUpdateStatus)
	mux.HandleFunc("DELETE /api/invoice/delete", h.handleDeleteInvoice)
	mux.HandleFunc("GET /api/invoice/health", h.handleHealth)
	
	log.Println("Invoice API routes configured")
}

func (h *Handlers) handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string              `json:"user_id"`
		ClientID  string              `json:"client_id"`
		ProjectID string              `json:"project_id"`
		Items     []types.InvoiceItem `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.ClientID == "" || len(req.Items) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	invoice, err := h.service.CreateInvoice(req.UserID, req.ClientID, req.ProjectID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    invoice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGenerateFromTimeEntries(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"user_id"`
		ClientID  string `json:"client_id"`
		ProjectID string `json:"project_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.ClientID == "" || req.ProjectID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		return
	}

	invoice, err := h.service.GenerateFromTimeEntries(req.UserID, req.ClientID, req.ProjectID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    invoice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetInvoices(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	invoices, err := h.service.GetInvoices(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    invoices,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetClientInvoices(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "client_id required", http.StatusBadRequest)
		return
	}

	invoices, err := h.service.GetInvoicesByClient(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    invoices,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InvoiceID string `json:"invoice_id"`
		Status    string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.InvoiceID == "" || req.Status == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	invoice, err := h.service.UpdateInvoiceStatus(req.InvoiceID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    invoice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleDeleteInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.URL.Query().Get("invoice_id")
	if invoiceID == "" {
		http.Error(w, "invoice_id required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteInvoice(invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Message: "Invoice deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"module":    "invoice",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}