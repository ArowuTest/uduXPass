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

type paymentRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewPaymentRepository(db *sqlx.DB) repositories.PaymentRepository {
	return &paymentRepository{db: db}
}

func NewPaymentRepositoryWithTx(tx *sqlx.Tx) repositories.PaymentRepository {
	return &paymentRepository{db: tx}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	query := `
		INSERT INTO payments (
			id, order_id, provider, provider_transaction_id, 
			amount, currency, status, 
			webhook_received_at, created_at, updated_at
		) VALUES (
			:id, :order_id, :provider, :provider_transaction_id,
			:amount, :currency, :status,
			:webhook_received_at, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, payment)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "reference") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "provider_transaction_id") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
			}
		}
		return fmt.Errorf("failed to create payment: %w", err)
	}
	
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	query := `
		SELECT p.id, p.order_id, p.provider, p.provider_transaction_id,
			   p.reference, p.amount, p.currency, p.status, p.payment_method,
			   p.provider_response, p.webhook_data, p.failure_reason,
			   p.processed_at, p.meta_info, p.created_at, p.updated_at,
			   o.code as order_code, o.email as order_email, o.total_amount as order_total,
			   e.name as event_title, e.slug as event_slug
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE p.id = $1 AND o.is_active = true AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}
	
	return &payment, nil
}

func (r *paymentRepository) GetByReference(ctx context.Context, reference string) (*entities.Payment, error) {
	var payment entities.Payment
	query := `
		SELECT p.id, p.order_id, p.provider, p.provider_transaction_id,
			   p.reference, p.amount, p.currency, p.status, p.payment_method,
			   p.provider_response, p.webhook_data, p.failure_reason,
			   p.processed_at, p.meta_info, p.created_at, p.updated_at,
			   o.code as order_code, o.email as order_email, o.total_amount as order_total,
			   e.name as event_title, e.slug as event_slug
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE p.reference = $1 AND o.is_active = true AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &payment, query, reference)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by reference: %w", err)
	}
	
	return &payment, nil
}

func (r *paymentRepository) GetByProviderTransactionID(ctx context.Context, provider entities.PaymentMethod, transactionID string) (*entities.Payment, error) {
	var payment entities.Payment
	query := `
		SELECT p.id, p.order_id, p.provider, p.provider_transaction_id,
			   p.reference, p.amount, p.currency, p.status, p.payment_method,
			   p.provider_response, p.webhook_data, p.failure_reason,
			   p.processed_at, p.meta_info, p.created_at, p.updated_at,
			   o.code as order_code, o.email as order_email, o.total_amount as order_total,
			   e.name as event_title, e.slug as event_slug
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE p.provider = $1 AND p.provider_transaction_id = $2 AND o.is_active = true AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &payment, query, provider, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by provider transaction ID: %w", err)
	}
	
	return &payment, nil
}

func (r *paymentRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	query := `
		SELECT p.id, p.order_id, p.provider, p.provider_transaction_id,
			   p.reference, p.amount, p.currency, p.status, p.payment_method,
			   p.provider_response, p.webhook_data, p.failure_reason,
			   p.processed_at, p.meta_info, p.created_at, p.updated_at,
			   o.code as order_code, o.email as order_email, o.total_amount as order_total,
			   e.name as event_title, e.slug as event_slug
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE p.order_id = $1 AND o.is_active = true AND e.is_active = true
		ORDER BY p.created_at DESC`
	
	err := r.db.SelectContext(ctx, &payments, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by order: %w", err)
	}
	
	return payments, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	payment.UpdatedAt = time.Now()
	
	query := `
		UPDATE payments SET
			order_id = :order_id,
			provider = :provider,
			provider_transaction_id = :provider_transaction_id,
			reference = :reference,
			amount = :amount,
			currency = :currency,
			status = :status,
			payment_method = :payment_method,
			provider_response = :provider_response,
			webhook_data = :webhook_data,
			failure_reason = :failure_reason,
			processed_at = :processed_at,
			meta_info = :meta_info,
			updated_at = :updated_at
		WHERE id = :id`
	
	result, err := r.db.NamedExecContext(ctx, query, payment)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "reference") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "provider_transaction_id") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				if strings.Contains(pqErr.Detail, "order_id") {
					return entities.ErrOrderNotFound
				}
			}
		}
		return fmt.Errorf("failed to update payment: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	
	return nil
}

