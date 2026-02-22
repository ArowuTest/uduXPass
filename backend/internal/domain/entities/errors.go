package entities

import "errors"

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	Message  string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Details map[string]interface{}
}

func (e *ValidationError) Error() string {
	return e.Message
}

// BusinessRuleError represents a business rule violation
type BusinessRuleError struct {
	Rule    string
	Message string
	Details map[string]interface{}
}

func (e *BusinessRuleError) Error() string {
	return e.Message
}

// ConflictError represents a resource conflict error
type ConflictError struct {
	Resource string
	Message  string
	Details  map[string]interface{}
}

func (e *ConflictError) Error() string {
	return e.Message
}

// Common domain errors
var (
	// General errors
	ErrConflictError        = errors.New("resource conflict")
	ErrNotFoundError        = errors.New("resource not found")
	
	// User errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserEmailExists      = errors.New("user with this email already exists")
	ErrUserPhoneExists      = errors.New("user with this phone number already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAccountLocked        = errors.New("account is locked")
	ErrAccountDeactivated   = errors.New("account is deactivated")
	ErrAccountNotVerified   = errors.New("account is not verified")
	ErrAccountSuspended     = errors.New("account is suspended")

	// Admin errors
	ErrAdminNotFound        = errors.New("admin not found")
	ErrAdminAlreadyExists   = errors.New("admin already exists")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInvalidRole          = errors.New("invalid role")

	// Event errors
	ErrEventNotFound        = errors.New("event not found")
	ErrEventAlreadyExists   = errors.New("event already exists")
	ErrEventNotActive       = errors.New("event is not active")
	ErrEventSoldOut         = errors.New("event is sold out")
	ErrEventCancelled       = errors.New("event is cancelled")
	ErrEventExpired         = errors.New("event has expired")

	// Ticket errors
	ErrTicketNotFound       = errors.New("ticket not found")
	ErrTicketAlreadyUsed    = errors.New("ticket already used")
	ErrTicketAlreadyRedeemed = errors.New("ticket already redeemed")
	ErrTicketExpired        = errors.New("ticket has expired")
	ErrTicketCancelled      = errors.New("ticket is cancelled")
	ErrTicketInvalid        = errors.New("ticket is invalid")
	ErrInsufficientTickets  = errors.New("insufficient tickets available")

	// Order errors
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderAlreadyExists   = errors.New("order already exists")
	ErrOrderCodeExists      = errors.New("order code already exists")
	ErrOrderExpired         = errors.New("order has expired")
	ErrOrderCancelled       = errors.New("order is cancelled")
	ErrOrderCompleted       = errors.New("order is already completed")
	ErrInvalidOrderStatus   = errors.New("invalid order status")

	// Payment errors
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentFailed        = errors.New("payment failed")
	ErrPaymentCancelled     = errors.New("payment cancelled")
	ErrPaymentExpired       = errors.New("payment expired")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPaymentProcessingError = errors.New("payment processing error")

	// Organizer errors
	ErrOrganizerNotFound    = errors.New("organizer not found")
	ErrOrganizerAlreadyExists = errors.New("organizer already exists")
	ErrOrganizerNotActive   = errors.New("organizer is not active")

	// Tour errors
	ErrTourNotFound         = errors.New("tour not found")
	ErrTourAlreadyExists    = errors.New("tour already exists")
	ErrTourNotActive        = errors.New("tour is not active")

	// Ticket Tier errors
	ErrTicketTierNotFound   = errors.New("ticket tier not found")
	ErrTicketTierSoldOut    = errors.New("ticket tier is sold out")
	ErrTicketTierNotActive  = errors.New("ticket tier is not active")

	// OTP errors
	ErrOTPNotFound          = errors.New("OTP not found")
	ErrOTPExpired           = errors.New("OTP has expired")
	ErrOTPAlreadyUsed       = errors.New("OTP already used")
	ErrInvalidOTP           = errors.New("invalid OTP")
	ErrOTPTokenNotFound     = errors.New("OTP token not found")
	ErrInvalidToken         = errors.New("invalid token")

	// Scanner errors
	ErrScannerNotFound      = errors.New("scanner not found")
	ErrScannerNotActive     = errors.New("scanner is not active")
	
	// Inventory errors
	ErrInsufficientInventory    = errors.New("insufficient inventory")
	ErrInventoryHoldNotFound    = errors.New("inventory hold not found")
	ErrInventoryHoldExpired     = errors.New("inventory hold expired")

	// General validation errors
	ErrValidationError       = errors.New("validation error")
	ErrInvalidInput          = errors.New("invalid input")
)

// Custom error types for compatibility
var InsufficientInventoryError = ErrInsufficientInventory

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
		Details: nil,
	}
}

// NewValidationErrorWithDetails creates a new validation error with details
func NewValidationErrorWithDetails(field, message string, details map[string]interface{}) error {
	return &ValidationError{
		Field:   field,
		Message: message,
		Details: details,
	}
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(rule, message string, details map[string]interface{}) error {
	return &BusinessRuleError{
		Rule:    rule,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, message string) error {
	return &NotFoundError{
		Resource: resource,
		Message:  message,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(resource, message string, details map[string]interface{}) error {
	return &ConflictError{
		Resource: resource,
		Message:  message,
		Details:  details,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsBusinessRuleError checks if an error is a business rule error
func IsBusinessRuleError(err error) bool {
	_, ok := err.(*BusinessRuleError)
	return ok
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

