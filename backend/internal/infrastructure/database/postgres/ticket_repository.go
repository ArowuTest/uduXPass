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

// ticketRepository implements repositories.TicketRepository.
//
// Schema design note:
//   The `tickets` table uses `order_line_id` as its only FK (→ order_lines.id).
//   To reach orders, ticket_tiers, and events we JOIN through order_lines:
//     tickets t
//       → order_lines ol  ON t.order_line_id = ol.id
//       → ticket_tiers tt ON ol.ticket_tier_id = tt.id
//       → events e         ON tt.event_id = e.id
//       → orders o         ON ol.order_id = o.id
//
//   Ticket entity DB fields:
//     id, order_line_id, serial_number, qr_code_data, qr_code_image_url,
//     status, redeemed_at, redeemed_by, created_at, updated_at

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

// ticketSelectColumns is the canonical SELECT list for a single ticket row.
// It selects only columns that exist in the tickets table and map to the Ticket entity.
const ticketSelectColumns = `
	t.id,
	t.order_line_id,
	t.serial_number,
	t.qr_code_data,
	t.qr_code_image_url,
	t.status,
	t.redeemed_at,
	t.redeemed_by,
	t.created_at,
	t.updated_at`

// ticketJoinClause is the standard JOIN chain from tickets through to orders.
const ticketJoinClause = `
	JOIN order_lines ol ON t.order_line_id = ol.id
	JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
	JOIN events e ON tt.event_id = e.id
	JOIN orders o ON ol.order_id = o.id`

// Create inserts a single ticket using only the entity's actual DB columns.
func (r *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, order_line_id, serial_number, qr_code_data, qr_code_image_url,
			status, redeemed_at, redeemed_by, created_at, updated_at
		) VALUES (
			:id, :order_line_id, :serial_number, :qr_code_data, :qr_code_image_url,
			:status, :redeemed_at, :redeemed_by, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "serial_number") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "qr_code_data") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_line_id") {
					return entities.ErrOrderNotFound
				}
			}
		}
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	return nil
}

// CreateBatch inserts multiple tickets in a single named-exec call.
func (r *ticketRepository) CreateBatch(ctx context.Context, tickets []*entities.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}

	// sqlx NamedExecContext with a slice performs a single multi-row INSERT.
	query := `
		INSERT INTO tickets (
			id, order_line_id, serial_number, qr_code_data, qr_code_image_url,
			status, redeemed_at, redeemed_by, created_at, updated_at
		) VALUES (
			:id, :order_line_id, :serial_number, :qr_code_data, :qr_code_image_url,
			:status, :redeemed_at, :redeemed_by, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, tickets)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "serial_number") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "qr_code_data") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_line_id") {
					return entities.ErrOrderNotFound
				}
			}
		}
		return fmt.Errorf("failed to create tickets batch: %w", err)
	}

	return nil
}

// GetByID retrieves a ticket by its primary key.
func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Ticket, error) {
	var ticket entities.Ticket
	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE t.id = $1
		  AND tt.is_active = true
		  AND e.is_active = true
		  AND o.is_active = true`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	return &ticket, nil
}

// GetByCode retrieves a ticket by its serial_number (the human-readable code).
func (r *ticketRepository) GetByCode(ctx context.Context, code string) (*entities.Ticket, error) {
	var ticket entities.Ticket
	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE t.serial_number = $1
		  AND tt.is_active = true
		  AND e.is_active = true
		  AND o.is_active = true`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.GetContext(ctx, &ticket, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by code: %w", err)
	}

	return &ticket, nil
}

// GetBySerialNumber is an alias for GetByCode.
func (r *ticketRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*entities.Ticket, error) {
	return r.GetByCode(ctx, serialNumber)
}

// GetByQRCode retrieves a ticket by its QR code data string.
// This is the primary lookup used by the scanner validation flow.
// It checks that the associated order is in 'paid' status and the event is published.
func (r *ticketRepository) GetByQRCode(ctx context.Context, qrCode string) (*entities.Ticket, error) {
	var ticket entities.Ticket

	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE t.qr_code_data = $1
		  AND t.status != 'cancelled'
		  AND o.status = 'paid'
		  AND tt.is_active = true
		  AND e.status = 'published'`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.GetContext(ctx, &ticket, query, qrCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by QR code: %w", err)
	}

	return &ticket, nil
}

// GetByOrder retrieves all tickets for a given order, joining through order_lines.
func (r *ticketRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE ol.order_id = $1
		  AND tt.is_active = true
		  AND e.is_active = true
		  AND o.is_active = true
		ORDER BY t.created_at ASC`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.SelectContext(ctx, &tickets, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets by order: %w", err)
	}

	return tickets, nil
}

// GetByOrderLine retrieves all tickets for a specific order line.
func (r *ticketRepository) GetByOrderLine(ctx context.Context, orderLineID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE t.order_line_id = $1
		  AND tt.is_active = true
		  AND e.is_active = true
		  AND o.is_active = true
		ORDER BY t.created_at ASC`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.SelectContext(ctx, &tickets, query, orderLineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets by order line: %w", err)
	}

	return tickets, nil
}

