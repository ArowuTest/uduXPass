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

type ticketRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewTicketRepository(db *sqlx.DB) repositories.TicketRepository {
	return &ticketRepository{db: db}
}

func NewTicketRepositoryWithTx(tx *sqlx.Tx) repositories.TicketRepository {
	return &ticketRepository{db: tx}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, order_id, ticket_tier_id, code, qr_code, status, 
			attendee_name, attendee_email, attendee_phone, 
			redeemed_at, redeemed_by, meta_info, created_at, updated_at
		) VALUES (
			:id, :order_id, :ticket_tier_id, :code, :qr_code, :status,
			:attendee_name, :attendee_email, :attendee_phone,
			:redeemed_at, :redeemed_by, :meta_info, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "code") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "qr_code") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
				if strings.Contains(pqErr.Detail, "ticket_tier_id") {
					return entities.ErrNotFoundError
				}
			}
		}
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	
	return nil
}

// CreateBatch creates multiple tickets in a batch
func (r *ticketRepository) CreateBatch(ctx context.Context, tickets []*entities.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}
	
	query := `
		INSERT INTO tickets (
			id, order_id, ticket_tier_id, code, qr_code, status, 
			attendee_name, attendee_email, attendee_phone, 
			redeemed_at, redeemed_by, meta_info, created_at, updated_at
		) VALUES (
			:id, :order_id, :ticket_tier_id, :code, :qr_code, :status,
			:attendee_name, :attendee_email, :attendee_phone,
			:redeemed_at, :redeemed_by, :meta_info, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, tickets)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "code") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "qr_code") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
				if strings.Contains(pqErr.Detail, "ticket_tier_id") {
					return entities.ErrNotFoundError
				}
			}
		}
		return fmt.Errorf("failed to create tickets batch: %w", err)
	}
	
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Ticket, error) {
	var ticket entities.Ticket
	query := `
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE t.id = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}
	
	return &ticket, nil
}

func (r *ticketRepository) GetByCode(ctx context.Context, code string) (*entities.Ticket, error) {
	var ticket entities.Ticket
	query := `
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE t.code = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &ticket, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by code: %w", err)
	}
	
	return &ticket, nil
}

// GetBySerialNumber retrieves a ticket by its serial number (code)
func (r *ticketRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*entities.Ticket, error) {
	return r.GetByCode(ctx, serialNumber)
}

	func (r *ticketRepository) GetByQRCode(ctx context.Context, qrCode string) (*entities.Ticket, error) {
		var ticket entities.Ticket
		
		// Enhanced query with JOINs to get complete ticket information
		query := `
			SELECT 
				t.id, 
				t.order_line_id, 
				t.serial_number, 
				t.qr_code_data, 
				t.status,
				t.redeemed_at, 
				t.redeemed_by, 
				t.created_at, 
				t.updated_at
			FROM tickets t
			INNER JOIN order_lines ol ON t.order_line_id = ol.id
			INNER JOIN orders o ON ol.order_id = o.id
			INNER JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
			INNER JOIN events e ON tt.event_id = e.id
			WHERE t.qr_code_data = $1
				AND t.status != 'cancelled'
				AND o.status = 'paid'
				AND tt.is_active = true
				AND e.status = 'published'`
		
		err := r.db.GetContext(ctx, &ticket, query, qrCode)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, entities.ErrTicketNotFound
			}
			return nil, fmt.Errorf("failed to get ticket by QR code: %w", err)
		}
		
		return &ticket, nil
	}

func (r *ticketRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	query := `
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE t.order_id = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true
		ORDER BY t.created_at ASC`
	
	err := r.db.SelectContext(ctx, &tickets, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets by order: %w", err)
	}
	
	return tickets, nil
}

// GetByOrderLine retrieves tickets for a specific order line
func (r *ticketRepository) GetByOrderLine(ctx context.Context, orderLineID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	query := `
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		JOIN order_lines ol ON t.order_id = ol.order_id AND t.ticket_tier_id = ol.ticket_tier_id
		WHERE ol.id = $1 AND tt.is_active = true AND e.is_active = true AND o.is_active = true
		ORDER BY t.created_at ASC`
	
	err := r.db.SelectContext(ctx, &tickets, query, orderLineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets by order line: %w", err)
	}
	
	return tickets, nil
}

func (r *ticketRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.EventID = &eventID
	return r.List(ctx, filter)
}

func (r *ticketRepository) GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.TicketTierID = &ticketTierID
	return r.List(ctx, filter)
}

// GetByUser retrieves tickets for a specific user
func (r *ticketRepository) GetByUser(ctx context.Context, userID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.UserID = &userID
	return r.List(ctx, filter)
}

// GetTicketStats retrieves ticket statistics
func (r *ticketRepository) GetTicketStats(ctx context.Context, filter repositories.TicketStatsFilter) (*repositories.TicketStats, error) {
	var stats repositories.TicketStats
	
	whereConditions := []string{"t.is_active = true", "tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}
	
	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_tickets,
			COUNT(CASE WHEN t.status = 'active' THEN 1 END) as active_tickets,
			COUNT(CASE WHEN t.status = 'redeemed' THEN 1 END) as redeemed_tickets,
			COUNT(CASE WHEN t.status = 'cancelled' THEN 1 END) as cancelled_tickets,
			COALESCE(SUM(tt.price), 0) as total_value
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket stats: %w", err)
	}
	
	return &stats, nil
}

// GetUpcoming retrieves upcoming tickets for a user
func (r *ticketRepository) GetUpcoming(ctx context.Context, userID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	
	query := `
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE o.user_id = $1 
		  AND t.status = 'active'
		  AND e.event_date > NOW()
		  AND tt.is_active = true 
		  AND e.is_active = true 
		  AND o.is_active = true
		ORDER BY e.event_date ASC
		LIMIT 10`
	
	err := r.db.SelectContext(ctx, &tickets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming tickets: %w", err)
	}
	
	return tickets, nil
}

// MarkRedeemed marks a ticket as redeemed
func (r *ticketRepository) MarkRedeemed(ctx context.Context, ticketID uuid.UUID, redeemedBy string) error {
	now := time.Now()
	
	query := `
		UPDATE tickets 
		SET status = 'redeemed', redeemed_at = $1, redeemed_by = $2, updated_at = $1
		WHERE id = $3 AND status = 'active' AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, now, redeemedBy, ticketID)
	if err != nil {
		return fmt.Errorf("failed to mark ticket as redeemed: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketNotFound
	}
	
	return nil
}

