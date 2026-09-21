package mailbox

import "github.com/Janso123/go-jmap"

// Get a list of mailbox IDs based on filter and sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.3
type Query struct {
	jmap.Query[Mailbox]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`

	SortAsTree bool `json:"sortAsTree,omitzero"`

	FilterAsTree bool `json:"filterAsTree,omitzero"`
}

// QueryResponse is the result of Mailbox/query.
type QueryResponse = jmap.QueryResponse
