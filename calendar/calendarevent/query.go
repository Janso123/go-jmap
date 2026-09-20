package calendarevent

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
)

// Query gets a list of event IDs based on filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.11
type Query struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Filter Filter `json:"filter,omitempty"`

	Sort []*SortComparator `json:"sort,omitempty"`

	Position int64 `json:"position,omitempty"`

	Anchor jmap.ID `json:"anchor,omitempty"`

	AnchorOffset int64 `json:"anchorOffset,omitempty"`

	Limit uint64 `json:"limit,omitempty"`

	CalculateTotal bool `json:"calculateTotal,omitempty"`

	ExpandRecurrences bool `json:"expandRecurrences,omitempty"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitempty"`
}

func (m *Query) Name() string { return "CalendarEvent/query" }

func (m *Query) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type QueryResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	QueryState string `json:"queryState,omitempty"`

	CanCalculateChanges bool `json:"canCalculateChanges,omitempty"`

	Position uint64 `json:"position,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`

	Total uint64 `json:"total,omitempty"`

	Limit uint64 `json:"limit,omitempty"`
}

func newQueryResponse() jmap.MethodResponse { return &QueryResponse{} }

type SortComparator struct {
	Property string `json:"property,omitempty"`

	IsAscending bool `json:"isAscending"`

	Collation jmap.CollationAlgo `json:"collation,omitempty"`
}
