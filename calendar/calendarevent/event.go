package calendarevent

import (
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// CalendarEvent is a JSCalendar Event with JMAP Calendars metadata.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5
type CalendarEvent struct {
	ID jmap.Optional[jmap.ID] `json:"id,omitzero"`

	BaseEventID jmap.Optional[jmap.ID] `json:"baseEventId,omitzero"`

	CalendarIDs jmap.Optional[map[jmap.ID]bool] `json:"calendarIds,omitzero"`

	IsDraft jmap.Optional[bool] `json:"isDraft,omitzero"`

	IsOrigin jmap.Optional[bool] `json:"isOrigin,omitzero"`

	UTCStart jmap.Optional[jscalendar.UTCDateTime] `json:"utcStart,omitzero"`

	UTCEnd jmap.Optional[jscalendar.UTCDateTime] `json:"utcEnd,omitzero"`

	UseDefaultAlerts bool `json:"useDefaultAlerts,omitzero"`

	MayInviteSelf bool `json:"mayInviteSelf,omitzero"`

	MayInviteOthers bool `json:"mayInviteOthers,omitzero"`

	HideAttendees bool `json:"hideAttendees,omitzero"`

	ICalendar string `json:"iCalendar,omitzero"`

	jscalendar.Event
}

// calendarEventJSON is CalendarEvent without custom marshalers so json/v2 can
// apply Event's Extra embed together with the JMAP metadata keys.
type calendarEventJSON struct {
	ID               jmap.Optional[jmap.ID]                `json:"id,omitzero"`
	BaseEventID      jmap.Optional[jmap.ID]                `json:"baseEventId,omitzero"`
	CalendarIDs      jmap.Optional[map[jmap.ID]bool]       `json:"calendarIds,omitzero"`
	IsDraft          jmap.Optional[bool]                   `json:"isDraft,omitzero"`
	IsOrigin         jmap.Optional[bool]                   `json:"isOrigin,omitzero"`
	UTCStart         jmap.Optional[jscalendar.UTCDateTime] `json:"utcStart,omitzero"`
	UTCEnd           jmap.Optional[jscalendar.UTCDateTime] `json:"utcEnd,omitzero"`
	UseDefaultAlerts bool                                  `json:"useDefaultAlerts,omitzero"`
	MayInviteSelf    bool                                  `json:"mayInviteSelf,omitzero"`
	MayInviteOthers  bool                                  `json:"mayInviteOthers,omitzero"`
	HideAttendees    bool                                  `json:"hideAttendees,omitzero"`
	ICalendar        string                                `json:"iCalendar,omitzero"`
	jscalendar.Event
}

func (e CalendarEvent) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(calendarEventJSON{
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
		Event:            e.Event,
	})
}

func (e *CalendarEvent) UnmarshalJSON(data []byte) error {
	var aux calendarEventJSON
	if err := jsonv2.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.ID = aux.ID
	e.BaseEventID = aux.BaseEventID
	e.CalendarIDs = aux.CalendarIDs
	e.IsDraft = aux.IsDraft
	e.IsOrigin = aux.IsOrigin
	e.UTCStart = aux.UTCStart
	e.UTCEnd = aux.UTCEnd
	e.UseDefaultAlerts = aux.UseDefaultAlerts
	e.MayInviteSelf = aux.MayInviteSelf
	e.MayInviteOthers = aux.MayInviteOthers
	e.HideAttendees = aux.HideAttendees
	e.ICalendar = aux.ICalendar
	e.Event = aux.Event
	return nil
}
