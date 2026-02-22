// updated
import {
  ApiResponse,
  AuthResponse,
  LoginCredentials,
  RegisterData,
  AdminAuthResponse,
  AdminUser,
  Event,
  CreateEventData,
  EventsQueryParams,
  EventStats,
  Order,
  CreateOrderData,
  OrdersQueryParams,
  InitiatePaymentData,
  PaymentResponse,
  Ticket,
  ScanTicketData,
  ScanResult,
  TicketsQueryParams,
  UserProfile,
  UpdateProfileData,
  ChangePasswordData,
  DashboardStats,
  SalesReport,
  SalesReportParams,
  RevenueReport,
  RevenueReportParams,
  ScannerStats,
  AdminUsersQueryParams,
  CreateAdminUserData,
  UpdateAdminUserData,
  HealthStatus,
  DatabaseHealth,
  ServicesHealth,
  PaginatedResponse
} from '../types/api';

// Import data transformers
import { transformBackendEventToFrontend, transformFrontendEventToBackend } from './dataTransformers';

// Production API configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';
console.log('[API Service] Base URL:', API_BASE_URL);

// Helper function to get auth headers
const getAuthHeaders = (): Record<string, string> => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

// FIXED: adminApiRequest function - was missing entirely
const adminApiRequest = async <T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> => {
  // FIXED: Use 'adminToken' to match AuthContext storage
  const token = localStorage.getItem('adminToken');
  const url = `${API_BASE_URL}${endpoint}`;
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    // Handle non-JSON responses
    const contentType = response.headers.get('content-type');
    let data: any;
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      data = { message: await response.text() };
    }

    if (!response.ok) {
      throw new Error(data.error || data.message || `HTTP error! status: ${response.status}`);
    }

    return { success: true, data };
  } catch (error) {
    console.error(`Admin API request failed: ${endpoint}`, error);
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    return { success: false, error: errorMessage };
  }
};

// Login-specific API request function (no auth headers)
const loginRequest = async <T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> => {
  const url = `${API_BASE_URL}${endpoint}`;
  console.log('[Login Request] URL:', url, 'Endpoint:', endpoint);
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    // Handle non-JSON responses
    const contentType = response.headers.get('content-type');
    let data: any;
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      data = { message: await response.text() };
    }

    if (!response.ok) {
      throw new Error(data.error || data.message || `HTTP error! status: ${response.status}`);
    }

    return { success: true, data };
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    return { success: false, error: errorMessage };
  }
};

// Generic API request function with proper typing
const apiRequest = async <T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> => {
  const url = `${API_BASE_URL}${endpoint}`;
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    // Handle non-JSON responses
    const contentType = response.headers.get('content-type');
    let data: any;
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      data = { message: await response.text() };
    }

    if (!response.ok) {
      throw new Error(data.error || data.message || `HTTP error! status: ${response.status}`);
    }

    return { success: true, data };
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    return { success: false, error: errorMessage };
  }
};

// Authentication API
export const authAPI = {
  login: async (credentials: LoginCredentials): Promise<ApiResponse<AuthResponse>> => {
    return loginRequest<AuthResponse>('/v1/auth/email/login', {
      method: 'POST',
      body: JSON.stringify(credentials)
    });
  },

  register: async (userData: RegisterData): Promise<ApiResponse<AuthResponse>> => {
    return loginRequest<AuthResponse>('/v1/auth/email/register', {
      method: 'POST',
      body: JSON.stringify(userData)
    });
  },

  logout: async (): Promise<ApiResponse<void>> => {
    const result = await apiRequest<void>('/auth/logout', { method: 'POST' });
    // Clear local storage regardless of API response
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    return result;
  },

  refreshToken: async (refreshToken: string): Promise<ApiResponse<AuthResponse>> => {
    return apiRequest<AuthResponse>('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken })
    });
  },

  verifyEmail: async (token: string): Promise<ApiResponse<void>> => {
    return apiRequest<void>('/auth/verify-email', {
      method: 'POST',
      body: JSON.stringify({ token })
    });
  },

  resetPassword: async (email: string): Promise<ApiResponse<void>> => {
    return apiRequest<void>('/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify({ email })
    });
  },

  confirmResetPassword: async (token: string, newPassword: string): Promise<ApiResponse<void>> => {
    return apiRequest<void>('/auth/confirm-reset-password', {
      method: 'POST',
      body: JSON.stringify({ token, new_password: newPassword })
    });
  }
};

