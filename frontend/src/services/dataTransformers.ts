// Data transformers to handle API response mapping
import { Event, TicketTier } from '../types/api';

// Transform backend event response to frontend Event interface
export const transformBackendEventToFrontend = (backendEvent: any): Event => {
  return {
    id: backendEvent.id,
    organizer_id: backendEvent.organizer_id || 'default-organizer',
    tour_id: backendEvent.tour_id,
    name: backendEvent.name,
    slug: backendEvent.slug || backendEvent.name?.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, ''),
    description: backendEvent.description,
    event_date: backendEvent.startDate || backendEvent.eventDate || backendEvent.event_date,
    doors_open: backendEvent.doorsOpen || backendEvent.doors_open,
    venue_name: backendEvent.venue?.name || backendEvent.venue_name,
    venue_address: backendEvent.venue?.address || backendEvent.venue_address,
    venue_city: backendEvent.venue?.city || backendEvent.venue_city,
    venue_state: backendEvent.venue?.state || backendEvent.venue_state,
    venue_country: backendEvent.venue?.country || backendEvent.venue_country || 'Nigeria',
    venue_capacity: backendEvent.venue?.capacity || backendEvent.capacity || backendEvent.venue_capacity,
    venue_latitude: backendEvent.venue?.latitude || backendEvent.venue_latitude,
    venue_longitude: backendEvent.venue?.longitude || backendEvent.venue_longitude,
    event_image_url: backendEvent.eventImageUrl || backendEvent.event_image_url,
    status: backendEvent.status || 'published',
    sale_start: backendEvent.saleStart || backendEvent.sale_start,
    sale_end: backendEvent.saleEnd || backendEvent.sale_end,
    sales_end_date: backendEvent.salesEndDate || backendEvent.sales_end_date,
    is_active: backendEvent.is_active !== undefined ? backendEvent.is_active : true,
    settings: backendEvent.settings || {},
    created_at: backendEvent.created_at || backendEvent.createdAt || new Date().toISOString(),
    updated_at: backendEvent.updated_at || backendEvent.updatedAt || new Date().toISOString(),
    
    // Transform nested objects
    organizer: backendEvent.organizer,
    tour: backendEvent.tour,
    ticket_tiers: (backendEvent.ticketTiers || backendEvent.ticket_tiers)?.map(transformBackendTicketTierToFrontend) || [],
    
    // Additional fields that might be in backend response
    revenue: backendEvent.revenue,
    ticketsSold: backendEvent.ticketsSold || backendEvent.tickets_sold,
    rating: backendEvent.rating
  };
};

// Transform backend ticket tier response to frontend TicketTier interface
export const transformBackendTicketTierToFrontend = (backendTier: any): TicketTier => {
  return {
    id: backendTier.id,
    event_id: backendTier.event_id || backendTier.eventId,
    name: backendTier.name,
    description: backendTier.description,
    price: typeof backendTier.price === 'string' ? parseFloat(backendTier.price) : backendTier.price,
    currency: backendTier.currency || 'NGN',
    quota: backendTier.totalQuantity || backendTier.quota,
    max_per_order: backendTier.maxPerOrder || backendTier.max_per_order || 4,
    min_per_order: backendTier.minPerOrder || backendTier.min_per_order || 1,
    sale_start: backendTier.saleStart || backendTier.sale_start,
    sale_end: backendTier.saleEnd || backendTier.sale_end,
    position: backendTier.position || 0,
    is_active: backendTier.is_active !== undefined ? backendTier.is_active : true,
    settings: backendTier.settings || {},
    created_at: backendTier.created_at || backendTier.createdAt || new Date().toISOString(),
    updated_at: backendTier.updated_at || backendTier.updatedAt || new Date().toISOString(),
    
    // Additional fields from backend
    availableQuantity: backendTier.availableQuantity,
    totalQuantity: backendTier.totalQuantity
  };
};

// Transform frontend Event to backend format for API calls
export const transformFrontendEventToBackend = (frontendEvent: Partial<Event>): any => {
  return {
    id: frontendEvent.id,
    organizer_id: frontendEvent.organizer_id,
    tour_id: frontendEvent.tour_id,
    name: frontendEvent.name,
    slug: frontendEvent.slug,
    description: frontendEvent.description,
    startDate: frontendEvent.event_date,
    eventDate: frontendEvent.event_date,
    doorsOpen: frontendEvent.doors_open,
    venue: {
      name: frontendEvent.venue_name,
      address: frontendEvent.venue_address,
      city: frontendEvent.venue_city,
      state: frontendEvent.venue_state,
      country: frontendEvent.venue_country,
      capacity: frontendEvent.venue_capacity,
      latitude: frontendEvent.venue_latitude,
      longitude: frontendEvent.venue_longitude
    },
    eventImageUrl: frontendEvent.event_image_url,
    status: frontendEvent.status,
    saleStart: frontendEvent.sale_start,
    saleEnd: frontendEvent.sale_end,
    salesEndDate: frontendEvent.sales_end_date,
    is_active: frontendEvent.is_active,
    settings: frontendEvent.settings,
    created_at: frontendEvent.created_at,
    updated_at: frontendEvent.updated_at,
    
    // Transform ticket tiers
    ticketTiers: frontendEvent.ticket_tiers?.map(tier => ({
      id: tier.id,
      eventId: tier.event_id,
      name: tier.name,
      description: tier.description,
      price: tier.price,
      currency: tier.currency,
      totalQuantity: tier.quota,
      maxPerOrder: tier.max_per_order,
      minPerOrder: tier.min_per_order,
      saleStart: tier.sale_start,
      saleEnd: tier.sale_end,
      position: tier.position,
      is_active: tier.is_active,
      settings: tier.settings
    }))
  };
};
