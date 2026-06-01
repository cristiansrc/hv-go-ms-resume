package response

// ReferenceResponse represents the API response for a Reference.
type ReferenceResponse struct {
	ID              int64   `json:"id"`
	FullName        string  `json:"fullName"`
	Position        string  `json:"position"`
	Company         *string `json:"company,omitempty"`
	CompanyEng      *string `json:"companyEng,omitempty"`
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Relationship    *string `json:"relationship,omitempty"`
	RelationshipEng *string `json:"relationshipEng,omitempty"`
}
