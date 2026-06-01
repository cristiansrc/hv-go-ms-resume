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

// HomeService implements the HomeUseCase interface.
type HomeService struct {
	repo      output.HomeRepository
	labelRepo output.LabelRepository
	imageRepo output.ImageUrlRepository
}

// NewHomeService creates a new HomeService.
func NewHomeService(repo output.HomeRepository, labelRepo output.LabelRepository, imageRepo output.ImageUrlRepository) input.HomeUseCase {
	return &HomeService{repo: repo, labelRepo: labelRepo, imageRepo: imageRepo}
}

// GetByID retrieves Home by ID.
func (s *HomeService) GetByID(ctx context.Context, id int64) (*response.HomeResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapHomeToResponse(ctx, data)
}

// Update modifies Home fields.
func (s *HomeService) Update(ctx context.Context, id int64, req *request.HomeRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Home{
		ID:                    id,
		Greeting:              req.Greeting,
		GreetingEng:           req.GreetingEng,
		ImageURLID:            req.ImageURLID,
		ButtonWorkLabel:       req.ButtonWorkLabel,
		ButtonWorkLabelEng:    req.ButtonWorkLabelEng,
		ButtonContactLabel:    req.ButtonContactLabel,
		ButtonContactLabelEng: req.ButtonContactLabelEng,
		CreatedAt:             existing.CreatedAt,
		UpdatedAt:             now,
	}

	return s.repo.Update(ctx, entity)
}

func (s *HomeService) mapHomeToResponse(ctx context.Context, e *entity.Home) (*response.HomeResponse, error) {
	resp := &response.HomeResponse{
		ID:                    e.ID,
		Greeting:              e.Greeting,
		GreetingEng:           e.GreetingEng,
		ButtonWorkLabel:       e.ButtonWorkLabel,
		ButtonWorkLabelEng:    e.ButtonWorkLabelEng,
		ButtonContactLabel:    e.ButtonContactLabel,
		ButtonContactLabelEng: e.ButtonContactLabelEng,
	}

	// Map image if present (direct FK reference)
	if e.ImageURLID != nil {
		img, err := s.imageRepo.GetByID(ctx, *e.ImageURLID)
		if err == nil {
			mapped := mapImageUrlToResponse(img)
			resp.ImageURL = &mapped
		}
	}

	// Labels (N:N relation via home_label table) - currently repos do not support
	// filtering by home_id, so return empty slice.
	resp.Labels = []response.LabelResponse{}

	return resp, nil
}
