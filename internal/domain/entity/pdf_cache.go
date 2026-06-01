package entity

// PdfCache represents a cached PDF entry.
type PdfCache struct {
	ID        int64  
	Language  string 
	Template  string 
	DataHash  string 
	S3Key     string 
	CreatedAt string 
	FileSize  *int64 
}
