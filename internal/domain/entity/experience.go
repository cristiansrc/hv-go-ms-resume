package entity

// Experience represents a work experience entry.
type Experience struct {
	ID                     int64   
	YearStart              string  
	YearEnd                *string 
	Company                string  
	Location               *string 
	LocationEng            *string 
	Position               *string 
	PositionEng            *string 
	Summary                *string 
	SummaryEng             *string 
	SummaryPdf             *string 
	SummaryPdfEng          *string 
	DescriptionItemsPdf    *string 
	DescriptionItemsPdfEng *string 
	CreatedAt              string  
	UpdatedAt              string  
	DeletedAt              *string 
}
