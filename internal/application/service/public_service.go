package service

import (
	"context"
	"sort"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// PublicService implements the PublicUseCase interface for public read-only operations.
type PublicService struct {
	basicDataRepo      output.BasicDataRepository
	homeRepo           output.HomeRepository
	imageRepo          output.ImageUrlRepository
	labelRepo          output.LabelRepository
	blogRepo           output.BlogRepository
	blogTypeRepo       output.BlogTypeRepository
	experienceRepo     output.ExperienceRepository
	educationRepo      output.EducationRepository
	courseRepo         output.CourseRepository
	certificationRepo  output.CertificationRepository
	languageRepo       output.LanguageRepository
	referenceRepo      output.ReferenceRepository
	customSectionRepo  output.CustomSectionRepository
	skillRepo          output.SkillRepository
	skillSonRepo       output.SkillSonRepository
	altchaPort         output.AltchaPort
}

// NewPublicService creates a new PublicService.
func NewPublicService(
	basicDataRepo output.BasicDataRepository,
	homeRepo output.HomeRepository,
	imageRepo output.ImageUrlRepository,
	labelRepo output.LabelRepository,
	blogRepo output.BlogRepository,
	blogTypeRepo output.BlogTypeRepository,
	experienceRepo output.ExperienceRepository,
	educationRepo output.EducationRepository,
	courseRepo output.CourseRepository,
	certificationRepo output.CertificationRepository,
	languageRepo output.LanguageRepository,
	referenceRepo output.ReferenceRepository,
	customSectionRepo output.CustomSectionRepository,
	skillRepo output.SkillRepository,
	skillSonRepo output.SkillSonRepository,
	altchaPort output.AltchaPort,
) input.PublicUseCase {
	return &PublicService{
		basicDataRepo:     basicDataRepo,
		homeRepo:          homeRepo,
		imageRepo:         imageRepo,
		labelRepo:         labelRepo,
		blogRepo:          blogRepo,
		blogTypeRepo:      blogTypeRepo,
		experienceRepo:    experienceRepo,
		educationRepo:     educationRepo,
		courseRepo:        courseRepo,
		certificationRepo: certificationRepo,
		languageRepo:      languageRepo,
		referenceRepo:     referenceRepo,
		customSectionRepo: customSectionRepo,
		skillRepo:         skillRepo,
		skillSonRepo:      skillSonRepo,
		altchaPort:        altchaPort,
	}
}

// GetInfoPage retrieves the aggregated public resume information.
func (s *PublicService) GetInfoPage(ctx context.Context) (*response.InfoPageResponse, error) {
	// Basic data (single row, id=1)
	basicData, err := s.basicDataRepo.GetByID(ctx, 1)
	if err != nil {
		return nil, err
	}

	// Home (single row, id=1)
	home, err := s.homeRepo.GetByID(ctx, 1)
	if err != nil {
		return nil, err
	}

	// Experiences
	experiences, err := s.experienceRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Educations
	educations, err := s.educationRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Labels (flat slice)
	labels, err := s.labelRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Altcha challenge
	altchaChallenge, err := s.altchaPort.GenerateChallenge()
	if err != nil {
		return nil, err
	}

	// Map home with nested labels
	homeResp := mapHomeToResponse(home, &labels, s.imageRepo, ctx)

	// Skills with nested SkillSons
	skills, err := s.skillRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	skillResponses := make([]response.SkillResponse, 0, len(skills))
	for _, sk := range skills {
		sons, err := s.skillSonRepo.ListBySkillID(ctx, sk.ID)
		if err != nil {
			return nil, err
		}
		sonResponses := make([]response.SkillSonResponse, len(sons))
		for j, son := range sons {
			sonResponses[j] = response.SkillSonResponse{
				ID:      son.ID,
				Name:    son.Name,
				NameEng: son.NameEng,
			}
		}
		skillResponses = append(skillResponses, response.SkillResponse{
			ID:      sk.ID,
			Name:    sk.Name,
			NameEng: sk.NameEng,
			Sons:    sonResponses,
		})
	}

	// Map experiences (SkillSons omitted)
	expResponses := make([]response.ExperienceResponse, len(experiences))
	for i, exp := range experiences {
		mapped, err := mapExperienceToResponse(&exp)
		if err != nil {
			return nil, err
		}
		expResponses[i] = *mapped
	}

	// Map educations
	eduResponses := make([]response.EducationResponse, len(educations))
	for i, edu := range educations {
		mapped, err := mapEducationToResponse(&edu)
		if err != nil {
			return nil, err
		}
		eduResponses[i] = *mapped
	}

	return &response.InfoPageResponse{
		Home:      homeResp,
		BasicData: mapBasicDataToResponse(basicData),
		Skills:    skillResponses,
		Experiences: expResponses,
		Educations:  eduResponses,
		AltchaChallenge: &response.AltchaChallengeResponse{
			Algorithm: altchaChallenge.Algorithm,
			Challenge: altchaChallenge.Challenge,
			Salt:      altchaChallenge.Salt,
			Signature: altchaChallenge.Signature,
		},
	}, nil
}

// GetPublicBlogPage retrieves a paginated list of public blog entries.
func (s *PublicService) GetPublicBlogPage(ctx context.Context, page, size int, sortBy string) (*response.BlogPageResponse, error) {
	allBlogs, err := s.blogRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Sort by ID descending
	sort.Slice(allBlogs, func(i, j int) bool {
		return allBlogs[i].ID > allBlogs[j].ID
	})

	// Paginate in memory
	paginated, totalPages, totalElements := paginate(allBlogs, page, size)

	// Map to responses
	content := make([]response.BlogResponse, len(paginated))
	for i, blog := range paginated {
		mapped, err := s.mapBlogToResponse(ctx, &blog)
		if err != nil {
			return nil, err
		}
		content[i] = *mapped
	}

	totalPagesInt := totalPages
	numberOfElements := len(content)
	last := page >= totalPagesInt-1
	first := page == 0
	empty := len(content) == 0

	// Determine sort info based on sortBy parameter
	sorted := false
	if sortBy != "" {
		sorted = true
	}

	return &response.BlogPageResponse{
		Content: content,
		Pageable: response.PageableInfo{
			Sort: response.SortInfo{
				Empty:    !sorted,
				Sorted:   sorted,
				Unsorted: !sorted,
			},
			Offset:     int64(page * size),
			PageSize:   size,
			PageNumber: page,
			Paged:      true,
			Unpaged:    false,
		},
		Last:             last,
		TotalPages:       totalPagesInt,
		TotalElements:    totalElements,
		Size:             size,
		Number:           page,
		Sort:             response.SortInfo{Empty: !sorted, Sorted: sorted, Unsorted: !sorted},
		First:            first,
		NumberOfElements: numberOfElements,
		Empty:            empty,
	}, nil
}

// GetPublicBlogByID retrieves a single public blog entry by ID.
func (s *PublicService) GetPublicBlogByID(ctx context.Context, id int64) (*response.BlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapBlogToResponse(ctx, blog)
}

// GetPublicBlogTypes retrieves all blog types.
func (s *PublicService) GetPublicBlogTypes(ctx context.Context) ([]response.BlogTypeResponse, error) {
	items, err := s.blogTypeRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.BlogTypeResponse, len(items))
	for i, item := range items {
		result[i] = mapBlogTypeToResponse(&item)
	}
	return result, nil
}

