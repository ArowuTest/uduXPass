// Core Entity Types
export interface User {
  id: string;
  email: string;
  phone?: string;
  firstName?: string;
  lastName?: string;
  isEmailVerified: boolean;
  isPhoneVerified: boolean;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
}

export interface AdminUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: AdminRole;
  permissions: AdminPermission[];
  isActive: boolean;
  lastLoginAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Event {
  id: string;
  name: string;
  description?: string;
  venue: string;
  venueLocation?: string;
  startDate: string;
  endDate: string;
  salesStartDate: string;
  salesEndDate: string;
  timezone: string;
  currency: string;
  status: EventStatus;
  isPublic: boolean;
  maxCapacity?: number;
  imageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TicketTier {
  id: string;
  eventId: string;
  name: string;
  description?: string;
  price: number;
  currency: string;
  maxQuantity: number;
  soldQuantity: number;
  availableQuantity: number;
  minPurchase: number;
  maxPurchase: number;
  salesStartDate: string;
  salesEndDate: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Order {
  id: string;
  code: string;
  eventId: string;
  userId?: string;
  email: string;
  phone?: string;
  firstName?: string;
  lastName?: string;
  customerFirstName?: string;
  customerLastName?: string;
  customerEmail?: string;
  customerPhone?: string;
  status: OrderStatus;
  totalAmount: number;
  currency: string;
  paymentMethod?: PaymentMethod;
  paymentId?: string;
  notes?: string;
  confirmedAt?: string;
  cancelledAt?: string;
  expiresAt: string;
  orderLines: OrderLine[];
  createdAt: string;
  updatedAt: string;
}

export interface OrderLine {
  id: string;
  orderId: string;
  ticketTierId: string;
  ticketTier: TicketTier;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  createdAt: string;
  updatedAt: string;
}

export interface Ticket {
  id: string;
  orderLineId: string;
  orderLine: OrderLine;
  code: string;
  qrCode: string;
  status: TicketStatus;
  scannedAt?: string;
  scannedBy?: string;
  createdAt: string;
  updatedAt: string;
}

// Enum Types
export enum UserRole {
  USER = 'user',
  ADMIN = 'admin'
}

export enum AdminRole {
  SUPER_ADMIN = 'super_admin',
  ADMIN = 'admin',
  MODERATOR = 'moderator'
}

export enum AdminPermission {
  USERS_VIEW = 'users_view',
  USERS_CREATE = 'users_create',
  USERS_UPDATE = 'users_update',
  USERS_DELETE = 'users_delete',
  EVENTS_VIEW = 'events_view',
  EVENTS_CREATE = 'events_create',
  EVENTS_UPDATE = 'events_update',
  EVENTS_DELETE = 'events_delete',
  ORDERS_VIEW = 'orders_view',
  ORDERS_UPDATE = 'orders_update',
  ORDERS_DELETE = 'orders_delete',
  TICKETS_VIEW = 'tickets_view',
  TICKETS_UPDATE = 'tickets_update',
  SCANNERS_VIEW = 'scanners_view',
  SCANNERS_CREATE = 'scanners_create',
  SCANNERS_UPDATE = 'scanners_update',
  SCANNERS_DELETE = 'scanners_delete',
  ANALYTICS_VIEW = 'analytics_view',
  REPORTS_VIEW = 'reports_view',
  SETTINGS_UPDATE = 'settings_update',
  ADMIN_CREATE = 'admin_create',
  ADMIN_UPDATE = 'admin_update',
  ADMIN_DELETE = 'admin_delete'
}

export enum EventStatus {
  DRAFT = 'draft',
  PUBLISHED = 'published',
  CANCELLED = 'cancelled',
  COMPLETED = 'completed'
}

export enum OrderStatus {
  PENDING = 'pending',
  CONFIRMED = 'confirmed',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired'
}

export enum TicketStatus {
  VALID = 'valid',
  USED = 'used',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired'
}

export enum PaymentMethod {
  CARD = 'card',
  MOMO = 'momo',
  BANK_TRANSFER = 'bank_transfer'
}

// API Response Types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T = any> {
  data: T[];
  pagination: {
    totalCount: number;
    totalPages: number;
    currentPage: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
}

// Authentication Types
export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export interface AdminAuthResponse {
  accessToken: string;
  refreshToken: string;
  admin: AdminUser;
}

// Form Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
}

export interface CreateEventRequest {
  name: string;
  description?: string;
  venue: string;
  venueLocation?: string;
  startDate: string;
  endDate: string;
  salesStartDate: string;
  salesEndDate: string;
  timezone: string;
  currency: string;
  isPublic: boolean;
  maxCapacity?: number;
  imageUrl?: string;
}

export interface CreateOrderRequest {
  eventId: string;
  email: string;
  phone?: string;
  firstName?: string;
  lastName?: string;
  orderLines: {
    ticketTierId: string;
    quantity: number;
  }[];
}

// Component Props Types
export interface BaseComponentProps {
  className?: string;
  children?: React.ReactNode;
}

export interface ButtonProps extends BaseComponentProps {
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
  size?: 'default' | 'sm' | 'lg' | 'icon';
  disabled?: boolean;
  onClick?: () => void;
  type?: 'button' | 'submit' | 'reset';
}

export interface InputProps extends BaseComponentProps {
  type?: string;
  placeholder?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
  required?: boolean;
}

// Context Types
export interface AuthContextType {
  user: User | null;
  admin: AdminUser | null;
  isAuthenticated: boolean;
  isAdminAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  adminLogin: (credentials: LoginRequest) => Promise<void>;
  adminLogout: () => void;
}

export interface ThemeContextType {
  theme: 'light' | 'dark' | 'system';
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
}

// Cart Types
export interface CartItem {
  id: string
  eventId: string
  ticketTier: TicketTier
  quantity: number
  addedAt: string
}

export interface CartContextType {
  items: CartItem[]
  isOpen: boolean
  addItem: (eventId: string, ticketTier: TicketTier, quantity?: number) => void
  removeItem: (itemId: string) => void
  updateQuantity: (itemId: string, quantity: number) => void
  clearCart: () => void
  getTotalItems: () => number
  getTotalPrice: () => number
  getItemsByEvent: (eventId: string) => CartItem[]
  hasItems: () => boolean
  openCart: () => void
  closeCart: () => void
  toggleCart: () => void
}

// Utility Types
export type Optional<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;
export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>;

