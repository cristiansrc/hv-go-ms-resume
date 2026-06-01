package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// ReferenceUseCase defines the operations for Reference.
type ReferenceUseCase interface {
	List(ctx context.Context) ([]response.ReferenceResponse, error)
	GetByID(ctx context.Context, id int64) (*response.ReferenceResponse, error)
	Create(ctx context.Context, req *request.ReferenceRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.ReferenceRequest) error
	Delete(ctx context.Context, id int64) error
}
