package entity

// Home represents the home page content (single row).
type Home struct {
	ID                   int64   
	Greeting             string  
	GreetingEng          string  
	ImageURLID           *int64  
	ButtonWorkLabel      string  
	ButtonWorkLabelEng   string  
	ButtonContactLabel   string  
	ButtonContactLabelEng string 
	CreatedAt            string  
	UpdatedAt            string  
	DeletedAt            *string 
}
