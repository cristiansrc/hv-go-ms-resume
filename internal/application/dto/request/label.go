package request

// LabelRequest represents the request body for creating/updating a Label.
type LabelRequest struct {
	Name    string `json:"name" validate:"required"`
	NameEng string `json:"nameEng" validate:"required"`
}
