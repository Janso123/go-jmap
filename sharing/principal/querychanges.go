package principal

import "github.com/Janso123/go-jmap"

// Get changes on a principal query.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.4
type QueryChanges struct {
	jmap.QueryChanges[Principal]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of Principal/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
