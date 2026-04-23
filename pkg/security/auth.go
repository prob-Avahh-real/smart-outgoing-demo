package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AdminToken  string
	ExpiresAt   time.Time
	Environment string
}

// NewAuthConfig creates a new authentication configuration
func NewAuthConfig(adminToken string, environment string) *AuthConfig {
	// Debug logging
	fmt.Printf("[DEBUG] NewAuthConfig called with adminToken=%s, environment=%s\n", adminToken, environment)

	// Generate token if not provided or using default in production environment
	if adminToken == "" {
		adminToken = generateSecureToken(32)
		fmt.Printf("[DEBUG] Generated new token (was empty)\n")
	}
	// Only replace demo_admin_token in production for security
	if (adminToken == "demo_admin_token") && (environment == "prod" || environment == "production") {
		adminToken = generateSecureToken(32)
		fmt.Printf("[DEBUG] Replaced demo_admin_token with secure token in production\n")
	}

	// Set TTL based on environment
	var ttl time.Duration
	switch environment {
	case "prod", "production":
		ttl = 1 * time.Hour
	case "test", "testing":
		ttl = 8 * time.Hour
	default:
		ttl = 24 * time.Hour
	}

	fmt.Printf("[DEBUG] Final AdminToken=%s, Environment=%s, TTL=%v\n", adminToken, environment, ttl)

	return &AuthConfig{
		AdminToken:  adminToken,
		ExpiresAt:   time.Now().Add(ttl),
		Environment: environment,
	}
}

// ValidateToken validates the provided token
func (ac *AuthConfig) ValidateToken(token string) error {
	if token == "" {
		return errors.New("token required")
	}

	if time.Now().After(ac.ExpiresAt) {
		return errors.New("token expired")
	}

	if token != ac.AdminToken {
		return errors.New("invalid token")
	}

	return nil
}

// GetToken returns the admin token
func (ac *AuthConfig) GetToken() string {
	return ac.AdminToken
}

// GetExpiry returns the token expiry time
func (ac *AuthConfig) GetExpiry() time.Time {
	return ac.ExpiresAt
}

// IsExpired checks if the token is expired
func (ac *AuthConfig) IsExpired() bool {
	return time.Now().After(ac.ExpiresAt)
}

// GetTTL returns the remaining time until expiry
func (ac *AuthConfig) GetTTL() time.Duration {
	return time.Until(ac.ExpiresAt)
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to less secure method if crypto/rand fails
		return fallbackToken(length)
	}

	// Use URL-safe base64 encoding without padding
	token := base64.URLEncoding.EncodeToString(bytes)
	return token[:length]
}

// fallbackToken provides a fallback token generation method
func fallbackToken(length int) string {
	// Use timestamp and process ID as fallback
	timestamp := time.Now().UnixNano()
	pid := os.Getpid()

	// Simple hash-like combination
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte((timestamp + int64(pid*i)) % 256)
	}

	return base64.URLEncoding.EncodeToString(result)[:length]
}

// GetEnvironmentFromEnv gets environment from environment variables
func GetEnvironmentFromEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = "dev" // default to development
	}
	return env
}
