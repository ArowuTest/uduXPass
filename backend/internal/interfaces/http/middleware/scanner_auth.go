package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/pkg/security"
)

// ScannerAuthMiddleware validates scanner JWT tokens
func ScannerAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := security.ValidateJWT(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Check if it's an access token (not refresh token)
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token type",
			})
			c.Abort()
			return
		}

	// Extract scanner information from claims
	scannerIDStr, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token: missing user ID",
		})
		c.Abort()
		return
	}

		scannerID, err := uuid.Parse(scannerIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid scanner ID format",
			})
			c.Abort()
			return
		}

		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)
		permissions, _ := claims["permissions"].([]interface{})

		// Convert permissions to string slice
		var permissionStrings []string
		if permissions != nil {
			for _, perm := range permissions {
				if permStr, ok := perm.(string); ok {
					permissionStrings = append(permissionStrings, permStr)
				}
			}
		}

		// Set scanner information in context
		c.Set("scanner_id", scannerID)
		c.Set("scanner_username", username)
		c.Set("scanner_role", role)
		c.Set("scanner_permissions", permissionStrings)

		// Continue to next handler
		c.Next()
	}
}

// ScannerPermissionMiddleware checks if scanner has required permission
func ScannerPermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("scanner_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Scanner permissions not found",
			})
			c.Abort()
			return
		}

		permissionStrings, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Invalid permissions format",
			})
			c.Abort()
			return
		}

		// Check if scanner has the required permission
		hasPermission := false
		for _, perm := range permissionStrings {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions",
				"required_permission": requiredPermission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ScannerRoleMiddleware checks if scanner has required role
func ScannerRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("scanner_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Scanner role not found",
			})
			c.Abort()
			return
		}

		scannerRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Invalid role format",
			})
			c.Abort()
			return
		}

		// Check if scanner has one of the allowed roles
		hasRole := false
		for _, allowedRole := range allowedRoles {
			if scannerRole == allowedRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient role privileges",
				"required_roles": allowedRoles,
				"current_role":   scannerRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

