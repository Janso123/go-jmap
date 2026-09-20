package contactcard

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/contacts"
)

// Query gets a list of contact card IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3
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

func (m *Query) Name() string { return "ContactCard/query" }

func (m *Query) Requires() []jmap.URI { return []jmap.URI{contacts.URI} }

type QueryResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	QueryState string `json:"queryState,omitempty"`

	CanCalculateChanges bool `json:"canCalculateChanges,omitempty"`

	Position uint64 `json:"position,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`

	Total uint64 `json:"total,omitempty"`

	Limit uint64 `json:"limit,omitempty"`
}

func newQueryResponse() jmap.MethodResponse { return &QueryResponse{} }

// SortComparator defines ContactCard query sorting.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3.2
type SortComparator struct {
	Property string `json:"property,omitempty"`

	IsAscending bool `json:"isAscending"`

	Collation jmap.CollationAlgo `json:"collation,omitempty"`
}
