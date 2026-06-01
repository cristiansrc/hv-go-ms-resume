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

// SkillSonService implements the SkillSonUseCase interface.
type SkillSonService struct {
	repo output.SkillSonRepository
}

// NewSkillSonService creates a new SkillSonService.
func NewSkillSonService(repo output.SkillSonRepository) input.SkillSonUseCase {
	return &SkillSonService{repo: repo}
}

// List retrieves all SkillSons.
func (s *SkillSonService) List(ctx context.Context) ([]response.SkillSonResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.SkillSonResponse, len(items))
	for i, item := range items {
		result[i] = *mapSkillSonToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a SkillSon by ID.
func (s *SkillSonService) GetByID(ctx context.Context, id int64) (*response.SkillSonResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapSkillSonToResponse(item), nil
}

// Create creates a new SkillSon.
func (s *SkillSonService) Create(ctx context.Context, req *request.SkillSonRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.SkillSon{
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: now,
		UpdatedAt: now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing SkillSon.
func (s *SkillSonService) Update(ctx context.Context, id int64, req *request.SkillSonRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.SkillSon{
		ID:        id,
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a SkillSon by ID.
func (s *SkillSonService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapSkillSonToResponse maps a SkillSon entity to a SkillSonResponse.
func mapSkillSonToResponse(e *entity.SkillSon) *response.SkillSonResponse {
	return &response.SkillSonResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
	}
}
