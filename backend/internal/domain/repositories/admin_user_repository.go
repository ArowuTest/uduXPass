package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// AdminUserRepository defines the interface for admin user data access
type AdminUserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, adminUser *entities.AdminUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*entities.AdminUser, error)
	Update(ctx context.Context, adminUser *entities.AdminUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List and search operations
	List(ctx context.Context, limit, offset int) ([]*entities.AdminUser, error)
	ListByRole(ctx context.Context, role entities.AdminRole, limit, offset int) ([]*entities.AdminUser, error)
	ListByStatus(ctx context.Context, status entities.AdminStatus, limit, offset int) ([]*entities.AdminUser, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.AdminUser, error)
	Count(ctx context.Context) (int64, error)
	CountWithFilter(ctx context.Context, filter AdminUserFilter) (int64, error)
	CountByRole(ctx context.Context, role entities.AdminRole) (int64, error)
	CountByStatus(ctx context.Context, status entities.AdminStatus) (int64, error)
	
	// Authentication operations
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, id uuid.UUID) error
	LockAccount(ctx context.Context, id uuid.UUID, lockUntil *time.Time) error
	UnlockAccount(ctx context.Context, id uuid.UUID) error
	
	// Permission operations
	UpdatePermissions(ctx context.Context, id uuid.UUID, permissions []entities.AdminPermission) error
	HasPermission(ctx context.Context, id uuid.UUID, permission entities.AdminPermission) (bool, error)
	
	// Status operations
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.AdminStatus) error
	
	// Two-factor authentication
	EnableTwoFactor(ctx context.Context, id uuid.UUID, secret string) error
	DisableTwoFactor(ctx context.Context, id uuid.UUID) error
	UpdateTwoFactorSecret(ctx context.Context, id uuid.UUID, secret string) error
	
	// Audit operations
	GetLoginHistory(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.AdminLoginHistory, error)
	LogLogin(ctx context.Context, adminID uuid.UUID, ipAddress, userAgent string, success bool) error
	
	// Bulk operations
	BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status entities.AdminStatus) error
	BulkDelete(ctx context.Context, ids []uuid.UUID) error
	
	// Analytics
	GetActiveAdminsCount(ctx context.Context) (int64, error)
	GetAdminsByCreatedDate(ctx context.Context, startDate, endDate time.Time) ([]*entities.AdminUser, error)
	GetMostActiveAdmins(ctx context.Context, limit int) ([]*entities.AdminUser, error)
}

// AdminLoginHistory represents admin login history
type AdminLoginHistory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AdminID   uuid.UUID `json:"admin_id" db:"admin_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Success   bool      `json:"success" db:"success"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AdminUserFilter represents filtering options for admin users
type AdminUserFilter struct {
	Role          *entities.AdminRole   `json:"role,omitempty"`
	Status        *entities.AdminStatus `json:"status,omitempty"`
	IsActive      *bool                 `json:"is_active,omitempty"`
	Email         *string               `json:"email,omitempty"`
	CreatedAt     *time.Time            `json:"created_at,omitempty"`
	Search        string                `json:"search,omitempty"`
	SortBy        string                `json:"sort_by,omitempty"`
	SortDirection string                `json:"sort_direction,omitempty"`
	Limit         int                   `json:"limit,omitempty"`
	Offset        int                   `json:"offset,omitempty"`
}

// AdminUserStats represents statistics about admin users
type AdminUserStats struct {
	TotalAdmins    int64                        `json:"total_admins"`
	ActiveAdmins   int64                        `json:"active_admins"`
	InactiveAdmins int64                        `json:"inactive_admins"`
	LockedAdmins   int64                        `json:"locked_admins"`
	SuperAdmins    int64                        `json:"super_admins"`
	RegularAdmins  int64                        `json:"regular_admins"`
	AdminsByRole   map[entities.AdminRole]int64 `json:"admins_by_role"`
}

// LoginRecord represents a login record for audit purposes
type LoginRecord struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AdminID   uuid.UUID `json:"admin_id" db:"admin_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Success   bool      `json:"success" db:"success"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