func (r *paymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payments WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	
	return nil
}

func (r *paymentRepository) List(ctx context.Context, filter repositories.PaymentFilter) ([]*entities.Payment, *repositories.PaginationResult, error) {
	var payments []*entities.Payment
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"o.is_active = true", "e.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrderID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.order_id = $%d", argIndex))
		args = append(args, *filter.OrderID)
		argIndex++
	}
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.Provider != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.provider = $%d", argIndex))
		args = append(args, *filter.Provider)
		argIndex++
	}
	
	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	
	if filter.PaymentMethod != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.payment_method = $%d", argIndex))
		args = append(args, *filter.PaymentMethod)
		argIndex++
	}
	
	if filter.Currency != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.currency = $%d", argIndex))
		args = append(args, filter.Currency)
		argIndex++
	}
	
	if filter.MinAmount != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.amount >= $%d", argIndex))
		args = append(args, *filter.MinAmount)
		argIndex++
	}
	
	if filter.MaxAmount != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.amount <= $%d", argIndex))
		args = append(args, *filter.MaxAmount)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(p.reference ILIKE $%d OR p.provider_transaction_id ILIKE $%d OR o.email ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	if filter.ProcessedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.processed_at >= $%d", argIndex))
		args = append(args, *filter.ProcessedFrom)
		argIndex++
	}
	
	if filter.ProcessedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.processed_at <= $%d", argIndex))
		args = append(args, *filter.ProcessedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count payments: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "p.created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("p.%s %s", filter.SortBy, direction)
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT p.id, p.order_id, p.provider, p.provider_transaction_id,
			   p.reference, p.amount, p.currency, p.status, p.payment_method,
			   p.provider_response, p.webhook_data, p.failure_reason,
			   p.processed_at, p.meta_info, p.created_at, p.updated_at,
			   o.code as order_code, o.email as order_email, o.total_amount as order_total,
			   e.name as event_title, e.slug as event_slug
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list payments: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return payments, pagination, nil
}

func (r *paymentRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM payments WHERE id = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check payment existence: %w", err)
	}
	
	return exists, nil
}

