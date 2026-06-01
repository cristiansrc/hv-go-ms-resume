package request

// BlogRequest represents the request body for creating/updating a Blog entry.
type BlogRequest struct {
	Title              string  `json:"title" validate:"required"`
	TitleEng           string  `json:"titleEng" validate:"required"`
	CleanURLTitle      *string `json:"cleanUrlTitle,omitempty"`
	DescriptionShort   string  `json:"descriptionShort" validate:"required"`
	Description        string  `json:"description" validate:"required"`
	DescriptionShortEng string `json:"descriptionShortEng" validate:"required"`
	DescriptionEng     string  `json:"descriptionEng" validate:"required"`
	ImageURLID         *int64  `json:"imageUrlId,omitempty"`
	VideoURLID         *int64  `json:"videoUrlId,omitempty"`
	BlogTypeID         *int64  `json:"blogTypeId,omitempty"`
}
