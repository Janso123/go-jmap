package quota

import "github.com/Janso123/go-jmap"

// Get a list of quota IDs based on filter and sort criteria.
type Query struct {
	jmap.Query[Quota]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of Quota/query.
type QueryResponse = jmap.QueryResponse
