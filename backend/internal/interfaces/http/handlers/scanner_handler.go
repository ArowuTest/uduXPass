package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/internal/usecases/scanner"
)

// ScannerHandler handles scanner-related HTTP requests
type ScannerHandler struct {
	scannerService *scanner.ScannerAuthService
	repoManager    repositories.RepositoryManager
}

// NewScannerHandler creates a new scanner handler
func NewScannerHandler(scannerService *scanner.ScannerAuthService, repoManager repositories.RepositoryManager) *ScannerHandler {
	return &ScannerHandler{
		scannerService: scannerService,
		repoManager:    repoManager,
	}
}

// Login handles scanner authentication
func (h *ScannerHandler) Login(c *gin.Context) {
	var req entities.ScannerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Get client IP and user agent for audit logging
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Attempt authentication
	response, err := h.scannerService.Login(c.Request.Context(), req.Username, req.Password, clientIP, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication failed",
			"error":   err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *ScannerHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token is required",
		})
		return
	}

	response, err := h.scannerService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Token refresh failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles scanner logout
func (h *ScannerHandler) Logout(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	err := h.scannerService.Logout(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Logout failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetProfile returns the scanner's profile information
func (h *ScannerHandler) GetProfile(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	scanner, err := h.repoManager.ScannerUsers().GetByID(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get scanner profile",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    scanner,
	})
}

// GetAssignedEvents returns events assigned to the scanner
func (h *ScannerHandler) GetAssignedEvents(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	events, err := h.repoManager.ScannerUsers().GetAssignedEvents(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get assigned events",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": events,
		},
	})
}

// StartSession starts a new scanning session
func (h *ScannerHandler) StartSession(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	var req entities.ScannerSessionStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	session, err := h.scannerService.StartSession(c.Request.Context(), scannerID.(uuid.UUID), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start session",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    session,
		"message": "Session started successfully",
	})
}

// EndSession ends the current scanning session
func (h *ScannerHandler) EndSession(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	// Get active session
	session, err := h.repoManager.ScannerUsers().GetActiveSession(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No active session found",
		})
		return
	}

	err = h.scannerService.EndSession(c.Request.Context(), session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to end session",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Session ended successfully",
	})
}

// ValidateTicket validates a ticket and records the scan
func (h *ScannerHandler) ValidateTicket(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	var req entities.TicketValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Get active session
	session, err := h.repoManager.ScannerUsers().GetActiveSession(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No active scanning session. Please start a session first.",
		})
		return
	}

	// Parse event ID
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid event ID format",
		})
		return
	}

	response, err := h.scannerService.ValidateTicket(c.Request.Context(), scannerID.(uuid.UUID), session.ID, eventID, req.TicketCode, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ticket validation failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStats returns scanner statistics
func (h *ScannerHandler) GetStats(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	// Optional event ID filter
	var eventID *uuid.UUID
	if eventIDStr := c.Query("event_id"); eventIDStr != "" {
		if parsed, err := uuid.Parse(eventIDStr); err == nil {
			eventID = &parsed
		}
	}

	stats, err := h.repoManager.ScannerUsers().GetScannerStats(c.Request.Context(), scannerID.(uuid.UUID), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get scanner statistics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetValidationHistory returns ticket validation history
func (h *ScannerHandler) GetValidationHistory(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	scannerUUID := scannerID.(uuid.UUID)
	filter := &repositories.TicketValidationFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  page,
			Limit: limit,
		},
		ScannerID: &scannerUUID,
	}

	// Optional filters
	if sessionIDStr := c.Query("session_id"); sessionIDStr != "" {
		if sessionID, err := uuid.Parse(sessionIDStr); err == nil {
			filter.SessionID = &sessionID
		}
	}

	if result := c.Query("validation_result"); result != "" {
		filter.ValidationResult = result
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}

	validations, pagination, err := h.repoManager.ScannerUsers().GetValidationHistory(c.Request.Context(), scannerID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get validation history",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       validations,
		"pagination": pagination,
	})
}

// GetCurrentSession returns the current active session
func (h *ScannerHandler) GetCurrentSession(c *gin.Context) {
	scannerID, exists := c.Get("scanner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Scanner not authenticated",
		})
		return
	}

	session, err := h.repoManager.ScannerUsers().GetActiveSession(c.Request.Context(), scannerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No active session found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    session,
	})
}

