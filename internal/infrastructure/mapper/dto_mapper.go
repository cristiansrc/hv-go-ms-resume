package mapper

import (
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

// ToCustomSectionResponse converts a domain CustomSection to the API response format
// mapping visible INTEGER to boolean.
type CustomSectionResponse struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	TitleEng      string  `json:"titleEng"`
	Content       *string `json:"content,omitempty"`
	ContentEng    *string `json:"contentEng,omitempty"`
	SummaryPdf    *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng *string `json:"summaryPdfEng,omitempty"`
	Visible       bool    `json:"visible"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

// MapCustomSectionToResponse converts domain entity to API response.
func MapCustomSectionToResponse(e *entity.CustomSection) *CustomSectionResponse {
	if e == nil {
		return nil
	}
	return &CustomSectionResponse{
		ID:            e.ID,
		Title:         e.Title,
		TitleEng:      e.TitleEng,
		Content:       e.Content,
		ContentEng:    e.ContentEng,
		SummaryPdf:    e.SummaryPdf,
		SummaryPdfEng: e.SummaryPdfEng,
		Visible:       e.Visible,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// MapCustomSectionsToResponse converts a slice of domain entities.
func MapCustomSectionsToResponse(entities []entity.CustomSection) []CustomSectionResponse {
	result := make([]CustomSectionResponse, len(entities))
	for i, e := range entities {
		result[i] = *MapCustomSectionToResponse(&e)
	}
	return result
}
