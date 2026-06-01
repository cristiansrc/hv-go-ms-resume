package response

// CertificationResponse represents the API response for a Certification.
type CertificationResponse struct {
	ID                     int64   `json:"id"`
	Name                   string  `json:"name"`
	NameEng                string  `json:"nameEng"`
	IssuingOrganization    string  `json:"issuingOrganization"`
	IssuingOrganizationEng string  `json:"issuingOrganizationEng"`
	IssueDate              string  `json:"issueDate"`
	ExpirationDate         *string `json:"expirationDate,omitempty"`
	VerificationURL        *string `json:"verificationUrl,omitempty"`
	CredentialID           *string `json:"credentialId,omitempty"`
	Description            *string `json:"description,omitempty"`
	DescriptionEng         *string `json:"descriptionEng,omitempty"`
	SummaryPdf             *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng          *string `json:"summaryPdfEng,omitempty"`
}
