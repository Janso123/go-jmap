package jmap

import "time"

// Bool returns a pointer to v for optional boolean JSON fields.
//
//go:fix inline
func Bool(v bool) *bool { return new(v) }

// IDPtr returns a pointer to id for optional ID JSON fields.
//
//go:fix inline
func IDPtr(id ID) *ID { return new(id) }

// UTCDatePtr returns a pointer to a UTCDate for optional date JSON fields.
func UTCDatePtr(t time.Time) *UTCDate {
	d := UTCDate(t)
	return &d
}

// DatePtr returns a pointer to a Date for optional RFC 8620 Date fields.
func DatePtr(t time.Time) *Date {
	d := Date(t)
	return &d
}

// UintPtr returns a pointer to n for optional UnsignedInt JSON fields.
//
//go:fix inline
func UintPtr(n UnsignedInt) *UnsignedInt { return new(n) }

// Uint64Ptr returns a pointer to n for optional UnsignedInt JSON fields.
func Uint64Ptr(n uint64) *UnsignedInt {
	u := UnsignedInt(n)
	return &u
}

// SomeUint returns a set Optional holding n.
func SomeUint(n UnsignedInt) Optional[UnsignedInt] { return Some(n) }

// SomeInt returns a set Optional holding n.
func SomeInt(n Int) Optional[Int] { return Some(n) }
