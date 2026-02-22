package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider represents the authentication provider type
type AuthProvider string

const (
	AuthProviderEmail AuthProvider = "email"
	AuthProviderMoMo  AuthProvider = "momo"
)

// User represents a unified user management entity for all authentication types
type User struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	Email          *string                `json:"email,omitempty" db:"email"`
	Phone          *string                `json:"phone,omitempty" db:"phone_number"`
	Password       *string                `json:"-" db:"password"`           // Added missing field
	PasswordHash   *string                `json:"-" db:"password_hash"`
	FirstName      *string                `json:"first_name,omitempty" db:"first_name"`
	LastName       *string                `json:"last_name,omitempty" db:"last_name"`
	AuthProvider   AuthProvider           `json:"auth_provider" db:"auth_provider"`
	MoMoID         *string                `json:"momo_id,omitempty" db:"momo_id"`
	EmailVerified  bool                   `json:"email_verified" db:"email_verified"`
	PhoneVerified  bool                   `json:"phone_verified" db:"phone_verified"`
	IsActive       bool                   `json:"is_active" db:"is_active"`
	LastLogin      *time.Time             `json:"last_login,omitempty" db:"last_login"`
	Settings       map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user (generic constructor)
func NewUser(email, firstName, lastName string) *User {
	return NewEmailUser(email, firstName, lastName)
}

// NewEmailUser creates a new user with email authentication
func NewEmailUser(email, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		Email:         &email,
		FirstName:     &firstName,
		LastName:      &lastName,
		AuthProvider:  AuthProviderEmail,
		EmailVerified: false,
		PhoneVerified: false,
		IsActive:      true,
		Settings:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewEmailUserWithPassword creates a new user with email authentication and password
func NewEmailUserWithPassword(email, passwordHash, firstName, lastName, phoneNumber string) *User {
	now := time.Now()
	user := &User{
		ID:            uuid.New(),
		Email:         &email,
		PasswordHash:  &passwordHash,
		FirstName:     &firstName,
		LastName:      &lastName,
		AuthProvider:  AuthProviderEmail,
		EmailVerified: false,
		PhoneVerified: false,
		IsActive:      true,
		Settings:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	if phoneNumber != "" {
		user.Phone = &phoneNumber
	}
	
	return user
}

// NewMoMoUser creates a new user with Mobile Money authentication
func NewMoMoUser(phone, momoID string) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		Phone:         &phone,
		MoMoID:        &momoID,
		AuthProvider:  AuthProviderMoMo,
		EmailVerified: false,
		PhoneVerified: true, // MoMo users are phone verified by default
		IsActive:      true,
		Settings:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// GetDisplayName returns the user's display name
func (u *User) GetDisplayName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.Email != nil {
		return *u.Email
	}
	if u.Phone != nil {
		return *u.Phone
	}
	return "User"
}

// GetContactInfo returns the primary contact information
func (u *User) GetContactInfo() string {
	if u.Email != nil {
		return *u.Email
	}
	if u.Phone != nil {
		return *u.Phone
	}
	return ""
}

// IsEmailUser checks if the user uses email authentication
func (u *User) IsEmailUser() bool {
	return u.AuthProvider == AuthProviderEmail
}

// IsMoMoUser checks if the user uses Mobile Money authentication
func (u *User) IsMoMoUser() bool {
	return u.AuthProvider == AuthProviderMoMo
}

// CanLogin checks if the user can login
func (u *User) CanLogin() bool {
	return u.IsActive
}

// MarkEmailVerified marks the user's email as verified
func (u *User) MarkEmailVerified() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}

// MarkPhoneVerified marks the user's phone as verified
func (u *User) MarkPhoneVerified() {
	u.PhoneVerified = true
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Validate validates the user
func (u *User) Validate() error {
	if u.AuthProvider == AuthProviderEmail {
		if u.Email == nil || *u.Email == "" {
			return NewValidationError("email", "email is required for email authentication")
		}
	}
	
	if u.AuthProvider == AuthProviderMoMo {
		if u.Phone == nil || *u.Phone == "" {
			return NewValidationError("phone", "phone is required for MoMo authentication")
		}
		if u.MoMoID == nil || *u.MoMoID == "" {
			return NewValidationError("momo_id", "MoMo ID is required for MoMo authentication")
		}
	}
	
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName *string) {
	if firstName != nil {
		u.FirstName = firstName
	}
	if lastName != nil {
		u.LastName = lastName
	}
	u.UpdatedAt = time.Now()
}

// SetPassword sets the user's password hash
func (u *User) SetPassword(passwordHash string) {
	u.PasswordHash = &passwordHash
	u.UpdatedAt = time.Now()
}

// HasPassword checks if the user has a password set
func (u *User) HasPassword() bool {
	return u.PasswordHash != nil && *u.PasswordHash != ""
}

