import axios from 'axios';

// Backend API base URL
// In development: Vite proxy forwards /v1/* to http://localhost:3000
// In production: set VITE_API_BASE_URL to the deployed backend URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// ─── Request interceptor: attach access token ─────────────────────────────────
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('scanner_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// ─── Response interceptor: handle 401 with token refresh ──────────────────────
let isRefreshing = false;
let failedQueue: Array<{ resolve: (v: string) => void; reject: (e: unknown) => void }> = [];

function processQueue(error: unknown, token: string | null = null) {
  failedQueue.forEach((prom) => {
    if (error) prom.reject(error);
    else prom.resolve(token!);
  });
  failedQueue = [];
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return api(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = localStorage.getItem('scanner_refresh_token');
      if (!refreshToken) {
        // No refresh token — clear auth and redirect
        localStorage.removeItem('scanner_token');
        localStorage.removeItem('scanner_refresh_token');
        localStorage.removeItem('scanner_user');
        window.location.href = '/login';
        return Promise.reject(error);
      }

      try {
        const response = await axios.post(`${API_BASE_URL}/v1/scanner/auth/refresh`, {}, {
          headers: { 'X-Refresh-Token': refreshToken },
        });
        const newToken = response.data?.access_token || response.data?.data?.access_token;
        if (!newToken) throw new Error('No token in refresh response');

        localStorage.setItem('scanner_token', newToken);
        if (response.data?.refresh_token) {
          localStorage.setItem('scanner_refresh_token', response.data.refresh_token);
        }
        api.defaults.headers.common.Authorization = `Bearer ${newToken}`;
        processQueue(null, newToken);
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        localStorage.removeItem('scanner_token');
        localStorage.removeItem('scanner_refresh_token');
        localStorage.removeItem('scanner_user');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

// ─── Types ─────────────────────────────────────────────────────────────────────

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
  event_name?: string;
  scanner_id: string;
  start_time: string;
  end_time: string | null;
  scans_count: number;
  valid_scans: number;
  invalid_scans: number;
  total_revenue: number;
  is_active: boolean;
  status: string;
  location?: string;
  notes?: string;
  event?: Event;
}

export interface CreateSessionRequest {
  event_id: string;
  location?: string;
  notes?: string;
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
 * ValidateTicketResponse matches the backend TicketValidationResponse struct
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

export interface ScannerStats {
  scanner_id: string;
  total_sessions: number;
  total_scans: number;
  valid_scans: number;
  invalid_scans: number;
  total_revenue: number;
  success_rate: number;
  last_active_at: string;
  events_assigned: number;
}

// Legacy alias
export type SessionStats = ScannerStats;

// ─── API Methods ───────────────────────────────────────────────────────────────
export const scannerApi = {
  // Authentication
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/v1/scanner/auth/login', data);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('scanner_token');
    localStorage.removeItem('scanner_refresh_token');
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
      location: [e.venue_name, e.venue_city].filter(Boolean).join(', ') || e.location || '',
      status: e.status || 'published',
    }));
  },

  // Sessions
  // Backend: POST /v1/scanner/session/start with { event_id: uuid }
  createSession: async (data: CreateSessionRequest): Promise<ScanningSession> => {
    const response = await api.post('/v1/scanner/session/start', {
      event_id: data.event_id,
      ...(data.location && { location: data.location }),
      ...(data.notes && { notes: data.notes }),
    });
    return response.data?.data || response.data;
  },

  // Backend: GET /v1/scanner/session/current
  getCurrentSession: async (): Promise<ScanningSession | null> => {
    try {
      const response = await api.get('/v1/scanner/session/current');
      return response.data?.data || null;
    } catch (err: any) {
      // 404 means no active session
      if (err?.response?.status === 404) return null;
      throw err;
    }
  },

  // Legacy alias
  getActiveSessions: async (): Promise<ScanningSession[]> => {
    const session = await scannerApi.getCurrentSession();
    return session ? [session] : [];
  },

  getAllSessions: async (): Promise<ScanningSession[]> => {
    const session = await scannerApi.getCurrentSession();
    return session ? [session] : [];
  },

  // Backend: POST /v1/scanner/session/end
  endSession: async (): Promise<void> => {
    await api.post('/v1/scanner/session/end', {});
  },

  // Backend: GET /v1/scanner/stats
  getStats: async (): Promise<ScannerStats> => {
    const response = await api.get('/v1/scanner/stats');
    return response.data?.data || response.data;
  },

  // Legacy alias
  getSessionStats: async (): Promise<SessionStats> => {
    return scannerApi.getStats();
  },

  // Ticket Validation
  // Backend: POST /v1/scanner/validate with { ticket_code: string, event_id: uuid, notes?: string }
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
