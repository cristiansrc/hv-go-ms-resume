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

// SkillService implements the SkillUseCase interface.
type SkillService struct {
	repo output.SkillRepository
}

// NewSkillService creates a new SkillService.
func NewSkillService(repo output.SkillRepository) input.SkillUseCase {
	return &SkillService{repo: repo}
}

// List retrieves all Skills.
func (s *SkillService) List(ctx context.Context) ([]response.SkillResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.SkillResponse, len(items))
	for i, item := range items {
		result[i] = *mapSkillToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a Skill by ID.
func (s *SkillService) GetByID(ctx context.Context, id int64) (*response.SkillResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapSkillToResponse(item), nil
}

// Create creates a new Skill.
func (s *SkillService) Create(ctx context.Context, req *request.SkillRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Skill{
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

// Update modifies an existing Skill.
func (s *SkillService) Update(ctx context.Context, id int64, req *request.SkillRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Skill{
		ID:        id,
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a Skill by ID.
func (s *SkillService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapSkillToResponse maps a Skill entity to a SkillResponse.
// Sons are omitted (empty slice) as the repository does not handle the N:N relation.
func mapSkillToResponse(e *entity.Skill) *response.SkillResponse {
	return &response.SkillResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
		Sons:    []response.SkillSonResponse{},
	}
}
