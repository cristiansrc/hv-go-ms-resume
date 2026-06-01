package entity

// Blog represents a blog article.
type Blog struct {
	ID                 int64   
	Title              string  
	TitleEng           string  
	CleanURLTitle      *string 
	DescriptionShort   string  
	Description        string  
	DescriptionShortEng string 
	DescriptionEng     string  
	ImageURLID         *int64  
	VideoURLID         *int64  
	BlogTypeID         *int64  
	CreatedAt          string  
	UpdatedAt          string  
	DeletedAt          *string 
}
