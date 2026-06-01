package request

// ReferenceRequest represents the request body for creating/updating a Reference.
type ReferenceRequest struct {
	FullName        string  `json:"fullName" validate:"required"`
	Position        string  `json:"position" validate:"required"`
	Company         *string `json:"company,omitempty"`
	CompanyEng      *string `json:"companyEng,omitempty"`
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Relationship    *string `json:"relationship,omitempty"`
	RelationshipEng *string `json:"relationshipEng,omitempty"`
}
