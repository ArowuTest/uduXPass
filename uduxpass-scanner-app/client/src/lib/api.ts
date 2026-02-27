import axios from 'axios';

// Backend API base URL - matches the uduXPass backend
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://8080-i3vhsavkuyc73e9syb280-50ae409d.us2.manus.computer';

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
  location: string;
  notes: string;
  start_time: string;
  end_time: string | null;
  status: string;
  event?: Event;
}

export interface CreateSessionRequest {
  event_id: string;
  location: string;
  notes?: string;
}

export interface ValidateTicketRequest {
  qr_code_data: string;
  session_id: string;
}

export interface ValidateTicketResponse {
  success: boolean;
  message?: string;
  error?: string;
  data?: {
    ticket: {
      id: string;
      qr_code: string;
      status: string;
      order_id: string;
      ticket_tier_id: string;
      user_id?: string;
      validated_at?: string;
      created_at: string;
      updated_at: string;
    };
    validated_at?: string;
    validation?: {
      id: string;
      ticket_id: string;
      validated_at: string;
      validated_by: string;
    };
  };
}

export interface SessionStats {
  total_scanned: number;
  valid_tickets: number;
  invalid_tickets: number;
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

  // Events
  getEvents: async (): Promise<Event[]> => {
    const response = await api.get('/v1/scanner/events');
    // Backend returns {success: true, data: {events: [...]}}
    const events = response.data?.data?.events || response.data?.events || [];
    // Map backend fields to frontend Event interface
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
  createSession: async (data: CreateSessionRequest): Promise<ScanningSession> => {
    const response = await api.post('/v1/scanner/session/start', data);
    return response.data;
  },

  getActiveSessions: async (): Promise<ScanningSession[]> => {
    const response = await api.get('/v1/scanner/session/current');
    const session = response.data;
    return session ? [session] : [];
  },

  getAllSessions: async (): Promise<ScanningSession[]> => {
    const response = await api.get('/v1/scanner/session/current');
    const session = response.data;
    return session ? [session] : [];
  },

  endSession: async (sessionId: string): Promise<void> => {
    await api.post('/v1/scanner/session/end', { session_id: sessionId });
  },

  getSessionStats: async (sessionId: string): Promise<SessionStats> => {
    const response = await api.get('/v1/scanner/stats');
    return response.data;
  },

  // Ticket Validation
  validateTicket: async (data: ValidateTicketRequest): Promise<ValidateTicketResponse> => {
    const response = await api.post('/v1/scanner/validate', {
      qr_code: data.qr_code_data,
      session_id: data.session_id
    });
    return response.data;
  },
};

export default api;
