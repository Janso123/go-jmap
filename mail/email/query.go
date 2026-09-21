package email

import "github.com/Janso123/go-jmap"

// Get list of email IDs based on filter and sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.4
type Query struct {
	jmap.Query[Email]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`

	CollapseThreads bool `json:"collapseThreads,omitzero"`
}

// QueryResponse is the result of Email/query.
type QueryResponse = jmap.QueryResponse
