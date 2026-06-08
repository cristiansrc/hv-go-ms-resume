package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type mockBasicDataUseCase struct {
	input.BasicDataUseCase
	getByIDFn func(ctx context.Context, id int64) (*response.BasicDataResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.BasicDataRequest) error
}

func (m *mockBasicDataUseCase) GetByID(ctx context.Context, id int64) (*response.BasicDataResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockBasicDataUseCase) Update(ctx context.Context, id int64, req *request.BasicDataRequest) error {
	return m.updateFn(ctx, id, req)
}

type mockLabelUseCase struct {
	input.LabelUseCase
	listFn   func(ctx context.Context) ([]response.LabelResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.LabelResponse, error)
	createFn  func(ctx context.Context, req *request.LabelRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.LabelRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockLabelUseCase) List(ctx context.Context) ([]response.LabelResponse, error) {
	return m.listFn(ctx)
}
func (m *mockLabelUseCase) GetByID(ctx context.Context, id int64) (*response.LabelResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockLabelUseCase) Create(ctx context.Context, req *request.LabelRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockLabelUseCase) Update(ctx context.Context, id int64, req *request.LabelRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockLabelUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

// ---------------------------------------------------------------------------
// Health Handler Tests
// ---------------------------------------------------------------------------

func TestHealthHandler_ServeHTTP(t *testing.T) {
	h := NewHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", resp.Status)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", resp.Version)
	}
	if resp.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

// ---------------------------------------------------------------------------
// BasicData Handler Tests
// ---------------------------------------------------------------------------

func TestBasicDataHandler_GetByID_Success(t *testing.T) {
	expected := &response.BasicDataResponse{
		ID:        1,
		FirstName: "John",
		Email:     "john@example.com",
	}
	useCase := &mockBasicDataUseCase{
		getByIDFn: func(ctx context.Context, id int64) (*response.BasicDataResponse, error) {
			return expected, nil
		},
	}
	h := NewBasicDataHandler(useCase)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/basic-data/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/basic-data/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp response.BasicDataResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.FirstName != expected.FirstName {
		t.Errorf("FirstName: got %s, want %s", resp.FirstName, expected.FirstName)
	}
}

func TestBasicDataHandler_GetByID_InvalidID(t *testing.T) {
	h := NewBasicDataHandler(&mockBasicDataUseCase{})

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/basic-data/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/basic-data/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp ApiErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Code != "INVALID_ID" {
		t.Errorf("expected code INVALID_ID, got %s", errResp.Code)
	}
}

func TestBasicDataHandler_GetByID_NotFound(t *testing.T) {
	useCase := &mockBasicDataUseCase{
		getByIDFn: func(ctx context.Context, id int64) (*response.BasicDataResponse, error) {
			return nil, fmt.Errorf("basic_data not found: %w", entity.ErrNotFound)
		},
	}
	h := NewBasicDataHandler(useCase)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/basic-data/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/basic-data/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestBasicDataHandler_Update_Success(t *testing.T) {
	useCase := &mockBasicDataUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.BasicDataRequest) error {
			return nil
		},
	}
	h := NewBasicDataHandler(useCase)

	r := chi.NewRouter()
	r.Put("/v1/ms-resume/basic-data/{id}", h.Update)

	body := `{"firstName":"John","firstSurName":"Doe","dateBirth":"1990-01-01","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/basic-data/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestBasicDataHandler_Update_InvalidJSON(t *testing.T) {
	h := NewBasicDataHandler(&mockBasicDataUseCase{})

	r := chi.NewRouter()
	r.Put("/v1/ms-resume/basic-data/{id}", h.Update)

	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/basic-data/1", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Label Handler Tests (full CRUD)
// ---------------------------------------------------------------------------

func TestLabelHandler_List_Success(t *testing.T) {
	expected := []response.LabelResponse{
		{ID: 1, Name: "Backend", NameEng: "Backend"},
		{ID: 2, Name: "Frontend", NameEng: "Frontend"},
	}
	useCase := &mockLabelUseCase{
		listFn: func(ctx context.Context) ([]response.LabelResponse, error) {
			return expected, nil
		},
	}
	h := NewLabelHandler(useCase)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/label", h.List)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/label", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []response.LabelResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 items, got %d", len(resp))
	}
}

func TestLabelHandler_GetByID_Success(t *testing.T) {
	expected := &response.LabelResponse{ID: 1, Name: "Backend", NameEng: "Backend"}
	useCase := &mockLabelUseCase{
		getByIDFn: func(ctx context.Context, id int64) (*response.LabelResponse, error) {
			return expected, nil
		},
	}
	h := NewLabelHandler(useCase)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/label/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/label/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp response.LabelResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != expected.Name {
		t.Errorf("Name: got %s, want %s", resp.Name, expected.Name)
	}
}

func TestLabelHandler_GetByID_InvalidID(t *testing.T) {
	h := NewLabelHandler(&mockLabelUseCase{})
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/label/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/label/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestLabelHandler_Create_Success(t *testing.T) {
	useCase := &mockLabelUseCase{
		createFn: func(ctx context.Context, req *request.LabelRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewLabelHandler(useCase)

	r := chi.NewRouter()
	r.Post("/v1/ms-resume/label", h.Create)

	body := `{"name":"Backend","nameEng":"Backend"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/label", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp response.CreateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != 1 {
		t.Errorf("ID: got %d, want 1", resp.ID)
	}
}

