package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Json map[string]interface{}

func (data Json) Value() (driver.Value, error) {
	return json.Marshal(data)
}

func (data *Json) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &data)
}

type JsonArray []Json

func (data JsonArray) Value() (driver.Value, error) {
	return json.Marshal(data)
}

func (data *JsonArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &data)
}