// GetPublicBlogTypeByID retrieves a blog type by ID.
func (s *PublicService) GetPublicBlogTypeByID(ctx context.Context, id int64) (*response.BlogTypeResponse, error) {
	item, err := s.blogTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapBlogTypeToResponse(item)
	return &resp, nil
}

// GetPublicCourses retrieves all courses.
func (s *PublicService) GetPublicCourses(ctx context.Context) ([]response.CourseResponse, error) {
	items, err := s.courseRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CourseResponse, len(items))
	for i, item := range items {
		result[i] = *mapCourseToResponse(&item)
	}
	return result, nil
}

// GetPublicCertifications retrieves all certifications.
func (s *PublicService) GetPublicCertifications(ctx context.Context) ([]response.CertificationResponse, error) {
	items, err := s.certificationRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CertificationResponse, len(items))
	for i, item := range items {
		result[i] = *mapCertificationToResponse(&item)
	}
	return result, nil
}

// GetPublicLanguages retrieves all languages.
func (s *PublicService) GetPublicLanguages(ctx context.Context) ([]response.LanguageResponse, error) {
	items, err := s.languageRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.LanguageResponse, len(items))
	for i, item := range items {
		result[i] = *mapLanguageToResponse(&item)
	}
	return result, nil
}

// GetPublicReferences retrieves all references.
func (s *PublicService) GetPublicReferences(ctx context.Context) ([]response.ReferenceResponse, error) {
	items, err := s.referenceRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.ReferenceResponse, len(items))
	for i, item := range items {
		result[i] = *mapReferenceToResponse(&item)
	}
	return result, nil
}

