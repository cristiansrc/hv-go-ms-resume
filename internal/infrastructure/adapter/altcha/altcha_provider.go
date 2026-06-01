package altcha

import (
	"encoding/json"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
)

// Provider implements the AltchaPort interface.
type Provider struct {
	hmacKey string
}

// NewProvider creates a new Altcha provider.
func NewProvider(hmacKey string) *Provider {
	if hmacKey == "" {
		hmacKey = "default-altcha-hmac-key-change-in-production"
	}
	return &Provider{hmacKey: hmacKey}
}

// GenerateChallenge creates a new Altcha challenge.
func (p *Provider) GenerateChallenge() (*output.AltchaChallenge, error) {
	challenge, err := altcha.CreateChallenge(altcha.ChallengeOptions{
		HMACKey:   p.hmacKey,
		MaxNumber: 100000,
	})
	if err != nil {
		return nil, err
	}

	return &output.AltchaChallenge{
		Algorithm: string(challenge.Algorithm),
		Challenge: challenge.Challenge,
		Salt:      challenge.Salt,
		Signature: challenge.Signature,
	}, nil
}

// ValidateSolution verifies an Altcha solution.
// The solution is expected to be a JSON-encoded Altcha payload from the client.
func (p *Provider) ValidateSolution(solution string) bool {
	var payload altcha.Payload
	if err := json.Unmarshal([]byte(solution), &payload); err != nil {
		return false
	}

	ok, err := altcha.VerifySolution(payload, p.hmacKey, false)
	if err != nil {
		return false
	}
	return ok
}
