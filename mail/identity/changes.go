package identity

import "github.com/Janso123/go-jmap"

// Get identity changes
// https://www.rfc-editor.org/rfc/rfc8621.html#section-6.2
type Changes struct {
	jmap.Changes[Identity]
}

// ChangesResponse is the result of Identity/changes.
type ChangesResponse = jmap.ChangesResponse
