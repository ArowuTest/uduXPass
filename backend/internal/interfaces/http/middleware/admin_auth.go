package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/pkg/jwt"
)

// AdminAuthMiddleware handles admin authentication
type AdminAuthMiddleware struct {
	jwtService *jwt.JWTService
	adminRepo  repositories.AdminUserRepository
}

// NewAdminAuthMiddleware creates a new admin authentication middleware
func NewAdminAuthMiddleware(jwtService *jwt.JWTService, adminRepo repositories.AdminUserRepository) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		jwtService: jwtService,
		adminRepo:  adminRepo,
	}
}

// RequireAuth middleware that requires admin authentication
func (m *AdminAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Missing authorization header",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}
		
		// Check if token has Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}
		
		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Missing token",
				"message": "Bearer token is required",
			})
			c.Abort()
			return
		}
		
		// Validate token
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "The provided token is invalid or expired",
			})
			c.Abort()
			return
		}
		
		// Parse admin ID
		adminID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Token contains invalid admin ID",
			})
			c.Abort()
			return
		}
		
		// Get admin user
		admin, err := m.adminRepo.GetByID(c.Request.Context(), adminID)
		if err != nil {
			if err == entities.ErrAdminUserNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Admin not found",
					"message": "The admin user associated with this token was not found",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Authentication failed",
					"message": "An error occurred during authentication",
				})
			}
			c.Abort()
			return
		}
		
		// Check if admin is active
		if !admin.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Account deactivated",
				"message": "Your admin account has been deactivated",
			})
			c.Abort()
			return
		}
		
		// Check if admin is locked
		if admin.IsLocked() {
			c.JSON(http.StatusLocked, gin.H{
				"error":   "Account locked",
				"message": "Your admin account is temporarily locked",
			})
			c.Abort()
			return
		}
		
		// Set admin information in context
		c.Set("admin_id", admin.ID)
		c.Set("admin_email", admin.Email)
		c.Set("admin_role", admin.Role)
		fmt.Printf("Permissions type: %T\n", admin.Permissions)
c.Set("admin_permissions", admin.Permissions)
		c.Set("admin", admin)
		
		c.Next()
	}
}

// RequirePermission middleware that requires specific permission
func (m *AdminAuthMiddleware) RequirePermission(permission entities.AdminPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if admin is authenticated
		permissions, exists := c.Get("admin_permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Admin authentication required",
			})
			c.Abort()
			return
		}
		
		// Check permission
		adminPermissions, ok := permissions.(entities.AdminPermissionArray)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Permission check failed",
				"message": "Unable to verify permissions",
			})
			c.Abort()
			return
		}
		
		hasPermission := false
		for _, p := range adminPermissions {
			if p == permission {
				hasPermission = true
				break
			}
		}
		
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "You don't have permission to perform this action",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (m *AdminAuthMiddleware) RequireRole(role entities.AdminRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if admin is authenticated
		adminRole, exists := c.Get("admin_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Admin authentication required",
			})
			c.Abort()
			return
		}
		
		// Check role
		if adminRole.(entities.AdminRole) != role {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient role",
				"message": "You don't have the required role to perform this action",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireAnyRole middleware that requires any of the specified roles
func (m *AdminAuthMiddleware) RequireAnyRole(roles ...entities.AdminRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if admin is authenticated
		adminRole, exists := c.Get("admin_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Admin authentication required",
			})
			c.Abort()
			return
		}
		
		// Check if admin has any of the required roles
		hasRole := false
		currentRole := adminRole.(entities.AdminRole)
		for _, role := range roles {
			if currentRole == role {
				hasRole = true
				break
			}
		}
		
		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient role",
				"message": "You don't have any of the required roles to perform this action",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireAnyPermission middleware that requires any of the specified permissions
func (m *AdminAuthMiddleware) RequireAnyPermission(requiredPermissions ...entities.AdminPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if admin is authenticated
		permissions, exists := c.Get("admin_permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Admin authentication required",
			})
			c.Abort()
			return
		}
		
		// Check permissions
		adminPermissions, ok := permissions.(entities.AdminPermissionArray)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Permission check failed",
				"message": "Unable to verify permissions",
			})
			c.Abort()
			return
		}
		
		hasPermission := false
		for _, adminPerm := range adminPermissions {
			for _, requiredPerm := range requiredPermissions {
				if adminPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}
		
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "You don't have any of the required permissions to perform this action",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireSuperAdmin middleware that requires super admin role
func (m *AdminAuthMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return m.RequireRole(entities.AdminRoleSuperAdmin)
}

// OptionalAuth middleware that optionally authenticates admin (doesn't fail if no auth)
func (m *AdminAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		
		// Check if token has Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}
		
		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.Next()
			return
		}
		
		// Validate token
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}
		
		// Parse admin ID
		adminID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Next()
			return
		}
		
		// Get admin user
		admin, err := m.adminRepo.GetByID(c.Request.Context(), adminID)
		if err != nil {
			c.Next()
			return
		}
		
		// Check if admin is active and not locked
		if !admin.IsActive || admin.IsLocked() {
			c.Next()
			return
		}
		
		// Set admin information in context
		c.Set("admin_id", admin.ID)
		c.Set("admin_email", admin.Email)
		c.Set("admin_role", admin.Role)
		fmt.Printf("Permissions type: %T\n", admin.Permissions)
c.Set("admin_permissions", admin.Permissions)
		c.Set("admin", admin)
		
		c.Next()
	}
}

// ActivityLogger middleware that logs admin activities
func (m *AdminAuthMiddleware) ActivityLogger(action string, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continue with request
		c.Next()
		
		// Log activity after request completion
		adminID, exists := c.Get("admin_id")
		if !exists {
			return
		}
		
		// Get resource ID if available
		var resourceID *uuid.UUID
		if id := c.Param("id"); id != "" {
			if parsedID, err := uuid.Parse(id); err == nil {
				resourceID = &parsedID
			}
		}
		
		// Log the activity (this would typically be done asynchronously)
		go func() {
			// Implementation would depend on having an activity log service
			// For now, we'll skip the actual logging implementation
			_ = adminID    // TODO: Use in actual logging implementation
			_ = resourceID // TODO: Use in actual logging implementation
		}()
	}
}

// RateLimitByAdmin middleware that applies rate limiting per admin user
func (m *AdminAuthMiddleware) RateLimitByAdmin(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("admin_id")
		if !exists {
			c.Next()
			return
		}
		
		// Implementation would depend on having a rate limiting service
		// For now, we'll skip the actual rate limiting implementation
		_ = adminID // TODO: Use in actual rate limiting implementation
		c.Next()
	}
}

