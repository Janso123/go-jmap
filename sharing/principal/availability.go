package principal

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
	"git.sr.ht/~rockorager/go-jmap/calendar/calendarevent"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
	"git.sr.ht/~rockorager/go-jmap/sharing"
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
	UTCStart jscalendar.UTCDateTime `json:"utcStart,omitempty"`

	UTCEnd jscalendar.UTCDateTime `json:"utcEnd,omitempty"`

	BusyStatus BusyStatus `json:"busyStatus,omitempty"`

	Event *calendarevent.CalendarEvent `json:"event,omitempty"`

	AccountID *jmap.ID `json:"accountId,omitempty"`
}

// GetAvailability calculates scheduling availability for a principal.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-2.2
type GetAvailability struct {
	Account jmap.ID `json:"accountId,omitempty"`

	ID jmap.ID `json:"id,omitempty"`

	UTCStart jscalendar.UTCDateTime `json:"utcStart,omitempty"`

	UTCEnd jscalendar.UTCDateTime `json:"utcEnd,omitempty"`

	ShowDetails bool `json:"showDetails,omitempty"`

	EventProperties []string `json:"eventProperties,omitempty"`
}

func (m *GetAvailability) Name() string { return "Principal/getAvailability" }

func (m *GetAvailability) Requires() []jmap.URI {
	return []jmap.URI{sharing.URI, calendar.AvailabilityURI}
}

type GetAvailabilityResponse struct {
	List []*BusyPeriod `json:"list,omitempty"`
}

func newGetAvailabilityResponse() jmap.MethodResponse { return &GetAvailabilityResponse{} }
