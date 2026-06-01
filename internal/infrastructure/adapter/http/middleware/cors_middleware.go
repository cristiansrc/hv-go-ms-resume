package middleware

import (
	"net/http"
)

// CORSMiddleware handles Cross-Origin Resource Sharing.
type CORSMiddleware struct {
	allowedOrigins []string
}

// NewCORSMiddleware creates a new CORSMiddleware.
func NewCORSMiddleware(allowedOrigins ...string) *CORSMiddleware {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	return &CORSMiddleware{allowedOrigins: allowedOrigins}
}

// CORS is an HTTP middleware that sets CORS headers.
func (m *CORSMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
