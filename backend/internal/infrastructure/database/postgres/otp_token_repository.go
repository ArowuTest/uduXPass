package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

type otpTokenRepository struct {
	db DB
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// NewOTPTokenRepository creates a new OTP token repository
func NewOTPTokenRepository(db *sqlx.DB) repositories.OTPTokenRepository {
	return &otpTokenRepository{
		db: db,
	}
}

func (r *otpTokenRepository) Create(ctx context.Context, token *entities.OTPToken) error {
	query := `
		INSERT INTO otp_tokens (
			id, phone, email, code, purpose, status, expires_at,
			attempt_count, max_attempts, ip_address, user_agent,
			created_at, updated_at
		) VALUES (
			:id, :phone, :email, :code, :purpose, :status, :expires_at,
			:attempt_count, :max_attempts, :ip_address, :user_agent,
			:created_at, :updated_at
		)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to create OTP token: %w", err)
	}
	
	return nil
}

func (r *otpTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OTPToken, error) {
	var token entities.OTPToken
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at,
			   used_at, attempt_count, max_attempts, ip_address, user_agent,
			   created_at, updated_at
		FROM otp_tokens 
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &token, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to get OTP token by ID: %w", err)
	}
	
	return &token, nil
}

func (r *otpTokenRepository) GetByIdentifierAndPurpose(ctx context.Context, identifier string, purpose entities.OTPPurpose) (*entities.OTPToken, error) {
	var token entities.OTPToken
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at,
			   used_at, attempt_count, max_attempts, ip_address, user_agent,
			   created_at, updated_at
		FROM otp_tokens 
		WHERE phone = $1 AND purpose = $2 AND status = 'pending' AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	err := r.db.GetContext(ctx, &token, query, identifier, string(purpose))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to get OTP token by identifier and purpose: %w", err)
	}
	
	return &token, nil
}

func (r *otpTokenRepository) GetByToken(ctx context.Context, tokenValue string) (*entities.OTPToken, error) {
	var token entities.OTPToken
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at,
			   used_at, attempt_count, max_attempts, ip_address, user_agent,
			   created_at, updated_at
		FROM otp_tokens 
		WHERE code = $1
	`
	
	err := r.db.GetContext(ctx, &token, query, tokenValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to get OTP token by token: %w", err)
	}
	
	return &token, nil
}

func (r *otpTokenRepository) Update(ctx context.Context, token *entities.OTPToken) error {
	query := `
		UPDATE otp_tokens 
		SET phone = :phone, email = :email, code = :code, purpose = :purpose,
		    status = :status, expires_at = :expires_at,
		    attempt_count = :attempt_count, max_attempts = :max_attempts,
		    ip_address = :ip_address, user_agent = :user_agent,
		    updated_at = :updated_at
		WHERE id = :id
	`
	
	result, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to update OTP token: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}

func (r *otpTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM otp_tokens WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete OTP token: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}

func (r *otpTokenRepository) GetByCode(ctx context.Context, code string) (*entities.OTPToken, error) {
	var token entities.OTPToken
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at,
			   used_at, attempt_count, max_attempts, ip_address, user_agent,
			   created_at, updated_at
		FROM otp_tokens 
		WHERE code = $1
	`
	
	err := r.db.GetContext(ctx, &token, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to get OTP token by code: %w", err)
	}
	
	return &token, nil
}

func (r *otpTokenRepository) List(ctx context.Context, filter *repositories.OTPTokenFilter) ([]*entities.OTPToken, *repositories.PaginationResult, error) {
	var tokens []*entities.OTPToken
	var totalCount int
	
	// Build query conditions
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if filter.Phone != nil {
		conditions = append(conditions, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *filter.Phone)
		argIndex++
	}
	
	if filter.Purpose != nil {
		conditions = append(conditions, fmt.Sprintf("purpose = $%d", argIndex))
		args = append(args, *filter.Purpose)
		argIndex++
	}
	
	if filter.IsUsed != nil {
		// Use status column instead of deprecated is_used
		if *filter.IsUsed {
			conditions = append(conditions, "status = 'used'")
		} else {
			conditions = append(conditions, "status != 'used'")
		}
	}
	
	if filter.ExpiredOnly != nil && *filter.ExpiredOnly {
		conditions = append(conditions, "expires_at < NOW()")
	}
	
	if filter.ActiveOnly != nil && *filter.ActiveOnly {
		conditions = append(conditions, "expires_at > NOW() AND status = 'pending'")
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}
	
	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM otp_tokens %s", whereClause)
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count OTP tokens: %w", err)
	}
	
	// Build main query
	query := fmt.Sprintf(`
		SELECT id, phone, email, code, purpose, status, expires_at,
			   used_at, attempt_count, max_attempts, ip_address, user_agent,
			   created_at, updated_at
		FROM otp_tokens %s
	`, whereClause)
	
	// Add sorting
	if filter.SortBy != nil {
		query += fmt.Sprintf(" ORDER BY %s", *filter.SortBy)
		if filter.SortOrder != nil {
			query += " " + *filter.SortOrder
		}
	} else {
		query += " ORDER BY created_at DESC"
	}
	
	// Add pagination
	limit := 20
	if filter.Limit != nil {
		limit = *filter.Limit
	}
	
	offset := 0
	if filter.Page != nil && *filter.Page > 1 {
		offset = (*filter.Page - 1) * limit
	}
	
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	
	err = r.db.SelectContext(ctx, &tokens, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list OTP tokens: %w", err)
	}
	
	// Calculate pagination
	totalPages := (totalCount + limit - 1) / limit
	currentPage := 1
	if filter.Page != nil {
		currentPage = *filter.Page
	}
	
	paginationResult := &repositories.PaginationResult{
		Total:      totalCount,
		TotalPages: totalPages,
		Page:       currentPage,
		Limit:      limit,
		HasNext:    currentPage < totalPages,
		HasPrev:    currentPage > 1,
	}
	
	return tokens, paginationResult, nil
}

func (r *otpTokenRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM otp_tokens WHERE id = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check OTP token existence: %w", err)
	}
	
	return exists, nil
}

func (r *otpTokenRepository) ExistsByToken(ctx context.Context, tokenValue string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM otp_tokens WHERE code = $1)`
	
	err := r.db.GetContext(ctx, &exists, query, tokenValue)
	if err != nil {
		return false, fmt.Errorf("failed to check OTP token value existence: %w", err)
	}
	
	return exists, nil
}

func (r *otpTokenRepository) MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	query := `
		UPDATE otp_tokens 
		SET status = 'used', used_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to mark OTP token as used: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}

func (r *otpTokenRepository) CleanupExpiredTokens(ctx context.Context) (int, error) {
	query := `DELETE FROM otp_tokens WHERE expires_at < NOW()`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired OTP tokens: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return int(rowsAffected), nil
}

func (r *otpTokenRepository) CountActiveTokensByPhone(ctx context.Context, phone string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM otp_tokens 
		WHERE phone = $1 AND status = 'pending' AND expires_at > NOW()
	`
	
	err := r.db.GetContext(ctx, &count, query, phone)
	if err != nil {
		return 0, fmt.Errorf("failed to count active tokens by phone: %w", err)
	}
	
	return count, nil
}

