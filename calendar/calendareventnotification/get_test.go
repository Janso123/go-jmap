package calendareventnotification

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account:    "u1",
		IDs:        jmap.Some([]jmap.ID{"cn1"}),
		Properties: jmap.Some([]string{"type", "calendarEventId", "event"}),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEventNotification/get",{"accountId":"u1","ids":["cn1"],"properties":["type","calendarEventId","event"]},"0"]]}`,
		string(data))
}

func TestGetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Get{}).Requires())
}
