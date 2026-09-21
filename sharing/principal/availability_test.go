package principal

import (
	"encoding/json"
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
		EventProperties: []string{"id", "title", "calendarIds"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
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
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	availabilityResp, ok := resp.Responses[0].Args.(*GetAvailabilityResponse)
	require.True(t, ok)
	require.Len(t, availabilityResp.List, 1)

	period := availabilityResp.List[0]
	assert.Equal(t, jscalendar.UTCDateTime("2026-09-19T08:00:00Z"), period.UTCStart)
	assert.Equal(t, jscalendar.UTCDateTime("2026-09-19T09:00:00Z"), period.UTCEnd)
	assert.Equal(t, BusyStatusConfirmed, period.BusyStatus)
	assert.Nil(t, period.Event)
	assert.Nil(t, period.AccountID)
}
