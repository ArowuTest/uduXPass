package payments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/internal/domain/services"
	"github.com/uduxpass/backend/internal/infrastructure/payments"
	"github.com/uduxpass/backend/pkg/qrcode"
)

// PaymentService handles payment processing use cases
type PaymentService struct {
	paymentRepo       repositories.PaymentRepository
	orderRepo         repositories.OrderRepository
	orderLineRepo     repositories.OrderLineRepository
	ticketRepo        repositories.TicketRepository
	inventoryHoldRepo repositories.InventoryHoldRepository
	eventRepo         repositories.EventRepository
	momoProvider      payments.MoMoProvider
	paystackProvider  payments.PaystackProvider
	unitOfWork        repositories.UnitOfWork
	qrGenerator       *qrcode.Generator
	emailService      services.EmailService
	jwtSecret         []byte
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	orderLineRepo repositories.OrderLineRepository,
	ticketRepo repositories.TicketRepository,
	inventoryHoldRepo repositories.InventoryHoldRepository,
	eventRepo repositories.EventRepository,
	momoProvider payments.MoMoProvider,
	paystackProvider payments.PaystackProvider,
	unitOfWork repositories.UnitOfWork,
	emailService services.EmailService,
	jwtSecret string,
) *PaymentService {
	return &PaymentService{
		paymentRepo:       paymentRepo,
		orderRepo:         orderRepo,
		orderLineRepo:     orderLineRepo,
		ticketRepo:        ticketRepo,
		inventoryHoldRepo: inventoryHoldRepo,
		eventRepo:         eventRepo,
		momoProvider:      momoProvider,
		paystackProvider:  paystackProvider,
		unitOfWork:        unitOfWork,
		qrGenerator:       qrcode.NewGenerator(),
		emailService:      emailService,
		jwtSecret:         []byte(jwtSecret),
	}
}

// InitiatePaymentRequest represents the request to initiate payment
type InitiatePaymentRequest struct {
	OrderID       uuid.UUID                `json:"order_id" validate:"required"`
	PaymentMethod entities.PaymentMethod   `json:"payment_method" validate:"required"`
	CustomerInfo  PaymentCustomerInfo      `json:"customer_info" validate:"required"`
	CallbackURL   string                   `json:"callback_url,omitempty"`
}

