package handlers

import (
	"net/http"

	"smart-outgoing-demo/pkg/security"

	"github.com/gin-gonic/gin"
)

// RequireAuth requires authentication for protected routes using improved token system
func RequireAuth(authConfig *security.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if token is expired
		if authConfig.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":          "admin token expired",
				"token_required": true,
			})
			c.Abort()
			return
		}

		// Get token from various sources
		token := extractToken(c)

		// Validate token
		if err := authConfig.ValidateToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":          "invalid admin token",
				"token_required": true,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAuthWithJWT requires authentication using improved JWT-like token system
func RequireAuthWithJWT(jwtConfig *security.ImprovedTokenConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":          "token required",
				"token_required": true,
			})
			c.Abort()
			return
		}

		// Validate improved token
		claims, err := jwtConfig.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":          "invalid or expired token",
				"token_required": true,
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// extractToken extracts token from request headers or query parameters
func extractToken(c *gin.Context) string {
	// Try x-admin-token header first
	token := c.GetHeader("x-admin-token")
	if token != "" {
		return token
	}

	// Try x-user-id header (for improved token system)
	token = c.GetHeader("x-user-token")
	if token != "" {
		return token
	}

	// Try query parameter
	token = c.Query("token")
	if token != "" {
		return token
	}

	// Try Authorization header (Bearer token)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Extract Bearer token
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
		return authHeader
	}

	return ""
}
