package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/middleware"
)

// ApiErrorResponse represents the standard API error response.
type ApiErrorResponse struct {
	Timestamp string          `json:"timestamp"`
	Status    int             `json:"status"`
	Error     string          `json:"error"`
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Path      string          `json:"path"`
	TraceID   string          `json:"trace_id"`
	Details   []ApiErrorDetail `json:"details,omitempty"`
}

// ApiErrorDetail represents a field-level error detail.
type ApiErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteError sends a standardized JSON error response.
func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	traceID := ""
	if rid, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
		traceID = rid
	}

	resp := ApiErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
		Error:     http.StatusText(status),
		Code:      code,
		Message:   message,
		Path:      r.URL.Path,
		TraceID:   traceID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// WriteValidationError sends a standardized validation error response.
func WriteValidationError(w http.ResponseWriter, r *http.Request, details []ApiErrorDetail) {
	traceID := ""
	if rid, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
		traceID = rid
	}

	resp := ApiErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    http.StatusBadRequest,
		Error:     http.StatusText(http.StatusBadRequest),
		Code:      "VALIDATION_ERROR",
		Message:   "The request contains invalid fields.",
		Path:      r.URL.Path,
		TraceID:   traceID,
		Details:   details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}
