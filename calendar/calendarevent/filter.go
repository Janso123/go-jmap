package calendarevent

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
)

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// FilterCondition defines CalendarEvent query filters.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.11.1
type FilterCondition struct {
	InCalendar jmap.ID `json:"inCalendar,omitempty"`

	After jscalendar.LocalDateTime `json:"after,omitempty"`

	Before jscalendar.LocalDateTime `json:"before,omitempty"`

	Text string `json:"text,omitempty"`

	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`

	Location string `json:"location,omitempty"`

	Owner string `json:"owner,omitempty"`

	Attendee string `json:"attendee,omitempty"`

	UID string `json:"uid,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
