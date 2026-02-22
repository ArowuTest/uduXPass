package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/pkg/jwt"
	"github.com/uduxpass/backend/pkg/security"
)

// AdminAuthService handles admin authentication operations
type AdminAuthService struct {
	adminRepo    repositories.AdminUserRepository
	jwtService   jwt.Service
	passwordSvc  security.PasswordService
	lockoutDuration time.Duration
	maxLoginAttempts int
}

// NewAdminAuthService creates a new admin authentication service
func NewAdminAuthService(
	adminRepo repositories.AdminUserRepository,
	jwtService jwt.Service,
	passwordSvc security.PasswordService,
) *AdminAuthService {
	return &AdminAuthService{
		adminRepo:        adminRepo,
		jwtService:       jwtService,
		passwordSvc:      passwordSvc,
		lockoutDuration:  30 * time.Minute, // 30 minutes lockout
		maxLoginAttempts: 5,
	}
}

// AdminLoginRequest represents an admin login request
type AdminLoginRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"` // Alias for email
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// GetEmail returns email from either field
func (r *AdminLoginRequest) GetEmail() string {
	if r.Email != "" {
		return r.Email
	}
	return r.Username
}

// AdminLoginResponse represents an admin login response
type AdminLoginResponse struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	Admin        *entities.AdminUser   `json:"admin"`
	ExpiresIn    int64                `json:"expires_in"`
	Permissions  []string             `json:"permissions"`
}

// CreateAdminRequest represents a create admin request
type CreateAdminRequest struct {
	Email       string                     `json:"email" validate:"required,email"`
	Password    string                     `json:"password" validate:"required,min=8"`
	FirstName   string                     `json:"first_name" validate:"required"`
	LastName    string                     `json:"last_name" validate:"required"`
	Role        entities.AdminRole         `json:"role" validate:"required"`
	Permissions []entities.AdminPermission `json:"permissions,omitempty"`
}

