package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// SkillSonHandler handles HTTP requests for the SkillSon entity.
type SkillSonHandler struct {
	useCase input.SkillSonUseCase
}

// NewSkillSonHandler creates a new SkillSonHandler.
func NewSkillSonHandler(useCase input.SkillSonUseCase) *SkillSonHandler {
	return &SkillSonHandler{useCase: useCase}
}

// List handles GET /v1/ms-resume/skill-sons
func (h *SkillSonHandler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.useCase.List(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list skill sons")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetByID handles GET /v1/ms-resume/skill-sons/{id}
func (h *SkillSonHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	data, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Skill son not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Create handles POST /v1/ms-resume/skill-sons
func (h *SkillSonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.SkillSonRequest
	if !decodeAndValidate(w, r, &req) {
		return
	}
	resp, err := h.useCase.Create(r.Context(), &req)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create skill son")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Update handles PUT /v1/ms-resume/skill-sons/{id}
func (h *SkillSonHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	var req request.SkillSonRequest
	if !decodeAndValidate(w, r, &req) {
		return
	}
	if err := h.useCase.Update(r.Context(), id, &req); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Skill son not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update skill son")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete handles DELETE /v1/ms-resume/skill-sons/{id}
func (h *SkillSonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Skill son not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete skill son")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
