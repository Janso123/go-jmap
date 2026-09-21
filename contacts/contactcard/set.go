package contactcard

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys contact cards.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.5
type Set struct {
	jmap.Set[ContactCard]
}

// SetResponse is the result of ContactCard/set.
type SetResponse = jmap.SetResponse[ContactCard]
