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

	// Courses
	courses, err := s.courseRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	courseItems := make([]output.CVSectionItem, len(courses))
	for i, c := range courses {
		courseItems[i] = output.CVSectionItem{
			Title:       s.getLocalizedString(&c.Name, &c.NameEng, language),
			Subtitle:    s.getLocalizedString(&c.Institution, &c.InstitutionEng, language),
			Date:        c.CompletionDate,
			Description: s.getLocalizedString(c.Description, c.DescriptionEng, language),
		}
	}
	cvData.Sections["courses"] = []output.CVSection{
		{Name: "Courses", Items: courseItems},
	}

	// Certifications
	certifications, err := s.certificationRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get certifications: %w", err)
	}
	certItems := make([]output.CVSectionItem, len(certifications))
	for i, c := range certifications {
		certItems[i] = output.CVSectionItem{
			Title:       s.getLocalizedString(&c.Name, &c.NameEng, language),
			Subtitle:    s.getLocalizedString(&c.IssuingOrganization, &c.IssuingOrganizationEng, language),
			Date:        c.IssueDate,
			Description: s.getLocalizedString(c.Description, c.DescriptionEng, language),
		}
	}
	cvData.Sections["certifications"] = []output.CVSection{
		{Name: "Certifications", Items: certItems},
	}

	// Languages
	languages, err := s.languageRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}
	langItems := make([]output.CVSectionItem, len(languages))
	for i, l := range languages {
		langItems[i] = output.CVSectionItem{
			Title:       s.getLocalizedString(&l.Language, &l.LanguageEng, language),
			Description: fmt.Sprintf("Reading: %s, Writing: %s, Speaking: %s", l.ReadingLevel, l.WritingLevel, l.SpeakingLevel),
		}
	}
	cvData.Sections["languages"] = []output.CVSection{
		{Name: "Languages", Items: langItems},
	}

	// References
	references, err := s.referenceRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}
	refItems := make([]output.CVSectionItem, len(references))
	for i, r := range references {
		refItems[i] = output.CVSectionItem{
			Title:       r.FullName,
			Subtitle:    r.Position,
			Description: s.getLocalizedString(r.Company, r.CompanyEng, language),
		}
	}
	cvData.Sections["references"] = []output.CVSection{
		{Name: "References", Items: refItems},
	}

	// CustomSections (only visible)
	allCustomSections, err := s.customSectionRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom sections: %w", err)
	}
	var visibleSections []entity.CustomSection
	for _, cs := range allCustomSections {
		if cs.Visible {
			visibleSections = append(visibleSections, cs)
		}
	}
	csItems := make([]output.CVSectionItem, len(visibleSections))
	for i, cs := range visibleSections {
		csItems[i] = output.CVSectionItem{
			Title:       s.getLocalizedString(&cs.Title, &cs.TitleEng, language),
			Description: s.getLocalizedString(cs.Content, cs.ContentEng, language),
		}
	}
	cvData.Sections["custom_sections"] = []output.CVSection{
		{Name: "Custom Sections", Items: csItems},
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
