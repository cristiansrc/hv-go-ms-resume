package request

// SkillSonRequest represents the request body for creating/updating a SkillSon.
type SkillSonRequest struct {
	Name    string `json:"name" validate:"required"`
	NameEng string `json:"nameEng" validate:"required"`
}
