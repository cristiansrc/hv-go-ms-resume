package request

// CertificationRequest represents the request body for creating/updating a Certification.
type CertificationRequest struct {
	Name                   string  `json:"name" validate:"required"`
	NameEng                string  `json:"nameEng" validate:"required"`
	IssuingOrganization    string  `json:"issuingOrganization" validate:"required"`
	IssuingOrganizationEng string  `json:"issuingOrganizationEng" validate:"required"`
	IssueDate              string  `json:"issueDate" validate:"required"`
	ExpirationDate         *string `json:"expirationDate,omitempty"`
	VerificationURL        *string `json:"verificationUrl,omitempty"`
	CredentialID           *string `json:"credentialId,omitempty"`
	Description            *string `json:"description,omitempty"`
	DescriptionEng         *string `json:"descriptionEng,omitempty"`
	SummaryPdf             *string `json:"summaryPdf,omitempty"`
	SummaryPdfEng          *string `json:"summaryPdfEng,omitempty"`
}
