package principal

import "github.com/Janso123/go-jmap"

// Get a list of principal IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type Query struct {
	jmap.Query[Principal]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryResponse is the result of Principal/query.
type QueryResponse = jmap.QueryResponse

type SortComparator struct {
	// The name of the property on the Principal objects to compare.
	Property string `json:"property,omitempty"`

	// If true, sort in ascending order.
	IsAscending bool `json:"isAscending"`

	// The identifier, as registered in the collation registry defined in
	// RFC4790, for the algorithm to use when comparing the order of strings.
	Collation jmap.CollationAlgo `json:"collation,omitempty"`
}
