package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// ImageUrlUseCase defines the operations for ImageUrl.
type ImageUrlUseCase interface {
	List(ctx context.Context) ([]response.ImageUrlResponse, error)
	GetByID(ctx context.Context, id int64) (*response.ImageUrlResponse, error)
	Create(ctx context.Context, req *request.ImageURLRequest) (*response.CreateResponse, error)
	Delete(ctx context.Context, id int64) error
}
