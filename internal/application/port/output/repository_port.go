package output

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// BasicDataRepository defines the persistence contract for BasicData.
type BasicDataRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.BasicData, error)
	Update(ctx context.Context, data *entity.BasicData) error
}

// HomeRepository defines the persistence contract for Home.
type HomeRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Home, error)
	Update(ctx context.Context, data *entity.Home) error
}

// LabelRepository defines the persistence contract for Label.
type LabelRepository interface {
	List(ctx context.Context) ([]entity.Label, error)
	GetByID(ctx context.Context, id int64) (*entity.Label, error)
	Create(ctx context.Context, data *entity.Label) (int64, error)
	Update(ctx context.Context, data *entity.Label) error
	Delete(ctx context.Context, id int64) error
}

// ImageUrlRepository defines the persistence contract for ImageUrl.
type ImageUrlRepository interface {
	List(ctx context.Context) ([]entity.ImageUrl, error)
	GetByID(ctx context.Context, id int64) (*entity.ImageUrl, error)
	Create(ctx context.Context, data *entity.ImageUrl) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// VideoUrlRepository defines the persistence contract for VideoUrl.
type VideoUrlRepository interface {
	List(ctx context.Context) ([]entity.VideoUrl, error)
	GetByID(ctx context.Context, id int64) (*entity.VideoUrl, error)
	Create(ctx context.Context, data *entity.VideoUrl) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// BlogRepository defines the persistence contract for Blog.
type BlogRepository interface {
	List(ctx context.Context) ([]entity.Blog, error)
	GetByID(ctx context.Context, id int64) (*entity.Blog, error)
	Create(ctx context.Context, data *entity.Blog) (int64, error)
	Update(ctx context.Context, data *entity.Blog) error
	Delete(ctx context.Context, id int64) error
}

// BlogTypeRepository defines the persistence contract for BlogType.
type BlogTypeRepository interface {
	List(ctx context.Context) ([]entity.BlogType, error)
	GetByID(ctx context.Context, id int64) (*entity.BlogType, error)
	Create(ctx context.Context, data *entity.BlogType) (int64, error)
	Update(ctx context.Context, data *entity.BlogType) error
	Delete(ctx context.Context, id int64) error
}

// SkillTypeRepository defines the persistence contract for SkillType.
type SkillTypeRepository interface {
	List(ctx context.Context) ([]entity.SkillType, error)
	GetByID(ctx context.Context, id int64) (*entity.SkillType, error)
	Create(ctx context.Context, data *entity.SkillType) (int64, error)
	Update(ctx context.Context, data *entity.SkillType) error
	Delete(ctx context.Context, id int64) error
}

// SkillRepository defines the persistence contract for Skill.
type SkillRepository interface {
	List(ctx context.Context) ([]entity.Skill, error)
	GetByID(ctx context.Context, id int64) (*entity.Skill, error)
	Create(ctx context.Context, data *entity.Skill) (int64, error)
	Update(ctx context.Context, data *entity.Skill) error
	Delete(ctx context.Context, id int64) error
}

// SkillSonRepository defines the persistence contract for SkillSon.
type SkillSonRepository interface {
	List(ctx context.Context) ([]entity.SkillSon, error)
	GetByID(ctx context.Context, id int64) (*entity.SkillSon, error)
	ListBySkillID(ctx context.Context, skillID int64) ([]entity.SkillSon, error)
	Create(ctx context.Context, data *entity.SkillSon) (int64, error)
	Update(ctx context.Context, data *entity.SkillSon) error
	Delete(ctx context.Context, id int64) error
}

// ExperienceRepository defines the persistence contract for Experience.
type ExperienceRepository interface {
	List(ctx context.Context) ([]entity.Experience, error)
	GetByID(ctx context.Context, id int64) (*entity.Experience, error)
	Create(ctx context.Context, data *entity.Experience) (int64, error)
	Update(ctx context.Context, data *entity.Experience) error
	Delete(ctx context.Context, id int64) error
}

// EducationRepository defines the persistence contract for Education.
type EducationRepository interface {
	List(ctx context.Context) ([]entity.Education, error)
	GetByID(ctx context.Context, id int64) (*entity.Education, error)
	Create(ctx context.Context, data *entity.Education) (int64, error)
	Update(ctx context.Context, data *entity.Education) error
	Delete(ctx context.Context, id int64) error
}

// FuturedProjectRepository defines the persistence contract for FuturedProject.
type FuturedProjectRepository interface {
	List(ctx context.Context) ([]entity.FuturedProject, error)
	GetByID(ctx context.Context, id int64) (*entity.FuturedProject, error)
	Create(ctx context.Context, data *entity.FuturedProject) (int64, error)
	Update(ctx context.Context, data *entity.FuturedProject) error
	Delete(ctx context.Context, id int64) error
}

// CourseRepository defines the persistence contract for Course.
type CourseRepository interface {
	List(ctx context.Context) ([]entity.Course, error)
	GetByID(ctx context.Context, id int64) (*entity.Course, error)
	Create(ctx context.Context, data *entity.Course) (int64, error)
	Update(ctx context.Context, data *entity.Course) error
	Delete(ctx context.Context, id int64) error
}

// CertificationRepository defines the persistence contract for Certification.
type CertificationRepository interface {
	List(ctx context.Context) ([]entity.Certification, error)
	GetByID(ctx context.Context, id int64) (*entity.Certification, error)
	Create(ctx context.Context, data *entity.Certification) (int64, error)
	Update(ctx context.Context, data *entity.Certification) error
	Delete(ctx context.Context, id int64) error
}

// LanguageRepository defines the persistence contract for Language.
type LanguageRepository interface {
	List(ctx context.Context) ([]entity.Language, error)
	GetByID(ctx context.Context, id int64) (*entity.Language, error)
	Create(ctx context.Context, data *entity.Language) (int64, error)
	Update(ctx context.Context, data *entity.Language) error
	Delete(ctx context.Context, id int64) error
}

// ReferenceRepository defines the persistence contract for Reference.
type ReferenceRepository interface {
	List(ctx context.Context) ([]entity.Reference, error)
	GetByID(ctx context.Context, id int64) (*entity.Reference, error)
	Create(ctx context.Context, data *entity.Reference) (int64, error)
	Update(ctx context.Context, data *entity.Reference) error
	Delete(ctx context.Context, id int64) error
}

// CustomSectionRepository defines the persistence contract for CustomSection.
type CustomSectionRepository interface {
	List(ctx context.Context) ([]entity.CustomSection, error)
	GetByID(ctx context.Context, id int64) (*entity.CustomSection, error)
	Create(ctx context.Context, data *entity.CustomSection) (int64, error)
	Update(ctx context.Context, data *entity.CustomSection) error
	Delete(ctx context.Context, id int64) error
}

// UserCredentialsRepository defines the persistence contract for admin credentials.
type UserCredentialsRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.UserCredentials, error)
}

// PdfCacheRepository defines the persistence contract for PDF cache.
type PdfCacheRepository interface {
	GetByLanguageAndTemplate(ctx context.Context, language, template string) (*entity.PdfCache, error)
	Upsert(ctx context.Context, data *entity.PdfCache) error
	Delete(ctx context.Context, id int64) error
}
