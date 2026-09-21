package jmap

import (
	"fmt"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// An Invocation represents method calls and responses
type Invocation struct {
	// The name of the method call or response
	Name string
	// Object containing the named arguments for the method or response
	Args any
	// Arbitrary string set by client, echoed back with responses
	CallID string
}

// MarshalJSON keeps encoding/json callers working; json/v2 prefers MarshalJSONTo.
func (i *Invocation) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(i)
}

func (i *Invocation) MarshalJSONTo(enc *jsontext.Encoder) error {
	args := i.Args
	if args == nil {
		args = map[string]any{}
	}
	return jsonv2.MarshalEncode(enc, []any{i.Name, args, i.CallID})
}

// UnmarshalJSON keeps encoding/json callers working; json/v2 prefers UnmarshalJSONFrom.
func (i *Invocation) UnmarshalJSON(data []byte) error {
	return jsonv2.Unmarshal(data, i)
}

func (i *Invocation) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var raw []jsontext.Value
	if err := jsonv2.UnmarshalDecode(dec, &raw); err != nil {
		return err
	}
	if len(raw) != 3 {
		return fmt.Errorf("Not enough values in invocation")
	}
	if err := jsonv2.Unmarshal(raw[0], &i.Name); err != nil {
		return err
	}
	newFn, ok := methods[i.Name]
	if !ok {
		i.Args = &UnknownResponse{Name: i.Name, Raw: raw[1]}
		return jsonv2.Unmarshal(raw[2], &i.CallID)
	}
	i.Args = newFn()
	if err := jsonv2.Unmarshal(raw[1], i.Args); err != nil {
		return err
	}
	return jsonv2.Unmarshal(raw[2], &i.CallID)
}
