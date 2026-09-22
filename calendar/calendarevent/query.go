package calendarevent

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// Query gets a list of event IDs based on filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.11
type Query struct {
	jmap.Query[CalendarEvent]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`

	ExpandRecurrences bool `json:"expandRecurrences,omitzero"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitzero"`
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

// QueryResponse is the result of CalendarEvent/query.
type QueryResponse = jmap.QueryResponse
