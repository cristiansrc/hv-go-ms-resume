package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/ctxkeys"
)

// RequestIDMiddleware adds a unique request ID to each request.
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware creates a new RequestIDMiddleware.
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// AttachRequestID is an HTTP middleware that adds X-Request-ID to responses and context.
func (m *RequestIDMiddleware) AttachRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), ctxkeys.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
