package contactcard

import "github.com/Janso123/go-jmap"

// Get contact card details.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.1
type Get struct {
	jmap.Get[ContactCard]
}

// GetResponse is the result of ContactCard/get.
type GetResponse = jmap.GetResponse[ContactCard]
