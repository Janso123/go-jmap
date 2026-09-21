package mailbox

import "github.com/Janso123/go-jmap"

// Get changes on a mailbox query
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.4
type QueryChanges struct {
	jmap.QueryChanges[Mailbox]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of Mailbox/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
