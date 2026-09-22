package calendars_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/calendar/calendars"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/stretchr/testify/require"
)

func TestCalendarTimeZoneNullRoundTrip(t *testing.T) {
	t.Parallel()
	var cal calendars.Calendar
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"c1","name":"Work","timeZone":null}`), &cal))
	require.True(t, cal.TimeZone.IsNull())

	b, err := jsonv2.Marshal(&cal)
	require.NoError(t, err)
	require.Contains(t, string(b), `"timeZone":null`)
}

func TestCalendarTimeZoneValue(t *testing.T) {
	t.Parallel()
	var cal calendars.Calendar
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"timeZone":"Europe/Warsaw"}`), &cal))
	v, ok := cal.TimeZone.Value()
	require.True(t, ok)
	require.Equal(t, jscalendar.TimeZoneID("Europe/Warsaw"), v)
	require.False(t, cal.TimeZone.IsNull())
}
