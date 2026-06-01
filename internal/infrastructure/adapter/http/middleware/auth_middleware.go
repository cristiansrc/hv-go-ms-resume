package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	// ClaimsContextKey is the key for storing JWT claims in the request context.
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware validates JWT tokens on protected routes.
type AuthMiddleware struct {
	jwtService *JWTService
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(jwtService *JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate is an HTTP middleware that validates JWT Bearer tokens.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, r, "Missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w, r, "Invalid authorization header format")
			return
		}

		portClaims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			writeUnauthorized(w, r, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, portClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	requestID := ""
	if rid, ok := r.Context().Value(RequestIDKey).(string); ok {
		requestID = rid
	}

	resp := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"status":    http.StatusUnauthorized,
		"error":     "Unauthorized",
		"code":      "UNAUTHORIZED",
		"message":   message,
		"path":      r.URL.Path,
		"trace_id":  requestID,
		"details":   nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(resp)
}
