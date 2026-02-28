package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrderLine represents a line item in an order
type OrderLine struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrderID        uuid.UUID `json:"order_id" db:"order_id"`
	TicketTierID   uuid.UUID `json:"ticket_tier_id" db:"ticket_tier_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	UnitPrice      float64   `json:"unit_price" db:"unit_price"`
	Subtotal       float64   `json:"subtotal" db:"subtotal"`
	TotalPrice     float64   `json:"total_price" db:"total_price"` // Alias for subtotal (DB compatibility)
	Fees           float64   `json:"fees" db:"fees"`
	Taxes          float64   `json:"taxes" db:"taxes"`
	DiscountAmount float64   `json:"discount_amount" db:"discount_amount"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	// Denormalised fields (populated by JOIN queries in GetByOrder/GetByID)
	TicketTierName        string  `json:"ticket_tier_name,omitempty" db:"ticket_tier_name"`
	TicketTierDescription string  `json:"ticket_tier_description,omitempty" db:"ticket_tier_description"`
	TicketTierPrice       float64 `json:"ticket_tier_price,omitempty" db:"ticket_tier_price"`
	TicketTierQuota       int     `json:"ticket_tier_quota,omitempty" db:"ticket_tier_quota"`
	TicketTierEventID     string  `json:"ticket_tier_event_id,omitempty" db:"ticket_tier_event_id"`
	EventTitle            string  `json:"event_title,omitempty" db:"event_title"`
	EventSlug             string  `json:"event_slug,omitempty" db:"event_slug"`

	// Relations
	Order      *Order      `json:"order,omitempty"`
	TicketTier *TicketTier `json:"ticket_tier,omitempty"`
}

// NewOrderLine creates a new order line
func NewOrderLine(orderID, ticketTierID uuid.UUID, quantity int, unitPrice float64) *OrderLine {
	subtotal := unitPrice * float64(quantity)
	return &OrderLine{
		ID:             uuid.New(),
		OrderID:        orderID,
		TicketTierID:   ticketTierID,
		Quantity:       quantity,
		UnitPrice:      unitPrice,
		Subtotal:       subtotal,
		TotalPrice:     subtotal, // Keep in sync with subtotal
		Fees:           0,
		Taxes:          0,
		DiscountAmount: 0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// GetTotal returns the total price for this line item including fees and taxes minus discounts
func (ol *OrderLine) GetTotal() float64 {
	return ol.Subtotal + ol.Fees + ol.Taxes - ol.DiscountAmount
}

// GetGrossTotal returns the subtotal before fees, taxes, and discounts
func (ol *OrderLine) GetGrossTotal() float64 {
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

	if ol.Fees < 0 {
		return NewValidationError("fees", "fees cannot be negative")
	}

	if ol.Taxes < 0 {
		return NewValidationError("taxes", "taxes cannot be negative")
	}

	if ol.DiscountAmount < 0 {
		return NewValidationError("discount_amount", "discount amount cannot be negative")
	}

	return nil
}