func TestLabelHandler_Create_InvalidJSON(t *testing.T) {
	h := NewLabelHandler(&mockLabelUseCase{})
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/label", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/label", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestLabelHandler_Update_Success(t *testing.T) {
	useCase := &mockLabelUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.LabelRequest) error {
			return nil
		},
	}
	h := NewLabelHandler(useCase)

	r := chi.NewRouter()
	r.Put("/v1/ms-resume/label/{id}", h.Update)

	body := `{"name":"Backend","nameEng":"Backend"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/label/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestLabelHandler_Delete_Success(t *testing.T) {
	useCase := &mockLabelUseCase{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	h := NewLabelHandler(useCase)

	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/label/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/label/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestLabelHandler_Delete_InvalidID(t *testing.T) {
	h := NewLabelHandler(&mockLabelUseCase{})
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/label/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/label/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Mocks for Course, Certification, Language, Reference, CustomSection
// ---------------------------------------------------------------------------

type mockCourseUseCase struct {
	input.CourseUseCase
	listFn    func(ctx context.Context) ([]response.CourseResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.CourseResponse, error)
	createFn  func(ctx context.Context, req *request.CourseRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.CourseRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockCourseUseCase) List(ctx context.Context) ([]response.CourseResponse, error) {
	return m.listFn(ctx)
}
func (m *mockCourseUseCase) GetByID(ctx context.Context, id int64) (*response.CourseResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCourseUseCase) Create(ctx context.Context, req *request.CourseRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockCourseUseCase) Update(ctx context.Context, id int64, req *request.CourseRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockCourseUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type mockCertificationUseCase struct {
	input.CertificationUseCase
	listFn    func(ctx context.Context) ([]response.CertificationResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.CertificationResponse, error)
	createFn  func(ctx context.Context, req *request.CertificationRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.CertificationRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockCertificationUseCase) List(ctx context.Context) ([]response.CertificationResponse, error) {
	return m.listFn(ctx)
}
func (m *mockCertificationUseCase) GetByID(ctx context.Context, id int64) (*response.CertificationResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCertificationUseCase) Create(ctx context.Context, req *request.CertificationRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockCertificationUseCase) Update(ctx context.Context, id int64, req *request.CertificationRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockCertificationUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type mockLanguageUseCase struct {
	input.LanguageUseCase
	listFn    func(ctx context.Context) ([]response.LanguageResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.LanguageResponse, error)
	createFn  func(ctx context.Context, req *request.LanguageRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.LanguageRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockLanguageUseCase) List(ctx context.Context) ([]response.LanguageResponse, error) {
	return m.listFn(ctx)
}
func (m *mockLanguageUseCase) GetByID(ctx context.Context, id int64) (*response.LanguageResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockLanguageUseCase) Create(ctx context.Context, req *request.LanguageRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockLanguageUseCase) Update(ctx context.Context, id int64, req *request.LanguageRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockLanguageUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type mockReferenceUseCase struct {
	input.ReferenceUseCase
	listFn    func(ctx context.Context) ([]response.ReferenceResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.ReferenceResponse, error)
	createFn  func(ctx context.Context, req *request.ReferenceRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.ReferenceRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockReferenceUseCase) List(ctx context.Context) ([]response.ReferenceResponse, error) {
	return m.listFn(ctx)
}
func (m *mockReferenceUseCase) GetByID(ctx context.Context, id int64) (*response.ReferenceResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockReferenceUseCase) Create(ctx context.Context, req *request.ReferenceRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockReferenceUseCase) Update(ctx context.Context, id int64, req *request.ReferenceRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockReferenceUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type mockCustomSectionUseCase struct {
	input.CustomSectionUseCase
	listFn    func(ctx context.Context) ([]response.CustomSectionResponse, error)
	getByIDFn func(ctx context.Context, id int64) (*response.CustomSectionResponse, error)
	createFn  func(ctx context.Context, req *request.CustomSectionRequest) (*response.CreateResponse, error)
	updateFn  func(ctx context.Context, id int64, req *request.CustomSectionRequest) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockCustomSectionUseCase) List(ctx context.Context) ([]response.CustomSectionResponse, error) {
	return m.listFn(ctx)
}
func (m *mockCustomSectionUseCase) GetByID(ctx context.Context, id int64) (*response.CustomSectionResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCustomSectionUseCase) Create(ctx context.Context, req *request.CustomSectionRequest) (*response.CreateResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockCustomSectionUseCase) Update(ctx context.Context, id int64, req *request.CustomSectionRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockCustomSectionUseCase) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

// ---------------------------------------------------------------------------
// Course Handler Tests
// ---------------------------------------------------------------------------

func TestCourseHandler_List_Success(t *testing.T) {
	expected := []response.CourseResponse{
		{ID: 1, Name: "Go Basics", Institution: "Udemy", CompletionDate: "2024-01-01"},
	}
	useCase := &mockCourseUseCase{
		listFn: func(ctx context.Context) ([]response.CourseResponse, error) {
			return expected, nil
		},
	}
	h := NewCourseHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/course", h.List)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/course", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var resp []response.CourseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp))
	}
}

func TestCourseHandler_GetByID_Success(t *testing.T) {
	expected := &response.CourseResponse{ID: 1, Name: "Go Basics"}
	useCase := &mockCourseUseCase{
		getByIDFn: func(ctx context.Context, id int64) (*response.CourseResponse, error) {
			return expected, nil
		},
	}
	h := NewCourseHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/course/{id}", h.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/course/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCourseHandler_Create_Success(t *testing.T) {
	useCase := &mockCourseUseCase{
		createFn: func(ctx context.Context, req *request.CourseRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewCourseHandler(useCase)
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/course", h.Create)
	body := `{"name":"Go Basics","nameEng":"Go Basics","institution":"Udemy","institutionEng":"Udemy","completionDate":"2024-01-01"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/course", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestCourseHandler_Create_InvalidJSON(t *testing.T) {
	h := NewCourseHandler(&mockCourseUseCase{})
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/course", h.Create)
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/course", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCourseHandler_Update_Success(t *testing.T) {
	useCase := &mockCourseUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.CourseRequest) error {
			return nil
		},
	}
	h := NewCourseHandler(useCase)
	r := chi.NewRouter()
	r.Put("/v1/ms-resume/course/{id}", h.Update)
	body := `{"name":"Go Basics","nameEng":"Go Basics","institution":"Udemy","institutionEng":"Udemy","completionDate":"2024-01-01"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/course/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestCourseHandler_Delete_Success(t *testing.T) {
	useCase := &mockCourseUseCase{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	h := NewCourseHandler(useCase)
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/course/{id}", h.Delete)
	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/course/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Certification Handler Tests
// ---------------------------------------------------------------------------

func TestCertificationHandler_List_Success(t *testing.T) {
	expected := []response.CertificationResponse{
		{ID: 1, Name: "AWS Certified", IssuingOrganization: "Amazon", IssueDate: "2024-06-01"},
	}
	useCase := &mockCertificationUseCase{
		listFn: func(ctx context.Context) ([]response.CertificationResponse, error) {
			return expected, nil
		},
	}
	h := NewCertificationHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/certification", h.List)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/certification", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCertificationHandler_Create_Success(t *testing.T) {
	useCase := &mockCertificationUseCase{
		createFn: func(ctx context.Context, req *request.CertificationRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewCertificationHandler(useCase)
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/certification", h.Create)
	body := `{"name":"AWS Certified","nameEng":"AWS Certified","issuingOrganization":"Amazon","issuingOrganizationEng":"Amazon","issueDate":"2024-06-01"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/certification", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestCertificationHandler_Update_Success(t *testing.T) {
	useCase := &mockCertificationUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.CertificationRequest) error {
			return nil
		},
	}
	h := NewCertificationHandler(useCase)
	r := chi.NewRouter()
	r.Put("/v1/ms-resume/certification/{id}", h.Update)
	body := `{"name":"AWS Certified","nameEng":"AWS Certified","issuingOrganization":"Amazon","issuingOrganizationEng":"Amazon","issueDate":"2024-06-01"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/certification/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestCertificationHandler_Delete_Success(t *testing.T) {
	useCase := &mockCertificationUseCase{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	h := NewCertificationHandler(useCase)
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/certification/{id}", h.Delete)
	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/certification/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Language Handler Tests
// ---------------------------------------------------------------------------

func TestLanguageHandler_List_Success(t *testing.T) {
	expected := []response.LanguageResponse{
		{ID: 1, Language: "English", ReadingLevel: "C1", WritingLevel: "C1", SpeakingLevel: "C1"},
	}
	useCase := &mockLanguageUseCase{
		listFn: func(ctx context.Context) ([]response.LanguageResponse, error) { return expected, nil },
	}
	h := NewLanguageHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/language", h.List)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/language", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestLanguageHandler_Create_Success(t *testing.T) {
	useCase := &mockLanguageUseCase{
		createFn: func(ctx context.Context, req *request.LanguageRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewLanguageHandler(useCase)
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/language", h.Create)
	body := `{"language":"English","languageEng":"English","readingLevel":"C1","writingLevel":"C1","speakingLevel":"C1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/language", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestLanguageHandler_Update_Success(t *testing.T) {
	useCase := &mockLanguageUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.LanguageRequest) error { return nil },
	}
	h := NewLanguageHandler(useCase)
	r := chi.NewRouter()
	r.Put("/v1/ms-resume/language/{id}", h.Update)
	body := `{"language":"English","languageEng":"English","readingLevel":"C1","writingLevel":"C1","speakingLevel":"C1"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/language/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestLanguageHandler_Delete_Success(t *testing.T) {
	useCase := &mockLanguageUseCase{
		deleteFn: func(ctx context.Context, id int64) error { return nil },
	}
	h := NewLanguageHandler(useCase)
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/language/{id}", h.Delete)
	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/language/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Reference Handler Tests
// ---------------------------------------------------------------------------

func TestReferenceHandler_List_Success(t *testing.T) {
	expected := []response.ReferenceResponse{
		{ID: 1, FullName: "Jane Smith", Position: "Manager"},
	}
	useCase := &mockReferenceUseCase{
		listFn: func(ctx context.Context) ([]response.ReferenceResponse, error) { return expected, nil },
	}
	h := NewReferenceHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/reference", h.List)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/reference", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestReferenceHandler_Create_Success(t *testing.T) {
	useCase := &mockReferenceUseCase{
		createFn: func(ctx context.Context, req *request.ReferenceRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewReferenceHandler(useCase)
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/reference", h.Create)
	body := `{"fullName":"Jane Smith","position":"Manager"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/reference", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestReferenceHandler_Update_Success(t *testing.T) {
	useCase := &mockReferenceUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.ReferenceRequest) error { return nil },
	}
	h := NewReferenceHandler(useCase)
	r := chi.NewRouter()
	r.Put("/v1/ms-resume/reference/{id}", h.Update)
	body := `{"fullName":"Jane Smith","position":"Manager"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/reference/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestReferenceHandler_Delete_Success(t *testing.T) {
	useCase := &mockReferenceUseCase{
		deleteFn: func(ctx context.Context, id int64) error { return nil },
	}
	h := NewReferenceHandler(useCase)
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/reference/{id}", h.Delete)
	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/reference/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// CustomSection Handler Tests
// ---------------------------------------------------------------------------

func TestCustomSectionHandler_List_Success(t *testing.T) {
	expected := []response.CustomSectionResponse{
		{ID: 1, Title: "Projects", Visible: true},
	}
	useCase := &mockCustomSectionUseCase{
		listFn: func(ctx context.Context) ([]response.CustomSectionResponse, error) { return expected, nil },
	}
	h := NewCustomSectionHandler(useCase)
	r := chi.NewRouter()
	r.Get("/v1/ms-resume/custom-section", h.List)
	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/custom-section", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCustomSectionHandler_Create_Success(t *testing.T) {
	useCase := &mockCustomSectionUseCase{
		createFn: func(ctx context.Context, req *request.CustomSectionRequest) (*response.CreateResponse, error) {
			return &response.CreateResponse{ID: 1}, nil
		},
	}
	h := NewCustomSectionHandler(useCase)
	r := chi.NewRouter()
	r.Post("/v1/ms-resume/custom-section", h.Create)
	body := `{"title":"Projects","titleEng":"Projects","visible":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/ms-resume/custom-section", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestCustomSectionHandler_Update_Success(t *testing.T) {
	useCase := &mockCustomSectionUseCase{
		updateFn: func(ctx context.Context, id int64, req *request.CustomSectionRequest) error { return nil },
	}
	h := NewCustomSectionHandler(useCase)
	r := chi.NewRouter()
	r.Put("/v1/ms-resume/custom-section/{id}", h.Update)
	body := `{"title":"Projects","titleEng":"Projects","visible":true}`
	req := httptest.NewRequest(http.MethodPut, "/v1/ms-resume/custom-section/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestCustomSectionHandler_Delete_Success(t *testing.T) {
	useCase := &mockCustomSectionUseCase{
		deleteFn: func(ctx context.Context, id int64) error { return nil },
	}
	h := NewCustomSectionHandler(useCase)
	r := chi.NewRouter()
	r.Delete("/v1/ms-resume/custom-section/{id}", h.Delete)
	req := httptest.NewRequest(http.MethodDelete, "/v1/ms-resume/custom-section/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Error Handler Tests
// ---------------------------------------------------------------------------

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	WriteError(w, r, http.StatusBadRequest, "TEST_ERROR", "Test error message")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp ApiErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Code != "TEST_ERROR" {
		t.Errorf("expected code TEST_ERROR, got %s", errResp.Code)
	}
	if errResp.Message != "Test error message" {
		t.Errorf("expected message 'Test error message', got '%s'", errResp.Message)
	}
	if errResp.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, errResp.Status)
	}
	if errResp.Path != "/test" {
		t.Errorf("expected path '/test', got '%s'", errResp.Path)
	}
	if errResp.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
	// TraceID is empty by default when middleware is not configured
	// (middleware extracts it from context, which is not set in unit tests)
}

