package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// OrganizerRepository defines the interface for organizer persistence operations
type OrganizerRepository interface {
	// Create creates a new organizer
	Create(ctx context.Context, organizer *entities.Organizer) error
	
	// GetByID retrieves an organizer by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Organizer, error)
	
	// GetBySlug retrieves an organizer by slug
	GetBySlug(ctx context.Context, slug string) (*entities.Organizer, error)
	
	// GetByEmail retrieves an organizer by email
	GetByEmail(ctx context.Context, email string) (*entities.Organizer, error)
	
	// Update updates an existing organizer
	Update(ctx context.Context, organizer *entities.Organizer) error
	
	// Delete soft deletes an organizer
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves organizers with pagination and filtering
	List(ctx context.Context, filter OrganizerFilter) ([]*entities.Organizer, *PaginationResult, error)
	
	// Exists checks if an organizer exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsBySlug checks if an organizer exists by slug
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	
	// ExistsByEmail checks if an organizer exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// OrganizerFilter defines filtering options for organizer queries
type OrganizerFilter struct {
	// Pagination
	Page  int
	Limit int
	
	// Filtering
	IsActive    *bool
	Search      string // Search in name, email
	Country     string // Filter by country
	City        string // Filter by city
	CreatedFrom *time.Time // Filter by creation date from
	CreatedTo   *time.Time // Filter by creation date to
	
	// Sorting
	SortBy    string // name, email, created_at
	SortOrder string // asc, desc
}

// OrganizerStats represents statistics for organizers
type OrganizerStats struct {
	TotalOrganizers       int     `json:"total_organizers" db:"total_organizers"`
	ActiveOrganizers      int     `json:"active_organizers" db:"active_organizers"`
	VerifiedOrganizers    int     `json:"verified_organizers" db:"verified_organizers"`
	TotalEvents           int     `json:"total_events" db:"total_events"`
	TotalTicketsSold      int     `json:"total_tickets_sold" db:"total_tickets_sold"`
	TotalRevenue          float64 `json:"total_revenue" db:"total_revenue"`
	AvgEventsPerOrganizer float64 `json:"avg_events_per_organizer" db:"avg_events_per_organizer"`
	AvgRevenuePerOrganizer float64 `json:"avg_revenue_per_organizer" db:"avg_revenue_per_organizer"`
}

