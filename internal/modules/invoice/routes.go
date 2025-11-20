package invoice

import (
	"log"
	"net/http"
)

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
