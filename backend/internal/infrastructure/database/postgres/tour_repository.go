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

type tourRepository struct {
	db interface {
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewTourRepository(db *sqlx.DB) repositories.TourRepository {
	return &tourRepository{db: db}
}

func NewTourRepositoryWithTx(tx *sqlx.Tx) repositories.TourRepository {
	return &tourRepository{db: tx}
}

func (r *tourRepository) Create(ctx context.Context, tour *entities.Tour) error {
	query := `
		INSERT INTO tours (
			id, organizer_id, name, slug, description, event_date, 
			end_date, is_active, image_url, meta_info, 
			created_at, updated_at
		) VALUES (
			:id, :organizer_id, :name, :slug, :description, :event_date,
			:end_date, :is_active, :image_url, :meta_info,
			:created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, tour)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				return entities.ErrNotFoundError
			}
		}
		return fmt.Errorf("failed to create tour: %w", err)
	}
	
	return nil
}

func (r *tourRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Tour, error) {
	var tour entities.Tour
	query := `
		SELECT t.id, t.organizer_id, t.name, t.slug, t.description, 
			   t.event_date, t.end_date, t.is_active, t.image_url, 
			   t.meta_info, t.created_at, t.updated_at,
			   o.name as organizer_name, o.slug as organizer_slug
		FROM tours t
		JOIN organizers o ON t.organizer_id = o.id
		WHERE t.id = $1 AND t.is_active = true AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &tour, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get tour by ID: %w", err)
	}
	
	return &tour, nil
}

func (r *tourRepository) GetBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (*entities.Tour, error) {
	var tour entities.Tour
	query := `
		SELECT t.id, t.organizer_id, t.name, t.slug, t.description, 
			   t.event_date, t.end_date, t.is_active, t.image_url, 
			   t.meta_info, t.created_at, t.updated_at,
			   o.name as organizer_name, o.slug as organizer_slug
		FROM tours t
		JOIN organizers o ON t.organizer_id = o.id
		WHERE t.organizer_id = $1 AND t.slug = $2 AND t.is_active = true AND o.is_active = true`
	
	err := r.db.GetContext(ctx, &tour, query, organizerID, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get tour by slug: %w", err)
	}
	
	return &tour, nil
}

func (r *tourRepository) GetByOrganizer(ctx context.Context, organizerID uuid.UUID, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	filter.OrganizerID = &organizerID
	return r.List(ctx, filter)
}

func (r *tourRepository) Update(ctx context.Context, tour *entities.Tour) error {
	tour.UpdatedAt = time.Now()
	
	query := `
		UPDATE tours SET
			organizer_id = :organizer_id,
			name = :name,
			slug = :slug,
			description = :description,
			event_date = :event_date,
			end_date = :end_date,
			is_active = :is_active,
			image_url = :image_url,
			meta_info = :meta_info,
			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, tour)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
			case "23503": // foreign_key_violation
				return entities.ErrNotFoundError
			}
		}
		return fmt.Errorf("failed to update tour: %w", err)
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

func (r *tourRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tours SET is_active = false WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tour: %w", err)
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

func (r *tourRepository) List(ctx context.Context, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	var tours []*entities.Tour
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"t.is_active = true", "o.is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrganizerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.organizer_id = $%d", argIndex))
		args = append(args, *filter.OrganizerID)
		argIndex++
	}
	
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(t.name ILIKE $%d OR t.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}
	
	if filter.StartDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.event_date >= $%d", argIndex))
		args = append(args, *filter.StartDateFrom)
		argIndex++
	}
	
	if filter.StartDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.event_date <= $%d", argIndex))
		args = append(args, *filter.StartDateTo)
		argIndex++
	}
	
	if filter.EndDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.end_date >= $%d", argIndex))
		args = append(args, *filter.EndDateFrom)
		argIndex++
	}
	
	if filter.EndDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("t.end_date <= $%d", argIndex))
		args = append(args, *filter.EndDateTo)
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
	
	whereClause := strings.Join(whereConditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM tours t 
		JOIN organizers o ON t.organizer_id = o.id 
		WHERE %s`, whereClause)
	
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count tours: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "t.event_date ASC"
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
		SELECT t.id, t.organizer_id, t.name, t.slug, t.description, 
			   t.event_date, t.end_date, t.is_active, t.image_url, 
			   t.meta_info, t.created_at, t.updated_at,
			   o.name as organizer_name, o.slug as organizer_slug,
			   COUNT(e.id) as event_count
		FROM tours t
		JOIN organizers o ON t.organizer_id = o.id
		LEFT JOIN events e ON t.id = e.tour_id AND e.is_active = true
		WHERE %s 
		GROUP BY t.id, t.organizer_id, t.name, t.slug, t.description, 
				 t.event_date, t.end_date, t.is_active, t.image_url, 
				 t.meta_info, t.created_at, t.updated_at,
				 o.name, o.slug
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &tours, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tours: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return tours, pagination, nil
}

