package entity

// UserCredentials represents admin credentials (single row).
type UserCredentials struct {
	ID           int64  
	Username     string 
	PasswordHash string 
	CreatedAt    string 
	UpdatedAt    string 
}
