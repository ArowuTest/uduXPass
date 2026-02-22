package events

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

// EventService handles event management use cases
type EventService struct {
	eventRepo      repositories.EventRepository
	tourRepo       repositories.TourRepository
	organizerRepo  repositories.OrganizerRepository
	ticketTierRepo repositories.TicketTierRepository
	unitOfWork     repositories.UnitOfWork
}

// NewEventService creates a new event service
func NewEventService(
	eventRepo repositories.EventRepository,
	tourRepo repositories.TourRepository,
	organizerRepo repositories.OrganizerRepository,
	ticketTierRepo repositories.TicketTierRepository,
	unitOfWork repositories.UnitOfWork,
) *EventService {
	return &EventService{
		eventRepo:      eventRepo,
		tourRepo:       tourRepo,
		organizerRepo:  organizerRepo,
		ticketTierRepo: ticketTierRepo,
		unitOfWork:     unitOfWork,
	}
}

// TicketTierRequest represents a ticket tier to be created with an event
type TicketTierRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description,omitempty"`
	Price       float64    `json:"price" validate:"required,gte=0"`
	Quota       int        `json:"quota" validate:"required,gt=0"`
	MinPurchase int        `json:"min_purchase,omitempty"`
	MaxPurchase int        `json:"max_purchase,omitempty"`
	SaleStart   *time.Time `json:"sale_start,omitempty"`
	SaleEnd     *time.Time `json:"sale_end,omitempty"`
}

// CreateEventRequest represents the request to create an event
type CreateEventRequest struct {
	OrganizerID     uuid.UUID             `json:"organizer_id" validate:"required"`
	TourID          *uuid.UUID            `json:"tour_id,omitempty"`
	Name            string                `json:"name" validate:"required"`
	Slug            string                `json:"slug" validate:"required"`
	Description     *string               `json:"description,omitempty"`
	EventDate       time.Time             `json:"event_date" validate:"required"`
	DoorsOpen       *time.Time            `json:"doors_open,omitempty"`
	VenueName       string                `json:"venue_name" validate:"required"`
	VenueAddress    string                `json:"venue_address" validate:"required"`
	VenueCity       string                `json:"venue_city" validate:"required"`
	VenueState      *string               `json:"venue_state,omitempty"`
	VenueCountry    string                `json:"venue_country" validate:"required"`
	VenueCapacity   *int                  `json:"venue_capacity,omitempty"`
	VenueLatitude   *float64              `json:"venue_latitude,omitempty"`
	VenueLongitude  *float64              `json:"venue_longitude,omitempty"`
	EventImageURL   *string               `json:"event_image_url,omitempty"`
	SaleStart       *time.Time            `json:"sale_start,omitempty"`
	SaleEnd         *time.Time            `json:"sale_end,omitempty"`
	TicketTiers     []TicketTierRequest   `json:"ticket_tiers,omitempty"`
}

