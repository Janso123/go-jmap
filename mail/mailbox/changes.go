package mailbox

import "github.com/Janso123/go-jmap"

// Get mailbox changes for the whole account
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.2
type Changes struct {
	jmap.Changes[Mailbox]
}

// ChangesResponse is the result of Mailbox/changes.
// updatedProperties is Mailbox-specific (RFC 8621 §2.2).
type ChangesResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	OldState string `json:"oldState,omitzero"`

	NewState string `json:"newState,omitzero"`

	HasMoreChanges bool `json:"hasMoreChanges,omitzero"`

	Created []jmap.ID `json:"created,omitzero"`

	Updated []jmap.ID `json:"updated,omitzero"`

	Destroyed []jmap.ID `json:"destroyed,omitzero"`

	UpdatedProperties []string `json:"updatedProperties,omitzero"`
}
