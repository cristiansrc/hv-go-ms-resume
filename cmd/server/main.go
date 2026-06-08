package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/service"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/altcha"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/client"
	adapterhttp "github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/handler"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/http/middleware"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/adapter/repository"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := config.NewLogger(cfg.LogLevel)
	logger.Info("starting hv-go-ms-resume", "version", cfg.Version, "port", cfg.Port)

	// Initialize database
	db, err := repository.NewDB(cfg.DatabasePath, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize JWT service
	jwtService := middleware.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration())

	// Initialize Altcha provider
	altchaProvider := altcha.NewProvider("")

	// Initialize repositories
	repos := &repositories{
		basicData:       repository.NewBasicDataRepo(db),
		home:            repository.NewHomeRepo(db),
		label:           repository.NewLabelRepo(db),
		imageUrl:        repository.NewImageUrlRepo(db),
		videoUrl:        repository.NewVideoUrlRepo(db),
		blog:            repository.NewBlogRepo(db),
		blogType:        repository.NewBlogTypeRepo(db),
		skillType:       repository.NewSkillTypeRepo(db),
		skill:           repository.NewSkillRepo(db),
		skillSon:        repository.NewSkillSonRepo(db),
		experience:      repository.NewExperienceRepo(db),
		education:       repository.NewEducationRepo(db),
		futuredProject:  repository.NewFuturedProjectRepo(db),
		course:          repository.NewCourseRepo(db),
		certification:   repository.NewCertificationRepo(db),
		language:        repository.NewLanguageRepo(db),
		reference:       repository.NewReferenceRepo(db),
		customSection:   repository.NewCustomSectionRepo(db),
		userCredentials: repository.NewUserCredentialsRepo(db),
		pdfCache:        repository.NewPdfCacheRepo(db),
	}

	// Initialize external clients
	var s3Client output.S3Port
	var renderCVClient output.RenderCVPort
	var telegramClient output.TelegramPort

	s3Client, err = client.NewS3Client(cfg.AWSRegion, cfg.AWSAccessKey, cfg.AWSSecretKey, cfg.AWSBucket)
	if err != nil {
		logger.Warn("failed to initialize S3 client, PDF caching will be disabled", "error", err)
		s3Client = nil
	}

	renderCVClient = client.NewRenderCVClient(cfg.RenderCVURL)

	telegramClient = client.NewTelegramClient(cfg.TelegramBotToken, cfg.TelegramChatID, logger)

	// Initialize use cases (services)
	authUseCase := service.NewAuthService(
		repos.userCredentials,
		altchaProvider,
		jwtService,
	)

	basicDataUseCase := service.NewBasicDataService(repos.basicData)
	homeUseCase := service.NewHomeService(repos.home, repos.label, repos.imageUrl)
	labelUseCase := service.NewLabelService(repos.label)
	imageUrlUseCase := service.NewImageUrlService(repos.imageUrl)
	videoUrlUseCase := service.NewVideoUrlService(repos.videoUrl)
	blogUseCase := service.NewBlogService(repos.blog, repos.imageUrl, repos.videoUrl, repos.blogType)
	blogTypeUseCase := service.NewBlogTypeService(repos.blogType)
	skillTypeUseCase := service.NewSkillTypeService(repos.skillType)
	skillUseCase := service.NewSkillService(repos.skill)
	skillSonUseCase := service.NewSkillSonService(repos.skillSon)
	experienceUseCase := service.NewExperienceService(repos.experience)
	educationUseCase := service.NewEducationService(repos.education)
	futuredProjectUseCase := service.NewFuturedProjectService(repos.futuredProject)
	courseUseCase := service.NewCourseService(repos.course)
	certificationUseCase := service.NewCertificationService(repos.certification)
	languageUseCase := service.NewLanguageService(repos.language)
	referenceUseCase := service.NewReferenceService(repos.reference)
	customSectionUseCase := service.NewCustomSectionService(repos.customSection)

	pdfUseCase := service.NewPdfService(
		repos.basicData,
		repos.experience,
		repos.education,
		repos.course,
		repos.certification,
		repos.language,
		repos.reference,
		repos.customSection,
		repos.pdfCache,
		renderCVClient,
		s3Client,
		logger,
	)

	contactUseCase := service.NewContactService(altchaProvider, telegramClient, logger)

	publicUseCase := service.NewPublicService(
		repos.basicData,
		repos.home,
		repos.imageUrl,
		repos.label,
		repos.blog,
		repos.blogType,
		repos.experience,
		repos.education,
		repos.course,
		repos.certification,
		repos.language,
		repos.reference,
		repos.customSection,
		repos.skill,
		repos.skillSon,
		repos.videoUrl,
		altchaProvider,
	)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(
		cfg.RateLimitPublic,
		cfg.RateLimitLogin,
		cfg.RateLimitAuthenticated,
	)
	corsMiddleware := middleware.NewCORSMiddleware()
	recoveryMiddleware := middleware.NewRecoveryMiddleware(logger)
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	requestIDMiddleware := middleware.NewRequestIDMiddleware()

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(cfg.Version)
	authHandler := handler.NewAuthHandler(authUseCase, altchaProvider, jwtService)
	basicDataHandler := handler.NewBasicDataHandler(basicDataUseCase)
	homeHandler := handler.NewHomeHandler(homeUseCase)
	labelHandler := handler.NewLabelHandler(labelUseCase)
	imageUrlHandler := handler.NewImageUrlHandler(imageUrlUseCase)
	videoUrlHandler := handler.NewVideoUrlHandler(videoUrlUseCase)
	blogHandler := handler.NewBlogHandler(blogUseCase)
	blogTypeHandler := handler.NewBlogTypeHandler(blogTypeUseCase)
	skillTypeHandler := handler.NewSkillTypeHandler(skillTypeUseCase)
	skillHandler := handler.NewSkillHandler(skillUseCase)
	skillSonHandler := handler.NewSkillSonHandler(skillSonUseCase)
	experienceHandler := handler.NewExperienceHandler(experienceUseCase)
	educationHandler := handler.NewEducationHandler(educationUseCase)
	futuredProjectHandler := handler.NewFuturedProjectHandler(futuredProjectUseCase)
	courseHandler := handler.NewCourseHandler(courseUseCase)
	certificationHandler := handler.NewCertificationHandler(certificationUseCase)
	languageHandler := handler.NewLanguageHandler(languageUseCase)
	referenceHandler := handler.NewReferenceHandler(referenceUseCase)
	customSectionHandler := handler.NewCustomSectionHandler(customSectionUseCase)

	publicHandler := handler.NewPublicHandler(
		publicUseCase,
		pdfUseCase,
		contactUseCase,
		authHandler,
	)

	// Create router
	router := adapterhttp.NewRouter(&adapterhttp.RouterConfig{
		Logger:                logger,
		Version:               cfg.Version,
		HealthHandler:         healthHandler,
		AuthHandler:           authHandler,
		PublicHandler:         publicHandler,
		BasicDataHandler:      basicDataHandler,
		HomeHandler:           homeHandler,
		LabelHandler:          labelHandler,
		ImageUrlHandler:       imageUrlHandler,
		VideoUrlHandler:       videoUrlHandler,
		BlogHandler:           blogHandler,
		BlogTypeHandler:       blogTypeHandler,
		SkillTypeHandler:      skillTypeHandler,
		SkillHandler:          skillHandler,
		SkillSonHandler:       skillSonHandler,
		ExperienceHandler:     experienceHandler,
		EducationHandler:      educationHandler,
		FuturedProjectHandler: futuredProjectHandler,
		CourseHandler:         courseHandler,
		CertificationHandler:  certificationHandler,
		LanguageHandler:       languageHandler,
		ReferenceHandler:      referenceHandler,
		CustomSectionHandler:  customSectionHandler,
		AuthMiddleware:        authMiddleware,
		RateLimit:             rateLimitMiddleware,
		CORS:                  corsMiddleware,
		Recovery:              recoveryMiddleware,
		Logging:               loggingMiddleware,
		RequestID:             requestIDMiddleware,
	})

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", "error", err)
		}
	}()

	logger.Info("server listening", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

type repositories struct {
	basicData       output.BasicDataRepository
	home            output.HomeRepository
	label           output.LabelRepository
	imageUrl        output.ImageUrlRepository
	videoUrl        output.VideoUrlRepository
	blog            output.BlogRepository
	blogType        output.BlogTypeRepository
	skillType       output.SkillTypeRepository
	skill           output.SkillRepository
	skillSon        output.SkillSonRepository
	experience      output.ExperienceRepository
	education       output.EducationRepository
	futuredProject  output.FuturedProjectRepository
	course          output.CourseRepository
	certification   output.CertificationRepository
	language        output.LanguageRepository
	reference       output.ReferenceRepository
	customSection   output.CustomSectionRepository
	userCredentials output.UserCredentialsRepository
	pdfCache        output.PdfCacheRepository
}
