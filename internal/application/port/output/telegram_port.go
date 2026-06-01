package output

import "context"

// TelegramPort defines the interface for sending Telegram messages.
type TelegramPort interface {
	SendMessage(ctx context.Context, message string) error
}
