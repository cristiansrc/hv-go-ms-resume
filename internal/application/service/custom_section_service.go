package service

import (
	"context"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// CustomSectionService implements the CustomSectionUseCase interface.
type CustomSectionService struct {
	repo output.CustomSectionRepository
}

// NewCustomSectionService creates a new CustomSectionService.
func NewCustomSectionService(repo output.CustomSectionRepository) input.CustomSectionUseCase {
	return &CustomSectionService{repo: repo}
}

// List retrieves all CustomSections.
func (s *CustomSectionService) List(ctx context.Context) ([]response.CustomSectionResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CustomSectionResponse, len(items))
	for i, item := range items {
		result[i] = *mapCustomSectionToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a CustomSection by ID.
func (s *CustomSectionService) GetByID(ctx context.Context, id int64) (*response.CustomSectionResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapCustomSectionToResponse(item), nil
}

// Create creates a new CustomSection.
func (s *CustomSectionService) Create(ctx context.Context, req *request.CustomSectionRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.CustomSection{
		Title:         req.Title,
		TitleEng:      req.TitleEng,
		Content:       req.Content,
		ContentEng:    req.ContentEng,
		SummaryPdf:    req.SummaryPdf,
		SummaryPdfEng: req.SummaryPdfEng,
		Visible:       req.Visible,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing CustomSection.
func (s *CustomSectionService) Update(ctx context.Context, id int64, req *request.CustomSectionRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.CustomSection{
		ID:            id,
		Title:         req.Title,
		TitleEng:      req.TitleEng,
		Content:       req.Content,
		ContentEng:    req.ContentEng,
		SummaryPdf:    req.SummaryPdf,
		SummaryPdfEng: req.SummaryPdfEng,
		Visible:       req.Visible,
		CreatedAt:     existing.CreatedAt,
		UpdatedAt:     now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a CustomSection by ID.
func (s *CustomSectionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapCustomSectionToResponse maps a CustomSection entity to a CustomSectionResponse.
func mapCustomSectionToResponse(e *entity.CustomSection) *response.CustomSectionResponse {
	return &response.CustomSectionResponse{
		ID:            e.ID,
		Title:         e.Title,
		TitleEng:      e.TitleEng,
		Content:       e.Content,
		ContentEng:    e.ContentEng,
		SummaryPdf:    e.SummaryPdf,
		SummaryPdfEng: e.SummaryPdfEng,
		Visible:       e.Visible,
	}
}
