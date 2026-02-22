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

type scannerUserRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}
}

// NewScannerUserRepository creates a new scanner user repository
func NewScannerUserRepository(db *sqlx.DB) repositories.ScannerUserRepository {
	return &scannerUserRepository{db: db}
}

// NewScannerUserRepositoryWithTx creates a new scanner user repository with a transaction
func NewScannerUserRepositoryWithTx(tx *sqlx.Tx) repositories.ScannerUserRepository {
	return &scannerUserRepository{db: tx}
}

// Authentication methods
func (r *scannerUserRepository) GetByUsername(ctx context.Context, username string) (*entities.ScannerUser, error) {
	query := `
		SELECT id, username, email, password_hash, name, role, status, 
			   created_at, last_login, created_by, login_attempts, locked_until, must_change_password
		FROM scanner_users 
		WHERE username = $1`
	
	var scanner entities.ScannerUser
	var name sql.NullString
	var lastLogin sql.NullTime
	var createdBy sql.NullString
	var lockedUntil sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&scanner.ID,
		&scanner.Username,
		&scanner.Email,
		&scanner.PasswordHash,
		&name,
		&scanner.Role,
		&scanner.Status,
		&scanner.CreatedAt,
		&lastLogin,
		&createdBy,
		&scanner.LoginAttempts,
		&lockedUntil,
		&scanner.MustChangePassword,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get scanner user by username: %w", err)
	}
	
	// Handle nullable fields
	if name.Valid {
		scanner.Name = name.String
	}
	if lastLogin.Valid {
		scanner.LastLogin = &lastLogin.Time
	}
	if createdBy.Valid {
		if createdByUUID, err := uuid.Parse(createdBy.String); err == nil {
			scanner.CreatedBy = &createdByUUID
		}
	}
	if lockedUntil.Valid {
		scanner.LockedUntil = &lockedUntil.Time
	}
	
	// Set UpdatedAt to CreatedAt since we don't have updated_at in the database
	scanner.UpdatedAt = scanner.CreatedAt
	
	return &scanner, nil
}

func (r *scannerUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ScannerUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, status, 
			   created_at, updated_at, last_login_at, created_by, updated_by
		FROM scanner_users 
		WHERE id = $1 AND is_active = true`
	
	var scanner entities.ScannerUser
	err := r.db.GetContext(ctx, &scanner, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get scanner user by ID: %w", err)
	}
	
	return &scanner, nil
}

func (r *scannerUserRepository) GetByEmail(ctx context.Context, email string) (*entities.ScannerUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, status, 
			   created_at, updated_at, last_login_at, created_by, updated_by
		FROM scanner_users 
		WHERE email = $1 AND is_active = true`
	
	var scanner entities.ScannerUser
	err := r.db.GetContext(ctx, &scanner, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get scanner user by email: %w", err)
	}
	
	return &scanner, nil
}

