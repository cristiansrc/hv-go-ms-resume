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

// LanguageService implements the LanguageUseCase interface.
type LanguageService struct {
	repo output.LanguageRepository
}

// NewLanguageService creates a new LanguageService.
func NewLanguageService(repo output.LanguageRepository) input.LanguageUseCase {
	return &LanguageService{repo: repo}
}

// List retrieves all Languages.
func (s *LanguageService) List(ctx context.Context) ([]response.LanguageResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.LanguageResponse, len(items))
	for i, item := range items {
		result[i] = *mapLanguageToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a Language by ID.
func (s *LanguageService) GetByID(ctx context.Context, id int64) (*response.LanguageResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapLanguageToResponse(item), nil
}

// Create creates a new Language.
func (s *LanguageService) Create(ctx context.Context, req *request.LanguageRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Language{
		Language:      req.Language,
		LanguageEng:   req.LanguageEng,
		ReadingLevel:  req.ReadingLevel,
		WritingLevel:  req.WritingLevel,
		SpeakingLevel: req.SpeakingLevel,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Language.
func (s *LanguageService) Update(ctx context.Context, id int64, req *request.LanguageRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Language{
		ID:            id,
		Language:      req.Language,
		LanguageEng:   req.LanguageEng,
		ReadingLevel:  req.ReadingLevel,
		WritingLevel:  req.WritingLevel,
		SpeakingLevel: req.SpeakingLevel,
		CreatedAt:     existing.CreatedAt,
		UpdatedAt:     now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a Language by ID.
func (s *LanguageService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapLanguageToResponse maps a Language entity to a LanguageResponse.
func mapLanguageToResponse(e *entity.Language) *response.LanguageResponse {
	return &response.LanguageResponse{
		ID:            e.ID,
		Language:      e.Language,
		LanguageEng:   e.LanguageEng,
		ReadingLevel:  e.ReadingLevel,
		WritingLevel:  e.WritingLevel,
		SpeakingLevel: e.SpeakingLevel,
	}
}
