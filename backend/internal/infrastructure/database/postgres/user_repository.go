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

type userRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	}
}

func NewUserRepository(db *sqlx.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func NewUserRepositoryWithTx(tx *sqlx.Tx) repositories.UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			id, first_name, last_name, email, phone, password_hash, 
			auth_provider, momo_id, email_verified, phone_verified, 
			is_active, created_at, updated_at
		) VALUES (
			:id, :first_name, :last_name, :email, :phone, :password_hash,
			:auth_provider, :momo_id, :email_verified, :phone_verified,
			:is_active, :created_at, :updated_at
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "email") {
					return entities.ErrUserEmailExists
				}
				if strings.Contains(pqErr.Detail, "phone") {
					return entities.ErrUserPhoneExists
				}
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE id = $1 AND is_active = true`
	
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE email = $1 AND is_active = true`
	
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE phone = $1 AND is_active = true`
	
	err := r.db.GetContext(ctx, &user, query, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByMoMoID(ctx context.Context, momoID string) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE momo_id = $1 AND is_active = true`
	
	err := r.db.GetContext(ctx, &user, query, momoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by MoMo ID: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByIdentifier(ctx context.Context, identifier string) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE (email = $1 OR phone = $1 OR momo_id = $1) AND is_active = true`
	
	err := r.db.GetContext(ctx, &user, query, identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by identifier: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()
	
	query := `
		UPDATE users SET
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			phone = :phone,
			password_hash = :password_hash,
			auth_provider = :auth_provider,
			momo_id = :momo_id,
			email_verified = :email_verified,
			phone_verified = :phone_verified,
			is_active = :is_active,
			last_login = :last_login,
			updated_at = :updated_at
		WHERE id = :id AND is_active = true`
	
	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "email") {
					return entities.ErrUserEmailExists
				}
				if strings.Contains(pqErr.Detail, "phone") {
					return entities.ErrUserPhoneExists
				}
			}
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) List(ctx context.Context, filter repositories.UserFilter) ([]*entities.User, *repositories.PaginationResult, error) {
	var users []*entities.User
	var totalCount int
	
	// Build WHERE clause
	whereConditions := []string{"is_active = true"}
	args := []interface{}{}
	argIndex := 1
	
	if filter.AuthProvider != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("auth_provider = $%d", argIndex))
		args = append(args, *filter.AuthProvider)
		argIndex++
	}
	
	if filter.EmailVerified != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("email_verified = $%d", argIndex))
		args = append(args, *filter.EmailVerified)
		argIndex++
	}
	
	if filter.PhoneVerified != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("phone_verified = $%d", argIndex))
		args = append(args, *filter.PhoneVerified)
		argIndex++
	}
	
	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(email ILIKE $%d OR phone ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereClause)
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
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
	
	// Build main query with pagination
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, phone, password_hash,
			   auth_provider, momo_id, email_verified, phone_verified,
			   is_active, last_login, created_at, updated_at
		FROM users 
		WHERE %s 
		ORDER BY %s 
		LIMIT $%d OFFSET $%d`, whereClause, orderBy, argIndex, argIndex+1)
	
	args = append(args, filter.Limit, offset)
	
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	pagination := &repositories.PaginationResult{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
	}
	
	return users, pagination, nil
}

func (r *userRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	return exists, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check user email existence: %w", err)
	}
	
	return exists, nil
}

func (r *userRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, phone)
	if err != nil {
		return false, fmt.Errorf("failed to check user phone existence: %w", err)
	}
	
	return exists, nil
}

func (r *userRepository) ExistsByMoMoID(ctx context.Context, momoID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE momo_id = $1 AND is_active = true)`
	
	err := r.db.GetContext(ctx, &exists, query, momoID)
	if err != nil {
		return false, fmt.Errorf("failed to check user MoMo ID existence: %w", err)
	}
	
	return exists, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET email_verified = true, updated_at = NOW() WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) VerifyPhone(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET phone_verified = true, updated_at = NOW() WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to verify phone: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*repositories.UserStats, error) {
	var stats repositories.UserStats
	
	query := `
		SELECT 
			$1 as user_id,
			COALESCE(COUNT(o.id), 0) as total_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'paid' THEN 1 END), 0) as paid_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'cancelled' THEN 1 END), 0) as cancelled_orders,
			COALESCE(COUNT(CASE WHEN o.status = 'refunded' THEN 1 END), 0) as refunded_orders,
			COALESCE(COUNT(t.id), 0) as total_tickets,
			COALESCE(COUNT(CASE WHEN t.status = 'redeemed' THEN 1 END), 0) as redeemed_tickets,
			COALESCE(SUM(CASE WHEN o.status = 'paid' THEN o.total_amount ELSE 0 END), 0) as total_spent,
			COALESCE(AVG(CASE WHEN o.status = 'paid' THEN o.total_amount END), 0) as average_order_value,
			MIN(CASE WHEN o.status = 'paid' THEN o.created_at END) as first_order_at,
			MAX(CASE WHEN o.status = 'paid' THEN o.created_at END) as last_order_at,
			(
				SELECT e.venue_city 
				FROM orders o2 
				JOIN events e ON o2.event_id = e.id 
				WHERE o2.user_id = $1 AND o2.status = 'paid'
				GROUP BY e.venue_city 
				ORDER BY COUNT(*) DESC 
				LIMIT 1
			) as favorite_venue_city,
			COALESCE(COUNT(DISTINCT CASE WHEN o.status = 'paid' THEN o.event_id END), 0) as events_attended
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		LEFT JOIN tickets t ON o.id = t.order_id
		WHERE u.id = $1 AND u.is_active = true
		GROUP BY u.id`
	
	err := r.db.GetContext(ctx, &stats, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	
	return &stats, nil
}


// IncrementFailedAttempts increments the failed login attempts counter
func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = failed_login_attempts + 1,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment failed login attempts: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}


// ResetFailedAttempts resets the failed login attempts counter
func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = 0,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset failed login attempts: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	
	return nil
}

