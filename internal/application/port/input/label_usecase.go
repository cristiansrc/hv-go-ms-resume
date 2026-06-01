package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// LabelUseCase defines the operations for Label.
type LabelUseCase interface {
	List(ctx context.Context) ([]response.LabelResponse, error)
	GetByID(ctx context.Context, id int64) (*response.LabelResponse, error)
	Create(ctx context.Context, req *request.LabelRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.LabelRequest) error
	Delete(ctx context.Context, id int64) error
}
