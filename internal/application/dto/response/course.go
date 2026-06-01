package response

// CourseResponse represents the API response for a Course entry.
type CourseResponse struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	NameEng         string  `json:"nameEng"`
	Institution     string  `json:"institution"`
	InstitutionEng  string  `json:"institutionEng"`
	CompletionDate  string  `json:"completionDate"`
	Description     *string `json:"description,omitempty"`
	DescriptionEng  *string `json:"descriptionEng,omitempty"`
	SummaryPdf      *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng   *string `json:"summaryPdfEng,omitempty"`
	CertificateURL  *string `json:"certificateUrl,omitempty"`
}
