package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// BlogTypeUseCase defines the operations for BlogType.
type BlogTypeUseCase interface {
	List(ctx context.Context) ([]response.BlogTypeResponse, error)
	GetByID(ctx context.Context, id int64) (*response.BlogTypeResponse, error)
	Create(ctx context.Context, req *request.BlogTypeRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.BlogTypeRequest) error
	Delete(ctx context.Context, id int64) error
}
