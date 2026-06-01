package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// PublicUseCase defines the public read-only operations.
type PublicUseCase interface {
	GetInfoPage(ctx context.Context) (*response.InfoPageResponse, error)
	GetPublicBlogPage(ctx context.Context, page, size int, sort string) (*response.BlogPageResponse, error)
	GetPublicBlogByID(ctx context.Context, id int64) (*response.BlogResponse, error)
	GetPublicBlogTypes(ctx context.Context) ([]response.BlogTypeResponse, error)
	GetPublicBlogTypeByID(ctx context.Context, id int64) (*response.BlogTypeResponse, error)
	GetPublicCourses(ctx context.Context) ([]response.CourseResponse, error)
	GetPublicCertifications(ctx context.Context) ([]response.CertificationResponse, error)
	GetPublicLanguages(ctx context.Context) ([]response.LanguageResponse, error)
	GetPublicReferences(ctx context.Context) ([]response.ReferenceResponse, error)
	GetPublicCustomSections(ctx context.Context) ([]response.CustomSectionResponse, error)
}
