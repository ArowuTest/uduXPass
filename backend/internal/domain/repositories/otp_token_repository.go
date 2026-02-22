package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

// OTPTokenRepository defines the interface for OTP token data operations
type OTPTokenRepository interface {
	// Create creates a new OTP token
	Create(ctx context.Context, token *entities.OTPToken) error

	// GetByID retrieves an OTP token by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OTPToken, error)

	// GetByPhone retrieves OTP tokens by phone number
	GetByPhone(ctx context.Context, phone string) ([]*entities.OTPToken, error)

	// GetActiveByPhone retrieves active (pending) OTP tokens by phone number
	GetActiveByPhone(ctx context.Context, phone string) ([]*entities.OTPToken, error)

	// GetByPhoneAndPurpose retrieves OTP tokens by phone and purpose
	GetByPhoneAndPurpose(ctx context.Context, phone string, purpose entities.OTPPurpose) ([]*entities.OTPToken, error)

	// GetActiveByPhoneAndPurpose retrieves active OTP tokens by phone and purpose
	GetActiveByPhoneAndPurpose(ctx context.Context, phone string, purpose entities.OTPPurpose) (*entities.OTPToken, error)

	// GetByCode retrieves an OTP token by code
	GetByCode(ctx context.Context, code string) (*entities.OTPToken, error)

	// Update updates an existing OTP token
	Update(ctx context.Context, token *entities.OTPToken) error

	// MarkUsed marks an OTP token as used
	MarkUsed(ctx context.Context, id uuid.UUID) error

	// MarkExpired marks an OTP token as expired
	MarkExpired(ctx context.Context, id uuid.UUID) error

	// MarkCancelled marks an OTP token as cancelled
	MarkCancelled(ctx context.Context, id uuid.UUID) error

	// IncrementAttempt increments the attempt count for an OTP token
	IncrementAttempt(ctx context.Context, id uuid.UUID) error

	// Delete deletes an OTP token
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByPhone deletes all OTP tokens for a phone number
	DeleteByPhone(ctx context.Context, phone string) error

	// DeleteExpired deletes all expired OTP tokens
	DeleteExpired(ctx context.Context) error

	// DeleteOlderThan deletes OTP tokens older than the specified duration
	DeleteOlderThan(ctx context.Context, duration time.Duration) error

	// GetExpiredTokens retrieves all expired tokens
	GetExpiredTokens(ctx context.Context) ([]*entities.OTPToken, error)

	// GetTokensByStatus retrieves tokens by status
	GetTokensByStatus(ctx context.Context, status entities.OTPStatus) ([]*entities.OTPToken, error)

	// GetTokensByPurpose retrieves tokens by purpose
	GetTokensByPurpose(ctx context.Context, purpose entities.OTPPurpose) ([]*entities.OTPToken, error)

	// CountActiveTokensByPhone counts active tokens for a phone number
	CountActiveTokensByPhone(ctx context.Context, phone string) (int, error)

	// CountTokensByPhoneAndTimeRange counts tokens for a phone in a time range
	CountTokensByPhoneAndTimeRange(ctx context.Context, phone string, start, end time.Time) (int, error)

	// GetRecentTokensByPhone retrieves recent tokens for a phone number
	GetRecentTokensByPhone(ctx context.Context, phone string, limit int) ([]*entities.OTPToken, error)

	// ValidateAndMarkUsed validates an OTP code and marks it as used if valid
	ValidateAndMarkUsed(ctx context.Context, phone string, code string, purpose entities.OTPPurpose) (*entities.OTPToken, error)

	// CleanupExpiredTokens removes expired tokens (for background cleanup)
	CleanupExpiredTokens(ctx context.Context) (int, error)

	// GetTokenStats retrieves statistics about OTP tokens
	GetTokenStats(ctx context.Context) (*OTPTokenStats, error)
}

// OTPTokenStats represents statistics about OTP tokens
type OTPTokenStats struct {
	TotalTokens    int `json:"total_tokens"`
	PendingTokens  int `json:"pending_tokens"`
	UsedTokens     int `json:"used_tokens"`
	ExpiredTokens  int `json:"expired_tokens"`
	CancelledTokens int `json:"cancelled_tokens"`
}

// NewOTPToken creates a new OTP token (helper function)
func NewOTPToken(phone string, code string, purpose entities.OTPPurpose, expiryDuration time.Duration) *entities.OTPToken {
	return entities.NewOTPToken(phone, code, purpose, expiryDuration)
}

// OTP Purpose constants for easy access
const (
	OTPPurposeLogin            = entities.OTPPurposeLogin
	OTPPurposeRegistration     = entities.OTPPurposeRegistration
	OTPPurposePasswordReset    = entities.OTPPurposePasswordReset
	OTPPurposePhoneVerification = entities.OTPPurposePhoneVerification
	OTPPurposeEmailVerification = entities.OTPPurposeEmailVerification
	OTPPurposePaymentVerification = entities.OTPPurposePaymentVerification
)

// OTPTokenFilter represents filtering options for OTP tokens
type OTPTokenFilter struct {
	Phone       *string                `json:"phone,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Identifier  *string                `json:"identifier,omitempty"`
	Purpose     *entities.OTPPurpose   `json:"purpose,omitempty"`
	Status      *entities.OTPStatus    `json:"status,omitempty"`
	IsUsed      *bool                  `json:"is_used,omitempty"`
	ExpiredOnly *bool                  `json:"expired_only,omitempty"`
	ActiveOnly  *bool                  `json:"active_only,omitempty"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	CreatedFrom *time.Time             `json:"created_from,omitempty"`
	CreatedTo   *time.Time             `json:"created_to,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	ExpiresFrom *time.Time             `json:"expires_from,omitempty"`
	ExpiresTo   *time.Time             `json:"expires_to,omitempty"`
	
	// Pagination and sorting
	Page      *int    `json:"page,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
	SortBy    *string `json:"sort_by,omitempty"`
	SortOrder *string `json:"sort_order,omitempty"`
}

// OTPValidationResult represents the result of OTP validation
type OTPValidationResult struct {
	Valid     bool                `json:"valid"`
	Token     *entities.OTPToken  `json:"token,omitempty"`
	Error     string              `json:"error,omitempty"`
	Attempts  int                 `json:"attempts"`
}

