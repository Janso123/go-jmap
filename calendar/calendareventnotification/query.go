package calendareventnotification

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Query gets notification ids matching the filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.4
type Query struct {
	jmap.Query[CalendarEventNotification]

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

// QueryResponse is the result of CalendarEventNotification/query.
type QueryResponse = jmap.QueryResponse
