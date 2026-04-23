package security

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
)

// ImprovedTokenConfig holds improved token configuration
type ImprovedTokenConfig struct {
	SecretKey string
	Issuer    string
	TokenTTL  time.Duration
}

// TokenClaims represents token claims
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Issuer    string    `json:"iss"`
}

// NewImprovedTokenConfig creates a new improved token configuration
func NewImprovedTokenConfig(secretKey, issuer string) *ImprovedTokenConfig {
	if secretKey == "" {
		secretKey = generateSecureToken(64)
	}
	if issuer == "" {
		issuer = "ai-parking-system"
	}
	return &ImprovedTokenConfig{
		SecretKey: secretKey,
		Issuer:    issuer,
		TokenTTL:  24 * time.Hour,
	}
}

// GenerateToken generates an improved token with user information
func (t *ImprovedTokenConfig) GenerateToken(userID, role string) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(t.TokenTTL),
		Issuer:    t.Issuer,
	}

	// Create token signature
	signature := t.createSignature(claims)

	// Encode claims and signature
	tokenData := base64.URLEncoding.EncodeToString([]byte(userID + "|" + role + "|" + claims.IssuedAt.Format(time.RFC3339) + "|" + signature))

	return tokenData, nil
}

// ValidateToken validates an improved token
func (t *ImprovedTokenConfig) ValidateToken(tokenString string) (*TokenClaims, error) {
	// Decode token
	tokenData, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, errors.New("invalid token encoding")
	}

	// Parse token parts (simplified format: userID|role|issuedAt|signature)
	parts := splitString(string(tokenData), "|")
	if len(parts) < 4 {
		return nil, errors.New("invalid token format")
	}

	issuedAt, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, errors.New("invalid timestamp")
	}

	claims := TokenClaims{
		UserID:    parts[0],
		Role:      parts[1],
		IssuedAt:  issuedAt,
		ExpiresAt: issuedAt.Add(t.TokenTTL),
		Issuer:    t.Issuer,
	}

	// Check expiration
	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	// Verify signature
	expectedSignature := t.createSignature(claims)
	if parts[3] != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	return &claims, nil
}

// createSignature creates a signature for the token claims
func (t *ImprovedTokenConfig) createSignature(claims TokenClaims) string {
	data := claims.UserID + claims.Role + claims.IssuedAt.Format(time.RFC3339) + t.SecretKey
	hash := sha256.Sum256([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])[:32]
}

// splitString is a helper function to split string
func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}
