package entities

import (
	"time"

	"github.com/google/uuid"
)

// OTPPurpose represents the purpose of an OTP token
type OTPPurpose string

const (
	OTPPurposeLogin            OTPPurpose = "login"
	OTPPurposeRegistration     OTPPurpose = "registration"
	OTPPurposePasswordReset    OTPPurpose = "password_reset"
	OTPPurposePhoneVerification OTPPurpose = "phone_verification"
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
	OTPPurposePaymentVerification OTPPurpose = "payment_verification"
)

// OTPStatus represents the status of an OTP token
type OTPStatus string

const (
	OTPStatusPending   OTPStatus = "pending"
	OTPStatusUsed      OTPStatus = "used"
	OTPStatusExpired   OTPStatus = "expired"
	OTPStatusCancelled OTPStatus = "cancelled"
)

// OTPToken represents an OTP token entity
type OTPToken struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Phone       string     `json:"phone" db:"phone"`
	Email       *string    `json:"email,omitempty" db:"email"`
	Code        string     `json:"code" db:"code"`
	Purpose     OTPPurpose `json:"purpose" db:"purpose"`
	Status      OTPStatus  `json:"status" db:"status"`
	ExpiresAt   time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt      *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	AttemptCount int       `json:"attempt_count" db:"attempt_count"`
	MaxAttempts  int       `json:"max_attempts" db:"max_attempts"`
	UserAgent   *string    `json:"user_agent,omitempty" db:"user_agent"`
	IPAddress   *string    `json:"ip_address,omitempty" db:"ip_address"`
}

// NewOTPToken creates a new OTP token
func NewOTPToken(phone string, code string, purpose OTPPurpose, expiryDuration time.Duration) *OTPToken {
	now := time.Now()
	return &OTPToken{
		ID:          uuid.New(),
		Phone:       phone,
		Code:        code,
		Purpose:     purpose,
		Status:      OTPStatusPending,
		ExpiresAt:   now.Add(expiryDuration),
		CreatedAt:   now,
		UpdatedAt:   now,
		AttemptCount: 0,
		MaxAttempts:  3, // Default max attempts
	}
}

// NewOTPTokenWithEmail creates a new OTP token with email
func NewOTPTokenWithEmail(phone string, email string, code string, purpose OTPPurpose, expiryDuration time.Duration) *OTPToken {
	token := NewOTPToken(phone, code, purpose, expiryDuration)
	token.Email = &email
	return token
}

// IsExpired checks if the OTP token is expired
func (o *OTPToken) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsUsed checks if the OTP token has been used
func (o *OTPToken) IsUsed() bool {
	return o.Status == OTPStatusUsed
}

// IsValid checks if the OTP token is valid (not expired, not used, not cancelled)
func (o *OTPToken) IsValid() bool {
	return o.Status == OTPStatusPending && !o.IsExpired()
}

// CanAttempt checks if more attempts are allowed
func (o *OTPToken) CanAttempt() bool {
	return o.AttemptCount < o.MaxAttempts
}

// IncrementAttempt increments the attempt count
func (o *OTPToken) IncrementAttempt() {
	o.AttemptCount++
	o.UpdatedAt = time.Now()
}

// MarkAsUsed marks the OTP token as used
func (o *OTPToken) MarkAsUsed() {
	now := time.Now()
	o.Status = OTPStatusUsed
	o.UsedAt = &now
	o.UpdatedAt = now
}

// MarkAsExpired marks the OTP token as expired
func (o *OTPToken) MarkAsExpired() {
	o.Status = OTPStatusExpired
	o.UpdatedAt = time.Now()
}

// MarkAsCancelled marks the OTP token as cancelled
func (o *OTPToken) MarkAsCancelled() {
	o.Status = OTPStatusCancelled
	o.UpdatedAt = time.Now()
}

// SetUserAgent sets the user agent
func (o *OTPToken) SetUserAgent(userAgent string) {
	o.UserAgent = &userAgent
	o.UpdatedAt = time.Now()
}

// SetIPAddress sets the IP address
func (o *OTPToken) SetIPAddress(ipAddress string) {
	o.IPAddress = &ipAddress
	o.UpdatedAt = time.Now()
}

// Validate validates the OTP token fields
func (o *OTPToken) Validate() error {
	if o.Phone == "" {
		return ErrInvalidInput
	}

	if o.Code == "" {
		return ErrInvalidInput
	}

	if o.Purpose == "" {
		return ErrInvalidInput
	}

	if o.ExpiresAt.IsZero() {
		return ErrInvalidInput
	}

	if o.MaxAttempts <= 0 {
		return ErrInvalidInput
	}

	return nil
}

// GetRemainingAttempts returns the number of remaining attempts
func (o *OTPToken) GetRemainingAttempts() int {
	remaining := o.MaxAttempts - o.AttemptCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetTimeUntilExpiry returns the duration until expiry
func (o *OTPToken) GetTimeUntilExpiry() time.Duration {
	if o.IsExpired() {
		return 0
	}
	return time.Until(o.ExpiresAt)
}

