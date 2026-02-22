package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

type adminUserRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
		QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	}
}

// NewAdminUserRepository creates a new admin user repository
func NewAdminUserRepository(db *sqlx.DB) repositories.AdminUserRepository {
	return &adminUserRepository{db: db}
}

// NewAdminUserRepositoryWithTx creates a new admin user repository with a transaction
func NewAdminUserRepositoryWithTx(tx *sqlx.Tx) repositories.AdminUserRepository {
	return &adminUserRepository{db: tx}
}

// Create creates a new admin user
func (r *adminUserRepository) Create(ctx context.Context, admin *entities.AdminUser) error {
	query := `
		INSERT INTO admin_users (
			id, email, password_hash, first_name, last_name, role, permissions,
			is_active, two_factor_enabled, failed_login_attempts, login_attempts,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`
	
	_, err := r.db.ExecContext(ctx, query,
		admin.ID, admin.Email, admin.PasswordHash, admin.FirstName, admin.LastName,
		admin.Role, pq.Array(admin.Permissions), admin.IsActive, admin.TwoFactorEnabled,
		admin.LoginAttempts, admin.CreatedAt, admin.UpdatedAt,
	)
	
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "admin_users_email_key" {
					return entities.ErrAdminUserAlreadyExists
				}
			}
		}
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	
	return nil
}

// GetByID retrieves an admin user by ID
func (r *adminUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret,
			   login_attempts, locked_until, last_login, created_at, updated_at
		FROM admin_users
		WHERE id = $1`
	
	admin := &entities.AdminUser{}
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(admin)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrAdminUserNotFound
		}
		return nil, fmt.Errorf("failed to get admin user by ID: %w", err)
	}
	
	return admin, nil
}

// GetByEmail retrieves an admin user by email
func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret,
			   login_attempts, locked_until, last_login, created_at, updated_at
		FROM admin_users
		WHERE email = $1`
	
	admin := &entities.AdminUser{}
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(admin)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrAdminUserNotFound
		}
		return nil, fmt.Errorf("failed to get admin user by email: %w", err)
	}
	
	return admin, nil
}

// Update updates an admin user
func (r *adminUserRepository) Update(ctx context.Context, admin *entities.AdminUser) error {
	query := `
		UPDATE admin_users
		SET email = $2, password_hash = $3, first_name = $4, last_name = $5,
			role = $6, permissions = $7, is_active = $8, two_factor_enabled = $9,
			two_factor_secret = $10, failed_login_attempts = $11, login_attempts = $12,
			locked_until = $13, last_login = $14, last_login_at = $15, updated_at = $16
		WHERE id = $1`
	
	admin.UpdatedAt = time.Now().UTC()
	
	result, err := r.db.ExecContext(ctx, query,
		admin.ID, admin.Email, admin.PasswordHash, admin.FirstName, admin.LastName,
		admin.Role, pq.Array(admin.Permissions), admin.IsActive, admin.TwoFactorEnabled,
		admin.TwoFactorSecret, admin.LoginAttempts,
		admin.LockedUntil, admin.LastLogin, admin.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update admin user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}

// Delete deletes an admin user
func (r *adminUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM admin_users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete admin user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}

// List lists admin users with filtering and pagination
func (r *adminUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list admin users: %w", err)
		}
		
		return admins, nil
	}
	
	// ListByRole lists admin users by role
func (r *adminUserRepository) CountWithFilter(ctx context.Context, filter repositories.AdminUserFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM admin_users WHERE 1=1`
	
	args := []interface{}{}
	argIndex := 1
	
	// Apply filters
	if filter.Email != nil && *filter.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", argIndex)
		args = append(args, "%"+*filter.Email+"%")
		argIndex++
	}
	
	if filter.Role != nil && string(*filter.Role) != "" {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, string(*filter.Role))
		argIndex++
	}
	
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	var count int64
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users: %w", err)
	}
	
	return count, nil
}

// Count counts all admin users (backward compatibility)
func (r *adminUserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM admin_users`
	
	err := r.db.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users: %w", err)
	}
	
	return count, nil
}

// ValidateCredentials validates admin credentialss
func (r *adminUserRepository) ValidateCredentials(ctx context.Context, email, passwordHash string) (*entities.AdminUser, error) {
	admin, err := r.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if admin.PasswordHash != passwordHash {
		return nil, entities.ErrInvalidCredentials
	}
	
	return admin, nil
}

