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