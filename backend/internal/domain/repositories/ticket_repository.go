package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// TicketRepository defines the interface for ticket persistence operations
type TicketRepository interface {
	// Create creates a new ticket
	Create(ctx context.Context, ticket *entities.Ticket) error
	
	// CreateBatch creates multiple tickets in a batch
	CreateBatch(ctx context.Context, tickets []*entities.Ticket) error
	
	// GetByID retrieves a ticket by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Ticket, error)
	
	// GetBySerialNumber retrieves a ticket by serial number
	GetBySerialNumber(ctx context.Context, serialNumber string) (*entities.Ticket, error)
	
	// GetByQRCode retrieves a ticket by QR code data
	GetByQRCode(ctx context.Context, qrCodeData string) (*entities.Ticket, error)
	
	// Update updates an existing ticket
	Update(ctx context.Context, ticket *entities.Ticket) error
	
	// Delete deletes a ticket
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves tickets with pagination and filtering
	List(ctx context.Context, filter TicketFilter) ([]*entities.Ticket, *PaginationResult, error)
	
	// GetByOrderLine retrieves tickets for a specific order line
	GetByOrderLine(ctx context.Context, orderLineID uuid.UUID) ([]*entities.Ticket, error)
	
	// GetByOrder retrieves tickets for a specific order
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Ticket, error)
	
	// GetByUser retrieves tickets for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID, filter TicketFilter) ([]*entities.Ticket, *PaginationResult, error)
	
	// GetByEvent retrieves tickets for a specific event
	GetByEvent(ctx context.Context, eventID uuid.UUID, filter TicketFilter) ([]*entities.Ticket, *PaginationResult, error)
	
	// GetUpcoming retrieves tickets for upcoming events
	GetUpcoming(ctx context.Context, userID uuid.UUID) ([]*entities.Ticket, error)
	
	// UpdateStatus updates the ticket status
	UpdateStatus(ctx context.Context, ticketID uuid.UUID, status entities.TicketStatus) error
	
	// MarkRedeemed marks a ticket as redeemed
	MarkRedeemed(ctx context.Context, ticketID uuid.UUID, redeemedBy string) error
	
	// MarkVoided marks a ticket as voided
	MarkVoided(ctx context.Context, ticketID uuid.UUID) error
	
	// GetTicketStats retrieves statistics for tickets
	GetTicketStats(ctx context.Context, filter TicketStatsFilter) (*TicketStats, error)
	
	// Exists checks if a ticket exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsBySerialNumber checks if a ticket exists by serial number
	ExistsBySerialNumber(ctx context.Context, serialNumber string) (bool, error)
	
	// ValidateForRedemption validates a ticket for redemption
	ValidateForRedemption(ctx context.Context, qrCodeData string, eventID uuid.UUID) (*TicketValidationResult, error)
}

// TicketTierRepository defines the interface for ticket tier persistence operations
type TicketTierRepository interface {
	// Create creates a new ticket tier
	Create(ctx context.Context, ticketTier *entities.TicketTier) error
	
	// GetByID retrieves a ticket tier by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TicketTier, error)
	
	// Update updates an existing ticket tier
	Update(ctx context.Context, ticketTier *entities.TicketTier) error
	
	// Delete soft deletes a ticket tier
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves ticket tiers with pagination and filtering
	List(ctx context.Context, filter TicketTierFilter) ([]*entities.TicketTier, *PaginationResult, error)
	
	// GetByEvent retrieves ticket tiers for a specific event
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]*entities.TicketTier, error)
	
	// GetActiveByEvent retrieves active ticket tiers for a specific event
	GetActiveByEvent(ctx context.Context, eventID uuid.UUID) ([]*entities.TicketTier, error)
	
	// GetAvailability retrieves availability information for ticket tiers
	GetAvailability(ctx context.Context, eventID uuid.UUID) ([]*TicketTierAvailability, error)
	
	// UpdatePosition updates the display position of ticket tiers
	UpdatePosition(ctx context.Context, tierPositions map[uuid.UUID]int) error
	
	// Exists checks if a ticket tier exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// GetTierStats retrieves statistics for a ticket tier
	GetTierStats(ctx context.Context, tierID uuid.UUID) (*TicketTierStats, error)
	
	// GetAvailableQuantity retrieves the available quantity for a ticket tier
	GetAvailableQuantity(ctx context.Context, ticketTierID uuid.UUID) (int, error)

	// IncrementSold atomically increments the sold count for a ticket tier by the given quantity
	IncrementSold(ctx context.Context, tierID uuid.UUID, quantity int) error
}

// PaymentRepository defines the interface for payment persistence operations
type PaymentRepository interface {
	// Create creates a new payment
	Create(ctx context.Context, payment *entities.Payment) error
	
	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	
	// GetByProviderTransactionID retrieves a payment by provider transaction ID
	GetByProviderTransactionID(ctx context.Context, provider entities.PaymentMethod, transactionID string) (*entities.Payment, error)
	
	// Update updates an existing payment
	Update(ctx context.Context, payment *entities.Payment) error
	
	// Delete deletes a payment
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves payments with pagination and filtering
	List(ctx context.Context, filter PaymentFilter) ([]*entities.Payment, *PaginationResult, error)
	
	// GetByOrder retrieves payments for a specific order
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error)
	
	// GetByProvider retrieves payments for a specific provider
	GetByProvider(ctx context.Context, provider entities.PaymentMethod, filter PaymentFilter) ([]*entities.Payment, *PaginationResult, error)
	
	// UpdateStatus updates the payment status
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entities.PaymentStatus) error
	
	// MarkWebhookReceived marks when a webhook was received
	MarkWebhookReceived(ctx context.Context, paymentID uuid.UUID) error
	
	// GetPaymentStats retrieves payment statistics
	GetPaymentStats(ctx context.Context, filter PaymentStatsFilter) (*PaymentStats, error)
	
	// Exists checks if a payment exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsByProviderTransactionID checks if a payment exists by provider transaction ID
	ExistsByProviderTransactionID(ctx context.Context, provider entities.PaymentMethod, transactionID string) (bool, error)
}

