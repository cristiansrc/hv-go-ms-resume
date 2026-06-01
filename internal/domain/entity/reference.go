package entity

// Reference represents a professional reference.
type Reference struct {
	ID              int64   
	FullName        string  
	Position        string  
	Company         *string 
	CompanyEng      *string 
	Email           *string 
	Phone           *string 
	Relationship    *string 
	RelationshipEng *string 
	Order           int     
	CreatedAt       string  
	UpdatedAt       string  
	DeletedAt       *string 
}
