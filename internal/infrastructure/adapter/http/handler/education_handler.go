package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
)

// EducationHandler handles HTTP requests for the Education entity.
type EducationHandler struct {
	useCase input.EducationUseCase
}

// NewEducationHandler creates a new EducationHandler.
func NewEducationHandler(useCase input.EducationUseCase) *EducationHandler {
	return &EducationHandler{useCase: useCase}
}

// List handles GET /v1/ms-resume/educations
func (h *EducationHandler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.useCase.List(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list educations")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetByID handles GET /v1/ms-resume/educations/{id}
func (h *EducationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	data, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Education not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Create handles POST /v1/ms-resume/educations
func (h *EducationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.EducationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}
	resp, err := h.useCase.Create(r.Context(), &req)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create education")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update handles PUT /v1/ms-resume/educations/{id}
func (h *EducationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	var req request.EducationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}
	if err := h.useCase.Update(r.Context(), id, &req); err != nil {
		WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Education not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete handles DELETE /v1/ms-resume/educations/{id}
func (h *EducationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	if err := h.useCase.Delete(r.Context(), id); err != nil {
		WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Education not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
