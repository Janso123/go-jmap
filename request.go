package jmap

import (
	"context"
	"fmt"
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
