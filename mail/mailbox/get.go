package mailbox

import "github.com/Janso123/go-jmap"

// Get mailbox details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.1
type Get struct {
	jmap.Get[Mailbox]
}

// GetResponse is the result of Mailbox/get.
type GetResponse = jmap.GetResponse[Mailbox]