// CreateEventResponse represents the response from event creation
type CreateEventResponse struct {
	Event *EventInfo `json:"event"`
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	// Verify organizer exists
	exists, err := s.organizerRepo.Exists(ctx, req.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organizer existence: %w", err)
	}
	if !exists {
		return nil, entities.NewNotFoundError("organizer", "organizer not found")
	}
	
	// Verify tour exists if provided
	if req.TourID != nil {
		exists, err := s.tourRepo.Exists(ctx, *req.TourID)
		if err != nil {
			return nil, fmt.Errorf("failed to check tour existence: %w", err)
		}
		if !exists {
			return nil, entities.NewNotFoundError("tour", "tour not found")
		}
	}
	
	// Check if event slug already exists for this organizer
	exists, err = s.eventRepo.ExistsBySlug(ctx, req.OrganizerID, req.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check event slug existence: %w", err)
	}
	if exists {
		return nil, entities.NewConflictError("event", "event with this slug already exists for this organizer", nil)
	}
	
	// Create event entity
	event := entities.NewEvent(
		req.OrganizerID,
		req.Name,
		req.Slug,
		req.EventDate,
		req.VenueName,
		req.VenueAddress,
		req.VenueCity,
		req.VenueCountry,
	)
	
	// Set optional fields
	// TourID removed from schema
	// if req.TourID != nil {
	// 	event.SetTour(*req.TourID)
	// }
	if req.Description != nil {
		event.Description = req.Description
	}
	if req.DoorsOpen != nil {
		event.DoorsOpen = req.DoorsOpen
	}
	if req.VenueState != nil {
		event.VenueState = req.VenueState
	}
	if req.VenueCapacity != nil {
		event.VenueCapacity = req.VenueCapacity
	}
	// Venue coordinates removed from schema
	// if req.VenueLatitude != nil && req.VenueLongitude != nil {
	// 	event.SetVenueLocation(*req.VenueLatitude, *req.VenueLongitude)
	// }
	if req.EventImageURL != nil {
		event.SetImage(*req.EventImageURL)
	}
	if req.SaleStart != nil && req.SaleEnd != nil {
		if err := event.SetSalePeriod(*req.SaleStart, *req.SaleEnd); err != nil {
			return nil, err
		}
	}
	
	// Validate event
	if err := event.Validate(); err != nil {
		return nil, err
	}
	
	// Begin transaction for atomic event creation
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	// Create event within transaction
	if err = tx.Events().Create(tx.Context(), event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	
	// Create ticket tiers if provided (strategic implementation)
	if len(req.TicketTiers) > 0 {
		for _, tierReq := range req.TicketTiers {
			// Create ticket tier entity
			tier := entities.NewTicketTier(event.ID, tierReq.Name, tierReq.Price)
			
			// Set optional fields
			if tierReq.Description != nil {
				tier.Description = tierReq.Description
			}
			if tierReq.Quota > 0 {
				if err := tier.SetQuota(tierReq.Quota); err != nil {
					return nil, fmt.Errorf("invalid tier quota: %w", err)
				}
			}
			if tierReq.MinPurchase > 0 && tierReq.MaxPurchase > 0 {
				if err := tier.SetPurchaseLimits(tierReq.MinPurchase, tierReq.MaxPurchase); err != nil {
					return nil, fmt.Errorf("invalid purchase limits: %w", err)
				}
			} else if tierReq.MinPurchase > 0 {
				tier.MinPurchase = tierReq.MinPurchase
			} else if tierReq.MaxPurchase > 0 {
				tier.MaxPurchase = tierReq.MaxPurchase
			}
			if tierReq.SaleStart != nil && tierReq.SaleEnd != nil {
				if err := tier.SetSalePeriod(*tierReq.SaleStart, *tierReq.SaleEnd); err != nil {
					return nil, fmt.Errorf("invalid sale period: %w", err)
				}
			}
			
			// Validate tier
			if err := tier.Validate(); err != nil {
				return nil, fmt.Errorf("invalid ticket tier '%s': %w", tierReq.Name, err)
			}
			
			// Create tier within transaction
			if err = tx.TicketTiers().Create(tx.Context(), tier); err != nil {
				return nil, fmt.Errorf("failed to create ticket tier '%s': %w", tierReq.Name, err)
			}
		}
	}
	
	// Commit transaction (atomic: event + tiers)
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return &CreateEventResponse{
		Event: mapEventToEventInfo(event),
	}, nil
}

// GetPublicEventsRequest represents the request to get public events
type GetPublicEventsRequest struct {
	City          string     `json:"city,omitempty"`
	Country       string     `json:"country,omitempty"`
	TourID        *uuid.UUID `json:"tour_id,omitempty"`
	Search        string     `json:"search,omitempty"`
	EventDateFrom *time.Time `json:"event_date_from,omitempty"`
	EventDateTo   *time.Time `json:"event_date_to,omitempty"`
	Page          int        `json:"page"`
	Limit         int        `json:"limit"`
	SortBy        string     `json:"sort_by"`
	SortOrder     string     `json:"sort_order"`
}

// GetPublicEventsResponse represents the response from getting public events
type GetPublicEventsResponse struct {
	Events     []*PublicEventInfo           `json:"events"`
	Pagination *repositories.PaginationResult `json:"pagination"`
}

// GetPublicEvents retrieves public events with filtering and pagination
func (s *EventService) GetPublicEvents(ctx context.Context, req *GetPublicEventsRequest) (*GetPublicEventsResponse, error) {
	filter := repositories.PublicEventFilter{
		BaseFilter: repositories.BaseFilter{
			Page:      req.Page,
			Limit:     req.Limit,
			SortBy:    req.SortBy,
			SortOrder: repositories.SortOrder(req.SortOrder),
		},
		City:               req.City,
		Country:            req.Country,
		TourID:             req.TourID,
		Search:             req.Search,
		EventDateFrom:      req.EventDateFrom,
		EventDateTo:        req.EventDateTo,
		IncludeTour:        true,
		IncludeTicketTiers: true,
		IncludeMinMaxPrice: true,
	}
	
	// Validate filter
	if err := filter.BaseFilter.Validate(); err != nil {
		return nil, err
	}
	
	events, pagination, err := s.eventRepo.ListPublic(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get public events: %w", err)
	}
	
	publicEvents := make([]*PublicEventInfo, len(events))
	for i, event := range events {
		publicEvents[i] = mapEventToPublicEventInfo(event)
	}
	
	return &GetPublicEventsResponse{
		Events:     publicEvents,
		Pagination: pagination,
	}, nil
}

// GetEventDetailsRequest represents the request to get event details
type GetEventDetailsRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required"`
}

