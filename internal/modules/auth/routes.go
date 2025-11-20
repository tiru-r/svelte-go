package auth

import (
	"log"
	"net/http"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/verify", h.VerifyToken)
	mux.HandleFunc("/profile", h.GetProfile)
	mux.HandleFunc("/refresh", h.RefreshToken)

	log.Println("Auth API routes configured")
}
