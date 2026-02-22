package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uduxpass/backend/internal/usecases/auth"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterEmailUser handles email user registration
func (h *AuthHandler) RegisterEmailUser(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate request
	if err := validateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}
	
	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, response)
}

// LoginEmailUser handles email user login
func (h *AuthHandler) LoginEmailUser(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate request
	if err := validateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}
	
	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

	// InitiateMoMoAuth handles MoMo authentication initiation
	func (h *AuthHandler) InitiateMoMoAuth(c *gin.Context) {
	type InitiateMoMoAuthRequest struct {
		Phone string `json:"phone" validate:"required"`
	}
	
	var req InitiateMoMoAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate request
	if err := validateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}
	
	// For now, use email-based OTP as a stub
	// In production, this would integrate with MoMo API
	err := h.authService.SendOTP(c.Request.Context(), req.Phone+"@momo.temp")
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to phone number",
		"phone":   req.Phone,
	})
}

// VerifyMoMoOTP handles MoMo OTP verification
func (h *AuthHandler) VerifyMoMoOTP(c *gin.Context) {
	type VerifyMoMoOTPRequest struct {
		Phone string `json:"phone" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}
	
	var req VerifyMoMoOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate request
	if err := validateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}
	
	// For now, use email-based OTP verification as a stub
	// In production, this would integrate with MoMo API
	err := h.authService.VerifyOTP(c.Request.Context(), req.Phone+"@momo.temp", req.OTP)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"phone":   req.Phone,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate request
	if err := validateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}
	
	// TODO: Implement token refresh logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Token refresh not implemented yet",
	})
}