// Filter definitions
type TicketFilter struct {
	BaseFilter
	
	// Filtering
	OrderID       *uuid.UUID
	OrderLineID   *uuid.UUID
	EventID       *uuid.UUID
	TicketTierID  *uuid.UUID
	UserID        *uuid.UUID
	Status        *entities.TicketStatus
	AttendeeEmail string
	
	// Search
	Search string
	
	// Date filtering
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	RedeemedFrom *time.Time
	RedeemedTo   *time.Time
	
	// Include related data
	IncludeOrder      bool
	IncludeTicketTier bool
	IncludeEvent      bool
	IncludeOrderLine  bool
}

type TicketTierFilter struct {
	BaseFilter
	
	// Filtering
	EventID   *uuid.UUID
	IsActive  *bool
	OnSale    *bool
	Search    string // Search in name
	
	// Price filtering
	MinPrice  *float64
	MaxPrice  *float64
	
	// Date filtering
	SaleStartFrom *time.Time
	SaleStartTo   *time.Time
	SaleEndFrom   *time.Time
	SaleEndTo     *time.Time
	
	// Availability filtering
	AvailableOnly bool
	
	// Include related data
	IncludeEvent        bool
	IncludeAvailability bool
	IncludeStats        bool
}

// TicketTierStatsFilter for filtering ticket tier statistics
type TicketTierStatsFilter struct {
	EventID *uuid.UUID
	TierID  *uuid.UUID
}

// TicketTierStats represents statistics for a ticket tier
type TicketTierStats struct {
	TierID        uuid.UUID `json:"tier_id" db:"tier_id"`
	TierName      string    `json:"tier_name" db:"tier_name"`
	Price         float64   `json:"price" db:"price"`
	Currency      string    `json:"currency" db:"currency"`
	Capacity      *int      `json:"capacity" db:"capacity"`
	SoldCount     int       `json:"sold_count" db:"sold_count"`
	Revenue       float64   `json:"revenue" db:"revenue"`
	ReservedCount int       `json:"reserved_count" db:"reserved_count"`
	AvailableCount int      `json:"available_count" db:"available_count"`
}

// TicketStatsFilter for filtering ticket statistics
type TicketStatsFilter struct {
	EventID *uuid.UUID
	UserID  *uuid.UUID
	Status  *entities.TicketStatus
}

// TicketStats represents ticket statistics
type TicketStats struct {
	TotalTickets     int     `json:"total_tickets" db:"total_tickets"`
	ActiveTickets    int     `json:"active_tickets" db:"active_tickets"`
	RedeemedTickets  int     `json:"redeemed_tickets" db:"redeemed_tickets"`
	CancelledTickets int     `json:"cancelled_tickets" db:"cancelled_tickets"`
	TotalValue       float64 `json:"total_value" db:"total_value"`
}

type PaymentFilter struct {
	BaseFilter
	
	// Filtering
	OrderID       *uuid.UUID
	EventID       *uuid.UUID
	Provider      *entities.PaymentMethod
	PaymentMethod *entities.PaymentMethod
	Status        *entities.PaymentStatus
	Currency      string
	Search        string
	
	// Date filtering
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	ProcessedFrom *time.Time
	ProcessedTo   *time.Time
	
	// Amount filtering
	MinAmount   *float64
	MaxAmount   *float64
	
	// Include related data
	IncludeOrder bool
}

// Statistics and availability types
type TicketTierAvailability struct {
	TicketTierID uuid.UUID `json:"ticket_tier_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	Quota        *int      `json:"quota"`
	Sold         int       `json:"sold"`
	Reserved     int       `json:"reserved"`
	Available    int       `json:"available"`
	IsOnSale     bool      `json:"is_on_sale"`
	SaleStatus   string    `json:"sale_status"`
}

type PaymentStats struct {
	TotalPayments     int     `json:"total_payments"`
	CompletedPayments int     `json:"completed_payments"`
	FailedPayments    int     `json:"failed_payments"`
	TotalAmount       float64 `json:"total_amount"`
	CompletedAmount   float64 `json:"completed_amount"`
	SuccessRate       float64 `json:"success_rate"`
	AverageAmount     float64 `json:"average_amount"`
}

type PaymentStatsFilter struct {
	Provider  *entities.PaymentMethod
	DateFrom  *time.Time
	DateTo    *time.Time
	EventID   *uuid.UUID
}

type TicketValidationResult struct {
	Valid           bool                  `json:"valid"`
	Ticket          *entities.Ticket      `json:"ticket,omitempty"`
	Status          string                `json:"status"`
	Message         string                `json:"message"`
	CustomerName    string                `json:"customer_name,omitempty"`
	TicketTierName  string                `json:"ticket_tier_name,omitempty"`
	AlreadyRedeemed bool                  `json:"already_redeemed"`
	RedemptionTime  *time.Time            `json:"redemption_time,omitempty"`
}

