package model

import (
	"time"
)

// NullTime represents an time.Time value that may be null, represented as an empty string in text, and null in JSON.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

const timeFormat = "2006-01-02 15:04:05"

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *NullTime) UnmarshalText(text []byte) error {
	n.Valid = false

	if len(text) == 0 {
		return nil
	}

	t, err := time.Parse(timeFormat, string(text))
	if err != nil {
		return err
	}

	n.Valid = true
	n.Time = t
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n NullTime) MarshalText() ([]byte, error) {
	var s string

	if n.Valid {
		s = n.Time.Format(timeFormat)
	}

	return []byte(s), nil
}

// MarshalJSON implements the encoding/json.MarshalJSON interface.
func (n NullTime) MarshalJSON() ([]byte, error) {
	s := "null"

	if n.Valid {
		s = `"` + n.Time.Format(timeFormat) + `"`
	}

	return []byte(s), nil
}
