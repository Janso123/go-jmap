package calendarevent

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryInvokeWithExpandedRecurrences(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Query{
		Account: "u1",
		Filter: &FilterCondition{
			InCalendar: "cal1",
			After:      "2026-01-01T00:00:00",
			Before:     "2026-02-01T00:00:00",
		},
		Sort: []*SortComparator{
			{Property: "start", IsAscending: true},
		},
		Limit:             10,
		ExpandRecurrences: true,
		TimeZone:          jscalendar.TimeZoneID("Europe/Warsaw"),
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEvent/query",{"accountId":"u1","limit":10,"filter":{"inCalendar":"cal1","after":"2026-01-01T00:00:00","before":"2026-02-01T00:00:00"},"sort":[{"property":"start","isAscending":true}],"expandRecurrences":true,"timeZone":"Europe/Warsaw"},"0"]]}`,
		string(data))
}

func TestQueryRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Query{}).Requires())
}
