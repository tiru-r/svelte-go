package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"datastar-go/templates"
)

// LoginPage serves the login page
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Login().Render(r.Context(), w)
}

// RegisterPage serves the registration page
func (h *Handlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Register().Render(r.Context(), w)
}

// VerifyAuth verifies the current user's authentication status
func (h *Handlers) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for auth token in cookie or header
	token := ""
	
	// First try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}
	
	// Then try cookie
	if token == "" {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"user": nil,
		})
		return
	}

	// TODO: Verify token with auth service
	// For now, return a mock authenticated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    "user123",
			"name":  "John Doe",
			"email": "john@example.com",
		},
	})
}