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

// VideoUrlService implements the VideoUrlUseCase interface.
type VideoUrlService struct {
	repo output.VideoUrlRepository
}

// NewVideoUrlService creates a new VideoUrlService.
func NewVideoUrlService(repo output.VideoUrlRepository) input.VideoUrlUseCase {
	return &VideoUrlService{repo: repo}
}

// List returns all video URLs.
func (s *VideoUrlService) List(ctx context.Context) ([]response.VideoUrlResponse, error) {
	data, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.VideoUrlResponse, len(data))
	for i, item := range data {
		resp[i] = mapVideoUrlToResponse(&item)
	}
	return resp, nil
}

// GetByID retrieves a VideoUrl by ID.
func (s *VideoUrlService) GetByID(ctx context.Context, id int64) (*response.VideoUrlResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapVideoUrlToResponse(data)
	return &resp, nil
}

// Create creates a new VideoUrl.
func (s *VideoUrlService) Create(ctx context.Context, req *request.VideoURLRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.VideoUrl{
		Name:      req.Name,
		NameEng:   req.NameEng,
		URL:       req.URL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Delete removes a VideoUrl.
func (s *VideoUrlService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func mapVideoUrlToResponse(e *entity.VideoUrl) response.VideoUrlResponse {
	return response.VideoUrlResponse{
		ID:      e.ID,
		Name:    e.Name,
		NameEng: e.NameEng,
		URL:     e.URL,
	}
}
