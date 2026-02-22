/**
 * Enterprise-grade Event Data Transformer
 * 
 * This utility provides robust data transformation between API responses
 * and frontend interfaces, ensuring compatibility across different API versions
 * and providing a stable interface for the frontend components.
 * 
 * @author uduXPass Platform Team
 * @version 1.0.0
 */

import { Event } from '../types/api';

/**
 * Raw event data structure from API (camelCase)
 */
interface RawEventData {
  id: string;
  name: string;
  slug?: string;
  description?: string;
  eventDate: string;
  doorsOpen?: string;
  eventImageUrl?: string;
  status: string;
  saleStart?: string;
  saleEnd?: string;
  settings?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
  venue?: {
    name: string;
    address: string;
    city: string;
    state?: string;
    country: string;
    capacity?: number;
    latitude?: number;
    longitude?: number;
  };
  organizer_id?: string;
  tour_id?: string;
  ticket_tiers?: any[];
  orders?: any[];
  tickets?: any[];
}

/**
 * Transforms raw API event data to frontend Event interface
 * 
 * @param rawEvent - Raw event data from API
 * @returns Transformed event data matching frontend interface
 */
export function transformEventData(rawEvent: RawEventData): Event {
  return {
    id: rawEvent.id,
    organizer_id: rawEvent.organizer_id || '',
    tour_id: rawEvent.tour_id,
    name: rawEvent.name,
    slug: rawEvent.slug || generateSlug(rawEvent.name),
    description: rawEvent.description,
    event_date: rawEvent.eventDate,
    doors_open: rawEvent.doorsOpen,
    venue_name: rawEvent.venue?.name || 'TBD',
    venue_address: rawEvent.venue?.address || '',
    venue_city: rawEvent.venue?.city || '',
    venue_state: rawEvent.venue?.state,
    venue_country: rawEvent.venue?.country || 'Nigeria',
    venue_capacity: rawEvent.venue?.capacity,
    venue_latitude: rawEvent.venue?.latitude,
    venue_longitude: rawEvent.venue?.longitude,
    event_image_url: rawEvent.eventImageUrl,
    status: rawEvent.status as any,
    sale_start: rawEvent.saleStart,
    sale_end: rawEvent.saleEnd,
    sales_end_date: rawEvent.saleEnd, // Alias for compatibility
    is_active: rawEvent.status !== 'cancelled',
    settings: rawEvent.settings || {},
    created_at: rawEvent.createdAt,
    updated_at: rawEvent.updatedAt,
    
    // Relations (if provided)
    ticket_tiers: rawEvent.ticket_tiers,
    orders: rawEvent.orders,
    tickets: rawEvent.tickets
  };
}

/**
 * Transforms an array of raw event data
 * 
 * @param rawEvents - Array of raw event data from API
 * @returns Array of transformed events
 */
export function transformEventsArray(rawEvents: RawEventData[]): Event[] {
  return rawEvents.map(transformEventData);
}

/**
 * Generates a URL-friendly slug from event name
 * 
 * @param name - Event name
 * @returns URL-friendly slug
 */
function generateSlug(name: string): string {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');
}

/**
 * Validates if the raw event data has required fields
 * 
 * @param rawEvent - Raw event data to validate
 * @returns True if valid, false otherwise
 */
export function validateEventData(rawEvent: any): rawEvent is RawEventData {
  return (
    rawEvent &&
    typeof rawEvent.id === 'string' &&
    typeof rawEvent.name === 'string' &&
    typeof rawEvent.eventDate === 'string' &&
    typeof rawEvent.status === 'string'
  );
}

/**
 * Safely transforms event data with validation
 * 
 * @param rawEvent - Raw event data from API
 * @returns Transformed event or null if invalid
 */
export function safeTransformEventData(rawEvent: any): Event | null {
  if (!validateEventData(rawEvent)) {
    console.warn('Invalid event data received:', rawEvent);
    return null;
  }
  
  try {
    return transformEventData(rawEvent);
  } catch (error) {
    console.error('Error transforming event data:', error, rawEvent);
    return null;
  }
}

/**
 * Transforms API response containing events
 * 
 * @param apiResponse - Full API response
 * @returns Transformed events array
 */
export function transformEventsResponse(apiResponse: any): Event[] {
  console.log('Transformer received:', apiResponse);
  console.log('apiResponse.data:', apiResponse?.data);
  console.log('apiResponse.data keys:', apiResponse?.data ? Object.keys(apiResponse.data) : 'no data');
  
  // Check if we have the expected structure
  if (!apiResponse?.success) {
    console.warn('API response not successful:', apiResponse);
    return [];
  }
  
  if (!apiResponse?.data) {
    console.warn('No data in API response:', apiResponse);
    return [];
  }
  
  // Handle double-nested structure: apiResponse.data.data.events
  let eventsData = apiResponse.data;
  
  // Check if there's a double nesting (data.data.events)
  if (apiResponse.data.data && apiResponse.data.data.events) {
    console.log('Detected double-nested structure, using data.data.events');
    eventsData = apiResponse.data.data;
  }
  
  console.log('Using eventsData:', eventsData);
  console.log('eventsData keys:', Object.keys(eventsData));
  console.log('eventsData.events exists:', !!eventsData.events);
  console.log('eventsData.events type:', typeof eventsData.events);
  console.log('eventsData.events is array:', Array.isArray(eventsData.events));
  
  if (!eventsData.events || !Array.isArray(eventsData.events)) {
    console.warn('No events array found in:', eventsData);
    return [];
  }
  
  console.log(`Processing ${eventsData.events.length} raw events`);
  
  const transformedEvents = eventsData.events
    .map(safeTransformEventData)
    .filter((event): event is Event => event !== null);
  
  console.log(`Successfully transformed ${transformedEvents.length} events`);
  
  return transformedEvents;
}

