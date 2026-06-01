package request

// CourseRequest represents the request body for creating/updating a Course.
type CourseRequest struct {
	Name           string  `json:"name" validate:"required"`
	NameEng        string  `json:"nameEng" validate:"required"`
	Institution    string  `json:"institution" validate:"required"`
	InstitutionEng string  `json:"institutionEng" validate:"required"`
	CompletionDate string  `json:"completionDate" validate:"required"`
	Description    *string `json:"description,omitempty"`
	DescriptionEng *string `json:"descriptionEng,omitempty"`
	SummaryPdf     *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng  *string `json:"summaryPdfEng,omitempty"`
	CertificateURL *string `json:"certificateUrl,omitempty"`
}