func (r *otpTokenRepository) CountTokensByPhoneAndTimeRange(ctx context.Context, phone string, from, to time.Time) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM otp_tokens 
		WHERE phone = $1 AND created_at >= $2 AND created_at <= $3
	`
	
	err := r.db.GetContext(ctx, &count, query, phone, from, to)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens by phone and time range: %w", err)
	}
	
	return count, nil
}

func (r *otpTokenRepository) DeleteByPhone(ctx context.Context, phone string) error {
	query := `DELETE FROM otp_tokens WHERE phone = $1`
	
	_, err := r.db.ExecContext(ctx, query, phone)
	if err != nil {
		return fmt.Errorf("failed to delete OTP tokens by phone: %w", err)
	}
	
	return nil
}


// DeleteExpired deletes all expired OTP tokens
func (r *otpTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM otp_tokens WHERE expires_at < NOW()`
	
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired OTP tokens: %w", err)
	}
	
	return nil
}


// DeleteOlderThan deletes OTP tokens older than the specified duration
func (r *otpTokenRepository) DeleteOlderThan(ctx context.Context, duration time.Duration) error {
	query := `DELETE FROM otp_tokens WHERE created_at < NOW() - INTERVAL '%d seconds'`
	
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, int(duration.Seconds())))
	if err != nil {
		return fmt.Errorf("failed to delete old OTP tokens: %w", err)
	}
	
	return nil
}


// GetActiveByPhone retrieves active OTP tokens by phone number
func (r *otpTokenRepository) GetActiveByPhone(ctx context.Context, phone string) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1 
		  AND expires_at > NOW() 
		  AND status = 'pending'
		ORDER BY created_at DESC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get active OTP tokens by phone: %w", err)
	}
	
	return tokens, nil
}


