package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrderLine represents a line item in an order
type OrderLine struct {
	ID           uuid.UUID `json:"id" db:"id"`
	OrderID      uuid.UUID `json:"order_id" db:"order_id"`
	TicketTierID uuid.UUID `json:"ticket_tier_id" db:"ticket_tier_id"`
	Quantity     int       `json:"quantity" db:"quantity"`
	UnitPrice    float64   `json:"unit_price" db:"unit_price"`
	Subtotal     float64   `json:"subtotal" db:"subtotal"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	Order      *Order      `json:"order,omitempty"`
	TicketTier *TicketTier `json:"ticket_tier,omitempty"`
}

// NewOrderLine creates a new order line
func NewOrderLine(orderID, ticketTierID uuid.UUID, quantity int, unitPrice float64) *OrderLine {
	subtotal := unitPrice * float64(quantity)
	return &OrderLine{
		ID:           uuid.New(),
		OrderID:      orderID,
		TicketTierID: ticketTierID,
		Quantity:     quantity,
		UnitPrice:    unitPrice,
		Subtotal:     subtotal,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// GetTotal returns the total price for this line item
func (ol *OrderLine) GetTotal() float64 {
	return ol.Subtotal
}

// Validate validates the order line
func (ol *OrderLine) Validate() error {
	if ol.Quantity <= 0 {
		return NewValidationError("quantity", "quantity must be greater than zero")
	}

	if ol.UnitPrice < 0 {
		return NewValidationError("unit_price", "unit price cannot be negative")
	}

	if ol.Subtotal < 0 {
		return NewValidationError("subtotal", "subtotal cannot be negative")
	}

	return nil
}