func (r *paymentRepository) ExistsByReference(ctx context.Context, reference string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM payments WHERE reference = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, reference)
	if err != nil {
		return false, fmt.Errorf("failed to check payment reference existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByProviderTransactionID checks if a payment exists by provider transaction ID
func (r *paymentRepository) ExistsByProviderTransactionID(ctx context.Context, provider entities.PaymentMethod, transactionID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM payments WHERE provider = $1 AND provider_transaction_id = $2)`
	
	err := r.db.GetContext(ctx, &exists, query, provider, transactionID)
	if err != nil {
		return false, fmt.Errorf("failed to check payment provider transaction ID existence: %w", err)
	}
	
	return exists, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entities.PaymentStatus) error {
	now := time.Now()
	query := `
		UPDATE payments 
		SET status = $1, processed_at = $2, updated_at = $2
		WHERE id = $3`
	
	result, err := r.db.ExecContext(ctx, query, status, now, paymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	
	return nil
}

// MarkWebhookReceived marks when a webhook was received
func (r *paymentRepository) MarkWebhookReceived(ctx context.Context, paymentID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE payments 
		SET updated_at = $1
		WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, now, paymentID)
	if err != nil {
		return fmt.Errorf("failed to mark webhook received: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	
	return nil
}

func (r *paymentRepository) UpdateWebhookData(ctx context.Context, paymentID uuid.UUID, webhookData map[string]interface{}) error {
	query := `
		UPDATE payments 
		SET webhook_data = $1, updated_at = NOW()
		WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, webhookData, paymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment webhook data: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	
	return nil
}

func (r *paymentRepository) GetSuccessful(ctx context.Context, filter repositories.PaymentFilter) ([]*entities.Payment, *repositories.PaginationResult, error) {
	status := entities.PaymentStatusCompleted
	filter.Status = &status
	return r.List(ctx, filter)
}

func (r *paymentRepository) GetFailed(ctx context.Context, filter repositories.PaymentFilter) ([]*entities.Payment, *repositories.PaginationResult, error) {
	status := entities.PaymentStatusFailed
	filter.Status = &status
	return r.List(ctx, filter)
}

func (r *paymentRepository) GetPending(ctx context.Context, filter repositories.PaymentFilter) ([]*entities.Payment, *repositories.PaginationResult, error) {
	status := entities.PaymentStatusPending
	filter.Status = &status
	return r.List(ctx, filter)
}

func (r *paymentRepository) GetByProvider(ctx context.Context, provider entities.PaymentMethod, filter repositories.PaymentFilter) ([]*entities.Payment, *repositories.PaginationResult, error) {
	filter.Provider = &provider
	return r.List(ctx, filter)
}

// GetPaymentStats retrieves payment statistics
func (r *paymentRepository) GetPaymentStats(ctx context.Context, filter repositories.PaymentStatsFilter) (*repositories.PaymentStats, error) {
	var stats repositories.PaymentStats
	
	// Build WHERE clause for stats
	whereConditions := []string{"o.is_active = true", "e.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}
	
	if filter.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_payments,
			COUNT(CASE WHEN p.status = 'completed' THEN 1 END) as successful_payments,
			COUNT(CASE WHEN p.status = 'failed' THEN 1 END) as failed_payments,
			COALESCE(SUM(CASE WHEN p.status = 'completed' THEN p.amount ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(CASE WHEN p.status = 'completed' THEN p.amount END), 0) as average_payment_amount
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment stats: %w", err)
	}
	
	return &stats, nil
}

func (r *paymentRepository) GetStats(ctx context.Context, filter repositories.PaymentFilter) (*repositories.PaymentStats, error) {
	var stats repositories.PaymentStats
	
	// Build WHERE clause for stats
	whereConditions := []string{"o.is_active = true", "e.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.EventID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("o.event_id = $%d", argIndex))
		args = append(args, *filter.EventID)
		argIndex++
	}
	
	if filter.Provider != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.provider = $%d", argIndex))
		args = append(args, *filter.Provider)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			COALESCE(COUNT(*), 0) as total_payments,
			COALESCE(COUNT(CASE WHEN p.status = 'success' THEN 1 END), 0) as successful_payments,
			COALESCE(COUNT(CASE WHEN p.status = 'failed' THEN 1 END), 0) as failed_payments,
			COALESCE(COUNT(CASE WHEN p.status = 'pending' THEN 1 END), 0) as pending_payments,
			COALESCE(COUNT(CASE WHEN p.status = 'cancelled' THEN 1 END), 0) as cancelled_payments,
			COALESCE(SUM(CASE WHEN p.status = 'success' THEN p.amount ELSE 0 END), 0) as total_successful_amount,
			COALESCE(SUM(p.amount), 0) as total_attempted_amount,
			COALESCE(AVG(CASE WHEN p.status = 'success' THEN p.amount END), 0) as avg_successful_amount,
			COALESCE(MIN(CASE WHEN p.status = 'success' THEN p.amount END), 0) as min_successful_amount,
			COALESCE(MAX(CASE WHEN p.status = 'success' THEN p.amount END), 0) as max_successful_amount,
			COALESCE(COUNT(DISTINCT p.provider), 0) as unique_providers,
			COALESCE(COUNT(DISTINCT p.currency), 0) as unique_currencies,
			MIN(CASE WHEN p.status = 'success' THEN p.processed_at END) as first_successful_payment_at,
			MAX(CASE WHEN p.status = 'success' THEN p.processed_at END) as last_successful_payment_at,
			ROUND(
				CASE 
					WHEN COUNT(*) > 0 THEN 
						(COUNT(CASE WHEN p.status = 'success' THEN 1 END)::float / COUNT(*)::float) * 100
					ELSE 0 
				END, 2
			) as success_rate
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN events e ON o.event_id = e.id
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment stats: %w", err)
	}
	
	return &stats, nil
}

