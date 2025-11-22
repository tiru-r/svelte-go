package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserKey     contextKey = "user"
	ClaimsKey   contextKey = "claims"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(service *Service) *Middleware {
	return &Middleware{service: service}
}

func (m *Middleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.service.VerifyToken(tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := m.service.GetUserByID(claims.UserID)
		if err != nil {
			log.Printf("User not found: %v", err)
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			http.Error(w, "Account deactivated", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserKey, user)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			if tokenString != "" {
				claims, err := m.service.VerifyToken(tokenString)
				if err == nil {
					user, err := m.service.GetUserByID(claims.UserID)
					if err == nil && user.IsActive {
						ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
						ctx = context.WithValue(ctx, UserKey, user)
						ctx = context.WithValue(ctx, ClaimsKey, claims)

						r.Header.Set("X-User-ID", claims.UserID)
						r.Header.Set("X-Username", claims.Username)
						r.Header.Set("X-Email", claims.Email)

						r = r.WithContext(ctx)
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	}
}

func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return r.Header.Get("X-User-ID")
}

func GetUser(r *http.Request) *User {
	if user, ok := r.Context().Value(UserKey).(*User); ok {
		return user
	}
	return nil
}

func GetClaims(r *http.Request) *AuthClaims {
	if claims, ok := r.Context().Value(ClaimsKey).(*AuthClaims); ok {
		return claims
	}
	return nil
}

// RequireWebAuth handles auth for web pages (supports both cookies and Authorization header)
func (m *Middleware) RequireWebAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		
		// Try Authorization header first (for API calls)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
		} else {
			// Try cookie (for web pages)
			if cookie, err := r.Cookie("auth_token"); err == nil {
				tokenString = cookie.Value
			}
		}
		
		if tokenString == "" {
			// Redirect to login for web requests, return 401 for API requests
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}

		claims, err := m.service.VerifyToken(tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false, // Set to true in production with HTTPS
				SameSite: http.SameSiteStrictMode,
			})
			
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}

		user, err := m.service.GetUserByID(claims.UserID)
		if err != nil {
			log.Printf("User not found: %v", err)
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}

		if !user.IsActive {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Account deactivated", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserKey, user)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}