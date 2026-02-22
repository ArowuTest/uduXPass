package payments

import (
	"time"
	
	"github.com/uduxpass/backend/internal/domain/entities"
)

// PaystackResponse represents the response from Paystack API
type PaystackResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// PaystackVerifyResponse represents Paystack payment verification response
type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64     `json:"id"`
		Domain          string    `json:"domain"`
		Status          string    `json:"status"`
		Reference       string    `json:"reference"`
		Amount          int64     `json:"amount"`
		Message         *string   `json:"message"`
		GatewayResponse string    `json:"gateway_response"`
		PaidAt          *string   `json:"paid_at"`
		CreatedAt       string    `json:"created_at"`
		Channel         string    `json:"channel"`
		Currency        string    `json:"currency"`
		Fees            int64     `json:"fees"`
		Customer  struct {
			ID           int64  `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
			Phone        string `json:"phone"`
		} `json:"customer"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			Brand             string `json:"brand"`
		} `json:"authorization"`
	} `json:"data"`
}

// PaymentResponse represents a unified payment response
type PaymentResponse struct {
	ID            string                 `json:"id"`
	Reference     string                 `json:"reference"`
	Status        PaymentStatus          `json:"status"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Method        entities.PaymentMethod `json:"method"`
	GatewayRef    string                 `json:"gateway_ref,omitempty"`
	TransactionID string                 `json:"transaction_id,omitempty"` // Provider-specific transaction ID
	PaymentURL    string                 `json:"payment_url,omitempty"`    // URL for payment completion
	Message       string                 `json:"message,omitempty"`
	PaidAt        *time.Time             `json:"paid_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
