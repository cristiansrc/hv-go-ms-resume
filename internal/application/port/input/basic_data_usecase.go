package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// BasicDataUseCase defines the operations for BasicData.
type BasicDataUseCase interface {
	GetByID(ctx context.Context, id int64) (*response.BasicDataResponse, error)
	Update(ctx context.Context, id int64, req *request.BasicDataRequest) error
}