// InitiatePaymentResponse represents the response from payment initiation
type InitiatePaymentResponse struct {
	PaymentID         uuid.UUID              `json:"payment_id"`
	PaymentMethod     entities.PaymentMethod `json:"payment_method"`
	Amount            float64                `json:"amount"`
	Currency          string                 `json:"currency"`
	Status            entities.PaymentStatus `json:"status"`
	AuthorizationURL  *string                `json:"authorization_url,omitempty"`
	PaymentReference  string                 `json:"payment_reference"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
	Instructions      *string                `json:"instructions,omitempty"`
}

// InitiatePayment initiates payment for an order
func (s *PaymentService) InitiatePayment(ctx context.Context, req *InitiatePaymentRequest) (*InitiatePaymentResponse, error) {
	// Get order
	fmt.Printf("[DEBUG] InitiatePayment: Looking for order ID: %s\n", req.OrderID.String())
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		fmt.Printf("[DEBUG] InitiatePayment: GetByID failed with error: %v\n", err)
		return nil, entities.NewNotFoundError("order", "order not found")
	}
	fmt.Printf("[DEBUG] InitiatePayment: Order found: %s (status: %s, active: %v)\n", order.Code, order.Status, order.IsActive)
	fmt.Printf("[DEBUG] InitiatePayment: ExpiresAt: %v, Now: %v, IsExpired: %v\n", order.ExpiresAt, time.Now(), order.IsExpired())

	// Verify order can be paid
	if !order.CanBePaid() {
		fmt.Printf("[DEBUG] InitiatePayment: CanBePaid returned false (Status: %s, IsExpired: %v)\n", order.Status, order.IsExpired())
		return nil, entities.NewBusinessRuleError("payment_not_allowed", "order cannot be paid (expired or already paid)", nil)
	}

	// Start transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create payment entity
	payment := entities.NewPayment(order.ID, req.PaymentMethod, order.TotalAmount, order.Currency)

	// Validate payment
	if err := payment.Validate(); err != nil {
		return nil, err
	}

	// Create payment record
	if err := tx.Payments().Create(tx.Context(), payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update order payment method
	order.SetPaymentMethod(req.PaymentMethod)
	if err := tx.Orders().Update(tx.Context(), order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Initiate payment with provider
	var response *InitiatePaymentResponse
	switch req.PaymentMethod {
	case entities.PaymentMethodMoMo:
		response, err = s.initiateMoMoPayment(ctx, payment, order, req.CustomerInfo)
	case entities.PaymentMethodPaystack:
		response, err = s.initiatePaystackPayment(ctx, payment, order, req.CustomerInfo, req.CallbackURL)
	default:
		return nil, entities.NewValidationError("payment_method", "unsupported payment method")
	}

	if err != nil {
		// Update payment status to failed
		payment.MarkFailed()
		s.paymentRepo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	return response, nil
}

// initiateMoMoPayment initiates MoMo payment
func (s *PaymentService) initiateMoMoPayment(ctx context.Context, payment *entities.Payment, order *entities.Order, customerInfo PaymentCustomerInfo) (*InitiatePaymentResponse, error) {
	// Convert to infrastructure type
	momoReq := payments.MoMoPaymentRequest{
		Amount:       payment.Amount,
		Currency:     payment.Currency,
		ExternalID:   payment.ID.String(),
		Phone:        customerInfo.Phone,
		PayerMessage: fmt.Sprintf("Payment for order %s", order.Code),
		PayeeNote:    fmt.Sprintf("uduXPass ticket payment - Order %s", order.Code),
	}

	momoResp, err := s.momoProvider.RequestToPay(ctx, momoReq)
	if err != nil {
		return nil, fmt.Errorf("MoMo payment initiation failed: %w", err)
	}

	// Update payment with provider transaction ID
	payment.SetProviderTransactionID(momoResp.TransactionID)
	payment.UpdateProviderResponse(map[string]interface{}{
		"transaction_id": momoResp.TransactionID,
		"status":         momoResp.Status,
		"message":        momoResp.Message,
	})

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	instructions := "Please approve the payment request on your mobile phone"
	expiresAt := time.Now().Add(5 * time.Minute) // MoMo payments expire in 5 minutes

	return &InitiatePaymentResponse{
		PaymentID:        payment.ID,
		PaymentMethod:    payment.Provider,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		Status:           payment.Status,
		PaymentReference: momoResp.TransactionID,
		ExpiresAt:        &expiresAt,
		Instructions:     &instructions,
	}, nil
}

// initiatePaystackPayment initiates Paystack payment
func (s *PaymentService) initiatePaystackPayment(ctx context.Context, payment *entities.Payment, order *entities.Order, customerInfo PaymentCustomerInfo, callbackURL string) (*InitiatePaymentResponse, error) {
	// Convert to infrastructure type
	paystackReq := payments.PaystackPaymentRequest{
		Amount:      payment.Amount * 100, // Paystack expects amount in kobo
		Currency:    payment.Currency,
		Email:       customerInfo.Email,
		Reference:   payment.ID.String(),
		CallbackURL: callbackURL,
		Metadata: map[string]string{
			"order_id":      order.ID.String(),
			"order_code":    order.Code,
			"customer_name": fmt.Sprintf("%s %s", customerInfo.FirstName, customerInfo.LastName),
		},
	}

	paystackResp, err := s.paystackProvider.InitializeTransaction(ctx, paystackReq)
	if err != nil {
		return nil, fmt.Errorf("Paystack payment initiation failed: %w", err)
	}

	// Update payment with provider transaction ID
	payment.SetProviderTransactionID(paystackResp.Reference)
	payment.UpdateProviderResponse(map[string]interface{}{
		"reference":         paystackResp.Reference,
		"authorization_url": paystackResp.PaymentURL,
		"access_code":       paystackResp.TransactionID,
	})

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return &InitiatePaymentResponse{
		PaymentID:        payment.ID,
		PaymentMethod:    payment.Provider,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		Status:           payment.Status,
		AuthorizationURL: &paystackResp.PaymentURL,
		PaymentReference: paystackResp.Reference,
	}, nil
}

// VerifyPaymentRequest represents the request to verify payment
type VerifyPaymentRequest struct {
	PaymentID uuid.UUID `json:"payment_id" validate:"required"`
}

// VerifyPaymentResponse represents the response from payment verification
type VerifyPaymentResponse struct {
	PaymentID        uuid.UUID              `json:"payment_id"`
	Status           entities.PaymentStatus `json:"status"`
	Amount           float64                `json:"amount"`
	Currency         string                 `json:"currency"`
	PaidAt           *time.Time             `json:"paid_at,omitempty"`
	OrderID          uuid.UUID              `json:"order_id"`
	OrderStatus      entities.OrderStatus   `json:"order_status"`
	TicketsGenerated bool                   `json:"tickets_generated"`
}

// VerifyPayment verifies payment status with the provider
func (s *PaymentService) VerifyPayment(ctx context.Context, req *VerifyPaymentRequest) (*VerifyPaymentResponse, error) {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, entities.NewNotFoundError("payment", "payment not found")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, entities.NewNotFoundError("order", "order not found")
	}

	// Verify with provider
	var verified bool
	var providerResponse map[string]interface{}

	switch payment.Provider {
	case entities.PaymentMethodMoMo:
		verified, providerResponse, err = s.verifyMoMoPayment(ctx, payment)
	case entities.PaymentMethodPaystack:
		verified, providerResponse, err = s.verifyPaystackPayment(ctx, payment)
	default:
		return nil, entities.NewValidationError("provider", "unsupported payment provider")
	}

	if err != nil {
		return nil, fmt.Errorf("payment verification failed: %w", err)
	}

	// Update payment with provider response
	payment.UpdateProviderResponse(providerResponse)

	ticketsGenerated := false

	if verified && payment.Status != entities.PaymentStatusCompleted {
		// Start transaction for payment completion
		tx, err := s.unitOfWork.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		// Mark payment as completed
		if err := payment.MarkCompleted(); err != nil {
			return nil, err
		}

		// Mark order as paid
		order.MarkPaid()

		// Update payment and order
		if err := tx.Payments().Update(tx.Context(), payment); err != nil {
			return nil, fmt.Errorf("failed to update payment: %w", err)
		}

		if err := tx.Orders().Update(tx.Context(), order); err != nil {
			return nil, fmt.Errorf("failed to update order: %w", err)
		}

		// Generate tickets
		if err := s.generateTickets(tx.Context(), tx, order); err != nil {
			return nil, fmt.Errorf("failed to generate tickets: %w", err)
		}
		ticketsGenerated = true

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	var paidAt *time.Time
	if payment.Status == entities.PaymentStatusCompleted {
		paidAt = &payment.UpdatedAt
	}

	return &VerifyPaymentResponse{
		PaymentID:        payment.ID,
		Status:           payment.Status,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		PaidAt:           paidAt,
		OrderID:          order.ID,
		OrderStatus:      order.Status,
		TicketsGenerated: ticketsGenerated,
	}, nil
}

// verifyMoMoPayment verifies MoMo payment status
func (s *PaymentService) verifyMoMoPayment(ctx context.Context, payment *entities.Payment) (bool, map[string]interface{}, error) {
	if payment.ProviderTransactionID == nil {
		return false, nil, fmt.Errorf("no provider transaction ID")
	}

	status, err := s.momoProvider.GetTransactionStatus(ctx, *payment.ProviderTransactionID)
	if err != nil {
		return false, nil, err
	}

	verified := status.Status == "SUCCESSFUL"
	response := map[string]interface{}{
		"status":         status.Status,
		"message":        status.Message,
		"transaction_id": status.TransactionID,
		"verified_at":    time.Now(),
	}

	return verified, response, nil
}

// verifyPaystackPayment verifies Paystack payment status
func (s *PaymentService) verifyPaystackPayment(ctx context.Context, payment *entities.Payment) (bool, map[string]interface{}, error) {
	if payment.ProviderTransactionID == nil {
		return false, nil, fmt.Errorf("no provider transaction ID")
	}

	status, err := s.paystackProvider.VerifyPayment(*payment.ProviderTransactionID)
	if err != nil {
		return false, nil, err
	}

	verified := status.Status == "success"
	response := map[string]interface{}{
		"status":         status.Status,
		"message":        status.Message,
		"transaction_id": status.TransactionID,
		"verified_at":    time.Now(),
	}

	return verified, response, nil
}

// ticketJWTClaims defines the claims embedded in each ticket's QR code JWT
type ticketJWTClaims struct {
	TicketID     string `json:"tid"`
	EventID      string `json:"eid"`
	SerialNumber string `json:"sn"`
	OrderLineID  string `json:"olid"`
	jwt.RegisteredClaims
}

// generateTickets generates tickets for a paid order.
// For each order line it:
//  1. Looks up the ticket tier to get the event ID
//  2. Creates N tickets (N = line.Quantity) with human-readable serial numbers
//  3. Signs each ticket's QR code data as a JWT (HS256, no expiry — tickets are permanent)
//  4. Atomically increments the tier's sold count
func (s *PaymentService) generateTickets(ctx context.Context, tx repositories.Transaction, order *entities.Order) error {
	// Get order lines
	orderLines, err := tx.OrderLines().GetByOrder(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("failed to get order lines: %w", err)
	}

	var allTickets []*entities.Ticket

	for _, line := range orderLines {
		// Fetch the ticket tier to get the event ID (needed for JWT claims)
		tier, err := tx.TicketTiers().GetByID(ctx, line.TicketTierID)
		if err != nil {
			return fmt.Errorf("failed to get ticket tier %s: %w", line.TicketTierID, err)
		}

		// Derive a short event code from the event slug or ID for the serial number
		eventCode := deriveEventCode(tier.EventID)

		var lineTickets []*entities.Ticket
		for i := 0; i < line.Quantity; i++ {
			ticketID := uuid.New()

			// Human-readable serial: UDUX-{EVENTCODE}-{RANDOM6}
			serialNumber := generateTicketSerialNumber(eventCode, ticketID)

			// QR code data is a signed JWT — scanner verifies signature before any DB lookup
			qrCodeData, err := s.signTicketJWT(ticketID, tier.EventID, serialNumber, line.ID)
			if err != nil {
				return fmt.Errorf("failed to sign ticket JWT for line %s ticket %d: %w", line.ID, i+1, err)
			}

			ticket := entities.NewTicket(line.ID, serialNumber, qrCodeData)
			ticket.ID = ticketID // Use the pre-generated ID so JWT matches

			// Generate QR code image as base64 (stored for PDF/email delivery)
			qrImageBase64, err := s.qrGenerator.GenerateQRCodeBase64(qrCodeData)
			if err != nil {
				// Log but don't fail — frontend can regenerate client-side if needed
				fmt.Printf("Warning: failed to generate QR image for ticket %s: %v\n", serialNumber, err)
			} else {
				ticket.QRCodeImageURL = &qrImageBase64
			}

			if err := ticket.Validate(); err != nil {
				return fmt.Errorf("invalid ticket %s: %w", serialNumber, err)
			}

			lineTickets = append(lineTickets, ticket)
		}

		allTickets = append(allTickets, lineTickets...)

		// Atomically increment the sold count on the tier — prevents race conditions
		if err := tx.TicketTiers().IncrementSold(ctx, line.TicketTierID, line.Quantity); err != nil {
			return fmt.Errorf("failed to increment sold count for tier %s: %w", line.TicketTierID, err)
		}
	}

	// Batch insert all tickets in a single query
	if len(allTickets) > 0 {
		if err := tx.Tickets().CreateBatch(ctx, allTickets); err != nil {
			return fmt.Errorf("failed to create tickets: %w", err)
		}
	}

	// Send ticket PDF email to customer (async — don't fail if email fails)
	// IMPORTANT: Use s.eventRepo and s.orderLineRepo (not the transaction) to avoid
	// using a closed/committed transaction connection in the goroutine.
	ticketsCopy := make([]*entities.Ticket, len(allTickets))
	copy(ticketsCopy, allTickets)
	go func() {
		eventUUID, err := uuid.Parse(order.EventID)
		if err != nil {
			fmt.Printf("Warning: invalid event ID for order %s: %v\n", order.Code, err)
			return
		}

		event, err := s.eventRepo.GetByID(context.Background(), eventUUID)
		if err != nil {
			fmt.Printf("Warning: failed to fetch event for order %s: %v\n", order.Code, err)
			return
		}

		orderLinesForEmail, err := s.orderLineRepo.GetByOrderID(context.Background(), order.ID)
		if err != nil {
			fmt.Printf("Warning: failed to fetch order lines for order %s: %v\n", order.Code, err)
			return
		}

		if err := s.emailService.SendTicketPDFEmail(context.Background(), order, ticketsCopy, orderLinesForEmail, event); err != nil {
			fmt.Printf("Warning: failed to send ticket PDF email for order %s: %v\n", order.Code, err)
		}
	}()

	return nil
}

// signTicketJWT creates a signed HS256 JWT for embedding in the ticket QR code.
// The JWT contains the ticket ID, event ID, serial number, and order line ID.
// Tickets do not expire — validity is controlled by the ticket's status in the DB.
func (s *PaymentService) signTicketJWT(ticketID, eventID uuid.UUID, serialNumber string, orderLineID uuid.UUID) (string, error) {
	claims := ticketJWTClaims{
		TicketID:     ticketID.String(),
		EventID:      eventID.String(),
		SerialNumber: serialNumber,
		OrderLineID:  orderLineID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "uduxpass-tickets",
			Subject:  ticketID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generateTicketSerialNumber generates a human-readable serial number.
// Format: UDUX-{EVENTCODE}-{TICKETID_SHORT}
// Example: UDUX-DAVI-A3F2B1
func generateTicketSerialNumber(eventCode string, ticketID uuid.UUID) string {
	// Use the first 6 hex chars of the ticket UUID as the unique suffix
	shortID := strings.ToUpper(ticketID.String()[:6])
	return fmt.Sprintf("UDUX-%s-%s", eventCode, shortID)
}

// deriveEventCode derives a short 4-character event code from the event UUID.
// In production this would use the event slug, but since we only have the tier here
// we use the first 4 hex chars of the event UUID.
func deriveEventCode(eventID uuid.UUID) string {
	return strings.ToUpper(eventID.String()[:4])
}

// WebhookRequest represents a payment webhook request
type WebhookRequest struct {
	Provider entities.PaymentMethod `json:"provider"`
	Event    string                 `json:"event"`
	Data     map[string]interface{} `json:"data"`
}

// WebhookResponse represents a payment webhook response
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleWebhook handles payment provider webhooks
func (s *PaymentService) HandleWebhook(ctx context.Context, req *WebhookRequest) (*WebhookResponse, error) {
	switch req.Provider {
	case entities.PaymentMethodMoMo:
		return s.handleMoMoWebhook(ctx, req)
	case entities.PaymentMethodPaystack:
		return s.handlePaystackWebhook(ctx, req)
	default:
		return nil, entities.NewValidationError("provider", "unsupported payment provider")
	}
}

// handleMoMoWebhook handles MoMo payment webhooks
func (s *PaymentService) handleMoMoWebhook(ctx context.Context, req *WebhookRequest) (*WebhookResponse, error) {
	transactionID, ok := req.Data["transaction_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing transaction_id in webhook data")
	}

	payment, err := s.paymentRepo.GetByProviderTransactionID(ctx, entities.PaymentMethodMoMo, transactionID)
	if err != nil {
		return nil, fmt.Errorf("payment not found for transaction ID: %s", transactionID)
	}

	payment.MarkWebhookReceived()
	payment.UpdateProviderResponse(req.Data)

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	if req.Event == "payment.success" {
		_, err := s.VerifyPayment(ctx, &VerifyPaymentRequest{PaymentID: payment.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to verify payment: %w", err)
		}
	}

	return &WebhookResponse{
		Status:  "success",
		Message: "Webhook processed successfully",
	}, nil
}

// handlePaystackWebhook handles Paystack payment webhooks
func (s *PaymentService) handlePaystackWebhook(ctx context.Context, req *WebhookRequest) (*WebhookResponse, error) {
	data, ok := req.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing data in webhook")
	}

	reference, ok := data["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("missing reference in webhook data")
	}

	payment, err := s.paymentRepo.GetByProviderTransactionID(ctx, entities.PaymentMethodPaystack, reference)
	if err != nil {
		return nil, fmt.Errorf("payment not found for reference: %s", reference)
	}

	payment.MarkWebhookReceived()
	payment.UpdateProviderResponse(req.Data)

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	if req.Event == "charge.success" {
		_, err := s.VerifyPayment(ctx, &VerifyPaymentRequest{PaymentID: payment.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to verify payment: %w", err)
		}
	}

	return &WebhookResponse{
		Status:  "success",
		Message: "Webhook processed successfully",
	}, nil
}

// ConfirmPaymentManuallyRequest is used by admins to manually confirm a payment
// (e.g., in dev/test environments without a live payment gateway, or for cash payments).
type ConfirmPaymentManuallyRequest struct {
	OrderID          uuid.UUID `json:"order_id"`
	PaymentReference string    `json:"payment_reference"`
	Amount           float64   `json:"amount"`
}

// ConfirmPaymentManuallyResponse is the response after manual payment confirmation.
type ConfirmPaymentManuallyResponse struct {
	OrderID          uuid.UUID `json:"order_id"`
	OrderStatus      string    `json:"order_status"`
	TicketsGenerated int       `json:"tickets_generated"`
	Message          string    `json:"message"`
}

// ConfirmPaymentManually allows an admin to manually mark an order as paid and
// trigger ticket generation. This bypasses the payment provider and is intended
// for use in development/test environments or for cash/manual payment scenarios.
func (s *PaymentService) ConfirmPaymentManually(ctx context.Context, req *ConfirmPaymentManuallyRequest) (*ConfirmPaymentManuallyResponse, error) {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, entities.NewNotFoundError("order", "order not found")
	}

	// Only pending or confirmed orders can be manually confirmed
	if order.Status != entities.OrderStatusPending && order.Status != entities.OrderStatusConfirmed {
		return nil, entities.NewValidationError("status", fmt.Sprintf("order is already %s and cannot be confirmed", order.Status))
	}

	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Mark order as paid
	order.MarkPaid()
	if req.PaymentReference != "" {
		order.PaymentReference = &req.PaymentReference
	}
	if err := tx.Orders().Update(tx.Context(), order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Generate tickets
	if err := s.generateTickets(tx.Context(), tx, order); err != nil {
		return nil, fmt.Errorf("failed to generate tickets: %w", err)
	}

	// Count generated tickets
	tickets, _ := tx.Tickets().GetByOrder(tx.Context(), order.ID)
	ticketCount := len(tickets)

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &ConfirmPaymentManuallyResponse{
		OrderID:          order.ID,
		OrderStatus:      string(order.Status),
		TicketsGenerated: ticketCount,
		Message:          fmt.Sprintf("Payment confirmed manually. %d ticket(s) generated.", ticketCount),
	}, nil
}
