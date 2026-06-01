package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// FuturedProjectUseCase defines the operations for FuturedProject.
type FuturedProjectUseCase interface {
	List(ctx context.Context) ([]response.FuturedProjectResponse, error)
	GetByID(ctx context.Context, id int64) (*response.FuturedProjectResponse, error)
	Create(ctx context.Context, req *request.FuturedProjectRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.FuturedProjectRequest) error
	Delete(ctx context.Context, id int64) error
}
