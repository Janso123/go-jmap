package principal

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys principals.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type Set struct {
	jmap.Set[Principal]
}

// SetResponse is the result of Principal/set.
type SetResponse = jmap.SetResponse[Principal]
