package response

// LanguageResponse represents the API response for a Language proficiency entry.
type LanguageResponse struct {
	ID            int64  `json:"id"`
	Language      string `json:"language"`
	LanguageEng   string `json:"languageEng"`
	ReadingLevel  string `json:"readingLevel"`
	WritingLevel  string `json:"writingLevel"`
	SpeakingLevel string `json:"speakingLevel"`
}
