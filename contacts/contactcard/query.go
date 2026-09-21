package contactcard

import "github.com/Janso123/go-jmap"

// Query gets a list of contact card IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3
type Query struct {
	jmap.Query[ContactCard]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of ContactCard/query.
type QueryResponse = jmap.QueryResponse
