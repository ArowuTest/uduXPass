package repositories

import (
	"context"
	"database/sql/driver"
)

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewPaginationResult creates a new pagination result
func NewPaginationResult(page, limit, total int) *PaginationResult {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	
	return &PaginationResult{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction defines the interface for database transactions
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
	
	// Context returns the transaction context
	Context() context.Context
	
	// Organizers returns the organizer repository within this transaction
	Organizers() OrganizerRepository
	
	// Tours returns the tour repository within this transaction
	Tours() TourRepository
	
	// Events returns the event repository within this transaction
	Events() EventRepository
	
	// TicketTiers returns the ticket tier repository within this transaction
	TicketTiers() TicketTierRepository
	
	// Users returns the user repository within this transaction
	Users() UserRepository
	
	// Orders returns the order repository within this transaction
	Orders() OrderRepository
	
	// OrderLines returns the order line repository within this transaction
	OrderLines() OrderLineRepository
	
	// Tickets returns the ticket repository within this transaction
	Tickets() TicketRepository
	
	// Payments returns the payment repository within this transaction
	Payments() PaymentRepository
	
	// InventoryHolds returns the inventory hold repository within this transaction
	InventoryHolds() InventoryHoldRepository
	
	// OTPTokens returns the OTP token repository within this transaction
	OTPTokens() OTPTokenRepository
}

// RepositoryManager defines the interface for accessing all repositories
type RepositoryManager interface {
	// UnitOfWork returns the unit of work for transaction management
	UnitOfWork() UnitOfWork
	
	// Organizers returns the organizer repository
	Organizers() OrganizerRepository
	
	// Tours returns the tour repository
	Tours() TourRepository
	
	// Events returns the event repository
	Events() EventRepository
	
	// TicketTiers returns the ticket tier repository
	TicketTiers() TicketTierRepository
	
	// Users returns the user repository
	Users() UserRepository
	
	// Orders returns the order repository
	Orders() OrderRepository
	
	// OrderLines returns the order line repository
	OrderLines() OrderLineRepository
	
	// Tickets returns the ticket repository
	Tickets() TicketRepository
	
	// Payments returns the payment repository
	Payments() PaymentRepository
	
	// InventoryHolds returns the inventory hold repository
	InventoryHolds() InventoryHoldRepository
	
	// OTPTokens returns the OTP token repository
	OTPTokens() OTPTokenRepository
	
	// ScannerUsers returns the scanner user repository
	ScannerUsers() ScannerUserRepository
	
	// Close closes all repository connections
	Close() error
}

// SortOrder represents sort order options
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// Validate validates the sort order
func (so SortOrder) Validate() bool {
	return so == SortOrderAsc || so == SortOrderDesc
}

// String returns the string representation
func (so SortOrder) String() string {
	return string(so)
}

// Value implements the driver.Valuer interface
func (so SortOrder) Value() (driver.Value, error) {
	return string(so), nil
}

// BaseFilter provides common filtering options
type BaseFilter struct {
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
	SortBy    string    `json:"sort_by"`
	SortOrder SortOrder `json:"sort_order"`
}

// Validate validates the base filter
func (bf *BaseFilter) Validate() error {
	if bf.Page < 1 {
		bf.Page = 1
	}
	if bf.Limit < 1 {
		bf.Limit = 20
	}
	if bf.Limit > 100 {
		bf.Limit = 100
	}
	if bf.SortOrder == "" {
		bf.SortOrder = SortOrderDesc
	}
	if !bf.SortOrder.Validate() {
		bf.SortOrder = SortOrderDesc
	}
	return nil
}

// GetOffset calculates the offset for pagination
func (bf *BaseFilter) GetOffset() int {
	return (bf.Page - 1) * bf.Limit
}

