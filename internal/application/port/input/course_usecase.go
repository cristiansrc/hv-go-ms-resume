package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// CourseUseCase defines the operations for Course.
type CourseUseCase interface {
	List(ctx context.Context) ([]response.CourseResponse, error)
	GetByID(ctx context.Context, id int64) (*response.CourseResponse, error)
	Create(ctx context.Context, req *request.CourseRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.CourseRequest) error
	Delete(ctx context.Context, id int64) error
}
