package jscalendar

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventExtraRoundTrip(t *testing.T) {
	const input = `{"@type":"Event","uid":"u1","x-unknown":1}`
	var ev Event
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &ev))
	require.Equal(t, "u1", ev.UID)
	_, ok := ev.Extra["x-unknown"]
	require.True(t, ok)
	out, err := jsonv2.Marshal(ev)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
}

func TestEventRoundTripMinimal(t *testing.T) {
	const input = `{
		"@type":"Event",
		"uid":"a8df6573-0474-496d-8496-033ad45d7fea",
		"title":"Some event",
		"start":"2020-01-15T13:00:00",
		"timeZone":"America/New_York",
		"duration":"PT1H"
	}`

	var event Event
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &event))

	data, err := jsonv2.Marshal(&event)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestEventRoundTripNestedMaps(t *testing.T) {
	const input = `{
		"@type":"Event",
		"uid":"715ed4c5-3cf5-427f-927c-db40cdd63894",
		"title":"Departmental meeting",
		"start":"2025-01-07T14:00:00",
		"timeZone":"Australia/Melbourne",
		"duration":"PT1H",
		"organizerCalendarAddress":"mailto:zoe@foobar.example.com",
		"participants":{
			"zoe":{
				"@type":"Participant",
				"name":"Zoe Zelda",
				"calendarAddress":"mailto:zoe@foobar.example.com",
				"participationStatus":"accepted",
				"roles":{"owner":true,"chair":true}
			}
		},
		"alerts":{
			"a1":{
				"@type":"Alert",
				"trigger":{
					"@type":"OffsetTrigger",
					"offset":"-PT15M",
					"relativeTo":"start"
				},
				"action":"display"
			}
		},
		"recurrenceRule":{
			"@type":"RecurrenceRule",
			"frequency":"weekly",
			"byDay":[{"@type":"NDay","day":"tu"}],
			"count":3
		},
		"recurrenceOverrides":{
			"2025-01-14T14:00:00":{
				"title":"Departmental meeting moved",
				"start":"2025-01-14T15:00:00"
			}
		}
	}`

	var event Event
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &event))

	data, err := jsonv2.Marshal(&event)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestEventRoundTripVendorExtensions(t *testing.T) {
	const input = `{
		"@type":"Event",
		"uid":"urn:uuid:vendor-extension",
		"title":"Ada Sync",
		"start":"2026-01-02T03:04:05",
		"example.com:colorHint":"violet",
		"participants":{
			"ada":{
				"name":"Ada Lovelace",
				"calendarAddress":"mailto:ada@example.com",
				"example.com:desk":"Babbage-42"
			}
		}
	}`

	var event Event
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &event))

	data, err := jsonv2.Marshal(&event)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestLinkDisplayIsStringBooleanMap(t *testing.T) {
	const input = `{"@type":"Event","uid":"u1","links":{"l1":{"@type":"Link","href":"https://x/i.png","rel":"icon","display":{"badge":true}}}}`
	var ev Event
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &ev))
	require.Equal(t, true, ev.Links["l1"].Display["badge"])
	out, err := jsonv2.Marshal(ev)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
}

func TestEventPrivacyAndFreeBusy(t *testing.T) {
	ev := Event{Type: "Event", UID: "u1", Privacy: PrivacyPrivate, FreeBusyStatus: FreeBusyBusy, Priority: 5}
	b, err := jsonv2.Marshal(ev)
	require.NoError(t, err)
	require.JSONEq(t, `{"@type":"Event","uid":"u1","priority":5,"freeBusyStatus":"busy","privacy":"private"}`, string(b))
}
