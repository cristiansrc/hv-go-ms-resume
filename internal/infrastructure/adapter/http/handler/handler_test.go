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
			return nil, fmt.Errorf("basic_data not found")
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