func (r *tourRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tours WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check tour existence: %w", err)
	}
	
	return exists, nil
}

func (r *tourRepository) ExistsBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tours WHERE organizer_id = $1 AND slug = $2 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, organizerID, slug)
	if err != nil {
		return false, fmt.Errorf("failed to check tour slug existence: %w", err)
	}
	
	return exists, nil
}

func (r *tourRepository) GetActive(ctx context.Context, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	isActive := true
	filter.IsActive = &isActive
	return r.List(ctx, filter)
}

// GetOngoing retrieves ongoing tours (currently happening)
func (r *tourRepository) GetOngoing(ctx context.Context, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	now := time.Now()
	filter.StartDateTo = &now
	filter.EndDateFrom = &now
	isActive := true
	filter.IsActive = &isActive
	return r.List(ctx, filter)
}

// GetTourStats retrieves statistics for a specific tour
func (r *tourRepository) GetTourStats(ctx context.Context, tourID uuid.UUID) (*repositories.TourStats, error) {
	var stats repositories.TourStats
	
	query := `
		SELECT 
			t.id as tour_id,
			t.name as tour_name,
			COUNT(e.id) as total_events,
			COUNT(CASE WHEN e.status = 'published' THEN 1 END) as published_events,
			COUNT(CASE WHEN e.status = 'published' AND e.event_date > NOW() THEN 1 END) as on_sale_events,
			COALESCE(SUM(CASE WHEN o.status IN ('confirmed', 'paid') THEN ol.total_amount ELSE 0 END), 0) as total_revenue,
			COUNT(DISTINCT CASE WHEN o.status IN ('confirmed', 'paid') THEN o.id END) as total_orders,
			COUNT(DISTINCT CASE WHEN tk.status = 'redeemed' THEN tk.id END) as redeemed_tickets
		FROM tours t
		LEFT JOIN events e ON t.id = e.tour_id AND e.is_active = true
		LEFT JOIN ticket_tiers tt ON e.id = tt.event_id AND tt.is_active = true
		LEFT JOIN order_lines ol ON tt.id = ol.ticket_tier_id
		LEFT JOIN orders o ON ol.order_id = o.id AND o.is_active = true
		LEFT JOIN tickets tk ON ol.id = tk.order_line_id
		WHERE t.id = $1 AND t.is_active = true
		GROUP BY t.id, t.name`
	
	err := r.db.GetContext(ctx, &stats, query, tourID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get tour stats: %w", err)
	}
	
	return &stats, nil
}

func (r *tourRepository) GetUpcoming(ctx context.Context, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	now := time.Now()
	filter.StartDateFrom = &now
	return r.List(ctx, filter)
}

func (r *tourRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, filter repositories.TourFilter) ([]*entities.Tour, *repositories.PaginationResult, error) {
	filter.StartDateFrom = &startDate
	filter.EndDateTo = &endDate
	return r.List(ctx, filter)
}

func (r *tourRepository) GetStats(ctx context.Context, tourID uuid.UUID) (*repositories.TourStats, error) {
	var stats repositories.TourStats
	
	query := `
		SELECT 
			t.id as tour_id,
			t.name as tour_name,
			COALESCE(COUNT(DISTINCT e.id), 0) as total_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'published' THEN e.id END), 0) as published_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'draft' THEN e.id END), 0) as draft_events,
			COALESCE(COUNT(DISTINCT CASE WHEN e.status = 'cancelled' THEN e.id END), 0) as cancelled_events,
			COALESCE(COUNT(DISTINCT o.id), 0) as total_orders,
			COALESCE(COUNT(DISTINCT CASE WHEN o.status = 'paid' THEN o.id END), 0) as paid_orders,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN o.total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(COUNT(DISTINCT tk.id), 0) as total_tickets_sold,
			COALESCE(COUNT(DISTINCT CASE WHEN tk.status = 'redeemed' THEN tk.id END), 0) as tickets_redeemed,
			MIN(CASE WHEN e.status = 'published' THEN e.event_date END) as first_event_date,
			MAX(CASE WHEN e.status = 'published' THEN e.event_date END) as last_event_date,
			MIN(CASE WHEN o.status = 'paid' THEN o.created_at END) as first_sale_at,
			MAX(CASE WHEN o.status = 'paid' THEN o.created_at END) as last_sale_at
		FROM tours t
		LEFT JOIN events e ON t.id = e.tour_id AND e.is_active = true
		LEFT JOIN orders o ON e.id = o.event_id AND o.is_active = true
		LEFT JOIN tickets tk ON o.id = tk.order_id
		WHERE t.id = $1 AND t.is_active = true
		GROUP BY t.id, t.name`
	
	err := r.db.GetContext(ctx, &stats, query, tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tour stats: %w", err)
	}
	
	return &stats, nil
}

