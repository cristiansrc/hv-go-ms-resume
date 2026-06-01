package request

// EducationRequest represents the request body for creating/updating an Education entry.
type EducationRequest struct {
	Institution    string   `json:"institution" validate:"required"`
	Area           string   `json:"area" validate:"required"`
	AreaEng        string   `json:"areaEng" validate:"required"`
	Degree         string   `json:"degree" validate:"required"`
	DegreeEng      string   `json:"degreeEng" validate:"required"`
	StartDate      string   `json:"startDate" validate:"required"`
	EndDate        *string  `json:"endDate,omitempty"`
	Location       string   `json:"location" validate:"required"`
	LocationEng    string   `json:"locationEng" validate:"required"`
	Highlights     []string `json:"highlights,omitempty"`
	HighlightsEng  []string `json:"highlightsEng,omitempty"`
}