// CreateAdminResponse represents a create admin response
type CreateAdminResponse struct {
	Admin   *entities.AdminUser `json:"admin"`
	Message string             `json:"message"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	AdminID         uuid.UUID `json:"admin_id" validate:"required"`
	CurrentPassword string    `json:"current_password" validate:"required"`
	NewPassword     string    `json:"new_password" validate:"required,min=8"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Login authenticates an admin user
func (s *AdminAuthService) Login(ctx context.Context, req *AdminLoginRequest) (*AdminLoginResponse, error) {
	// Get email from either email or username field
	email := req.GetEmail()
	if email == "" {
		return nil, entities.ErrInvalidCredentials
	}
	
	// Get admin by email
	admin, err := s.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, entities.ErrAdminNotFound) {
			return nil, entities.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	// Check account status
	if !admin.IsActive {
		return nil, entities.ErrAccountDeactivated
	}

	// Check if account is locked
	if admin.LockedUntil != nil && time.Now().Before(*admin.LockedUntil) {
		return nil, entities.ErrAccountLocked
	}

	// Unlock account if lockout period has expired
	if admin.LockedUntil != nil && time.Now().After(*admin.LockedUntil) {
		admin.LockedUntil = nil
		admin.LoginAttempts = 0
		_ = s.adminRepo.Update(ctx, admin)
	}

	// Verify password
	valid, err := s.passwordSvc.VerifyPassword(req.Password, admin.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	if !valid {
		// Increment login attempts
		admin.LoginAttempts++

		// Lock account if max attempts reached
		if admin.LoginAttempts >= s.maxLoginAttempts {
			lockoutUntil := time.Now().Add(s.lockoutDuration)
			admin.LockedUntil = &lockoutUntil
		}

		_ = s.adminRepo.Update(ctx, admin)
		return nil, entities.ErrInvalidCredentials
	}

	// Reset login attempts on successful login
	admin.LoginAttempts = 0
	now := time.Now()
	admin.LastLogin = &now
	admin.LockedUntil = nil
	_ = s.adminRepo.Update(ctx, admin)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(admin.ID, string(admin.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(admin.ID, string(admin.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        admin,
		ExpiresIn:    3600, // 1 hour
		Permissions:  admin.GetPermissionStrings(),
	}, nil
}

// CreateAdmin creates a new admin user
func (s *AdminAuthService) CreateAdmin(ctx context.Context, req *CreateAdminRequest) (*CreateAdminResponse, error) {
	// Check if admin already exists
	existingAdmin, err := s.adminRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, entities.ErrAdminNotFound) {
		return nil, fmt.Errorf("failed to check existing admin: %w", err)
	}

	if existingAdmin != nil {
		return nil, entities.ErrAdminAlreadyExists
	}

	// Validate password strength
	if err := s.passwordSvc.ValidatePasswordStrength(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash password
	hashedPassword, err := s.passwordSvc.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	admin := entities.NewAdminUser(req.Email, req.FirstName, req.LastName, req.Role, nil)
	
	// Set the hashed password
	admin.PasswordHash = hashedPassword

	// Set permissions if provided
	if len(req.Permissions) > 0 {
		admin.Permissions = req.Permissions
	}

	// Save admin
	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	return &CreateAdminResponse{
		Admin:   admin,
		Message: "Admin created successfully",
	}, nil
}

// ChangePassword changes an admin's password
func (s *AdminAuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	// Get admin
	admin, err := s.adminRepo.GetByID(ctx, req.AdminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	// Verify current password
	valid, err := s.passwordSvc.VerifyPassword(req.CurrentPassword, admin.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to verify current password: %w", err)
	}

	if !valid {
		return entities.ErrInvalidCredentials
	}

	// Validate new password strength
	if err := s.passwordSvc.ValidatePasswordStrength(req.NewPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.passwordSvc.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	admin.PasswordHash = hashedPassword
	now := time.Now()
	admin.UpdatedAt = now

	return s.adminRepo.Update(ctx, admin)
}

// RefreshToken refreshes an admin's access token
func (s *AdminAuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AdminLoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Parse admin ID from claims
	adminID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid admin ID in token: %w", err)
	}

	// Get admin
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	// Check admin status
	if !admin.IsActive {
		return nil, entities.ErrAccountDeactivated
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(admin.ID, string(admin.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(admin.ID, string(admin.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Admin:        admin,
		ExpiresIn:    3600, // 1 hour
		Permissions:  admin.GetPermissionStrings(),
	}, nil
}

// ValidateToken validates an admin JWT token
func (s *AdminAuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtService.ValidateAccessToken(tokenString)
}

// GetAdminByToken gets an admin by JWT token
func (s *AdminAuthService) GetAdminByToken(ctx context.Context, tokenString string) (*entities.AdminUser, error) {
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	// Parse admin ID from claims
	adminID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid admin ID in token: %w", err)
	}
	
	return s.adminRepo.GetByID(ctx, adminID)
}

// ResetPassword resets an admin's password (admin-only operation)
func (s *AdminAuthService) ResetPassword(ctx context.Context, adminID uuid.UUID, newPassword string) error {
	// Validate password strength
	if err := s.passwordSvc.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash password
	hashedPassword, err := s.passwordSvc.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Get admin
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	// Update password
	admin.PasswordHash = hashedPassword
	now := time.Now()
	admin.UpdatedAt = now
	admin.LoginAttempts = 0
	admin.LockedUntil = nil

	return s.adminRepo.Update(ctx, admin)
}

// LockAdmin locks an admin account
func (s *AdminAuthService) LockAdmin(ctx context.Context, adminID uuid.UUID, duration time.Duration) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	if duration > 0 {
		lockoutUntil := time.Now().Add(duration)
		admin.LockedUntil = &lockoutUntil
	} else {
		// Permanent lock - set far future date
		lockoutUntil := time.Now().Add(100 * 365 * 24 * time.Hour)
		admin.LockedUntil = &lockoutUntil
	}

	return s.adminRepo.Update(ctx, admin)
}

// UnlockAdmin unlocks an admin account
func (s *AdminAuthService) UnlockAdmin(ctx context.Context, adminID uuid.UUID) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	admin.LockedUntil = nil
	admin.LoginAttempts = 0

	return s.adminRepo.Update(ctx, admin)
}

// DeactivateAdmin deactivates an admin account
func (s *AdminAuthService) DeactivateAdmin(ctx context.Context, adminID uuid.UUID) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	admin.IsActive = false

	return s.adminRepo.Update(ctx, admin)
}

// ActivateAdmin activates an admin account
func (s *AdminAuthService) ActivateAdmin(ctx context.Context, adminID uuid.UUID) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	admin.IsActive = true
	admin.LockedUntil = nil
	admin.LoginAttempts = 0

	return s.adminRepo.Update(ctx, admin)
}


// UpdateAdminRequest represents an update admin request
type UpdateAdminRequest struct {
	AdminID     uuid.UUID                  `json:"admin_id" validate:"required"`
	FirstName   *string                    `json:"first_name,omitempty"`
	LastName    *string                    `json:"last_name,omitempty"`
	Email       *string                    `json:"email,omitempty"`
	Role        *entities.AdminRole        `json:"role,omitempty"`
	Permissions []entities.AdminPermission `json:"permissions,omitempty"`
	IsActive    *bool                      `json:"is_active,omitempty"`
}

// UpdateAdminResponse represents an update admin response
type UpdateAdminResponse struct {
	Admin   *entities.AdminUser `json:"admin"`
	Message string             `json:"message"`
}

// ListAdminsRequest represents a list admins request
type ListAdminsRequest struct {
	Page   int                    `json:"page,omitempty"`
	Limit  int                    `json:"limit,omitempty"`
	Role   *entities.AdminRole    `json:"role,omitempty"`
	Status *entities.AdminStatus  `json:"status,omitempty"`
	Search string                 `json:"search,omitempty"`
}

// ListAdminsResponse represents a list admins response
type ListAdminsResponse struct {
	Admins     []*entities.AdminUser `json:"admins"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// AdminStatsResponse represents admin statistics
type AdminStatsResponse struct {
	TotalAdmins    int64                        `json:"total_admins"`
	ActiveAdmins   int64                        `json:"active_admins"`
	InactiveAdmins int64                        `json:"inactive_admins"`
	AdminsByRole   map[entities.AdminRole]int64 `json:"admins_by_role"`
}

// GetAdmin retrieves an admin by ID
func (s *AdminAuthService) GetAdmin(ctx context.Context, adminID uuid.UUID) (*entities.AdminUser, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}
	
	return admin, nil
}

// UpdateAdmin updates an admin user
func (s *AdminAuthService) UpdateAdmin(ctx context.Context, req *UpdateAdminRequest) (*UpdateAdminResponse, error) {
	// Get existing admin
	admin, err := s.adminRepo.GetByID(ctx, req.AdminID)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	// Update fields if provided
	if req.FirstName != nil {
		admin.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		admin.LastName = *req.LastName
	}
	if req.Email != nil {
		admin.Email = *req.Email
	}
	if req.Role != nil {
		admin.Role = *req.Role
	}
	if req.Permissions != nil {
		admin.Permissions = req.Permissions
	}
	if req.IsActive != nil {
		admin.IsActive = *req.IsActive
	}

	admin.UpdatedAt = time.Now()

	// Update in repository
	err = s.adminRepo.Update(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("failed to update admin: %w", err)
	}

	return &UpdateAdminResponse{
		Admin:   admin,
		Message: "Admin updated successfully",
	}, nil
}

// ListAdmins retrieves a list of admins with filtering
func (s *AdminAuthService) ListAdmins(ctx context.Context, req *ListAdminsRequest) (*ListAdminsResponse, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get admins
	admins, err := s.adminRepo.List(ctx, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list admins: %w", err)
	}

	// Get total count
	total, err := s.adminRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count admins: %w", err)
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &ListAdminsResponse{
		Admins:     admins,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteAdmin soft deletes an admin user
func (s *AdminAuthService) DeleteAdmin(ctx context.Context, adminID uuid.UUID) error {
	// Check if admin exists
	_, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	// Delete admin
	err = s.adminRepo.Delete(ctx, adminID)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	return nil
}

// GetAdminStats retrieves admin statistics
func (s *AdminAuthService) GetAdminStats(ctx context.Context) (*AdminStatsResponse, error) {
	// Get total count
	total, err := s.adminRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total admin count: %w", err)
	}

	// Get active count
	active, err := s.adminRepo.CountByStatus(ctx, entities.AdminStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get active admin count: %w", err)
	}

	// Get inactive count
	inactive := total - active

	// Get counts by role
	adminsByRole := make(map[entities.AdminRole]int64)
	roles := []entities.AdminRole{
		entities.AdminRoleSuperAdmin,
		entities.AdminRoleEventManager,
		entities.AdminRoleSupport,
		entities.AdminRoleAnalyst,
	}

	for _, role := range roles {
		count, err := s.adminRepo.CountByRole(ctx, role)
		if err != nil {
			return nil, fmt.Errorf("failed to get count for role %s: %w", role, err)
		}
		adminsByRole[role] = count
	}

	return &AdminStatsResponse{
		TotalAdmins:    total,
		ActiveAdmins:   active,
		InactiveAdmins: inactive,
		AdminsByRole:   adminsByRole,
	}, nil
}

