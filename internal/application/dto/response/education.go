package response

// EducationResponse represents the API response for an Education entry.
type EducationResponse struct {
	ID            int64    `json:"id"`
	Institution   string   `json:"institution"`
	Area          string   `json:"area"`
	AreaEng       string   `json:"areaEng"`
	Degree        string   `json:"degree"`
	DegreeEng     string   `json:"degreeEng"`
	StartDate     string   `json:"startDate"`
	EndDate       *string  `json:"endDate,omitempty"`
	Location      string   `json:"location"`
	LocationEng   string   `json:"locationEng"`
	Highlights    []string `json:"highlights,omitempty"`
	HighlightsEng []string `json:"highlightsEng,omitempty"`
}