// RecordLoginAttempt records a failed login attempt
func (r *adminUserRepository) RecordLoginAttempt(ctx context.Context, email string) error {
	query := `
		UPDATE admin_users 
		SET failed_login_attempts = failed_login_attempts + 1,
			login_attempts = login_attempts + 1,
			locked_until = CASE 
				WHEN failed_login_attempts + 1 >= 5 THEN NOW() + INTERVAL '30 minutes'
				ELSE locked_until
			END,
			updated_at = NOW()
		WHERE email = $1`
	
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

// RecordSuccessfulLogin records a successful login
func (r *adminUserRepository) RecordSuccessfulLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET last_login = NOW(),
			last_login_at = NOW(),
			failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// UnlockAccount unlocks an admin account
func (r *adminUserRepository) UnlockAccount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// UpdateRole updates an admin user's role
func (r *adminUserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role entities.AdminRole) error {
	query := `
		UPDATE admin_users 
		SET role = $2, 
			permissions = $3,
			updated_at = NOW()
		WHERE id = $1`
	
	permissions := entities.GetRolePermissions(role)
	_, err := r.db.ExecContext(ctx, query, id, role, pq.Array(permissions))
	return err
}

// UpdatePermissions updates an admin user's permissions
func (r *adminUserRepository) UpdatePermissions(ctx context.Context, id uuid.UUID, permissions []entities.AdminPermission) error {
	query := `
		UPDATE admin_users 
		SET permissions = $2,
			updated_at = NOW()
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id, pq.Array(permissions))
	return err
}

// GetAdminsByRole gets all admins with a specific role
func (r *adminUserRepository) GetAdminsByRole(ctx context.Context, role entities.AdminRole) ([]*entities.AdminUser, error) {
	return r.ListByRole(ctx, role, 100, 0) // Default limit and offset
}

// EnableTwoFactor enables two-factor authentication
func (r *adminUserRepository) EnableTwoFactor(ctx context.Context, id uuid.UUID, secret string) error {
	query := `
		UPDATE admin_users 
		SET two_factor_enabled = true,
			two_factor_secret = $2,
			updated_at = NOW()
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id, secret)
	return err
}

// DisableTwoFactor disables two-factor authentication
func (r *adminUserRepository) DisableTwoFactor(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET two_factor_enabled = false,
			two_factor_secret = NULL,
			updated_at = NOW()
		WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetTwoFactorSecret gets the two-factor secret for an admin
func (r *adminUserRepository) GetTwoFactorSecret(ctx context.Context, id uuid.UUID) (string, error) {
	query := `SELECT two_factor_secret FROM admin_users WHERE id = $1`
	
	var secret sql.NullString
	err := r.db.QueryRowxContext(ctx, query, id).Scan(&secret)
	if err != nil {
		return "", err
	}
	
	return secret.String, nil
}

// GetRecentLogins gets admins with recent logins
func (r *adminUserRepository) GetRecentLogins(ctx context.Context, since time.Time) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE last_login >= $1
		ORDER BY last_login DESC`
	
	rows, err := r.db.QueryxContext(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logins: %w", err)
	}
	defer rows.Close()
	
	var admins []*entities.AdminUser
	for rows.Next() {
		admin := &entities.AdminUser{}
		if err := rows.StructScan(admin); err != nil {
			return nil, fmt.Errorf("failed to scan admin user: %w", err)
		}
		admins = append(admins, admin)
	}
	
	return admins, nil
}

// GetLockedAccounts gets all locked admin accounts
func (r *adminUserRepository) GetLockedAccounts(ctx context.Context) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE locked_until > NOW()
		ORDER BY locked_until DESC`
	
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get locked accounts: %w", err)
	}
	defer rows.Close()
	
	var admins []*entities.AdminUser
	for rows.Next() {
		admin := &entities.AdminUser{}
		if err := rows.StructScan(admin); err != nil {
			return nil, fmt.Errorf("failed to scan admin user: %w", err)
		}
		admins = append(admins, admin)
	}
	
	return admins, nil
}

