package service

import (
	"context"
	"testing"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type mockBasicDataRepo struct {
	output.BasicDataRepository
	getByIDFn func(ctx context.Context, id int64) (*entity.BasicData, error)
}

func (m *mockBasicDataRepo) GetByID(ctx context.Context, id int64) (*entity.BasicData, error) {
	return m.getByIDFn(ctx, id)
}

type mockHomeRepo struct {
	output.HomeRepository
	getByIDFn func(ctx context.Context, id int64) (*entity.Home, error)
}

func (m *mockHomeRepo) GetByID(ctx context.Context, id int64) (*entity.Home, error) {
	return m.getByIDFn(ctx, id)
}

type mockLabelRepo struct {
	output.LabelRepository
	listFn func(ctx context.Context) ([]entity.Label, error)
}

func (m *mockLabelRepo) List(ctx context.Context) ([]entity.Label, error) {
	return m.listFn(ctx)
}

type mockImageRepo struct {
	output.ImageUrlRepository
	getByIDFn func(ctx context.Context, id int64) (*entity.ImageUrl, error)
}

func (m *mockImageRepo) GetByID(ctx context.Context, id int64) (*entity.ImageUrl, error) {
	return m.getByIDFn(ctx, id)
}

type mockExperienceRepo struct {
	output.ExperienceRepository
	listFn func(ctx context.Context) ([]entity.Experience, error)
}

func (m *mockExperienceRepo) List(ctx context.Context) ([]entity.Experience, error) {
	return m.listFn(ctx)
}

type mockEducationRepo struct {
	output.EducationRepository
	listFn func(ctx context.Context) ([]entity.Education, error)
}

func (m *mockEducationRepo) List(ctx context.Context) ([]entity.Education, error) {
	return m.listFn(ctx)
}

type mockSkillRepo struct {
	output.SkillRepository
	listFn func(ctx context.Context) ([]entity.Skill, error)
}

func (m *mockSkillRepo) List(ctx context.Context) ([]entity.Skill, error) {
	return m.listFn(ctx)
}

type mockSkillSonRepo struct {
	output.SkillSonRepository
	listBySkillIDFn func(ctx context.Context, skillID int64) ([]entity.SkillSon, error)
}

func (m *mockSkillSonRepo) ListBySkillID(ctx context.Context, skillID int64) ([]entity.SkillSon, error) {
	return m.listBySkillIDFn(ctx, skillID)
}

type mockCourseRepo struct {
	output.CourseRepository
	listFn func(ctx context.Context) ([]entity.Course, error)
}

func (m *mockCourseRepo) List(ctx context.Context) ([]entity.Course, error) {
	return m.listFn(ctx)
}

type mockCertificationRepo struct {
	output.CertificationRepository
	listFn func(ctx context.Context) ([]entity.Certification, error)
}

func (m *mockCertificationRepo) List(ctx context.Context) ([]entity.Certification, error) {
	return m.listFn(ctx)
}

type mockLanguageRepo struct {
	output.LanguageRepository
	listFn func(ctx context.Context) ([]entity.Language, error)
}

func (m *mockLanguageRepo) List(ctx context.Context) ([]entity.Language, error) {
	return m.listFn(ctx)
}

type mockReferenceRepo struct {
	output.ReferenceRepository
	listFn func(ctx context.Context) ([]entity.Reference, error)
}

func (m *mockReferenceRepo) List(ctx context.Context) ([]entity.Reference, error) {
	return m.listFn(ctx)
}

type mockCustomSectionRepo struct {
	output.CustomSectionRepository
	listFn func(ctx context.Context) ([]entity.CustomSection, error)
}

func (m *mockCustomSectionRepo) List(ctx context.Context) ([]entity.CustomSection, error) {
	return m.listFn(ctx)
}

type mockAltchaPort struct {
	output.AltchaPort
	generateChallengeFn func() (*output.AltchaChallenge, error)
}

func (m *mockAltchaPort) GenerateChallenge() (*output.AltchaChallenge, error) {
	return m.generateChallengeFn()
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newMockGetInfoPagePublicService(
	basicDataRepo output.BasicDataRepository,
	homeRepo output.HomeRepository,
	imageRepo output.ImageUrlRepository,
	labelRepo output.LabelRepository,
	experienceRepo output.ExperienceRepository,
	educationRepo output.EducationRepository,
	skillRepo output.SkillRepository,
	skillSonRepo output.SkillSonRepository,
	courseRepo output.CourseRepository,
	certificationRepo output.CertificationRepository,
	languageRepo output.LanguageRepository,
	referenceRepo output.ReferenceRepository,
	customSectionRepo output.CustomSectionRepository,
	altchaPort output.AltchaPort,
) *PublicService {
	return &PublicService{
		basicDataRepo:     basicDataRepo,
		homeRepo:          homeRepo,
		imageRepo:         imageRepo,
		labelRepo:         labelRepo,
		experienceRepo:    experienceRepo,
		educationRepo:     educationRepo,
		skillRepo:         skillRepo,
		skillSonRepo:      skillSonRepo,
		courseRepo:        courseRepo,
		certificationRepo: certificationRepo,
		languageRepo:      languageRepo,
		referenceRepo:     referenceRepo,
		customSectionRepo: customSectionRepo,
		altchaPort:        altchaPort,
	}
}

func defaultBasicData() *entity.BasicData {
	return &entity.BasicData{
		ID: 1, FirstName: "John", FirstSurname: "Doe", Email: "john@example.com",
	}
}

func defaultHome() *entity.Home {
	return &entity.Home{ID: 1, Greeting: "Hello", GreetingEng: "Hello EN"}
}

func defaultLabel() []entity.Label {
	return []entity.Label{{ID: 1, Name: "Backend", NameEng: "Backend"}}
}

func defaultExperiences() []entity.Experience {
	return []entity.Experience{
		{ID: 1, YearStart: "2020", Company: "Acme"},
	}
}

func defaultEducations() []entity.Education {
	return []entity.Education{
		{ID: 1, Institution: "MIT", Area: "CS", Degree: "BSc"},
	}
}

func defaultSkills() []entity.Skill {
	return []entity.Skill{{ID: 1, Name: "Go", NameEng: "Go"}}
}

func defaultSkillSons() []entity.SkillSon {
	return []entity.SkillSon{{ID: 1, Name: "Fiber", NameEng: "Fiber"}}
}

// ---------------------------------------------------------------------------
// Tests: GetInfoPage
// ---------------------------------------------------------------------------

func TestGetInfoPage_Success_WithAllNewEntities(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) {
			return []entity.Course{
				{ID: 1, Name: "Go Basics", Institution: "Udemy", CompletionDate: "2024-01-01"},
			}, nil
		},
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) {
			return []entity.Certification{
				{ID: 1, Name: "AWS Certified", IssuingOrganization: "Amazon", IssueDate: "2024-06-01"},
			}, nil
		},
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) {
			return []entity.Language{
				{ID: 1, Language: "English", ReadingLevel: "C1", WritingLevel: "C1", SpeakingLevel: "C1"},
			}, nil
		},
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) {
			return []entity.Reference{
				{ID: 1, FullName: "Jane Smith", Position: "Manager", Company: strPtr("Acme")},
			}, nil
		},
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) {
			return []entity.CustomSection{
				{ID: 1, Title: "Projects", Visible: true},
				{ID: 2, Title: "Hidden", Visible: false},
				{ID: 3, Title: "Volunteer", Visible: true},
			}, nil
		},
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	// Verify the 5 new fields are present
	if len(resp.Courses) != 1 {
		t.Errorf("expected 1 course, got %d", len(resp.Courses))
	}
	if len(resp.Certifications) != 1 {
		t.Errorf("expected 1 certification, got %d", len(resp.Certifications))
	}
	if len(resp.Languages) != 1 {
		t.Errorf("expected 1 language, got %d", len(resp.Languages))
	}
	if len(resp.References) != 1 {
		t.Errorf("expected 1 reference, got %d", len(resp.References))
	}
	// CustomSection: only visible ones (2 out of 3)
	if len(resp.CustomSections) != 2 {
		t.Errorf("expected 2 custom sections (visible only), got %d", len(resp.CustomSections))
	}
	if resp.CustomSections[0].Title != "Projects" {
		t.Errorf("expected first custom section 'Projects', got '%s'", resp.CustomSections[0].Title)
	}
	if resp.CustomSections[1].Title != "Volunteer" {
		t.Errorf("expected second custom section 'Volunteer', got '%s'", resp.CustomSections[1].Title)
	}

	// Verify existing fields are still present
	if resp.Home == nil {
		t.Error("expected home to be present")
	}
	if resp.BasicData == nil {
		t.Error("expected basicData to be present")
	}
	if resp.Skills == nil {
		t.Error("expected skills to be present")
	}
	if resp.Experiences == nil {
		t.Error("expected experiences to be present")
	}
	if resp.Educations == nil {
		t.Error("expected educations to be present")
	}
	if resp.AltchaChallenge == nil {
		t.Error("expected altchaChallenge to be present")
	}
}

