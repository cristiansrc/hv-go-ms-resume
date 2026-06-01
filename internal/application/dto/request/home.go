package request

// HomeRequest represents the request body for creating/updating a Home section.
type HomeRequest struct {
	Greeting              string  `json:"greeting" validate:"required"`
	GreetingEng           string  `json:"greetingEng" validate:"required"`
	ImageURLID            *int64  `json:"imageUrlId,omitempty"`
	ButtonWorkLabel       string  `json:"buttonWorkLabel" validate:"required"`
	ButtonWorkLabelEng    string  `json:"buttonWorkLabelEng" validate:"required"`
	ButtonContactLabel    string  `json:"buttonContactLabel" validate:"required"`
	ButtonContactLabelEng string  `json:"buttonContactLabelEng" validate:"required"`
	LabelIDs              []int64 `json:"labelIds,omitempty"`
}
