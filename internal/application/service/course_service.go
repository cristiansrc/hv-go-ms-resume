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

// CourseService implements the CourseUseCase interface.
type CourseService struct {
	repo output.CourseRepository
}

// NewCourseService creates a new CourseService.
func NewCourseService(repo output.CourseRepository) input.CourseUseCase {
	return &CourseService{repo: repo}
}

// List retrieves all Courses.
func (s *CourseService) List(ctx context.Context) ([]response.CourseResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CourseResponse, len(items))
	for i, item := range items {
		result[i] = *mapCourseToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a Course by ID.
func (s *CourseService) GetByID(ctx context.Context, id int64) (*response.CourseResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapCourseToResponse(item), nil
}

// Create creates a new Course.
func (s *CourseService) Create(ctx context.Context, req *request.CourseRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Course{
		Name:            req.Name,
		NameEng:         req.NameEng,
		Institution:     req.Institution,
		InstitutionEng:  req.InstitutionEng,
		CompletionDate:  req.CompletionDate,
		Description:     req.Description,
		DescriptionEng:  req.DescriptionEng,
		SummaryPdf:      req.SummaryPdf,
		SummaryPdfEng:   req.SummaryPdfEng,
		CertificateURL:  req.CertificateURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Course.
func (s *CourseService) Update(ctx context.Context, id int64, req *request.CourseRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Course{
		ID:              id,
		Name:            req.Name,
		NameEng:         req.NameEng,
		Institution:     req.Institution,
		InstitutionEng:  req.InstitutionEng,
		CompletionDate:  req.CompletionDate,
		Description:     req.Description,
		DescriptionEng:  req.DescriptionEng,
		SummaryPdf:      req.SummaryPdf,
		SummaryPdfEng:   req.SummaryPdfEng,
		CertificateURL:  req.CertificateURL,
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a Course by ID.
func (s *CourseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapCourseToResponse maps a Course entity to a CourseResponse.
func mapCourseToResponse(e *entity.Course) *response.CourseResponse {
	return &response.CourseResponse{
		ID:              e.ID,
		Name:            e.Name,
		NameEng:         e.NameEng,
		Institution:     e.Institution,
		InstitutionEng:  e.InstitutionEng,
		CompletionDate:  e.CompletionDate,
		Description:     e.Description,
		DescriptionEng:  e.DescriptionEng,
		SummaryPdf:      e.SummaryPdf,
		SummaryPdfEng:   e.SummaryPdfEng,
		CertificateURL:  e.CertificateURL,
	}
}
