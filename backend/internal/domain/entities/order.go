package entities

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusExpired   OrderStatus = "expired"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodMoMo     PaymentMethod = "momo"
	PaymentMethodPaystack PaymentMethod = "paystack"
	PaymentMethodCard     PaymentMethod = "card"
)

// Order represents shopping cart and purchase transactions
type Order struct {
	ID                 uuid.UUID              `json:"id" db:"id"`
	Code               string                 `json:"code" db:"code"`
	EventID            string                 `json:"event_id" db:"event_id"` // Changed to string to match usage
	UserID             *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Email              string                 `json:"email" db:"email"`
	Phone              *string                `json:"phone,omitempty" db:"phone"`
	FirstName          *string                `json:"first_name,omitempty" db:"first_name"`
	LastName           *string                `json:"last_name,omitempty" db:"last_name"`
	
	// Customer fields that were missing
	CustomerFirstName  string                 `json:"customer_first_name" db:"customer_first_name"`
	CustomerLastName   string                 `json:"customer_last_name" db:"customer_last_name"`
	CustomerEmail      string                 `json:"customer_email" db:"customer_email"`
	CustomerPhone      string                 `json:"customer_phone" db:"customer_phone"`
	
	Status             OrderStatus            `json:"status" db:"status"`
	TotalAmount        float64                `json:"total_amount" db:"total_amount"`
	Currency           string                 `json:"currency" db:"currency"`
	PaymentMethod      *PaymentMethod         `json:"payment_method,omitempty" db:"payment_method"`
	PaymentID          *string                `json:"payment_id,omitempty" db:"payment_id"`
	PaymentReference   *string                `json:"payment_reference,omitempty" db:"payment_reference"`
	Notes              *string                `json:"notes,omitempty" db:"notes"`
	ConfirmedAt        *time.Time             `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CancelledAt        *time.Time             `json:"cancelled_at,omitempty" db:"cancelled_at"`
	PaidAt             *time.Time             `json:"paid_at,omitempty" db:"paid_at"`
	ExpiresAt          time.Time              `json:"expires_at" db:"expires_at"`
	Secret             string                 `json:"-" db:"secret"`
	Locale             string                 `json:"locale" db:"locale"`
	Comment            *string                `json:"comment,omitempty" db:"comment"`
	MetaInfo           map[string]interface{} `json:"meta_info" db:"meta_info"`
	IsActive           bool                   `json:"is_active" db:"is_active"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`

	// Relations
	Event      *Event       `json:"event,omitempty"`
	User       *User        `json:"user,omitempty"`
	OrderLines []OrderLine  `json:"order_lines,omitempty"`
	Payments   []Payment    `json:"payments,omitempty"`
}

// NewOrder creates a new order with generated code and secret
func NewOrder(eventID string, email string) *Order {
	return &Order{
		ID:        uuid.New(),
		Code:      generateOrderCode(),
		EventID:   eventID,
		Email:     email,
		Status:    OrderStatusPending,
		Currency:  "NGN",
		ExpiresAt: time.Now().Add(15 * time.Minute).UTC(), // Add 15 minutes THEN convert to UTC
		Secret:    generateSecret(),
		Locale:    "en",
		IsActive:  true, // Orders are active by default
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// generateOrderCode generates a unique order code
func generateOrderCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "ORD-" + hex.EncodeToString(bytes)
}

// generateSecret generates a secret for order verification
func generateSecret() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// IsExpired checks if the order has expired
func (o *Order) IsExpired() bool {
	return time.Now().UTC().After(o.ExpiresAt)
}

// CanBePaid checks if the order can be paid
func (o *Order) CanBePaid() bool {
	return o.Status == OrderStatusPending && !o.IsExpired()
}

// MarkAsPaid marks the order as paid
func (o *Order) MarkAsPaid() error {
	if !o.CanBePaid() {
		return NewBusinessRuleError("payment", "order cannot be paid", map[string]interface{}{
			"status":     o.Status,
			"expired":    o.IsExpired(),
			"expires_at": o.ExpiresAt,
		})
	}
	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if o.Status == OrderStatusPaid {
		return NewBusinessRuleError("cancellation", "paid orders cannot be cancelled", nil)
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// Refund refunds the order
func (o *Order) Refund() error {
	if o.Status != OrderStatusPaid {
		return NewBusinessRuleError("refund", "only paid orders can be refunded", nil)
	}
	o.Status = OrderStatusRefunded
	o.UpdatedAt = time.Now()
	return nil
}

// CalculateTotal calculates the total amount from order lines
func (o *Order) CalculateTotal() {
	total := 0.0
	for _, line := range o.OrderLines {
		total += line.Subtotal
	}
	o.TotalAmount = total
	o.UpdatedAt = time.Now()
}

// AddOrderLine adds an order line to the order
func (o *Order) AddOrderLine(ticketTierID uuid.UUID, quantity int, price float64) error {
	if o.Status != OrderStatusPending {
		return NewBusinessRuleError("order_modification", "only pending orders can be modified", nil)
	}

	subtotal := price * float64(quantity)
	orderLine := &OrderLine{
		ID:           uuid.New(),
		OrderID:      o.ID,
		TicketTierID: ticketTierID,
		Quantity:     quantity,
		UnitPrice:    price,
		Subtotal:     subtotal,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	o.OrderLines = append(o.OrderLines, *orderLine)
	o.CalculateTotal()
	return nil
}

// GetTotalTickets returns the total number of tickets in the order
func (o *Order) GetTotalTickets() int {
	total := 0
	for _, line := range o.OrderLines {
		total += line.Quantity
	}
	return total
}

// Validate validates the order
func (o *Order) Validate() error {
	if o.Email == "" {
		return NewValidationError("email", "email is required")
	}

	if o.TotalAmount <= 0 {
		return NewValidationError("total_amount", "total amount must be greater than zero")
	}

	if len(o.OrderLines) == 0 {
		return NewValidationError("order_lines", "order must have at least one order line")
	}

	return nil
}

// SetPaymentMethod sets the payment method for the order
func (o *Order) SetPaymentMethod(method PaymentMethod) {
	o.PaymentMethod = &method
}

// MarkPaid marks the order as paid
func (o *Order) MarkPaid() {
	o.Status = OrderStatusPaid
	now := time.Now()
	o.ConfirmedAt = &now
}

