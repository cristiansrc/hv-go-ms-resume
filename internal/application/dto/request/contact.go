package request

// ContactRequest represents the request body for the contact form.
type ContactRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required"`
	Altcha  string `json:"altcha" validate:"required"`
}
