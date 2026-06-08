package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/input"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// ErrInvalidAltcha is a sentinel error for invalid altcha solution.
var ErrInvalidAltcha = errors.New("invalid altcha solution")

// ContactService implements the ContactUseCase interface.
type ContactService struct {
	altchaPort   output.AltchaPort
	telegramPort output.TelegramPort
	logger       *slog.Logger
}

// NewContactService creates a new ContactService.
func NewContactService(
	altchaPort output.AltchaPort,
	telegramPort output.TelegramPort,
	logger *slog.Logger,
) input.ContactUseCase {
	return &ContactService{
		altchaPort:   altchaPort,
		telegramPort: telegramPort,
		logger:       logger,
	}
}

// SubmitContact validates the Altcha challenge and sends a Telegram notification.
func (s *ContactService) SubmitContact(ctx context.Context, name, email, message, altcha string) error {
	if !s.altchaPort.ValidateSolution(altcha) {
		return ErrInvalidAltcha
	}

	telegramMessage := fmt.Sprintf(
		"<b>New Contact Message</b>\n\n<b>Name:</b> %s\n<b>Email:</b> %s\n<b>Message:</b>\n%s",
		name, email, message,
	)

	// Telegram error is logged but not returned to the user
	if err := s.telegramPort.SendMessage(ctx, telegramMessage); err != nil {
		s.logger.Error("failed to send telegram notification", "error", err)
	}

	return nil
}
