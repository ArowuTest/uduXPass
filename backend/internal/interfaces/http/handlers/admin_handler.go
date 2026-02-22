package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/internal/usecases/admin"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminAuthService *admin.AdminAuthService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminAuthService *admin.AdminAuthService) *AdminHandler {
	return &AdminHandler{
		adminAuthService: adminAuthService,
	}
}

// Login handles admin login requests
func (h *AdminHandler) Login(c *gin.Context) {
	var req admin.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Add IP address and user agent
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	response, err := h.adminAuthService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case entities.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			})
		case entities.ErrAccountLocked:
			c.JSON(http.StatusLocked, ErrorResponse{
				Error:   "Account locked",
				Message: "Account is temporarily locked due to multiple failed login attempts",
			})
		case entities.ErrAccountDeactivated:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Account deactivated",
				Message: "Your account has been deactivated",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Login failed",
				Message: "An error occurred during login",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// ChangePassword handles admin password change requests
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Admin authentication required",
		})
		return
	}
	
	var req admin.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	req.AdminID = adminID.(uuid.UUID)
	
	err := h.adminAuthService.ChangePassword(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case entities.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid current password",
				Message: "The current password is incorrect",
			})
		default:
			if validationErr, ok := err.(*entities.ValidationError); ok {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Validation failed",
					Message: validationErr.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Password change failed",
				Message: "An error occurred while changing password",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// CreateAdmin handles admin creation requests
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	// Check permission
	if !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You don't have permission to create admin users",
		})
		return
	}
	var req admin.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	response, err := h.adminAuthService.CreateAdmin(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case entities.ErrAdminUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Admin already exists",
				Message: "An admin user with this email already exists",
			})
		default:
			if validationErr, ok := err.(*entities.ValidationError); ok {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Validation failed",
					Message: validationErr.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Admin creation failed",
				Message: "An error occurred while creating admin user",
			})
		}
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// GetAdmin handles admin retrieval requests
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	adminIDParam := c.Param("id")
	adminID, err := uuid.Parse(adminIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid admin ID",
			Message: "The provided admin ID is not valid",
		})
		return
	}
	
	// Check if requesting own profile or has permission
	currentAdminID, _ := c.Get("admin_id")
	if adminID != currentAdminID.(uuid.UUID) && !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You can only view your own profile or need admin management permission",
		})
		return
	}
	
	admin, err := h.adminAuthService.GetAdmin(c.Request.Context(), adminID)
	if err != nil {
		switch err {
		case entities.ErrAdminUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Admin not found",
				Message: "The requested admin user was not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to get admin",
				Message: "An error occurred while retrieving admin user",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, admin)
}

