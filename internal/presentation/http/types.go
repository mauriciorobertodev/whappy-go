package http

import "time"

type ISOTime struct {
	time.Time
}

func NewISOTime(t *time.Time) *ISOTime {
	if t == nil {
		return nil
	}

	return &ISOTime{*t}
}

func (t *ISOTime) MarshalJSON() ([]byte, error) {
	if t == nil || t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.UTC().Format(time.RFC3339) + `"`), nil
}
