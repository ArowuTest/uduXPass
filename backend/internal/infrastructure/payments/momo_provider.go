package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uduxpass/backend/internal/domain/entities"
 )

// MoMoProvider implements payment processing for Mobile Money
type MoMoProvider struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

// NewMoMoProvider creates a new Mobile Money payment provider
func NewMoMoProvider(apiKey, apiSecret string ) *MoMoProvider {
	return &MoMoProvider{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.momo.com/v1", // Replace with actual MoMo API URL
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// InitializePayment initializes a payment with Mobile Money
func (m *MoMoProvider ) InitializePayment(order *entities.Order) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/payments/initialize", m.baseURL)
	
	payload := map[string]interface{}{
		"amount":      order.TotalAmount,
		"currency":    "NGN",
		"reference":   order.ID.String(),
		"phone":       order.CustomerPhone,
		"description": fmt.Sprintf("Payment for order %s", order.ID.String()),
		"metadata": map[string]interface{}{
			"order_id": order.ID.String(),
		},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload ))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var momoResp map[string]interface{}
	if err := json.Unmarshal(body, &momoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	status, ok := momoResp["status"].(string)
	if !ok || status != "success" {
		message, _ := momoResp["message"].(string)
		return nil, fmt.Errorf("momo error: %s", message)
	}
	
	transactionID, ok := momoResp["transaction_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing transaction_id in response")
	}
	
	return &PaymentResponse{
		ID:        transactionID,
		Reference: order.ID.String(),
		Status:    PaymentStatusPending,
		Amount:    order.TotalAmount,
		Currency:  "NGN",
		Method:    entities.PaymentMethodMoMo,
		Message:   "Payment initiated. Please complete on your mobile device.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"phone": order.CustomerPhone,
		},
	}, nil
}

// VerifyPayment verifies a payment with Mobile Money
func (m *MoMoProvider) VerifyPayment(reference string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/payments/verify/%s", m.baseURL, reference)
	
	req, err := http.NewRequest("GET", url, nil )
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var verifyResp map[string]interface{}
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	// Convert status
	var status PaymentStatus
	statusStr, _ := verifyResp["status"].(string)
	switch statusStr {
	case "success", "completed":
		status = PaymentStatusSuccess
	case "failed":
		status = PaymentStatusFailed
	case "cancelled":
		status = PaymentStatusCancelled
	default:
		status = PaymentStatusPending
	}
	
	// Parse paid_at timestamp
	var paidAt *time.Time
	if paidAtStr, ok := verifyResp["paid_at"].(string); ok && paidAtStr != "" {
		if t, err := time.Parse(time.RFC3339, paidAtStr); err == nil {
			paidAt = &t
		}
	}
	
	amount, _ := verifyResp["amount"].(float64)
	currency, _ := verifyResp["currency"].(string)
	if currency == "" {
		currency = "NGN"
	}
	
	return &PaymentResponse{
		ID:        reference,
		Reference: reference,
		Status:    status,
		Amount:    amount,
		Currency:  currency,
		Method:    entities.PaymentMethodMoMo,
		PaidAt:    paidAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ProcessWebhook processes Mobile Money webhook events
func (m *MoMoProvider) ProcessWebhook(payload []byte) (*PaymentResponse, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}
	
	reference, ok := event["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("missing reference in webhook")
	}
	
	// Verify the payment to ensure authenticity
	return m.VerifyPayment(reference)
}

// RequestToPay initiates a payment request
func (m *MoMoProvider) RequestToPay(ctx context.Context, request MoMoPaymentRequest) (*PaymentResponse, error) {
	// Stub implementation for testing
	return &PaymentResponse{
		Status:        "pending",
		Reference:     "momo_" + time.Now().Format("20060102150405"),
		TransactionID: "txn_" + time.Now().Format("20060102150405"),
		PaymentURL:    "https://momo.test/pay/123",
		Message:       "Payment initiated successfully",
	}, nil
}

// GetTransactionStatus gets the status of a transaction
func (m *MoMoProvider) GetTransactionStatus(ctx context.Context, reference string) (*PaymentResponse, error) {
	// Stub implementation for testing
	return &PaymentResponse{
		Status:        "completed",
		Reference:     reference,
		TransactionID: "txn_" + reference,
		Message:       "Payment completed successfully",
	}, nil
}

// MoMoPaymentRequest represents a Mobile Money payment request (local definition)
type MoMoPaymentRequest struct {
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	ExternalID   string  `json:"external_id"`
	Phone        string  `json:"phone"`
	PayerMessage string  `json:"payer_message"`
	PayeeNote    string  `json:"payee_note"`
}

