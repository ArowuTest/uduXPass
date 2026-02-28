package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

type orderLineRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewOrderLineRepository(db *sqlx.DB) repositories.OrderLineRepository {
	return &orderLineRepository{db: db}
}

func NewOrderLineRepositoryWithTx(tx *sqlx.Tx) repositories.OrderLineRepository {
	return &orderLineRepository{db: tx}
}

func (r *orderLineRepository) Create(ctx context.Context, orderLine *entities.OrderLine) error {
	query := `
		INSERT INTO order_lines (
			id, order_id, ticket_tier_id, quantity, unit_price, 
			subtotal, total_price, fees, taxes, discount_amount, created_at, updated_at
		) VALUES (
			:id, :order_id, :ticket_tier_id, :quantity, :unit_price,
			:subtotal, :total_price, :fees, :taxes, :discount_amount, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, orderLine)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
				if strings.Contains(pqErr.Detail, "ticket_tier_id") {
					return entities.ErrNotFoundError
				}
			case "23514": // check_constraint_violation
				return entities.ErrValidationError
			}
		}
		return fmt.Errorf("failed to create order line: %w", err)
	}
	
	return nil
}

func (r *orderLineRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderLine, error) {
	var orderLine entities.OrderLine
	query := `
		SELECT ol.id, ol.order_id, ol.ticket_tier_id, ol.quantity, 
			   ol.unit_price, ol.subtotal, ol.created_at,
			   tt.name as ticket_tier_name, tt.description as ticket_tier_description,
			   e.name as event_title, e.slug as event_slug
		FROM order_lines ol
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		WHERE ol.id = $1 AND tt.is_active = true AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &orderLine, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get order line by ID: %w", err)
	}
	
	return &orderLine, nil
}

func (r *orderLineRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderLine, error) {
	var orderLines []*entities.OrderLine
	query := `
		SELECT ol.id, ol.order_id, ol.ticket_tier_id, ol.quantity, 
			   ol.unit_price, ol.subtotal, ol.created_at,
			   tt.name as ticket_tier_name, tt.description as ticket_tier_description,
			   e.name as event_title, e.slug as event_slug
		FROM order_lines ol
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		WHERE ol.order_id = $1 AND tt.is_active = true AND e.is_active = true
		ORDER BY ol.created_at ASC`
	
	err := r.db.SelectContext(ctx, &orderLines, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines by order: %w", err)
	}
	
	return orderLines, nil
}

// GetOrderTotal calculates the total amount for an order
func (r *orderLineRepository) GetOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(subtotal), 0) 
		FROM order_lines 
		WHERE order_id = $1`
	
	err := r.db.GetContext(ctx, &total, query, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to get order total: %w", err)
	}
	
	return total, nil
}

// GetOrderQuantity calculates the total quantity for an order
func (r *orderLineRepository) GetOrderQuantity(ctx context.Context, orderID uuid.UUID) (int, error) {
	var quantity int
	query := `
		SELECT COALESCE(SUM(quantity), 0) 
		FROM order_lines 
		WHERE order_id = $1`
	
	err := r.db.GetContext(ctx, &quantity, query, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to get order quantity: %w", err)
	}
	
	return quantity, nil
}

// DeleteByOrder deletes all order lines for an order
func (r *orderLineRepository) DeleteByOrder(ctx context.Context, orderID uuid.UUID) error {
	query := `DELETE FROM order_lines WHERE order_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order lines by order: %w", err)
	}
	
	return nil
}

func (r *orderLineRepository) GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter repositories.OrderLineFilter) ([]*entities.OrderLine, *repositories.PaginationResult, error) {
	filter.TicketTierID = &ticketTierID
	return r.List(ctx, filter)
}

// GetByEvent retrieves order lines for a specific event
func (r *orderLineRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, filter repositories.OrderLineFilter) ([]*entities.OrderLine, *repositories.PaginationResult, error) {
	filter.EventID = &eventID
	return r.List(ctx, filter)
}