func TestGetInfoPage_EmptyTables_ReturnsEmptySlices(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	// All new repos return empty lists
	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) { return []entity.CustomSection{}, nil },
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	// All new fields should be empty slices, not nil
	if resp.Courses == nil {
		t.Error("expected non-nil Courses slice (empty)")
	}
	if resp.Certifications == nil {
		t.Error("expected non-nil Certifications slice (empty)")
	}
	if resp.Languages == nil {
		t.Error("expected non-nil Languages slice (empty)")
	}
	if resp.References == nil {
		t.Error("expected non-nil References slice (empty)")
	}
	if resp.CustomSections == nil {
		t.Error("expected non-nil CustomSections slice (empty)")
	}
	if len(resp.Courses) != 0 {
		t.Errorf("expected 0 courses, got %d", len(resp.Courses))
	}
	if len(resp.Certifications) != 0 {
		t.Errorf("expected 0 certifications, got %d", len(resp.Certifications))
	}
}

func TestGetInfoPage_ErrorInCertificationRepo_ReturnsError(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) {
			return nil, assertAnError("cert repo error")
		},
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) { return []entity.CustomSection{}, nil },
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err == nil {
		t.Fatal("expected error when certification repo fails, got nil")
	}
	if resp != nil {
		t.Error("expected nil response when error occurs")
	}
}

