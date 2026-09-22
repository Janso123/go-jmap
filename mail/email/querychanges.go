package email

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get changes to an email query since a given state
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.5
type QueryChanges struct {
	jmap.QueryChanges[Email]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`

	CollapseThreads bool `json:"collapseThreads,omitzero"`
}

func (q *QueryChanges) Requires() []jmap.URI {
	return mailRequires(FilterNeedsSMIME(q.Filter))
}

type queryChangesAlias QueryChanges

func (q *QueryChanges) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s struct {
		queryChangesAlias
		Filter jsontext.Value `json:"filter"`
	}
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	f, err := jmap.UnmarshalFilter[FilterCondition](s.Filter)
	if err != nil {
		return err
	}
	*q = QueryChanges(s.queryChangesAlias)
	q.Filter = f
	return nil
}

// QueryChangesResponse is the result of Email/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
