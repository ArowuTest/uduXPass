package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

type orderRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
	realDB *sqlx.DB // Keep reference to actual DB for transaction management
}

func NewOrderRepository(db *sqlx.DB) repositories.OrderRepository {
	return &orderRepository{db: db, realDB: db}
}

func NewOrderRepositoryWithTx(tx *sqlx.Tx) repositories.OrderRepository {
	// When using a transaction, we don't need realDB since we're already in a transaction
	return &orderRepository{db: tx, realDB: nil}
}

	func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
		// Check if we're already in a transaction (realDB is nil when using WithTx)
		var tx interface {
			NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		}
		var shouldCommit bool
		
		if r.realDB != nil {
			// Not in a transaction, start one
			newTx, err := r.realDB.BeginTxx(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to begin transaction: %w", err)
			}
			defer newTx.Rollback()
			tx = newTx
			shouldCommit = true
		} else {
			// Already in a transaction, use the existing one
			tx = r.db
			shouldCommit = false
		}

	// Insert order
	orderQuery := `
		INSERT INTO orders (
			id, user_id, event_id, code, secret, status, total_amount, 
			currency, email, phone, first_name, last_name,
			customer_email, customer_phone, customer_first_name, 
			customer_last_name, expires_at, created_at, updated_at
		) VALUES (
			:id, :user_id, :event_id, :code, :secret, :status, :total_amount,
			:currency, :email, :phone, :first_name, :last_name,
			:customer_email, :customer_phone, :customer_first_name,
			:customer_last_name, :expires_at, :created_at, :updated_at
		)`
	
		_, err := tx.NamedExecContext(ctx, orderQuery, order)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "code") {
					return entities.ErrOrderCodeExists
				}
			}
		}
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order lines
	if len(order.OrderLines) > 0 {
		lineQuery := `
			INSERT INTO order_lines (
				id, order_id, ticket_tier_id, quantity, unit_price, 
				total_price, created_at, updated_at
			) VALUES (
				:id, :order_id, :ticket_tier_id, :quantity, :unit_price,
				:total_price, :created_at, :updated_at
			)`
		
		for _, line := range order.OrderLines {
			_, err = tx.NamedExecContext(ctx, lineQuery, line)
			if err != nil {
				return fmt.Errorf("failed to create order line: %w", err)
			}
		}
	}

		// Only commit if we started the transaction
		if shouldCommit {
			if commitTx, ok := tx.(*sqlx.Tx); ok {
				return commitTx.Commit()
			}
		}
		return nil
	}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	
	orderQuery := `
		SELECT o.id, o.user_id, o.event_id, o.code, o.secret, o.status,
			   o.total_amount, o.currency, o.customer_email, o.customer_phone,
			   o.customer_first_name, o.customer_last_name, o.payment_method,
			   o.payment_reference, o.paid_at, o.expires_at, o.created_at, o.updated_at,
			   o.is_active
		FROM orders o
		WHERE o.id = $1 AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &order, orderQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	
	// Convert timestamps to UTC (database stores without timezone)
	order.ExpiresAt = order.ExpiresAt.UTC()
	order.CreatedAt = order.CreatedAt.UTC()
	order.UpdatedAt = order.UpdatedAt.UTC()
	if order.PaidAt != nil {
		utcTime := order.PaidAt.UTC()
		order.PaidAt = &utcTime
	}

	// Load order lines
	lines, err := r.getOrderLines(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load order lines: %w", err)
	}
	// Convert []*entities.OrderLine to []entities.OrderLine
	order.OrderLines = make([]entities.OrderLine, len(lines))
	for i, line := range lines {
		order.OrderLines[i] = *line
	}

	return &order, nil
}

func (r *orderRepository) GetByCode(ctx context.Context, code string) (*entities.Order, error) {
	var order entities.Order
	
	orderQuery := `
		SELECT o.id, o.user_id, o.event_id, o.code, o.secret, o.status,
			   o.total_amount, o.currency, o.customer_email, o.customer_phone,
			   o.customer_first_name, o.customer_last_name, o.payment_method,
			   o.payment_reference, o.paid_at, o.expires_at, o.created_at, o.updated_at,
			   o.is_active
		FROM orders o
		WHERE o.code = $1 AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &order, orderQuery, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by code: %w", err)
	}

	// Load order lines
	lines, err := r.getOrderLines(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load order lines: %w", err)
	}
	// Convert []*entities.OrderLine to []entities.OrderLine
	order.OrderLines = make([]entities.OrderLine, len(lines))
	for i, line := range lines {
		order.OrderLines[i] = *line
	}

	return &order, nil
}

