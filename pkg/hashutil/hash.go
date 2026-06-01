package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ComputeHash computes a deterministic SHA-256 hash from the given data.
// The data is serialized to JSON to ensure deterministic output.
func ComputeHash(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data for hashing: %w", err)
	}

	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}
