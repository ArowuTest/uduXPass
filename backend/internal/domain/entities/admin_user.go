package entities

import (
	"time"

	"github.com/google/uuid"
)

// AdminRole represents different admin access levels
type AdminRole string

const (
	AdminRoleSuperAdmin       AdminRole = "super_admin"       // Full system access
	AdminRoleAdmin            AdminRole = "admin"             // All functions except system settings
	AdminRoleEventManager     AdminRole = "event_manager"     // Event and organizer management
	AdminRoleSupport          AdminRole = "support"           // Customer support and user management
	AdminRoleAnalyst          AdminRole = "analyst"           // Read-only access to reports and analytics
	AdminRoleScannerOperator  AdminRole = "scanner_operator"  // Scanner operations only
)

// AdminStatus represents admin account status
type AdminStatus string

const (
	AdminStatusActive   AdminStatus = "active"
	AdminStatusInactive AdminStatus = "inactive"
	AdminStatusLocked   AdminStatus = "locked"
	AdminStatusSuspended AdminStatus = "suspended"
)

// AdminPermission represents specific permissions
type AdminPermission string

const (
	// System Management
	PermissionSystemSettings     AdminPermission = "system_settings"
	PermissionUserManagement     AdminPermission = "user_management"
	PermissionAdminManagement    AdminPermission = "admin_management"
	
	// Event Management
	PermissionEventCreate        AdminPermission = "event_create"
	PermissionEventEdit          AdminPermission = "event_edit"
	PermissionEventDelete        AdminPermission = "event_delete"
	PermissionEventPublish       AdminPermission = "event_publish"
	AdminPermissionEventsCreate  AdminPermission = "events_create"
	AdminPermissionEventsUpdate  AdminPermission = "events_update"
	AdminPermissionEventsDelete  AdminPermission = "events_delete"
	AdminPermissionEventsView    AdminPermission = "events_view"
	
	// Organizer Management
	PermissionOrganizerCreate    AdminPermission = "organizer_create"
	PermissionOrganizerEdit      AdminPermission = "organizer_edit"
	PermissionOrganizerDelete    AdminPermission = "organizer_delete"
	PermissionOrganizerApprove   AdminPermission = "organizer_approve"
	
	// Order Management
	PermissionOrderView          AdminPermission = "order_view"
	PermissionOrderEdit          AdminPermission = "order_edit"
	PermissionOrderRefund        AdminPermission = "order_refund"
	PermissionOrderCancel        AdminPermission = "order_cancel"
	
	// Payment Management
	PermissionPaymentView        AdminPermission = "payment_view"
	PermissionPaymentProcess     AdminPermission = "payment_process"
	PermissionPaymentRefund      AdminPermission = "payment_refund"
	
	// Customer Support
	PermissionSupportTickets     AdminPermission = "support_tickets"
	PermissionCustomerView       AdminPermission = "customer_view"
	PermissionCustomerEdit       AdminPermission = "customer_edit"
	
	// Analytics and Reports
	PermissionAnalyticsView      AdminPermission = "analytics_view"
	PermissionReportsGenerate    AdminPermission = "reports_generate"
	PermissionReportsExport      AdminPermission = "reports_export"
	
	// Scanner Management
	PermissionScannerManage      AdminPermission = "scanner_manage"
	PermissionScannerView        AdminPermission = "scanner_view"
	
	// Additional specific permissions for middleware
	AdminPermissionUsersView     AdminPermission = "users_view"
	AdminPermissionUsersUpdate   AdminPermission = "users_update"
	AdminPermissionOrdersView    AdminPermission = "orders_view"
	AdminPermissionOrdersUpdate  AdminPermission = "orders_update"
	AdminPermissionAnalyticsView AdminPermission = "analytics_view"
	AdminPermissionScannersView  AdminPermission = "scanners_view"
	AdminPermissionScannersUpdate AdminPermission = "scanners_update"
	AdminPermissionTicketsValidate AdminPermission = "tickets_validate"
)

