// Core API response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// Authentication types
export type AuthProvider = 'email' | 'momo';

export interface User {
  id: string;
  email?: string;
  phone?: string;
  first_name?: string;
  last_name?: string;
  auth_provider: AuthProvider;
  momo_id?: string;
  email_verified: boolean;
  phone_verified: boolean;
  is_active: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Admin types
export interface AdminUser {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

export interface AdminAuthResponse {
  admin: AdminUser;
  access_token: string;
  expires_in: number;
}

// Event types
export type EventStatus = 'draft' | 'published' | 'on_sale' | 'sold_out' | 'cancelled' | 'completed';

export interface Organizer {
  id: string;
  name: string;
  slug: string;
  email: string;
  phone?: string;
  website_url?: string;
  logo_url?: string;
  description?: string;
  address?: string;
  city?: string;
  state?: string;
  country: string;
  is_active: boolean;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface Tour {
  id: string;
  organizer_id: string;
  name: string;
  slug: string;
  artist_name: string;
  description?: string;
  tour_image_url?: string;
  start_date?: string;
  end_date?: string;
  is_active: boolean;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface Event {
  id: string;
  organizer_id: string;
  tour_id?: string;
  name: string;
  slug: string;
  description?: string;
  event_date: string;
  doors_open?: string;
  venue_name: string;
  venue_address: string;
  venue_city: string;
  venue_state?: string;
  venue_country: string;
  venue_capacity?: number;
  venue_latitude?: number;
  venue_longitude?: number;
  event_image_url?: string;
  status: EventStatus;
  sale_start?: string;
  sale_end?: string;
  sales_end_date?: string;
  is_active: boolean;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
  
  // Relations
  organizer?: Organizer;
  tour?: Tour;
  ticket_tiers?: TicketTier[];
  orders?: Order[];
  tickets?: Ticket[];
}

export interface CreateEventData {
  name: string;
  slug: string;
  description?: string;
  event_date: string;
  doors_open?: string;
  venue_name: string;
  venue_address: string;
  venue_city: string;
  venue_state?: string;
  venue_country: string;
  venue_capacity?: number;
  venue_latitude?: number;
  venue_longitude?: number;
  event_image_url?: string;
  sale_start?: string;
  sale_end?: string;
  tour_id?: string;
}

export interface EventsQueryParams {
  page?: number;
  limit?: number;
  status?: EventStatus;
  organizer_id?: string;
  tour_id?: string;
  city?: string;
  country?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
}

// Ticket Tier types
export interface TicketTier {
  id: string;
  event_id: string;
  name: string;
  description?: string;
  price: number;
  currency: string;
  quota?: number;
  max_per_order: number;
  min_per_order: number;
  sale_start?: string;
  sale_end?: string;
  position: number;
  is_active: boolean;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
}

// Order types
export type OrderStatus = 'pending' | 'paid' | 'confirmed' | 'expired' | 'cancelled' | 'refunded';
export type PaymentMethod = 'momo' | 'paystack' | 'card';

export interface OrderLine {
  id: string;
  order_id: string;
  ticket_tier_id: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  created_at: string;
  
  // Relations
  ticket_tier?: TicketTier;
}

export interface Order {
  id: string;
  code: string;
  event_id: string;
  user_id?: string;
  email: string;
  phone?: string;
  first_name?: string;
  last_name?: string;
  customer_first_name: string;
  customer_last_name: string;
  customer_email: string;
  customer_phone: string;
  status: OrderStatus;
  total_amount: number;
  currency: string;
  payment_method?: PaymentMethod;
  payment_id?: string;
  notes?: string;
  confirmed_at?: string;
  cancelled_at?: string;
  expires_at: string;
  locale: string;
  comment?: string;
  meta_info: Record<string, any>;
  created_at: string;
  updated_at: string;
  
  // Relations
  event?: Event;
  user?: User;
  order_lines?: OrderLine[];
  payments?: Payment[];
  tickets?: Ticket[];
}

export interface CreateOrderData {
  event_id: string;
  customer_info: {
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
  };
  order_lines: {
    ticket_tier_id: string;
    quantity: number;
  }[];
  payment_method?: PaymentMethod;
  notes?: string;
}

export interface OrdersQueryParams {
  page?: number;
  limit?: number;
  status?: OrderStatus;
  event_id?: string;
  user_id?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  date_from?: string;
  date_to?: string;
}

// Payment types
export type PaymentStatus = 'pending' | 'completed' | 'failed' | 'cancelled' | 'refunded';

export interface Payment {
  id: string;
  order_id: string;
  provider: PaymentMethod;
  provider_transaction_id?: string;
  amount: number;
  currency: string;
  status: PaymentStatus;
  provider_response: Record<string, any>;
  webhook_received_at?: string;
  created_at: string;
  updated_at: string;
}

export interface InitiatePaymentData {
  order_id: string;
  payment_method: PaymentMethod;
  return_url?: string;
  callback_url?: string;
}

export interface PaymentResponse {
  payment_id: string;
  payment_url?: string;
  reference: string;
  status: PaymentStatus;
  expires_at?: string;
}

// Ticket types
export type TicketStatus = 'active' | 'redeemed' | 'voided';

export interface Ticket {
  id: string;
  order_line_id: string;
  serial_number: string;
  qr_code_data: string;
  status: TicketStatus;
  redeemed_at?: string;
  redeemed_by?: string;
  created_at: string;
  updated_at: string;
  
  // Relations
  order_line?: OrderLine;
}

export interface ScanTicketData {
  qr_code_data: string;
  scanner_id?: string;
  scanner_name?: string;
}

export interface ScanResult {
  valid: boolean;
  ticket?: Ticket;
  message: string;
  scan_timestamp: string;
  scanner_id?: string;
}

export interface TicketsQueryParams {
  page?: number;
  limit?: number;
  status?: TicketStatus;
  event_id?: string;
  order_id?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
}

// Analytics types
export interface DashboardStats {
  total_events: number;
  active_events: number;
  total_orders: number;
  total_revenue: number;
  total_tickets_sold: number;
  total_tickets_scanned: number;
  revenue_this_month: number;
  orders_this_month: number;
  top_events: Array<{
    event_id: string;
    event_name: string;
    revenue: number;
    tickets_sold: number;
  }>;
  recent_orders: Order[];
}

export interface EventStats {
  event_id: string;
  event_name: string;
  total_capacity: number;
  tickets_sold: number;
  tickets_scanned: number;
  revenue: number;
  orders_count: number;
  conversion_rate: number;
  tier_breakdown: Array<{
    tier_name: string;
    sold: number;
    revenue: number;
  }>;
  daily_sales: Array<{
    date: string;
    tickets: number;
    revenue: number;
  }>;
}

export interface SalesReportParams {
  event_id?: string;
  date_from?: string;
  date_to?: string;
  group_by?: 'day' | 'week' | 'month';
}

export interface SalesReport {
  period: string;
  total_revenue: number;
  total_tickets: number;
  total_orders: number;
  average_order_value: number;
  breakdown: Array<{
    period: string;
    revenue: number;
    tickets: number;
    orders: number;
  }>;
}

export interface RevenueReportParams {
  event_id?: string;
  date_from?: string;
  date_to?: string;
  currency?: string;
}

export interface RevenueReport {
  total_revenue: number;
  currency: string;
  payment_method_breakdown: Array<{
    method: PaymentMethod;
    revenue: number;
    percentage: number;
  }>;
  tier_breakdown: Array<{
    tier_name: string;
    revenue: number;
    percentage: number;
  }>;
  daily_revenue: Array<{
    date: string;
    revenue: number;
  }>;
}

export interface ScannerStats {
  total_scans: number;
  valid_scans: number;
  invalid_scans: number;
  success_rate: number;
  scans_today: number;
  active_scanners: number;
  recent_scans: Array<{
    ticket_id: string;
    event_name: string;
    scan_time: string;
    status: 'valid' | 'invalid';
  }>;
}

// User management types
export interface UserProfile {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  email_verified: boolean;
  phone_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpdateProfileData {
  first_name?: string;
  last_name?: string;
  phone?: string;
}

export interface ChangePasswordData {
  current_password: string;
  new_password: string;
  confirm_password: string;
}

// Admin user management types
export interface AdminUsersQueryParams {
  page?: number;
  limit?: number;
  role?: string;
  is_active?: boolean;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
}

export interface CreateAdminUserData {
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  password: string;
}

export interface UpdateAdminUserData {
  first_name?: string;
  last_name?: string;
  role?: string;
  is_active?: boolean;
}

// Health check types
export interface HealthStatus {
  status: 'healthy' | 'unhealthy';
  timestamp: string;
  version?: string;
  uptime?: number;
}

export interface DatabaseHealth {
  status: 'healthy' | 'unhealthy';
  connection_pool: {
    active: number;
    idle: number;
    max: number;
  };
  response_time: number;
}

export interface ServicesHealth {
  database: DatabaseHealth;
  redis?: {
    status: 'healthy' | 'unhealthy';
    response_time: number;
  };
  payment_providers?: {
    paystack: 'healthy' | 'unhealthy';
    momo: 'healthy' | 'unhealthy';
  };
}

// Error types
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
  field?: string;
}

export interface ValidationError {
  field: string;
  message: string;
  code: string;
}

// Pagination types
export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

// Generic query parameters
export interface BaseQueryParams {
  page?: number;
  limit?: number;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
}