func TestWriteError_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/internal", nil)
	WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var errResp ApiErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Code != "INTERNAL_ERROR" {
		t.Errorf("expected code INTERNAL_ERROR, got %s", errResp.Code)
	}
}

// ---------------------------------------------------------------------------
// Mock PublicUseCase
// ---------------------------------------------------------------------------

type mockPublicUseCase struct {
	input.PublicUseCase
	getInfoPageFn func(ctx context.Context) (*response.InfoPageResponse, error)
}

func (m *mockPublicUseCase) GetInfoPage(ctx context.Context) (*response.InfoPageResponse, error) {
	return m.getInfoPageFn(ctx)
}

// ---------------------------------------------------------------------------
// Public Handler Integration Tests
// ---------------------------------------------------------------------------

func TestPublicHandler_GetInfoPage_Returns11Fields(t *testing.T) {
	// Build a fully populated InfoPageResponse with all 11 fields
	homeResp := &response.HomeResponse{
		ID: 1, Greeting: "Hello", Labels: []response.LabelResponse{{ID: 1, Name: "Backend", NameEng: "Backend"}},
	}
	basicDataResp := &response.BasicDataResponse{
		ID: 1, FirstName: "John", FirstSurname: "Doe", Email: "john@example.com",
	}
	skillsResp := []response.SkillResponse{
		{ID: 1, Name: "Go", NameEng: "Go", Sons: []response.SkillSonResponse{{ID: 1, Name: "Fiber", NameEng: "Fiber"}}},
	}
	expResp := []response.ExperienceResponse{
		{ID: 1, YearStart: "2020", Company: "Acme", SkillSons: []response.SkillSonResponse{}},
	}
	eduResp := []response.EducationResponse{
		{ID: 1, Institution: "MIT", Area: "CS", Degree: "BSc"},
	}
	altchaResp := &response.AltchaChallengeResponse{
		Algorithm: "SHA-256", Challenge: "abc123", Salt: "salty", Signature: "signed",
	}
	courseResp := []response.CourseResponse{
		{ID: 1, Name: "Go Basics", Institution: "Udemy", CompletionDate: "2024-01-01"},
	}
	certResp := []response.CertificationResponse{
		{ID: 1, Name: "AWS Certified", IssuingOrganization: "Amazon", IssueDate: "2024-06-01"},
	}
	langResp := []response.LanguageResponse{
		{ID: 1, Language: "English", ReadingLevel: "C1", WritingLevel: "C1", SpeakingLevel: "C1"},
	}
	refResp := []response.ReferenceResponse{
		{ID: 1, FullName: "Jane Smith", Position: "Manager"},
	}
	csResp := []response.CustomSectionResponse{
		{ID: 1, Title: "Projects", Visible: true},
	}

	useCase := &mockPublicUseCase{
		getInfoPageFn: func(ctx context.Context) (*response.InfoPageResponse, error) {
			return &response.InfoPageResponse{
				Home:            homeResp,
				BasicData:       basicDataResp,
				Skills:          skillsResp,
				Experiences:     expResp,
				Educations:      eduResp,
				AltchaChallenge: altchaResp,
				Courses:         courseResp,
				Certifications:  certResp,
				Languages:       langResp,
				References:      refResp,
				CustomSections:  csResp,
			}, nil
		},
	}

	// We need the other use cases for the handler constructor.
	// Use minimal mocks for PDF and Contact use cases.
	pdfUseCase := &mockPdfUseCase{}
	contactUseCase := &mockContactUseCase{}
	authHandler := NewAuthHandler(nil, nil, nil)

	pubHandler := NewPublicHandler(useCase, pdfUseCase, contactUseCase, authHandler)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/public/info-page", pubHandler.GetInfoPage)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/public/info-page", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify all 11 fields exist
	expectedFields := []string{
		"home", "basicData", "skills", "experiences", "educations", "altchaChallenge",
		"courses", "certifications", "languages", "references", "customSections",
	}
	for _, field := range expectedFields {
		if _, ok := resp[field]; !ok {
			t.Errorf("expected field '%s' to be present in response", field)
		}
	}

	// Verify the 5 new fields have the correct types (arrays)
	for _, field := range []string{"courses", "certifications", "languages", "references", "customSections"} {
		arr, ok := resp[field].([]interface{})
		if !ok {
			t.Errorf("expected field '%s' to be an array, got %T", field, resp[field])
			continue
		}
		if len(arr) == 0 {
			t.Errorf("expected field '%s' to have items", field)
		}
	}

	// Verify the 6 existing fields are present
	if _, ok := resp["home"]; !ok {
		t.Error("expected field 'home' to be present")
	}
	if _, ok := resp["basicData"]; !ok {
		t.Error("expected field 'basicData' to be present")
	}
	if _, ok := resp["skills"]; !ok {
		t.Error("expected field 'skills' to be present")
	}
}

