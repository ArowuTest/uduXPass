package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/infrastructure/database"
	"github.com/uduxpass/backend/internal/infrastructure/email"
	"github.com/uduxpass/backend/internal/infrastructure/payments"
	"github.com/uduxpass/backend/internal/interfaces/http/handlers"
	"github.com/uduxpass/backend/internal/usecases/admin"
	"github.com/uduxpass/backend/internal/usecases/auth"
	"github.com/uduxpass/backend/internal/usecases/events"
	"github.com/uduxpass/backend/internal/usecases/orders"
	paymentservice "github.com/uduxpass/backend/internal/usecases/payments"
	"github.com/uduxpass/backend/internal/usecases/scanner"
	"github.com/uduxpass/backend/pkg/jwt"
	"github.com/uduxpass/backend/pkg/security"
)

// Config holds server configuration
type Config struct {
	Host               string
	Port               string
	Environment        string
	JWTSecret          string
	CORSAllowedOrigins string
}

// Server represents the HTTP server
type Server struct {
	config     *Config
	router     *gin.Engine
	httpServer *http.Server
	dbManager  *database.DatabaseManager
	
	// Services
	jwtService      jwt.Service
	passwordService security.PasswordService
	authService     *auth.AuthService
	adminAuthService *admin.AdminAuthService
	eventService    *events.EventService
	orderService    *orders.OrderService
	paymentService  *paymentservice.PaymentService
	scannerAuthService *scanner.ScannerAuthService
	
	// Handlers
	authHandler    *handlers.AuthHandler
	adminHandler   *handlers.AdminHandlerExtended
	scannerHandler *handlers.ScannerHandler
	orderHandler   *handlers.OrderHandler
}

