package entities

import (
	"time"

	"github.com/google/uuid"
)

// Tour represents a grouping mechanism for related events
type Tour struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	OrganizerID  uuid.UUID              `json:"organizer_id" db:"organizer_id"`
	Name         string                 `json:"name" db:"name"`
	Slug         string                 `json:"slug" db:"slug"`
	Description  *string                `json:"description,omitempty" db:"description"`
	ArtistName   string                 `json:"artist_name" db:"artist_name"`
	TourImageURL *string                `json:"tour_image_url,omitempty" db:"tour_image_url"`
	StartDate    *time.Time             `json:"start_date,omitempty" db:"start_date"`
	EndDate      *time.Time             `json:"end_date,omitempty" db:"end_date"`
	Settings     map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	IsActive     bool                   `json:"is_active" db:"is_active"`
}

// NewTour creates a new tour with default values
func NewTour(organizerID uuid.UUID, name, slug, artistName string) *Tour {
	now := time.Now()
	return &Tour{
		ID:          uuid.New(),
		OrganizerID: organizerID,
		Name:        name,
		Slug:        slug,
		ArtistName:  artistName,
		Settings:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}
}

// Validate performs business rule validation for the tour
func (t *Tour) Validate() error {
	if t.Name == "" {
		return NewValidationError("name", "name is required")
	}
	if t.Slug == "" {
		return NewValidationError("slug", "slug is required")
	}
	if t.ArtistName == "" {
		return NewValidationError("artist_name", "artist name is required")
	}
	if len(t.Name) > 255 {
		return NewValidationError("name", "name must be 255 characters or less")
	}
	if len(t.Slug) > 100 {
		return NewValidationError("slug", "slug must be 100 characters or less")
	}
	if len(t.ArtistName) > 255 {
		return NewValidationError("artist_name", "artist name must be 255 characters or less")
	}
	
	// Validate date range if both dates are provided
	if t.StartDate != nil && t.EndDate != nil && t.StartDate.After(*t.EndDate) {
		return NewValidationError("date_range", "start date must be before end date")
	}
	
	return nil
}

// SetDates sets the tour start and end dates
func (t *Tour) SetDates(startDate, endDate time.Time) error {
	if startDate.After(endDate) {
		return NewValidationError("date_range", "start date must be before end date")
	}
	t.StartDate = &startDate
	t.EndDate = &endDate
	t.UpdatedAt = time.Now()
	return nil
}

// SetImage sets the tour image URL
func (t *Tour) SetImage(imageURL string) {
	t.TourImageURL = &imageURL
	t.UpdatedAt = time.Now()
}

// UpdateSettings updates the tour settings
func (t *Tour) UpdateSettings(settings map[string]interface{}) {
	if t.Settings == nil {
		t.Settings = make(map[string]interface{})
	}
	for key, value := range settings {
		t.Settings[key] = value
	}
	t.UpdatedAt = time.Now()
}

// IsActive checks if the tour is currently active
func (t *Tour) IsActiveTour() bool {
	return t.IsActive
}

// Deactivate marks the tour as inactive
func (t *Tour) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

// Activate marks the tour as active
func (t *Tour) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

// IsOngoing checks if the tour is currently ongoing
func (t *Tour) IsOngoing() bool {
	if t.StartDate == nil || t.EndDate == nil {
		return false
	}
	now := time.Now()
	return now.After(*t.StartDate) && now.Before(*t.EndDate)
}

// IsUpcoming checks if the tour is upcoming
func (t *Tour) IsUpcoming() bool {
	if t.StartDate == nil {
		return false
	}
	return time.Now().Before(*t.StartDate)
}

// IsCompleted checks if the tour is completed
func (t *Tour) IsCompleted() bool {
	if t.EndDate == nil {
		return false
	}
	return time.Now().After(*t.EndDate)
}

