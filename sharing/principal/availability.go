package principal

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/calendarevent"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/Janso123/go-jmap/sharing"
)

// BusyStatus describes why a principal is unavailable for scheduling.
type BusyStatus string

const (
	BusyStatusConfirmed   BusyStatus = "confirmed"
	BusyStatusTentative   BusyStatus = "tentative"
	BusyStatusUnavailable BusyStatus = "unavailable"
)

// BusyPeriod describes a span of time when a principal is unavailable.
type BusyPeriod struct {
	UTCStart jscalendar.UTCDateTime `json:"utcStart,omitzero"`

	UTCEnd jscalendar.UTCDateTime `json:"utcEnd,omitzero"`

	BusyStatus BusyStatus `json:"busyStatus,omitzero"`

	Event jmap.Optional[calendarevent.CalendarEvent] `json:"event,omitzero"`

	AccountID jmap.Optional[jmap.ID] `json:"accountId,omitzero"`
}

// GetAvailability calculates scheduling availability for a principal.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-2.2
type GetAvailability struct {
	Account jmap.ID `json:"accountId,omitzero"`

	ID jmap.ID `json:"id,omitzero"`

	UTCStart jscalendar.UTCDateTime `json:"utcStart,omitzero"`

	UTCEnd jscalendar.UTCDateTime `json:"utcEnd,omitzero"`

	ShowDetails bool `json:"showDetails"`

	EventProperties jmap.Optional[[]string] `json:"eventProperties,omitzero"`
}

func (m *GetAvailability) Name() string { return "Principal/getAvailability" }

func (m *GetAvailability) Requires() []jmap.URI {
	return []jmap.URI{sharing.URI, calendar.AvailabilityURI}
}

type GetAvailabilityResponse struct {
	List []*BusyPeriod `json:"list,omitzero"`
}

func newGetAvailabilityResponse() jmap.MethodResponse { return &GetAvailabilityResponse{} }
