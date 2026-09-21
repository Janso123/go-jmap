package calendar_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/stretchr/testify/require"
)

func TestCalendarAlertEventTypes(t *testing.T) {
	require.Equal(t, jmap.EventType("CalendarEvent"), calendar.CalendarEvent)
	require.Equal(t, jmap.EventType("CalendarAlert"), calendar.CalendarAlertEvent)
	require.Equal(t, jmap.EventType("ParticipantIdentity"), calendar.ParticipantIdentity)
}

func TestCalendarAlertFromSSEData(t *testing.T) {
	raw := `{"@type":"CalendarAlert","accountId":"a","calendarEventId":"e","uid":"u","alertId":"al"}`
	var a calendar.CalendarAlert
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &a))
	require.Equal(t, jmap.ID("e"), a.CalendarEventID)
}
