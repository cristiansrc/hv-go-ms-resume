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

// BasicDataHandler handles HTTP requests for the BasicData entity.
type BasicDataHandler struct {
	useCase input.BasicDataUseCase
}

// NewBasicDataHandler creates a new BasicDataHandler.
func NewBasicDataHandler(useCase input.BasicDataUseCase) *BasicDataHandler {
	return &BasicDataHandler{useCase: useCase}
}

// GetByID handles GET /v1/ms-resume/basic-data/{id}
func (h *BasicDataHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	data, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Basic data not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Update handles PUT /v1/ms-resume/basic-data/{id}
func (h *BasicDataHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	var req request.BasicDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.useCase.Update(r.Context(), id, &req); err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update basic data")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
