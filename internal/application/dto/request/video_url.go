package request

// VideoURLRequest represents the request body for creating/updating a VideoURL.
type VideoURLRequest struct {
	Name    string `json:"name" validate:"required"`
	NameEng string `json:"nameEng" validate:"required"`
	URL     string `json:"url" validate:"required"`
}
