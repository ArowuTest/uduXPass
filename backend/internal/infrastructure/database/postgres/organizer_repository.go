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

type organizerRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewOrganizerRepository(db *sqlx.DB) repositories.OrganizerRepository {
	return &organizerRepository{db: db}
}

func NewOrganizerRepositoryWithTx(tx *sqlx.Tx) repositories.OrganizerRepository {
	return &organizerRepository{db: tx}
}

// selectColumns returns the columns that exist in both the entity struct and the database table
// Actual table columns: id, name, slug, email, phone, website_url, logo_url, description, address, city, state, country, is_active, settings, created_at, updated_at
const organizerSelectColumns = `id, name, slug, description, website_url AS website, email, phone, logo_url, is_active, created_at, updated_at`

func (r *organizerRepository) Create(ctx context.Context, organizer *entities.Organizer) error {
	query := `
		INSERT INTO organizers (
			id, name, slug, description, website_url, email, phone, 
			logo_url, is_active, 
			created_at, updated_at
		) VALUES (
			:id, :name, :slug, :description, :website, :email, :phone,
			:logo_url, :is_active,
			:created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, organizer)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "email") {
					return entities.ErrConflictError
				}
			}
		}
		return fmt.Errorf("failed to create organizer: %w", err)
	}
	
	return nil
}

func (r *organizerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Organizer, error) {
	var organizer entities.Organizer
	query := fmt.Sprintf(`
		SELECT %s
		FROM organizers 
		WHERE id = $1 AND is_active = true`, organizerSelectColumns)
	
	err := r.db.GetContext(ctx, &organizer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get organizer by ID: %w", err)
	}
	
	return &organizer, nil
}

func (r *organizerRepository) GetBySlug(ctx context.Context, slug string) (*entities.Organizer, error) {
	var organizer entities.Organizer
	query := fmt.Sprintf(`
		SELECT %s
		FROM organizers 
		WHERE slug = $1 AND is_active = true`, organizerSelectColumns)
	
	err := r.db.GetContext(ctx, &organizer, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get organizer by slug: %w", err)
	}
	
	return &organizer, nil
}

func (r *organizerRepository) GetByEmail(ctx context.Context, email string) (*entities.Organizer, error) {
	var organizer entities.Organizer
	query := fmt.Sprintf(`
		SELECT %s
		FROM organizers 
		WHERE email = $1 AND is_active = true`, organizerSelectColumns)
	
	err := r.db.GetContext(ctx, &organizer, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get organizer by email: %w", err)
	}
	
	return &organizer, nil
}

func (r *organizerRepository) Update(ctx context.Context, organizer *entities.Organizer) error {
	organizer.UpdatedAt = time.Now()
	
	query := `
		UPDATE organizers SET
			name = :name,
			slug = :slug,
			description = :description,
			website_url = :website,
			email = :email,
			phone = :phone,
			logo_url = :logo_url,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, organizer)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
				if strings.Contains(pqErr.Detail, "email") {
					return entities.ErrConflictError
				}
			}
		}
		return fmt.Errorf("failed to update organizer: %w", err)
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

func (r *organizerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE organizers SET is_active = false WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organizer: %w", err)
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

func (r *organizerRepository) List(ctx context.Context, filter repositories.OrganizerFilter) ([]*entities.Organizer, *repositories.PaginationResult, error) {
	var organizers []*entities.Organizer
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	if filter.Country != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, filter.Country)
		argIndex++
	}
	
	if filter.City != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, filter.City)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}
	
	if filter.CreatedFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedFrom)
		argIndex++
	}
	
	if filter.CreatedTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedTo)
		argIndex++
	}
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM organizers WHERE %s", whereClause)
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count organizers: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.SortBy, direction)
	}
	
	// Handle pagination defaults
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT %s
		FROM organizers 
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, organizerSelectColumns, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &organizers, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list organizers: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return organizers, pagination, nil
}

func (r *organizerRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizers WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check organizer existence: %w", err)
	}
	
	return exists, nil
}

func (r *organizerRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizers WHERE slug = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, slug)
	if err != nil {
		return false, fmt.Errorf("failed to check organizer slug existence: %w", err)
	}
	
	return exists, nil
}

func (r *organizerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizers WHERE email = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check organizer email existence: %w", err)
	}
	
	return exists, nil
}

func (r *organizerRepository) GetStats(ctx context.Context, organizerID uuid.UUID) (*repositories.OrganizerStats, error) {
	var stats repositories.OrganizerStats
	
	query := `
		SELECT 
			$1 as organizer_id,
			COALESCE(COUNT(DISTINCT e.id), 0) as total_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'published' THEN e.id END), 0) as published_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'draft' THEN e.id END), 0) as draft_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'cancelled' THEN e.id END), 0) as cancelled_events,
			COALESCE(COUNT(DISTINCT o.id), 0) as total_orders,
			COALESCE(COUNT(DISTINCT CASE WHEN o.status = 'paid' THEN o.id END), 0) as paid_orders,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN o.total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(COUNT(DISTINCT t.id), 0) as total_tickets_sold,
			COALESCE(COUNT(DISTINCT CASE WHEN t.status = 'redeemed' THEN t.id END), 0) as tickets_redeemed,
			MIN(CASE WHEN e.status = 'published' THEN e.created_at END) as first_event_at,
			MAX(CASE WHEN e.status = 'published' THEN e.created_at END) as last_event_at
		FROM organizers org
		LEFT JOIN events e ON org.id = e.organizer_id AND e.is_active = true
		LEFT JOIN orders o ON e.id = o.event_id AND o.is_active = true
		LEFT JOIN tickets t ON o.id = t.order_id
		WHERE org.id = $1 AND org.is_active = true
		GROUP BY org.id`
	
	err := r.db.GetContext(ctx, &stats, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer stats: %w", err)
	}
	
	return &stats, nil
}
