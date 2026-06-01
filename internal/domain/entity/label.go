package entity

// Label represents a category label.
type Label struct {
	ID        int64   
	Name      string  
	NameEng   string  
	Order     int     
	CreatedAt string  
	UpdatedAt string  
	DeletedAt *string 
}
