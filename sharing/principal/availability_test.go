package principal

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/Janso123/go-jmap/sharing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAvailabilityInvoke(t *testing.T) {
	req := &jmap.Request{}
	utcStart := jscalendar.UTCDateTime("2026-09-19T08:00:00Z")
	utcEnd := jscalendar.UTCDateTime("2026-09-19T10:00:00Z")

	id := req.Invoke(&GetAvailability{
		Account:         "u1",
		ID:              "p1",
		UTCStart:        utcStart,
		UTCEnd:          utcEnd,
		ShowDetails:     true,
		EventProperties: jmap.Some([]string{"id", "title", "calendarIds"}),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"using":["urn:ietf:params:jmap:principals","urn:ietf:params:jmap:principals:availability"],"methodCalls":[["Principal/getAvailability",{"accountId":"u1","id":"p1","utcStart":"2026-09-19T08:00:00Z","utcEnd":"2026-09-19T10:00:00Z","showDetails":true,"eventProperties":["id","title","calendarIds"]},"0"]]}`,
		string(data))
}

func TestGetAvailabilityRequiresSharingAndAvailabilityCapabilities(t *testing.T) {
	assert.ElementsMatch(t,
		[]jmap.URI{sharing.URI, calendar.AvailabilityURI},
		(&GetAvailability{}).Requires())
}

func TestGetAvailabilityResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Principal/getAvailability",{"list":[{"utcStart":"2026-09-19T08:00:00Z","utcEnd":"2026-09-19T09:00:00Z","busyStatus":"confirmed","event":null,"accountId":null}]},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	availabilityResp, ok := resp.Responses[0].Args.(*GetAvailabilityResponse)
	require.True(t, ok)
	require.Len(t, availabilityResp.List, 1)

	period := availabilityResp.List[0]
	assert.Equal(t, jscalendar.UTCDateTime("2026-09-19T08:00:00Z"), period.UTCStart)
	assert.Equal(t, jscalendar.UTCDateTime("2026-09-19T09:00:00Z"), period.UTCEnd)
	assert.Equal(t, BusyStatusConfirmed, period.BusyStatus)
	assert.True(t, period.Event.IsNull())
	assert.True(t, period.AccountID.IsNull())
}

func TestBusyPeriodNullEventAndAccountID(t *testing.T) {
	var period BusyPeriod
	raw := `{"utcStart":"2026-09-19T08:00:00Z","utcEnd":"2026-09-19T09:00:00Z","busyStatus":"busy","event":null,"accountId":null}`
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &period))
	require.True(t, period.Event.IsNull())
	require.True(t, period.AccountID.IsNull())

	b, err := jsonv2.Marshal(&period)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(b))
}

func TestGetAvailabilityNullEventPropertiesAndShowDetailsFalse(t *testing.T) {
	b, err := jsonv2.Marshal(&GetAvailability{
		EventProperties: jmap.Null[[]string](),
	})
	require.NoError(t, err)
	require.JSONEq(t, `{"eventProperties":null,"showDetails":false}`, string(b))
}

func TestPrincipalCalendarsCapability(t *testing.T) {
	var p Principal
	raw := `{
		"capabilities": {
			"urn:ietf:params:jmap:calendars": {
				"accountId": "a1",
				"mayGetAvailability": true,
				"mayShareWith": false,
				"calendarAddress": "mailto:ada@example.com"
			}
		}
	}`
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &p))
	cap, ok, err := p.CalendarsCapability()
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, cap)
	id, idOK := cap.AccountID.Value()
	require.True(t, idOK)
	require.Equal(t, jmap.ID("a1"), id)
	require.True(t, cap.MayGetAvailability)
	require.False(t, cap.MayShareWith)
	require.Equal(t, "mailto:ada@example.com", cap.CalendarAddress)

	var absent Principal
	missing, ok, err := absent.CalendarsCapability()
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, missing)

	var nullAccount Principal
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"capabilities":{"urn:ietf:params:jmap:calendars":{"accountId":null,"mayGetAvailability":false,"mayShareWith":true,"calendarAddress":"mailto:room@example.com"}}}`), &nullAccount))
	nullCap, ok, err := nullAccount.CalendarsCapability()
	require.NoError(t, err)
	require.True(t, ok)
	require.True(t, nullCap.AccountID.IsNull())
	require.False(t, nullCap.MayGetAvailability)
	require.True(t, nullCap.MayShareWith)
	require.Equal(t, "mailto:room@example.com", nullCap.CalendarAddress)
}
