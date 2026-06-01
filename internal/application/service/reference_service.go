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

// ReferenceService implements the ReferenceUseCase interface.
type ReferenceService struct {
	repo output.ReferenceRepository
}

// NewReferenceService creates a new ReferenceService.
func NewReferenceService(repo output.ReferenceRepository) input.ReferenceUseCase {
	return &ReferenceService{repo: repo}
}

// List retrieves all References.
func (s *ReferenceService) List(ctx context.Context) ([]response.ReferenceResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.ReferenceResponse, len(items))
	for i, item := range items {
		result[i] = *mapReferenceToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a Reference by ID.
func (s *ReferenceService) GetByID(ctx context.Context, id int64) (*response.ReferenceResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapReferenceToResponse(item), nil
}

// Create creates a new Reference.
func (s *ReferenceService) Create(ctx context.Context, req *request.ReferenceRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Reference{
		FullName:        req.FullName,
		Position:        req.Position,
		Company:         req.Company,
		CompanyEng:      req.CompanyEng,
		Email:           req.Email,
		Phone:           req.Phone,
		Relationship:    req.Relationship,
		RelationshipEng: req.RelationshipEng,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Reference.
func (s *ReferenceService) Update(ctx context.Context, id int64, req *request.ReferenceRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Reference{
		ID:              id,
		FullName:        req.FullName,
		Position:        req.Position,
		Company:         req.Company,
		CompanyEng:      req.CompanyEng,
		Email:           req.Email,
		Phone:           req.Phone,
		Relationship:    req.Relationship,
		RelationshipEng: req.RelationshipEng,
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a Reference by ID.
func (s *ReferenceService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapReferenceToResponse maps a Reference entity to a ReferenceResponse.
func mapReferenceToResponse(e *entity.Reference) *response.ReferenceResponse {
	return &response.ReferenceResponse{
		ID:              e.ID,
		FullName:        e.FullName,
		Position:        e.Position,
		Company:         e.Company,
		CompanyEng:      e.CompanyEng,
		Email:           e.Email,
		Phone:           e.Phone,
		Relationship:    e.Relationship,
		RelationshipEng: e.RelationshipEng,
	}
}
