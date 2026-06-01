package entity

// Certification represents a professional certification.
type Certification struct {
	ID                     int64   
	Name                   string  
	NameEng                string  
	IssuingOrganization    string  
	IssuingOrganizationEng string  
	IssueDate              string  
	ExpirationDate         *string 
	VerificationURL        *string 
	CredentialID           *string 
	Description            *string 
	DescriptionEng         *string 
	SummaryPdf             *string 
	SummaryPdfEng          *string 
	Order                  int     
	CreatedAt              string  
	UpdatedAt              string  
	DeletedAt              *string 
}
