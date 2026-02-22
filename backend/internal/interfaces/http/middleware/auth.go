package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/pkg/jwt"
)

// JWTAuth creates a JWT authentication middleware
func JWTAuth(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}
		
		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}
		
		token := tokenParts[1]
		
		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		
		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			c.Abort()
			return
		}
		
		// Set user information in context
		c.Set("user_id", userID)
		c.Set("user_identifier", claims.Identifier)
		c.Set("jwt_claims", claims)
		
		c.Next()
	}
}

// OptionalJWTAuth creates an optional JWT authentication middleware
func OptionalJWTAuth(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		
		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}
		
		token := tokenParts[1]
		
		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}
		
		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.Next()
			return
		}
		
		// Set user information in context
		c.Set("user_id", userID)
		c.Set("user_identifier", claims.Identifier)
		c.Set("jwt_claims", claims)
		
		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserIdentifier extracts user identifier from context
func GetUserIdentifier(c *gin.Context) (string, bool) {
	identifier, exists := c.Get("user_identifier")
	if !exists {
		return "", false
	}
	
	id, ok := identifier.(string)
	return id, ok
}

// RequireUserID ensures a user ID is present in the context
func RequireUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		c.Abort()
		return uuid.Nil, false
	}
	
	return userID, true
}

