package calendarevent

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// Query gets a list of event IDs based on filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.11
type Query struct {
	jmap.Query[CalendarEvent]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`

	ExpandRecurrences bool `json:"expandRecurrences,omitzero"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitzero"`
}

// QueryResponse is the result of CalendarEvent/query.
type QueryResponse = jmap.QueryResponse

type SortComparator struct {
	Property string `json:"property,omitzero"`

	// IsAscending has no omitzero: RFC default is true, so Desc must emit false.
	IsAscending bool `json:"isAscending"`

	Collation jmap.CollationAlgo `json:"collation,omitzero"`
}
