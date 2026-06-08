package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/dto/request"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authUseCase  input.AuthUseCase
	altchaPort   output.AltchaPort
	jwtPort      output.JWTTokenPort
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(
	authUseCase input.AuthUseCase,
	altchaPort output.AltchaPort,
	jwtPort output.JWTTokenPort,
) *AuthHandler {
	return &AuthHandler{
		authUseCase:  authUseCase,
		altchaPort:   altchaPort,
		jwtPort:      jwtPort,
	}
}

// Login handles POST /v1/ms-resume/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "invalid altcha") || strings.Contains(msg, "invalid credentials") {
			WriteError(w, r, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "Invalid credentials or Altcha")
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetChallenge handles GET /v1/ms-resume/public/challenge
func (h *AuthHandler) GetChallenge(w http.ResponseWriter, r *http.Request) {
	challenge, err := h.altchaPort.GenerateChallenge()
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate challenge")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}
