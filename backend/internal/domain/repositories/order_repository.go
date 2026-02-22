package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	// Create creates a new order
	Create(ctx context.Context, order *entities.Order) error
	
	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	
	// GetByCode retrieves an order by code
	GetByCode(ctx context.Context, code string) (*entities.Order, error)
	
	// GetBySecret retrieves an order by secret (for guest access)
	GetBySecret(ctx context.Context, secret string) (*entities.Order, error)
	
	// Update updates an existing order
	Update(ctx context.Context, order *entities.Order) error
	
	// Delete soft deletes an order
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves orders with pagination and filtering
	List(ctx context.Context, filter OrderFilter) ([]*entities.Order, *PaginationResult, error)
	
	// GetByUser retrieves orders for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID, filter OrderFilter) ([]*entities.Order, *PaginationResult, error)
	
	// GetByUserID retrieves orders for a specific user (alias for GetByUser)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error)
	
	// GetByEvent retrieves orders for a specific event
	GetByEvent(ctx context.Context, eventID uuid.UUID, filter OrderFilter) ([]*entities.Order, *PaginationResult, error)
	
	// GetByEventID retrieves orders for a specific event (alias for GetByEvent)
	GetByEventID(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*entities.Order, error)
	
	// GetByEmail retrieves orders for a specific email
	GetByEmail(ctx context.Context, email string, filter OrderFilter) ([]*entities.Order, *PaginationResult, error)
	
	// GetExpired retrieves expired orders
	GetExpired(ctx context.Context, filter OrderFilter) ([]*entities.Order, *PaginationResult, error)
	
	// GetExpiredOrders retrieves expired orders (alias for GetExpired)
	GetExpiredOrders(ctx context.Context) ([]*entities.Order, error)
	
	// UpdateStatus updates the order status
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error
	
	// MarkExpired marks orders as expired based on expiration time
	MarkExpired(ctx context.Context) (int, error)
	
	// GetOrderStats retrieves statistics for an order
	GetOrderStats(ctx context.Context, orderID uuid.UUID) (*OrderStats, error)
	
	// Exists checks if an order exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsByCode checks if an order exists by code
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

// OrderFilter defines filtering options for order queries
type OrderFilter struct {
	BaseFilter
	
	// Filtering
	UserID        *uuid.UUID
	EventID       *uuid.UUID
	Status        *entities.OrderStatus
	PaymentMethod *entities.PaymentMethod
	Email         string
	Phone         string
	Search        string // Search in code, email, customer name
	
	// Date filtering
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	ExpiresFrom   *time.Time
	ExpiresTo     *time.Time
	
	// Amount filtering
	MinAmount     *float64
	MaxAmount     *float64
	
	// Include related data
	IncludeLines    bool
	IncludeEvent    bool
	IncludeUser     bool
	IncludePayments bool
	IncludeTickets  bool
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders       int     `json:"total_orders" db:"total_orders"`
	PaidOrders        int     `json:"paid_orders" db:"paid_orders"`
	PendingOrders     int     `json:"pending_orders" db:"pending_orders"`
	CancelledOrders   int     `json:"cancelled_orders" db:"cancelled_orders"`
	ExpiredOrders     int     `json:"expired_orders" db:"expired_orders"`
	TotalRevenue      float64 `json:"total_revenue" db:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value" db:"average_order_value"`
	ConversionRate    float64 `json:"conversion_rate" db:"conversion_rate"`
}

