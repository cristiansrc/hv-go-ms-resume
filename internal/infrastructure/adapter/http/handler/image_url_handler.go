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

// ImageUrlHandler handles HTTP requests for the ImageUrl entity.
type ImageUrlHandler struct {
	useCase input.ImageUrlUseCase
}

// NewImageUrlHandler creates a new ImageUrlHandler.
func NewImageUrlHandler(useCase input.ImageUrlUseCase) *ImageUrlHandler {
	return &ImageUrlHandler{useCase: useCase}
}

// List handles GET /v1/ms-resume/image-urls
func (h *ImageUrlHandler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.useCase.List(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetByID handles GET /v1/ms-resume/image-urls/{id}
func (h *ImageUrlHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	data, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Image URL not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Create handles POST /v1/ms-resume/image-urls
func (h *ImageUrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.ImageURLRequest
	if !decodeAndValidate(w, r, &req) {
		return
	}

	created, err := h.useCase.Create(r.Context(), &req)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create image URL")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// Delete handles DELETE /v1/ms-resume/image-urls/{id}
func (h *ImageUrlHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Image URL not found")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete image URL")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
