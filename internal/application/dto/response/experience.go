package response

// ExperienceResponse represents the API response for an Experience entry.
type ExperienceResponse struct {
	ID                     int64              `json:"id"`
	YearStart              string             `json:"yearStart"`
	YearEnd                *string            `json:"yearEnd,omitempty"`
	Company                string             `json:"company"`
	Location               *string            `json:"location,omitempty"`
	LocationEng            *string            `json:"locationEng,omitempty"`
	Position               *string            `json:"position,omitempty"`
	PositionEng            *string            `json:"positionEng,omitempty"`
	Summary                *string            `json:"summary,omitempty"`
	SummaryEng             *string            `json:"summaryEng,omitempty"`
	SummaryPdf             *string            `json:"summaryPdf,omitempty"`
	SummaryPdfEng          *string            `json:"summaryPdfEng,omitempty"`
	DescriptionItemsPdf    []string           `json:"descriptionItemsPdf,omitempty"`
	DescriptionItemsPdfEng []string           `json:"descriptionItemsPdfEng,omitempty"`
	SkillSons              []SkillSonResponse `json:"skillSons,omitempty"`
}
