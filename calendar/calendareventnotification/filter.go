package calendareventnotification

import (
	"time"

	"git.sr.ht/~rockorager/go-jmap"
)

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// FilterCondition filters event notifications by creation time, type, or event.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.4.1
type FilterCondition struct {
	After *time.Time `json:"after,omitempty"`

	Before *time.Time `json:"before,omitempty"`

	Type Type `json:"type,omitempty"`

	CalendarEventIDs []jmap.ID `json:"calendarEventIds,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
