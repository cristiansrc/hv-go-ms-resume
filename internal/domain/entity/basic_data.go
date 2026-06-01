package entity

// BasicData represents the main personal information (single row).
type BasicData struct {
	ID                int64   
	FirstName         string  
	OthersName        *string 
	FirstSurname      string  
	OthersSurname     *string 
	DateBirth         string  
	Located           *string 
	LocatedEng        *string 
	StartWorkingDate  *string 
	Greeting          *string 
	GreetingEng       *string 
	Email             string  
	Instagram         *string 
	Linkedin          *string 
	X                 *string 
	Github            *string 
	Description       *string 
	DescriptionEng    *string 
	DescriptionPdf    *string 
	DescriptionPdfEng *string 
	Wrapper           *string 
	WrapperEng        *string 
	CreatedAt         string  
	UpdatedAt         string  
	DeletedAt         *string 
}
