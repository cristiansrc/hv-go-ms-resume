package response

// FuturedProjectResponse represents the API response for a FuturedProject.
type FuturedProjectResponse struct {
	ID                int64               `json:"id"`
	Name              string              `json:"name"`
	NameEng           string              `json:"nameEng"`
	DescriptionShort  string              `json:"descriptionShort"`
	Description       string              `json:"description"`
	DescriptionShortEng string             `json:"descriptionShortEng"`
	DescriptionEng    string              `json:"descriptionEng"`
	Experience        *ExperienceResponse `json:"experience,omitempty"`
	ImageListURL      *ImageUrlResponse   `json:"imageListUrl,omitempty"`
	ImageURL          *ImageUrlResponse   `json:"imageUrl,omitempty"`
}