// NewServer creates a new HTTP server with proper dependency injection
func NewServer(config *Config, dbManager *database.DatabaseManager) *Server {
	// Initialize services
	jwtService := jwt.NewJWTService(
		config.JWTSecret,
		24*time.Hour,  // Access token TTL (extended for E2E testing)
		168*time.Hour, // Refresh token TTL (7 days)
		"uduxpass",
	)
	
	passwordService := security.NewBcryptPasswordService(security.BcryptConfig{
		Cost: 10, // Default bcrypt cost
	})
	
	// Initialize use case services
	authService := auth.NewAuthService(
		dbManager.Users(),
		dbManager.OTPTokens(),
		jwtService,
		passwordService,
	)
	
	adminAuthService := admin.NewAdminAuthService(
		dbManager.AdminUsers(),
		jwtService,
		passwordService,
	)
	
	eventService := events.NewEventService(
		dbManager.Events(),
		dbManager.Tours(),
		dbManager.Organizers(),
		dbManager.TicketTiers(),
		dbManager.UnitOfWork(),
	)
	
	orderService := orders.NewOrderService(
		dbManager.Orders(),
		dbManager.OrderLines(),
		dbManager.InventoryHolds(),
		dbManager.Events(),
		dbManager.TicketTiers(),
		dbManager.Users(),
	)
	
	// Initialize email service
	emailService := email.NewSMTPEmailService()
	
	// Initialize payment providers
	// Get Paystack secret key from environment
	paystackSecretKey := getEnv("PAYSTACK_SECRET_KEY", "sk_test_b748a89ad84f35c2c46cffc3581e1d7b8f6b4b3e")
	paystackProvider := payments.NewPaystackProvider(paystackSecretKey)
	
	// MoMo provider - placeholder for now
	var momoProvider payments.MoMoProvider
	
	paymentService := paymentservice.NewPaymentService(
		dbManager.Payments(),
		dbManager.Orders(),
		dbManager.Tickets(),
		dbManager.InventoryHolds(),
		momoProvider,
		*paystackProvider,
		dbManager.UnitOfWork(),
		emailService,
	)
	
	scannerAuthService := scanner.NewScannerAuthService(
		dbManager,
		config.JWTSecret,
	)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandlerExtended(
		adminAuthService,
		eventService,
		dbManager.Users(),
		dbManager.Orders(),
		dbManager.Tickets(),
		dbManager.ScannerUsers(),
	)
	scannerHandler := handlers.NewScannerHandler(
		scannerAuthService,
		dbManager,
	)
	
	// Set Gin mode based on environment
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	server := &Server{
		config:             config,
		router:             gin.New(),
		dbManager:          dbManager,
		jwtService:         jwtService,
		passwordService:    passwordService,
		authService:        authService,
		adminAuthService:   adminAuthService,
		eventService:       eventService,
		orderService:       orderService,
		paymentService:     paymentService,
		scannerAuthService: scannerAuthService,
		authHandler:        authHandler,
		adminHandler:       adminHandler,
		scannerHandler:     scannerHandler,
		orderHandler:       handlers.NewOrderHandler(orderService, paymentService),
	}
	
	server.setupMiddleware()
	server.setupRoutes()
	
	return server
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())
	
	// Logger middleware
	s.router.Use(gin.Logger())
	
	// CORS middleware - Custom implementation for production-ready configuration
	s.router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Determine if origin is allowed
		allowOrigin := ""
		if s.config.CORSAllowedOrigins != "" {
			if s.config.CORSAllowedOrigins == "*" {
				// Allow all origins (development only)
				allowOrigin = "*"
			} else {
				// Check if origin is in allowed list
				allowedOrigins := strings.Split(s.config.CORSAllowedOrigins, ",")
				for _, allowed := range allowedOrigins {
					if strings.TrimSpace(allowed) == origin {
						allowOrigin = origin
						break
					}
				}
			}
		} else {
			// Default behavior based on environment
			if s.config.Environment == "development" || s.config.Environment == "sandbox" {
				// Allow all origins in development/sandbox
				allowOrigin = "*"
			} else {
				// Check against default allowed origins
				defaultOrigins := []string{
					"http://localhost:3000",
					"http://localhost:5173",
					"http://localhost:8080",
				}
				for _, allowed := range defaultOrigins {
					if allowed == origin {
						allowOrigin = origin
						break
					}
				}
			}
		}
		
		// Set CORS headers if origin is allowed
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Max-Age", "43200") // 12 hours
			
			// Only set Allow-Credentials if not using wildcard
			if allowOrigin != "*" {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.handleHealth)
	
	// API v1 routes
	v1 := s.router.Group("/v1")
	{
		// Public authentication routes
		auth := v1.Group("/auth")
		{
			// Email authentication
			email := auth.Group("/email")
			{
				email.POST("/register", s.authHandler.RegisterEmailUser)
				email.POST("/login", s.authHandler.LoginEmailUser)
			}
			
			// MoMo authentication
			momo := auth.Group("/momo")
			{
				momo.POST("/initiate", s.authHandler.InitiateMoMoAuth)
				momo.POST("/verify", s.authHandler.VerifyMoMoOTP)
			}
			
			// Token refresh
			auth.POST("/refresh", s.authHandler.RefreshToken)
		}
		
		// Public events routes
		events := v1.Group("/events")
		{
			events.GET("", s.handleGetEvents)
			events.GET("/:id", s.handleGetEvent)
		}
		
		// Protected user routes
		user := v1.Group("/user")
		user.Use(s.authMiddleware())
		{
			user.GET("/profile", s.handleGetProfile)
			user.PUT("/profile", s.handleUpdateProfile)
			user.GET("/orders", s.handleGetUserOrders)
			user.GET("/tickets", s.handleGetUserTickets)
		}
		
		// Order routes
		orders := v1.Group("/orders")
		orders.Use(s.authMiddleware())
		{
			orders.POST("", s.handleCreateOrder)
			orders.GET("/:id", s.handleGetOrder)
			orders.POST("/:id/cancel", s.handleCancelOrder)
		}
		
		// Payment routes
		payments := v1.Group("/payments")
		payments.Use(s.authMiddleware())
		{
			payments.POST("/initiate", s.handleInitiatePayment)
			payments.GET("/:id/verify", s.handleVerifyPayment)
		}
		
		// Webhook routes (no auth required)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/momo", s.handleMomoWebhook)
			webhooks.POST("/paystack", s.handlePaystackWebhook)
		}
		
		// Scanner routes
		scanner := v1.Group("/scanner")
		{
			// Scanner authentication routes (no auth required)
			scannerAuth := scanner.Group("/auth")
			{
				scannerAuth.POST("/login", s.scannerHandler.Login)
				scannerAuth.POST("/refresh", s.scannerHandler.RefreshToken)
			}
			
			// Protected scanner routes
			scannerProtected := scanner.Group("")
			scannerProtected.Use(s.scannerAuthMiddleware())
			{
				scannerProtected.POST("/logout", s.scannerHandler.Logout)
				scannerProtected.GET("/profile", s.scannerHandler.GetProfile)
				scannerProtected.GET("/events", s.scannerHandler.GetAssignedEvents)
				scannerProtected.POST("/session/start", s.scannerHandler.StartSession)
				scannerProtected.POST("/session/end", s.scannerHandler.EndSession)
				scannerProtected.GET("/session/current", s.scannerHandler.GetCurrentSession)
				scannerProtected.POST("/validate", s.scannerHandler.ValidateTicket)
				scannerProtected.GET("/stats", s.scannerHandler.GetStats)
				scannerProtected.GET("/validation-history", s.scannerHandler.GetValidationHistory)
			}
		}
		
		// Admin routes
		admin := v1.Group("/admin")
		{
			// Admin authentication routes (no auth required)
			adminAuth := admin.Group("/auth")
			{
				adminAuth.POST("/login", s.adminHandler.Login)
			}
			
			// Protected admin routes
			adminProtected := admin.Group("")
			adminProtected.Use(s.adminAuthMiddleware())
			{
				// Event management
				adminProtected.GET("/events", s.adminHandler.GetEvents)
				adminProtected.POST("/events", s.adminHandler.CreateEvent)
				adminProtected.GET("/events/:id", s.adminHandler.GetEvent)
				adminProtected.PUT("/events/:id", s.adminHandler.UpdateEvent)
				adminProtected.DELETE("/events/:id", s.adminHandler.DeleteEvent)
				adminProtected.POST("/events/:id/publish", s.adminHandler.PublishEvent)
				adminProtected.GET("/events/:id/analytics", s.adminHandler.GetEventAnalytics)
				
				// User management
				adminProtected.GET("/users", s.adminHandler.GetUsers)
				adminProtected.POST("/users", s.adminHandler.CreateUser)
				adminProtected.GET("/users/:id", s.adminHandler.GetUser)
				adminProtected.PUT("/users/:id", s.adminHandler.UpdateUser)
				adminProtected.DELETE("/users/:id", s.adminHandler.DeleteUser)
				
				// Order management
				adminProtected.GET("/orders", s.adminHandler.GetOrders)
				adminProtected.GET("/orders/:id", s.adminHandler.GetOrder)
				adminProtected.PUT("/orders/:id", s.adminHandler.UpdateOrder)
				adminProtected.DELETE("/orders/:id", s.adminHandler.DeleteOrder)
				
				// Ticket management
				adminProtected.GET("/tickets", s.adminHandler.GetTickets)
				adminProtected.GET("/tickets/:id", s.adminHandler.GetTicket)
				adminProtected.PUT("/tickets/:id", s.adminHandler.UpdateTicket)
				adminProtected.POST("/tickets/:id/validate", s.adminHandler.ValidateTicket)
				
				// Analytics and reports
				adminProtected.GET("/analytics/dashboard", s.adminHandler.GetDashboard)
				adminProtected.GET("/analytics/events", s.adminHandler.GetEventAnalytics)
				adminProtected.GET("/analytics/sales", s.adminHandler.GetSalesAnalytics)
				adminProtected.GET("/analytics/users", s.adminHandler.GetUserAnalytics)
				
				// Scanner user management
				adminProtected.GET("/scanner-users", s.adminHandler.GetScannerUsers)
				adminProtected.POST("/scanner-users", s.adminHandler.CreateScannerUser)
				adminProtected.GET("/scanner-users/:id", s.adminHandler.GetScannerUser)
				adminProtected.PUT("/scanner-users/:id", s.adminHandler.UpdateScannerUser)
				adminProtected.DELETE("/scanner-users/:id", s.adminHandler.DeleteScannerUser)
				adminProtected.POST("/scanner-users/:id/reset-password", s.adminHandler.ResetScannerUserPassword)
				adminProtected.PUT("/scanner-users/:id/status", s.adminHandler.UpdateScannerUserStatus)
				
				// Settings
				adminProtected.GET("/settings", s.adminHandler.GetSettings)
				adminProtected.PUT("/settings", s.adminHandler.UpdateSettings)
				
				// CSV Export routes
				adminProtected.GET("/export/users", s.adminHandler.ExportUsers)
				adminProtected.GET("/export/orders", s.adminHandler.ExportOrders)
				adminProtected.GET("/export/events", s.adminHandler.ExportEvents)
				adminProtected.GET("/export/tickets", s.adminHandler.ExportTickets)
			}
		}
	}
	
	// Admin routes without /v1 prefix for compatibility
	adminCompat := s.router.Group("/admin")
	adminCompat.Use(s.adminAuthMiddleware())
	{
		adminCompat.GET("/analytics/dashboard", s.adminHandler.GetDashboard)
		adminCompat.GET("/events", s.adminHandler.GetEvents)
		adminCompat.GET("/users", s.adminHandler.GetUsers)
		adminCompat.GET("/orders", s.adminHandler.GetOrders)
	}
}

