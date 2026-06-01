package response

// BlogResponse represents the API response for a Blog entry.
type BlogResponse struct {
	ID                 int64              `json:"id"`
	Title              string             `json:"title"`
	TitleEng           string             `json:"titleEng"`
	CleanURLTitle      *string            `json:"cleanUrlTitle,omitempty"`
	DescriptionShort   string             `json:"descriptionShort"`
	Description        string             `json:"description"`
	DescriptionShortEng string            `json:"descriptionShortEng"`
	DescriptionEng     string             `json:"descriptionEng"`
	ImageURL           *ImageUrlResponse  `json:"imageUrl,omitempty"`
	VideoURL           *VideoUrlResponse  `json:"videoUrl,omitempty"`
	BlogType           *BlogTypeResponse  `json:"blogType,omitempty"`
}
