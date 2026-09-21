package calendareventnotification

import "github.com/Janso123/go-jmap"

// FilterCondition filters event notifications by creation time, type, or event.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.4.1
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	After *jmap.UTCDate `json:"after,omitzero"`

	Before *jmap.UTCDate `json:"before,omitzero"`

	Type Type `json:"type,omitzero"`

	CalendarEventIDs []jmap.ID `json:"calendarEventIds,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
