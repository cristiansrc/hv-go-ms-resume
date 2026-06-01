package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
	"github.com/cristiansrc/hv-go-ms-resume/pkg/hashutil"
)

// Sentinel errors for PDF service.
var (
	ErrInvalidLanguage      = errors.New("invalid language")
	ErrInvalidTemplate      = errors.New("invalid template")
	ErrRenderCVUnavailable  = errors.New("render service unavailable")
)

// PdfService implements the PdfUseCase interface with caching.
type PdfService struct {
	basicDataRepo      output.BasicDataRepository
	experienceRepo     output.ExperienceRepository
	educationRepo      output.EducationRepository
	courseRepo         output.CourseRepository
	certificationRepo  output.CertificationRepository
	languageRepo       output.LanguageRepository
	referenceRepo      output.ReferenceRepository
	customSectionRepo  output.CustomSectionRepository
	pdfCacheRepo       output.PdfCacheRepository
	renderCVPort       output.RenderCVPort
	s3Port             output.S3Port
	logger             *slog.Logger
}

// NewPdfService creates a new PdfService.
func NewPdfService(
	basicDataRepo output.BasicDataRepository,
	experienceRepo output.ExperienceRepository,
	educationRepo output.EducationRepository,
	courseRepo output.CourseRepository,
	certificationRepo output.CertificationRepository,
	languageRepo output.LanguageRepository,
	referenceRepo output.ReferenceRepository,
	customSectionRepo output.CustomSectionRepository,
	pdfCacheRepo output.PdfCacheRepository,
	renderCVPort output.RenderCVPort,
	s3Port output.S3Port,
	logger *slog.Logger,
) input.PdfUseCase {
	return &PdfService{
		basicDataRepo:     basicDataRepo,
		experienceRepo:    experienceRepo,
		educationRepo:     educationRepo,
		courseRepo:        courseRepo,
		certificationRepo: certificationRepo,
		languageRepo:      languageRepo,
		referenceRepo:     referenceRepo,
		customSectionRepo: customSectionRepo,
		pdfCacheRepo:      pdfCacheRepo,
		renderCVPort:      renderCVPort,
		s3Port:            s3Port,
		logger:            logger,
	}
}

// GetCurriculumPDF returns a PDF for the given language and template.
// It uses a hash-based cache to avoid regenerating the PDF when data hasn't changed.
func (s *PdfService) GetCurriculumPDF(ctx context.Context, language string, template string) (io.ReadCloser, string, error) {
	// Validate language and template
	validLanguages := map[string]bool{"english": true, "spanish": true}
	validTemplates := map[string]bool{
		"engineeringclassic": true, "engineeringresumes": true,
		"moderncv": true, "sb2nov": true,
	}

	if !validLanguages[language] {
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidLanguage, language)
	}
	if !validTemplates[template] {
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidTemplate, template)
	}

	// Collect all data needed for the PDF
	cvData, err := s.collectCvData(ctx, language)
	if err != nil {
		return nil, "", fmt.Errorf("failed to collect CV data: %w", err)
	}

	// Compute hash of the data for cache comparison
	dataHash, err := hashutil.ComputeHash(cvData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compute hash: %w", err)
	}

	// Check cache
	cached, err := s.pdfCacheRepo.GetByLanguageAndTemplate(ctx, language, template)
	if err != nil {
		s.logger.Error("failed to check pdf cache", "error", err)
	}
	if cached != nil && cached.DataHash == dataHash && s.s3Port != nil {
		// Cache hit - return from S3
		reader, err := s.s3Port.GetObject(ctx, cached.S3Key)
		if err != nil {
			s.logger.Error("failed to get cached PDF from S3, will regenerate", "error", err)
		} else {
			return reader, cached.DataHash, nil
		}
	}

	// Cache miss - generate via RenderCV
	pdfReader, err := s.renderCVPort.RenderPDF(ctx, cvData, language, template)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrRenderCVUnavailable, err)
	}
	defer pdfReader.Close()

	// Read all PDF bytes
	pdfBytes, err := io.ReadAll(pdfReader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read PDF: %w", err)
	}

	// Save to S3 (only if S3 is available)
	s3Key := fmt.Sprintf("pdfs/%s/%s/%s.pdf", language, template, dataHash)
	if s.s3Port != nil {
		if err := s.s3Port.PutObject(ctx, s3Key, bytes.NewReader(pdfBytes), "application/pdf"); err != nil {
			s.logger.Error("failed to save PDF to S3, returning without cache", "error", err)
		} else {
			// Delete old cache if hash changed
			if cached != nil && cached.DataHash != dataHash {
				if err := s.s3Port.DeleteObject(ctx, cached.S3Key); err != nil {
					s.logger.Error("failed to delete old cached PDF", "error", err)
				}
			}

			fileSize := int64(len(pdfBytes))
			// Update cache record
			cacheEntry := &entity.PdfCache{
				Language: language,
				Template: template,
				DataHash: dataHash,
				S3Key:    s3Key,
				FileSize: &fileSize,
			}
			if err := s.pdfCacheRepo.Upsert(ctx, cacheEntry); err != nil {
				s.logger.Error("failed to upsert pdf cache", "error", err)
			}
		}
	}

	return io.NopCloser(bytes.NewReader(pdfBytes)), dataHash, nil
}

func (s *PdfService) collectCvData(ctx context.Context, language string) (*output.CvData, error) {
	basicData, err := s.basicDataRepo.GetByID(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic data: %w", err)
	}

	cvData := &output.CvData{
		Name:     basicData.FirstName + " " + basicData.FirstSurname,
		Email:    basicData.Email,
		Location: s.getLocalizedString(basicData.Located, basicData.LocatedEng, language),
		Sections: make(map[string][]output.CVSection),
	}

	if basicData.Wrapper != nil {
		cvData.Headline = s.getLocalizedString(basicData.Wrapper, basicData.WrapperEng, language)
	}

	return cvData, nil
}

func (s *PdfService) getLocalizedString(es, en *string, language string) string {
	if language == "spanish" && es != nil {
		return *es
	}
	if en != nil {
		return *en
	}
	if es != nil {
		return *es
	}
	return ""
}
