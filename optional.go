package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// Optional is a JSON field that is omitted, JSON null, or a value.
type Optional[T any] struct {
	value T
	null  bool
	set   bool
}

// Some returns an Optional that marshals v.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, set: true}
}

// Null returns an Optional that marshals as JSON null.
func Null[T any]() Optional[T] {
	return Optional[T]{null: true, set: true}
}

// IsZero reports whether the Optional should be omitted (unset).
func (o Optional[T]) IsZero() bool { return !o.set }

// IsNull reports whether the Optional is an explicit JSON null.
func (o Optional[T]) IsNull() bool { return o.set && o.null }

// Value returns the inner value and whether it is a non-null set value.
func (o Optional[T]) Value() (T, bool) {
	if !o.set || o.null {
		var z T
		return z, false
	}
	return o.value, true
}

func (o Optional[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.null {
		return enc.WriteToken(jsontext.Null)
	}
	return jsonv2.MarshalEncode(enc, o.value)
}

func (o *Optional[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == 'n' {
		if _, err := dec.ReadValue(); err != nil {
			return err
		}
		*o = Null[T]()
		return nil
	}
	var v T
	if err := jsonv2.UnmarshalDecode(dec, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.set || o.null {
		return []byte("null"), nil
	}
	return jsonv2.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = Null[T]()
		return nil
	}
	var v T
	if err := jsonv2.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}
