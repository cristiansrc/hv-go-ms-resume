package entity

// Course represents a completed course.
type Course struct {
	ID              int64   
	Name            string  
	NameEng         string  
	Institution     string  
	InstitutionEng  string  
	CompletionDate  string  
	Description     *string 
	DescriptionEng  *string 
	SummaryPdf      *string 
	SummaryPdfEng   *string 
	CertificateURL  *string 
	Order           int     
	CreatedAt       string  
	UpdatedAt       string  
	DeletedAt       *string 
}
