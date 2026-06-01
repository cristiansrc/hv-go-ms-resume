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

// ImageUrlService implements the ImageUrlUseCase interface.
type ImageUrlService struct {
	repo output.ImageUrlRepository
}

// NewImageUrlService creates a new ImageUrlService.
func NewImageUrlService(repo output.ImageUrlRepository) input.ImageUrlUseCase {
	return &ImageUrlService{repo: repo}
}

// List returns all image URLs.
func (s *ImageUrlService) List(ctx context.Context) ([]response.ImageUrlResponse, error) {
	data, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.ImageUrlResponse, len(data))
	for i, item := range data {
		resp[i] = mapImageUrlToResponse(&item)
	}
	return resp, nil
}

// GetByID retrieves an ImageUrl by ID.
func (s *ImageUrlService) GetByID(ctx context.Context, id int64) (*response.ImageUrlResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapImageUrlToResponse(data)
	return &resp, nil
}

// Create creates a new ImageUrl.
func (s *ImageUrlService) Create(ctx context.Context, req *request.ImageURLRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.ImageUrl{
		Name:      req.Name,
		NameEng:   req.NameEng,
		URL:       req.File,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Delete removes an ImageUrl.
func (s *ImageUrlService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func mapImageUrlToResponse(e *entity.ImageUrl) response.ImageUrlResponse {
	return response.ImageUrlResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
		URL:     e.URL,
	}
}
