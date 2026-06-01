package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// BlogUseCase defines the operations for Blog.
type BlogUseCase interface {
	List(ctx context.Context) ([]response.BlogResponse, error)
	GetByID(ctx context.Context, id int64) (*response.BlogResponse, error)
	Create(ctx context.Context, req *request.BlogRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.BlogRequest) error
	Delete(ctx context.Context, id int64) error
}
