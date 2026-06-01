package input

import "context"

// ContactUseCase defines the contact form operations.
type ContactUseCase interface {
	SubmitContact(ctx context.Context, name, email, message, altcha string) error
}
