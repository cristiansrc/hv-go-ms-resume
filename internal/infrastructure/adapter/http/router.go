package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/handler"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/middleware"
)

// RouterConfig holds all dependencies for the router.
type RouterConfig struct {
	Logger                *slog.Logger
	Version               string
	HealthHandler         *handler.HealthHandler
	AuthHandler           *handler.AuthHandler
	PublicHandler         *handler.PublicHandler
	BasicDataHandler      *handler.BasicDataHandler
	HomeHandler           *handler.HomeHandler
	LabelHandler          *handler.LabelHandler
	ImageUrlHandler       *handler.ImageUrlHandler
	VideoUrlHandler       *handler.VideoUrlHandler
	BlogHandler           *handler.BlogHandler
	BlogTypeHandler       *handler.BlogTypeHandler
	SkillTypeHandler      *handler.SkillTypeHandler
	SkillHandler          *handler.SkillHandler
	SkillSonHandler       *handler.SkillSonHandler
	ExperienceHandler     *handler.ExperienceHandler
	EducationHandler      *handler.EducationHandler
	FuturedProjectHandler *handler.FuturedProjectHandler
	CourseHandler         *handler.CourseHandler
	CertificationHandler  *handler.CertificationHandler
	LanguageHandler       *handler.LanguageHandler
	ReferenceHandler      *handler.ReferenceHandler
	CustomSectionHandler  *handler.CustomSectionHandler
	AuthMiddleware        *middleware.AuthMiddleware
	RateLimit             *middleware.RateLimitMiddleware
	CORS                  *middleware.CORSMiddleware
	Recovery              *middleware.RecoveryMiddleware
	Logging               *middleware.LoggingMiddleware
	RequestID             *middleware.RequestIDMiddleware
}

