package entities

import (
	"time"

	"github.com/google/uuid"
)

// TicketTier represents product definitions for tickets
type TicketTier struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	EventID      uuid.UUID  `json:"event_id" db:"event_id"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description,omitempty" db:"description"`
	Price        float64    `json:"price" db:"price"`
	Currency     string     `json:"currency" db:"currency"`
	Quota        int        `json:"quota" db:"quota"`
	Sold         int        `json:"sold" db:"sold"`
	MinPurchase  int        `json:"min_purchase" db:"min_purchase"`
	MaxPurchase  int        `json:"max_purchase" db:"max_purchase"`
	SaleStart    *time.Time `json:"sale_start,omitempty" db:"sale_start"`
	SaleEnd      *time.Time `json:"sale_end,omitempty" db:"sale_end"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// NewTicketTier creates a new ticket tier with default values
func NewTicketTier(eventID uuid.UUID, name string, price float64) *TicketTier {
	now := time.Now()
	return &TicketTier{
		ID:          uuid.New(),
		EventID:     eventID,
		Name:        name,
		Price:       price,
		Quota:       100,
		Sold:        0,
		MaxPurchase: 10,
		MinPurchase: 1,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate performs business rule validation for the ticket tier
func (tt *TicketTier) Validate() error {
	if tt.Name == "" {
		return NewValidationError("name", "name is required")
	}
	if tt.Price < 0 {
		return NewValidationError("price", "price must be non-negative")
	}
	if tt.MaxPurchase <= 0 {
		return NewValidationError("max_purchase", "max purchase must be positive")
	}
	if tt.MinPurchase <= 0 {
		return NewValidationError("min_purchase", "min purchase must be positive")
	}
	if tt.MinPurchase > tt.MaxPurchase {
		return NewValidationError("purchase_limits", "min purchase cannot exceed max purchase")
	}
	if tt.Quota <= 0 {
		return NewValidationError("quota", "quota must be positive")
	}
	
	// Validate sale dates if provided
	if tt.SaleStart != nil && tt.SaleEnd != nil && tt.SaleStart.After(*tt.SaleEnd) {
		return NewValidationError("sale_dates", "sale start must be before sale end")
	}
	
	return nil
}

// SetQuota sets the quota limit for the ticket tier
func (tt *TicketTier) SetQuota(quota int) error {
	if quota <= 0 {
		return NewValidationError("quota", "quota must be positive")
	}
	tt.Quota = quota
	tt.UpdatedAt = time.Now()
	return nil
}

// SetSalePeriod sets the sale start and end times
func (tt *TicketTier) SetSalePeriod(start, end time.Time) error {
	if start.After(end) {
		return NewValidationError("sale_period", "sale start must be before sale end")
	}
	tt.SaleStart = &start
	tt.SaleEnd = &end
	tt.UpdatedAt = time.Now()
	return nil
}

// SetPurchaseLimits sets the minimum and maximum tickets per purchase
func (tt *TicketTier) SetPurchaseLimits(min, max int) error {
	if min <= 0 {
		return NewValidationError("min_purchase", "min purchase must be positive")
	}
	if max <= 0 {
		return NewValidationError("max_purchase", "max purchase must be positive")
	}
	if min > max {
		return NewValidationError("purchase_limits", "min purchase cannot exceed max purchase")
	}
	tt.MinPurchase = min
	tt.MaxPurchase = max
	tt.UpdatedAt = time.Now()
	return nil
}

// UpdatePrice updates the ticket tier price
func (tt *TicketTier) UpdatePrice(price float64) error {
	if price < 0 {
		return NewValidationError("price", "price must be non-negative")
	}
	tt.Price = price
	tt.UpdatedAt = time.Now()
	return nil
}

// IsOnSale checks if the ticket tier is currently on sale
func (tt *TicketTier) IsOnSale() bool {
	if !tt.IsActive {
		return false
	}
	
	now := time.Now()
	
	// Check sale start time
	if tt.SaleStart != nil && now.Before(*tt.SaleStart) {
		return false
	}
	
	// Check sale end time
	if tt.SaleEnd != nil && now.After(*tt.SaleEnd) {
		return false
	}
	
	return true
}

// IsValidQuantity checks if the requested quantity is valid for this tier
func (tt *TicketTier) IsValidQuantity(quantity int) bool {
	return quantity >= tt.MinPurchase && quantity <= tt.MaxPurchase
}

// GetAvailableQuantity returns the number of tickets still available
func (tt *TicketTier) GetAvailableQuantity() int {
	return tt.Quota - tt.Sold
}

// GetQuota returns the total quota
func (tt *TicketTier) GetQuota() int {
	return tt.Quota
}

// Activate marks the ticket tier as active
func (tt *TicketTier) Activate() {
	tt.IsActive = true
	tt.UpdatedAt = time.Now()
}

// Deactivate marks the ticket tier as inactive
func (tt *TicketTier) Deactivate() {
	tt.IsActive = false
	tt.UpdatedAt = time.Now()
}

// IsFree checks if the ticket tier is free
func (tt *TicketTier) IsFree() bool {
	return tt.Price == 0
}
