package entities

import (
	"time"

	"github.com/google/uuid"
)

// Organizer represents a top-level entity for multi-tenancy
type Organizer struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Slug        string                 `json:"slug" db:"slug"`
	Email       string                 `json:"email" db:"email"`
	Phone       *string                `json:"phone,omitempty" db:"phone"`
	Website     *string                `json:"website,omitempty" db:"website"`
	LogoURL     *string                `json:"logo_url,omitempty" db:"logo_url"`
	Description *string                `json:"description,omitempty" db:"description"`
	Settings    map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	IsActive    bool                   `json:"is_active" db:"is_active"`
}

// NewOrganizer creates a new organizer with default values
func NewOrganizer(name, slug, email string) *Organizer {
	now := time.Now()
	return &Organizer{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		Email:     email,
		Settings:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}
}

// Validate performs business rule validation for the organizer
func (o *Organizer) Validate() error {
	if o.Name == "" {
		return NewValidationError("name", "name is required")
	}
	if o.Slug == "" {
		return NewValidationError("slug", "slug is required")
	}
	if o.Email == "" {
		return NewValidationError("email", "email is required")
	}
	if len(o.Name) > 255 {
		return NewValidationError("name", "name must be 255 characters or less")
	}
	if len(o.Slug) > 100 {
		return NewValidationError("slug", "slug must be 100 characters or less")
	}
	return nil
}

// UpdateSettings updates the organizer settings
func (o *Organizer) UpdateSettings(settings map[string]interface{}) {
	if o.Settings == nil {
		o.Settings = make(map[string]interface{})
	}
	for key, value := range settings {
		o.Settings[key] = value
	}
	o.UpdatedAt = time.Now()
}

// SetLogo sets the logo URL for the organizer
func (o *Organizer) SetLogo(logoURL string) {
	o.LogoURL = &logoURL
	o.UpdatedAt = time.Now()
}

// Deactivate marks the organizer as inactive
func (o *Organizer) Deactivate() {
	o.IsActive = false
	o.UpdatedAt = time.Now()
}

// Activate marks the organizer as active
func (o *Organizer) Activate() {
	o.IsActive = true
	o.UpdatedAt = time.Now()
}

