package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// EventRepository defines the interface for event persistence operations
type EventRepository interface {
	// Create creates a new event
	Create(ctx context.Context, event *entities.Event) error
	
	// GetByID retrieves an event by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Event, error)
	
	// GetBySlug retrieves an event by organizer ID and slug
	GetBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (*entities.Event, error)
	
	// Update updates an existing event
	Update(ctx context.Context, event *entities.Event) error
	
	// Delete soft deletes an event
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves events with pagination and filtering
	List(ctx context.Context, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// ListPublic retrieves public events (published/on_sale) with pagination and filtering
	ListPublic(ctx context.Context, filter PublicEventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// GetByOrganizer retrieves events for a specific organizer
	GetByOrganizer(ctx context.Context, organizerID uuid.UUID, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// GetByTour retrieves events for a specific tour
	GetByTour(ctx context.Context, tourID uuid.UUID, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// GetUpcoming retrieves upcoming events
	GetUpcoming(ctx context.Context, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// GetByCity retrieves events in a specific city
	GetByCity(ctx context.Context, city string, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// GetByDateRange retrieves events within a date range
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, filter EventFilter) ([]*entities.Event, *PaginationResult, error)
	
	// Exists checks if an event exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsBySlug checks if an event exists by organizer ID and slug
	ExistsBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (bool, error)
	
	// GetEventStats retrieves statistics for an event
	GetEventStats(ctx context.Context, eventID uuid.UUID) (*EventStats, error)
	
	// UpdateStatus updates the event status
	UpdateStatus(ctx context.Context, eventID uuid.UUID, status entities.EventStatus) error
}

// EventFilter defines filtering options for event queries
type EventFilter struct {
	BaseFilter
	
	// Filtering
	OrganizerID *uuid.UUID
	TourID      *uuid.UUID
	Status      *entities.EventStatus
	City        string
	Country     string
	IsActive    *bool
	Search      string // Search in name, venue_name, venue_city
	Tags        []string // Filter by tags
	
	// Date filtering
	EventDateFrom *time.Time
	EventDateTo   *time.Time
	SaleStart     *time.Time
	SaleEnd       *time.Time
	
	// Include related data
	IncludeTour         bool
	IncludeOrganizer    bool
	IncludeTicketTiers  bool
	IncludeStats        bool
}

// PublicEventFilter defines filtering options for public event queries
type PublicEventFilter struct {
	BaseFilter
	
	// Public filtering (only published/on_sale events)
	City        string
	Country     string
	TourID      *uuid.UUID
	Search      string
	
	// Date filtering
	EventDateFrom *time.Time
	EventDateTo   *time.Time
	
	// Include related data
	IncludeTour        bool
	IncludeTicketTiers bool
	IncludeMinMaxPrice bool
}

// EventStats represents event statistics
type EventStats struct {
	EventID           uuid.UUID `json:"event_id"`
	TotalTickets      int       `json:"total_tickets"`
	SoldTickets       int       `json:"sold_tickets"`
	AvailableTickets  int       `json:"available_tickets"`
	ReservedTickets   int       `json:"reserved_tickets"`
	TotalRevenue      float64   `json:"total_revenue"`
	PendingRevenue    float64   `json:"pending_revenue"`
	ConfirmedRevenue  float64   `json:"confirmed_revenue"`
	RefundedRevenue   float64   `json:"refunded_revenue"`
	TotalOrders       int       `json:"total_orders"`
	PaidOrders        int       `json:"paid_orders"`
	PendingOrders     int       `json:"pending_orders"`
	CancelledOrders   int       `json:"cancelled_orders"`
	RefundedOrders    int       `json:"refunded_orders"`
	RedeemedTickets   int       `json:"redeemed_tickets"`
	VoidedTickets     int       `json:"voided_tickets"`
	LastSaleAt        *time.Time `json:"last_sale_at"`
	LastRedemptionAt  *time.Time `json:"last_redemption_at"`
}

// TourRepository defines the interface for tour persistence operations
type TourRepository interface {
	// Create creates a new tour
	Create(ctx context.Context, tour *entities.Tour) error
	
	// GetByID retrieves a tour by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Tour, error)
	
	// GetBySlug retrieves a tour by organizer ID and slug
	GetBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (*entities.Tour, error)
	
	// Update updates an existing tour
	Update(ctx context.Context, tour *entities.Tour) error
	
	// Delete soft deletes a tour
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves tours with pagination and filtering
	List(ctx context.Context, filter TourFilter) ([]*entities.Tour, *PaginationResult, error)
	
	// GetByOrganizer retrieves tours for a specific organizer
	GetByOrganizer(ctx context.Context, organizerID uuid.UUID, filter TourFilter) ([]*entities.Tour, *PaginationResult, error)
	
	// GetActive retrieves active tours
	GetActive(ctx context.Context, filter TourFilter) ([]*entities.Tour, *PaginationResult, error)
	
	// GetUpcoming retrieves upcoming tours
	GetUpcoming(ctx context.Context, filter TourFilter) ([]*entities.Tour, *PaginationResult, error)
	
	// GetOngoing retrieves ongoing tours
	GetOngoing(ctx context.Context, filter TourFilter) ([]*entities.Tour, *PaginationResult, error)
	
	// Exists checks if a tour exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsBySlug checks if a tour exists by organizer ID and slug
	ExistsBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (bool, error)
	
	// GetTourStats retrieves statistics for a tour
	GetTourStats(ctx context.Context, tourID uuid.UUID) (*TourStats, error)
}

// TourFilter defines filtering options for tour queries
type TourFilter struct {
	BaseFilter
	
	// Filtering
	OrganizerID *uuid.UUID
	ArtistName  string
	IsActive    *bool
	Search      string // Search in name, artist_name
	
	// Date filtering
	StartDateFrom *time.Time
	StartDateTo   *time.Time
	EndDateFrom   *time.Time
	EndDateTo     *time.Time
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	
	// Include related data
	IncludeOrganizer bool
	IncludeEvents    bool
	IncludeStats     bool
}

// TourStats represents tour statistics
type TourStats struct {
	TourID           uuid.UUID  `json:"tour_id"`
	TotalEvents      int        `json:"total_events"`
	PublishedEvents  int        `json:"published_events"`
	OnSaleEvents     int        `json:"on_sale_events"`
	SoldOutEvents    int        `json:"sold_out_events"`
	CompletedEvents  int        `json:"completed_events"`
	CancelledEvents  int        `json:"cancelled_events"`
	TotalTickets     int        `json:"total_tickets"`
	SoldTickets      int        `json:"sold_tickets"`
	TotalRevenue     float64    `json:"total_revenue"`
	ConfirmedRevenue float64    `json:"confirmed_revenue"`
	FirstEventDate   *time.Time `json:"first_event_date"`
	LastEventDate    *time.Time `json:"last_event_date"`
}

