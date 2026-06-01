package response

// SkillResponse represents the API response for a Skill with nested SkillSons.
type SkillResponse struct {
	ID      int64              `json:"id"`
	Name    string             `json:"name"`
	NameEng string             `json:"nameEng"`
	Sons    []SkillSonResponse `json:"sons,omitempty"`
}
