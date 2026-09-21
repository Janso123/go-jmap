package emailsubmission

import "github.com/Janso123/go-jmap"

// Get changes over an email submission query
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.4
type QueryChanges struct {
	jmap.QueryChanges[EmailSubmission]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of EmailSubmission/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
