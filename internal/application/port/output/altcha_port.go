package output

// AltchaChallenge represents an Altcha challenge for the client.
type AltchaChallenge struct {
	Algorithm string `json:"algorithm"`
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
	Signature string `json:"signature"`
}

// AltchaPort defines the interface for Altcha challenge generation and validation.
type AltchaPort interface {
	GenerateChallenge() (*AltchaChallenge, error)
	ValidateSolution(solution string) bool
}
