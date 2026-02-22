package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/uduxpass/backend/internal/domain/entities"
 )

// PaystackProvider implements payment processing for Paystack
type PaystackProvider struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

// NewPaystackProvider creates a new Paystack payment provider
func NewPaystackProvider(secretKey string ) *PaystackProvider {
	return &PaystackProvider{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// InitializePayment initializes a payment with Paystack
func (p *PaystackProvider ) InitializePayment(order *entities.Order) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", p.baseURL)
	
	// Convert amount to kobo (Paystack uses kobo for NGN)
	amountInKobo := int64(order.TotalAmount * 100)
	
	payload := map[string]interface{}{
		"email":     order.CustomerEmail,
		"amount":    amountInKobo,
		"reference": order.ID.String(),
		"currency":  "NGN",
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
	
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var paystackResp PaystackResponse
	if err := json.Unmarshal(body, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if !paystackResp.Status {
		return nil, fmt.Errorf("paystack error: %s", paystackResp.Message)
	}
	
	// Extract authorization URL from response
	authURL, ok := paystackResp.Data["authorization_url"].(string)
	if !ok {
		return nil, fmt.Errorf("missing authorization_url in response")
	}
	
	accessCode, ok := paystackResp.Data["access_code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing access_code in response")
	}
	
	return &PaymentResponse{
		ID:        accessCode,
		Reference: order.ID.String(),
		Status:    PaymentStatusPending,
		Amount:    order.TotalAmount,
		Currency:  "NGN",
		Method:    entities.PaymentMethodCard,
		Message:   authURL, // Store authorization URL in message field
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"authorization_url": authURL,
			"access_code":      accessCode,
		},
	}, nil
}

// VerifyPayment verifies a payment with Paystack
func (p *PaystackProvider) VerifyPayment(reference string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", p.baseURL, reference)
	
	req, err := http.NewRequest("GET", url, nil )
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var verifyResp PaystackVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if !verifyResp.Status {
		return nil, fmt.Errorf("paystack error: %s", verifyResp.Message)
	}
	
	// Convert status
	var status PaymentStatus
	switch verifyResp.Data.Status {
	case "success":
		status = PaymentStatusSuccess
	case "failed":
		status = PaymentStatusFailed
	case "abandoned":
		status = PaymentStatusCancelled
	default:
		status = PaymentStatusPending
	}
	
	// Parse paid_at timestamp
	var paidAt *time.Time
	if verifyResp.Data.PaidAt != nil && *verifyResp.Data.PaidAt != "" {
		if t, err := time.Parse(time.RFC3339, *verifyResp.Data.PaidAt); err == nil {
			paidAt = &t
		}
	}
	
	// Convert amount from kobo to naira
	amount := float64(verifyResp.Data.Amount) / 100
	
	return &PaymentResponse{
		ID:        strconv.FormatInt(verifyResp.Data.ID, 10),
		Reference: verifyResp.Data.Reference,
		Status:    status,
		Amount:    amount,
		Currency:  verifyResp.Data.Currency,
		Method:    entities.PaymentMethodCard,
		GatewayRef: verifyResp.Data.GatewayResponse,
		PaidAt:    paidAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"channel": verifyResp.Data.Channel,
			"fees":    verifyResp.Data.Fees,
		},
	}, nil
}

// ProcessWebhook processes Paystack webhook events
func (p *PaystackProvider) ProcessWebhook(payload []byte) (*PaymentResponse, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}
	
	eventType, ok := event["event"].(string)
	if !ok {
		return nil, fmt.Errorf("missing event type in webhook")
	}
	
	if eventType != "charge.success" {
		return nil, fmt.Errorf("unsupported event type: %s", eventType)
	}
	
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing data in webhook")
	}
	
	reference, ok := data["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("missing reference in webhook data")
	}
	
	// Verify the payment to ensure authenticity
	return p.VerifyPayment(reference)
}

// InitializeTransaction initializes a transaction with Paystack
func (p *PaystackProvider) InitializeTransaction(ctx context.Context, request PaystackPaymentRequest) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", p.baseURL)
	
	// Convert amount to kobo (Paystack uses kobo for NGN)
	amountInKobo := int64(request.Amount * 100)
	
	payload := map[string]interface{}{
		"email":        request.Email,
		"amount":       amountInKobo,
		"reference":    request.Reference,
		"currency":     request.Currency,
		"callback_url": request.CallbackURL,
		"metadata":     request.Metadata,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var paystackResp PaystackResponse
	if err := json.Unmarshal(body, &paystackResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if !paystackResp.Status {
		return nil, fmt.Errorf("paystack error: %s", paystackResp.Message)
	}
	
	// Extract authorization URL from response
	authURL, _ := paystackResp.Data["authorization_url"].(string)
	accessCode, _ := paystackResp.Data["access_code"].(string)
	
	return &PaymentResponse{
		ID:            accessCode,
		Reference:     request.Reference,
		Status:        "pending",
		Amount:        request.Amount,
		Currency:      request.Currency,
		Method:        entities.PaymentMethodCard,
		PaymentURL:    authURL,
		TransactionID: accessCode,
		Message:       "Transaction initialized successfully",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}


// PaystackPaymentRequest represents a Paystack payment request (local definition)
type PaystackPaymentRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Email       string            `json:"email"`
	Reference   string            `json:"reference"`
	CallbackURL string            `json:"callback_url"`
	Metadata    map[string]string `json:"metadata"`
}

