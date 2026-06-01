package entity

// Language represents a language proficiency entry.
type Language struct {
	ID            int64  
	Language      string 
	LanguageEng   string 
	ReadingLevel  string 
	WritingLevel  string 
	SpeakingLevel string 
	Order         int    
	CreatedAt     string 
	UpdatedAt     string 
	DeletedAt     *string 
}
