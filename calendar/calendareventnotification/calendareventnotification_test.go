package calendareventnotification

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:fix inline
func boolPtr(v bool) *bool { return new(v) }

func TestCalendarEventNotificationMarshal(t *testing.T) {
	created := time.Date(2026, time.September, 19, 16, 0, 0, 0, time.UTC)
	data, err := jsonv2.Marshal(&CalendarEventNotification{
		ID:      "cn1",
		Created: jmap.UTCDatePtr(created),
		ChangedBy: &Person{
			Name:            "Jane Doe",
			Email:           jmap.Some("jane@example.com"),
			PrincipalID:     jmap.Some(jmap.ID("p1")),
			CalendarAddress: jmap.Some("mailto:jane@example.com"),
		},
		Comment:         jmap.Some("Rescheduled due to travel"),
		Type:            TypeUpdated,
		CalendarEventID: "ev1",
		IsDraft:         new(true),
		Event: &jscalendar.Event{
			UID:   "urn:uuid:ev1",
			Title: "Planning",
			Start: "2026-09-19T10:00:00",
		},
		EventPatch: jmap.Patch{
			"title": "Updated planning",
		},
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id":"cn1",
		"created":"2026-09-19T16:00:00Z",
		"changedBy":{
			"name":"Jane Doe",
			"email":"jane@example.com",
			"principalId":"p1",
			"calendarAddress":"mailto:jane@example.com"
		},
		"comment":"Rescheduled due to travel",
		"type":"updated",
		"calendarEventId":"ev1",
		"isDraft":true,
		"event":{
			"uid":"urn:uuid:ev1",
			"title":"Planning",
			"start":"2026-09-19T10:00:00"
		},
		"eventPatch":{"title":"Updated planning"}
	}`, string(data))
}
