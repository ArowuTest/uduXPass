import axios from 'axios';

// Backend API base URL - matches the uduXPass backend
// In production, set VITE_API_BASE_URL to the deployed backend URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('scanner_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Types
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  access_token: string;
  refresh_token: string;
  scanner: {
    id: string;
    username: string;
    name: string;
    email: string;
    role: string;
    status: string;
  };
  expires_in: number;
  message: string;
}

export interface Event {
  id: string;
  name: string;
  description: string;
  start_time: string;
  end_time: string;
  location: string;
  status: string;
}

export interface ScanningSession {
  id: string;
  event_id: string;
  scanner_id: string;
  start_time: string;
  end_time: string | null;
  scans_count: number;
  valid_scans: number;
  invalid_scans: number;
  total_revenue: number;
  is_active: boolean;
  notes?: string;
  event?: Event;
}

export interface CreateSessionRequest {
  event_id: string;
}

/**
 * ValidateTicketRequest matches the backend TicketValidationRequest struct:
 * - ticket_code: the JWT QR code data string
 * - event_id: UUID of the event being scanned
 * - notes: optional notes about this scan
 */
export interface ValidateTicketRequest {
  ticket_code: string;
  event_id: string;
  notes?: string;
}

/**
 * ValidateTicketResponse matches the backend TicketValidationResponse struct:
 * - success: whether the API call succeeded
 * - valid: whether the ticket is valid for entry
 * - message: human-readable result message
 * - ticket_id: UUID of the validated ticket
 * - serial_number: e.g. "UDUX-3B9E-FAC204"
 * - ticket_type: tier name
 * - holder_name: attendee name
 * - validation_time: ISO timestamp
 * - already_validated: true if ticket was already scanned
 */
export interface ValidateTicketResponse {
  success: boolean;
  valid: boolean;
  message: string;
  ticket_id?: string;
  serial_number?: string;
  ticket_type?: string;
  holder_name?: string;
  validation_time: string;
  already_validated: boolean;
}

export interface SessionStats {
  scans_count: number;
  valid_scans: number;
  invalid_scans: number;
  total_revenue: number;
}

// API Methods
export const scannerApi = {
  // Authentication
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/v1/scanner/auth/login', data);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('scanner_token');
    localStorage.removeItem('scanner_user');
  },

  // Events - backend returns {success: true, data: {events: [...]}}
  getEvents: async (): Promise<Event[]> => {
    const response = await api.get('/v1/scanner/events');
    const events = response.data?.data?.events || response.data?.events || [];
    return events.map((e: any) => ({
      id: e.event_id || e.id,
      name: e.event_name || e.name,
      description: e.description || '',
      start_time: e.event_date || e.start_time,
      end_time: e.end_time || '',
      location: e.venue_name || e.location || '',
      status: e.status || 'published'
    }));
  },

  // Sessions
  // Backend: POST /v1/scanner/session/start with { event_id: uuid }
  // Response: { success: true, data: ScannerSession, message: "Session started successfully" }
  createSession: async (data: CreateSessionRequest): Promise<ScanningSession> => {
    const response = await api.post('/v1/scanner/session/start', { event_id: data.event_id });
    // Backend wraps in { success, data, message }
    return response.data?.data || response.data;
  },

  // Backend: GET /v1/scanner/session/current
  // Response: { success: true, data: ScannerSession }
  getActiveSessions: async (): Promise<ScanningSession[]> => {
    const response = await api.get('/v1/scanner/session/current');
    const session = response.data?.data || null;
    return session ? [session] : [];
  },

  getAllSessions: async (): Promise<ScanningSession[]> => {
    const response = await api.get('/v1/scanner/session/current');
    const session = response.data?.data || null;
    return session ? [session] : [];
  },

  // Backend: POST /v1/scanner/session/end
  endSession: async (): Promise<void> => {
    await api.post('/v1/scanner/session/end', {});
  },

  getSessionStats: async (): Promise<SessionStats> => {
    const response = await api.get('/v1/scanner/stats');
    return response.data?.data || response.data;
  },

  // Ticket Validation
  // Backend: POST /v1/scanner/validate with { ticket_code: string, event_id: uuid, notes?: string }
  // Response: TicketValidationResponse (flat, no data wrapper)
  validateTicket: async (data: ValidateTicketRequest): Promise<ValidateTicketResponse> => {
    const response = await api.post('/v1/scanner/validate', {
      ticket_code: data.ticket_code,
      event_id: data.event_id,
      notes: data.notes,
    });
    return response.data;
  },
};

export default api;
