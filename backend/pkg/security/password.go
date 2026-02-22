package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService defines the interface for password hashing and verification
type PasswordService interface {
	// HashPassword hashes a plain text password
	HashPassword(password string) (string, error)
	
	// VerifyPassword verifies a plain text password against a hashed password
	// Returns (valid bool, error)
	VerifyPassword(password, hashedPassword string) (bool, error)
	
	// ValidatePasswordStrength validates password strength requirements
	ValidatePasswordStrength(password string) error
}

// BcryptConfig holds configuration for bcrypt password service
type BcryptConfig struct {
	Cost int // bcrypt cost factor (default: 10)
}

// BcryptPasswordService implements PasswordService using bcrypt
type BcryptPasswordService struct {
	cost int
}

// NewBcryptPasswordService creates a new bcrypt password service
func NewBcryptPasswordService(config BcryptConfig) PasswordService {
	cost := config.Cost
	if cost == 0 {
		cost = bcrypt.DefaultCost // Use default cost if not specified
	}
	return &BcryptPasswordService{
		cost: cost,
	}
}

// HashPassword hashes a plain text password using bcrypt
func (s *BcryptPasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
// Returns (valid bool, error)
func (s *BcryptPasswordService) VerifyPassword(password, hashedPassword string) (bool, error) {
	if hashedPassword == "" {
		return false, fmt.Errorf("hashed password cannot be empty")
	}
	if password == "" {
		return false, fmt.Errorf("password cannot be empty")
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil // Invalid password, but not an error
		}
		return false, fmt.Errorf("failed to verify password: %w", err)
	}
	
	return true, nil
}

// ValidatePasswordStrength validates password strength requirements
func (s *BcryptPasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters long")
	}
	
	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	
	return nil
}
