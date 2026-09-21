package mailbox

import "github.com/Janso123/go-jmap"

// Create, delete & modify mailboxes
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.5
type Set struct {
	jmap.Set[Mailbox]

	OnDestroyRemoveEmails bool `json:"onDestroyRemoveEmails,omitzero"`
}

// SetResponse is the result of Mailbox/set.
type SetResponse = jmap.SetResponse[Mailbox]
