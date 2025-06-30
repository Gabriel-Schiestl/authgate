package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("falha ao converter valor para bytes: %v", value)
	}

	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	return json.Marshal(s)
}