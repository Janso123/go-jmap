package calendarevent

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// FilterCondition defines CalendarEvent query filters.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.11.1
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	InCalendar jmap.ID `json:"inCalendar,omitzero"`

	After jscalendar.LocalDateTime `json:"after,omitzero"`

	Before jscalendar.LocalDateTime `json:"before,omitzero"`

	Text string `json:"text,omitzero"`

	Title string `json:"title,omitzero"`

	Description string `json:"description,omitzero"`

	Location string `json:"location,omitzero"`

	Owner string `json:"owner,omitzero"`

	Attendee string `json:"attendee,omitzero"`

	UID string `json:"uid,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
