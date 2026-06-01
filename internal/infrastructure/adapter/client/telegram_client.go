package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// TelegramClient implements the TelegramPort interface.
type TelegramClient struct {
	botToken string
	chatID   string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewTelegramClient creates a new TelegramClient.
func NewTelegramClient(botToken, chatID string, logger *slog.Logger) *TelegramClient {
	return &TelegramClient{
		botToken: botToken,
		chatID:   chatID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// SendMessage sends a message to the configured Telegram chat.
func (c *TelegramClient) SendMessage(ctx context.Context, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)

	payload := map[string]string{
		"chat_id":    c.chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Telegram payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create Telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("telegram send error", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("telegram returned non-OK", "status", resp.StatusCode)
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

	return nil
}

var _ output.TelegramPort = (*TelegramClient)(nil)