// Middleware functions

// authMiddleware validates JWT tokens for regular users
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		
		token := parts[1]
		
		// Validate token
		claims, err := s.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
		// Set user context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		
		c.Next()
	}
}

// adminMiddleware checks if user has admin role
func (s *Server) adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		c.Abort()
		return
	}

	token := parts[1]

	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	if claims.Role != "super_admin" && claims.Role != "admin" && claims.Role != "event_manager" && claims.Role != "support" && claims.Role != "analyst" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		c.Abort()
		return
	}

	c.Set("adminID", claims.UserID)
	c.Set("adminRole", claims.Role)

	c.Next()
	}
}

// scannerAuthMiddleware validates JWT tokens for scanner users
func (s *Server) scannerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		
		token := parts[1]
		
		// Validate token
		claims, err := s.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
	// Check if role is scanner-related
	if claims.Role != "scanner" && claims.Role != "scanner_operator" && claims.Role != "scanner_supervisor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Scanner access required"})
		c.Abort()
		return
	}
		
		// Parse scanner ID to UUID
		scannerID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid scanner ID"})
			c.Abort()
			return
		}
		
		// Set scanner context
		c.Set("scanner_id", scannerID)
		c.Set("scanner_role", claims.Role)
		
		c.Next()
	}
}

// Handler functions (temporary implementations until proper handlers are connected)