// GetByEvent retrieves tickets for a specific event using the List method.
func (r *ticketRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.EventID = &eventID
	return r.List(ctx, filter)
}

// GetByTicketTier retrieves tickets for a specific ticket tier using the List method.
func (r *ticketRepository) GetByTicketTier(ctx context.Context, ticketTierID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.TicketTierID = &ticketTierID
	return r.List(ctx, filter)
}

// GetByUser retrieves tickets for a specific user using the List method.
func (r *ticketRepository) GetByUser(ctx context.Context, userID uuid.UUID, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	filter.UserID = &userID
	return r.List(ctx, filter)
}

// GetUpcoming retrieves upcoming active tickets for a user.
func (r *ticketRepository) GetUpcoming(ctx context.Context, userID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket

	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE o.user_id = $1
		  AND t.status = 'active'
		  AND e.event_date > NOW()
		  AND tt.is_active = true
		  AND e.is_active = true
		  AND o.is_active = true
		ORDER BY e.event_date ASC
		LIMIT 10`,
		ticketSelectColumns, ticketJoinClause)

	err := r.db.SelectContext(ctx, &tickets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming tickets: %w", err)
	}

	return tickets, nil
}

// GetTicketStats retrieves aggregate ticket statistics.
func (r *ticketRepository) GetTicketStats(ctx context.Context, filter repositories.TicketStatsFilter) (*repositories.TicketStats, error) {
	var stats repositories.TicketStats

	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1

	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
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
		%s
		WHERE %s`, ticketJoinClause, whereClause)

	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket stats: %w", err)
	}

	return &stats, nil
}

// MarkRedeemed marks a ticket as redeemed. Only succeeds if the ticket is currently 'active'.
func (r *ticketRepository) MarkRedeemed(ctx context.Context, ticketID uuid.UUID, redeemedBy string) error {
	now := time.Now()

	query := `
		UPDATE tickets
		SET status = 'redeemed', redeemed_at = $1, redeemed_by = $2, updated_at = $1
		WHERE id = $3 AND status = 'active'`

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

// MarkVoided marks a ticket as cancelled (voided). Only succeeds if the ticket is 'active' or 'redeemed'.
func (r *ticketRepository) MarkVoided(ctx context.Context, ticketID uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE tickets
		SET status = 'cancelled', updated_at = $1
		WHERE id = $2 AND status IN ('active', 'redeemed')`

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

// ValidateForRedemption validates whether a ticket identified by QR code data can be redeemed
// for a specific event. It performs all business-rule checks and returns a structured result.
func (r *ticketRepository) ValidateForRedemption(ctx context.Context, qrCode string, eventID uuid.UUID) (*repositories.TicketValidationResult, error) {
	var result repositories.TicketValidationResult

	// We JOIN through order_lines since tickets.order_line_id is the only FK.
	query := `
		SELECT
			t.id as ticket_id,
			t.status,
			o.code as order_code,
			tt.name as ticket_tier_name,
			e.name as event_title,
			e.event_date as event_event_date,
			e.end_date as event_end_date
		FROM tickets t
		JOIN order_lines ol ON t.order_line_id = ol.id
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ol.order_id = o.id
		WHERE t.qr_code_data = $1
		  AND e.id = $2
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

	switch result.Status {
	case "redeemed":
		result.Valid = false
		result.AlreadyRedeemed = true
		result.Message = "Ticket already redeemed"
	case "cancelled", "refunded", "transferred":
		result.Valid = false
		result.Message = fmt.Sprintf("Ticket is %s", result.Status)
	case "active":
		result.Valid = true
		result.Message = "Ticket is valid for redemption"
	default:
		result.Valid = false
		result.Message = "Ticket is not active"
	}

	return &result, nil
}

// Update updates a ticket's mutable fields using the entity's actual DB columns.
func (r *ticketRepository) Update(ctx context.Context, ticket *entities.Ticket) error {
	ticket.UpdatedAt = time.Now()

	query := `
		UPDATE tickets SET
			qr_code_data    = :qr_code_data,
			qr_code_image_url = :qr_code_image_url,
			status          = :status,
			redeemed_at     = :redeemed_at,
			redeemed_by     = :redeemed_by,
			updated_at      = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "qr_code_data") {
					return entities.ErrConflictError
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

// Delete permanently removes a ticket by ID.
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

// List retrieves tickets with filtering and pagination.
// All filters that reference order/event/tier data JOIN through order_lines.
func (r *ticketRepository) List(ctx context.Context, filter repositories.TicketFilter) ([]*entities.Ticket, *repositories.PaginationResult, error) {
	var tickets []*entities.Ticket
	var totalCount int

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

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(t.serial_number ILIKE $%d OR t.qr_code_data ILIKE $%d)", argIndex, argIndex))
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

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM tickets t
		%s
		WHERE %s`, ticketJoinClause, whereClause)

	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count tickets: %w", err)
	}

	orderBy := "t.created_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("t.%s %s", filter.SortBy, direction)
	}

	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT %s
		FROM tickets t
		%s
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		ticketSelectColumns, ticketJoinClause, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, filter.Limit, offset)

	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}

	return tickets, pagination, nil
}