// CRUD operations
func (r *scannerUserRepository) Create(ctx context.Context, scanner *entities.ScannerUser) error {
	if scanner.ID == uuid.Nil {
		scanner.ID = uuid.New()
	}
	
	now := time.Now()
	scanner.CreatedAt = now
	scanner.UpdatedAt = now
	
	query := `
		INSERT INTO scanner_users (id, username, email, password_hash, first_name, last_name, 
								  role, status, created_at, updated_at, created_by, updated_by)
		VALUES (:id, :username, :email, :password_hash, :first_name, :last_name, 
				:role, :status, :created_at, :updated_at, :created_by, :updated_by)`
	
	_, err := r.db.NamedExecContext(ctx, query, scanner)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return entities.ErrConflictError
			}
		}
		return fmt.Errorf("failed to create scanner user: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) Update(ctx context.Context, scanner *entities.ScannerUser) error {
	scanner.UpdatedAt = time.Now()
	
	query := `
		UPDATE scanner_users 
		SET username = :username, email = :email, password_hash = :password_hash, 
			first_name = :first_name, last_name = :last_name, role = :role, 
			status = :status, updated_at = :updated_at, updated_by = :updated_by,
			last_login_at = :last_login_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, scanner)
	if err != nil {
		return fmt.Errorf("failed to update scanner user: %w", err)
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

func (r *scannerUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE scanner_users 
		SET is_active = false 
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scanner user: %w", err)
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

// List and search
func (r *scannerUserRepository) List(ctx context.Context, filter *repositories.ScannerUserFilter) ([]*entities.ScannerUser, *repositories.PaginationResult, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	conditions = append(conditions, "is_active = true")
	
	if filter != nil {
		if filter.Role != nil {
			conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
			args = append(args, *filter.Role)
			argIndex++
		}
		
		if filter.Status != nil {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}
		
		if filter.Search != "" {
			conditions = append(conditions, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
			args = append(args, "%"+filter.Search+"%")
			argIndex++
		}
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM scanner_users %s", whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count scanner users: %w", err)
	}
	
	// Build main query with pagination
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, first_name, last_name, role, status, 
			   created_at, updated_at, last_login_at, created_by, updated_by
		FROM scanner_users %s
		ORDER BY created_at DESC`, whereClause)
	
	if filter != nil && filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
		
		if filter.GetOffset() > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.GetOffset())
		}
	}
	
	var scanners []*entities.ScannerUser
	err = r.db.SelectContext(ctx, &scanners, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list scanner users: %w", err)
	}
	
	pagination := repositories.NewPaginationResult(filter.Page, filter.Limit, total)
	
	return scanners, pagination, nil
}

func (r *scannerUserRepository) Search(ctx context.Context, query string, filter *repositories.ScannerUserFilter) ([]*entities.ScannerUser, *repositories.PaginationResult, error) {
	if filter == nil {
		filter = &repositories.ScannerUserFilter{}
	}
	filter.Search = query
	return r.List(ctx, filter)
}

