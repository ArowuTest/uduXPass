package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

// OrderService handles order operations
type OrderService struct {
	orderRepo         repositories.OrderRepository
	orderLineRepo     repositories.OrderLineRepository
	inventoryHoldRepo repositories.InventoryHoldRepository
	eventRepo         repositories.EventRepository
	ticketTierRepo    repositories.TicketTierRepository
	userRepo          repositories.UserRepository
	holdDuration      time.Duration
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repositories.OrderRepository,
	orderLineRepo repositories.OrderLineRepository,
	inventoryHoldRepo repositories.InventoryHoldRepository,
	eventRepo repositories.EventRepository,
	ticketTierRepo repositories.TicketTierRepository,
	userRepo repositories.UserRepository,
) *OrderService {
	return &OrderService{
		orderRepo:         orderRepo,
		orderLineRepo:     orderLineRepo,
		inventoryHoldRepo: inventoryHoldRepo,
		eventRepo:         eventRepo,
		ticketTierRepo:    ticketTierRepo,
		userRepo:          userRepo,
		holdDuration:      15 * time.Minute, // 15 minutes hold
	}
}

// CreateOrderRequest represents a create order request
type CreateOrderRequest struct {
	UserID      uuid.UUID              `json:"user_id" validate:"required"`
	EventID     uuid.UUID              `json:"event_id" validate:"required"`
	OrderLines  []CreateOrderLineItem  `json:"order_lines" validate:"required,min=1"`
	CustomerInfo *CustomerInfo         `json:"customer_info,omitempty"`
}

// CreateOrderLineItem represents an order line item
type CreateOrderLineItem struct {
	TicketTierID uuid.UUID `json:"ticket_tier_id" validate:"required"`
	Quantity     int       `json:"quantity" validate:"required,min=1"`
}

// CustomerInfo represents customer information
type CustomerInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// CreateOrderResponse represents a create order response
type CreateOrderResponse struct {
	Order       *entities.Order      `json:"order"`
	OrderLines  []*entities.OrderLine `json:"order_lines"`
	TotalAmount float64              `json:"total_amount"`
	ExpiresAt   time.Time            `json:"expires_at"`
}

// UpdateOrderRequest represents an update order request
type UpdateOrderRequest struct {
	Status      *entities.OrderStatus `json:"status,omitempty"`
	PaymentID   *uuid.UUID           `json:"payment_id,omitempty"`
	Notes       *string              `json:"notes,omitempty"`
}

// CreateOrder creates a new order with inventory holds
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate event exists and is active
	event, err := s.eventRepo.GetByID(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event.Status != entities.EventStatusPublished {
		return nil, entities.ErrEventNotActive
	}

	// Check if event sales are still open
	if event.SaleEnd != nil && time.Now().After(*event.SaleEnd) {
		return nil, entities.ErrEventExpired
	}

	// Create order
	customerEmail := ""
	if req.CustomerInfo != nil && req.CustomerInfo.Email != "" {
		customerEmail = req.CustomerInfo.Email
	} else if user.Email != nil {
		customerEmail = *user.Email
	}
	
	order := entities.NewOrder(req.EventID.String(), customerEmail)
	
	// Set customer info if provided
	if req.CustomerInfo != nil {
		order.CustomerFirstName = req.CustomerInfo.FirstName
		order.CustomerLastName = req.CustomerInfo.LastName
		order.CustomerEmail = req.CustomerInfo.Email
		order.CustomerPhone = req.CustomerInfo.Phone
	} else {
		// Use user info as fallback
		if user.FirstName != nil {
			order.CustomerFirstName = *user.FirstName
		}
		if user.LastName != nil {
			order.CustomerLastName = *user.LastName
		}
		if user.Email != nil {
			order.CustomerEmail = *user.Email
		}
		if user.Phone != nil {
			order.CustomerPhone = *user.Phone
		}
	}

	// Set order expiry (in UTC to match database storage)
	expiresAt := time.Now().UTC().Add(s.holdDuration)
	order.ExpiresAt = expiresAt

	// Save order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	var orderLines []*entities.OrderLine
	var totalAmount float64

	// Process each order line
	for _, lineItem := range req.OrderLines {
		// Get ticket tier
		ticketTier, err := s.ticketTierRepo.GetByID(ctx, lineItem.TicketTierID)
		if err != nil {
			return nil, fmt.Errorf("failed to get ticket tier: %w", err)
		}

		// Check availability
		available, err := s.ticketTierRepo.GetAvailableQuantity(ctx, lineItem.TicketTierID)
		if err != nil {
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}

		if available < lineItem.Quantity {
			return nil, entities.ErrInsufficientTickets
		}

		// Create inventory hold
		inventoryHold := entities.NewInventoryHold(
			order.ID,
			lineItem.TicketTierID,
			lineItem.Quantity,
			s.holdDuration,
		)

		if err := s.inventoryHoldRepo.Create(ctx, inventoryHold); err != nil {
			return nil, fmt.Errorf("failed to create inventory hold: %w", err)
		}

		// Create order line
		orderLine := entities.NewOrderLine(
			order.ID,
			lineItem.TicketTierID,
			lineItem.Quantity,
			ticketTier.Price,
		)

		if err := s.orderLineRepo.Create(ctx, orderLine); err != nil {
			return nil, fmt.Errorf("failed to create order line: %w", err)
		}

		orderLines = append(orderLines, orderLine)
		totalAmount += float64(lineItem.Quantity) * ticketTier.Price
	}

	// Update order total
	order.TotalAmount = totalAmount
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order total: %w", err)
	}

	return &CreateOrderResponse{
		Order:       order,
		OrderLines:  orderLines,
		TotalAmount: totalAmount,
		ExpiresAt:   expiresAt,
	}, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*entities.Order, error) {
	return s.orderRepo.GetByID(ctx, orderID)
}

