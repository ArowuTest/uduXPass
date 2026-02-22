package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// InventoryHoldFilter represents filtering options for inventory holds
type InventoryHoldFilter struct {
	BaseFilter
	OrderID       *uuid.UUID `json:"order_id,omitempty"`
	TicketTierID  *uuid.UUID `json:"ticket_tier_id,omitempty"`
	EventID       *uuid.UUID `json:"event_id,omitempty"`
	MinQuantity   *int       `json:"min_quantity,omitempty"`
	MaxQuantity   *int       `json:"max_quantity,omitempty"`
	ExpiresFrom   *time.Time `json:"expires_from,omitempty"`
	ExpiresTo     *time.Time `json:"expires_to,omitempty"`
	CreatedFrom   *time.Time `json:"created_from,omitempty"`
	CreatedTo     *time.Time `json:"created_to,omitempty"`
	ExpiredOnly   bool       `json:"expired_only,omitempty"`
	ActiveOnly    bool       `json:"active_only,omitempty"`
}

// InventoryHoldStats represents statistics for inventory holds
type InventoryHoldStats struct {
	TotalHolds            int       `json:"total_holds" db:"total_holds"`
	ActiveHolds           int       `json:"active_holds" db:"active_holds"`
	ExpiredHolds          int       `json:"expired_holds" db:"expired_holds"`
	TotalQuantityHeld     int       `json:"total_quantity_held" db:"total_quantity_held"`
	ActiveQuantityHeld    int       `json:"active_quantity_held" db:"active_quantity_held"`
	ExpiredQuantityHeld   int       `json:"expired_quantity_held" db:"expired_quantity_held"`
	AvgQuantityPerHold    float64   `json:"avg_quantity_per_hold" db:"avg_quantity_per_hold"`
	MinQuantityPerHold    int       `json:"min_quantity_per_hold" db:"min_quantity_per_hold"`
	MaxQuantityPerHold    int       `json:"max_quantity_per_hold" db:"max_quantity_per_hold"`
	UniqueTicketTiers     int       `json:"unique_ticket_tiers" db:"unique_ticket_tiers"`
	UniqueOrders          int       `json:"unique_orders" db:"unique_orders"`
	FirstHoldCreatedAt    *time.Time `json:"first_hold_created_at" db:"first_hold_created_at"`
	LastHoldCreatedAt     *time.Time `json:"last_hold_created_at" db:"last_hold_created_at"`
	NextExpiryAt          *time.Time `json:"next_expiry_at" db:"next_expiry_at"`
	LatestExpiryAt        *time.Time `json:"latest_expiry_at" db:"latest_expiry_at"`
}

// InventoryHoldRepository defines the interface for inventory hold data access
type InventoryHoldRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, hold *entities.InventoryHold) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.InventoryHold, error)
	Update(ctx context.Context, hold *entities.InventoryHold) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter InventoryHoldFilter) ([]*entities.InventoryHold, *PaginationResult, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// Order-specific operations
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.InventoryHold, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.InventoryHold, error)
	DeleteByOrder(ctx context.Context, orderID uuid.UUID) error

	// Ticket tier operations
	GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter InventoryHoldFilter) ([]*entities.InventoryHold, *PaginationResult, error)
	GetTotalHeldQuantity(ctx context.Context, ticketTierID uuid.UUID) (int, error)

	// Expiry management
	GetExpired(ctx context.Context, filter InventoryHoldFilter) ([]*entities.InventoryHold, *PaginationResult, error)
	GetActive(ctx context.Context, filter InventoryHoldFilter) ([]*entities.InventoryHold, *PaginationResult, error)
	CleanupExpired(ctx context.Context) (int, error)
	ExtendExpiry(ctx context.Context, holdID uuid.UUID, newExpiryTime time.Time) error
	GetHoldsByExpiry(ctx context.Context, expiryTime time.Time, filter InventoryHoldFilter) ([]*entities.InventoryHold, *PaginationResult, error)

	// Statistics and analytics
	GetStats(ctx context.Context, filter InventoryHoldFilter) (*InventoryHoldStats, error)
	
	// Additional methods needed by order service
	DeleteBySessionID(ctx context.Context, sessionID string) error
	GetTotalHeld(ctx context.Context, ticketTierID uuid.UUID) (int, error)
}

