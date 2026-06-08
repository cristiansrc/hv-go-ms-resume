package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/service"
)

// PublicHandler handles public read-only endpoints.
type PublicHandler struct {
	publicUseCase   input.PublicUseCase
	pdfUseCase      input.PdfUseCase
	contactUseCase  input.ContactUseCase
	authHandler     *AuthHandler
}

// NewPublicHandler creates a new PublicHandler.
func NewPublicHandler(
	publicUseCase input.PublicUseCase,
	pdfUseCase input.PdfUseCase,
	contactUseCase input.ContactUseCase,
	authHandler *AuthHandler,
) *PublicHandler {
	return &PublicHandler{
		publicUseCase:   publicUseCase,
		pdfUseCase:     pdfUseCase,
		contactUseCase:  contactUseCase,
		authHandler:    authHandler,
	}
}

// TemplatesResponse represents the available PDF templates.
type TemplatesResponse struct {
	Templates []TemplateItem `json:"templates"`
}

// TemplateItem represents a single template option.
type TemplateItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LanguagesResponse represents the available PDF languages.
type LanguagesResponse struct {
	Languages []LanguageItem `json:"languages"`
}

// LanguageItem represents a single language option.
type LanguageItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ContactRequest represents a contact form submission.
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Altcha  string `json:"altcha"`
}

// GetTemplates handles GET /v1/ms-resume/public/templates
func (h *PublicHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	resp := TemplatesResponse{
		Templates: []TemplateItem{
			{ID: "engineeringclassic", Name: "Engineering Classic"},
			{ID: "engineeringresumes", Name: "Engineering Resumes"},
			{ID: "moderncv", Name: "Modern CV"},
			{ID: "sb2nov", Name: "SB2Nov"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetLanguages handles GET /v1/ms-resume/public/languages
func (h *PublicHandler) GetLanguages(w http.ResponseWriter, r *http.Request) {
	resp := LanguagesResponse{
		Languages: []LanguageItem{
			{Code: "english", Name: "English"},
			{Code: "spanish", Name: "Spanish"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetInfoPage handles GET /v1/ms-resume/public/info-page
func (h *PublicHandler) GetInfoPage(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetInfoPage(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get info page")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetCurriculum handles GET /v1/ms-resume/public/curriculum/{language}
func (h *PublicHandler) GetCurriculum(w http.ResponseWriter, r *http.Request) {
	language := chi.URLParam(r, "language")
	template := r.URL.Query().Get("template")
	if template == "" {
		template = "engineeringclassic"
	}

	reader, hash, err := h.pdfUseCase.GetCurriculumPDF(r.Context(), language, template)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLanguage), errors.Is(err, service.ErrInvalidTemplate):
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid language or template parameter")
		case errors.Is(err, service.ErrRenderCVUnavailable):
			WriteError(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "PDF render service is currently unavailable")
		default:
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		}
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=curriculum_"+hash[:8]+".pdf")
	io.Copy(w, reader)
}

// SubmitContact handles POST /v1/ms-resume/public/contact
func (h *PublicHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Validate required fields
	var validationErrors []ApiErrorDetail
	if req.Name == "" {
		validationErrors = append(validationErrors, ApiErrorDetail{Field: "name", Code: "FIELD_REQUIRED", Message: "Name is required"})
	}
	if req.Email == "" {
		validationErrors = append(validationErrors, ApiErrorDetail{Field: "email", Code: "FIELD_REQUIRED", Message: "Email is required"})
	} else if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		validationErrors = append(validationErrors, ApiErrorDetail{Field: "email", Code: "INVALID_EMAIL", Message: "Email must be a valid email address"})
	}
	if req.Message == "" {
		validationErrors = append(validationErrors, ApiErrorDetail{Field: "message", Code: "FIELD_REQUIRED", Message: "Message is required"})
	}
	if len(validationErrors) > 0 {
		WriteValidationError(w, r, validationErrors)
		return
	}

	if err := h.contactUseCase.SubmitContact(r.Context(), req.Name, req.Email, req.Message, req.Altcha); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "The request could not be processed")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetChallenge delegates to auth handler
func (h *PublicHandler) GetChallenge(w http.ResponseWriter, r *http.Request) {
	h.authHandler.GetChallenge(w, r)
}

// GetPublicBlog handles GET /v1/ms-resume/public/blog
func (h *PublicHandler) GetPublicBlog(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	sort := r.URL.Query().Get("sort")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	if size > 100 {
		size = 100
	}
	if sort == "" {
		sort = "id,desc"
	}

	data, err := h.publicUseCase.GetPublicBlogPage(r.Context(), page, size, sort)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get blog page")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicBlogByID handles GET /v1/ms-resume/public/blog/{id}
func (h *PublicHandler) GetPublicBlogByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	data, err := h.publicUseCase.GetPublicBlogByID(r.Context(), id)
	if err != nil {
		WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Blog not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicBlogTypes handles GET /v1/ms-resume/public/blog-type
func (h *PublicHandler) GetPublicBlogTypes(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicBlogTypes(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list blog types")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicBlogTypeByID handles GET /v1/ms-resume/public/blog-type/{id}
func (h *PublicHandler) GetPublicBlogTypeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid ID parameter")
		return
	}
	data, err := h.publicUseCase.GetPublicBlogTypeByID(r.Context(), id)
	if err != nil {
		WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "Blog type not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicCourses handles GET /v1/ms-resume/public/courses
func (h *PublicHandler) GetPublicCourses(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicCourses(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list courses")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicCertifications handles GET /v1/ms-resume/public/certifications
func (h *PublicHandler) GetPublicCertifications(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicCertifications(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list certifications")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicLanguages handles GET /v1/ms-resume/public/languages-data
func (h *PublicHandler) GetPublicLanguages(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicLanguages(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list languages")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicReferences handles GET /v1/ms-resume/public/references
func (h *PublicHandler) GetPublicReferences(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicReferences(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list references")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPublicCustomSections handles GET /v1/ms-resume/public/custom-sections
func (h *PublicHandler) GetPublicCustomSections(w http.ResponseWriter, r *http.Request) {
	data, err := h.publicUseCase.GetPublicCustomSections(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list custom sections")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
