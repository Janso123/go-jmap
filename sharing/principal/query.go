package principal

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/sharing"
)

// Get a list of principal IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type Query struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Filter Filter `json:"filter,omitempty"`

	Sort []*SortComparator `json:"sort,omitempty"`

	Position int64 `json:"position,omitempty"`

	Anchor jmap.ID `json:"anchor,omitempty"`

	AnchorOffset int64 `json:"anchorOffset,omitempty"`

	Limit uint64 `json:"limit,omitempty"`

	CalculateTotal bool `json:"calculateTotal,omitempty"`
}

func (m *Query) Name() string { return "Principal/query" }

func (m *Query) Requires() []jmap.URI { return []jmap.URI{sharing.URI} }

type QueryResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	QueryState string `json:"queryState,omitempty"`

	CanCalculateChanges bool `json:"canCalculateChanges,omitempty"`

	Position uint64 `json:"position,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`

	Total int64 `json:"total,omitempty"`

	Limit uint64 `json:"limit,omitempty"`
}

func newQueryResponse() jmap.MethodResponse { return &QueryResponse{} }

type SortComparator struct {
	// The name of the property on the Principal objects to compare.
	Property string `json:"property,omitempty"`

	// If true, sort in ascending order.
	IsAscending bool `json:"isAscending"`

	// The identifier, as registered in the collation registry defined in
	// RFC4790, for the algorithm to use when comparing the order of strings.
	Collation jmap.CollationAlgo `json:"collation,omitempty"`
}
