package response

// SkillTypeResponse represents the API response for a SkillType with nested Skills.
type SkillTypeResponse struct {
	ID      int64           `json:"id"`
	Name    string          `json:"name"`
	NameEng string          `json:"nameEng"`
	Skills  []SkillResponse `json:"skills,omitempty"`
}
