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

type ticketTierRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewTicketTierRepository(db *sqlx.DB) repositories.TicketTierRepository {
	return &ticketTierRepository{db: db}
}

func NewTicketTierRepositoryWithTx(tx *sqlx.Tx) repositories.TicketTierRepository {
	return &ticketTierRepository{db: tx}
}

func (r *ticketTierRepository) Create(ctx context.Context, tier *entities.TicketTier) error {
	query := `
		INSERT INTO ticket_tiers (
			id, event_id, name, description, price, currency, 
			quota, max_per_order, sale_start, sale_end, 
			is_active, position, created_at, updated_at
		) VALUES (
			:id, :event_id, :name, :description, :price, :currency,
			:quota, :max_per_order, :sale_start, :sale_end,
			:is_active, :position, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, tier)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return entities.ErrConflictError
			case "23503": // foreign_key_violation
				return entities.ErrEventNotFound
			}
		}
		return fmt.Errorf("failed to create ticket tier: %w", err)
	}
	
	return nil
}

func (r *ticketTierRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.TicketTier, error) {
	var tier entities.TicketTier
	query := `
		SELECT tt.id, tt.event_id, tt.name, tt.description, tt.price, tt.currency,
			   tt.quota, tt.sold, tt.min_purchase, tt.max_purchase, tt.sale_start, tt.sale_end,
			   tt.is_active, tt.created_at, tt.updated_at
		FROM ticket_tiers tt
		WHERE tt.id = $1 AND tt.is_active = true`
	
	err := r.db.GetContext(ctx, &tier, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get ticket tier by ID: %w", err)
	}
	
	return &tier, nil
}

func (r *ticketTierRepository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]*entities.TicketTier, error) {
	filter := repositories.TicketTierFilter{
		EventID: &eventID,
	}
	tiers, _, err := r.List(ctx, filter)
	return tiers, err
}

// GetActiveByEvent retrieves active ticket tiers for a specific event
func (r *ticketTierRepository) GetActiveByEvent(ctx context.Context, eventID uuid.UUID) ([]*entities.TicketTier, error) {
	var tiers []*entities.TicketTier
	
	query := `
		SELECT tt.id, tt.event_id, tt.name, tt.description, tt.price,
			   tt.quota, tt.sold, tt.min_purchase, tt.max_purchase,
			   tt.sale_start, tt.sale_end, tt.is_active,
			   tt.created_at, tt.updated_at
		FROM ticket_tiers tt
		WHERE tt.event_id = $1 
		AND tt.is_active = true
		AND (tt.sale_start IS NULL OR tt.sale_start <= NOW())
		AND (tt.sale_end IS NULL OR tt.sale_end >= NOW())
		ORDER BY tt.created_at ASC`
	
	err := r.db.SelectContext(ctx, &tiers, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active ticket tiers by event: %w", err)
	}
	
	return tiers, nil
}

func (r *ticketTierRepository) Update(ctx context.Context, tier *entities.TicketTier) error {
	tier.UpdatedAt = time.Now()
	
	query := `
		UPDATE ticket_tiers SET
			event_id = :event_id,
			name = :name,
			description = :description,
			price = :price,
			currency = :currency,
			quota = :quota,
			max_per_order = :max_per_order,
			sale_start = :sale_start,
			sale_end = :sale_end,
			is_active = :is_active,
			position = :position,
			meta_info = :meta_info,
			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, tier)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				return entities.ErrEventNotFound
			}
		}
		return fmt.Errorf("failed to update ticket tier: %w", err)
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

