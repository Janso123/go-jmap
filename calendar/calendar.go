package calendar

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

const (
	// URI is the JMAP Calendars capability.
	URI jmap.URI = "urn:ietf:params:jmap:calendars"

	// ParseURI is the CalendarEvent/parse capability.
	ParseURI jmap.URI = "urn:ietf:params:jmap:calendars:parse"

	// AvailabilityURI is the principals availability capability.
	AvailabilityURI jmap.URI = "urn:ietf:params:jmap:principals:availability"

	// CalendarEvent is the CalendarEvent push/data type.
	CalendarEvent jmap.EventType = "CalendarEvent"

	// CalendarAlertEvent is the CalendarAlert push/data type
	// (SSE event: calendarAlert).
	CalendarAlertEvent jmap.EventType = "CalendarAlert"

	// ParticipantIdentity is the ParticipantIdentity push/data type.
	ParticipantIdentity jmap.EventType = "ParticipantIdentity"
)

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterCapability(&ParseCapability{})
	jmap.RegisterCapability(&AvailabilityCapability{})
}

// Capability describes the JMAP Calendars capability.
// The same type is also used for the empty session capability object.
type Capability struct {
	MaxCalendarsPerEvent     *uint64                `json:"maxCalendarsPerEvent,omitempty"`
	MinDateTime              jscalendar.UTCDateTime `json:"minDateTime,omitempty"`
	MaxDateTime              jscalendar.UTCDateTime `json:"maxDateTime,omitempty"`
	MaxExpandedQueryDuration jscalendar.Duration    `json:"maxExpandedQueryDuration,omitempty"`
	MaxParticipantsPerEvent  *uint64                `json:"maxParticipantsPerEvent,omitempty"`
	MayCreateCalendar        bool                   `json:"mayCreateCalendar,omitempty"`
}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// ParseCapability describes the CalendarEvent/parse capability.
type ParseCapability struct{}

func (c *ParseCapability) URI() jmap.URI { return ParseURI }

func (c *ParseCapability) New() jmap.Capability { return &ParseCapability{} }

// AvailabilityCapability describes the principals availability capability.
// The same type is also used for the empty session capability object.
type AvailabilityCapability struct {
	MaxAvailabilityDuration jscalendar.Duration `json:"maxAvailabilityDuration,omitempty"`
}

func (c *AvailabilityCapability) URI() jmap.URI { return AvailabilityURI }

func (c *AvailabilityCapability) New() jmap.Capability { return &AvailabilityCapability{} }
