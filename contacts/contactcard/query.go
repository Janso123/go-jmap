package contactcard

import "github.com/Janso123/go-jmap"

// Query gets a list of contact card IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3
type Query struct {
	jmap.Query[ContactCard]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryResponse is the result of ContactCard/query.
type QueryResponse = jmap.QueryResponse

// SortComparator defines ContactCard query sorting.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3.2
type SortComparator struct {
	Property string `json:"property,omitzero"`

	// IsAscending has no omitzero: RFC default is true, so Desc must emit false.
	IsAscending bool `json:"isAscending"`

	Collation jmap.CollationAlgo `json:"collation,omitzero"`
}
