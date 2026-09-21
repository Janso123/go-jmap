package quota

import "github.com/Janso123/go-jmap"

// Get changes on a quota query.
type QueryChanges struct {
	jmap.QueryChanges[Quota]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of Quota/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
