package request

// SkillRequest represents the request body for creating/updating a Skill.
type SkillRequest struct {
	Name        string  `json:"name" validate:"required"`
	NameEng     string  `json:"nameEng" validate:"required"`
	SkillSonIDs []int64 `json:"skillSonIds,omitempty"`
}
