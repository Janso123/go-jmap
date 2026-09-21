package sieve

import "github.com/Janso123/go-jmap"

// Get a list of Sieve script IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.5
type Query struct {
	jmap.Query[SieveScript]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of SieveScript/query.
type QueryResponse = jmap.QueryResponse
