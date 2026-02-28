package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/internal/usecases/admin"
	"github.com/uduxpass/backend/internal/usecases/events"
)

// AdminHandlerExtended handles all admin-related HTTP requests
type AdminHandlerExtended struct {
	adminAuthService *admin.AdminAuthService
	eventService     *events.EventService
	userRepo         repositories.UserRepository
	orderRepo        repositories.OrderRepository
	ticketRepo       repositories.TicketRepository
	scannerUserRepo  repositories.ScannerUserRepository
	organizerRepo    repositories.OrganizerRepository
}

// NewAdminHandlerExtended creates a new extended admin handler
func NewAdminHandlerExtended(
	adminAuthService *admin.AdminAuthService,
	eventService *events.EventService,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	ticketRepo repositories.TicketRepository,
	scannerUserRepo repositories.ScannerUserRepository,
	organizerRepo repositories.OrganizerRepository,
) *AdminHandlerExtended {
	return &AdminHandlerExtended{
		adminAuthService: adminAuthService,
		eventService:     eventService,
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		ticketRepo:       ticketRepo,
		scannerUserRepo:  scannerUserRepo,
		organizerRepo:    organizerRepo,
	}
}

// Login handles admin login requests
func (h *AdminHandlerExtended) Login(c *gin.Context) {
	var req admin.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Add IP address and user agent
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	response, err := h.adminAuthService.Login(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// Event Management Methods

func (h *AdminHandlerExtended) GetEvents(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	// Use GetPublicEvents for now - TODO: Add admin-specific list method
	req := &events.GetPublicEventsRequest{
		Page:  page,
		Limit: limit,
	}
	
	response, err := h.eventService.GetPublicEvents(ctx, req)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": response.Events,
			"pagination": response.Pagination,
		},
	})
}

func (h *AdminHandlerExtended) CreateEvent(c *gin.Context) {
	var req events.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	
	response, err := h.eventService.CreateEvent(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *AdminHandlerExtended) GetEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	req := &events.GetEventDetailsRequest{
		EventID: eventID,
	}
	
	response, err := h.eventService.GetEventDetails(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}
	
	event := response.Event
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    event,
	})
}

func (h *AdminHandlerExtended) UpdateEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	_, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	// TODO: Implement UpdateEvent method in event service
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Update event not implemented yet"})
}

func (h *AdminHandlerExtended) DeleteEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	_, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	// TODO: Implement DeleteEvent method in event service
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Delete event not implemented yet"})
}

func (h *AdminHandlerExtended) PublishEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	req := &events.PublishEventRequest{
		EventID: eventID,
	}
	
	if _, err := h.eventService.PublishEvent(c.Request.Context(), req); err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event published successfully",
	})
}

func (h *AdminHandlerExtended) GetEventAnalytics(c *gin.Context) {
	eventIDStr := c.Param("id")
	_, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	// TODO: Implement GetEventAnalytics method in event service
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Event analytics not implemented yet"})
}

// User Management Methods

func (h *AdminHandlerExtended) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	filter := repositories.UserFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  page,
			Limit: limit,
		},
	}
	
	users, pagination, err := h.userRepo.List(ctx, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"pagination": pagination,
		},
	})
}

func (h *AdminHandlerExtended) CreateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func (h *AdminHandlerExtended) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Order Management Methods

func (h *AdminHandlerExtended) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	filter := repositories.OrderFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  page,
			Limit: limit,
		},
	}
	
	orders, pagination, err := h.orderRepo.List(ctx, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"orders": orders,
			"pagination": pagination,
		},
	})
}

func (h *AdminHandlerExtended) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	
	order, err := h.orderRepo.GetByID(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

func (h *AdminHandlerExtended) UpdateOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) DeleteOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Ticket Management Methods

func (h *AdminHandlerExtended) GetTickets(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	filter := repositories.TicketFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  page,
			Limit: limit,
		},
	}
	
	tickets, pagination, err := h.ticketRepo.List(ctx, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"tickets": tickets,
			"pagination": pagination,
		},
	})
}

func (h *AdminHandlerExtended) GetTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	ticket, err := h.ticketRepo.GetByID(c.Request.Context(), ticketID)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticket,
	})
}

func (h *AdminHandlerExtended) UpdateTicket(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) ValidateTicket(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Analytics Methods

func (h *AdminHandlerExtended) GetDashboard(c *gin.Context) {
	// Get dashboard statistics
	// TODO: Implement proper dashboard aggregation
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"totalEvents":  0,
			"totalOrders":  0,
			"totalTickets": 0,
			"totalRevenue": 0,
		},
	})
}

func (h *AdminHandlerExtended) GetSalesAnalytics(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) GetUserAnalytics(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Scanner User Management Methods

func (h *AdminHandlerExtended) GetScannerUsers(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	filter := &repositories.ScannerUserFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  page,
			Limit: limit,
		},
	}
	
	scannerUsers, pagination, err := h.scannerUserRepo.List(ctx, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"scannerUsers": scannerUsers,
			"pagination": pagination,
		},
	})
}

func (h *AdminHandlerExtended) CreateScannerUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) GetScannerUser(c *gin.Context) {
	scannerUserIDStr := c.Param("id")
	scannerUserID, err := uuid.Parse(scannerUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scanner user ID"})
		return
	}
	
	scannerUser, err := h.scannerUserRepo.GetByID(c.Request.Context(), scannerUserID)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    scannerUser,
	})
}

func (h *AdminHandlerExtended) UpdateScannerUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) DeleteScannerUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) ResetScannerUserPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) UpdateScannerUserStatus(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Settings Methods

func (h *AdminHandlerExtended) GetSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) UpdateSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// Export Methods

func (h *AdminHandlerExtended) ExportUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) ExportOrders(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) ExportEvents(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *AdminHandlerExtended) ExportTickets(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetOrganizers returns a list of all organizers
func (h *AdminHandlerExtended) GetOrganizers(c *gin.Context) {
	filter := repositories.OrganizerFilter{
		Page:  1,
		Limit: 100,
	}
	organizers, _, err := h.organizerRepo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    organizers,
	})
}
