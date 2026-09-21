package email

import "github.com/Janso123/go-jmap"

// Create, delete or update emails
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.6
type Set struct {
	jmap.Set[Email]
}

// SetResponse is the result of Email/set.
type SetResponse = jmap.SetResponse[Email]
