package request

// LoginRequest represents the login request body.
type LoginRequest struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Altcha   string `json:"altcha" validate:"required"`
}
