package input

import (
	"context"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
)

// AuthUseCase defines the authentication operations.
type AuthUseCase interface {
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
}
