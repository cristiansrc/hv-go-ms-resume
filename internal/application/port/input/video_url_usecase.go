package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// VideoUrlUseCase defines the operations for VideoUrl.
type VideoUrlUseCase interface {
	List(ctx context.Context) ([]response.VideoUrlResponse, error)
	GetByID(ctx context.Context, id int64) (*response.VideoUrlResponse, error)
	Create(ctx context.Context, req *request.VideoURLRequest) (*response.CreateResponse, error)
	Delete(ctx context.Context, id int64) error
}
