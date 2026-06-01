package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// EducationUseCase defines the operations for Education.
type EducationUseCase interface {
	List(ctx context.Context) ([]response.EducationResponse, error)
	GetByID(ctx context.Context, id int64) (*response.EducationResponse, error)
	Create(ctx context.Context, req *request.EducationRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.EducationRequest) error
	Delete(ctx context.Context, id int64) error
}
