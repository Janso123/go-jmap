package sharenotification

import "github.com/Janso123/go-jmap"

// Get changes on a share notification query.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.5
type QueryChanges struct {
	jmap.QueryChanges[ShareNotification]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of ShareNotification/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
