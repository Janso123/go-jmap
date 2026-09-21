package sharenotification

import "github.com/Janso123/go-jmap"

// Get a list of share notification IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.4
type Query struct {
	jmap.Query[ShareNotification]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of ShareNotification/query.
type QueryResponse = jmap.QueryResponse
