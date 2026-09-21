package principal

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/sharing"
)

// Get a list of principal IDs based on filter and sort criteria.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type Query struct {
	jmap.Query[Principal]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

func (q *Query) Requires() []jmap.URI {
	uris := []jmap.URI{sharing.URI}
	if filterHasCalendarAddress(q.Filter) {
		uris = append(uris, calendar.AvailabilityURI)
	}
	return uris
}

func filterHasCalendarAddress(f jmap.Filter) bool {
	if f == nil {
		return false
	}
	switch x := f.(type) {
	case *FilterCondition:
		return x.CalendarAddress != ""
	case FilterCondition:
		return x.CalendarAddress != ""
	case *jmap.FilterOperator:
		for _, c := range x.Conditions {
			if filterHasCalendarAddress(c) {
				return true
			}
		}
	case jmap.FilterOperator:
		for _, c := range x.Conditions {
			if filterHasCalendarAddress(c) {
				return true
			}
		}
	}
	return false
}

// QueryResponse is the result of Principal/query.
type QueryResponse = jmap.QueryResponse
