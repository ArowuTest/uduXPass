package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/usecases/orders"
	"github.com/uduxpass/backend/internal/usecases/payments"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService   *orders.OrderService
	paymentService *payments.PaymentService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(
	orderService *orders.OrderService,
	paymentService *payments.PaymentService,
) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		paymentService: paymentService,
	}
}

// CreateOrder handles order creation with payment initialization
// POST /v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req orders.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Set user ID from authenticated user
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}
	
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}
	
	req.UserID = userUUID

	// Create order with inventory holds
	orderResp, err := h.orderService.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Initialize payment with Paystack
	paymentReq := &payments.InitiatePaymentRequest{
		OrderID:       orderResp.Order.ID,
		PaymentMethod: "paystack", // Default to Paystack
		CustomerInfo: payments.PaymentCustomerInfo{
			Email:     orderResp.Order.CustomerEmail,
			Phone:     orderResp.Order.CustomerPhone,
			FirstName: orderResp.Order.CustomerFirstName,
			LastName:  orderResp.Order.CustomerLastName,
		},
		CallbackURL: c.Request.Host + "/v1/webhooks/paystack",
	}

	paymentResp, err := h.paymentService.InitiatePayment(c.Request.Context(), paymentReq)
	if err != nil {
		// Order created but payment initialization failed
		// The order will expire after hold duration if not paid
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Order created but payment initialization failed",
			"data": gin.H{
				"order":        orderResp.Order,
				"order_lines":  orderResp.OrderLines,
				"total_amount": orderResp.TotalAmount,
				"expires_at":   orderResp.ExpiresAt,
				"payment_error": err.Error(),
			},
		})
		return
	}

	// Success: Order created and payment initialized
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order created successfully",
		"data": gin.H{
			"order":         orderResp.Order,
			"order_lines":   orderResp.OrderLines,
			"total_amount":  orderResp.TotalAmount,
			"expires_at":    orderResp.ExpiresAt,
			"payment": gin.H{
				"payment_id":        paymentResp.PaymentID,
				"authorization_url": paymentResp.AuthorizationURL,
				"reference":         paymentResp.PaymentReference,
				"status":            paymentResp.Status,
				"expires_at":        paymentResp.ExpiresAt,
			},
		},
	})
}

// GetOrder retrieves an order by ID
// GET /v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	order, orderLines, err := h.orderService.GetOrderWithLines(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	// Verify order belongs to user (unless admin)
	role, _ := c.Get("role")
	if role != "admin" && role != "super_admin" {
		if order.UserID.String() != userID.(uuid.UUID).String() {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"order":       order,
			"order_lines": orderLines,
		},
	})
}

// GetUserOrders retrieves orders for the authenticated user
// GET /v1/orders
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseIntParam(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := parseIntParam(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID.(uuid.UUID), limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"orders": orders,
			"pagination": gin.H{
				"limit":  limit,
				"offset": offset,
			},
		},
	})
}

// Helper function to parse int parameters
func parseIntParam(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
