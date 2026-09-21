package identity

import "github.com/Janso123/go-jmap"

// Modify identities
// https://www.rfc-editor.org/rfc/rfc8621.html#section-6.3
type Set struct {
	jmap.Set[Identity]
}

// SetResponse is the result of Identity/set.
type SetResponse = jmap.SetResponse[Identity]