// AdminUser represents an administrative user with role-based permissions
type AdminUser struct {
	ID                  uuid.UUID              `json:"id" db:"id"`
	Email               string                 `json:"email" db:"email"`
	PasswordHash        string                 `json:"-" db:"password_hash"`
	FirstName           string                 `json:"first_name" db:"first_name"`
	LastName            string                 `json:"last_name" db:"last_name"`
	Role                AdminRole              `json:"role" db:"role"`
	Permissions         AdminPermissionArray   `json:"permissions" db:"permissions"`
	IsActive            bool                   `json:"is_active" db:"is_active"`
	LastLogin           *time.Time             `json:"last_login,omitempty" db:"last_login"`
	LoginAttempts       int                    `json:"login_attempts" db:"login_attempts"`
	LockedUntil         *time.Time             `json:"locked_until,omitempty" db:"locked_until"`
	MustChangePassword  bool                   `json:"must_change_password" db:"must_change_password"`
	TwoFactorEnabled    bool                   `json:"two_factor_enabled" db:"two_factor_enabled"`
	TwoFactorSecret     *string                `json:"-" db:"two_factor_secret"`
	CreatedBy           *uuid.UUID             `json:"created_by,omitempty" db:"created_by"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// Error definitions
var (
	ErrAdminUserNotFound      = NewNotFoundError("admin_user", "admin user not found")
	ErrAdminUserAlreadyExists = NewConflictError("admin_user", "admin user already exists", nil)
)

// NewAdminUser creates a new admin user with specified role
func NewAdminUser(email, firstName, lastName string, role AdminRole, createdBy *uuid.UUID) *AdminUser {
	now := time.Now()
	admin := &AdminUser{
		ID:                  uuid.New(),
		Email:               email,
		FirstName:           firstName,
		LastName:            lastName,
		Role:                role,
		Permissions:         GetRolePermissions(role),
		IsActive:            true,
		LoginAttempts:       0,
		MustChangePassword:  true,
		TwoFactorEnabled:    false,
		CreatedBy:           createdBy,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	return admin
}

// GetRolePermissions returns the default permissions for a role
func GetRolePermissions(role AdminRole) AdminPermissionArray {
	switch role {
	case AdminRoleSuperAdmin:
		return []AdminPermission{
			PermissionSystemSettings,
			PermissionUserManagement,
			PermissionAdminManagement,
			PermissionEventCreate,
			PermissionEventEdit,
			PermissionEventDelete,
			PermissionEventPublish,
			PermissionOrganizerCreate,
			PermissionOrganizerEdit,
			PermissionOrganizerDelete,
			PermissionOrganizerApprove,
			PermissionOrderView,
			PermissionOrderEdit,
			PermissionOrderRefund,
			PermissionOrderCancel,
			PermissionPaymentView,
			PermissionPaymentProcess,
			PermissionPaymentRefund,
			PermissionSupportTickets,
			PermissionCustomerView,
			PermissionCustomerEdit,
			PermissionAnalyticsView,
			PermissionReportsGenerate,
			PermissionReportsExport,
			PermissionScannerManage,
			PermissionScannerView,
		}
	case AdminRoleAdmin:
		return []AdminPermission{
			PermissionUserManagement,
			PermissionEventCreate,
			PermissionEventEdit,
			PermissionEventDelete,
			PermissionEventPublish,
			PermissionOrganizerCreate,
			PermissionOrganizerEdit,
			PermissionOrganizerApprove,
			PermissionOrderView,
			PermissionOrderEdit,
			PermissionOrderRefund,
			PermissionOrderCancel,
			PermissionPaymentView,
			PermissionPaymentProcess,
			PermissionPaymentRefund,
			PermissionSupportTickets,
			PermissionCustomerView,
			PermissionCustomerEdit,
			PermissionAnalyticsView,
			PermissionReportsGenerate,
			PermissionReportsExport,
			PermissionScannerManage,
			PermissionScannerView,
		}
	case AdminRoleEventManager:
		return []AdminPermission{
			PermissionEventCreate,
			PermissionEventEdit,
			PermissionEventDelete,
			PermissionEventPublish,
			PermissionOrganizerCreate,
			PermissionOrganizerEdit,
			PermissionOrganizerApprove,
			PermissionOrderView,
			PermissionOrderEdit,
			PermissionAnalyticsView,
			PermissionReportsGenerate,
			PermissionScannerView,
		}
	case AdminRoleSupport:
		return []AdminPermission{
			PermissionSupportTickets,
			PermissionCustomerView,
			PermissionCustomerEdit,
			PermissionOrderView,
			PermissionOrderEdit,
			PermissionOrderRefund,
			PermissionOrderCancel,
			PermissionPaymentView,
		}
	case AdminRoleAnalyst:
		return []AdminPermission{
			PermissionAnalyticsView,
			PermissionReportsGenerate,
			PermissionReportsExport,
			PermissionOrderView,
			PermissionPaymentView,
		}
	case AdminRoleScannerOperator:
		return []AdminPermission{
			PermissionScannerView,
		}
	default:
		return []AdminPermission{}
	}
}

// HasPermission checks if the admin user has a specific permission
func (au *AdminUser) HasPermission(permission AdminPermission) bool {
	for _, p := range au.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// IsLocked checks if the admin account is locked
func (au *AdminUser) IsLocked() bool {
	return au.LockedUntil != nil && time.Now().Before(*au.LockedUntil)
}

// CanLogin checks if the admin can login
func (au *AdminUser) CanLogin() bool {
	return au.IsActive && !au.IsLocked()
}

// IncrementFailedAttempts increments failed login attempts
func (au *AdminUser) IncrementFailedAttempts() {
	au.LoginAttempts++
	now := time.Now()
	au.UpdatedAt = now
	
	// Lock account after 5 failed attempts
	if au.LoginAttempts >= 5 {
		lockUntil := now.Add(30 * time.Minute)
		au.LockedUntil = &lockUntil
	}
}

// ResetFailedAttempts resets failed login attempts
func (au *AdminUser) ResetFailedAttempts() {
	au.LoginAttempts = 0
	au.LockedUntil = nil
	au.UpdatedAt = time.Now()
}

// UpdateLastLogin updates the last login timestamp
func (au *AdminUser) UpdateLastLogin() {
	now := time.Now()
	au.LastLogin = &now
	au.UpdatedAt = now
}

// Validate validates the admin user
func (au *AdminUser) Validate() error {
	if au.Email == "" {
		return NewValidationError("email", "email is required")
	}
	
	if au.FirstName == "" {
		return NewValidationError("first_name", "first name is required")
	}
	
	if au.LastName == "" {
		return NewValidationError("last_name", "last name is required")
	}
	
	if au.Role == "" {
		return NewValidationError("role", "role is required")
	}
	
	return nil
}

// GetFullName returns the full name of the admin user
func (au *AdminUser) GetFullName() string {
	return au.FirstName + " " + au.LastName
}

// GetPermissionStrings returns the permissions as a slice of strings
func (au *AdminUser) GetPermissionStrings() []string {
	permissions := make([]string, len(au.Permissions))
	for i, perm := range au.Permissions {
		permissions[i] = string(perm)
	}
	return permissions
}

