package calendarevent

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
)

// Get calendar event details.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.7
type Get struct {
	Account jmap.ID `json:"accountId,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`

	Properties []string `json:"properties,omitempty"`

	ReferenceIDs *jmap.ResultReference `json:"#ids,omitempty"`

	ReferenceProperties *jmap.ResultReference `json:"#properties,omitempty"`

	RecurrenceOverridesBefore *jscalendar.UTCDateTime `json:"recurrenceOverridesBefore,omitempty"`

	RecurrenceOverridesAfter *jscalendar.UTCDateTime `json:"recurrenceOverridesAfter,omitempty"`

	ReduceParticipants bool `json:"reduceParticipants,omitempty"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitempty"`
}

func (m *Get) Name() string { return "CalendarEvent/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	State string `json:"state,omitempty"`

	List []*CalendarEvent `json:"list,omitempty"`

	NotFound []jmap.ID `json:"notFound,omitempty"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