func TestPublicHandler_GetInfoPage_EmptyArrays(t *testing.T) {
	useCase := &mockPublicUseCase{
		getInfoPageFn: func(ctx context.Context) (*response.InfoPageResponse, error) {
			return &response.InfoPageResponse{
				Home:            &response.HomeResponse{ID: 1, Greeting: "Hello"},
				BasicData:       &response.BasicDataResponse{ID: 1, FirstName: "John", Email: "j@example.com"},
				Skills:          []response.SkillResponse{},
				Experiences:     []response.ExperienceResponse{},
				Educations:      []response.EducationResponse{},
				AltchaChallenge: &response.AltchaChallengeResponse{Algorithm: "SHA-256", Challenge: "x", Salt: "s", Signature: "sig"},
				Courses:         []response.CourseResponse{},
				Certifications:  []response.CertificationResponse{},
				Languages:       []response.LanguageResponse{},
				References:      []response.ReferenceResponse{},
				CustomSections:  []response.CustomSectionResponse{},
			}, nil
		},
	}

	pdfUseCase := &mockPdfUseCase{}
	contactUseCase := &mockContactUseCase{}
	authHandler := NewAuthHandler(nil, nil, nil)

	pubHandler := NewPublicHandler(useCase, pdfUseCase, contactUseCase, authHandler)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/public/info-page", pubHandler.GetInfoPage)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/public/info-page", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// All slice fields should appear as [] when empty (omitempty removed from slices)
	for _, field := range []string{"courses", "certifications", "languages", "references", "customSections"} {
		arr, ok := resp[field].([]interface{})
		if !ok {
			t.Errorf("expected field '%s' to be an array [], got %T", field, resp[field])
			continue
		}
		if len(arr) != 0 {
			t.Errorf("expected field '%s' to be empty, got %d items", field, len(arr))
		}
	}

	// Same for existing slice fields
	for _, field := range []string{"skills", "experiences", "educations"} {
		arr, ok := resp[field].([]interface{})
		if !ok {
			t.Errorf("expected field '%s' to be an array [], got %T", field, resp[field])
			continue
		}
		if len(arr) != 0 {
			t.Errorf("expected field '%s' to be empty, got %d items", field, len(arr))
		}
	}

	// Verify all 11 fields are present
	if len(resp) != 11 {
		t.Errorf("expected exactly 11 fields, got %d: %v", len(resp), resp)
	}
}

