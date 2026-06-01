package request

// FuturedProjectRequest represents the request body for creating/updating a FuturedProject.
type FuturedProjectRequest struct {
	Name               string  `json:"name" validate:"required"`
	NameEng            string  `json:"nameEng" validate:"required"`
	ExperienceID       int64   `json:"experienceId" validate:"required"`
	DescriptionShort   string  `json:"descriptionShort" validate:"required"`
	Description        string  `json:"description" validate:"required"`
	DescriptionShortEng string `json:"descriptionShortEng" validate:"required"`
	DescriptionEng     string  `json:"descriptionEng" validate:"required"`
	ImageListURLID     *int64  `json:"imageListUrlId,omitempty"`
	ImageURLID         *int64  `json:"imageUrlId,omitempty"`
}