// NewRouter creates and configures the Chi router with all routes and middleware.
func NewRouter(cfg *RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(cfg.CORS.CORS)
	r.Use(cfg.Recovery.Recover)
	r.Use(cfg.RequestID.AttachRequestID)
	r.Use(cfg.Logging.Log)
	r.Use(chimw.RealIP)

	// Health - no rate limiting
	r.Get("/health", cfg.HealthHandler.ServeHTTP)

	// Public routes with rate limiting
	r.Route("/v1/ms-resume", func(r chi.Router) {
		// Public endpoints (no auth)
		r.Group(func(r chi.Router) {
			r.Use(cfg.RateLimit.RateLimit("public"))

			r.Get("/public/info-page", cfg.PublicHandler.GetInfoPage)
			r.Get("/public/curriculum/{language}", cfg.PublicHandler.GetCurriculum)
			r.Post("/public/contact", cfg.PublicHandler.SubmitContact)
			r.Get("/public/challenge", cfg.PublicHandler.GetChallenge)
			r.Get("/public/blog", cfg.PublicHandler.GetPublicBlog)
			r.Get("/public/blog/{id}", cfg.PublicHandler.GetPublicBlogByID)
			r.Get("/public/blog-type", cfg.PublicHandler.GetPublicBlogTypes)
			r.Get("/public/blog-type/{id}", cfg.PublicHandler.GetPublicBlogTypeByID)
			r.Get("/public/templates", cfg.PublicHandler.GetTemplates)
			r.Get("/public/languages", cfg.PublicHandler.GetLanguages)
			r.Get("/public/courses", cfg.PublicHandler.GetPublicCourses)
			r.Get("/public/certifications", cfg.PublicHandler.GetPublicCertifications)
			r.Get("/public/languages-data", cfg.PublicHandler.GetPublicLanguages)
			r.Get("/public/references", cfg.PublicHandler.GetPublicReferences)
			r.Get("/public/custom-sections", cfg.PublicHandler.GetPublicCustomSections)
		})

		// Login with strict rate limiting
		r.Group(func(r chi.Router) {
			r.Use(cfg.RateLimit.RateLimit("login"))
			r.Post("/login", cfg.AuthHandler.Login)
		})

		// Protected routes (JWT required)
		r.Group(func(r chi.Router) {
			r.Use(cfg.AuthMiddleware.Authenticate)
			r.Use(cfg.RateLimit.RateLimit("auth"))

			// BasicData
			r.Get("/basic-data/{id}", cfg.BasicDataHandler.GetByID)
			r.Put("/basic-data/{id}", cfg.BasicDataHandler.Update)

			// Home
			r.Get("/home/{id}", cfg.HomeHandler.GetByID)
			r.Put("/home/{id}", cfg.HomeHandler.Update)

			// Label
			r.Get("/label", cfg.LabelHandler.List)
			r.Get("/label/{id}", cfg.LabelHandler.GetByID)
			r.Post("/label", cfg.LabelHandler.Create)
			r.Put("/label/{id}", cfg.LabelHandler.Update)
			r.Delete("/label/{id}", cfg.LabelHandler.Delete)

			// ImageUrl
			r.Get("/image-url", cfg.ImageUrlHandler.List)
			r.Get("/image-url/{id}", cfg.ImageUrlHandler.GetByID)
			r.Post("/image-url", cfg.ImageUrlHandler.Create)
			r.Delete("/image-url/{id}", cfg.ImageUrlHandler.Delete)

			// VideoUrl
			r.Get("/video-url", cfg.VideoUrlHandler.List)
			r.Get("/video-url/{id}", cfg.VideoUrlHandler.GetByID)
			r.Post("/video-url", cfg.VideoUrlHandler.Create)
			r.Delete("/video-url/{id}", cfg.VideoUrlHandler.Delete)

			// Blog
			r.Get("/blog", cfg.BlogHandler.List)
			r.Get("/blog/{id}", cfg.BlogHandler.GetByID)
			r.Post("/blog", cfg.BlogHandler.Create)
			r.Put("/blog/{id}", cfg.BlogHandler.Update)
			r.Delete("/blog/{id}", cfg.BlogHandler.Delete)

			// BlogType
			r.Get("/blog-type", cfg.BlogTypeHandler.List)
			r.Get("/blog-type/{id}", cfg.BlogTypeHandler.GetByID)
			r.Post("/blog-type", cfg.BlogTypeHandler.Create)
			r.Put("/blog-type/{id}", cfg.BlogTypeHandler.Update)
			r.Delete("/blog-type/{id}", cfg.BlogTypeHandler.Delete)

			// SkillType
			r.Get("/skill-type", cfg.SkillTypeHandler.List)
			r.Get("/skill-type/{id}", cfg.SkillTypeHandler.GetByID)
			r.Post("/skill-type", cfg.SkillTypeHandler.Create)
			r.Put("/skill-type/{id}", cfg.SkillTypeHandler.Update)
			r.Delete("/skill-type/{id}", cfg.SkillTypeHandler.Delete)

			// Skill
			r.Get("/skill", cfg.SkillHandler.List)
			r.Get("/skill/{id}", cfg.SkillHandler.GetByID)
			r.Post("/skill", cfg.SkillHandler.Create)
			r.Put("/skill/{id}", cfg.SkillHandler.Update)
			r.Delete("/skill/{id}", cfg.SkillHandler.Delete)

			// SkillSon
			r.Get("/skill-son", cfg.SkillSonHandler.List)
			r.Get("/skill-son/{id}", cfg.SkillSonHandler.GetByID)
			r.Post("/skill-son", cfg.SkillSonHandler.Create)
			r.Put("/skill-son/{id}", cfg.SkillSonHandler.Update)
			r.Delete("/skill-son/{id}", cfg.SkillSonHandler.Delete)

			// Experience
			r.Get("/experience", cfg.ExperienceHandler.List)
			r.Get("/experience/{id}", cfg.ExperienceHandler.GetByID)
			r.Post("/experience", cfg.ExperienceHandler.Create)
			r.Put("/experience/{id}", cfg.ExperienceHandler.Update)
			r.Delete("/experience/{id}", cfg.ExperienceHandler.Delete)

			// Education
			r.Get("/education", cfg.EducationHandler.List)
			r.Get("/education/{id}", cfg.EducationHandler.GetByID)
			r.Post("/education", cfg.EducationHandler.Create)
			r.Put("/education/{id}", cfg.EducationHandler.Update)
			r.Delete("/education/{id}", cfg.EducationHandler.Delete)

			// FuturedProject
			r.Get("/futured-project", cfg.FuturedProjectHandler.List)
			r.Get("/futured-project/{id}", cfg.FuturedProjectHandler.GetByID)
			r.Post("/futured-project", cfg.FuturedProjectHandler.Create)
			r.Put("/futured-project/{id}", cfg.FuturedProjectHandler.Update)
			r.Delete("/futured-project/{id}", cfg.FuturedProjectHandler.Delete)

			// Course
			r.Get("/course", cfg.CourseHandler.List)
			r.Get("/course/{id}", cfg.CourseHandler.GetByID)
			r.Post("/course", cfg.CourseHandler.Create)
			r.Put("/course/{id}", cfg.CourseHandler.Update)
			r.Delete("/course/{id}", cfg.CourseHandler.Delete)

			// Certification
			r.Get("/certification", cfg.CertificationHandler.List)
			r.Get("/certification/{id}", cfg.CertificationHandler.GetByID)
			r.Post("/certification", cfg.CertificationHandler.Create)
			r.Put("/certification/{id}", cfg.CertificationHandler.Update)
			r.Delete("/certification/{id}", cfg.CertificationHandler.Delete)

			// Language
			r.Get("/language", cfg.LanguageHandler.List)
			r.Get("/language/{id}", cfg.LanguageHandler.GetByID)
			r.Post("/language", cfg.LanguageHandler.Create)
			r.Put("/language/{id}", cfg.LanguageHandler.Update)
			r.Delete("/language/{id}", cfg.LanguageHandler.Delete)

			// Reference
			r.Get("/reference", cfg.ReferenceHandler.List)
			r.Get("/reference/{id}", cfg.ReferenceHandler.GetByID)
			r.Post("/reference", cfg.ReferenceHandler.Create)
			r.Put("/reference/{id}", cfg.ReferenceHandler.Update)
			r.Delete("/reference/{id}", cfg.ReferenceHandler.Delete)

			// CustomSection
			r.Get("/custom-section", cfg.CustomSectionHandler.List)
			r.Get("/custom-section/{id}", cfg.CustomSectionHandler.GetByID)
			r.Post("/custom-section", cfg.CustomSectionHandler.Create)
			r.Put("/custom-section/{id}", cfg.CustomSectionHandler.Update)
			r.Delete("/custom-section/{id}", cfg.CustomSectionHandler.Delete)
		})
	})

	return r
}