// GetEventSales retrieves sales statistics for an event
func (r *orderLineRepository) GetEventSales(ctx context.Context, eventID uuid.UUID) (*repositories.OrderLineStats, error) {
	filter := repositories.OrderLineFilter{
		EventID: &eventID,
	}
	return r.GetStats(ctx, filter)
}

// GetTicketTierSales retrieves sales statistics for a ticket tier
func (r *orderLineRepository) GetTicketTierSales(ctx context.Context, ticketTierID uuid.UUID) (*repositories.OrderLineStats, error) {
	filter := repositories.OrderLineFilter{
		TicketTierID: &ticketTierID,
	}
	return r.GetStats(ctx, filter)
}

func (r *orderLineRepository) Update(ctx context.Context, orderLine *entities.OrderLine) error {
	query := `
			UPDATE order_lines SET
				order_id = :order_id,
				ticket_tier_id = :ticket_tier_id,
				quantity = :quantity,
				unit_price = :unit_price,
				subtotal = :subtotal,
				fees = :fees,
				taxes = :taxes,
				discount_amount = :discount_amount,
				updated_at = :updated_at
			WHERE id = :id`
	
	result, err := r.db.NamedExecContext(ctx, query, orderLine)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
				if strings.Contains(pqErr.Detail, "ticket_tier_id") {
					return entities.ErrNotFoundError
				}
			case "23514": // check_constraint_violation
				return entities.ErrValidationError
			}
		}
		return fmt.Errorf("failed to update order line: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrNotFoundError
	}
	
	return nil
}

func (r *orderLineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM order_lines WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order line: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrNotFoundError
	}
	
	return nil
}

func (r *orderLineRepository) List(ctx context.Context, filter repositories.OrderLineFilter) ([]*entities.OrderLine, *repositories.PaginationResult, error) {
	var orderLines []*entities.OrderLine
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"tt.is_active = true", "e.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.MinQuantity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.quantity >= $%d", argIndex))
		args = append(args, *filter.MinQuantity)
		argIndex++
	}
	
	if filter.MaxQuantity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.quantity <= $%d", argIndex))
		args = append(args, *filter.MaxQuantity)
		argIndex++
	}
	
	if filter.MinPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.subtotal >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}
	
	if filter.MaxPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.subtotal <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM order_lines ol
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count order lines: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "ol.created_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("ol.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT ol.id, ol.order_id, ol.ticket_tier_id, ol.quantity, 
			   ol.unit_price, ol.subtotal, ol.created_at,
			   tt.name as ticket_tier_name, tt.description as ticket_tier_description,
			   e.name as event_title, e.slug as event_slug
		FROM order_lines ol
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &orderLines, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list order lines: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return orderLines, pagination, nil
}

func (r *orderLineRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM order_lines WHERE id = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check order line existence: %w", err)
	}
	
	return exists, nil
}

func (r *orderLineRepository) GetTotalQuantityByTicketTier(ctx context.Context, ticketTierID uuid.UUID) (int, error) {
	var totalQuantity int
	query := `
		SELECT COALESCE(SUM(ol.quantity), 0)
		FROM order_lines ol
		JOIN orders o ON ol.order_id = o.id
		WHERE ol.ticket_tier_id = $1 
		AND o.status IN ('paid', 'pending') 
		AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &totalQuantity, query, ticketTierID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total quantity by ticket tier: %w", err)
	}
	
	return totalQuantity, nil
}

func (r *orderLineRepository) GetTotalRevenueByTicketTier(ctx context.Context, ticketTierID uuid.UUID) (float64, error) {
	var totalRevenue float64
	query := `
		SELECT COALESCE(SUM(ol.subtotal), 0)
		FROM order_lines ol
		JOIN orders o ON ol.order_id = o.id
		WHERE ol.ticket_tier_id = $1 
		AND o.status = 'paid' 
		AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &totalRevenue, query, ticketTierID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue by ticket tier: %w", err)
	}
	
	return totalRevenue, nil
}

