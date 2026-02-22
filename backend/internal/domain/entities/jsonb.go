package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONB represents a PostgreSQL JSONB column
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for database writes
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database reads
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprintf("cannot scan %T into JSONB", value))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(map[string]interface{}(j))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (j *JSONB) UnmarshalJSON(data []byte) error {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	*j = result
	return nil
}

// Get returns a value from the JSONB map
func (j JSONB) Get(key string) interface{} {
	if j == nil {
		return nil
	}
	return j[key]
}

// Set sets a value in the JSONB map
func (j JSONB) Set(key string, value interface{}) {
	if j == nil {
		j = make(map[string]interface{})
	}
	j[key] = value
}

// Has checks if a key exists in the JSONB map
func (j JSONB) Has(key string) bool {
	if j == nil {
		return false
	}
	_, exists := j[key]
	return exists
}

// Delete removes a key from the JSONB map
func (j JSONB) Delete(key string) {
	if j != nil {
		delete(j, key)
	}
}

// Keys returns all keys in the JSONB map
func (j JSONB) Keys() []string {
	if j == nil {
		return []string{}
	}
	
	keys := make([]string, 0, len(j))
	for k := range j {
		keys = append(keys, k)
	}
	return keys
}

// IsEmpty checks if the JSONB map is empty
func (j JSONB) IsEmpty() bool {
	return j == nil || len(j) == 0
}

// Clone creates a deep copy of the JSONB map
func (j JSONB) Clone() JSONB {
	if j == nil {
		return make(JSONB)
	}
	
	clone := make(JSONB, len(j))
	for k, v := range j {
		clone[k] = v
	}
	return clone
}

