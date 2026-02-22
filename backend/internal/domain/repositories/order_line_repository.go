package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// OrderLineFilter represents filtering options for order lines
type OrderLineFilter struct {
	BaseFilter
	OrderID       *uuid.UUID `json:"order_id,omitempty"`
	TicketTierID  *uuid.UUID `json:"ticket_tier_id,omitempty"`
	EventID       *uuid.UUID `json:"event_id,omitempty"`
	MinQuantity   *int       `json:"min_quantity,omitempty"`
	MaxQuantity   *int       `json:"max_quantity,omitempty"`
	MinPrice      *float64   `json:"min_price,omitempty"`
	MaxPrice      *float64   `json:"max_price,omitempty"`
	CreatedFrom   *time.Time `json:"created_from,omitempty"`
	CreatedTo     *time.Time `json:"created_to,omitempty"`
}

// OrderLineStats represents statistics for order lines
type OrderLineStats struct {
	TotalLines         int     `json:"total_lines" db:"total_lines"`
	TotalQuantity      int     `json:"total_quantity" db:"total_quantity"`
	TotalAmount        float64 `json:"total_amount" db:"total_amount"`
	AvgQuantityPerLine float64 `json:"avg_quantity_per_line" db:"avg_quantity_per_line"`
	AvgPricePerLine    float64 `json:"avg_price_per_line" db:"avg_price_per_line"`
	MinQuantity        int     `json:"min_quantity" db:"min_quantity"`
	MaxQuantity        int     `json:"max_quantity" db:"max_quantity"`
	MinPrice           float64 `json:"min_price" db:"min_price"`
	MaxPrice           float64 `json:"max_price" db:"max_price"`
	UniqueTicketTiers  int     `json:"unique_ticket_tiers" db:"unique_ticket_tiers"`
	UniqueOrders       int     `json:"unique_orders" db:"unique_orders"`
	FirstLineCreatedAt *time.Time `json:"first_line_created_at" db:"first_line_created_at"`
	LastLineCreatedAt  *time.Time `json:"last_line_created_at" db:"last_line_created_at"`
}

// OrderLineRepository defines the interface for order line data access
type OrderLineRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, orderLine *entities.OrderLine) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderLine, error)
	Update(ctx context.Context, orderLine *entities.OrderLine) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter OrderLineFilter) ([]*entities.OrderLine, *PaginationResult, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// Order-specific operations
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderLine, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderLine, error)
	DeleteByOrder(ctx context.Context, orderID uuid.UUID) error
	GetOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetOrderQuantity(ctx context.Context, orderID uuid.UUID) (int, error)

	// Ticket tier operations
	GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter OrderLineFilter) ([]*entities.OrderLine, *PaginationResult, error)
	GetTicketTierSales(ctx context.Context, ticketTierID uuid.UUID) (*OrderLineStats, error)

	// Event operations
	GetByEvent(ctx context.Context, eventID uuid.UUID, filter OrderLineFilter) ([]*entities.OrderLine, *PaginationResult, error)
	GetEventSales(ctx context.Context, eventID uuid.UUID) (*OrderLineStats, error)

	// Bulk operations
	CreateBatch(ctx context.Context, orderLines []*entities.OrderLine) error
	UpdateBatch(ctx context.Context, orderLines []*entities.OrderLine) error

	// Statistics and analytics
	GetStats(ctx context.Context, filter OrderLineFilter) (*OrderLineStats, error)
}

