package request

// ExperienceRequest represents the request body for creating/updating an Experience entry.
type ExperienceRequest struct {
	YearStart              string   `json:"yearStart" validate:"required"`
	YearEnd                *string  `json:"yearEnd,omitempty"`
	Company                string   `json:"company" validate:"required"`
	Location               *string  `json:"location,omitempty"`
	LocationEng            *string  `json:"locationEng,omitempty"`
	Position               *string  `json:"position,omitempty"`
	PositionEng            *string  `json:"positionEng,omitempty"`
	Summary                *string  `json:"summary,omitempty"`
	SummaryEng             *string  `json:"summaryEng,omitempty"`
	SummaryPdf             *string  `json:"summaryPdf,omitempty"`
	SummaryPdfEng          *string  `json:"summaryPdfEng,omitempty"`
	DescriptionItemsPdf    []string `json:"descriptionItemsPdf,omitempty"`
	DescriptionItemsPdfEng []string `json:"descriptionItemsPdfEng,omitempty"`
	SkillSonIDs            []int64  `json:"skillSonIds,omitempty"`
}
