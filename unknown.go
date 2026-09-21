package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// UnknownResponse holds arguments for a method that was not registered.
// Invocation unmarshal uses this instead of returning an error so the rest of
// a Response can still be processed.
type UnknownResponse struct {
	Name string `json:"-"`
	Raw  jsontext.Value
}

// MarshalJSON keeps encoding/json callers working; json/v2 prefers MarshalJSONTo.
func (u UnknownResponse) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(u.Raw)
}

func (u UnknownResponse) MarshalJSONTo(enc *jsontext.Encoder) error {
	if len(u.Raw) == 0 {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteValue(u.Raw)
}
