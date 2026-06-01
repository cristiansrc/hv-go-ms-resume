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

// ExperienceService implements the ExperienceUseCase interface.
type ExperienceService struct {
	repo output.ExperienceRepository
}

// NewExperienceService creates a new ExperienceService.
func NewExperienceService(repo output.ExperienceRepository) input.ExperienceUseCase {
	return &ExperienceService{repo: repo}
}

// List retrieves all Experiences.
func (s *ExperienceService) List(ctx context.Context) ([]response.ExperienceResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.ExperienceResponse, len(items))
	for i, item := range items {
		resp, err := mapExperienceToResponse(&item)
		if err != nil {
			return nil, err
		}
		result[i] = *resp
	}
	return result, nil
}

// GetByID retrieves an Experience by ID.
func (s *ExperienceService) GetByID(ctx context.Context, id int64) (*response.ExperienceResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapExperienceToResponse(item)
}

// Create creates a new Experience.
func (s *ExperienceService) Create(ctx context.Context, req *request.ExperienceRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity, err := experienceRequestToEntity(req, now, now)
	if err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Experience.
func (s *ExperienceService) Update(ctx context.Context, id int64, req *request.ExperienceRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity, err := experienceRequestToEntity(req, existing.CreatedAt, now)
	if err != nil {
		return err
	}
	entity.ID = id
	return s.repo.Update(ctx, entity)
}

// Delete removes an Experience by ID.
func (s *ExperienceService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// experienceRequestToEntity builds an Experience entity from a request.
func experienceRequestToEntity(req *request.ExperienceRequest, createdAt, updatedAt string) (*entity.Experience, error) {
	var descItemsPdf *string
	if req.DescriptionItemsPdf != nil {
		b, err := json.Marshal(req.DescriptionItemsPdf)
		if err != nil {
			return nil, err
		}
		s := string(b)
		descItemsPdf = &s
	}

	var descItemsPdfEng *string
	if req.DescriptionItemsPdfEng != nil {
		b, err := json.Marshal(req.DescriptionItemsPdfEng)
		if err != nil {
			return nil, err
		}
		s := string(b)
		descItemsPdfEng = &s
	}

	return &entity.Experience{
		YearStart:              req.YearStart,
		YearEnd:                req.YearEnd,
		Company:                req.Company,
		Location:               req.Location,
		LocationEng:            req.LocationEng,
		Position:               req.Position,
		PositionEng:            req.PositionEng,
		Summary:                req.Summary,
		SummaryEng:             req.SummaryEng,
		SummaryPdf:             req.SummaryPdf,
		SummaryPdfEng:          req.SummaryPdfEng,
		DescriptionItemsPdf:    descItemsPdf,
		DescriptionItemsPdfEng: descItemsPdfEng,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}, nil
}

// mapExperienceToResponse maps an Experience entity to an ExperienceResponse.
// SkillSons are omitted (empty slice) as the repository does not handle the N:N relation.
func mapExperienceToResponse(e *entity.Experience) (*response.ExperienceResponse, error) {
	var descItemsPdf []string
	if e.DescriptionItemsPdf != nil && *e.DescriptionItemsPdf != "" {
		if err := json.Unmarshal([]byte(*e.DescriptionItemsPdf), &descItemsPdf); err != nil {
			return nil, err
		}
	}

	var descItemsPdfEng []string
	if e.DescriptionItemsPdfEng != nil && *e.DescriptionItemsPdfEng != "" {
		if err := json.Unmarshal([]byte(*e.DescriptionItemsPdfEng), &descItemsPdfEng); err != nil {
			return nil, err
		}
	}

	return &response.ExperienceResponse{
		ID:                     e.ID,
		YearStart:              e.YearStart,
		YearEnd:                e.YearEnd,
		Company:                e.Company,
		Location:               e.Location,
		LocationEng:            e.LocationEng,
		Position:               e.Position,
		PositionEng:            e.PositionEng,
		Summary:                e.Summary,
		SummaryEng:             e.SummaryEng,
		SummaryPdf:             e.SummaryPdf,
		SummaryPdfEng:          e.SummaryPdfEng,
		DescriptionItemsPdf:    descItemsPdf,
		DescriptionItemsPdfEng: descItemsPdfEng,
		SkillSons:              []response.SkillSonResponse{},
	}, nil
}