func (r *ticketTierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE ticket_tiers SET is_active = false WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket tier: %w", err)
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

func (r *ticketTierRepository) List(ctx context.Context, filter repositories.TicketTierFilter) ([]*entities.TicketTier, *repositories.PaginationResult, error) {
	var tiers []*entities.TicketTier
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"tt.is_active = true", "e.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	if filter.MinPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.price >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}
	
	if filter.MaxPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.price <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(tt.name ILIKE $%d OR tt.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}
	
	if filter.SaleStartFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.sale_start >= $%d", argIndex))
		args = append(args, *filter.SaleStartFrom)
		argIndex++
	}
	
	if filter.SaleStartTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.sale_start <= $%d", argIndex))
		args = append(args, *filter.SaleStartTo)
		argIndex++
	}
	
	if filter.SaleEndFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.sale_end >= $%d", argIndex))
		args = append(args, *filter.SaleEndFrom)
		argIndex++
	}
	
	if filter.SaleEndTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.sale_end <= $%d", argIndex))
		args = append(args, *filter.SaleEndTo)
		argIndex++
	}
	
	// Add availability filter if requested
	if filter.AvailableOnly {
		now := time.Now()
		whereConditions = append(whereConditions, fmt.Sprintf("tt.is_active = true AND tt.sale_start <= $%d AND tt.sale_end >= $%d", argIndex, argIndex+1))
		args = append(args, now, now)
		argIndex += 2
		
		// Check quota vs sold tickets
		whereConditions = append(whereConditions, `
			tt.quota > COALESCE((
				SELECT SUM(ol.quantity) 
				FROM order_lines ol 
				JOIN orders o ON ol.order_id = o.id 
				WHERE ol.ticket_tier_id = tt.id 
				AND o.status = 'paid' 
				AND o.is_active = true
			), 0)`)
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM ticket_tiers tt 
		JOIN events e ON tt.event_id = e.id 
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count ticket tiers: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "tt.position ASC, tt.created_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("tt.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT tt.id, tt.event_id, tt.name, tt.description, tt.price, tt.currency,
			   tt.quota, tt.max_per_order, tt.sale_start, tt.sale_end,
			   tt.is_active, tt.position, tt.meta_info, tt.created_at, tt.updated_at,
			   e.name as event_title, e.slug as event_slug,
			   COALESCE((
				   SELECT SUM(ol.quantity) 
				   FROM order_lines ol 
				   JOIN orders o ON ol.order_id = o.id 
				   WHERE ol.ticket_tier_id = tt.id 
				   AND o.status = 'paid' 
				   AND o.is_active = true
			   ), 0) as sold_count,
			   (tt.quota - COALESCE((
				   SELECT SUM(ol.quantity) 
				   FROM order_lines ol 
				   JOIN orders o ON ol.order_id = o.id 
				   WHERE ol.ticket_tier_id = tt.id 
				   AND o.status = 'paid' 
				   AND o.is_active = true
			   ), 0)) as available_count
		FROM ticket_tiers tt
		JOIN events e ON tt.event_id = e.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &tiers, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list ticket tiers: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return tiers, pagination, nil
}

func (r *ticketTierRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM ticket_tiers WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket tier existence: %w", err)
	}
	
	return exists, nil
}

func (r *ticketTierRepository) GetAvailable(ctx context.Context, eventID uuid.UUID) ([]*entities.TicketTier, error) {
	filter := repositories.TicketTierFilter{
		BaseFilter: repositories.BaseFilter{
			Page:  1,
			Limit: 100,
		},
		EventID:       &eventID,
		AvailableOnly: true,
	}
	
	tiers, _, err := r.List(ctx, filter)
	return tiers, err
}

func (r *ticketTierRepository) GetAvailability(ctx context.Context, eventID uuid.UUID) ([]*repositories.TicketTierAvailability, error) {
	var availability []*repositories.TicketTierAvailability
	
	query := `
		SELECT 
			tt.id as ticket_tier_id,
			tt.name,
			tt.price,
			tt.currency,
			tt.quota as quota,
			COALESCE(sold.count, 0) as sold,
			COALESCE(reserved.count, 0) as reserved,
			CASE 
				WHEN tt.quota IS NULL THEN 999999
				ELSE tt.quota - COALESCE(sold.count, 0) - COALESCE(reserved.count, 0)
			END as available,
			tt.is_active as is_on_sale,
			CASE 
				WHEN NOT tt.is_active THEN 'inactive'
				WHEN tt.sale_event_date IS NOT NULL AND tt.sale_event_date > NOW() THEN 'not_started'
				WHEN tt.sale_end_date IS NOT NULL AND tt.sale_end_date < NOW() THEN 'ended'
				WHEN tt.quota IS NOT NULL AND COALESCE(sold.count, 0) >= tt.quota THEN 'sold_out'
				ELSE 'available'
			END as sale_status
		FROM ticket_tiers tt
		LEFT JOIN (
			SELECT ol.ticket_tier_id, SUM(ol.quantity) as count
			FROM order_lines ol
			JOIN orders o ON ol.order_id = o.id
			WHERE o.status IN ('confirmed', 'paid') AND o.is_active = true
			GROUP BY ol.ticket_tier_id
		) sold ON tt.id = sold.ticket_tier_id
		LEFT JOIN (
			SELECT ih.ticket_tier_id, SUM(ih.quantity) as count
			FROM inventory_holds ih
			WHERE ih.expires_at > NOW()
			GROUP BY ih.ticket_tier_id
		) reserved ON tt.id = reserved.ticket_tier_id
		WHERE tt.event_id = $1 AND tt.is_active = true
		ORDER BY tt.position ASC, tt.created_at ASC`
	
	err := r.db.SelectContext(ctx, &availability, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket tier availability: %w", err)
	}
	
	return availability, nil
}

