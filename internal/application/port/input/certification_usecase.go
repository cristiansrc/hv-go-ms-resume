package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// CertificationUseCase defines the operations for Certification.
type CertificationUseCase interface {
	List(ctx context.Context) ([]response.CertificationResponse, error)
	GetByID(ctx context.Context, id int64) (*response.CertificationResponse, error)
	Create(ctx context.Context, req *request.CertificationRequest) (*response.CreateResponse, error)
	Update(ctx context.Context, id int64, req *request.CertificationRequest) error
	Delete(ctx context.Context, id int64) error
}
