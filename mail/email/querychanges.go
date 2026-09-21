package email

import "github.com/Janso123/go-jmap"

// Get changes to an email query since a given state
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.5
type QueryChanges struct {
	jmap.QueryChanges[Email]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`

	CollapseThreads bool `json:"collapseThreads,omitzero"`
}

func (q *QueryChanges) Requires() []jmap.URI {
	return mailRequires(filterNeedsSMIME(q.Filter))
}

// QueryChangesResponse is the result of Email/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