// MarkVoided marks a ticket as voided
func (r *ticketRepository) MarkVoided(ctx context.Context, ticketID uuid.UUID) error {
	now := time.Now()
	
	query := `
		UPDATE tickets 
		SET status = 'voided', updated_at = $1
		WHERE id = $2 AND status IN ('active', 'redeemed') AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, now, ticketID)
	if err != nil {
		return fmt.Errorf("failed to mark ticket as voided: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketNotFound
	}
	
	return nil
}

// ValidateForRedemption validates if a ticket can be redeemed
func (r *ticketRepository) ValidateForRedemption(ctx context.Context, qrCode string, eventID uuid.UUID) (*repositories.TicketValidationResult, error) {
	var result repositories.TicketValidationResult
	
	query := `
		SELECT 
			t.id as ticket_id,
			t.status,
			t.attendee_name,
			t.attendee_email,
			e.name as event_title,
			e.event_date as event_event_date,
			e.end_date as event_end_date,
			tt.name as ticket_tier_name,
			o.code as order_code
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE t.qr_code_data = $1 
		  AND e.id = $2
		  AND t.is_active = true 
		  AND tt.is_active = true 
		  AND e.is_active = true 
		  AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &result, query, qrCode, eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			result.Valid = false
			result.Message = "Ticket not found"
			return &result, nil
		}
		return nil, fmt.Errorf("failed to validate ticket: %w", err)
	}
	
	// Check if ticket is already redeemed
	if result.Status == "redeemed" {
		result.Valid = false
		result.AlreadyRedeemed = true
		result.Message = "Ticket already redeemed"
		return &result, nil
	}
	
	// Check if ticket is voided
	if result.Status == "voided" {
		result.Valid = false
		result.Message = "Ticket is voided"
		return &result, nil
	}
	
	// Check if ticket is active
	if result.Status != "active" {
		result.Valid = false
		result.Message = "Ticket is not active"
		return &result, nil
	}
	
	// Ticket is valid for redemption
	result.Valid = true
	result.Message = "Ticket is valid for redemption"
	
	return &result, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *entities.Ticket) error {
	ticket.UpdatedAt = time.Now()
	
	query := `
		UPDATE tickets SET
			order_id = :order_id,
			ticket_tier_id = :ticket_tier_id,
			code = :code,
			qr_code_data = :qr_code,
			status = :status,
			attendee_name = :attendee_name,
			attendee_email = :attendee_email,
			attendee_phone = :attendee_phone,
			redeemed_at = :redeemed_at,
			redeemed_by = :redeemed_by,
			meta_info = :meta_info,
			updated_at = :updated_at
		WHERE id = :id`
	
	result, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "code") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "qr_code") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
				if strings.Contains(pqErr.Detail, "ticket_tier_id") {
					return entities.ErrNotFoundError
				}
			}
		}
		return fmt.Errorf("failed to update ticket: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketNotFound
	}
	
	return nil
}

func (r *ticketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tickets WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketNotFound
	}
	
	return nil
}

func (r *ticketRepository) List(ctx context.Context, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	var tickets []*entities.Ticket
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	
	if filter.AttendeeEmail != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("t.attendee_email = $%d", argIndex))
		args = append(args, filter.AttendeeEmail)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(t.code ILIKE $%d OR t.attendee_name ILIKE $%d OR t.attendee_email ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	if filter.RedeemedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.redeemed_at >= $%d", argIndex))
		args = append(args, *filter.RedeemedFrom)
		argIndex++
	}
	
	if filter.RedeemedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.redeemed_at <= $%d", argIndex))
		args = append(args, *filter.RedeemedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count tickets: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "t.created_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("t.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT t.id, t.order_id, t.ticket_tier_id, t.code, t.qr_code_data, t.status,
			   t.attendee_name, t.attendee_email, t.attendee_phone,
			   t.redeemed_at, t.redeemed_by, t.meta_info, t.created_at, t.updated_at,
			   tt.name as ticket_tier_name, tt.price as ticket_tier_price,
			   e.name as event_title, e.slug as event_slug, e.event_date as event_event_date,
			   e.venue_name, e.venue_address, e.venue_city,
			   o.code as order_code, o.email as order_email
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return tickets, pagination, nil
}

