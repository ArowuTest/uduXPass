package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONBArray represents a PostgreSQL JSONB column that stores a JSON array.
// This is distinct from JSONB (which is map-based) and is used for fields
// like gallery_images that store ordered lists of strings.
type JSONBArray []interface{}

// Value implements the driver.Valuer interface for database writes.
func (j JSONBArray) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal([]interface{}(j))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements the sql.Scanner interface for database reads.
// Handles both JSON arrays and JSON objects (for backward compatibility).
func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = JSONBArray{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprintf("cannot scan %T into JSONBArray", value))
	}

	// Try to unmarshal as an array first
	var arr []interface{}
	if err := json.Unmarshal(bytes, &arr); err == nil {
		*j = JSONBArray(arr)
		return nil
	}

	// Fallback: try as a map (legacy data) — wrap values in array
	var obj map[string]interface{}
	if err := json.Unmarshal(bytes, &obj); err == nil {
		result := make(JSONBArray, 0, len(obj))
		for _, v := range obj {
			result = append(result, v)
		}
		*j = result
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into JSONBArray", string(bytes))
}

// MarshalJSON implements the json.Marshaler interface.
func (j JSONBArray) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]interface{}(j))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *JSONBArray) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*j = JSONBArray(arr)
	return nil
}

// Strings returns the array as a slice of strings, filtering non-string values.
func (j JSONBArray) Strings() []string {
	if j == nil {
		return []string{}
	}
	result := make([]string, 0, len(j))
	for _, v := range j {
		if s, ok := v.(string); ok && s != "" {
			result = append(result, s)
		}
	}
	return result
}

// IsEmpty returns true if the array is nil or has no elements.
func (j JSONBArray) IsEmpty() bool {
	return j == nil || len(j) == 0
}
