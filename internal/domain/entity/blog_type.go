package entity

// BlogType represents a blog category.
type BlogType struct {
	ID        int64   
	Name      string  
	NameEng   string  
	Order     int     
	CreatedAt string  
	UpdatedAt string  
	DeletedAt *string 
}
