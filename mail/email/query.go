package email

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get list of email IDs based on filter and sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.4
type Query struct {
	jmap.Query[Email]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`

	CollapseThreads bool `json:"collapseThreads,omitzero"`
}

func (q *Query) Requires() []jmap.URI {
	return mailRequires(FilterNeedsSMIME(q.Filter))
}

type queryAlias Query

func (q *Query) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s struct {
		queryAlias
		Filter jsontext.Value `json:"filter"`
	}
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	f, err := jmap.UnmarshalFilter[FilterCondition](s.Filter)
	if err != nil {
		return err
	}
	*q = Query(s.queryAlias)
	q.Filter = f
	return nil
}

// QueryResponse is the result of Email/query.
type QueryResponse = jmap.QueryResponse
