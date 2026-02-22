package payments

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending RefundStatus = "pending"
	RefundStatusSuccess RefundStatus = "success"
	RefundStatusFailed  RefundStatus = "failed"
)

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	InitializePayment(ctx context.Context, req *InitializePaymentRequest) (*InitializePaymentResponse, error)
	VerifyPayment(ctx context.Context, reference string) (*VerifyPaymentResponse, error)
	ProcessWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)
	RefundPayment(ctx context.Context, req *RefundPaymentRequest) (*RefundPaymentResponse, error)
	GetSupportedCurrencies() []string
	GetSupportedPaymentMethods() []string
}

// InitializePaymentRequest represents a payment initialization request
type InitializePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	EventID       uuid.UUID `json:"event_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Reference     string    `json:"reference"`
	CallbackURL   string    `json:"callback_url"`
	Description   string    `json:"description"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// InitializePaymentResponse represents a payment initialization response
type InitializePaymentResponse struct {
	PaymentID        string    `json:"payment_id"`
	PaymentReference string    `json:"payment_reference"`
	PaymentURL       string    `json:"payment_url"`
	AccessCode       string    `json:"access_code,omitempty"`
	ExpiresAt        time.Time `json:"expires_at"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// VerifyPaymentRequest represents a payment verification request
type VerifyPaymentRequest struct {
	PaymentReference string `json:"payment_reference"`
	Provider         string `json:"provider"`
}

// VerifyPaymentResponse represents a payment verification response
type VerifyPaymentResponse struct {
	PaymentID        string                 `json:"payment_id"`
	PaymentReference string                 `json:"payment_reference"`
	Status           PaymentStatus          `json:"status"`
	Amount           float64                `json:"amount"`
	Currency         string                 `json:"currency"`
	PaidAt           *time.Time             `json:"paid_at"`
	GatewayResponse  string                 `json:"gateway_response"`
	TransactionID    string                 `json:"transaction_id"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RefundPaymentRequest represents a payment refund request
type RefundPaymentRequest struct {
	PaymentReference string  `json:"payment_reference"`
	Amount           float64 `json:"amount"`
	Reason           string  `json:"reason,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RefundPaymentResponse represents a payment refund response
type RefundPaymentResponse struct {
	RefundID         string        `json:"refund_id"`
	PaymentReference string        `json:"payment_reference"`
	Status           RefundStatus  `json:"status"`
	Amount           float64       `json:"amount"`
	Currency         string        `json:"currency"`
	ProcessedAt      *time.Time    `json:"processed_at"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// WebhookEvent represents a webhook event from a payment provider
type WebhookEvent struct {
	Event            string                 `json:"event"`
	PaymentReference string                 `json:"payment_reference"`
	Status           PaymentStatus          `json:"status"`
	Data             map[string]interface{} `json:"data"`
	Provider         string                 `json:"provider"`
	Timestamp        time.Time              `json:"timestamp"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Provider    string   `json:"provider"`
	Currencies  []string `json:"currencies"`
	Countries   []string `json:"countries"`
	IsActive    bool     `json:"is_active"`
	Description string   `json:"description"`
	IconURL     string   `json:"icon_url,omitempty"`
}

// PaymentConfig represents payment configuration
type PaymentConfig struct {
	Provider    string            `json:"provider"`
	Environment string            `json:"environment"`
	Config      map[string]string `json:"config"`
	IsActive    bool              `json:"is_active"`
}

// PaymentFee represents payment processing fees
type PaymentFee struct {
	Provider     string  `json:"provider"`
	PaymentMethod string `json:"payment_method"`
	FixedFee     float64 `json:"fixed_fee"`
	PercentageFee float64 `json:"percentage_fee"`
	Currency     string  `json:"currency"`
}

// PaymentAnalytics represents payment analytics data
type PaymentAnalytics struct {
	TotalPayments    int64   `json:"total_payments"`
	SuccessfulPayments int64 `json:"successful_payments"`
	FailedPayments   int64   `json:"failed_payments"`
	TotalAmount      float64 `json:"total_amount"`
	AverageAmount    float64 `json:"average_amount"`
	Currency         string  `json:"currency"`
	Period           string  `json:"period"`
	ByProvider       map[string]PaymentProviderStats `json:"by_provider"`
	ByMethod         map[string]PaymentMethodStats   `json:"by_method"`
}

// PaymentProviderStats represents payment statistics by provider
type PaymentProviderStats struct {
	Provider         string  `json:"provider"`
	TotalPayments    int64   `json:"total_payments"`
	SuccessfulPayments int64 `json:"successful_payments"`
	FailedPayments   int64   `json:"failed_payments"`
	TotalAmount      float64 `json:"total_amount"`
	SuccessRate      float64 `json:"success_rate"`
}

// PaymentMethodStats represents payment statistics by method
type PaymentMethodStats struct {
	Method           string  `json:"method"`
	TotalPayments    int64   `json:"total_payments"`
	SuccessfulPayments int64 `json:"successful_payments"`
	FailedPayments   int64   `json:"failed_payments"`
	TotalAmount      float64 `json:"total_amount"`
	SuccessRate      float64 `json:"success_rate"`
}

// PaymentError represents a payment error
type PaymentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *PaymentError) Error() string {
	return e.Message
}

// NewPaymentError creates a new payment error
func NewPaymentError(code, message, details string) *PaymentError {
	return &PaymentError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common payment error codes
const (
	ErrCodeInvalidAmount        = "INVALID_AMOUNT"
	ErrCodeInvalidCurrency      = "INVALID_CURRENCY"
	ErrCodeInvalidPaymentMethod = "INVALID_PAYMENT_METHOD"
	ErrCodePaymentFailed        = "PAYMENT_FAILED"
	ErrCodePaymentExpired       = "PAYMENT_EXPIRED"
	ErrCodePaymentCancelled     = "PAYMENT_CANCELLED"
	ErrCodeInsufficientFunds    = "INSUFFICIENT_FUNDS"
	ErrCodeProviderError        = "PROVIDER_ERROR"
	ErrCodeNetworkError         = "NETWORK_ERROR"
	ErrCodeInvalidSignature     = "INVALID_SIGNATURE"
	ErrCodeWebhookError         = "WEBHOOK_ERROR"
	ErrCodeRefundFailed         = "REFUND_FAILED"
	ErrCodeRefundNotSupported   = "REFUND_NOT_SUPPORTED"
)

