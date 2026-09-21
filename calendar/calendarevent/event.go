package calendarevent

import (
	"encoding/json"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// CalendarEvent is a JSCalendar Event with JMAP Calendars metadata.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5
type CalendarEvent struct {
	ID *jmap.ID `json:"id,omitempty"`

	BaseEventID *jmap.ID `json:"baseEventId,omitempty"`

	CalendarIDs *map[jmap.ID]bool `json:"calendarIds,omitempty"`

	IsDraft *bool `json:"isDraft,omitempty"`

	IsOrigin *bool `json:"isOrigin,omitempty"`

	UTCStart *jscalendar.UTCDateTime `json:"utcStart,omitempty"`

	UTCEnd *jscalendar.UTCDateTime `json:"utcEnd,omitempty"`

	UseDefaultAlerts bool `json:"useDefaultAlerts,omitempty"`

	MayInviteSelf bool `json:"mayInviteSelf,omitempty"`

	MayInviteOthers bool `json:"mayInviteOthers,omitempty"`

	HideAttendees bool `json:"hideAttendees,omitempty"`

	ICalendar string `json:"iCalendar,omitempty"`

	jscalendar.Event
}

var calendarEventOverlayKeys = [...]string{
	"id",
	"baseEventId",
	"calendarIds",
	"isDraft",
	"isOrigin",
	"utcStart",
	"utcEnd",
	"useDefaultAlerts",
	"mayInviteSelf",
	"mayInviteOthers",
	"hideAttendees",
	"iCalendar",
}

func (e CalendarEvent) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(e.Event)
	if err != nil {
		return nil, err
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	type metadata struct {
		ID               *jmap.ID                `json:"id,omitempty"`
		BaseEventID      *jmap.ID                `json:"baseEventId,omitempty"`
		CalendarIDs      *map[jmap.ID]bool       `json:"calendarIds,omitempty"`
		IsDraft          *bool                   `json:"isDraft,omitempty"`
		IsOrigin         *bool                   `json:"isOrigin,omitempty"`
		UTCStart         *jscalendar.UTCDateTime `json:"utcStart,omitempty"`
		UTCEnd           *jscalendar.UTCDateTime `json:"utcEnd,omitempty"`
		UseDefaultAlerts bool                    `json:"useDefaultAlerts,omitempty"`
		MayInviteSelf    bool                    `json:"mayInviteSelf,omitempty"`
		MayInviteOthers  bool                    `json:"mayInviteOthers,omitempty"`
		HideAttendees    bool                    `json:"hideAttendees,omitempty"`
		ICalendar        string                  `json:"iCalendar,omitempty"`
	}

	meta := metadata{
		ID:               e.ID,
		BaseEventID:      e.BaseEventID,
		CalendarIDs:      e.CalendarIDs,
		IsDraft:          e.IsDraft,
		IsOrigin:         e.IsOrigin,
		UTCStart:         e.UTCStart,
		UTCEnd:           e.UTCEnd,
		UseDefaultAlerts: e.UseDefaultAlerts,
		MayInviteSelf:    e.MayInviteSelf,
		MayInviteOthers:  e.MayInviteOthers,
		HideAttendees:    e.HideAttendees,
		ICalendar:        e.ICalendar,
	}

	metaData, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	var metaObj map[string]json.RawMessage
	if err := json.Unmarshal(metaData, &metaObj); err != nil {
		return nil, err
	}

	for key, value := range metaObj {
		obj[key] = value
	}

	return json.Marshal(obj)
}

func (e *CalendarEvent) UnmarshalJSON(data []byte) error {
	type metadata struct {
		ID               *jmap.ID                `json:"id,omitempty"`
		BaseEventID      *jmap.ID                `json:"baseEventId,omitempty"`
		CalendarIDs      *map[jmap.ID]bool       `json:"calendarIds,omitempty"`
		IsDraft          *bool                   `json:"isDraft,omitempty"`
		IsOrigin         *bool                   `json:"isOrigin,omitempty"`
		UTCStart         *jscalendar.UTCDateTime `json:"utcStart,omitempty"`
		UTCEnd           *jscalendar.UTCDateTime `json:"utcEnd,omitempty"`
		UseDefaultAlerts bool                    `json:"useDefaultAlerts,omitempty"`
		MayInviteSelf    bool                    `json:"mayInviteSelf,omitempty"`
		MayInviteOthers  bool                    `json:"mayInviteOthers,omitempty"`
		HideAttendees    bool                    `json:"hideAttendees,omitempty"`
		ICalendar        string                  `json:"iCalendar,omitempty"`
	}

	var meta metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, key := range calendarEventOverlayKeys {
		delete(raw, key)
	}

	eventData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	var event jscalendar.Event
	if err := json.Unmarshal(eventData, &event); err != nil {
		return err
	}

	e.ID = meta.ID
	e.BaseEventID = meta.BaseEventID
	e.CalendarIDs = meta.CalendarIDs
	e.IsDraft = meta.IsDraft
	e.IsOrigin = meta.IsOrigin
	e.UTCStart = meta.UTCStart
	e.UTCEnd = meta.UTCEnd
	e.UseDefaultAlerts = meta.UseDefaultAlerts
	e.MayInviteSelf = meta.MayInviteSelf
	e.MayInviteOthers = meta.MayInviteOthers
	e.HideAttendees = meta.HideAttendees
	e.ICalendar = meta.ICalendar
	e.Event = event

	return nil
}
