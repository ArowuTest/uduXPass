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

type inventoryHoldRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewInventoryHoldRepository(db *sqlx.DB) repositories.InventoryHoldRepository {
	return &inventoryHoldRepository{db: db}
}

func NewInventoryHoldRepositoryWithTx(tx *sqlx.Tx) repositories.InventoryHoldRepository {
	return &inventoryHoldRepository{db: tx}
}

func (r *inventoryHoldRepository) Create(ctx context.Context, hold *entities.InventoryHold) error {
	query := `
		INSERT INTO inventory_holds (
			id, order_id, ticket_tier_id, quantity, 
			expires_at, created_at
		) VALUES (
			:id, :order_id, :ticket_tier_id, :quantity,
			:expires_at, :created_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, hold)
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
		return fmt.Errorf("failed to create inventory hold: %w", err)
	}
	
	return nil
}

func (r *inventoryHoldRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.InventoryHold, error) {
	var hold entities.InventoryHold
	query := `
		SELECT ih.id, ih.order_id, ih.ticket_tier_id, ih.quantity,
			   ih.expires_at, ih.created_at,
			   tt.name as ticket_tier_name, tt.quota as ticket_tier_capacity,
			   e.name as event_title, e.slug as event_slug,
			   o.code as order_code, o.status as order_status
		FROM inventory_holds ih
		JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ih.order_id = o.id
		WHERE ih.id = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &hold, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get inventory hold by ID: %w", err)
	}
	
	return &hold, nil
}

func (r *inventoryHoldRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.InventoryHold, error) {
	var holds []*entities.InventoryHold
	query := `
		SELECT ih.id, ih.order_id, ih.ticket_tier_id, ih.quantity,
			   ih.expires_at, ih.created_at,
			   tt.name as ticket_tier_name, tt.quota as ticket_tier_capacity,
			   e.name as event_title, e.slug as event_slug,
			   o.code as order_code, o.status as order_status
		FROM inventory_holds ih
		JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ih.order_id = o.id
		WHERE ih.order_id = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true
		ORDER BY ih.created_at ASC`
	
	err := r.db.SelectContext(ctx, &holds, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory holds by order: %w", err)
	}
	
	return holds, nil
}

func (r *inventoryHoldRepository) GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter repositories.InventoryHoldFilter) ([]*entities.InventoryHold, *repositories.PaginationResult, error) {
	filter.TicketTierID = &ticketTierID
	return r.List(ctx, filter)
}

