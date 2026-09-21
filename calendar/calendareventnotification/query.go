package calendareventnotification

import "github.com/Janso123/go-jmap"

// Query gets notification ids matching the filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.4
type Query struct {
	jmap.Query[CalendarEventNotification]

	Filter Filter `json:"filter,omitzero"`

	Sort []*SortComparator `json:"sort,omitzero"`
}

// QueryResponse is the result of CalendarEventNotification/query.
type QueryResponse = jmap.QueryResponse

type SortComparator struct {
	// The name of the property on the CalendarEventNotification objects to compare.
	Property string `json:"property,omitzero"`

	// IsAscending has no omitzero: RFC default is true, so Desc must emit false.
	IsAscending bool `json:"isAscending"`

	// The identifier, as registered in the collation registry defined in RFC4790.
	Collation jmap.CollationAlgo `json:"collation,omitzero"`
}
