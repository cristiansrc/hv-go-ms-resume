package response

// CustomSectionResponse represents the API response for a CustomSection.
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
