package output

// JWTClaims represents the claims extracted from a validated JWT.
type JWTClaims struct {
	Username string
}

// JWTTokenPort defines the interface for JWT token operations.
type JWTTokenPort interface {
	GenerateToken(username string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}