// UpdateAdmin handles admin update requests
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	adminIDParam := c.Param("id")
	adminID, err := uuid.Parse(adminIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid admin ID",
			Message: "The provided admin ID is not valid",
		})
		return
	}
	
	// Check if updating own profile or has permission
	currentAdminID, _ := c.Get("admin_id")
	if adminID != currentAdminID.(uuid.UUID) && !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You can only update your own profile or need admin management permission",
		})
		return
	}
	
	var req admin.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	req.AdminID = adminID
	
	response, err := h.adminAuthService.UpdateAdmin(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case entities.ErrAdminUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Admin not found",
				Message: "The requested admin user was not found",
			})
		default:
			if validationErr, ok := err.(*entities.ValidationError); ok {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Validation failed",
					Message: validationErr.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Admin update failed",
				Message: "An error occurred while updating admin user",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// ListAdmins handles admin listing requests
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	// Check permission
	if !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You don't have permission to list admin users",
		})
		return
	}
	
	// Parse query parameters
	filter := repositories.AdminUserFilter{}
	
	if role := c.Query("role"); role != "" {
		adminRole := entities.AdminRole(role)
		filter.Role = &adminRole
	}
	
	if active := c.Query("active"); active != "" {
		if isActive, err := strconv.ParseBool(active); err == nil {
			filter.IsActive = &isActive
		}
	}
	
	if email := c.Query("email"); email != "" {
		filter.Email = &email
	}
	
	// Parse pagination
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}
	
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}
	
	// Parse sorting
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	
	if sortDir := c.Query("sort_direction"); sortDir != "" {
		filter.SortDirection = sortDir
	}
	
	// Convert filter to request
	req := &admin.ListAdminsRequest{
		Page:   1, // Default page
		Limit:  filter.Limit,
		Search: filter.Search,
	}
	if filter.Role != nil {
		req.Role = filter.Role
	}
	if filter.Status != nil {
		req.Status = filter.Status
	}
	
	response, err := h.adminAuthService.ListAdmins(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list admins",
			Message: "An error occurred while retrieving admin users",
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// DeleteAdmin handles admin deletion requests
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	// Check permission
	if !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You don't have permission to delete admin users",
		})
		return
	}
	
	adminIDParam := c.Param("id")
	adminID, err := uuid.Parse(adminIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid admin ID",
			Message: "The provided admin ID is not valid",
		})
		return
	}
	
	// Prevent self-deletion
	currentAdminID, _ := c.Get("admin_id")
	if adminID == currentAdminID.(uuid.UUID) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Cannot delete self",
			Message: "You cannot delete your own admin account",
		})
		return
	}
	
	err = h.adminAuthService.DeleteAdmin(c.Request.Context(), adminID)
	if err != nil {
		switch err {
		case entities.ErrAdminUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Admin not found",
				Message: "The requested admin user was not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Admin deletion failed",
				Message: "An error occurred while deleting admin user",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin user deleted successfully",
	})
}

// UnlockAdmin handles admin account unlock requests
func (h *AdminHandler) UnlockAdmin(c *gin.Context) {
	// Check permission
	if !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You don't have permission to unlock admin accounts",
		})
		return
	}
	
	adminIDParam := c.Param("id")
	adminID, err := uuid.Parse(adminIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid admin ID",
			Message: "The provided admin ID is not valid",
		})
		return
	}
	
	err = h.adminAuthService.UnlockAdmin(c.Request.Context(), adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Unlock failed",
			Message: "An error occurred while unlocking admin account",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin account unlocked successfully",
	})
}

// GetAdminStats handles admin statistics requests
func (h *AdminHandler) GetAdminStats(c *gin.Context) {
	// Check permission
	if !h.hasPermission(c, entities.PermissionAdminManagement) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Message: "You don't have permission to view admin statistics",
		})
		return
	}
	
	stats, err := h.adminAuthService.GetAdminStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get stats",
			Message: "An error occurred while retrieving admin statistics",
		})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// RefreshToken handles token refresh requests
func (h *AdminHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	refreshReq := &admin.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
	
	response, err := h.adminAuthService.RefreshToken(c.Request.Context(), refreshReq)
	if err != nil {
		switch err {
		case entities.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid token",
				Message: "The refresh token is invalid or expired",
			})
		case entities.ErrAccountDeactivated:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Account deactivated",
				Message: "Your account has been deactivated",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Token refresh failed",
				Message: "An error occurred while refreshing token",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// GetProfile handles admin profile requests
func (h *AdminHandler) GetProfile(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Admin authentication required",
		})
		return
	}
	
	response, err := h.adminAuthService.GetAdmin(c.Request.Context(), adminID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get profile",
			Message: "An error occurred while retrieving your profile",
		})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// hasPermission checks if the current admin has a specific permission
func (h *AdminHandler) hasPermission(c *gin.Context, permission entities.AdminPermission) bool {
	permissions, exists := c.Get("admin_permissions")
	if !exists {
		return false
	}
	
	adminPermissions, ok := permissions.([]entities.AdminPermission)
	if !ok {
		return false
	}
	
	for _, p := range adminPermissions {
		if p == permission {
			return true
		}
	}
	
	return false
}

