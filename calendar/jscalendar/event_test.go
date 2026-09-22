package jscalendar

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
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

func TestEventRecurrenceNullStaysNull(t *testing.T) {
	t.Parallel()
	var ev Event
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"recurrenceRule":null,"recurrenceOverrides":null}`), &ev))
	require.True(t, ev.RecurrenceRule.IsNull())
	require.True(t, ev.RecurrenceOverrides.IsNull())
	b, err := jsonv2.Marshal(&ev)
	require.NoError(t, err)
	require.Contains(t, string(b), `"recurrenceRule":null`)
	require.Contains(t, string(b), `"recurrenceOverrides":null`)
}

func TestLinkCIDRoundTripsViaExtra(t *testing.T) {
	t.Parallel()
	const input = `{"cid":"abc","href":"https://e"}`
	var link Link
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &link))
	raw, ok := link.Extra["cid"]
	require.True(t, ok)
	require.JSONEq(t, `"abc"`, string(raw))
	require.Equal(t, "https://e", link.Href)
	out, err := jsonv2.Marshal(&link)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
}

func TestLinkSizeZero(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(Link{Size: new(jmap.UnsignedInt(0))})
	require.NoError(t, err)
	require.Contains(t, string(b), `"size":0`)
}

func TestLinkBlobIDNotOnlyInExtra(t *testing.T) {
	t.Parallel()
	var l Link
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"href":"https://ex/a","blobId":"b1"}`), &l))
	_, inExtra := l.Extra["blobId"]
	require.False(t, inExtra)
	require.Equal(t, jmap.ID("b1"), l.BlobID)
	b, err := jsonv2.Marshal(&l)
	require.NoError(t, err)
	require.Contains(t, string(b), `"blobId":"b1"`)
}
