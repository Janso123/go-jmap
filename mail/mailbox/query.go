package mailbox

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get a list of mailbox IDs based on filter and sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2.3
type Query struct {
	jmap.Query[Mailbox]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`

	SortAsTree bool `json:"sortAsTree,omitzero"`

	FilterAsTree bool `json:"filterAsTree,omitzero"`
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

// QueryResponse is the result of Mailbox/query.
type QueryResponse = jmap.QueryResponse
