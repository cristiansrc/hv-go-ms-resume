package output

import (
	"context"
	"io"
)

// CvData represents the data sent to RenderCV for PDF generation.
type CvData struct {
	Name        string                `json:"name"`
	Email       string                `json:"email"`
	Location    string                `json:"location"`
	Headline    string                `json:"headline"`
	SocialLinks []SocialNetwork       `json:"social_links"`
	Sections    map[string][]CVSection `json:"sections"`
}

// SocialNetwork represents a social media link.
type SocialNetwork struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	URL      string `json:"url"`
}

// CVSection represents a section in the CV.
type CVSection struct {
	Name     string            `json:"name"`
	Items    []CVSectionItem   `json:"items"`
}

// CVSectionItem represents a single item within a CV section.
type CVSectionItem struct {
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Date        string   `json:"date,omitempty"`
	Description string   `json:"description,omitempty"`
	Highlights  []string `json:"highlights,omitempty"`
}

// RenderCVPort defines the interface for calling the RenderCV service.
type RenderCVPort interface {
	RenderPDF(ctx context.Context, data *CvData, language string, template string) (io.ReadCloser, error)
	HealthCheck(ctx context.Context) error
}
