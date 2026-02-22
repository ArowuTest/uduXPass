package entities

import (
	"time"

	"github.com/google/uuid"
)

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusOnSale    EventStatus = "on_sale"
	EventStatusSoldOut   EventStatus = "sold_out"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

// Event represents an individual ticketed occasion
type Event struct {
	ID              uuid.UUID              `json:"id" db:"id"`
	OrganizerID     *uuid.UUID             `json:"organizer_id,omitempty" db:"organizer_id"`
	CategoryID      *uuid.UUID             `json:"category_id,omitempty" db:"category_id"`
	Name            string                 `json:"name" db:"name"`
	Slug            string                 `json:"slug" db:"slug"`
	Description     *string                `json:"description,omitempty" db:"description"`
	EventDate       time.Time              `json:"event_date" db:"event_date"`
	DoorsOpen       *time.Time             `json:"doors_open,omitempty" db:"doors_open"`
	VenueName       string                 `json:"venue_name" db:"venue_name"`
	VenueAddress    string                 `json:"venue_address" db:"venue_address"`
	VenueCity       string                 `json:"venue_city" db:"venue_city"`
	VenueState      *string                `json:"venue_state,omitempty" db:"venue_state"`
	VenueCountry    *string                `json:"venue_country,omitempty" db:"venue_country"`
	VenueCapacity   *int                   `json:"venue_capacity,omitempty" db:"venue_capacity"`
	EventImageURL   *string                `json:"event_image_url,omitempty" db:"event_image_url"`
	Status          EventStatus            `json:"status" db:"status"`
	SaleStart       *time.Time             `json:"sale_start,omitempty" db:"sale_start"`
	SaleEnd         *time.Time             `json:"sale_end,omitempty" db:"sale_end"`
	Settings        JSONB                  `json:"settings" db:"settings"`
	Currency        *string                `json:"currency,omitempty" db:"currency"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	IsActive        bool                   `json:"is_active" db:"is_active"`

	// Relations
	Organizer   *Organizer    `json:"organizer,omitempty"`
	Tour        *Tour         `json:"tour,omitempty"`
	TicketTiers []TicketTier  `json:"ticket_tiers,omitempty"`
	Orders      []Order       `json:"orders,omitempty"`
	Tickets     []Ticket      `json:"tickets,omitempty"`
}

// NewEvent creates a new event with default values
func NewEvent(organizerID uuid.UUID, name, slug string, eventDate time.Time, venueName, venueAddress, venueCity, venueCountry string) *Event {
	orgID := organizerID
	vcountry := venueCountry
	return &Event{
		ID:           uuid.New(),
		OrganizerID:  &orgID,
		Name:         name,
		Slug:         slug,
		EventDate:    eventDate,
		VenueName:    venueName,
		VenueAddress: venueAddress,
		VenueCity:    venueCity,
		VenueCountry: &vcountry,
		Status:       EventStatusDraft,
		Settings:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}
}

// IsOnSale checks if the event is currently on sale
func (e *Event) IsOnSale() bool {
	now := time.Now()
	
	if e.Status != EventStatusPublished && e.Status != EventStatusOnSale {
		return false
	}
	
	if e.SaleStart != nil && now.Before(*e.SaleStart) {
		return false
	}
	
	if e.SaleEnd != nil && now.After(*e.SaleEnd) {
		return false
	}
	
	return true
}

// CanBeEdited checks if the event can be edited
func (e *Event) CanBeEdited() bool {
	return e.Status == EventStatusDraft
}

// Publish changes the event status to published
func (e *Event) Publish() error {
	if e.Status != EventStatusDraft {
		return NewBusinessRuleError("business_rule", "only draft events can be published", nil)
	}
	e.Status = EventStatusPublished
	e.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the event
func (e *Event) Cancel() error {
	if e.Status == EventStatusCompleted {
		return NewBusinessRuleError("business_rule", "completed events cannot be cancelled", nil)
	}
	e.Status = EventStatusCancelled
	e.UpdatedAt = time.Now()
	return nil
}

// Complete marks the event as completed
func (e *Event) Complete() error {
	if e.Status == EventStatusCancelled {
		return NewBusinessRuleError("business_rule", "cancelled events cannot be completed", nil)
	}
	e.Status = EventStatusCompleted
	e.UpdatedAt = time.Now()
	return nil
}

// SetSoldOut marks the event as sold out
func (e *Event) SetSoldOut() error {
	if !e.IsOnSale() {
		return NewBusinessRuleError("business_rule", "only events on sale can be marked as sold out", nil)
	}
	e.Status = EventStatusSoldOut
	e.UpdatedAt = time.Now()
	return nil
}

// SetOnSale marks the event as on sale
func (e *Event) SetOnSale() error {
	if e.Status != EventStatusPublished {
		return NewBusinessRuleError("business_rule", "only published events can be put on sale", nil)
	}
	e.Status = EventStatusOnSale
	e.UpdatedAt = time.Now()
	return nil
}

// Validate validates the event
func (e *Event) Validate() error {
	if e.Name == "" {
		return NewValidationError("name", "event name is required")
	}
	
	if e.Slug == "" {
		return NewValidationError("slug", "event slug is required")
	}
	
	if e.VenueName == "" {
		return NewValidationError("venue_name", "venue name is required")
	}
	
	if e.VenueAddress == "" {
		return NewValidationError("venue_address", "venue address is required")
	}
	
	if e.VenueCity == "" {
		return NewValidationError("venue_city", "venue city is required")
	}
	
	if e.VenueCountry == nil || *e.VenueCountry == "" {
		return NewValidationError("venue_country", "venue country is required")
	}
	
	if e.EventDate.IsZero() {
		return NewValidationError("event_date", "event date is required")
	}
	
	if e.EventDate.Before(time.Now()) {
		return NewValidationError("event_date", "event date cannot be in the past")
	}
	
	return nil
}

// GetAvailableTickets returns the total number of available tickets
func (e *Event) GetAvailableTickets() int {
	total := 0
	for _, tier := range e.TicketTiers {
		if tier.IsActive {
			// Calculate available tickets (quota minus sold)
			total += (tier.Quota - tier.Sold)
		}
	}
	return total
}

// GetTotalCapacity returns the total capacity of all ticket tiers
func (e *Event) GetTotalCapacity() int {
	total := 0
	for _, tier := range e.TicketTiers {
		if tier.IsActive {
			total += tier.Quota
		}
	}
	return total
}

// GetSoldTickets returns the total number of sold tickets
func (e *Event) GetSoldTickets() int {
	// This would need to be calculated from actual orders/tickets
	// For now, return 0 as placeholder
	return 0
}

// GetRevenue returns the total revenue from ticket sales
func (e *Event) GetRevenue() float64 {
	total := 0.0
	for _, order := range e.Orders {
		if order.Status == OrderStatusPaid {
			total += order.TotalAmount
		}
	}
	return total
}

// SetTour sets the tour ID for the event (deprecated - tour_id not in schema)
// func (e *Event) SetTour(tourID uuid.UUID) {
// 	e.TourID = &tourID
// 	e.UpdatedAt = time.Now()
// }

// SetVenueLocation sets the venue coordinates (deprecated - coordinates not in schema)
// func (e *Event) SetVenueLocation(latitude, longitude float64) {
// 	e.VenueLatitude = &latitude
// 	e.VenueLongitude = &longitude
// 	e.UpdatedAt = time.Now()
// }

// SetImage sets the event image URL
func (e *Event) SetImage(imageURL string) {
	e.EventImageURL = &imageURL
	e.UpdatedAt = time.Now()
}

// SetSalePeriod sets the sale start and end times
func (e *Event) SetSalePeriod(saleStart, saleEnd time.Time) error {
	if saleEnd.Before(saleStart) {
		return NewValidationError("sale_period", "sale end time must be after sale start time")
	}
	
	e.SaleStart = &saleStart
	e.SaleEnd = &saleEnd
	e.UpdatedAt = time.Now()
	
	return nil
}

// StartSale starts ticket sales for the event
func (e *Event) StartSale() error {
	if e.Status == EventStatusOnSale {
		return NewValidationError("status", "event is already on sale")
	}
	
	if e.Status == EventStatusCancelled {
		return NewValidationError("status", "cannot start sale for cancelled event")
	}
	
	if e.Status == EventStatusCompleted {
		return NewValidationError("status", "cannot start sale for completed event")
	}
	
	e.Status = EventStatusOnSale
	e.UpdatedAt = time.Now()
	
	return nil
}