// Exists checks whether a ticket with the given ID exists.
func (r *ticketRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1)`

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket existence: %w", err)
	}

	return exists, nil
}

// ExistsByCode checks whether a ticket with the given serial_number exists.
func (r *ticketRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE serial_number = $1)`

	err := r.db.GetContext(ctx, &exists, query, code)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket serial number existence: %w", err)
	}

	return exists, nil
}

// ExistsByQRCode checks whether a ticket with the given QR code data exists.
func (r *ticketRepository) ExistsByQRCode(ctx context.Context, qrCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tickets WHERE qr_code_data = $1)`

	err := r.db.GetContext(ctx, &exists, query, qrCode)
	if err != nil {
		return false, fmt.Errorf("failed to check ticket QR code existence: %w", err)
	}

	return exists, nil
}

// ExistsBySerialNumber is an alias for ExistsByCode.
func (r *ticketRepository) ExistsBySerialNumber(ctx context.Context, serialNumber string) (bool, error) {
	return r.ExistsByCode(ctx, serialNumber)
}

// UpdateStatus updates only the status field of a ticket.
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

// Redeem marks a ticket as redeemed. Only succeeds if the ticket is currently 'active'.
// This is an alias for MarkRedeemed used by some service layers.
func (r *ticketRepository) Redeem(ctx context.Context, ticketID uuid.UUID, redeemedBy string) error {
	now := time.Now()
	query := `
		UPDATE tickets
		SET status = 'redeemed', redeemed_at = $1, redeemed_by = $2, updated_at = $1
		WHERE id = $3 AND status = 'active'`

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

// ValidateForEntry validates a ticket by its serial_number (code) for event entry.
// It performs all business-rule checks and returns a structured result.
func (r *ticketRepository) ValidateForEntry(ctx context.Context, code string) (*repositories.TicketValidationResult, error) {
	var result repositories.TicketValidationResult

	// JOIN through order_lines to reach ticket_tiers, events, and orders.
	query := `
		SELECT
			t.id as ticket_id,
			t.serial_number,
			t.status,
			o.code as order_code,
			o.status as order_status,
			tt.name as ticket_tier_name,
			e.id as event_id,
			e.name as event_title,
			e.event_date as event_event_date,
			e.end_date as event_end_date,
			e.venue_name,
			CASE
				WHEN t.status = 'redeemed'                              THEN false
				WHEN t.status != 'active'                               THEN false
				WHEN o.status != 'paid'                                 THEN false
				WHEN e.event_date > NOW() + INTERVAL '1 hour'           THEN false
				WHEN e.end_date IS NOT NULL AND e.end_date < NOW() - INTERVAL '1 hour' THEN false
				ELSE true
			END as is_valid,
			CASE
				WHEN t.status = 'redeemed'                              THEN 'Ticket already redeemed'
				WHEN t.status != 'active'                               THEN 'Ticket is not active'
				WHEN o.status != 'paid'                                 THEN 'Order not paid'
				WHEN e.event_date > NOW() + INTERVAL '1 hour'           THEN 'Event has not started yet'
				WHEN e.end_date IS NOT NULL AND e.end_date < NOW() - INTERVAL '1 hour' THEN 'Event has ended'
				ELSE 'Valid for entry'
			END as validation_message
		FROM tickets t
		JOIN order_lines ol ON t.order_line_id = ol.id
		JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
		JOIN events e ON tt.event_id = e.id
		JOIN orders o ON ol.order_id = o.id
		WHERE t.serial_number = $1
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

// GetStats retrieves aggregate ticket statistics with flexible filtering.
func (r *ticketRepository) GetStats(ctx context.Context, filter repositories.TicketFilter) (*repositories.TicketStats, error) {
	var stats repositories.TicketStats

	whereConditions := []string{"tt.is_active = true", "e.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1

	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("tt.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}

	if filter.TicketTierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.ticket_tier_id = $%d", argIndex))
		args = append(args, *filter.TicketTierID)
		argIndex++
	}

	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ol.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	query := fmt.Sprintf(`
		SELECT
			COALESCE(COUNT(*), 0) as total_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'active'    THEN 1 END), 0) as active_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'redeemed'  THEN 1 END), 0) as redeemed_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'cancelled' THEN 1 END), 0) as cancelled_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'refunded'  THEN 1 END), 0) as refunded_tickets,
			COALESCE(COUNT(CASE WHEN o.status = 'paid'      THEN 1 END), 0) as paid_tickets,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN tt.price ELSE 0 END), 0) as total_value,
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
		%s
		WHERE %s`, ticketJoinClause, whereClause)

	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket stats: %w", err)
	}

	return &stats, nil
}
