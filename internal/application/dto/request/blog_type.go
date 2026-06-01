package request

// BlogTypeRequest represents the request body for creating/updating a BlogType.
type BlogTypeRequest struct {
	Name    string `json:"name" validate:"required"`
	NameEng string `json:"nameEng" validate:"required"`
}
