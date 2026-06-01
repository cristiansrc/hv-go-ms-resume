package request

// SkillTypeRequest represents the request body for creating/updating a SkillType.
type SkillTypeRequest struct {
	Name     string  `json:"name" validate:"required"`
	NameEng  string  `json:"nameEng" validate:"required"`
	SkillIDs []int64 `json:"skillIds,omitempty"`
}