func (r *ticketRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket existence: %w", err)
	}
	
	return exists, nil
}

func (r *ticketRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE code = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, code)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket code existence: %w", err)
	}
	
	return exists, nil
}

func (r *ticketRepository) ExistsByQRCode(ctx context.Context, qrCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE qr_code_data_data = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, qrCode)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket QR code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsBySerialNumber checks if a ticket exists by serial number
func (r *ticketRepository) ExistsBySerialNumber(ctx context.Context, serialNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE code = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, serialNumber)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket serial number existence: %w", err)
	}
	
	return exists, nil
}

func (r *ticketRepository) UpdateStatus(ctx context.Context, ticketID uuid.UUID, status entities.TicketStatus) error {
	query := `
		UPDATE tickets 
		SET status = $1, updated_at = NOW()
		WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, status, ticketID)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketNotFound
	}
	
	return nil
}

func (r *ticketRepository) Redeem(ctx context.Context, ticketID uuid.UUID, redeemedBy string) error {
	now := time.Now()
	query := `
		UPDATE tickets 
		SET status = 'redeemed', redeemed_at = $1, redeemed_by = $2, updated_at = $1
		WHERE id = $3 AND status = 'valid'`
	
	result, err := r.db.ExecContext(ctx, query, now, redeemedBy, ticketID)
	if err != nil {
		return fmt.Errorf("failed to redeem ticket: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrTicketAlreadyRedeemed
	}
	
	return nil
}

func (r *ticketRepository) ValidateForEntry(ctx context.Context, code string) (*repositories.TicketValidationResult, error) {
	var result repositories.TicketValidationResult
	
	query := `
		SELECT 
			t.id as ticket_id,
			t.code,
			t.status,
			t.attendee_name,
			t.attendee_email,
			t.redeemed_at,
			t.redeemed_by,
			tt.name as ticket_tier_name,
			e.id as event_id,
			e.name as event_title,
			e.event_date as event_event_date,
			e.end_date as event_end_date,
			e.venue_name,
			o.code as order_code,
			o.status as order_status,
			CASE 
				WHEN t.status = 'redeemed' THEN false
				WHEN t.status != 'valid' THEN false
				WHEN o.status != 'paid' THEN false
				WHEN e.event_date > NOW() + INTERVAL '1 hour' THEN false
				WHEN e.end_date < NOW() - INTERVAL '1 hour' THEN false
				ELSE true
			END as is_valid,
			CASE 
				WHEN t.status = 'redeemed' THEN 'Ticket already redeemed'
				WHEN t.status != 'valid' THEN 'Ticket is not valid'
				WHEN o.status != 'paid' THEN 'Order not paid'
				WHEN e.event_date > NOW() + INTERVAL '1 hour' THEN 'Event has not started yet'
				WHEN e.end_date < NOW() - INTERVAL '1 hour' THEN 'Event has ended'
				ELSE 'Valid for entry'
			END as validation_message
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE t.code = $1 
		AND tt.is_active = true 
		AND e.is_active = true 
		AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &result, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return &repositories.TicketValidationResult{
				Valid:   false,
				Message: "Ticket not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to validate ticket for entry: %w", err)
	}
	
	return &result, nil
}

func (r *ticketRepository) GetStats(ctx context.Context, filter repositories.TicketFilter) (*repositories.TicketStats, error) {
	var stats repositories.TicketStats
	
	// Build WHERE clause for stats
	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COALESCE(COUNT(*), 0) as total_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'valid' THEN 1 END), 0) as valid_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'redeemed' THEN 1 END), 0) as redeemed_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'cancelled' THEN 1 END), 0) as cancelled_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'refunded' THEN 1 END), 0) as refunded_tickets,
			COALESCE(COUNT(CASE WHEN o.status = 'paid' THEN 1 END), 0) as paid_tickets,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN tt.price ELSE 0 END), 0) as total_value,
			COALESCE(COUNT(DISTINCT t.attendee_email), 0) as unique_attendees,
			MIN(CASE WHEN t.status = 'redeemed' THEN t.redeemed_at END) as first_redemption_at,
			MAX(CASE WHEN t.status = 'redeemed' THEN t.redeemed_at END) as last_redemption_at,
			ROUND(
				CASE 
					WHEN COUNT(*) > 0 THEN 
						(COUNT(CASE WHEN t.status = 'redeemed' THEN 1 END)::float / COUNT(*)::float) * 100
					ELSE 0 
				END, 2
			) as redemption_rate
		FROM tickets t
		JOIN ticket_tiers tt ON t.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON t.order_id = o.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket stats: %w", err)
	}
	
	return &stats, nil
}

