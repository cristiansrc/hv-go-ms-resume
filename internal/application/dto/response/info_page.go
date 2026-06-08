package response

// InfoPageResponse represents the aggregated response for the resume information page.
type InfoPageResponse struct {
	Home            *HomeResponse              `json:"home,omitempty"`
	BasicData       *BasicDataResponse         `json:"basicData,omitempty"`
	Skills          []SkillResponse            `json:"skills"`
	Experiences     []ExperienceResponse       `json:"experiences"`
	Educations      []EducationResponse        `json:"educations"`
	AltchaChallenge *AltchaChallengeResponse   `json:"altchaChallenge,omitempty"`
	Courses         []CourseResponse           `json:"courses"`
	Certifications  []CertificationResponse    `json:"certifications"`
	Languages       []LanguageResponse         `json:"languages"`
	References      []ReferenceResponse        `json:"references"`
	CustomSections  []CustomSectionResponse    `json:"customSections"`
}

// AltchaChallengeResponse represents an Altcha challenge for the client.
type AltchaChallengeResponse struct {
	Algorithm string `json:"algorithm"`
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
	Signature string `json:"signature"`
}
