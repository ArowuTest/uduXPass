package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/pkg/jwt"
	"github.com/uduxpass/backend/pkg/security"
)

// AuthService handles user authentication
type AuthService struct {
	userRepo    repositories.UserRepository
	otpRepo     repositories.OTPTokenRepository
	jwtService  jwt.Service
	passwordSvc security.PasswordService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repositories.UserRepository,
	otpRepo repositories.OTPTokenRepository,
	jwtService jwt.Service,
	passwordSvc security.PasswordService,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		otpRepo:     otpRepo,
		jwtService:  jwtService,
		passwordSvc: passwordSvc,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"firstName" validate:"required"` // Accept camelCase from frontend
	LastName    string `json:"lastName" validate:"required"`  // Accept camelCase from frontend
	PhoneNumber string `json:"phone_number"`
	Phone       string `json:"phone"` // Alias for phone_number for API compatibility
}

// GetPhone returns the phone number from either field
func (r *RegisterRequest) GetPhone() string {
	if r.PhoneNumber != "" {
		return r.PhoneNumber
	}
	return r.Phone
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    int64          `json:"expires_in"`
	User         *entities.User `json:"user"`
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != entities.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, entities.NewValidationError("email", "User with this email already exists")
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

	// Get phone number from either field
	phoneNumber := req.GetPhone()
	if phoneNumber == "" {
		return nil, entities.NewValidationError("phone", "Phone number is required")
	}

	// Create user
	user := entities.NewEmailUserWithPassword(
		req.Email,
		hashedPassword,
		req.FirstName,
		req.LastName,
		phoneNumber,
	)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         user,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == entities.ErrUserNotFound {
			return nil, entities.NewValidationError("email", "Invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, entities.NewValidationError("account", "Account is deactivated")
	}

	// Verify password
	if user.PasswordHash == nil {
		return nil, entities.NewValidationError("password", "Invalid email or password")
	}
	
	valid, err := s.passwordSvc.VerifyPassword(req.Password, *user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}
	
	if !valid {
		// Increment failed attempts
		s.userRepo.IncrementFailedAttempts(ctx, user.ID)
		return nil, entities.NewValidationError("password", "Invalid email or password")
	}

	// Reset failed attempts on successful login
	if err := s.userRepo.ResetFailedAttempts(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to reset failed attempts for user %s: %v\n", user.ID, err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login for user %s: %v\n", user.ID, err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         user,
	}, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, entities.NewValidationError("account", "Account is deactivated")
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         user,
	}, nil
}

// SendOTP sends an OTP to the user
func (s *AuthService) SendOTP(ctx context.Context, email string) error {
	// Get user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == entities.ErrUserNotFound {
			return entities.NewValidationError("email", "User not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Generate OTP
	otp := s.generateOTP()
	
	// Create OTP token
	phone := "unknown"
	if user.Phone != nil && *user.Phone != "" {
		phone = *user.Phone
	}
	
	otpToken := &entities.OTPToken{
		ID:        uuid.New(),
		Phone:     phone,
		Email:     user.Email,
		Code:      otp,
		Purpose:   entities.OTPPurposePasswordReset,
		Status:    entities.OTPStatusPending,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minutes
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AttemptCount: 0,
		MaxAttempts:  3,
	}

	// Save OTP token
	if err := s.otpRepo.Create(ctx, otpToken); err != nil {
		return fmt.Errorf("failed to save OTP token: %w", err)
	}

	// TODO: Send OTP via email/SMS
	fmt.Printf("OTP for %s: %s\n", email, otp)

	return nil
}

// VerifyOTP verifies an OTP
func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	// Get OTP token by code
	otpToken, err := s.otpRepo.GetByCode(ctx, otp)
	if err != nil {
		return entities.NewValidationError("otp", "Invalid or expired OTP")
	}

	// Validate OTP token
	if !otpToken.IsValid() {
		return entities.NewValidationError("otp", "Invalid or expired OTP")
	}

	// Check if the OTP is for password reset purpose
	if otpToken.Purpose != entities.OTPPurposePasswordReset {
		return entities.NewValidationError("otp", "Invalid OTP purpose")
	}

	// Check if the email matches (if OTP has email)
	if otpToken.Email != nil && *otpToken.Email != email {
		return entities.NewValidationError("otp", "OTP does not match the provided email")
	}

	// Mark OTP as used
	otpToken.MarkAsUsed()
	if err := s.otpRepo.Update(ctx, otpToken); err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Check if OTP is already used
	if otpToken.UsedAt != nil {
		return entities.NewValidationError("otp", "OTP has already been used")
	}

	// Mark OTP as used
	now := time.Now()
	otpToken.UsedAt = &now
	if err := s.otpRepo.Update(ctx, otpToken); err != nil {
		return fmt.Errorf("failed to update OTP token: %w", err)
	}

	return nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	// Verify OTP first
	if err := s.VerifyOTP(ctx, email, otp); err != nil {
		return err
	}

	// Get user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate new password
	if err := s.passwordSvc.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.passwordSvc.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// generateOTP generates a 6-digit OTP
func (s *AuthService) generateOTP() string {
	// Simple OTP generation - in production, use crypto/rand
	return fmt.Sprintf("%06d", time.Now().Unix()%1000000)
}
