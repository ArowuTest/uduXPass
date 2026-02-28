package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// ScannerUserRepository defines the interface for scanner user data operations
type ScannerUserRepository interface {
	// Authentication
	GetByUsername(ctx context.Context, username string) (*entities.ScannerUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ScannerUser, error)
	GetByEmail(ctx context.Context, email string) (*entities.ScannerUser, error)

	// CRUD operations
	Create(ctx context.Context, scanner *entities.ScannerUser) error
	Update(ctx context.Context, scanner *entities.ScannerUser) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List and search
	List(ctx context.Context, filter *ScannerUserFilter) ([]*entities.ScannerUser, *PaginationResult, error)
	Search(ctx context.Context, query string, filter *ScannerUserFilter) ([]*entities.ScannerUser, *PaginationResult, error)

	// Event assignments
	AssignToEvent(ctx context.Context, scannerID, eventID, assignedBy uuid.UUID) error
	UnassignFromEvent(ctx context.Context, scannerID, eventID uuid.UUID) error
	GetAssignedEvents(ctx context.Context, scannerID uuid.UUID) ([]*entities.ScannerAssignedEvent, error)
	GetEventScanners(ctx context.Context, eventID uuid.UUID) ([]*entities.ScannerUser, error)

	// Session management
	CreateSession(ctx context.Context, session *entities.ScannerSession) error
	GetActiveSession(ctx context.Context, scannerID uuid.UUID) (*entities.ScannerSession, error)
	EndSession(ctx context.Context, sessionID uuid.UUID) error
	UpdateSessionStats(ctx context.Context, sessionID uuid.UUID, scansCount, validScans, invalidScans int, totalRevenue float64) error

	// Audit and logging
	LogActivity(ctx context.Context, log *entities.ScannerAuditLog) error
	RecordLogin(ctx context.Context, loginHistory *entities.ScannerLoginHistory) error
	GetLoginHistory(ctx context.Context, scannerID uuid.UUID, limit int) ([]*entities.ScannerLoginHistory, error)
	GetAuditLog(ctx context.Context, scannerID uuid.UUID, filter *ScannerAuditFilter) ([]*entities.ScannerAuditLog, *PaginationResult, error)

	// Ticket validation
	ValidateTicket(ctx context.Context, validation *entities.TicketValidation) error
	GetValidationHistory(ctx context.Context, scannerID uuid.UUID, filter *TicketValidationFilter) ([]*entities.TicketValidation, *PaginationResult, error)

	// Statistics
	GetScannerStats(ctx context.Context, scannerID uuid.UUID, eventID *uuid.UUID) (*ScannerStats, error)
	GetEventScanStats(ctx context.Context, eventID uuid.UUID) (*EventScanStats, error)
}

// ScannerUserFilter represents filtering options for scanner users
type ScannerUserFilter struct {
	BaseFilter
	Role      *entities.ScannerRole   `json:"role,omitempty"`
	Status    *entities.ScannerStatus `json:"status,omitempty"`
	EventID   *uuid.UUID              `json:"event_id,omitempty"`
	CreatedBy *uuid.UUID              `json:"created_by,omitempty"`
	Search    string                  `json:"search,omitempty"`
}

// ScannerAuditFilter represents filtering options for scanner audit logs
type ScannerAuditFilter struct {
	BaseFilter
	ScannerID    *uuid.UUID `json:"scanner_id,omitempty"`
	SessionID    *uuid.UUID `json:"session_id,omitempty"`
	Action       string     `json:"action,omitempty"`
	ResourceType string     `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	DateFrom     *string    `json:"date_from,omitempty"`
	DateTo       *string    `json:"date_to,omitempty"`
}

// TicketValidationFilter represents filtering options for ticket validations
type TicketValidationFilter struct {
	BaseFilter
	ScannerID        *uuid.UUID `json:"scanner_id,omitempty"`
	SessionID        *uuid.UUID `json:"session_id,omitempty"`
	ValidationResult string     `json:"validation_result,omitempty"`
	DateFrom         *string    `json:"date_from,omitempty"`
	DateTo           *string    `json:"date_to,omitempty"`
}

// ScannerStats represents statistics for a scanner
// db tags must match the SQL column aliases in GetScannerStats query
type ScannerStats struct {
	ScannerID      uuid.UUID  `json:"scanner_id" db:"scanner_id"`
	TotalSessions  int        `json:"total_sessions" db:"total_sessions"`
	TotalScans     int        `json:"total_scans" db:"total_scans"`
	ValidScans     int        `json:"valid_scans" db:"valid_scans"`
	InvalidScans   int        `json:"invalid_scans" db:"invalid_scans"`
	TotalRevenue   float64    `json:"total_revenue" db:"total_revenue"`
	SuccessRate    float64    `json:"success_rate" db:"success_rate"`
	LastActiveAt   *time.Time `json:"last_active_at,omitempty" db:"last_active_at"`
	EventsAssigned int        `json:"events_assigned" db:"-"`
}

// EventScanStats represents scanning statistics for an event
// db tags must match the SQL column aliases in GetEventScanStats query
type EventScanStats struct {
	EventID          uuid.UUID          `json:"event_id" db:"event_id"`
	TotalScanners    int                `json:"total_scanners" db:"total_scanners"`
	ActiveScanners   int                `json:"active_scanners" db:"active_scanners"`
	TotalScans       int                `json:"total_scans" db:"total_scans"`
	ValidScans       int                `json:"valid_scans" db:"valid_scans"`
	InvalidScans     int                `json:"invalid_scans" db:"invalid_scans"`
	TotalRevenue     float64            `json:"total_revenue" db:"total_revenue"`
	SuccessRate      float64            `json:"success_rate" db:"success_rate"`
	PeakScanTime     *string            `json:"peak_scan_time,omitempty" db:"-"`
	ScannerBreakdown []ScannerEventStats `json:"scanner_breakdown" db:"-"`
}

// ScannerEventStats represents statistics for a scanner within a specific event
// db tags must match the SQL column aliases in GetEventScanStats breakdown query
type ScannerEventStats struct {
	ScannerID    uuid.UUID `json:"scanner_id" db:"scanner_id"`
	ScannerName  string    `json:"scanner_name" db:"scanner_name"`
	TotalScans   int       `json:"total_scans" db:"total_scans"`
	ValidScans   int       `json:"valid_scans" db:"valid_scans"`
	InvalidScans int       `json:"invalid_scans" db:"invalid_scans"`
	Revenue      float64   `json:"revenue" db:"revenue"`
	SuccessRate  float64   `json:"success_rate" db:"success_rate"`
	LastScanAt   *time.Time `json:"last_scan_at,omitempty" db:"last_scan_at"`
}
