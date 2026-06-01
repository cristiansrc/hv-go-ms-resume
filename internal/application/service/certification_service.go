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

// CertificationService implements the CertificationUseCase interface.
type CertificationService struct {
	repo output.CertificationRepository
}

// NewCertificationService creates a new CertificationService.
func NewCertificationService(repo output.CertificationRepository) input.CertificationUseCase {
	return &CertificationService{repo: repo}
}

// List retrieves all Certifications.
func (s *CertificationService) List(ctx context.Context) ([]response.CertificationResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CertificationResponse, len(items))
	for i, item := range items {
		result[i] = *mapCertificationToResponse(&item)
	}
	return result, nil
}

// GetByID retrieves a Certification by ID.
func (s *CertificationService) GetByID(ctx context.Context, id int64) (*response.CertificationResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapCertificationToResponse(item), nil
}

// Create creates a new Certification.
func (s *CertificationService) Create(ctx context.Context, req *request.CertificationRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Certification{
		Name:                   req.Name,
		NameEng:                req.NameEng,
		IssuingOrganization:    req.IssuingOrganization,
		IssuingOrganizationEng: req.IssuingOrganizationEng,
		IssueDate:              req.IssueDate,
		ExpirationDate:         req.ExpirationDate,
		VerificationURL:        req.VerificationURL,
		CredentialID:           req.CredentialID,
		Description:            req.Description,
		DescriptionEng:         req.DescriptionEng,
		SummaryPdf:             req.SummaryPdf,
		SummaryPdfEng:          req.SummaryPdfEng,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies an existing Certification.
func (s *CertificationService) Update(ctx context.Context, id int64, req *request.CertificationRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Certification{
		ID:                     id,
		Name:                   req.Name,
		NameEng:                req.NameEng,
		IssuingOrganization:    req.IssuingOrganization,
		IssuingOrganizationEng: req.IssuingOrganizationEng,
		IssueDate:              req.IssueDate,
		ExpirationDate:         req.ExpirationDate,
		VerificationURL:        req.VerificationURL,
		CredentialID:           req.CredentialID,
		Description:            req.Description,
		DescriptionEng:         req.DescriptionEng,
		SummaryPdf:             req.SummaryPdf,
		SummaryPdfEng:          req.SummaryPdfEng,
		CreatedAt:              existing.CreatedAt,
		UpdatedAt:              now,
	}
	return s.repo.Update(ctx, entity)
}

// Delete removes a Certification by ID.
func (s *CertificationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// mapCertificationToResponse maps a Certification entity to a CertificationResponse.
func mapCertificationToResponse(e *entity.Certification) *response.CertificationResponse {
	return &response.CertificationResponse{
		ID:                     e.ID,
		Name:                   e.Name,
		NameEng:                e.NameEng,
		IssuingOrganization:    e.IssuingOrganization,
		IssuingOrganizationEng: e.IssuingOrganizationEng,
		IssueDate:              e.IssueDate,
		ExpirationDate:         e.ExpirationDate,
		VerificationURL:        e.VerificationURL,
		CredentialID:           e.CredentialID,
		Description:            e.Description,
		DescriptionEng:         e.DescriptionEng,
		SummaryPdf:             e.SummaryPdf,
		SummaryPdfEng:          e.SummaryPdfEng,
	}
}
