package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// SkillTypeUseCase defines the operations for SkillType.
type SkillTypeUseCase interface {
	List(ctx context.Context) ([]response.SkillTypeResponse, error)
	GetByID(ctx context.Context, id int64) (*response.SkillTypeResponse, error)
	Create(ctx context.Context, req *request.SkillTypeRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.SkillTypeRequest) error
	Delete(ctx context.Context, id int64) error
}