func (r *orderLineRepository) GetTotalQuantityByEvent(ctx context.Context, eventID uuid.UUID) (int, error) {
	var totalQuantity int
	query := `
		SELECT COALESCE(SUM(ol.quantity), 0)
		FROM order_lines ol
		JOIN orders o ON ol.order_id = o.id
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		WHERE tt.event_id = $1 
		AND o.status IN ('paid', 'pending') 
		AND o.is_active = true
		AND tt.is_active = true`
	
	err := r.db.GetContext(ctx, &totalQuantity, query, eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total quantity by event: %w", err)
	}
	
	return totalQuantity, nil
}

func (r *orderLineRepository) GetTotalRevenueByEvent(ctx context.Context, eventID uuid.UUID) (float64, error) {
	var totalRevenue float64
	query := `
		SELECT COALESCE(SUM(ol.subtotal), 0)
		FROM order_lines ol
		JOIN orders o ON ol.order_id = o.id
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		WHERE tt.event_id = $1 
		AND o.status = 'paid' 
		AND o.is_active = true
		AND tt.is_active = true`
	
	err := r.db.GetContext(ctx, &totalRevenue, query, eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue by event: %w", err)
	}
	
	return totalRevenue, nil
}

func (r *orderLineRepository) GetStats(ctx context.Context, filter repositories.OrderLineFilter) (*repositories.OrderLineStats, error) {
	var stats repositories.OrderLineStats
	
	// Build WHERE clause for stats
	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COALESCE(COUNT(*), 0) as total_lines,
			COALESCE(SUM(ol.quantity), 0) as total_quantity,
			COALESCE(SUM(ol.subtotal), 0) as total_amount,
			COALESCE(AVG(ol.quantity), 0) as avg_quantity_per_line,
			COALESCE(AVG(ol.subtotal), 0) as avg_amount_per_line,
			COALESCE(MIN(ol.unit_price), 0) as min_unit_price,
			COALESCE(MAX(ol.unit_price), 0) as max_unit_price,
			COALESCE(COUNT(DISTINCT ol.ticket_tier_id), 0) as unique_ticket_tiers,
			COALESCE(COUNT(DISTINCT CASE WHEN o.status = 'paid' THEN ol.id END), 0) as paid_lines,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN ol.subtotal ELSE 0 END), 0) as paid_amount
		FROM order_lines ol
		JOIN orders o ON ol.order_id = o.id
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get order line stats: %w", err)
	}
	
	return &stats, nil
}


// CreateBatch creates multiple order lines in a single transaction
func (r *orderLineRepository) CreateBatch(ctx context.Context, orderLines []*entities.OrderLine) error {
	if len(orderLines) == 0 {
		return nil
	}
	
	query := `
		INSERT INTO order_lines (
			id, order_id, ticket_tier_id, quantity, unit_price, subtotal,
			fees, taxes, discount_amount, created_at, updated_at
		) VALUES (
			:id, :order_id, :ticket_tier_id, :quantity, :unit_price, :subtotal,
			:fees, :taxes, :discount_amount, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, orderLines)
	if err != nil {
		return fmt.Errorf("failed to create order lines batch: %w", err)
	}
	
	return nil
}

// UpdateBatch updates multiple order lines in a single transaction
func (r *orderLineRepository) UpdateBatch(ctx context.Context, orderLines []*entities.OrderLine) error {
	if len(orderLines) == 0 {
		return nil
	}
	
	query := `
		UPDATE order_lines SET
			quantity = :quantity,
			unit_price = :unit_price,
			subtotal = :subtotal,
			fees = :fees,
			taxes = :taxes,
			discount_amount = :discount_amount,
			updated_at = :updated_at
		WHERE id = :id`
	
	_, err := r.db.NamedExecContext(ctx, query, orderLines)
	if err != nil {
		return fmt.Errorf("failed to update order lines batch: %w", err)
	}
	
	return nil
}


// GetByOrderID retrieves order lines by order ID (alias for GetByOrder)
func (r *orderLineRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderLine, error) {
	return r.GetByOrder(ctx, orderID)
}

