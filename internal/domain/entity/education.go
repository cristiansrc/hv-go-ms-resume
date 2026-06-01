package entity

// Education represents an education entry.
type Education struct {
	ID           int64   
	Institution  string  
	Area         string  
	AreaEng      string  
	Degree       string  
	DegreeEng    string  
	StartDate    string  
	EndDate      *string 
	Location     string  
	LocationEng  string  
	Highlights   *string 
	HighlightsEng *string 
	Order        int     
	CreatedAt    string  
	UpdatedAt    string  
	DeletedAt    *string 
}
