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

// BlogTypeService implements the BlogTypeUseCase interface.
type BlogTypeService struct {
	repo output.BlogTypeRepository
}

// NewBlogTypeService creates a new BlogTypeService.
func NewBlogTypeService(repo output.BlogTypeRepository) input.BlogTypeUseCase {
	return &BlogTypeService{repo: repo}
}

// List returns all blog types.
func (s *BlogTypeService) List(ctx context.Context) ([]response.BlogTypeResponse, error) {
	data, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.BlogTypeResponse, len(data))
	for i, item := range data {
		resp[i] = mapBlogTypeToResponse(&item)
	}
	return resp, nil
}

// GetByID retrieves a BlogType by ID.
func (s *BlogTypeService) GetByID(ctx context.Context, id int64) (*response.BlogTypeResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapBlogTypeToResponse(data)
	return &resp, nil
}

// Create creates a new BlogType.
func (s *BlogTypeService) Create(ctx context.Context, req *request.BlogTypeRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.BlogType{
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

// Update modifies a BlogType.
func (s *BlogTypeService) Update(ctx context.Context, id int64, req *request.BlogTypeRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.BlogType{
		ID:        id,
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}

	return s.repo.Update(ctx, entity)
}

// Delete removes a BlogType.
func (s *BlogTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func mapBlogTypeToResponse(e *entity.BlogType) response.BlogTypeResponse {
	return response.BlogTypeResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
	}
}
