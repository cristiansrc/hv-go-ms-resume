package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// HomeUseCase defines the operations for Home.
type HomeUseCase interface {
	GetByID(ctx context.Context, id int64) (*response.HomeResponse, error)
	Update(ctx context.Context, id int64, req *request.HomeRequest) error
}
