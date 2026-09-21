package jmap

import (
	"context"
	"fmt"
)

// Call invokes a single method and returns the typed response for that call.
func Call[T MethodResponse](ctx context.Context, c *Client, m Method) (T, error) {
	var zero T
	req := &Request{Context: ctx}
	id := req.Invoke(m)
	resp, err := c.Do(ctx, req)
	if err != nil {
		return zero, err
	}
	return As[T](resp, id)
}

// ByCallID returns the invocation with the given call ID, if present.
func (r *Response) ByCallID(id string) (*Invocation, bool) {
	for _, inv := range r.Responses {
		if inv != nil && inv.CallID == id {
			return inv, true
		}
	}
	return nil, false
}

// As returns the typed Args for the invocation with callID.
func As[T MethodResponse](r *Response, callID string) (T, error) {
	var zero T
	inv, ok := r.ByCallID(callID)
	if !ok {
		return zero, fmt.Errorf("no response for call id %q", callID)
	}
	if me, ok := inv.Args.(*MethodError); ok {
		return zero, me
	}
	v, ok := inv.Args.(T)
	if !ok {
		return zero, fmt.Errorf("response %q has type %T, want %T", callID, inv.Args, zero)
	}
	return v, nil
}