// GetOrderWithLines retrieves an order with its order lines
func (s *OrderService) GetOrderWithLines(ctx context.Context, orderID uuid.UUID) (*entities.Order, []*entities.OrderLine, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	orderLines, err := s.orderLineRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	return order, orderLines, nil
}

// GetUserOrders retrieves orders for a user
func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetEventOrders retrieves orders for an event
func (s *OrderService) GetEventOrders(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	return s.orderRepo.GetByEventID(ctx, eventID, limit, offset)
}

// UpdateOrder updates an order
func (s *OrderService) UpdateOrder(ctx context.Context, orderID uuid.UUID, req *UpdateOrderRequest) (*entities.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Status != nil {
		order.Status = *req.Status
	}

	if req.PaymentID != nil {
		paymentIDStr := req.PaymentID.String()
		order.PaymentID = &paymentIDStr
	}

	if req.Notes != nil {
		order.Notes = req.Notes
	}

	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// ConfirmOrder confirms an order and converts holds to actual allocations
func (s *OrderService) ConfirmOrder(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != entities.OrderStatusPending {
		return fmt.Errorf("order is not in pending status")
	}

	// Update order status
	order.Status = entities.OrderStatusConfirmed
	paymentIDStr := paymentID.String()
	order.PaymentID = &paymentIDStr
	order.ConfirmedAt = &[]time.Time{time.Now()}[0]

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Convert inventory holds to actual allocations
	holds, err := s.inventoryHoldRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, hold := range holds {
		hold.Status = entities.InventoryHoldStatusConfirmed
		if err := s.inventoryHoldRepo.Update(ctx, hold); err != nil {
			return err
		}
	}

	return nil
}

// CancelOrder cancels an order and releases inventory holds
func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, reason string) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == entities.OrderStatusCancelled {
		return fmt.Errorf("order is already cancelled")
	}

	// Update order status
	order.Status = entities.OrderStatusCancelled
	order.CancelledAt = &[]time.Time{time.Now()}[0]
	if reason != "" {
		order.Notes = &reason
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Release inventory holds
	holds, err := s.inventoryHoldRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, hold := range holds {
		hold.Status = entities.InventoryHoldStatusReleased
		if err := s.inventoryHoldRepo.Update(ctx, hold); err != nil {
			return err
		}
	}

	return nil
}

// ExpireOrder expires an order and releases inventory holds
func (s *OrderService) ExpireOrder(ctx context.Context, orderID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != entities.OrderStatusPending {
		return fmt.Errorf("order is not in pending status")
	}

	// Update order status
	order.Status = entities.OrderStatusExpired

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Release inventory holds
	holds, err := s.inventoryHoldRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, hold := range holds {
		hold.Status = entities.InventoryHoldStatusExpired
		if err := s.inventoryHoldRepo.Update(ctx, hold); err != nil {
			return err
		}
	}

	return nil
}

// ProcessExpiredOrders processes expired orders (background job)
func (s *OrderService) ProcessExpiredOrders(ctx context.Context) error {
	expiredOrders, err := s.orderRepo.GetExpiredOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range expiredOrders {
		if err := s.ExpireOrder(ctx, order.ID); err != nil {
			// Log error but continue processing other orders
			fmt.Printf("Failed to expire order %s: %v\n", order.ID, err)
		}
	}

	return nil
}

// GetOrderStats retrieves order statistics
func (s *OrderService) GetOrderStats(ctx context.Context, eventID *uuid.UUID) (*OrderStats, error) {
	var eventUUID uuid.UUID
	if eventID != nil {
		eventUUID = *eventID
	}
	
	repoStats, err := s.orderRepo.GetOrderStats(ctx, eventUUID)
	if err != nil {
		return nil, err
	}
	
	// Convert from repositories.OrderStats to service OrderStats
	return &OrderStats{
		TotalOrders:     repoStats.TotalOrders,
		PendingOrders:   repoStats.PendingOrders,
		ConfirmedOrders: repoStats.PaidOrders, // Map PaidOrders to ConfirmedOrders
		CancelledOrders: repoStats.CancelledOrders,
		ExpiredOrders:   repoStats.ExpiredOrders,
		TotalRevenue:    repoStats.TotalRevenue,
	}, nil
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders     int     `json:"total_orders"`
	PendingOrders   int     `json:"pending_orders"`
	ConfirmedOrders int     `json:"confirmed_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
	ExpiredOrders   int     `json:"expired_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
}

