package sieve

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get a list of Sieve script IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.5
type Query struct {
	jmap.Query[SieveScript]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
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

// QueryResponse is the result of SieveScript/query.
type QueryResponse = jmap.QueryResponse
