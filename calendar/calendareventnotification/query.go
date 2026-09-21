package calendareventnotification

import "github.com/Janso123/go-jmap"

// Query gets notification ids matching the filter and sort criteria.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.4
type Query struct {
	jmap.Query[CalendarEventNotification]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryResponse is the result of CalendarEventNotification/query.
type QueryResponse = jmap.QueryResponse
