package expense

import (
	"log"
	"net/http"
)

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/expense/create", h.handleCreateExpense)
	mux.HandleFunc("GET /api/expense/list", h.handleGetExpenses)
	mux.HandleFunc("GET /api/expense/project", h.handleGetProjectExpenses)
	mux.HandleFunc("PUT /api/expense/update", h.handleUpdateExpense)
	mux.HandleFunc("DELETE /api/expense/delete", h.handleDeleteExpense)
	mux.HandleFunc("GET /api/expense/health", h.handleHealth)

	log.Println("Expense API routes configured")
}
