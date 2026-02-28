package entities

import (
	"time"

	"github.com/google/uuid"
)

// ScannerRole represents the role types for scanner users
type ScannerRole string

const (
	ScannerRoleOperator    ScannerRole = "scanner_operator"
	ScannerRoleLead        ScannerRole = "lead_scanner"
	ScannerRoleSupervisor  ScannerRole = "scanner_supervisor"
)

// ScannerStatus represents scanner account status
type ScannerStatus string

const (
	ScannerStatusActive    ScannerStatus = "active"
	ScannerStatusInactive  ScannerStatus = "inactive"
	ScannerStatusLocked    ScannerStatus = "locked"
	ScannerStatusSuspended ScannerStatus = "suspended"
)

// ScannerPermission represents specific permissions for scanner users
type ScannerPermission string

const (
	ScannerPermissionScanTickets      ScannerPermission = "scan_tickets"
	ScannerPermissionManualEntry      ScannerPermission = "manual_entry"
	ScannerPermissionEmergencyOverride ScannerPermission = "emergency_override"
	ScannerPermissionBulkScan         ScannerPermission = "bulk_scan"
	ScannerPermissionViewReports      ScannerPermission = "view_reports"
	ScannerPermissionManageSettings   ScannerPermission = "manage_settings"
)

// ScannerUser represents a scanner operator user in the system
type ScannerUser struct {
	ID                  uuid.UUID     `json:"id" db:"id"`
	Username            string        `json:"username" db:"username"`
	PasswordHash        string        `json:"-" db:"password_hash"`
	Name                string        `json:"name" db:"name"`
	Email               string        `json:"email" db:"email"`
	Role                ScannerRole   `json:"role" db:"role"`
	Permissions         []string      `json:"permissions" db:"permissions"`
	Status              ScannerStatus `json:"status" db:"status"`
	LastLogin           *time.Time    `json:"last_login,omitempty" db:"last_login"`
	LoginAttempts       int           `json:"-" db:"login_attempts"`
	LockedUntil         *time.Time    `json:"-" db:"locked_until"`
	MustChangePassword  bool          `json:"must_change_password" db:"must_change_password"`
	CreatedBy           *uuid.UUID    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt           time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at" db:"updated_at"`
}

// ScannerEventAssignment represents the assignment of a scanner to an event
type ScannerEventAssignment struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ScannerID  uuid.UUID `json:"scanner_id" db:"scanner_id"`
	EventID    uuid.UUID `json:"event_id" db:"event_id"`
	AssignedBy uuid.UUID `json:"assigned_by" db:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
	IsActive   bool      `json:"is_active" db:"is_active"`
}

// ScannerSession represents an active scanning session
type ScannerSession struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ScannerID    uuid.UUID  `json:"scanner_id" db:"scanner_id"`
	EventID      uuid.UUID  `json:"event_id" db:"event_id"`
	StartTime    time.Time  `json:"start_time" db:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty" db:"end_time"`
	ScansCount   int        `json:"scans_count" db:"scans_count"`
	ValidScans   int        `json:"valid_scans" db:"valid_scans"`
	InvalidScans int        `json:"invalid_scans" db:"invalid_scans"`
	TotalRevenue float64    `json:"total_revenue" db:"total_revenue"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	Notes        *string    `json:"notes,omitempty" db:"notes"`
}

// ScannerAuditLog represents audit trail for scanner actions
type ScannerAuditLog struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	ScannerID    uuid.UUID              `json:"scanner_id" db:"scanner_id"`
	SessionID    *uuid.UUID             `json:"session_id,omitempty" db:"session_id"`
	Action       string                 `json:"action" db:"action"`
	ResourceType *string                `json:"resource_type,omitempty" db:"resource_type"`
	ResourceID   *uuid.UUID             `json:"resource_id,omitempty" db:"resource_id"`
	Details      map[string]interface{} `json:"details" db:"details"`
	IPAddress    *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string                `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// ScannerLoginHistory represents login history for scanner users
type ScannerLoginHistory struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ScannerID uuid.UUID  `json:"scanner_id" db:"scanner_id"`
	IPAddress *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string    `json:"user_agent,omitempty" db:"user_agent"`
	Success   bool       `json:"success" db:"success"`
	LoginAt   time.Time  `json:"login_at" db:"login_at"`
	LogoutAt  *time.Time `json:"logout_at,omitempty" db:"logout_at"`
}

// TicketValidation represents a ticket validation attempt
type TicketValidation struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	TicketID             uuid.UUID `json:"ticket_id" db:"ticket_id"`
	ScannerID            uuid.UUID `json:"scanner_id" db:"scanner_id"`
	SessionID            uuid.UUID `json:"session_id" db:"session_id"`
	ValidationResult     string    `json:"validation_result" db:"validation_result"`
	ValidationTimestamp  time.Time `json:"validation_timestamp" db:"validation_timestamp"`
	Notes                *string   `json:"notes,omitempty" db:"notes"`
}

