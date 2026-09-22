package calendar_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/calendar/calendarevent"
	"github.com/Janso123/go-jmap/calendar/calendars"
	"github.com/stretchr/testify/require"
)

func TestDefaultAlertsWithTimeNullAndOmitted(t *testing.T) {
	t.Parallel()
	var cal calendars.Calendar
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"defaultAlertsWithTime":null}`), &cal))
	require.True(t, cal.DefaultAlertsWithTime.IsNull())
	b, err := jsonv2.Marshal(&cal)
	require.NoError(t, err)
	require.Contains(t, string(b), `"defaultAlertsWithTime":null`)

	var omitted calendars.Calendar
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"name":"c"}`), &omitted))
	require.True(t, omitted.DefaultAlertsWithTime.IsZero())
	b, err = jsonv2.Marshal(&omitted)
	require.NoError(t, err)
	require.NotContains(t, string(b), "defaultAlertsWithTime")
}

func TestCalendarDescriptionNull(t *testing.T) {
	t.Parallel()
	var cal calendars.Calendar
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"description":null}`), &cal))
	require.True(t, cal.Description.IsNull())
	b, err := jsonv2.Marshal(&cal)
	require.NoError(t, err)
	require.Contains(t, string(b), `"description":null`)
}

func TestParseNotFoundAndNotParsableNull(t *testing.T) {
	t.Parallel()
	var r calendarevent.ParseResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"notFound":null,"notParsable":null}`), &r))
	require.True(t, r.NotFound.IsNull())
	require.True(t, r.NotParsable.IsNull())
	b, err := jsonv2.Marshal(&r)
	require.NoError(t, err)
	require.Contains(t, string(b), `"notFound":null`)
	require.Contains(t, string(b), `"notParsable":null`)
}

func TestRightsMayReadFreeBusyFalseReemitted(t *testing.T) {
	t.Parallel()
	var rights calendars.Rights
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"mayReadFreeBusy":false}`), &rights))
	require.False(t, rights.MayReadFreeBusy)
	b, err := jsonv2.Marshal(&rights)
	require.NoError(t, err)
	require.Contains(t, string(b), `"mayReadFreeBusy":false`)
}

func TestCalendarCreateOmitsIsDefault(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&calendars.Calendar{Name: "c"})
	require.NoError(t, err)
	require.NotContains(t, string(b), "isDefault")
	require.Contains(t, string(b), `"name":"c"`)
}
