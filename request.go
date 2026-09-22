package jmap

import (
	"context"
	"fmt"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

type Request struct {
	// The context to make the request with
	Context context.Context `json:"-"`

	// The JMAP capabilities the request should use
	Using []URI `json:"using"`

	// A slice of methods the server will process. These will be processed
	// sequentially
	Calls []*Invocation `json:"methodCalls"`

	// A map of (client-specified) creation ID to the ID the server assigned
	// when a record was successfully created.
	CreatedIDs map[ID]ID `json:"createdIds,omitzero"`
}

// argsValidator is implemented by method arguments that can spot a local
// RFC 8620 violation before the request goes out. Get, Changes, Query,
// QueryChanges, Set, and Copy implement it, so the thin types embedding them do
// too.
type argsValidator interface {
	ValidateArgs() error
}

type requestPlain Request

// ValidateArgs reports the first local RFC 8620 violation among the arguments of
// the queued calls. The RFC 8887 WebSocket framing builds its own wire struct,
// so it calls this directly instead of marshaling the Request.
func (r *Request) ValidateArgs() error {
	for _, call := range r.Calls {
		if call == nil {
			continue
		}
		if v, ok := call.Args.(argsValidator); ok {
			if err := v.ValidateArgs(); err != nil {
				return fmt.Errorf("jmap: call %q (%s): %w", call.CallID, call.Name, err)
			}
		}
	}
	return nil
}

// MarshalJSONTo validates every call's arguments before encoding the request.
// Request is never embedded by the thin packages, so this method is not
// promoted anywhere and cannot swallow a wrapper's own fields.
func (r *Request) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := r.ValidateArgs(); err != nil {
		return err
	}
	return jsonv2.MarshalEncode(enc, (*requestPlain)(r))
}

// Invoke a method. Each call to Invoke will add the passed Method to the
// Request. The Requires method will be called and added to the request. The
// CallID of the Method is returned. CallIDs are assigned as the hex
// representation of the index of the call, eg "0"
func (r *Request) Invoke(m Method) string {
	i := &Invocation{
		Name:   m.Name(),
		Args:   m,
		CallID: fmt.Sprintf("%x", len(r.Calls)),
	}
	r.Calls = append(r.Calls, i)

	r.Using = mergeURIs(r.Using, m.Requires())
	return i.CallID
}

func mergeURIs(target []URI, opts []URI) []URI {
	seen := make(map[URI]bool, len(target)+len(opts))
	uris := make([]URI, 0, len(target)+len(opts))
	for _, k := range target {
		if seen[k] {
			continue
		}
		seen[k] = true
		uris = append(uris, k)
	}
	for _, k := range opts {
		if seen[k] {
			continue
		}
		seen[k] = true
		uris = append(uris, k)
	}
	return uris
}