// HasPermission checks if the scanner user has a specific permission
func (su *ScannerUser) HasPermission(permission ScannerPermission) bool {
	for _, p := range su.Permissions {
		if ScannerPermission(p) == permission {
			return true
		}
	}
	return false
}

// IsLocked checks if the scanner account is currently locked
func (su *ScannerUser) IsLocked() bool {
	if su.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*su.LockedUntil)
}

// GetDefaultPermissions returns default permissions for a scanner role
func GetDefaultScannerPermissions(role ScannerRole) []ScannerPermission {
	switch role {
	case ScannerRoleOperator:
		return []ScannerPermission{
			ScannerPermissionScanTickets,
			ScannerPermissionManualEntry,
		}
	case ScannerRoleLead:
		return []ScannerPermission{
			ScannerPermissionScanTickets,
			ScannerPermissionManualEntry,
			ScannerPermissionBulkScan,
			ScannerPermissionViewReports,
		}
	case ScannerRoleSupervisor:
		return []ScannerPermission{
			ScannerPermissionScanTickets,
			ScannerPermissionManualEntry,
			ScannerPermissionEmergencyOverride,
			ScannerPermissionBulkScan,
			ScannerPermissionViewReports,
			ScannerPermissionManageSettings,
		}
	default:
		return []ScannerPermission{ScannerPermissionScanTickets}
	}
}

// Request/Response DTOs

// ScannerLoginRequest represents a scanner login request
type ScannerLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ScannerLoginResponse represents a scanner login response
type ScannerLoginResponse struct {
	Success      bool         `json:"success"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	Scanner      *ScannerUser `json:"scanner,omitempty"`
	ExpiresIn    int64        `json:"expires_in,omitempty"`
	Message      string       `json:"message,omitempty"`
}

// ScannerCreateRequest represents a request to create a new scanner
type ScannerCreateRequest struct {
	Username    string              `json:"username" binding:"required"`
	Password    string              `json:"password" binding:"required,min=8"`
	Name        string              `json:"name" binding:"required"`
	Email       string              `json:"email" binding:"required,email"`
	Role        ScannerRole         `json:"role"`
	Permissions []ScannerPermission `json:"permissions"`
	EventIDs    []uuid.UUID         `json:"event_ids"`
}

// ScannerUpdateRequest represents a request to update a scanner
type ScannerUpdateRequest struct {
	Name        string              `json:"name"`
	Email       string              `json:"email"`
	Role        ScannerRole         `json:"role"`
	Permissions []ScannerPermission `json:"permissions"`
	EventIDs    []uuid.UUID         `json:"event_ids"`
	Status      *ScannerStatus      `json:"status"`
}

// ScannerAssignedEvent represents an event assigned to a scanner
type ScannerAssignedEvent struct {
	ScannerID   uuid.UUID   `json:"scanner_id" db:"scanner_id"`
	EventID     uuid.UUID   `json:"event_id" db:"event_id"`
	AssignedBy  uuid.UUID   `json:"assigned_by" db:"assigned_by"`
	AssignedAt  time.Time   `json:"assigned_at" db:"assigned_at"`
	EventName   string      `json:"event_name" db:"event_name"`
	EventDate   time.Time   `json:"event_date" db:"event_date"`
	VenueName   string      `json:"venue_name" db:"venue_name"`
	VenueCity   string      `json:"venue_city" db:"venue_city"`
	Status      EventStatus `json:"status" db:"event_status"`
}

// TicketValidationRequest represents a ticket validation request
type TicketValidationRequest struct {
	TicketCode string  `json:"ticket_code" binding:"required"`
	EventID    string  `json:"event_id" binding:"required"`
	Notes      *string `json:"notes,omitempty"`
}

// TicketValidationResponse represents a ticket validation response
type TicketValidationResponse struct {
	Success          bool      `json:"success"`
	Valid            bool      `json:"valid"`
	Message          string    `json:"message"`
	TicketID         *string   `json:"ticket_id,omitempty"`
	SerialNumber     *string   `json:"serial_number,omitempty"`
	TicketType       *string   `json:"ticket_type,omitempty"`
	HolderName       *string   `json:"holder_name,omitempty"`
	ValidationTime   time.Time `json:"validation_time"`
	AlreadyValidated bool      `json:"already_validated"`
}

// ScannerSessionStartRequest represents a request to start a scanning session
type ScannerSessionStartRequest struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
}

// ScannerSessionResponse represents a scanner session response
type ScannerSessionResponse struct {
	Success bool             `json:"success"`
	Session *ScannerSession  `json:"session,omitempty"`
	Message string           `json:"message,omitempty"`
}

