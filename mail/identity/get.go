package identity

import "github.com/Janso123/go-jmap"

// Get details identity details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-6.1
type Get struct {
	jmap.Get[Identity]
}

// GetResponse is the result of Identity/get.
type GetResponse = jmap.GetResponse[Identity]
