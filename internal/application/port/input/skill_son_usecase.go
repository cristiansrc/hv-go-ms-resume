package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// SkillSonUseCase defines the operations for SkillSon.
type SkillSonUseCase interface {
	List(ctx context.Context) ([]response.SkillSonResponse, error)
	GetByID(ctx context.Context, id int64) (*response.SkillSonResponse, error)
	Create(ctx context.Context, req *request.SkillSonRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.SkillSonRequest) error
	Delete(ctx context.Context, id int64) error
}
