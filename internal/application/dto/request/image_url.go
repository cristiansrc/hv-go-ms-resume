package request

// ImageURLRequest represents the request body for creating/updating an ImageURL.
type ImageURLRequest struct {
	Name    string `json:"name" validate:"required"`
	NameEng string `json:"nameEng" validate:"required"`
	File    string `json:"file" validate:"required"`
}