// Admin Authentication API
export const adminAuthAPI = {
  login: async (credentials: LoginCredentials): Promise<ApiResponse<AdminAuthResponse>> => {
    return loginRequest<AdminAuthResponse>('/v1/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials)
    });
  },

  me: async (): Promise<ApiResponse<AdminUser>> => {
    return adminApiRequest<AdminUser>('/admin/auth/me');
  },

  logout: async (): Promise<ApiResponse<void>> => {
    const result = await adminApiRequest<void>('/admin/auth/logout', { method: 'POST' });
    // Clear admin storage regardless of API response
    localStorage.removeItem('adminToken');
    localStorage.removeItem('adminData');
    return result;
  }
};

// Events API
export const eventsAPI = {
  // Get events with pagination
  getEvents: async (params?: EventsQueryParams): Promise<ApiResponse<PaginatedResponse<Event>>> => {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.search) queryParams.append('search', params.search);
    if (params?.city) queryParams.append('city', params.city);

    const endpoint = `/v1/events${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    const response = await apiRequest<any>(endpoint);
    
    // Transform backend response to frontend format
    if (response.success && response.data && response.data.events) {
      response.data.events = response.data.events.map(transformBackendEventToFrontend);
      // Rename 'events' to 'data' for PaginatedResponse format
      response.data = {
        data: response.data.events,
        pagination: response.data.pagination
      };
    }
    
    return response as ApiResponse<PaginatedResponse<Event>>;
  },

  getAll: async (params?: EventsQueryParams): Promise<ApiResponse<Event[]>> => {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.search) queryParams.append('search', params.search);
    if (params?.status) queryParams.append('status', params.status);

    const endpoint = `/v1/events${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    const response = await apiRequest<any>(endpoint);
    
    // Transform backend response to frontend format
    if (response.success && response.data) {
      if (response.data.events) {
        // Handle paginated response
        response.data.events = response.data.events.map(transformBackendEventToFrontend);
      } else if (Array.isArray(response.data)) {
        // Handle direct array response
        response.data = response.data.map(transformBackendEventToFrontend);
      }
    }
    
    return response as ApiResponse<Event[]>;
  },

  getById: async (id: string): Promise<ApiResponse<Event>> => {
    const response = await apiRequest<any>(`/v1/events/${id}`);
    
    // Transform backend response to frontend format
    if (response.success && response.data) {
      // Backend returns {success, data: {event data}}
      // apiRequest already wraps it, so response.data is the backend's full response
      // We need to extract response.data.data (the actual event)
      const eventData = response.data.data || response.data;
      response.data = transformBackendEventToFrontend(eventData);
    }
    
    return response as ApiResponse<Event>;
  },

  create: async (eventData: CreateEventData): Promise<ApiResponse<Event>> => {
    // Transform frontend data to backend format
    const backendEventData = transformFrontendEventToBackend(eventData as any);
    
    const response = await adminApiRequest<any>('/v1/admin/events', {
      method: 'POST',
      body: JSON.stringify(backendEventData)
    });
    
    // Transform backend response to frontend format
    if (response.success && response.data) {
      response.data = transformBackendEventToFrontend(response.data);
    }
    
    return response as ApiResponse<Event>;
  },

  update: async (id: string, eventData: Partial<CreateEventData>): Promise<ApiResponse<Event>> => {
    // Transform frontend data to backend format
    const backendEventData = transformFrontendEventToBackend(eventData as any);
    
    const response = await adminApiRequest<any>(`/v1/admin/events/${id}`, {
      method: 'PUT',
      body: JSON.stringify(backendEventData)
    });
    
    // Transform backend response to frontend format
    if (response.success && response.data) {
      response.data = transformBackendEventToFrontend(response.data);
    }
    
    return response as ApiResponse<Event>;
  },

  delete: async (id: string): Promise<ApiResponse<void>> => {
    return adminApiRequest<void>(`/v1/admin/events/${id}`, {
      method: 'DELETE'
    });
  },

  getStats: async (id: string): Promise<ApiResponse<EventStats>> => {
    return adminApiRequest<EventStats>(`/v1/admin/events/${id}/stats`);
  }
};