func TestPublicHandler_GetInfoPage_ErrorReturns500(t *testing.T) {
	useCase := &mockPublicUseCase{
		getInfoPageFn: func(ctx context.Context) (*response.InfoPageResponse, error) {
			return nil, fmt.Errorf("internal error")
		},
	}

	pdfUseCase := &mockPdfUseCase{}
	contactUseCase := &mockContactUseCase{}
	authHandler := NewAuthHandler(nil, nil, nil)

	pubHandler := NewPublicHandler(useCase, pdfUseCase, contactUseCase, authHandler)

	r := chi.NewRouter()
	r.Get("/v1/ms-resume/public/info-page", pubHandler.GetInfoPage)

	req := httptest.NewRequest(http.MethodGet, "/v1/ms-resume/public/info-page", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var errResp ApiErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Code != "INTERNAL_ERROR" {
		t.Errorf("expected code INTERNAL_ERROR, got %s", errResp.Code)
	}
	if errResp.Status != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", errResp.Status)
	}
}

// ---------------------------------------------------------------------------
// Minimal mocks for PdfUseCase and ContactUseCase
// ---------------------------------------------------------------------------

type mockPdfUseCase struct {
	input.PdfUseCase
}

type mockContactUseCase struct {
	input.ContactUseCase
}
