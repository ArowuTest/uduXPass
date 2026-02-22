package entities

import (
	"time"

	"github.com/google/uuid"
)

// InventoryHoldStatus represents the status of an inventory hold
type InventoryHoldStatus string

const (
	InventoryHoldStatusActive    InventoryHoldStatus = "active"
	InventoryHoldStatusConfirmed InventoryHoldStatus = "confirmed"
	InventoryHoldStatusReleased  InventoryHoldStatus = "released"
	InventoryHoldStatusExpired   InventoryHoldStatus = "expired"
)

// InventoryHold represents a temporary hold on ticket inventory
type InventoryHold struct {
	ID           uuid.UUID           `json:"id" db:"id"`
	OrderID      uuid.UUID           `json:"order_id" db:"order_id"`
	TicketTierID uuid.UUID           `json:"ticket_tier_id" db:"ticket_tier_id"`
	Quantity     int                 `json:"quantity" db:"quantity"`
	Status       InventoryHoldStatus `json:"status" db:"status"`
	ExpiresAt    time.Time           `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	
	// Related entities (for joins)
	TicketTierName     string `json:"ticket_tier_name,omitempty" db:"ticket_tier_name"`
	TicketTierCapacity int    `json:"ticket_tier_capacity,omitempty" db:"ticket_tier_capacity"`
	EventTitle         string `json:"event_title,omitempty" db:"event_title"`
	EventSlug          string `json:"event_slug,omitempty" db:"event_slug"`
	OrderCode          string `json:"order_code,omitempty" db:"order_code"`
	OrderStatus        string `json:"order_status,omitempty" db:"order_status"`
}

// NewInventoryHold creates a new inventory hold
func NewInventoryHold(orderID, ticketTierID uuid.UUID, quantity int, duration time.Duration) *InventoryHold {
	return &InventoryHold{
		ID:           uuid.New(),
		OrderID:      orderID,
		TicketTierID: ticketTierID,
		Quantity:     quantity,
		ExpiresAt:    time.Now().Add(duration),
		CreatedAt:    time.Now(),
	}
}

// IsExpired checks if the inventory hold has expired
func (ih *InventoryHold) IsExpired() bool {
	return time.Now().After(ih.ExpiresAt)
}

// IsActive checks if the inventory hold is still active
func (ih *InventoryHold) IsActive() bool {
	return !ih.IsExpired()
}

// TimeUntilExpiry returns the duration until the hold expires
func (ih *InventoryHold) TimeUntilExpiry() time.Duration {
	if ih.IsExpired() {
		return 0
	}
	return time.Until(ih.ExpiresAt)
}

// Extend extends the expiry time of the hold
func (ih *InventoryHold) Extend(duration time.Duration) {
	ih.ExpiresAt = ih.ExpiresAt.Add(duration)
}

// Validate validates the inventory hold
func (ih *InventoryHold) Validate() error {
	if ih.ID == uuid.Nil {
		return ErrValidationError
	}
	
	if ih.OrderID == uuid.Nil {
		return ErrValidationError
	}
	
	if ih.TicketTierID == uuid.Nil {
		return ErrValidationError
	}
	
	if ih.Quantity <= 0 {
		return ErrValidationError
	}
	
	if ih.ExpiresAt.Before(time.Now()) {
		return ErrValidationError
	}
	
	return nil
}