// Orders API
export const ordersAPI = {
  create: async (orderData: CreateOrderData): Promise<ApiResponse<Order>> => {
    return apiRequest<Order>('/v1/orders', {
      method: 'POST',
      body: JSON.stringify(orderData)
    });
  },

  getById: async (id: string): Promise<ApiResponse<Order>> => {
    return apiRequest<Order>(`/v1/orders/${id}`);
  },

  getUserOrders: async (): Promise<ApiResponse<Order[]>> => {
    return apiRequest<Order[]>('/v1/orders/user');
  },

  getAll: async (params?: OrdersQueryParams): Promise<ApiResponse<Order[]>> => {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.status) queryParams.append('status', params.status);
    if (params?.event_id) queryParams.append('event_id', params.event_id);

    const endpoint = `/v1/admin/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return adminApiRequest<Order[]>(endpoint);
  },

  updateStatus: async (id: string, status: string): Promise<ApiResponse<Order>> => {
    return adminApiRequest<Order>(`/admin/orders/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status })
    });
  },

  cancel: async (id: string): Promise<ApiResponse<Order>> => {
    return adminApiRequest<Order>(`/admin/orders/${id}/cancel`, {
      method: 'POST'
    });
  },

  refund: async (id: string): Promise<ApiResponse<Order>> => {
    return adminApiRequest<Order>(`/admin/orders/${id}/refund`, {
      method: 'POST'
    });
  }
};

// Payments API
export const paymentsAPI = {
  initiate: async (paymentData: InitiatePaymentData): Promise<ApiResponse<PaymentResponse>> => {
    return apiRequest<PaymentResponse>('/v1/payments/initiate', {
      method: 'POST',
      body: JSON.stringify(paymentData)
    });
  },

  verify: async (reference: string): Promise<ApiResponse<{ status: string; order: Order }>> => {
    return apiRequest<{ status: string; order: Order }>(`/v1/payments/verify/${reference}`);
  }
};

// Tickets API
export const ticketsAPI = {
  getUserTickets: async (params?: TicketsQueryParams): Promise<ApiResponse<Ticket[]>> => {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.status) queryParams.append('status', params.status);

    const endpoint = `/v1/tickets/user${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiRequest<Ticket[]>(endpoint);
  },

  getAll: async (params?: TicketsQueryParams): Promise<ApiResponse<Ticket[]>> => {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.status) queryParams.append('status', params.status);
    if (params?.event_id) queryParams.append('event_id', params.event_id);

    const endpoint = `/v1/admin/tickets${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return adminApiRequest<Ticket[]>(endpoint);
  },

  validateTicket: async (ticketId: string): Promise<ApiResponse<{ valid: boolean; ticket: Ticket }>> => {
    return adminApiRequest<{ valid: boolean; ticket: Ticket }>(`/v1/admin/tickets/${ticketId}/validate`, {
      method: 'POST'
    });
  },

  scanTicket: async (scanData: ScanTicketData): Promise<ApiResponse<ScanResult>> => {
    return adminApiRequest<ScanResult>('/v1/admin/tickets/scan', {
      method: 'POST',
      body: JSON.stringify(scanData)
    });
  }
};

// User API
export const userAPI = {
  getProfile: async (): Promise<ApiResponse<UserProfile>> => {
    return apiRequest<UserProfile>('/v1/user/profile');
  },

  updateProfile: async (userData: UpdateProfileData): Promise<ApiResponse<UserProfile>> => {
    return apiRequest<UserProfile>('/v1/user/profile', {
      method: 'PUT',
      body: JSON.stringify(userData)
    });
  },

  changePassword: async (passwordData: ChangePasswordData): Promise<ApiResponse<void>> => {
    return apiRequest<void>('/v1/user/change-password', {
      method: 'POST',
      body: JSON.stringify(passwordData)
    });
  }
};

