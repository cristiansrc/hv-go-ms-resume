package response

// HomeResponse represents the API response for the Home section.
type HomeResponse struct {
	ID                   int64            `json:"id"`
	Greeting             string           `json:"greeting"`
	GreetingEng          string           `json:"greetingEng"`
	ImageURL             *ImageUrlResponse `json:"imageUrl,omitempty"`
	ButtonWorkLabel      string           `json:"buttonWorkLabel"`
	ButtonWorkLabelEng   string           `json:"buttonWorkLabelEng"`
	ButtonContactLabel   string           `json:"buttonContactLabel"`
	ButtonContactLabelEng string          `json:"buttonContactLabelEng"`
	Labels               []LabelResponse  `json:"labels,omitempty"`
}
