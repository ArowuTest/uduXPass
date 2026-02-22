package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// AdminPermissionArray is a custom type that implements sql.Scanner and driver.Valuer
// for PostgreSQL JSONB arrays containing AdminPermission values.
// This allows seamless integration with sqlx and PostgreSQL JSONB type.
type AdminPermissionArray []AdminPermission

// Scan implements the sql.Scanner interface for reading PostgreSQL JSONB arrays
// This method is called when sqlx scans a database row into the struct
func (a *AdminPermissionArray) Scan(src interface{}) error {
	// Handle NULL values
	if src == nil {
		*a = AdminPermissionArray{}
		return nil
	}

	// Convert the source to bytes
	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("incompatible type for AdminPermissionArray: %T", src)
	}

	// Unmarshal JSON array into []string first
	var stringArray []string
	if err := json.Unmarshal(source, &stringArray); err != nil {
		return fmt.Errorf("failed to unmarshal admin permissions array: %w", err)
	}

	// Convert []string to []AdminPermission
	permissions := make([]AdminPermission, len(stringArray))
	for i, s := range stringArray {
		permissions[i] = AdminPermission(s)
	}

	*a = permissions
	return nil
}

// Value implements the driver.Valuer interface for writing to PostgreSQL JSONB arrays
// This method is called when sqlx writes the struct to the database
func (a AdminPermissionArray) Value() (driver.Value, error) {
	// Handle empty arrays
	if len(a) == 0 {
		return json.Marshal([]string{})
	}

	// Convert []AdminPermission to []string
	stringArray := make([]string, len(a))
	for i, p := range a {
		stringArray[i] = string(p)
	}

	// Marshal to JSON
	return json.Marshal(stringArray)
}

// MarshalJSON implements the json.Marshaler interface
// This ensures proper JSON serialization when the struct is returned in API responses
func (a AdminPermissionArray) MarshalJSON() ([]byte, error) {
	// Convert []AdminPermission to []string for JSON serialization
	stringArray := make([]string, len(a))
	for i, p := range a {
		stringArray[i] = string(p)
	}
	return json.Marshal(stringArray)
}

// UnmarshalJSON implements the json.Unmarshaler interface
// This ensures proper JSON deserialization when the struct is received in API requests
func (a *AdminPermissionArray) UnmarshalJSON(data []byte) error {
	var stringArray []string
	if err := json.Unmarshal(data, &stringArray); err != nil {
		return fmt.Errorf("failed to unmarshal admin permissions: %w", err)
	}

	permissions := make([]AdminPermission, len(stringArray))
	for i, s := range stringArray {
		permissions[i] = AdminPermission(s)
	}

	*a = permissions
	return nil
}

// String returns a string representation of the permissions array
// Useful for logging and debugging
func (a AdminPermissionArray) String() string {
	if len(a) == 0 {
		return "[]"
	}

	permissions := make([]string, len(a))
	for i, p := range a {
		permissions[i] = string(p)
	}

	return "[" + strings.Join(permissions, ", ") + "]"
}

// Contains checks if the array contains a specific permission
// Useful for permission checking in business logic
func (a AdminPermissionArray) Contains(permission AdminPermission) bool {
	for _, p := range a {
		if p == permission {
			return true
		}
	}
	return false
}

// Add adds a permission to the array if it doesn't already exist
// Returns true if the permission was added, false if it already existed
func (a *AdminPermissionArray) Add(permission AdminPermission) bool {
	if a.Contains(permission) {
		return false
	}
	*a = append(*a, permission)
	return true
}

// Remove removes a permission from the array if it exists
// Returns true if the permission was removed, false if it didn't exist
func (a *AdminPermissionArray) Remove(permission AdminPermission) bool {
	for i, p := range *a {
		if p == permission {
			*a = append((*a)[:i], (*a)[i+1:]...)
			return true
		}
	}
	return false
}

// ToSlice returns a regular []AdminPermission slice
// Useful when you need to work with the permissions as a standard slice
func (a AdminPermissionArray) ToSlice() []AdminPermission {
	return []AdminPermission(a)
}

// FromSlice creates an AdminPermissionArray from a []AdminPermission slice
// Useful for converting from standard slices to the custom type
func FromSlice(permissions []AdminPermission) AdminPermissionArray {
	return AdminPermissionArray(permissions)
}

// HasAny checks if the array contains any of the specified permissions
// Useful for permission checking when multiple permissions are acceptable
func (a AdminPermissionArray) HasAny(permissions ...AdminPermission) bool {
	for _, required := range permissions {
		if a.Contains(required) {
			return true
		}
	}
	return false
}

// HasAll checks if the array contains all of the specified permissions
// Useful for permission checking when all permissions are required
func (a AdminPermissionArray) HasAll(permissions ...AdminPermission) bool {
	for _, required := range permissions {
		if !a.Contains(required) {
			return false
		}
	}
	return true
}
