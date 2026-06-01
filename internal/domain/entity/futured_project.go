package entity

// FuturedProject represents a featured project entry.
type FuturedProject struct {
	ID                int64   
	Name              string  
	NameEng           string  
	DescriptionShort  string  
	Description       string  
	DescriptionShortEng string 
	DescriptionEng    string  
	ExperienceID      int64   
	ImageListURLID    *int64  
	ImageURLID        *int64  
	Order             int     
	CreatedAt         string  
	UpdatedAt         string  
	DeletedAt         *string 
}