// GetPublicCustomSections retrieves only visible custom sections.
func (s *PublicService) GetPublicCustomSections(ctx context.Context) ([]response.CustomSectionResponse, error) {
	all, err := s.customSectionRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var visibleSections []entity.CustomSection
	for _, cs := range all {
		if cs.Visible {
			visibleSections = append(visibleSections, cs)
		}
	}

	result := make([]response.CustomSectionResponse, len(visibleSections))
	for i, cs := range visibleSections {
		result[i] = *mapCustomSectionToResponse(&cs)
	}
	return result, nil
}

// mapBlogToResponse maps a Blog entity to a BlogResponse, resolving nested references.
// TODO: optimize with batch query - currently makes N+1 queries for ImageURL and BlogType
func (s *PublicService) mapBlogToResponse(ctx context.Context, e *entity.Blog) (*response.BlogResponse, error) {
	resp := &response.BlogResponse{
		ID:                 e.ID,
		Title:              e.Title,
		TitleEng:           e.TitleEng,
		CleanURLTitle:      e.CleanURLTitle,
		DescriptionShort:   e.DescriptionShort,
		Description:        e.Description,
		DescriptionShortEng: e.DescriptionShortEng,
		DescriptionEng:     e.DescriptionEng,
	}

	// Resolve ImageURL if present
	if e.ImageURLID != nil {
		img, err := s.imageRepo.GetByID(ctx, *e.ImageURLID)
		if err == nil {
			resp.ImageURL = &response.ImageUrlResponse{
				ID:      img.ID,
				Name:    img.Name,
				NameEng: img.NameEng,
				URL:     img.URL,
			}
		}
	}

	// Resolve BlogType if present
	if e.BlogTypeID != nil {
		bt, err := s.blogTypeRepo.GetByID(ctx, *e.BlogTypeID)
		if err == nil {
			resp.BlogType = &response.BlogTypeResponse{
				ID:      bt.ID,
				Name:    bt.Name,
				NameEng: bt.NameEng,
			}
		}
	}

	// VideoURL is not resolved as VideoUrlRepository is not available

	return resp, nil
}

// mapHomeToResponse maps a Home entity to a HomeResponse, resolving nested references.
func mapHomeToResponse(e *entity.Home, labels *[]entity.Label, imageRepo output.ImageUrlRepository, ctx context.Context) *response.HomeResponse {
	resp := &response.HomeResponse{
		ID:                    e.ID,
		Greeting:              e.Greeting,
		GreetingEng:           e.GreetingEng,
		ButtonWorkLabel:       e.ButtonWorkLabel,
		ButtonWorkLabelEng:    e.ButtonWorkLabelEng,
		ButtonContactLabel:    e.ButtonContactLabel,
		ButtonContactLabelEng: e.ButtonContactLabelEng,
	}

	// Resolve ImageURL if present
	if e.ImageURLID != nil {
		img, err := imageRepo.GetByID(ctx, *e.ImageURLID)
		if err == nil {
			resp.ImageURL = &response.ImageUrlResponse{
				ID:      img.ID,
				Name:    img.Name,
				NameEng: img.NameEng,
				URL:     img.URL,
			}
		}
	}

	// Map all labels as flat slice
	if labels != nil {
		labelResponses := make([]response.LabelResponse, len(*labels))
		for i, l := range *labels {
			labelResponses[i] = response.LabelResponse{
				ID:      l.ID,
				Name:    l.Name,
				NameEng: l.NameEng,
			}
		}
		resp.Labels = labelResponses
	}

	return resp
}

// paginate performs in-memory pagination on a slice of items.
func paginate[T any](items []T, page, size int) ([]T, int, int64) {
	total := int64(len(items))
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}
	start := page * size
	if start >= len(items) {
		return []T{}, totalPages, total
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], totalPages, total
}
