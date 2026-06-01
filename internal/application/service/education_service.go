package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// EducationService implements the EducationUseCase interface.
type EducationService struct {
	repo output.EducationRepository
}

// NewEducationService creates a new EducationService.
func NewEducationService(repo output.EducationRepository) input.EducationUseCase {
	return &EducationService{repo: repo}
}

// List retrieves all Education entries.
func (s *EducationService) List(ctx context.Context) ([]response.EducationResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.EducationResponse, len(items))
	for i, item := range items {
		resp, err := mapEducationToResponse(&item)
		if err != nil {
			return nil, err
		}
		result[i] = *resp
	}
	return result, nil
}

// GetByID retrieves an Education entry by ID.
func (s *EducationService) GetByID(ctx context.Context, id int64) (*response.EducationResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapEducationToResponse(item)
}

// Create creates a new Education entry.
func (s *EducationService) Create(ctx context.Context, req *request.EducationRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity, err := educationRequestToEntity(req, now, now)
	if err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Education entry.
func (s *EducationService) Update(ctx context.Context, id int64, req *request.EducationRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity, err := educationRequestToEntity(req, existing.CreatedAt, now)
	if err != nil {
		return err
	}
	entity.ID = id
	return s.repo.Update(ctx, entity)
}

// Delete removes an Education entry by ID.
func (s *EducationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// educationRequestToEntity builds an Education entity from a request.
func educationRequestToEntity(req *request.EducationRequest, createdAt, updatedAt string) (*entity.Education, error) {
	var highlights *string
	if req.Highlights != nil {
		b, err := json.Marshal(req.Highlights)
		if err != nil {
			return nil, err
		}
		s := string(b)
		highlights = &s
	}

	var highlightsEng *string
	if req.HighlightsEng != nil {
		b, err := json.Marshal(req.HighlightsEng)
		if err != nil {
			return nil, err
		}
		s := string(b)
		highlightsEng = &s
	}

	return &entity.Education{
		Institution:   req.Institution,
		Area:          req.Area,
		AreaEng:       req.AreaEng,
		Degree:        req.Degree,
		DegreeEng:     req.DegreeEng,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Location:      req.Location,
		LocationEng:   req.LocationEng,
		Highlights:    highlights,
		HighlightsEng: highlightsEng,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

// mapEducationToResponse maps an Education entity to an EducationResponse.
func mapEducationToResponse(e *entity.Education) (*response.EducationResponse, error) {
	var highlights []string
	if e.Highlights != nil && *e.Highlights != "" {
		if err := json.Unmarshal([]byte(*e.Highlights), &highlights); err != nil {
			return nil, err
		}
	}

	var highlightsEng []string
	if e.HighlightsEng != nil && *e.HighlightsEng != "" {
		if err := json.Unmarshal([]byte(*e.HighlightsEng), &highlightsEng); err != nil {
			return nil, err
		}
	}

	return &response.EducationResponse{
		ID:            e.ID,
		Institution:   e.Institution,
		Area:          e.Area,
		AreaEng:       e.AreaEng,
		Degree:        e.Degree,
		DegreeEng:     e.DegreeEng,
		StartDate:     e.StartDate,
		EndDate:       e.EndDate,
		Location:      e.Location,
		LocationEng:   e.LocationEng,
		Highlights:    highlights,
		HighlightsEng: highlightsEng,
	}, nil
}
