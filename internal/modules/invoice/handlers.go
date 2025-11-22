package invoice

import (
	"encoding/json"
	"net/http"
	"time"

	"datastar-go/internal/shared/types"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
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
		UserID      string             `json:"user_id"`
		ClientID    string             `json:"client_id"`
		ProjectID   string             `json:"project_id"`
		HourlyRate  float64            `json:"hourly_rate"`
		TimeEntries []*types.TimeEntry `json:"time_entries"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.ClientID == "" || req.ProjectID == "" || req.HourlyRate <= 0 || len(req.TimeEntries) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	invoice, err := h.service.GenerateFromTimeEntries(req.UserID, req.ClientID, req.ProjectID, req.HourlyRate, req.TimeEntries)
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
	response := map[string]any{
		"status":    "healthy",
		"module":    "invoice",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