// GetStats returns admin user statistics
func (r *adminUserRepository) GetStats(ctx context.Context) (*repositories.AdminUserStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_admins,
			COUNT(*) FILTER (WHERE is_active = true) as active_admins,
			COUNT(*) FILTER (WHERE is_active = false) as inactive_admins,
			COUNT(*) FILTER (WHERE locked_until > NOW()) as locked_admins,
			COUNT(*) FILTER (WHERE two_factor_enabled = true) as two_factor_enabled,
			COUNT(*) FILTER (WHERE last_login >= NOW() - INTERVAL '24 hours') as recent_logins
		FROM admin_users`
	
	stats := &repositories.AdminUserStats{}
	err := r.db.QueryRowxContext(ctx, query).StructScan(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin user stats: %w", err)
	}
	
	// Get role statistics
	roleQuery := `
		SELECT role, COUNT(*) as count
		FROM admin_users
		GROUP BY role`
	
	rows, err := r.db.QueryxContext(ctx, roleQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get role stats: %w", err)
	}
	defer rows.Close()
	
	stats.AdminsByRole = make(map[entities.AdminRole]int64)
	for rows.Next() {
		var role entities.AdminRole
		var count int64
		if err := rows.Scan(&role, &count); err != nil {
			return nil, fmt.Errorf("failed to scan role stats: %w", err)
		}
		stats.AdminsByRole[role] = count
	}
	
	return stats, nil
}

// GetLoginHistory gets login history for an admin
func (r *adminUserRepository) GetLoginHistory(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]*entities.AdminLoginHistory, error) {
	// This would require a separate login_history table in a real implementation
	// For now, return empty slice
	return []*entities.AdminLoginHistory{}, nil
}

// BulkUpdateStatus updates status for multiple admins
func (r *adminUserRepository) BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status entities.AdminStatus) error {
	if len(ids) == 0 {
		return nil
	}
	
	// Convert AdminStatus to boolean for database
	isActive := status == entities.AdminStatusActive
	
	query := `
		UPDATE admin_users 
		SET is_active = $1, updated_at = NOW()
		WHERE id = ANY($2)`
	
	_, err := r.db.ExecContext(ctx, query, isActive, pq.Array(ids))
	return err
}

// BulkUnlock unlocks multiple admin accounts
func (r *adminUserRepository) BulkUnlock(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	
	query := `
		UPDATE admin_users 
		SET failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = ANY($1)`
	
	_, err := r.db.ExecContext(ctx, query, pq.Array(ids))
	return err
}

// ExistsByEmail checks if an admin user exists by email
func (r *adminUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admin_users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRowxContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check admin user existence by email: %w", err)
	}
	
	return exists, nil
}

// ExistsByID checks if an admin user exists by ID
func (r *adminUserRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admin_users WHERE id = $1)`
	
	var exists bool
	err := r.db.QueryRowxContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check admin user existence by ID: %w", err)
	}
	
	return exists, nil
}

// Activate activates an admin user account
func (r *adminUserRepository) Activate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET status = 'active', updated_at = NOW()
		WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to activate admin user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// BulkDelete deletes multiple admin users by their IDs
func (r *adminUserRepository) BulkDelete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	
	query := `DELETE FROM admin_users WHERE id = ANY($1)`
	
	_, err := r.db.ExecContext(ctx, query, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("failed to bulk delete admin users: %w", err)
	}
	
	return nil
}


// CountByRole counts admin users by role
func (r *adminUserRepository) CountByRole(ctx context.Context, role entities.AdminRole) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM admin_users WHERE role = $1 AND is_active = true`
	
	err := r.db.QueryRowxContext(ctx, query, string(role)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users by role: %w", err)
	}
	
	return count, nil
}

// CountByStatus counts admin users by status
func (r *adminUserRepository) CountByStatus(ctx context.Context, status entities.AdminStatus) (int64, error) {
	var count int64
	var query string
	
	switch status {
	case entities.AdminStatusActive:
		query = `SELECT COUNT(*) FROM admin_users WHERE is_active = true AND is_active = true`
	case entities.AdminStatusInactive:
		query = `SELECT COUNT(*) FROM admin_users WHERE is_active = false AND is_active = true`
	case entities.AdminStatusLocked:
		query = `SELECT COUNT(*) FROM admin_users WHERE locked_until IS NOT NULL AND locked_until > NOW() AND is_active = true`
	default:
		return 0, fmt.Errorf("unsupported admin status: %s", status)
	}
	
	err := r.db.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users by status: %w", err)
	}
	
	return count, nil
}


// Deactivate deactivates an admin user account
func (r *adminUserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET is_active = false,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate admin user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// GetActiveAdminsCount returns the count of active admin users
func (r *adminUserRepository) GetActiveAdminsCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM admin_users WHERE is_active = true AND is_active = true`
	
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active admin count: %w", err)
	}
	
	return count, nil
}


// GetAdminsByCreatedDate retrieves admin users created within a date range
func (r *adminUserRepository) GetAdminsByCreatedDate(ctx context.Context, startDate, endDate time.Time) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions, 
			   is_active, last_login_at, failed_login_attempts, locked_until, 
			   two_factor_enabled, two_factor_secret, created_at, updated_at
		FROM admin_users 
		WHERE created_at >= $1 AND created_at <= $2 
		  AND is_active = true
		ORDER BY created_at DESC`
	
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get admins by created date: %w", err)
	}
	
	return admins, nil
}


// GetMostActiveAdmins retrieves the most active admin users
func (r *adminUserRepository) GetMostActiveAdmins(ctx context.Context, limit int) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions, 
			   is_active, last_login_at, failed_login_attempts, locked_until, 
			   two_factor_enabled, two_factor_secret, created_at, updated_at
		FROM admin_users 
		WHERE is_active = true AND is_active = true
		ORDER BY last_login_at DESC NULLS LAST
		LIMIT $1`
	
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most active admins: %w", err)
	}
	
	return admins, nil
}


