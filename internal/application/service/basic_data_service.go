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

// BasicDataService implements the BasicDataUseCase interface.
type BasicDataService struct {
	repo output.BasicDataRepository
}

// NewBasicDataService creates a new BasicDataService.
func NewBasicDataService(repo output.BasicDataRepository) input.BasicDataUseCase {
	return &BasicDataService{repo: repo}
}

// GetByID retrieves BasicData by ID.
func (s *BasicDataService) GetByID(ctx context.Context, id int64) (*response.BasicDataResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapBasicDataToResponse(data), nil
}

// Update modifies BasicData fields.
func (s *BasicDataService) Update(ctx context.Context, id int64, req *request.BasicDataRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	bd := &entity.BasicData{
		ID:                id,
		FirstName:         req.FirstName,
		OthersName:        req.OthersName,
		FirstSurname:      req.FirstSurname,
		OthersSurname:     req.OthersSurname,
		DateBirth:         req.DateBirth,
		Located:           req.Located,
		LocatedEng:        req.LocatedEng,
		StartWorkingDate:  req.StartWorkingDate,
		Greeting:          req.Greeting,
		GreetingEng:       req.GreetingEng,
		Email:             req.Email,
		Instagram:         req.Instagram,
		Linkedin:          req.Linkedin,
		X:                 req.X,
		Github:            req.Github,
		Description:       req.Description,
		DescriptionEng:    req.DescriptionEng,
		DescriptionPdf:    req.DescriptionPdf,
		DescriptionPdfEng: req.DescriptionPdfEng,
		Wrapper:           req.Wrapper,
		WrapperEng:        req.WrapperEng,
		CreatedAt:         existing.CreatedAt,
		UpdatedAt:         now,
	}

	return s.repo.Update(ctx, bd)
}

func mapBasicDataToResponse(e *entity.BasicData) *response.BasicDataResponse {
	return &response.BasicDataResponse{
		ID:                e.ID,
		FirstName:         e.FirstName,
		OthersName:        e.OthersName,
		FirstSurname:      e.FirstSurname,
		OthersSurname:     e.OthersSurname,
		DateBirth:         e.DateBirth,
		Located:           e.Located,
		LocatedEng:        e.LocatedEng,
		StartWorkingDate:  e.StartWorkingDate,
		Greeting:          e.Greeting,
		GreetingEng:       e.GreetingEng,
		Email:             e.Email,
		Instagram:         e.Instagram,
		Linkedin:          e.Linkedin,
		X:                 e.X,
		Github:            e.Github,
		Description:       e.Description,
		DescriptionEng:    e.DescriptionEng,
		DescriptionPdf:    e.DescriptionPdf,
		DescriptionPdfEng: e.DescriptionPdfEng,
		Wrapper:           e.Wrapper,
		WrapperEng:        e.WrapperEng,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}
