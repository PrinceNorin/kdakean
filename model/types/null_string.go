package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

type NullString sql.NullString

func (s *NullString) Scan(val interface{}) error {
	var i sql.NullString
	if err := i.Scan(val); err != nil {
		return err
	}

	if reflect.TypeOf(val) == nil {
		*s = NullString{i.String, false}
	} else {
		*s = NullString{i.String, true}
	}

	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if s.String != "" {
		return driver.Value(s.String), nil
	}
	return nil, nil
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *NullString) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if val != "" {
		s.String = val
		s.Valid = true
	} else {
		s.String = ""
		s.Valid = false
	}
	return err
}
