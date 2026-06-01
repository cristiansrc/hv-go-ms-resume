package response

// ImageUrlResponse represents the API response for an ImageUrl.
type ImageUrlResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	NameEng string `json:"nameEng"`
	URL     string `json:"url"`
}
