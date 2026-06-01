package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
)

// HomeHandler handles HTTP requests for the Home entity.
type HomeHandler struct {
	useCase input.HomeUseCase
}

// NewHomeHandler creates a new HomeHandler.
func NewHomeHandler(useCase input.HomeUseCase) *HomeHandler {
	return &HomeHandler{useCase: useCase}
}

// GetByID handles GET /v1/ms-resume/home/{id}
func (h *HomeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	data, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Home not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Update handles PUT /v1/ms-resume/home/{id}
func (h *HomeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	var req request.HomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.useCase.Update(r.Context(), id, &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Home not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update home")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
