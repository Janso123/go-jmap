package calendar

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

// CalendarAlert is a push notification for a triggered calendar alert
// (draft-ietf-jmap-calendars-29 §6.4). SSE event name: calendarAlert.
type CalendarAlert struct {
	Type string `json:"@type,omitzero"`

	AccountID jmap.ID `json:"accountId,omitzero"`

	CalendarEventID jmap.ID `json:"calendarEventId,omitzero"`

	UID string `json:"uid,omitzero"`

	RecurrenceID *jscalendar.LocalDateTime `json:"recurrenceId,omitzero"`

	AlertID string `json:"alertId,omitzero"`
}
