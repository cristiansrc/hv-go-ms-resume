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

// LabelService implements the LabelUseCase interface.
type LabelService struct {
	repo output.LabelRepository
}

// NewLabelService creates a new LabelService.
func NewLabelService(repo output.LabelRepository) input.LabelUseCase {
	return &LabelService{repo: repo}
}

// List returns all labels.
func (s *LabelService) List(ctx context.Context) ([]response.LabelResponse, error) {
	data, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.LabelResponse, len(data))
	for i, item := range data {
		resp[i] = mapLabelToResponse(&item)
	}
	return resp, nil
}

// GetByID retrieves a Label by ID.
func (s *LabelService) GetByID(ctx context.Context, id int64) (*response.LabelResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapLabelToResponse(data)
	return &resp, nil
}

// Create creates a new Label.
func (s *LabelService) Create(ctx context.Context, req *request.LabelRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Label{
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

// Update modifies a Label.
func (s *LabelService) Update(ctx context.Context, id int64, req *request.LabelRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Label{
		ID:        id,
		Name:      req.Name,
		NameEng:   req.NameEng,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}

	return s.repo.Update(ctx, entity)
}

// Delete removes a Label.
func (s *LabelService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func mapLabelToResponse(e *entity.Label) response.LabelResponse {
	return response.LabelResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
	}
}
