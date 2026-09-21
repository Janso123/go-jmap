package calendarevent

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// Get calendar event details.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.7
type Get struct {
	jmap.Get[CalendarEvent]

	RecurrenceOverridesBefore *jscalendar.UTCDateTime `json:"recurrenceOverridesBefore,omitzero"`

	RecurrenceOverridesAfter *jscalendar.UTCDateTime `json:"recurrenceOverridesAfter,omitzero"`

	ReduceParticipants bool `json:"reduceParticipants,omitzero"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitzero"`
}

// GetResponse is the result of CalendarEvent/get.
type GetResponse = jmap.GetResponse[CalendarEvent]
