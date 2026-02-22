package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/pkg/jwt"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtService jwt.Service
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireAdminAuth middleware that requires admin authentication
func (m *AuthMiddleware) RequireAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		if !m.isAdminRole(claims.Role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Set admin context
		c.Set("admin_id", claims.UserID)
		c.Set("admin_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequirePermission middleware that requires specific permission
func (m *AuthMiddleware) RequirePermission(permission entities.AdminPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		if !m.isAdminRole(claims.Role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Check if admin has required permission
		if !m.hasPermission(claims.Role, permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Set admin context
		c.Set("admin_id", claims.UserID)
		c.Set("admin_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth middleware that optionally authenticates users
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			// Don't abort, just continue without authentication
			c.Next()
			return
		}

		// Set user context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Check if user has required role
		if claims.Role != role {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient role permissions",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// extractToken extracts JWT token from Authorization header
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for Bearer token
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// isAdminRole checks if the role is an admin role
func (m *AuthMiddleware) isAdminRole(role string) bool {
	adminRoles := []string{
		string(entities.AdminRoleSuperAdmin),
		string(entities.AdminRoleAdmin),
		string(entities.AdminRoleEventManager),
		string(entities.AdminRoleScannerOperator),
		string(entities.AdminRoleAnalyst),
	}

	for _, adminRole := range adminRoles {
		if role == adminRole {
			return true
		}
	}

	return false
}

// hasPermission checks if the admin role has the required permission
func (m *AuthMiddleware) hasPermission(role string, permission entities.AdminPermission) bool {
	// Super admin has all permissions
	if role == string(entities.AdminRoleSuperAdmin) {
		return true
	}

	// Define role-based permissions
	rolePermissions := map[string][]entities.AdminPermission{
		string(entities.AdminRoleAdmin): {
			entities.AdminPermissionEventsCreate,
			entities.AdminPermissionEventsUpdate,
			entities.AdminPermissionEventsDelete,
			entities.AdminPermissionEventsView,
			entities.AdminPermissionUsersView,
			entities.AdminPermissionUsersUpdate,
			entities.AdminPermissionOrdersView,
			entities.AdminPermissionOrdersUpdate,
			entities.AdminPermissionAnalyticsView,
			entities.AdminPermissionScannersView,
			entities.AdminPermissionScannersUpdate,
		},
		string(entities.AdminRoleEventManager): {
			entities.AdminPermissionEventsCreate,
			entities.AdminPermissionEventsUpdate,
			entities.AdminPermissionEventsView,
			entities.AdminPermissionOrdersView,
			entities.AdminPermissionAnalyticsView,
		},
		string(entities.AdminRoleScannerOperator): {
			entities.AdminPermissionEventsView,
			entities.AdminPermissionScannersView,
			entities.AdminPermissionTicketsValidate,
		},
		string(entities.AdminRoleAnalyst): {
			entities.AdminPermissionEventsView,
			entities.AdminPermissionUsersView,
			entities.AdminPermissionOrdersView,
			entities.AdminPermissionAnalyticsView,
		},
	}

	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}

	return false
}

// GetUserIDString gets the authenticated user ID from context as string
func GetUserIDString(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetAdminID gets the authenticated admin ID from context
func GetAdminID(c *gin.Context) (string, bool) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		return "", false
	}

	adminIDStr, ok := adminID.(string)
	return adminIDStr, ok
}

// GetUserRole gets the authenticated user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	roleStr, ok := role.(string)
	return roleStr, ok
}

// GetClaims gets the JWT claims from context
func GetClaims(c *gin.Context) (*jwt.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claimsObj, ok := claims.(*jwt.Claims)
	return claimsObj, ok
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin checks if the authenticated user is an admin
func IsAdmin(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	if !exists {
		return false
	}

	middleware := &AuthMiddleware{}
	return middleware.isAdminRole(role)
}

// HasPermission checks if the authenticated admin has a specific permission
func HasPermission(c *gin.Context, permission entities.AdminPermission) bool {
	role, exists := GetUserRole(c)
	if !exists {
		return false
	}

	middleware := &AuthMiddleware{}
	return middleware.hasPermission(role, permission)
}

