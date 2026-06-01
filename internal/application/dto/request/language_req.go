package request

// LanguageRequest represents the request body for creating/updating a Language entry.
type LanguageRequest struct {
	Language     string `json:"language" validate:"required"`
	LanguageEng  string `json:"languageEng" validate:"required"`
	ReadingLevel string `json:"readingLevel" validate:"required"`
	WritingLevel string `json:"writingLevel" validate:"required"`
	SpeakingLevel string `json:"speakingLevel" validate:"required"`
}