// GetTierStats retrieves statistics for a specific ticket tier
func (r *ticketTierRepository) GetTierStats(ctx context.Context, tierID uuid.UUID) (*repositories.TicketTierStats, error) {
	var stats repositories.TicketTierStats
	
	query := `
		SELECT 
			tt.id as tier_id,
			tt.name as tier_name,
			tt.price,
			tt.currency,
			tt.quota,
			COALESCE(sold.count, 0) as sold_count,
			COALESCE(sold.revenue, 0) as revenue,
			COALESCE(reserved.count, 0) as reserved_count,
			CASE 
				WHEN tt.quota IS NULL THEN 999999
				ELSE tt.quota - COALESCE(sold.count, 0) - COALESCE(reserved.count, 0)
			END as available_count
		FROM ticket_tiers tt
		LEFT JOIN (
			SELECT ol.ticket_tier_id, SUM(ol.quantity) as count, SUM(ol.total_amount) as revenue
			FROM order_lines ol
			JOIN orders o ON ol.order_id = o.id
			WHERE o.status IN ('confirmed', 'paid') AND o.is_active = true
			GROUP BY ol.ticket_tier_id
		) sold ON tt.id = sold.ticket_tier_id
		LEFT JOIN (
			SELECT ih.ticket_tier_id, SUM(ih.quantity) as count
			FROM inventory_holds ih
			WHERE ih.expires_at > NOW()
			GROUP BY ih.ticket_tier_id
		) reserved ON tt.id = reserved.ticket_tier_id
		WHERE tt.id = $1 AND tt.is_active = true`
	
	err := r.db.GetContext(ctx, &stats, query, tierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get ticket tier stats: %w", err)
	}
	
	return &stats, nil
}

func (r *ticketTierRepository) UpdatePosition(ctx context.Context, positions map[uuid.UUID]int) error {
	// Use individual updates since we can't guarantee transaction support in interface
	query := `UPDATE ticket_tiers SET position = $1, updated_at = NOW() WHERE id = $2`
	
	for tierID, position := range positions {
		_, err := r.db.ExecContext(ctx, query, position, tierID)
		if err != nil {
			return fmt.Errorf("failed to update ticket tier position: %w", err)
		}
	}
	
	return nil
}

func (r *ticketTierRepository) UpdateCapacity(ctx context.Context, tierID uuid.UUID, quota int) error {
	query := `
		UPDATE ticket_tiers 
		SET quota = $1, updated_at = NOW()
		WHERE id = $2 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, quota, tierID)
	if err != nil {
		return fmt.Errorf("failed to update ticket tier quota: %w", err)
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

func (r *ticketTierRepository) GetStats(ctx context.Context, tierID uuid.UUID) (*repositories.TicketTierStats, error) {
	var stats repositories.TicketTierStats
	
	query := `
		SELECT 
			tt.id as tier_id,
			tt.name as tier_name,
			tt.quota,
			tt.price,
			COALESCE(COUNT(DISTINCT o.id), 0) as total_orders,
			COALESCE(COUNT(DISTINCT CASE WHEN o.status = 'paid' THEN o.id END), 0) as paid_orders,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN ol.quantity ELSE 0 END), 0) as tickets_sold,
			COALESCE(SUM(CASE WHEN o.status = 'pending' THEN ol.quantity ELSE 0 END), 0) as tickets_reserved,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN ol.total_price ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(CASE WHEN o.status = 'paid' THEN ol.quantity END), 0) as avg_quantity_per_order,
			MIN(CASE WHEN o.status = 'paid' THEN o.created_at END) as first_sale_at,
			MAX(CASE WHEN o.status = 'paid' THEN o.created_at END) as last_sale_at
		FROM ticket_tiers tt
		LEFT JOIN order_lines ol ON tt.id = ol.ticket_tier_id
		LEFT JOIN orders o ON ol.order_id = o.id AND o.is_active = true
		WHERE tt.id = $1 AND tt.is_active = true
		GROUP BY tt.id, tt.name, tt.quota, tt.price`
	
	err := r.db.GetContext(ctx, &stats, query, tierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket tier stats: %w", err)
	}
	
	return &stats, nil
}


// GetAvailableQuantity retrieves the available quantity for a ticket tier
func (r *ticketTierRepository) GetAvailableQuantity(ctx context.Context, ticketTierID uuid.UUID) (int, error) {
	query := `
		SELECT 
			tt.quota - COALESCE(SUM(
				CASE 
					WHEN o.status IN ('paid', 'pending') THEN ol.quantity 
					ELSE 0 
				END
			), 0) as available_quantity
		FROM ticket_tiers tt
		LEFT JOIN order_lines ol ON tt.id = ol.ticket_tier_id
		LEFT JOIN orders o ON ol.order_id = o.id AND o.is_active = true
		WHERE tt.id = $1 AND tt.is_active = true
		GROUP BY tt.id, tt.quota`
	
	var availableQuantity int
	err := r.db.GetContext(ctx, &availableQuantity, query, ticketTierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, entities.ErrTicketTierNotFound
		}
		return 0, fmt.Errorf("failed to get available quantity for ticket tier: %w", err)
	}
	
	// Ensure we don't return negative quantities
	if availableQuantity < 0 {
		availableQuantity = 0
	}
	
	return availableQuantity, nil
}