// FIXED: Analytics API with proper adminApiRequest
export const analyticsAPI = {
  getDashboard: async (): Promise<ApiResponse<DashboardStats>> => {
    return adminApiRequest<DashboardStats>('/v1/admin/analytics/dashboard');
  },

  getEventAnalytics: async (eventId: string): Promise<ApiResponse<EventStats>> => {
    return adminApiRequest<EventStats>(`/v1/admin/analytics/events/${eventId}`);
  },

  // ADDED: Missing analytics functions that frontend expects
  getSalesAnalytics: async (params: any = {}): Promise<ApiResponse<any>> => {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = `/v1/admin/analytics/sales-analytics${queryString ? `?${queryString}` : ''}`;
    return adminApiRequest<any>(endpoint);
  },

  getUserAnalytics: async (params: any = {}): Promise<ApiResponse<any>> => {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = `/v1/admin/analytics/user-analytics${queryString ? `?${queryString}` : ''}`;
    return adminApiRequest<any>(endpoint);
  },

  getSalesReport: async (params: SalesReportParams): Promise<ApiResponse<SalesReport>> => {
    const queryParams = new URLSearchParams();
    if (params.event_id) queryParams.append('event_id', params.event_id);
    if (params.date_from) queryParams.append('date_from', params.date_from);
    if (params.date_to) queryParams.append('date_to', params.date_to);
    if (params.group_by) queryParams.append('group_by', params.group_by);

    const endpoint = `/v1/admin/analytics/sales${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return adminApiRequest<SalesReport>(endpoint);
  },

  getRevenueReport: async (params: RevenueReportParams): Promise<ApiResponse<RevenueReport>> => {
    const queryParams = new URLSearchParams();
    if (params.event_id) queryParams.append('event_id', params.event_id);
    if (params.date_from) queryParams.append('date_from', params.date_from);
    if (params.date_to) queryParams.append('date_to', params.date_to);
    if (params.group_by) queryParams.append('group_by', params.group_by);

    const endpoint = `/v1/admin/analytics/revenue${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return adminApiRequest<RevenueReport>(endpoint);
  },

  getScannerStats: async (): Promise<ApiResponse<ScannerStats>> => {
    return adminApiRequest<ScannerStats>('/v1/admin/analytics/scanners');
  }
};

// Admin Users API
export const adminUsersAPI = {
  getAdminUsers: async (params: AdminUsersQueryParams = {}): Promise<ApiResponse<AdminUser[]>> => {
    const queryString = new URLSearchParams(params as Record<string, string>).toString();
    const endpoint = `/admin/users${queryString ? `?${queryString}` : ''}`;
    return adminApiRequest<AdminUser[]>(endpoint);
  },

  createAdminUser: async (userData: CreateAdminUserData): Promise<ApiResponse<AdminUser>> => {
    return adminApiRequest<AdminUser>('/admin/users', {
      method: 'POST',
      body: JSON.stringify(userData)
    });
  },

  updateAdminUser: async (id: string, userData: UpdateAdminUserData): Promise<ApiResponse<AdminUser>> => {
    return adminApiRequest<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(userData)
    });
  },

  deleteAdminUser: async (id: string): Promise<ApiResponse<void>> => {
    return adminApiRequest<void>(`/admin/users/${id}`, {
      method: 'DELETE'
    });
  }
};

// Scanners API
export const scannersAPI = {
  getAll: async (): Promise<ApiResponse<any[]>> => {
    return adminApiRequest<any[]>('/v1/admin/scanners');
  },

  create: async (scannerData: any): Promise<ApiResponse<any>> => {
    return adminApiRequest<any>('/v1/admin/scanners', {
      method: 'POST',
      body: JSON.stringify(scannerData)
    });
  },

  getById: async (id: string): Promise<ApiResponse<any>> => {
    return adminApiRequest<any>(`/v1/admin/scanners/${id}`);
  },

  update: async (id: string, scannerData: any): Promise<ApiResponse<any>> => {
    return adminApiRequest<any>(`/v1/admin/scanners/${id}`, {
      method: 'PUT',
      body: JSON.stringify(scannerData)
    });
  },

  delete: async (id: string): Promise<ApiResponse<void>> => {
    return adminApiRequest<void>(`/v1/admin/scanners/${id}`, {
      method: 'DELETE'
    });
  }
};

// Health Check API
export const healthAPI = {
  check: async (): Promise<ApiResponse<HealthStatus>> => {
    return apiRequest<HealthStatus>('/health');
  },

  checkDatabase: async (): Promise<ApiResponse<DatabaseHealth>> => {
    return apiRequest<DatabaseHealth>('/health/database');
  },

  checkServices: async (): Promise<ApiResponse<ServicesHealth>> => {
    return apiRequest<ServicesHealth>('/health/services');
  }
};

// Export all APIs as default export
export default {
  authAPI,
  adminAuthAPI,
  eventsAPI,
  ordersAPI,
  paymentsAPI,
  ticketsAPI,
  userAPI,
  analyticsAPI,
  adminUsersAPI,
  scannersAPI,
  healthAPI
};








// Scanner Users API
export const scannerUsersAPI = {
  getAll: async (): Promise<ApiResponse<any[]>> => {
    return adminApiRequest<any[]>("/v1/admin/scanner-users");
  },

  create: async (userData: any): Promise<ApiResponse<any>> => {
    return adminApiRequest<any>("/v1/admin/scanner-users", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  },

  update: async (id: string, userData: any): Promise<ApiResponse<any>> => {
    return adminApiRequest<any>(`/v1/admin/scanner-users/${id}`, {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  },

  delete: async (id: string): Promise<ApiResponse<void>> => {
    return adminApiRequest<void>(`/v1/admin/scanner-users/${id}`, {
      method: "DELETE",
    });
  },
};

