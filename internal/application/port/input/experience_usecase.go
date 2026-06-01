package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// ExperienceUseCase defines the operations for Experience.
type ExperienceUseCase interface {
	List(ctx context.Context) ([]response.ExperienceResponse, error)
	GetByID(ctx context.Context, id int64) (*response.ExperienceResponse, error)
	Create(ctx context.Context, req *request.ExperienceRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.ExperienceRequest) error
	Delete(ctx context.Context, id int64) error
}