func (r *orderRepository) GetBySecret(ctx context.Context, secret string) (*entities.Order, error) {
	var order entities.Order
	
	orderQuery := `
		SELECT o.id, o.user_id, o.event_id, o.code, o.secret, o.status,
			   o.total_amount, o.currency, o.customer_email, o.customer_phone,
			   o.customer_first_name, o.customer_last_name, o.payment_method,
			   o.payment_reference, o.paid_at, o.expires_at, o.created_at, o.updated_at,
			   e.name as event_title, e.event_date as event_event_date,
			   e.venue_name, e.venue_city
		FROM orders o
		JOIN events e ON o.event_id = e.id
		WHERE o.secret = $1 AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &order, orderQuery, secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by secret: %w", err)
	}

	// Load order lines
	lines, err := r.getOrderLines(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load order lines: %w", err)
	}
	// Convert []*entities.OrderLine to []entities.OrderLine
	order.OrderLines = make([]entities.OrderLine, len(lines))
	for i, line := range lines {
		order.OrderLines[i] = *line
	}

	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) error {
	order.UpdatedAt = time.Now()
	
	query := `
		UPDATE orders SET
			status = :status,
			total_amount = :total_amount,
			currency = :currency,
			customer_email = :customer_email,
			customer_phone = :customer_phone,
			customer_first_name = :customer_first_name,
			customer_last_name = :customer_last_name,
			payment_method = :payment_method,
			payment_reference = :payment_reference,
			paid_at = :paid_at,
			expires_at = :expires_at,
			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	
	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE orders SET is_active = false WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	
	return nil
}

func (r *orderRepository) List(ctx context.Context, filter repositories.OrderFilter) ([]*entities.Order, *repositories.PaginationResult, error) {
	var orders []*entities.Order
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}
	
	if filter.Email != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("o.customer_email = $%d", argIndex))
		args = append(args, filter.Email)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM orders o 
		JOIN events e ON o.event_id = e.id 
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count orders: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "o.created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("o.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT o.id, o.user_id, o.event_id, o.code, o.secret, o.status,
			   o.total_amount, o.currency, o.customer_email, o.customer_phone,
			   o.customer_first_name, o.customer_last_name, o.payment_method,
			   o.payment_reference, o.paid_at, o.expires_at, o.created_at, o.updated_at,
			   e.name as event_title, e.event_date as event_event_date,
			   e.venue_name, e.venue_city
		FROM orders o
		JOIN events e ON o.event_id = e.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	// Load order lines for each order
	for _, order := range orders {
		lines, err := r.getOrderLines(ctx, order.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load order lines: %w", err)
		}
		// Convert []*entities.OrderLine to []entities.OrderLine
		order.OrderLines = make([]entities.OrderLine, len(lines))
		for i, line := range lines {
			order.OrderLines[i] = *line
		}
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return orders, pagination, nil
}

func (r *orderRepository) GetByUser(ctx context.Context, userID uuid.UUID, filter repositories.OrderFilter) ([]*entities.Order, *repositories.PaginationResult, error) {
	filter.UserID = &userID
	return r.List(ctx, filter)
}

func (r *orderRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, filter repositories.OrderFilter) ([]*entities.Order, *repositories.PaginationResult, error) {
	filter.EventID = &eventID
	return r.List(ctx, filter)
}

func (r *orderRepository) GetByEmail(ctx context.Context, email string, filter repositories.OrderFilter) ([]*entities.Order, *repositories.PaginationResult, error) {
	filter.Email = email
	return r.List(ctx, filter)
}

