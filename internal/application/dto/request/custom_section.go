package request

// CustomSectionRequest represents the request body for creating/updating a CustomSection.
type CustomSectionRequest struct {
	Title        string  `json:"title" validate:"required"`
	TitleEng     string  `json:"titleEng" validate:"required"`
	Content      *string `json:"content,omitempty"`
	ContentEng   *string `json:"contentEng,omitempty"`
	SummaryPdf   *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng *string `json:"summaryPdfEng,omitempty"`
	Visible      bool    `json:"visible" validate:"required"`
}
