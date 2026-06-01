package response

// VideoUrlResponse represents the API response for a VideoUrl.
type VideoUrlResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	NameEng string `json:"nameEng"`
	URL     string `json:"url"`
}