// GetEventDetailsResponse represents the response from getting event details
type GetEventDetailsResponse struct {
	Event       *EventInfo                      `json:"event"`
	TicketTiers []*TicketTierAvailabilityInfo   `json:"ticket_tiers"`
	Stats       *repositories.EventStats        `json:"stats,omitempty"`
}

// GetEventDetails retrieves detailed information about an event
func (s *EventService) GetEventDetails(ctx context.Context, req *GetEventDetailsRequest) (*GetEventDetailsResponse, error) {
	// Get event
	event, err := s.eventRepo.GetByID(ctx, req.EventID)
	if err != nil {
		return nil, entities.NewNotFoundError("event", "event not found")
	}
	
	// Get ticket tier availability
	availability, err := s.ticketTierRepo.GetAvailability(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket tier availability: %w", err)
	}
	
	// Get event stats
	stats, err := s.eventRepo.GetEventStats(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event stats: %w", err)
	}
	
	ticketTiers := make([]*TicketTierAvailabilityInfo, len(availability))
	for i, tier := range availability {
		ticketTiers[i] = mapTicketTierAvailabilityToInfo(tier)
	}
	
	return &GetEventDetailsResponse{
		Event:       mapEventToEventInfo(event),
		TicketTiers: ticketTiers,
		Stats:       stats,
	}, nil
}

// PublishEventRequest represents the request to publish an event
type PublishEventRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required"`
}

// PublishEventResponse represents the response from publishing an event
type PublishEventResponse struct {
	Event   *EventInfo `json:"event"`
	Message string     `json:"message"`
}

// PublishEvent publishes an event to make it visible to the public
func (s *EventService) PublishEvent(ctx context.Context, req *PublishEventRequest) (*PublishEventResponse, error) {
	// Get event
	event, err := s.eventRepo.GetByID(ctx, req.EventID)
	if err != nil {
		return nil, entities.NewNotFoundError("event", "event not found")
	}
	
	// Publish event
	if err := event.Publish(); err != nil {
		return nil, err
	}
	
	// Update event
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}
	
	return &PublishEventResponse{
		Event:   mapEventToEventInfo(event),
		Message: "Event published successfully",
	}, nil
}

// StartSaleRequest represents the request to start ticket sales
type StartSaleRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required"`
}

// StartSaleResponse represents the response from starting ticket sales
type StartSaleResponse struct {
	Event   *EventInfo `json:"event"`
	Message string     `json:"message"`
}

// StartSale starts ticket sales for an event
func (s *EventService) StartSale(ctx context.Context, req *StartSaleRequest) (*StartSaleResponse, error) {
	// Get event
	event, err := s.eventRepo.GetByID(ctx, req.EventID)
	if err != nil {
		return nil, entities.NewNotFoundError("event", "event not found")
	}
	
	// Start sale
	if err := event.StartSale(); err != nil {
		return nil, err
	}
	
	// Update event
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}
	
	return &StartSaleResponse{
		Event:   mapEventToEventInfo(event),
		Message: "Ticket sales started successfully",
	}, nil
}

