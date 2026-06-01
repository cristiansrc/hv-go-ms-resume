package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// SkillUseCase defines the operations for Skill.
type SkillUseCase interface {
	List(ctx context.Context) ([]response.SkillResponse, error)
	GetByID(ctx context.Context, id int64) (*response.SkillResponse, error)
	Create(ctx context.Context, req *request.SkillRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.SkillRequest) error
	Delete(ctx context.Context, id int64) error
}
