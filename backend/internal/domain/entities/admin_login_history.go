package entities

import (
	"time"

	"github.com/google/uuid"
)

// AdminLoginHistory represents an admin login attempt record
type AdminLoginHistory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AdminID     uuid.UUID `json:"admin_id" db:"admin_id"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Success     bool      `json:"success" db:"success"`
	FailureReason string  `json:"failure_reason,omitempty" db:"failure_reason"`
	LoginTime   time.Time `json:"login_time" db:"login_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// NewAdminLoginHistory creates a new admin login history record
func NewAdminLoginHistory(adminID uuid.UUID, ipAddress, userAgent string, success bool, failureReason string) *AdminLoginHistory {
	return &AdminLoginHistory{
		ID:            uuid.New(),
		AdminID:       adminID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Success:       success,
		FailureReason: failureReason,
		LoginTime:     time.Now(),
		CreatedAt:     time.Now(),
	}
}

// Validate validates the admin login history
func (alh *AdminLoginHistory) Validate() error {
	if alh.AdminID == uuid.Nil {
		return NewValidationError("admin_id", "admin ID is required")
	}

	if alh.IPAddress == "" {
		return NewValidationError("ip_address", "IP address is required")
	}

	if alh.UserAgent == "" {
		return NewValidationError("user_agent", "user agent is required")
	}

	return nil
}

