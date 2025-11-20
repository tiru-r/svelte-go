package client

import (
	"log"
	"net/http"
)

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
