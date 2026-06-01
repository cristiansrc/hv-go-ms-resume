package entity

// CustomSection represents an editable custom section on the resume.
type CustomSection struct {
	ID            int64   
	Title         string  
	TitleEng      string  
	Content       *string 
	ContentEng    *string 
	SummaryPdf    *string 
	SummaryPdfEng *string 
	Order         int     
	Visible       bool    
	CreatedAt     string  
	UpdatedAt     string  
	DeletedAt     *string 
}