func (s *Server) handleHealth(c *gin.Context) {
	// Check database health
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	
	dbHealthy := true
	if err := s.dbManager.Health(ctx); err != nil {
		dbHealthy = false
	}
	
	status := "healthy"
	httpStatus := http.StatusOK
	if !dbHealthy {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}
	
	c.JSON(httpStatus, gin.H{
		"status": status,
		"database": dbHealthy,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleGetEvents(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Parse query parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	search := c.Query("search")
	city := c.Query("city")
	
	// Use event service to get public events
	req := &events.GetPublicEventsRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
		City:   city,
	}
	
	response, err := s.eventService.GetPublicEvents(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
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

func (s *Server) handleGetEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
	ctx := c.Request.Context()
	event, err := s.dbManager.Events().GetByID(ctx, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	
	// Fetch active ticket tiers for this event
	tiers, err := s.dbManager.TicketTiers().GetActiveByEvent(ctx, eventID)
	if err != nil {
		// Log error but don't fail the request if tiers can't be fetched
		fmt.Printf("Warning: Failed to fetch ticket tiers for event %s: %v\n", eventID, err)
		tiers = []*entities.TicketTier{}
	}
	
	// Create response with event and ticket tiers
	response := gin.H{
		"id":               event.ID,
		"name":             event.Name,
		"slug":             event.Slug,
		"description":      event.Description,
		"event_date":       event.EventDate,
		"doors_open":       event.DoorsOpen,
		"venue_name":       event.VenueName,
		"venue_address":    event.VenueAddress,
		"venue_city":       event.VenueCity,
		"venue_state":      event.VenueState,
		"venue_country":    event.VenueCountry,
		"venue_capacity":   event.VenueCapacity,
		"event_image_url":  event.EventImageURL,
		"status":           event.Status,
		"sale_start":       event.SaleStart,
		"sale_end":         event.SaleEnd,
		"currency":         event.Currency,
		"category_id":      event.CategoryID,
		"settings":         event.Settings,
		"is_active":        event.IsActive,
		"created_at":       event.CreatedAt,
		"updated_at":       event.UpdatedAt,
		"ticket_tiers":     tiers,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *Server) handleGetProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	ctx := c.Request.Context()
	user, err := s.dbManager.Users().GetByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func (s *Server) handleUpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handleGetUserOrders(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handleGetUserTickets(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handleCreateOrder(c *gin.Context) {
	s.orderHandler.CreateOrder(c)
}

func (s *Server) handleGetOrder(c *gin.Context) {
	s.orderHandler.GetOrder(c)
}

func (s *Server) handleCancelOrder(c *gin.Context) {
	// TODO: Implement order cancellation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Order cancellation not yet implemented"})
}

func (s *Server) handleInitiatePayment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handleVerifyPayment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handleMomoWebhook(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) handlePaystackWebhook(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
