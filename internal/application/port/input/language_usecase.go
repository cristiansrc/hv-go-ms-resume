package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// LanguageUseCase defines the operations for Language.
type LanguageUseCase interface {
	List(ctx context.Context) ([]response.LanguageResponse, error)
	GetByID(ctx context.Context, id int64) (*response.LanguageResponse, error)
	Create(ctx context.Context, req *request.LanguageRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.LanguageRequest) error
	Delete(ctx context.Context, id int64) error
}
