package contactcard

import "github.com/Janso123/go-jmap"

// QueryChanges gets changes to a contact card query since a given state.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.4
type QueryChanges struct {
	jmap.QueryChanges[ContactCard]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of ContactCard/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
