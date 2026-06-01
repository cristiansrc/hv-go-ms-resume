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

// SkillTypeService implements the SkillTypeUseCase interface.
type SkillTypeService struct {
	repo output.SkillTypeRepository
}

// NewSkillTypeService creates a new SkillTypeService.
func NewSkillTypeService(repo output.SkillTypeRepository) input.SkillTypeUseCase {
	return &SkillTypeService{repo: repo}
}

// List retrieves all SkillTypes.
func (s *SkillTypeService) List(ctx context.Context) ([]response.SkillTypeResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.SkillTypeResponse, len(items))
	for i, item := range items {
		result[i] = *mapSkillTypeToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a SkillType by ID.
func (s *SkillTypeService) GetByID(ctx context.Context, id int64) (*response.SkillTypeResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapSkillTypeToResponse(item), nil
}

// Create creates a new SkillType.
func (s *SkillTypeService) Create(ctx context.Context, req *request.SkillTypeRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.SkillType{
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

// Update modifies an existing SkillType.
func (s *SkillTypeService) Update(ctx context.Context, id int64, req *request.SkillTypeRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.SkillType{
		ID:        id,
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a SkillType by ID.
func (s *SkillTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapSkillTypeToResponse maps a SkillType entity to a SkillTypeResponse.
// Skills are omitted (empty slice) as the repository does not handle the N:N relation.
func mapSkillTypeToResponse(e *entity.SkillType) *response.SkillTypeResponse {
	return &response.SkillTypeResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
		Skills:  []response.SkillResponse{},
	}
}
