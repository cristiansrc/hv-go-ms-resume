package response

// SkillSonResponse represents the API response for a SkillSon.
type SkillSonResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	NameEng string `json:"nameEng"`
}
