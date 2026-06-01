package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler returns the health status of the service.
type HealthHandler struct {
	version   string
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version:   version,
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response body.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ServeHTTP handles GET /health requests.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