// GetActiveByPhoneAndPurpose retrieves active OTP tokens by phone and purpose
func (r *otpTokenRepository) GetActiveByPhoneAndPurpose(ctx context.Context, phone string, purpose entities.OTPPurpose) (*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1 
		  AND purpose = $2
		  AND expires_at > NOW() 
		  AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1`
	
	var token entities.OTPToken
	err := r.db.GetContext(ctx, &token, query, phone, string(purpose))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to get active OTP token by phone and purpose: %w", err)
	}
	
	return &token, nil
}


// GetByPhone retrieves OTP tokens by phone number
func (r *otpTokenRepository) GetByPhone(ctx context.Context, phone string) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1 
		ORDER BY created_at DESC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP tokens by phone: %w", err)
	}
	
	return tokens, nil
}


// GetByPhoneAndPurpose retrieves OTP tokens by phone and purpose
func (r *otpTokenRepository) GetByPhoneAndPurpose(ctx context.Context, phone string, purpose entities.OTPPurpose) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1 AND purpose = $2
		ORDER BY created_at DESC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, phone, string(purpose))
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP tokens by phone and purpose: %w", err)
	}
	
	return tokens, nil
}


// GetExpiredTokens retrieves all expired tokens
func (r *otpTokenRepository) GetExpiredTokens(ctx context.Context) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE expires_at <= NOW() AND status = 'pending'
		ORDER BY expires_at ASC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired tokens: %w", err)
	}
	
	return tokens, nil
}


// GetRecentTokensByPhone retrieves recent tokens for a phone number
func (r *otpTokenRepository) GetRecentTokensByPhone(ctx context.Context, phone string, limit int) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1
		ORDER BY created_at DESC
		LIMIT $2`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, phone, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tokens by phone: %w", err)
	}
	
	return tokens, nil
}


// GetTokenStats retrieves statistics about OTP tokens
func (r *otpTokenRepository) GetTokenStats(ctx context.Context) (*repositories.OTPTokenStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_tokens,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tokens,
			COUNT(CASE WHEN status = 'used' THEN 1 END) as used_tokens,
			COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_tokens,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_tokens
		FROM otp_tokens`
	
	var stats repositories.OTPTokenStats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get token stats: %w", err)
	}
	
	return &stats, nil
}


// GetTokensByPurpose retrieves tokens by purpose
func (r *otpTokenRepository) GetTokensByPurpose(ctx context.Context, purpose entities.OTPPurpose) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE purpose = $1
		ORDER BY created_at DESC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, string(purpose))
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens by purpose: %w", err)
	}
	
	return tokens, nil
}


// GetTokensByStatus retrieves tokens by status
func (r *otpTokenRepository) GetTokensByStatus(ctx context.Context, status entities.OTPStatus) ([]*entities.OTPToken, error) {
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE status = $1
		ORDER BY created_at DESC`
	
	var tokens []*entities.OTPToken
	err := r.db.SelectContext(ctx, &tokens, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens by status: %w", err)
	}
	
	return tokens, nil
}


// IncrementAttempt increments the attempt count for an OTP token
func (r *otpTokenRepository) IncrementAttempt(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_tokens 
		SET attempt_count = attempt_count + 1,
			updated_at = NOW()
		WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment OTP attempt: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}


// MarkCancelled marks an OTP token as cancelled
func (r *otpTokenRepository) MarkCancelled(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_tokens 
		SET status = 'cancelled',
			updated_at = NOW()
		WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP token as cancelled: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}


// MarkExpired marks an OTP token as expired
func (r *otpTokenRepository) MarkExpired(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_tokens 
		SET status = 'expired',
			updated_at = NOW()
		WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP token as expired: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}


// MarkUsed marks an OTP token as used
func (r *otpTokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_tokens 
		SET status = 'used',
			used_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP token as used: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return entities.ErrOTPTokenNotFound
	}
	
	return nil
}


// ValidateAndMarkUsed validates an OTP code and marks it as used if valid
func (r *otpTokenRepository) ValidateAndMarkUsed(ctx context.Context, phone string, code string, purpose entities.OTPPurpose) (*entities.OTPToken, error) {
	// First, find the active token
	query := `
		SELECT id, phone, email, code, purpose, status, expires_at, 
			   used_at, created_at, updated_at, attempt_count, max_attempts, user_agent, ip_address
		FROM otp_tokens 
		WHERE phone = $1 
		  AND code = $2
		  AND purpose = $3
		  AND status = 'pending'
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1`
	
	var token entities.OTPToken
	err := r.db.GetContext(ctx, &token, query, phone, code, string(purpose))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrOTPTokenNotFound
		}
		return nil, fmt.Errorf("failed to validate OTP token: %w", err)
	}
	
	// Check if max attempts exceeded
	if token.AttemptCount >= token.MaxAttempts {
		return nil, fmt.Errorf("maximum attempts exceeded")
	}
	
	// Mark as used
	err = r.MarkUsed(ctx, token.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}
	
	// Update the token status
	token.Status = entities.OTPStatusUsed
	
	return &token, nil
}

