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

// BlogService implements the BlogUseCase interface.
type BlogService struct {
	repo        output.BlogRepository
	imageRepo   output.ImageUrlRepository
	videoRepo   output.VideoUrlRepository
	blogTypeRepo output.BlogTypeRepository
}

// NewBlogService creates a new BlogService.
func NewBlogService(
	repo output.BlogRepository,
	imageRepo output.ImageUrlRepository,
	videoRepo output.VideoUrlRepository,
	blogTypeRepo output.BlogTypeRepository,
) input.BlogUseCase {
	return &BlogService{
		repo:         repo,
		imageRepo:    imageRepo,
		videoRepo:    videoRepo,
		blogTypeRepo: blogTypeRepo,
	}
}

// List returns all blog entries.
func (s *BlogService) List(ctx context.Context) ([]response.BlogResponse, error) {
	data, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.BlogResponse, len(data))
	for i, item := range data {
		mapped, err := s.mapBlogToResponse(ctx, &item)
		if err != nil {
			return nil, err
		}
		resp[i] = *mapped
	}
	return resp, nil
}

// GetByID retrieves a Blog by ID.
func (s *BlogService) GetByID(ctx context.Context, id int64) (*response.BlogResponse, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapBlogToResponse(ctx, data)
}

// Create creates a new Blog.
func (s *BlogService) Create(ctx context.Context, req *request.BlogRequest) (*response.CreateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Blog{
		Title:               req.Title,
		TitleEng:            req.TitleEng,
		CleanURLTitle:       req.CleanURLTitle,
		DescriptionShort:    req.DescriptionShort,
		Description:         req.Description,
		DescriptionShortEng: req.DescriptionShortEng,
		DescriptionEng:      req.DescriptionEng,
		ImageURLID:          req.ImageURLID,
		VideoURLID:          req.VideoURLID,
		BlogTypeID:          req.BlogTypeID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	id, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{ID: id}, nil
}

// Update modifies a Blog.
func (s *BlogService) Update(ctx context.Context, id int64, req *request.BlogRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entity := &entity.Blog{
		ID:                  id,
		Title:               req.Title,
		TitleEng:            req.TitleEng,
		CleanURLTitle:       req.CleanURLTitle,
		DescriptionShort:    req.DescriptionShort,
		Description:         req.Description,
		DescriptionShortEng: req.DescriptionShortEng,
		DescriptionEng:      req.DescriptionEng,
		ImageURLID:          req.ImageURLID,
		VideoURLID:          req.VideoURLID,
		BlogTypeID:          req.BlogTypeID,
		CreatedAt:           existing.CreatedAt,
		UpdatedAt:           now,
	}

	return s.repo.Update(ctx, entity)
}

// Delete removes a Blog.
func (s *BlogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *BlogService) mapBlogToResponse(ctx context.Context, e *entity.Blog) (*response.BlogResponse, error) {
	resp := &response.BlogResponse{
		ID:                  e.ID,
		Title:               e.Title,
		TitleEng:            e.TitleEng,
		CleanURLTitle:       e.CleanURLTitle,
		DescriptionShort:    e.DescriptionShort,
		Description:         e.Description,
		DescriptionShortEng: e.DescriptionShortEng,
		DescriptionEng:      e.DescriptionEng,
	}

	// Map image if present (direct FK reference)
	if e.ImageURLID != nil {
		img, err := s.imageRepo.GetByID(ctx, *e.ImageURLID)
		if err == nil {
			mapped := mapImageUrlToResponse(img)
			resp.ImageURL = &mapped
		}
	}

	// Map video if present (direct FK reference)
	if e.VideoURLID != nil {
		video, err := s.videoRepo.GetByID(ctx, *e.VideoURLID)
		if err == nil {
			mapped := mapVideoUrlToResponse(video)
			resp.VideoURL = &mapped
		}
	}

	// Map blog type if present (direct FK reference)
	if e.BlogTypeID != nil {
		bt, err := s.blogTypeRepo.GetByID(ctx, *e.BlogTypeID)
		if err == nil {
			mapped := mapBlogTypeToResponse(bt)
			resp.BlogType = &mapped
		}
	}

	return resp, nil
}
