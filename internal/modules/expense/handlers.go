package expense

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

func (h *Handlers) handleCreateExpense(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string  `json:"user_id"`
		ProjectID   string  `json:"project_id"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Description == "" || req.Amount <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	expense, err := h.service.CreateExpense(req.UserID, req.ProjectID, req.Category, req.Description, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    expense,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	expenses, err := h.service.GetExpenses(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    expenses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleGetProjectExpenses(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		http.Error(w, "project_id required", http.StatusBadRequest)
		return
	}

	expenses, err := h.service.GetExpensesByProject(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    expenses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleUpdateExpense(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ExpenseID string         `json:"expense_id"`
		Updates   map[string]any `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ExpenseID == "" {
		http.Error(w, "expense_id required", http.StatusBadRequest)
		return
	}

	expense, err := h.service.UpdateExpense(req.ExpenseID, req.Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Data:    expense,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("expense_id")
	if expenseID == "" {
		http.Error(w, "expense_id required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteExpense(expenseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Success: true,
		Message: "Expense deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"status":    "healthy",
		"module":    "expense",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