func TestGetInfoPage_ErrorInLanguageRepo_ReturnsError(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) {
			return nil, assertAnError("lang repo error")
		},
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) { return []entity.CustomSection{}, nil },
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err == nil {
		t.Fatal("expected error when language repo fails, got nil")
	}
	if resp != nil {
		t.Error("expected nil response when error occurs")
	}
}

func TestGetInfoPage_ErrorInReferenceRepo_ReturnsError(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) {
			return nil, assertAnError("ref repo error")
		},
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) { return []entity.CustomSection{}, nil },
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err == nil {
		t.Fatal("expected error when reference repo fails, got nil")
	}
	if resp != nil {
		t.Error("expected nil response when error occurs")
	}
}

func TestGetInfoPage_ErrorInCourseRepo_ReturnsError(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) {
			return nil, assertAnError("course repo error")
		},
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) { return []entity.CustomSection{}, nil },
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err == nil {
		t.Fatal("expected error when course repo fails, got nil")
	}
	if resp != nil {
		t.Error("expected nil response when error occurs")
	}
}

func TestGetInfoPage_ErrorInCustomSectionRepo_ReturnsError(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) {
			return nil, assertAnError("custom section repo error")
		},
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err == nil {
		t.Fatal("expected error when custom section repo fails, got nil")
	}
	if resp != nil {
		t.Error("expected nil response when error occurs")
	}
}

func TestGetInfoPage_CustomSection_OnlyVisible(t *testing.T) {
	altcha := &mockAltchaPort{
		generateChallengeFn: func() (*output.AltchaChallenge, error) {
			return &output.AltchaChallenge{Algorithm: "SHA-256", Challenge: "abc", Salt: "salt", Signature: "sig"}, nil
		},
	}

	courseRepo := &mockCourseRepo{
		listFn: func(ctx context.Context) ([]entity.Course, error) { return []entity.Course{}, nil },
	}
	certRepo := &mockCertificationRepo{
		listFn: func(ctx context.Context) ([]entity.Certification, error) { return []entity.Certification{}, nil },
	}
	langRepo := &mockLanguageRepo{
		listFn: func(ctx context.Context) ([]entity.Language, error) { return []entity.Language{}, nil },
	}
	refRepo := &mockReferenceRepo{
		listFn: func(ctx context.Context) ([]entity.Reference, error) { return []entity.Reference{}, nil },
	}
	// Only invisible custom sections → result should have 0 custom sections
	csRepo := &mockCustomSectionRepo{
		listFn: func(ctx context.Context) ([]entity.CustomSection, error) {
			return []entity.CustomSection{
				{ID: 1, Title: "Hidden1", Visible: false},
				{ID: 2, Title: "Hidden2", Visible: false},
			}, nil
		},
	}

	svc := newMockGetInfoPagePublicService(
		&mockBasicDataRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.BasicData, error) { return defaultBasicData(), nil }},
		&mockHomeRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Home, error) { return defaultHome(), nil }},
		&mockImageRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.ImageUrl, error) { return nil, nil }},
		&mockLabelRepo{listFn: func(ctx context.Context) ([]entity.Label, error) { return defaultLabel(), nil }},
		&mockExperienceRepo{listFn: func(ctx context.Context) ([]entity.Experience, error) { return defaultExperiences(), nil }},
		&mockEducationRepo{listFn: func(ctx context.Context) ([]entity.Education, error) { return defaultEducations(), nil }},
		&mockSkillRepo{listFn: func(ctx context.Context) ([]entity.Skill, error) { return defaultSkills(), nil }},
		&mockSkillSonRepo{listBySkillIDFn: func(ctx context.Context, skillID int64) ([]entity.SkillSon, error) { return defaultSkillSons(), nil }},
		courseRepo, certRepo, langRepo, refRepo, csRepo, altcha,
	)

	resp, err := svc.GetInfoPage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if len(resp.CustomSections) != 0 {
		t.Errorf("expected 0 custom sections (all invisible), got %d", len(resp.CustomSections))
	}
	// Should be an empty slice, not nil
	if resp.CustomSections == nil {
		t.Error("expected non-nil CustomSections slice (empty)")
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string {
	return &s
}

// assertAnError returns a simple error for testing error paths.
type assertAnError string

func (e assertAnError) Error() string {
	return string(e)
}
