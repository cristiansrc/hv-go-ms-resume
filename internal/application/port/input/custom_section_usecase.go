package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// CustomSectionUseCase defines the operations for CustomSection.
type CustomSectionUseCase interface {
	List(ctx context.Context) ([]response.CustomSectionResponse, error)
	GetByID(ctx context.Context, id int64) (*response.CustomSectionResponse, error)
	Create(ctx context.Context, req *request.CustomSectionRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.CustomSectionRequest) error
	Delete(ctx context.Context, id int64) error
}
