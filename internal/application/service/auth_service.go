package service

import (
	"context"
	"errors"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/response"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements the AuthUseCase interface.
type AuthService struct {
	userRepo   output.UserCredentialsRepository
	altchaPort output.AltchaPort
	jwtPort    output.JWTTokenPort
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo output.UserCredentialsRepository,
	altchaPort output.AltchaPort,
	jwtPort output.JWTTokenPort,
) input.AuthUseCase {
	return &AuthService{
		userRepo:   userRepo,
		altchaPort: altchaPort,
		jwtPort:    jwtPort,
	}
}

// Login authenticates a user and returns a JWT.
func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	if !s.altchaPort.ValidateSolution(req.Altcha) {
		return nil, errors.New("invalid altcha solution")
	}

	creds, err := s.userRepo.GetByUsername(ctx, req.User)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtPort.GenerateToken(req.User)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &response.LoginResponse{Token: token}, nil
}
