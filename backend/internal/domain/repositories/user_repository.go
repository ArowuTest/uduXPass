package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// GetByPhone retrieves a user by phone number
	GetByPhone(ctx context.Context, phone string) (*entities.User, error)
	
	// GetByMoMoID retrieves a user by MoMo ID
	GetByMoMoID(ctx context.Context, momoID string) (*entities.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves users with pagination and filtering
	List(ctx context.Context, filter UserFilter) ([]*entities.User, *PaginationResult, error)
	
	// Exists checks if a user exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	
	// ExistsByEmail checks if a user exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	
	// ExistsByPhone checks if a user exists by phone
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	
	// ExistsByMoMoID checks if a user exists by MoMo ID
	ExistsByMoMoID(ctx context.Context, momoID string) (bool, error)
	
	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	
	// VerifyEmail marks the user's email as verified
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	
	// VerifyPhone marks the user's phone as verified
	VerifyPhone(ctx context.Context, userID uuid.UUID) error
	
	// UpdatePassword updates the user's password hash
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	
	// IncrementFailedAttempts increments the failed login attempts counter
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	
	// ResetFailedAttempts resets the failed login attempts counter
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	
	// GetUserStats retrieves statistics for a user
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
}

// UserFilter defines filtering options for user queries
type UserFilter struct {
	BaseFilter
	
	// Filtering
	AuthProvider    *entities.AuthProvider
	EmailVerified   *bool
	PhoneVerified   *bool
	IsActive        *bool
	Search          string // Search in email, phone, first_name, last_name
	
	// Date filtering
	CreatedFrom     *time.Time
	CreatedTo       *time.Time
	LastLoginFrom   *time.Time
	LastLoginTo     *time.Time
}

// UserStats represents user statistics
type UserStats struct {
	UserID              uuid.UUID  `json:"user_id"`
	TotalOrders         int        `json:"total_orders"`
	PaidOrders          int        `json:"paid_orders"`
	CancelledOrders     int        `json:"cancelled_orders"`
	RefundedOrders      int        `json:"refunded_orders"`
	TotalTickets        int        `json:"total_tickets"`
	RedeemedTickets     int        `json:"redeemed_tickets"`
	TotalSpent          float64    `json:"total_spent"`
	AverageOrderValue   float64    `json:"average_order_value"`
	FirstOrderAt        *time.Time `json:"first_order_at"`
	LastOrderAt         *time.Time `json:"last_order_at"`
	FavoriteVenueCity   string     `json:"favorite_venue_city"`
	EventsAttended      int        `json:"events_attended"`
}

