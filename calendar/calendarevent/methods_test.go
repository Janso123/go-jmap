package calendarevent

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func idPtr(v jmap.ID) *jmap.ID { return &v }

func boolPtr(v bool) *bool { return &v }

func calendarIDsPtr(v map[jmap.ID]bool) *map[jmap.ID]bool { return &v }

func TestCalendarEventJSON(t *testing.T) {
	baseID := jmap.ID("base1")
	start := jscalendar.UTCDateTime("2026-03-01T08:00:00Z")
	end := jscalendar.UTCDateTime("2026-03-01T09:00:00Z")

	event := CalendarEvent{
		ID:          idPtr("ev1"),
		BaseEventID: &baseID,
		CalendarIDs: calendarIDsPtr(map[jmap.ID]bool{"cal1": true}),
		IsDraft:     boolPtr(true),
		IsOrigin:    boolPtr(true),
		UTCStart:    &start,
		UTCEnd:      &end,
		Event: jscalendar.Event{
			UID:      "urn:uuid:ev1",
			Title:    "Planning",
			Start:    "2026-03-01T09:00:00",
			TimeZone: "Europe/Warsaw",
			Duration: "PT1H",
		},
	}

	data, err := json.Marshal(&event)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id":"ev1",
		"baseEventId":"base1",
		"calendarIds":{"cal1":true},
		"isDraft":true,
		"isOrigin":true,
		"utcStart":"2026-03-01T08:00:00Z",
		"utcEnd":"2026-03-01T09:00:00Z",
		"uid":"urn:uuid:ev1",
		"title":"Planning",
		"start":"2026-03-01T09:00:00",
		"timeZone":"Europe/Warsaw",
		"duration":"PT1H"
	}`, string(data))

	var roundTrip CalendarEvent
	require.NoError(t, json.Unmarshal(data, &roundTrip))
	assert.Equal(t, event.ID, roundTrip.ID)
	assert.Equal(t, event.BaseEventID, roundTrip.BaseEventID)
	assert.Equal(t, event.CalendarIDs, roundTrip.CalendarIDs)
	assert.Equal(t, event.IsDraft, roundTrip.IsDraft)
	assert.Equal(t, event.IsOrigin, roundTrip.IsOrigin)
	assert.Equal(t, event.UTCStart, roundTrip.UTCStart)
	assert.Equal(t, event.UTCEnd, roundTrip.UTCEnd)
	assert.Equal(t, event.UID, roundTrip.UID)
	assert.Equal(t, event.Title, roundTrip.Title)
}

func TestCalendarEventMarshalDoesNotLeakClearedOverlayFields(t *testing.T) {
	const input = `{
		"id":"ev1",
		"calendarIds":{"cal1":true},
		"isDraft":true,
		"uid":"urn:uuid:ev1",
		"title":"Planning",
		"start":"2026-03-01T09:00:00",
		"example.com:colorHint":"violet"
	}`

	var event CalendarEvent
	require.NoError(t, json.Unmarshal([]byte(input), &event))

	event.ID = nil
	event.CalendarIDs = nil
	event.IsDraft = nil

	data, err := json.Marshal(&event)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"uid":"urn:uuid:ev1",
		"title":"Planning",
		"start":"2026-03-01T09:00:00",
		"example.com:colorHint":"violet"
	}`, string(data))
}

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}
	before := jscalendar.UTCDateTime("2026-02-01T00:00:00Z")
	after := jscalendar.UTCDateTime("2026-01-01T00:00:00Z")

	id := req.Invoke(&Get{
		Account:                   "u1",
		IDs:                       []jmap.ID{"ev1"},
		Properties:                []string{"title", "utcStart"},
		RecurrenceOverridesBefore: &before,
		RecurrenceOverridesAfter:  &after,
		ReduceParticipants:        true,
		TimeZone:                  "Europe/Warsaw",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/get",{"accountId":"u1","ids":["ev1"],"properties":["title","utcStart"],"recurrenceOverridesBefore":"2026-02-01T00:00:00Z","recurrenceOverridesAfter":"2026-01-01T00:00:00Z","reduceParticipants":true,"timeZone":"Europe/Warsaw"},"0"]]}`,
		string(data))
}

func TestGetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Get{}).Requires())
}

func TestChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Changes{
		Account:    "u1",
		SinceState: "s1",
		MaxChanges: 25,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/changes",{"accountId":"u1","sinceState":"s1","maxChanges":25},"0"]]}`,
		string(data))
}

func TestChangesRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Changes{}).Requires())
}

func TestSetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account: "u1",
		Create: map[jmap.ID]*CalendarEvent{
			"ev1": {
				CalendarIDs: calendarIDsPtr(map[jmap.ID]bool{"cal1": true}),
				IsDraft:     boolPtr(true),
				Event: jscalendar.Event{
					UID:      "urn:uuid:ev1",
					Title:    "Planning",
					Start:    "2026-03-01T09:00:00",
					TimeZone: "Europe/Warsaw",
					Duration: "PT1H",
				},
			},
		},
		SendSchedulingMessages: true,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/set",{"accountId":"u1","create":{"ev1":{"calendarIds":{"cal1":true},"isDraft":true,"uid":"urn:uuid:ev1","title":"Planning","start":"2026-03-01T09:00:00","timeZone":"Europe/Warsaw","duration":"PT1H"}},"sendSchedulingMessages":true},"0"]]}`,
		string(data))
}

func TestSetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Set{}).Requires())
}

func TestCopyInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Copy{
		FromAccount: "u1",
		Account:     "u2",
		Create: map[jmap.ID]*CalendarEvent{
			"ev1": {
				CalendarIDs: calendarIDsPtr(map[jmap.ID]bool{"cal2": true}),
				Event: jscalendar.Event{
					UID:   "urn:uuid:ev1-copy",
					Title: "Planning copy",
					Start: "2026-03-02T09:00:00",
				},
			},
		},
		OnSuccessDestroyOriginal: true,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/copy",{"fromAccountId":"u1","accountId":"u2","create":{"ev1":{"calendarIds":{"cal2":true},"uid":"urn:uuid:ev1-copy","title":"Planning copy","start":"2026-03-02T09:00:00"}},"onSuccessDestroyOriginal":true},"0"]]}`,
		string(data))
}

func TestCopyRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Copy{}).Requires())
}

func TestQueryChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&QueryChanges{
		Account:         "u1",
		Filter:          &FilterCondition{UID: "urn:uuid:ev1"},
		Sort:            []*SortComparator{{Property: "start", IsAscending: true}},
		SinceQueryState: "q1",
		MaxChanges:      10,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/queryChanges",{"accountId":"u1","filter":{"uid":"urn:uuid:ev1"},"sort":[{"property":"start","isAscending":true}],"sinceQueryState":"q1","maxChanges":10},"0"]]}`,
		string(data))
}

func TestQueryChangesRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&QueryChanges{}).Requires())
}

func TestParseInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Parse{
		Account:    "u1",
		BlobIDs:    []jmap.ID{"b1"},
		Properties: []string{"uid", "title", "calendarIds"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars:parse"],"methodCalls":[["CalendarEvent/parse",{"accountId":"u1","blobIds":["b1"],"properties":["uid","title","calendarIds"]},"0"]]}`,
		string(data))
}

func TestParseRequiresParseCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.ParseURI}, (&Parse{}).Requires())
}

func TestParseResponseUnmarshalNullMetadata(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["CalendarEvent/parse",{"accountId":"u1","parsed":{"b1":[{"id":null,"baseEventId":null,"calendarIds":null,"isDraft":null,"isOrigin":null,"uid":"urn:uuid:ev1","title":"Planning","start":"2026-03-01T09:00:00"}]}},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	parseResp, ok := resp.Responses[0].Args.(*ParseResponse)
	require.True(t, ok)
	require.Contains(t, parseResp.Parsed, jmap.ID("b1"))
	require.Len(t, parseResp.Parsed["b1"], 1)

	event := parseResp.Parsed["b1"][0]
	assert.Nil(t, event.ID)
	assert.Nil(t, event.BaseEventID)
	assert.Nil(t, event.CalendarIDs)
	assert.Nil(t, event.IsDraft)
	assert.Nil(t, event.IsOrigin)
	assert.Equal(t, "urn:uuid:ev1", event.UID)
	assert.Equal(t, "Planning", event.Title)
}