// HasPermission checks if an admin user has a specific permission
func (r *adminUserRepository) HasPermission(ctx context.Context, id uuid.UUID, permission entities.AdminPermission) (bool, error) {
	query := `
		SELECT permissions 
		FROM admin_users 
		WHERE id = $1 AND is_active = true AND is_active = true`
	
	var permissions []entities.AdminPermission
	err := r.db.GetContext(ctx, &permissions, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check admin permission: %w", err)
	}
	
	// Check if the permission exists in the permissions array
	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	
	return false, nil
}


// IncrementFailedAttempts increments the failed login attempts for an admin user
func (r *adminUserRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET failed_login_attempts = failed_login_attempts + 1,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment failed attempts: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// ListByRole lists admin users by role
func (r *adminUserRepository) ListByRole(ctx context.Context, role entities.AdminRole, limit, offset int) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE role = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, query, string(role), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list admin users by role: %w", err)
	}
	
	return admins, nil
}


// ListByStatus lists admin users by status
func (r *adminUserRepository) ListByStatus(ctx context.Context, status entities.AdminStatus, limit, offset int) ([]*entities.AdminUser, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE status = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, query, string(status), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list admin users by status: %w", err)
	}
	
	return admins, nil
}


// LockAccount locks an admin user account until a specific time
func (r *adminUserRepository) LockAccount(ctx context.Context, id uuid.UUID, lockUntil *time.Time) error {
	query := `
		UPDATE admin_users 
		SET locked_until = $2,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id, lockUntil)
	if err != nil {
		return fmt.Errorf("failed to lock admin account: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// LogLogin logs a login attempt for an admin user
func (r *adminUserRepository) LogLogin(ctx context.Context, adminID uuid.UUID, ipAddress, userAgent string, success bool) error {
	// This would typically insert into a separate admin_login_history table
	// For now, we'll just update the last_login_at field if successful
	if success {
		query := `
			UPDATE admin_users 
			SET last_login_at = NOW(),
				updated_at = NOW()
			WHERE id = $1 AND is_active = true`
		
		_, err := r.db.ExecContext(ctx, query, adminID)
		if err != nil {
			return fmt.Errorf("failed to log successful login: %w", err)
		}
	}
	
	// In a full implementation, this would also insert into admin_login_history table
	// with ipAddress, userAgent, success status, etc.
	
	return nil
}


// ResetFailedAttempts resets the failed login attempts for an admin user
func (r *adminUserRepository) ResetFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET failed_login_attempts = 0,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to reset failed attempts: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// Search searches admin users by query string
func (r *adminUserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.AdminUser, error) {
	searchQuery := `
		SELECT id, email, password_hash, first_name, last_name, role, permissions,
			   is_active, two_factor_enabled, two_factor_secret, failed_login_attempts,
			   login_attempts, locked_until, last_login, last_login_at, created_at, updated_at
		FROM admin_users
		WHERE (
			email ILIKE $1 OR 
			first_name ILIKE $1 OR 
			last_name ILIKE $1 OR
			CONCAT(first_name, ' ', last_name) ILIKE $1
		) AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	searchPattern := "%" + query + "%"
	var admins []*entities.AdminUser
	err := r.db.SelectContext(ctx, &admins, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search admin users: %w", err)
	}
	
	return admins, nil
}


// UpdateLastLogin updates the last login timestamp for an admin user
func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE admin_users 
		SET last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// UpdatePassword updates the password hash for an admin user
func (r *adminUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `
		UPDATE admin_users 
		SET password_hash = $2,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// UpdateStatus updates the status for an admin user
func (r *adminUserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.AdminStatus) error {
	query := `
		UPDATE admin_users 
		SET status = $2,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id, string(status))
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}


// UpdateTwoFactorSecret updates the two-factor authentication secret for an admin user
func (r *adminUserRepository) UpdateTwoFactorSecret(ctx context.Context, id uuid.UUID, secret string) error {
	query := `
		UPDATE admin_users 
		SET two_factor_secret = $2,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id, secret)
	if err != nil {
		return fmt.Errorf("failed to update two-factor secret: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrAdminUserNotFound
	}
	
	return nil
}