// Response types
type EventInfo struct {
	ID              uuid.UUID                `json:"id"`
	OrganizerID     *uuid.UUID               `json:"organizer_id,omitempty"`
	CategoryID      *uuid.UUID               `json:"category_id,omitempty"`
	Name            string                   `json:"name"`
	Slug            string                   `json:"slug"`
	Description     *string                  `json:"description,omitempty"`
	EventDate       time.Time                `json:"event_date"`
	DoorsOpen       *time.Time               `json:"doors_open,omitempty"`
	VenueName       string                   `json:"venue_name"`
	VenueAddress    string                   `json:"venue_address"`
	VenueCity       string                   `json:"venue_city"`
	VenueState      *string                  `json:"venue_state,omitempty"`
	VenueCountry    *string                  `json:"venue_country,omitempty"`
	VenueCapacity   *int                     `json:"venue_capacity,omitempty"`
	EventImageURL   *string                  `json:"event_image_url,omitempty"`
	Status          entities.EventStatus     `json:"status"`
	SaleStart       *time.Time               `json:"sale_start,omitempty"`
	SaleEnd         *time.Time               `json:"sale_end,omitempty"`
	Currency        *string                  `json:"currency,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	IsActive        bool                     `json:"is_active"`
}

type PublicEventInfo struct {
	ID             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	Slug           string               `json:"slug"`
	Description    *string              `json:"description,omitempty"`
	EventDate      time.Time            `json:"event_date"`
	DoorsOpen      *time.Time           `json:"doors_open,omitempty"`
	VenueName      string               `json:"venue_name"`
	VenueCity      string               `json:"venue_city"`
	VenueAddress   string               `json:"venue_address"`
	EventImageURL  *string              `json:"event_image_url,omitempty"`
	Status         entities.EventStatus `json:"status"`
	MinPrice       *float64             `json:"min_price,omitempty"`
	MaxPrice       *float64             `json:"max_price,omitempty"`
	Currency       string               `json:"currency"`
	TourInfo       *TourInfo            `json:"tour_info,omitempty"`
}

type TourInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ArtistName string    `json:"artist_name"`
}

type TicketTierAvailabilityInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Currency   string    `json:"currency"`
	Available  int       `json:"available"`
	IsOnSale   bool      `json:"is_on_sale"`
	SaleStatus string    `json:"sale_status"`
}

// Helper functions
func mapEventToEventInfo(event *entities.Event) *EventInfo {
	return &EventInfo{
		ID:             event.ID,
		OrganizerID:    event.OrganizerID,
		CategoryID:     event.CategoryID,
		Name:           event.Name,
		Slug:           event.Slug,
		Description:    event.Description,
		EventDate:      event.EventDate,
		DoorsOpen:      event.DoorsOpen,
		VenueName:      event.VenueName,
		VenueAddress:   event.VenueAddress,
		VenueCity:      event.VenueCity,
		VenueState:     event.VenueState,
		VenueCountry:   event.VenueCountry,
		VenueCapacity:  event.VenueCapacity,
		EventImageURL:  event.EventImageURL,
		Status:         event.Status,
		SaleStart:      event.SaleStart,
		SaleEnd:        event.SaleEnd,
		Currency:       event.Currency,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
		IsActive:       event.IsActive,
	}
}

func mapEventToPublicEventInfo(event *entities.Event) *PublicEventInfo {
	return &PublicEventInfo{
		ID:            event.ID,
		Name:          event.Name,
		Slug:          event.Slug,
		Description:   event.Description,
		EventDate:     event.EventDate,
		DoorsOpen:     event.DoorsOpen,
		VenueName:     event.VenueName,
		VenueCity:     event.VenueCity,
		VenueAddress:  event.VenueAddress,
		EventImageURL: event.EventImageURL,
		Status:        event.Status,
		Currency:      "NGN",
		// TODO: Calculate min/max prices from ticket tiers
	}
}

func mapTicketTierAvailabilityToInfo(tier *repositories.TicketTierAvailability) *TicketTierAvailabilityInfo {
	return &TicketTierAvailabilityInfo{
		ID:         tier.TicketTierID,
		Name:       tier.Name,
		Price:      tier.Price,
		Currency:   tier.Currency,
		Available:  tier.Available,
		IsOnSale:   tier.IsOnSale,
		SaleStatus: tier.SaleStatus,
	}
}

