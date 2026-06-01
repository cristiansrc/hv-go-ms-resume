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

// FuturedProjectService implements the FuturedProjectUseCase interface.
type FuturedProjectService struct {
	repo output.FuturedProjectRepository
}

// NewFuturedProjectService creates a new FuturedProjectService.
func NewFuturedProjectService(repo output.FuturedProjectRepository) input.FuturedProjectUseCase {
	return &FuturedProjectService{repo: repo}
}

// List retrieves all FuturedProjects.
func (s *FuturedProjectService) List(ctx context.Context) ([]response.FuturedProjectResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.FuturedProjectResponse, len(items))
	for i, item := range items {
		result[i] = *mapFuturedProjectToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a FuturedProject by ID.
func (s *FuturedProjectService) GetByID(ctx context.Context, id int64) (*response.FuturedProjectResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapFuturedProjectToResponse(item), nil
}

// Create creates a new FuturedProject.
func (s *FuturedProjectService) Create(ctx context.Context, req *request.FuturedProjectRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.FuturedProject{
		Name:               req.Name,
		NameEng:            req.NameEng,
		ExperienceID:       req.ExperienceID,
		DescriptionShort:   req.DescriptionShort,
		Description:        req.Description,
		DescriptionShortEng: req.DescriptionShortEng,
		DescriptionEng:     req.DescriptionEng,
		ImageListURLID:     req.ImageListURLID,
		ImageURLID:         req.ImageURLID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing FuturedProject.
func (s *FuturedProjectService) Update(ctx context.Context, id int64, req *request.FuturedProjectRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.FuturedProject{
		ID:                 id,
		Name:               req.Name,
		NameEng:            req.NameEng,
		ExperienceID:       req.ExperienceID,
		DescriptionShort:   req.DescriptionShort,
		Description:        req.Description,
		DescriptionShortEng: req.DescriptionShortEng,
		DescriptionEng:     req.DescriptionEng,
		ImageListURLID:     req.ImageListURLID,
		ImageURLID:         req.ImageURLID,
		CreatedAt:          existing.CreatedAt,
		UpdatedAt:          now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a FuturedProject by ID.
func (s *FuturedProjectService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapFuturedProjectToResponse maps a FuturedProject entity to a FuturedProjectResponse.
// Nested Experience and ImageUrl references are omitted as they require additional repository lookups.
func mapFuturedProjectToResponse(e *entity.FuturedProject) *response.FuturedProjectResponse {
	return &response.FuturedProjectResponse{
		ID:                 e.ID,
		Name:               e.Name,
		NameEng:            e.NameEng,
		DescriptionShort:   e.DescriptionShort,
		Description:        e.Description,
		DescriptionShortEng: e.DescriptionShortEng,
		DescriptionEng:     e.DescriptionEng,
		Experience:         nil,
		ImageListURL:       nil,
		ImageURL:           nil,
	}
}