func (r *orderRepository) GetExpired(ctx context.Context, filter repositories.OrderFilter) ([]*entities.Order, *repositories.PaginationResult, error) {
	var orders []*entities.Order
	var totalCount int
	
	// Build WHERE clause for expired orders
	whereConditions := []string{
		"o.is_active = true",
		"o.status = 'pending'",
		"o.expires_at < NOW()",
	}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM orders o 
		JOIN events e ON o.event_id = e.id 
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count expired orders: %w", err)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT o.id, o.user_id, o.event_id, o.code, o.secret, o.status,
			   o.total_amount, o.currency, o.customer_email, o.customer_phone,
			   o.customer_first_name, o.customer_last_name, o.payment_method,
			   o.payment_reference, o.paid_at, o.expires_at, o.created_at, o.updated_at,
			   e.name as event_title, e.event_date as event_event_date,
			   e.venue_name, e.venue_city
		FROM orders o
		JOIN events e ON o.event_id = e.id
		WHERE %s 
		ORDER BY o.expires_at ASC 
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list expired orders: %w", err)
	}
	
	// Load order lines for each order
	for _, order := range orders {
		lines, err := r.getOrderLines(ctx, order.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load order lines: %w", err)
		}
		// Convert []*entities.OrderLine to []entities.OrderLine
		order.OrderLines = make([]entities.OrderLine, len(lines))
		for i, line := range lines {
			order.OrderLines[i] = *line
		}
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return orders, pagination, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error {
	query := `
		UPDATE orders 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	
	return nil
}

func (r *orderRepository) MarkExpired(ctx context.Context) (int, error) {
	query := `
		UPDATE orders 
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'pending' AND expires_at < NOW() AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to mark orders as expired: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return int(rowsAffected), nil
}

func (r *orderRepository) GetOrderStats(ctx context.Context, orderID uuid.UUID) (*repositories.OrderStats, error) {
	var stats repositories.OrderStats
	
	query := `
		SELECT 
			$1 as event_id,
			COALESCE(COUNT(o.id), 0) as total_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'paid' THEN 1 END), 0) as paid_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'pending' THEN 1 END), 0) as pending_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'expired' THEN 1 END), 0) as expired_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'cancelled' THEN 1 END), 0) as cancelled_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'refunded' THEN 1 END), 0) as refunded_orders,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN o.total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(CASE WHEN o.status = 'paid' THEN o.total_amount END), 0) as average_order_value,
			COALESCE(SUM(ol.quantity), 0) as total_tickets_sold,
			MIN(CASE WHEN o.status = 'paid' THEN o.created_at END) as first_sale_at,
			MAX(CASE WHEN o.status = 'paid' THEN o.created_at END) as last_sale_at
		FROM events e
		LEFT JOIN orders o ON e.id = o.event_id AND o.is_active = true
		LEFT JOIN order_lines ol ON o.id = ol.order_id AND o.status = 'paid'
		WHERE e.id = $1 AND e.is_active = true
		GROUP BY e.id`
	
	err := r.db.GetContext(ctx, &stats, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}
	
	return &stats, nil
}

func (r *orderRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}
	
	return exists, nil
}

func (r *orderRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE code = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, code)
	if err != nil {
		return false, fmt.Errorf("failed to check order code existence: %w", err)
	}
	
	return exists, nil
}

// Helper function to load order lines
func (r *orderRepository) getOrderLines(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderLine, error) {
	var lines []*entities.OrderLine
	
	query := `
		SELECT ol.id, ol.order_id, ol.ticket_tier_id, ol.quantity, 
			   ol.unit_price, ol.subtotal, ol.created_at, ol.updated_at
		FROM order_lines ol
		WHERE ol.order_id = $1
		ORDER BY ol.created_at ASC`
	
	err := r.db.SelectContext(ctx, &lines, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}
	
	return lines, nil
}


// GetByEventID retrieves orders by event ID (simplified version of GetByEvent)
func (r *orderRepository) GetByEventID(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	var orders []*entities.Order
	query := `
		SELECT id, code, event_id, user_id, email, phone, first_name, last_name,
			   customer_first_name, customer_last_name, customer_email, customer_phone,
			   status, total_amount, currency, payment_method, payment_id, notes, confirmed_at,
			   expires_at, secret, locale, comment, meta_info, created_at, updated_at
		FROM orders 
		WHERE event_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	err := r.db.SelectContext(ctx, &orders, query, eventID.String(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by event ID: %w", err)
	}
	
	return orders, nil
}

// GetByUserID retrieves orders by user ID (simplified version of GetByUser)
func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	var orders []*entities.Order
	query := `
		SELECT id, code, event_id, user_id, email, phone, first_name, last_name,
			   customer_first_name, customer_last_name, customer_email, customer_phone,
			   status, total_amount, currency, payment_method, payment_id, notes, confirmed_at,
			   expires_at, secret, locale, comment, meta_info, created_at, updated_at
		FROM orders 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	err := r.db.SelectContext(ctx, &orders, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}
	
	return orders, nil
}


// GetExpiredOrders retrieves expired orders (simplified version of GetExpired)
func (r *orderRepository) GetExpiredOrders(ctx context.Context) ([]*entities.Order, error) {
	var orders []*entities.Order
	query := `
		SELECT id, code, event_id, user_id, email, phone, first_name, last_name,
			   customer_first_name, customer_last_name, customer_email, customer_phone,
			   status, total_amount, currency, payment_method, payment_id, notes, confirmed_at, cancelled_at,
			   expires_at, secret, locale, comment, meta_info, created_at, updated_at
		FROM orders 
		WHERE status = 'pending' AND expires_at <= NOW() AND is_active = true
		ORDER BY expires_at ASC`
	
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired orders: %w", err)
	}
	
	return orders, nil
}