func (r *inventoryHoldRepository) Update(ctx context.Context, hold *entities.InventoryHold) error {
	query := `
		UPDATE inventory_holds SET
			order_id = :order_id,
			ticket_tier_id = :ticket_tier_id,
			quantity = :quantity,
			expires_at = :expires_at
		WHERE id = :id`
	
	result, err := r.db.NamedExecContext(ctx, query, hold)
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
		return fmt.Errorf("failed to update inventory hold: %w", err)
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

func (r *inventoryHoldRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM inventory_holds WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory hold: %w", err)
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

func (r *inventoryHoldRepository) DeleteByOrder(ctx context.Context, orderID uuid.UUID) error {
	query := `DELETE FROM inventory_holds WHERE order_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete inventory holds by order: %w", err)
	}
	
	return nil
}

func (r *inventoryHoldRepository) List(ctx context.Context, filter repositories.InventoryHoldFilter) ([]*entities.InventoryHold, *repositories.PaginationResult, error) {
	var holds []*entities.InventoryHold
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.MinQuantity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.quantity >= $%d", argIndex))
		args = append(args, *filter.MinQuantity)
		argIndex++
	}
	
	if filter.MaxQuantity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.quantity <= $%d", argIndex))
		args = append(args, *filter.MaxQuantity)
		argIndex++
	}
	
	if filter.ExpiresFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.expires_at >= $%d", argIndex))
		args = append(args, *filter.ExpiresFrom)
		argIndex++
	}
	
	if filter.ExpiresTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.expires_at <= $%d", argIndex))
		args = append(args, *filter.ExpiresTo)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	if filter.ExpiredOnly {
		whereConditions = append(whereConditions, "ih.expires_at < NOW()")
	}
	
	if filter.ActiveOnly {
		whereConditions = append(whereConditions, "ih.expires_at >= NOW()")
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM inventory_holds ih
		JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ih.order_id = o.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count inventory holds: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "ih.expires_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("ih.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT ih.id, ih.order_id, ih.ticket_tier_id, ih.quantity,
			   ih.expires_at, ih.created_at,
			   tt.name as ticket_tier_name, tt.quota as ticket_tier_capacity,
			   e.name as event_title, e.slug as event_slug,
			   o.code as order_code, o.status as order_status
		FROM inventory_holds ih
		JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ih.order_id = o.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &holds, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list inventory holds: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return holds, pagination, nil
}

func (r *inventoryHoldRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM inventory_holds WHERE id = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check inventory hold existence: %w", err)
	}
	
	return exists, nil
}

func (r *inventoryHoldRepository) GetExpired(ctx context.Context, filter repositories.InventoryHoldFilter) ([]*entities.InventoryHold, *repositories.PaginationResult, error) {
	filter.ExpiredOnly = true
	return r.List(ctx, filter)
}

func (r *inventoryHoldRepository) GetActive(ctx context.Context, filter repositories.InventoryHoldFilter) ([]*entities.InventoryHold, *repositories.PaginationResult, error) {
	filter.ActiveOnly = true
	return r.List(ctx, filter)
}

func (r *inventoryHoldRepository) CleanupExpired(ctx context.Context) (int, error) {
	query := `DELETE FROM inventory_holds WHERE expires_at < NOW()`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired inventory holds: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return int(rowsAffected), nil
}

func (r *inventoryHoldRepository) ExtendExpiry(ctx context.Context, holdID uuid.UUID, newExpiryTime time.Time) error {
	query := `
		UPDATE inventory_holds 
		SET expires_at = $1
		WHERE id = $2 AND expires_at >= NOW()`
	
	result, err := r.db.ExecContext(ctx, query, newExpiryTime, holdID)
	if err != nil {
		return fmt.Errorf("failed to extend inventory hold expiry: %w", err)
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

func (r *inventoryHoldRepository) GetTotalHeldQuantity(ctx context.Context, ticketTierID uuid.UUID) (int, error) {
	var totalQuantity int
	query := `
		SELECT COALESCE(SUM(quantity), 0)
		FROM inventory_holds
		WHERE ticket_tier_id = $1 AND expires_at >= NOW()`
	
	err := r.db.GetContext(ctx, &totalQuantity, query, ticketTierID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total held quantity: %w", err)
	}
	
	return totalQuantity, nil
}

func (r *inventoryHoldRepository) GetHoldsByExpiry(ctx context.Context, expiryTime time.Time, filter repositories.InventoryHoldFilter) ([]*entities.InventoryHold, *repositories.PaginationResult, error) {
	filter.ExpiresTo = &expiryTime
	return r.List(ctx, filter)
}

func (r *inventoryHoldRepository) GetStats(ctx context.Context, filter repositories.InventoryHoldFilter) (*repositories.InventoryHoldStats, error) {
	var stats repositories.InventoryHoldStats
	
	// Build WHERE clause for stats
	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ih.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COALESCE(COUNT(*), 0) as total_holds,
			COALESCE(COUNT(CASE WHEN ih.expires_at >= NOW() THEN 1 END), 0) as active_holds,
			COALESCE(COUNT(CASE WHEN ih.expires_at < NOW() THEN 1 END), 0) as expired_holds,
			COALESCE(SUM(ih.quantity), 0) as total_quantity_held,
			COALESCE(SUM(CASE WHEN ih.expires_at >= NOW() THEN ih.quantity ELSE 0 END), 0) as active_quantity_held,
			COALESCE(SUM(CASE WHEN ih.expires_at < NOW() THEN ih.quantity ELSE 0 END), 0) as expired_quantity_held,
			COALESCE(AVG(ih.quantity), 0) as avg_quantity_per_hold,
			COALESCE(MIN(ih.quantity), 0) as min_quantity_per_hold,
			COALESCE(MAX(ih.quantity), 0) as max_quantity_per_hold,
			COALESCE(COUNT(DISTINCT ih.ticket_tier_id), 0) as unique_ticket_tiers,
			COALESCE(COUNT(DISTINCT ih.order_id), 0) as unique_orders,
			MIN(ih.created_at) as first_hold_created_at,
			MAX(ih.created_at) as last_hold_created_at,
			MIN(CASE WHEN ih.expires_at >= NOW() THEN ih.expires_at END) as next_expiry_at,
			MAX(ih.expires_at) as latest_expiry_at
		FROM inventory_holds ih
		JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ih.order_id = o.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory hold stats: %w", err)
	}
	
	return &stats, nil
}


// DeleteBySessionID deletes inventory holds by session ID
func (r *inventoryHoldRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	query := `DELETE FROM inventory_holds WHERE session_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete inventory holds by session ID: %w", err)
	}
	
	return nil
}

// GetTotalHeld returns the total quantity held for a ticket tier (alias for GetTotalHeldQuantity)
func (r *inventoryHoldRepository) GetTotalHeld(ctx context.Context, ticketTierID uuid.UUID) (int, error) {
	return r.GetTotalHeldQuantity(ctx, ticketTierID)
}


// GetByOrderID retrieves inventory holds by order ID (alias for GetByOrder)
func (r *inventoryHoldRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.InventoryHold, error) {
	return r.GetByOrder(ctx, orderID)
}

