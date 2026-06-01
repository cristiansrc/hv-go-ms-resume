package response

// InfoPageResponse represents the aggregated response for the resume information page.
type InfoPageResponse struct {
	Home           *HomeResponse           `json:"home,omitempty"`
	BasicData      *BasicDataResponse      `json:"basicData,omitempty"`
	Skills         []SkillResponse         `json:"skills,omitempty"`
	Experiences    []ExperienceResponse    `json:"experiences,omitempty"`
	Educations     []EducationResponse     `json:"educations,omitempty"`
	AltchaChallenge *AltchaChallengeResponse `json:"altchaChallenge,omitempty"`
}

// AltchaChallengeResponse represents an Altcha challenge for the client.
type AltchaChallengeResponse struct {
	Algorithm string `json:"algorithm"`
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
	Signature string `json:"signature"`
}
