package jmap

import "time"

// Bool returns a pointer to v for optional boolean JSON fields.
func Bool(v bool) *bool { return &v }

// IDPtr returns a pointer to id for optional ID JSON fields.
func IDPtr(id ID) *ID { return &id }

// UTCDatePtr returns a pointer to a UTCDate for optional date JSON fields.
func UTCDatePtr(t time.Time) *UTCDate {
	d := UTCDate(t)
	return &d
}
