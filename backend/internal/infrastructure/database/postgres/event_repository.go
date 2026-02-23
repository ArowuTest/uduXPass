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

type eventRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewEventRepository(db *sqlx.DB) repositories.EventRepository {
	return &eventRepository{db: db}
}

func NewEventRepositoryWithTx(tx *sqlx.Tx) repositories.EventRepository {
	return &eventRepository{db: tx}
}

func (r *eventRepository) Create(ctx context.Context, event *entities.Event) error {
		query := `
			INSERT INTO events (
				id, organizer_id, category_id, name, slug, description, 
				event_date, doors_open, venue_name, venue_address, 
				venue_city, venue_state, venue_country, venue_capacity, 
				event_image_url, status, sale_start, sale_end, 
				settings, is_active, created_at, updated_at
			) VALUES (
				:id, :organizer_id, :category_id, :name, :slug, :description,
				:event_date, :doors_open, :venue_name, :venue_address,
				:venue_city, :venue_state, :venue_country, :venue_capacity,
				:event_image_url, :status, :sale_start, :sale_end,
				:settings, :is_active, :created_at, :updated_at
			)`
	
	_, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
			}
		}
		return fmt.Errorf("failed to create event: %w", err)
	}
	
	return nil
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Event, error) {
	var event entities.Event
	query := `
		SELECT e.id, e.organizer_id, e.category_id, e.name, e.slug, e.description,
			   e.event_date, e.doors_open, e.venue_name, e.venue_address, 
			   e.venue_city, e.venue_state, e.venue_country, e.venue_capacity,
			   e.event_image_url, e.status, e.sale_start, e.sale_end, 
			   e.settings, e.created_at, e.updated_at, e.is_active
		FROM events e
		WHERE e.id = $1 AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	
	return &event, nil
}

func (r *eventRepository) GetBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (*entities.Event, error) {
	var event entities.Event
	query := `
		SELECT e.id, e.organizer_id, e.category_id, e.name, e.slug, e.description,
			   e.event_date, e.doors_open, e.venue_name, e.venue_address, 
			   e.venue_city, e.venue_state, e.venue_country, e.venue_capacity,
			   e.event_image_url, e.status, e.sale_start, e.sale_end, 
			   e.settings, e.created_at, e.updated_at, e.is_active
		FROM events e
		WHERE e.organizer_id = $1 AND e.slug = $2 AND e.is_active = true`
	
	err := r.db.GetContext(ctx, &event, query, organizerID, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get event by slug: %w", err)
	}
	
	return &event, nil
}

func (r *eventRepository) buildEventQuery(baseQuery string, filter repositories.EventFilter) (string, []interface{}) {
	query := baseQuery + " WHERE e.is_active = true"
	args := []interface{}{}
	argIndex := 1
	
	if filter.OrganizerID != nil {
		query += fmt.Sprintf(" AND e.organizer_id = $%d", argIndex)
		args = append(args, *filter.OrganizerID)
		argIndex++
	}
	
	if filter.TourID != nil {
		query += fmt.Sprintf(" AND e.tour_id = $%d", argIndex)
		args = append(args, *filter.TourID)
		argIndex++
	}
	
	if filter.Status != nil {
		query += fmt.Sprintf(" AND e.status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}
	
	if filter.City != "" {
		query += fmt.Sprintf(" AND e.venue_city ILIKE $%d", argIndex)
		args = append(args, "%"+filter.City+"%")
		argIndex++
	}
	
	if filter.Country != "" {
		query += fmt.Sprintf(" AND e.venue_country ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Country+"%")
		argIndex++
	}
	
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (e.name ILIKE $%d OR e.description ILIKE $%d OR e.venue_name ILIKE $%d)", argIndex, argIndex, argIndex)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}
	
	if filter.EventDateFrom != nil {
		query += fmt.Sprintf(" AND e.event_date >= $%d", argIndex)
		args = append(args, *filter.EventDateFrom)
		argIndex++
	}
	
	if filter.EventDateTo != nil {
		query += fmt.Sprintf(" AND e.event_date <= $%d", argIndex)
		args = append(args, *filter.EventDateTo)
		argIndex++
	}
	
	return query, args
}

func (r *eventRepository) List(ctx context.Context, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	var events []*entities.Event
	
	baseQuery := `
		SELECT e.id, e.organizer_id, e.category_id, e.name, e.slug, e.description,
			   e.event_date, e.doors_open, e.venue_name, e.venue_address, 
			   e.venue_city, e.venue_state, e.venue_country, e.venue_capacity,
			   e.event_image_url, e.status, e.sale_start, e.sale_end, 
			   e.settings, e.created_at, e.updated_at, e.is_active
		FROM events e`
	
	query, args := r.buildEventQuery(baseQuery, filter)
	
	// Add ordering
	if filter.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY e.%s %s", filter.SortBy, filter.SortOrder)
	} else {
		query += " ORDER BY e.event_date DESC"
	}
	
	// Count total records
	countQuery := "SELECT COUNT(*) FROM events e"
	countQuery, countArgs := r.buildEventQuery(countQuery, filter)
	
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count events: %w", err)
	}
	
	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, filter.Limit)
		
		offset := (filter.Page - 1) * filter.Limit
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
			args = append(args, offset)
		}
	}
	
	err = r.db.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list events: %w", err)
	}
	
	pagination := repositories.NewPaginationResult(filter.Page, filter.Limit, total)
	return events, pagination, nil
}

func (r *eventRepository) ListPublic(ctx context.Context, filter repositories.PublicEventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	var events []*entities.Event
	
	query := `
		SELECT e.id, e.organizer_id, e.category_id, e.name, e.slug, e.description,
			   e.event_date, e.doors_open, e.venue_name, e.venue_address, 
			   e.venue_city, e.venue_state, e.venue_country, e.venue_capacity,
			   e.event_image_url, e.status, e.sale_start, e.sale_end, 
			   e.settings, e.created_at, e.updated_at, e.is_active
		FROM events e
		WHERE e.is_active = true AND e.status IN ('published', 'on_sale')`
	
	args := []interface{}{}
	argIndex := 1
	
	if filter.City != "" {
		query += fmt.Sprintf(" AND e.venue_city ILIKE $%d", argIndex)
		args = append(args, "%"+filter.City+"%")
		argIndex++
	}
	
	if filter.Country != "" {
		query += fmt.Sprintf(" AND e.venue_country ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Country+"%")
		argIndex++
	}
	
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (e.name ILIKE $%d OR e.description ILIKE $%d)", argIndex, argIndex)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}
	
	// Count total records with same filters
	countQuery := "SELECT COUNT(*) FROM events e WHERE e.is_active = true AND e.status IN ('published', 'on_sale')"
	countArgs := []interface{}{}
	countArgIndex := 1
	
	if filter.City != "" {
		countQuery += fmt.Sprintf(" AND e.venue_city ILIKE $%d", countArgIndex)
		countArgs = append(countArgs, "%"+filter.City+"%")
		countArgIndex++
	}
	
	if filter.Country != "" {
		countQuery += fmt.Sprintf(" AND e.venue_country ILIKE $%d", countArgIndex)
		countArgs = append(countArgs, "%"+filter.Country+"%")
		countArgIndex++
	}
	
	if filter.Search != "" {
		countQuery += fmt.Sprintf(" AND (e.name ILIKE $%d OR e.description ILIKE $%d)", countArgIndex, countArgIndex)
		searchTerm := "%" + filter.Search + "%"
		countArgs = append(countArgs, searchTerm)
		countArgIndex++
	}
	
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count public events: %w", err)
	}
	
	// Add ordering and pagination
	query += " ORDER BY e.event_date ASC"
	
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, filter.Limit)
		
		offset := (filter.Page - 1) * filter.Limit
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
			args = append(args, offset)
		}
	}
	
	err = r.db.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list public events: %w", err)
	}
	
	pagination := repositories.NewPaginationResult(filter.Page, filter.Limit, total)
	return events, pagination, nil
}

func (r *eventRepository) GetByOrganizer(ctx context.Context, organizerID uuid.UUID, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	filter.OrganizerID = &organizerID
	return r.List(ctx, filter)
}

func (r *eventRepository) GetByTour(ctx context.Context, tourID uuid.UUID, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	filter.TourID = &tourID
	return r.List(ctx, filter)
}

func (r *eventRepository) GetUpcoming(ctx context.Context, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	now := time.Now()
	filter.EventDateFrom = &now
	return r.List(ctx, filter)
}

func (r *eventRepository) GetByCity(ctx context.Context, city string, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	filter.City = city
	return r.List(ctx, filter)
}

func (r *eventRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, filter repositories.EventFilter) ([]*entities.Event, *repositories.PaginationResult, error) {
	filter.EventDateFrom = &startDate
	filter.EventDateTo = &endDate
	return r.List(ctx, filter)
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) error {
	query := `
		UPDATE events SET
			name = :name,
			slug = :slug,
			description = :description,
			event_date = :event_date,
			doors_open = :doors_open,
			venue_name = :venue_name,
			venue_address = :venue_address,
			venue_city = :venue_city,
			venue_state = :venue_state,
			venue_country = :venue_country,
			venue_capacity = :venue_capacity,
			event_image_url = :event_image_url,
			status = :status,
			sale_start = :sale_start,
			sale_end = :sale_end,
			settings = :settings,

			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return entities.ErrConflictError
				}
			}
		}
		return fmt.Errorf("failed to update event: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrEventNotFound
	}
	
	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE events SET is_active = false, updated_at = NOW() WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrEventNotFound
	}
	
	return nil
}

func (r *eventRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if event exists: %w", err)
	}
	
	return exists, nil
}

func (r *eventRepository) ExistsBySlug(ctx context.Context, organizerID uuid.UUID, slug string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE organizer_id = $1 AND slug = $2 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, organizerID, slug)
	if err != nil {
		return false, fmt.Errorf("failed to check if event exists by slug: %w", err)
	}
	
	return exists, nil
}

func (r *eventRepository) GetEventStats(ctx context.Context, eventID uuid.UUID) (*repositories.EventStats, error) {
	// Implementation for event statistics
	return &repositories.EventStats{}, nil
}

func (r *eventRepository) UpdateStatus(ctx context.Context, eventID uuid.UUID, status entities.EventStatus) error {
	query := `UPDATE events SET status = $1, updated_at = NOW() WHERE id = $2 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, status, eventID)
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrEventNotFound
	}
	
	return nil
}

