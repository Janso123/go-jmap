package emailsubmission

import "github.com/Janso123/go-jmap"

// List email submission IDs based on filter and sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.3
type Query struct {
	jmap.Query[EmailSubmission]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of EmailSubmission/query.
type QueryResponse = jmap.QueryResponse
