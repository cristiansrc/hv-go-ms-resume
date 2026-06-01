package response

// LabelResponse represents the API response for a Label.
type LabelResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	NameEng string `json:"nameEng"`
}