// Event assignments
func (r *scannerUserRepository) AssignToEvent(ctx context.Context, scannerID, eventID, assignedBy uuid.UUID) error {
	query := `
		INSERT INTO scanner_event_assignments (scanner_id, event_id, assigned_by, assigned_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (scanner_id, event_id) DO NOTHING`
	
	_, err := r.db.ExecContext(ctx, query, scannerID, eventID, assignedBy)
	if err != nil {
		return fmt.Errorf("failed to assign scanner to event: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) UnassignFromEvent(ctx context.Context, scannerID, eventID uuid.UUID) error {
	query := `DELETE FROM scanner_event_assignments WHERE scanner_id = $1 AND event_id = $2`
	
	_, err := r.db.ExecContext(ctx, query, scannerID, eventID)
	if err != nil {
		return fmt.Errorf("failed to unassign scanner from event: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) GetAssignedEvents(ctx context.Context, scannerID uuid.UUID) ([]*entities.ScannerAssignedEvent, error) {
	query := `
		SELECT 
			sea.scanner_id,
			sea.event_id,
			sea.assigned_by,
			sea.assigned_at,
			e.name as event_name,
			e.event_date,
			e.venue_name,
			e.venue_city,
			e.status as event_status
		FROM scanner_event_assignments sea
		JOIN events e ON sea.event_id = e.id
		WHERE sea.scanner_id = $1
		ORDER BY e.event_date ASC`
	
	var assignments []*entities.ScannerAssignedEvent
	err := r.db.SelectContext(ctx, &assignments, query, scannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned events: %w", err)
	}
	
	return assignments, nil
}

func (r *scannerUserRepository) GetEventScanners(ctx context.Context, eventID uuid.UUID) ([]*entities.ScannerUser, error) {
	query := `
		SELECT 
			su.id, su.username, su.email, su.password_hash, su.first_name, su.last_name, 
			su.role, su.status, su.created_at, su.updated_at, su.last_login_at, 
			su.created_by, su.updated_by
		FROM scanner_users su
		JOIN scanner_event_assignments sea ON su.id = sea.scanner_id
		WHERE sea.event_id = $1 AND su.is_active = true
		ORDER BY su.username ASC`
	
	var scanners []*entities.ScannerUser
	err := r.db.SelectContext(ctx, &scanners, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event scanners: %w", err)
	}
	
	return scanners, nil
}

// Session management
func (r *scannerUserRepository) CreateSession(ctx context.Context, session *entities.ScannerSession) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	
	now := time.Now()
	session.StartTime = now
	
	query := `
		INSERT INTO scanner_sessions (id, scanner_id, event_id, start_time, scans_count, valid_scans, invalid_scans, total_revenue, is_active)
		VALUES (:id, :scanner_id, :event_id, :start_time, :scans_count, :valid_scans, :invalid_scans, :total_revenue, :is_active)`
	
	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("failed to create scanner session: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) GetActiveSession(ctx context.Context, scannerID uuid.UUID) (*entities.ScannerSession, error) {
	query := `
		SELECT id, scanner_id, event_id, start_time, end_time, scans_count, 
			   valid_scans, invalid_scans, total_revenue, is_active, notes
		FROM scanner_sessions 
		WHERE scanner_id = $1 AND end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1`
	
	var session entities.ScannerSession
	err := r.db.GetContext(ctx, &session, query, scannerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFoundError
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}
	
	return &session, nil
}

func (r *scannerUserRepository) EndSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		UPDATE scanner_sessions 
		SET ended_at = NOW() 
		WHERE id = $1 AND ended_at IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to end session: %w", err)
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

func (r *scannerUserRepository) UpdateSessionStats(ctx context.Context, sessionID uuid.UUID, scansCount, validScans, invalidScans int, totalRevenue float64) error {
	query := `
		UPDATE scanner_sessions 
		SET scans_count = $2, valid_scans = $3, invalid_scans = $4, total_revenue = $5
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, sessionID, scansCount, validScans, invalidScans, totalRevenue)
	if err != nil {
		return fmt.Errorf("failed to update session stats: %w", err)
	}
	
	return nil
}

// Audit and logging
func (r *scannerUserRepository) LogActivity(ctx context.Context, log *entities.ScannerAuditLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	
	log.CreatedAt = time.Now()
	
	query := `
		INSERT INTO scanner_audit_logs (id, scanner_id, session_id, action, resource_type, 
										resource_id, details, ip_address, user_agent, created_at)
		VALUES (:id, :scanner_id, :session_id, :action, :resource_type, :resource_id, 
				:details, :ip_address, :user_agent, :created_at)`
	
	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to log scanner activity: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) RecordLogin(ctx context.Context, loginHistory *entities.ScannerLoginHistory) error {
	if loginHistory.ID == uuid.Nil {
		loginHistory.ID = uuid.New()
	}
	
	loginHistory.LoginAt = time.Now()
	
	query := `
		INSERT INTO scanner_login_history (id, scanner_id, login_at, ip_address, user_agent, success)
		VALUES (:id, :scanner_id, :login_at, :ip_address, :user_agent, :success)`
	
	_, err := r.db.NamedExecContext(ctx, query, loginHistory)
	if err != nil {
		return fmt.Errorf("failed to record login history: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) GetLoginHistory(ctx context.Context, scannerID uuid.UUID, limit int) ([]*entities.ScannerLoginHistory, error) {
	query := `
		SELECT id, scanner_id, login_at, ip_address, user_agent, success
		FROM scanner_login_history 
		WHERE scanner_id = $1
		ORDER BY login_at DESC
		LIMIT $2`
	
	var history []*entities.ScannerLoginHistory
	err := r.db.SelectContext(ctx, &history, query, scannerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}
	
	return history, nil
}

func (r *scannerUserRepository) GetAuditLog(ctx context.Context, scannerID uuid.UUID, filter *repositories.ScannerAuditFilter) ([]*entities.ScannerAuditLog, *repositories.PaginationResult, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	conditions = append(conditions, fmt.Sprintf("scanner_id = $%d", argIndex))
	args = append(args, scannerID)
	argIndex++
	
	if filter != nil {
		if filter.SessionID != nil {
			conditions = append(conditions, fmt.Sprintf("session_id = $%d", argIndex))
			args = append(args, *filter.SessionID)
			argIndex++
		}
		
		if filter.Action != "" {
			conditions = append(conditions, fmt.Sprintf("action = $%d", argIndex))
			args = append(args, filter.Action)
			argIndex++
		}
		
		if filter.ResourceType != "" {
			conditions = append(conditions, fmt.Sprintf("resource_type = $%d", argIndex))
			args = append(args, filter.ResourceType)
			argIndex++
		}
	}
	
	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM scanner_audit_logs %s", whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count audit logs: %w", err)
	}
	
	// Build main query with pagination
	query := fmt.Sprintf(`
		SELECT id, scanner_id, session_id, action, resource_type, resource_id, 
			   details, ip_address, user_agent, created_at
		FROM scanner_audit_logs %s
		ORDER BY created_at DESC`, whereClause)
	
	if filter != nil && filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
		
		if filter.GetOffset() > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.GetOffset())
		}
	}
	
	var logs []*entities.ScannerAuditLog
	err = r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	
	pagination := repositories.NewPaginationResult(filter.Page, filter.Limit, total)
	
	return logs, pagination, nil
}

// Ticket validation
func (r *scannerUserRepository) ValidateTicket(ctx context.Context, validation *entities.TicketValidation) error {
	if validation.ID == uuid.Nil {
		validation.ID = uuid.New()
	}
	
	validation.ValidationTimestamp = time.Now()
	
	query := `
		INSERT INTO ticket_validations (id, ticket_id, scanner_id, session_id, validation_result, 
										validation_timestamp, notes)
		VALUES (:id, :ticket_id, :scanner_id, :session_id, :validation_result, 
				:validation_timestamp, :notes)`
	
	_, err := r.db.NamedExecContext(ctx, query, validation)
	if err != nil {
		return fmt.Errorf("failed to validate ticket: %w", err)
	}
	
	return nil
}

func (r *scannerUserRepository) GetValidationHistory(ctx context.Context, scannerID uuid.UUID, filter *repositories.TicketValidationFilter) ([]*entities.TicketValidation, *repositories.PaginationResult, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	conditions = append(conditions, fmt.Sprintf("scanner_id = $%d", argIndex))
	args = append(args, scannerID)
	argIndex++
	
	if filter != nil {
		if filter.SessionID != nil {
			conditions = append(conditions, fmt.Sprintf("session_id = $%d", argIndex))
			args = append(args, *filter.SessionID)
			argIndex++
		}
		
		if filter.ValidationResult != "" {
			conditions = append(conditions, fmt.Sprintf("validation_result = $%d", argIndex))
			args = append(args, filter.ValidationResult)
			argIndex++
		}
	}
	
	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ticket_validations %s", whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count validations: %w", err)
	}
	
	// Build main query with pagination
	query := fmt.Sprintf(`
		SELECT id, ticket_id, scanner_id, session_id, validation_result, 
			   error_message, validated_at, ip_address
		FROM ticket_validations %s
		ORDER BY validated_at DESC`, whereClause)
	
	if filter != nil && filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
		
		if filter.GetOffset() > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.GetOffset())
		}
	}
	
	var validations []*entities.TicketValidation
	err = r.db.SelectContext(ctx, &validations, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get validation history: %w", err)
	}
	
	pagination := repositories.NewPaginationResult(filter.Page, filter.Limit, total)
	
	return validations, pagination, nil
}

// Statistics
func (r *scannerUserRepository) GetScannerStats(ctx context.Context, scannerID uuid.UUID, eventID *uuid.UUID) (*repositories.ScannerStats, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	conditions = append(conditions, fmt.Sprintf("scanner_id = $%d", argIndex))
	args = append(args, scannerID)
	argIndex++
	
	if eventID != nil {
		conditions = append(conditions, fmt.Sprintf("event_id = $%d", argIndex))
		args = append(args, *eventID)
		argIndex++
	}
	
	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	
	query := fmt.Sprintf(`
		SELECT 
			scanner_id,
			COUNT(*) as total_sessions,
			COALESCE(SUM(scans_count), 0) as total_scans,
			COALESCE(SUM(valid_scans), 0) as valid_scans,
			COALESCE(SUM(invalid_scans), 0) as invalid_scans,
			COALESCE(SUM(total_revenue), 0) as total_revenue,
			CASE 
				WHEN SUM(scans_count) > 0 THEN (SUM(valid_scans)::float / SUM(scans_count)::float) * 100
				ELSE 0 
			END as success_rate,
			MAX(started_at) as last_active_at
		FROM scanner_sessions %s
		GROUP BY scanner_id`, whereClause)
	
	var stats repositories.ScannerStats
	err := r.db.GetContext(ctx, &stats, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty stats if no sessions found
			return &repositories.ScannerStats{
				ScannerID: scannerID,
			}, nil
		}
		return nil, fmt.Errorf("failed to get scanner stats: %w", err)
	}
	
	// Get events assigned count
	eventCountQuery := `SELECT COUNT(*) FROM scanner_event_assignments WHERE scanner_id = $1`
	err = r.db.GetContext(ctx, &stats.EventsAssigned, eventCountQuery, scannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events assigned count: %w", err)
	}
	
	return &stats, nil
}

func (r *scannerUserRepository) GetEventScanStats(ctx context.Context, eventID uuid.UUID) (*repositories.EventScanStats, error) {
	query := `
		SELECT 
			$1 as event_id,
			COUNT(DISTINCT scanner_id) as total_scanners,
			COUNT(DISTINCT CASE WHEN ended_at IS NULL THEN scanner_id END) as active_scanners,
			COALESCE(SUM(scans_count), 0) as total_scans,
			COALESCE(SUM(valid_scans), 0) as valid_scans,
			COALESCE(SUM(invalid_scans), 0) as invalid_scans,
			COALESCE(SUM(total_revenue), 0) as total_revenue,
			CASE 
				WHEN SUM(scans_count) > 0 THEN (SUM(valid_scans)::float / SUM(scans_count)::float) * 100
				ELSE 0 
			END as success_rate
		FROM scanner_sessions 
		WHERE event_id = $1`
	
	var stats repositories.EventScanStats
	err := r.db.GetContext(ctx, &stats, query, eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty stats if no sessions found
			return &repositories.EventScanStats{
				EventID: eventID,
			}, nil
		}
		return nil, fmt.Errorf("failed to get event scan stats: %w", err)
	}
	
	// Get scanner breakdown
	breakdownQuery := `
		SELECT 
			ss.scanner_id,
			CONCAT(su.first_name, ' ', su.last_name) as scanner_name,
			COALESCE(SUM(ss.scans_count), 0) as total_scans,
			COALESCE(SUM(ss.valid_scans), 0) as valid_scans,
			COALESCE(SUM(ss.invalid_scans), 0) as invalid_scans,
			COALESCE(SUM(ss.total_revenue), 0) as revenue,
			CASE 
				WHEN SUM(ss.scans_count) > 0 THEN (SUM(ss.valid_scans)::float / SUM(ss.scans_count)::float) * 100
				ELSE 0 
			END as success_rate,
			MAX(ss.started_at) as last_scan_at
		FROM scanner_sessions ss
		JOIN scanner_users su ON ss.scanner_id = su.id
		WHERE ss.event_id = $1
		GROUP BY ss.scanner_id, su.first_name, su.last_name
		ORDER BY total_scans DESC`
	
	var breakdown []repositories.ScannerEventStats
	err = r.db.SelectContext(ctx, &breakdown, breakdownQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scanner breakdown: %w", err)
	}
	
	stats.ScannerBreakdown = breakdown
	
	return &stats, nil
}

