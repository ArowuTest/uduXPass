package entities

import (
	"time"

	"github.com/google/uuid"
)

// TicketStatus represents the status of a ticket
type TicketStatus string

const (
	TicketStatusActive   TicketStatus = "active"
	TicketStatusRedeemed TicketStatus = "redeemed"
	TicketStatusVoided   TicketStatus = "voided"
)

// Ticket represents individual tickets generated from paid orders
type Ticket struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	OrderLineID     uuid.UUID     `json:"order_line_id" db:"order_line_id"`
	SerialNumber    string        `json:"serial_number" db:"serial_number"`
	QRCodeData      string        `json:"qr_code_data" db:"qr_code_data"`
	QRCodeImageURL  *string       `json:"qr_code_image_url,omitempty" db:"qr_code_image_url"`
	Status          TicketStatus  `json:"status" db:"status"`
	RedeemedAt      *time.Time    `json:"redeemed_at,omitempty" db:"redeemed_at"`
	RedeemedBy      *string       `json:"redeemed_by,omitempty" db:"redeemed_by"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

// NewTicket creates a new ticket with default values
func NewTicket(orderLineID uuid.UUID, serialNumber, qrCodeData string) *Ticket {
	now := time.Now()
	return &Ticket{
		ID:           uuid.New(),
		OrderLineID:  orderLineID,
		SerialNumber: serialNumber,
		QRCodeData:   qrCodeData,
		Status:       TicketStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Validate performs business rule validation for the ticket
func (t *Ticket) Validate() error {
	if t.SerialNumber == "" {
		return NewValidationError("serial_number", "serial number is required")
	}
	if t.QRCodeData == "" {
		return NewValidationError("qr_code_data", "QR code data is required")
	}
	if len(t.SerialNumber) > 50 {
		return NewValidationError("serial_number", "serial number must be 50 characters or less")
	}
	
	// Validate redemption data consistency
	if t.Status == TicketStatusRedeemed {
		if t.RedeemedAt == nil {
			return NewValidationError("redeemed_at", "redeemed tickets must have redemption timestamp")
		}
		if t.RedeemedBy == nil || *t.RedeemedBy == "" {
			return NewValidationError("redeemed_by", "redeemed tickets must have redeemer information")
		}
	} else {
		if t.RedeemedAt != nil {
			return NewValidationError("redeemed_at", "non-redeemed tickets cannot have redemption timestamp")
		}
		if t.RedeemedBy != nil {
			return NewValidationError("redeemed_by", "non-redeemed tickets cannot have redeemer information")
		}
	}
	
	return nil
}

// Redeem marks the ticket as redeemed
func (t *Ticket) Redeem(redeemedBy string) error {
	if t.Status != TicketStatusActive {
		return NewBusinessRuleError("business_rule", "only active tickets can be redeemed", nil)
	}
	
	now := time.Now()
	t.Status = TicketStatusRedeemed
	t.RedeemedAt = &now
	t.RedeemedBy = &redeemedBy
	t.UpdatedAt = now
	return nil
}

// Void marks the ticket as voided
func (t *Ticket) Void() error {
	if t.Status == TicketStatusRedeemed {
		return NewBusinessRuleError("business_rule", "redeemed tickets cannot be voided", nil)
	}
	
	t.Status = TicketStatusVoided
	t.UpdatedAt = time.Now()
	return nil
}

// Reactivate marks a voided ticket as active again
func (t *Ticket) Reactivate() error {
	if t.Status != TicketStatusVoided {
		return NewBusinessRuleError("business_rule", "only voided tickets can be reactivated", nil)
	}
	
	t.Status = TicketStatusActive
	t.UpdatedAt = time.Now()
	return nil
}

// IsActive checks if the ticket is active
func (t *Ticket) IsActive() bool {
	return t.Status == TicketStatusActive
}

// IsRedeemed checks if the ticket has been redeemed
func (t *Ticket) IsRedeemed() bool {
	return t.Status == TicketStatusRedeemed
}

// IsVoided checks if the ticket is voided
func (t *Ticket) IsVoided() bool {
	return t.Status == TicketStatusVoided
}

// CanBeRedeemed checks if the ticket can be redeemed
func (t *Ticket) CanBeRedeemed() bool {
	return t.Status == TicketStatusActive
}

// CanBeVoided checks if the ticket can be voided
func (t *Ticket) CanBeVoided() bool {
	return t.Status == TicketStatusActive
}

// GetRedeemedDuration returns how long ago the ticket was redeemed
func (t *Ticket) GetRedeemedDuration() *time.Duration {
	if t.RedeemedAt == nil {
		return nil
	}
	duration := time.Since(*t.RedeemedAt)
	return &duration
}

// Payment represents payment transaction tracking
type Payment struct {
	ID                      uuid.UUID              `json:"id" db:"id"`
	OrderID                 uuid.UUID              `json:"order_id" db:"order_id"`
	Provider                PaymentMethod          `json:"provider" db:"provider"`
	ProviderTransactionID   *string                `json:"provider_transaction_id,omitempty" db:"provider_transaction_id"`
	Amount                  float64                `json:"amount" db:"amount"`
	Currency                string                 `json:"currency" db:"currency"`
	Status                  PaymentStatus          `json:"status" db:"status"`
	ProviderResponse        map[string]interface{} `json:"provider_response" db:"provider_response"`
	WebhookReceivedAt       *time.Time             `json:"webhook_received_at,omitempty" db:"webhook_received_at"`
	CreatedAt               time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at" db:"updated_at"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// NewPayment creates a new payment with default values
func NewPayment(orderID uuid.UUID, provider PaymentMethod, amount float64, currency string) *Payment {
	now := time.Now().UTC()
	return &Payment{
		ID:               uuid.New(),
		OrderID:          orderID,
		Provider:         provider,
		Amount:           amount,
		Currency:         currency,
		Status:           PaymentStatusPending,
		ProviderResponse: nil, // Will be set when provider responds
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Validate performs business rule validation for the payment
func (p *Payment) Validate() error {
	if p.Amount < 0 {
		return NewValidationError("amount", "amount must be non-negative")
	}
	if p.Currency == "" {
		return NewValidationError("currency", "currency is required")
	}
	if p.Provider != PaymentMethodMoMo && p.Provider != PaymentMethodPaystack {
		return NewValidationError("provider", "invalid payment provider")
	}
	return nil
}

// SetProviderTransactionID sets the provider transaction ID
func (p *Payment) SetProviderTransactionID(transactionID string) {
	p.ProviderTransactionID = &transactionID
	p.UpdatedAt = time.Now()
}

// MarkCompleted marks the payment as completed
func (p *Payment) MarkCompleted() error {
	if p.Status != PaymentStatusPending {
		return NewBusinessRuleError("business_rule", "only pending payments can be marked as completed", nil)
	}
	p.Status = PaymentStatusCompleted
	p.UpdatedAt = time.Now()
	return nil
}

// MarkFailed marks the payment as failed
func (p *Payment) MarkFailed() error {
	if p.Status != PaymentStatusPending {
		return NewBusinessRuleError("business_rule", "only pending payments can be marked as failed", nil)
	}
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the payment
func (p *Payment) Cancel() error {
	if p.Status == PaymentStatusCompleted {
		return NewBusinessRuleError("business_rule", "completed payments cannot be cancelled", nil)
	}
	if p.Status == PaymentStatusRefunded {
		return NewBusinessRuleError("business_rule", "refunded payments cannot be cancelled", nil)
	}
	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
	return nil
}

// Refund marks the payment as refunded
func (p *Payment) Refund() error {
	if p.Status != PaymentStatusCompleted {
		return NewBusinessRuleError("business_rule", "only completed payments can be refunded", nil)
	}
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateProviderResponse updates the provider response data
func (p *Payment) UpdateProviderResponse(response map[string]interface{}) {
	if p.ProviderResponse == nil {
		p.ProviderResponse = make(map[string]interface{})
	}
	for key, value := range response {
		p.ProviderResponse[key] = value
	}
	p.UpdatedAt = time.Now()
}

// MarkWebhookReceived marks when the webhook was received
func (p *Payment) MarkWebhookReceived() {
	now := time.Now()
	p.WebhookReceivedAt = &now
	p.UpdatedAt = now
}

// IsCompleted checks if the payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsPending checks if the payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsFailed checks if the payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// CanBeRefunded checks if the payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusCompleted
}

