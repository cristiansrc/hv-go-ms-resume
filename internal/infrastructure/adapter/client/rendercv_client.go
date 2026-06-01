package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// RenderCVClient implements the RenderCVPort interface.
type RenderCVClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewRenderCVClient creates a new RenderCVClient.
func NewRenderCVClient(baseURL string) *RenderCVClient {
	return &RenderCVClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RenderPDF calls the RenderCV service to generate a PDF.
func (c *RenderCVClient) RenderPDF(ctx context.Context, data *output.CvData, language string, template string) (io.ReadCloser, error) {
	payload := map[string]interface{}{
		"data":     data,
		"language": language,
		"template": template,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RenderCV payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/render", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create RenderCV request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call RenderCV: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("RenderCV returned status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// HealthCheck checks the RenderCV service health.
func (c *RenderCVClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("RenderCV health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RenderCV health check returned status %d", resp.StatusCode)
	}

	return nil
}

var _ output.RenderCVPort = (*RenderCVClient)(nil)
