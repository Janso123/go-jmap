package calendarevent_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/calendarevent"
	"github.com/stretchr/testify/require"
)

func TestCalendarEventIsObject(t *testing.T) {
	var _ jmap.Object = calendarevent.CalendarEvent{}
	require.Equal(t, "CalendarEvent", calendarevent.CalendarEvent{}.JMAPType())
	require.Equal(t, []jmap.URI{calendar.URI}, calendarevent.CalendarEvent{}.Requires())
}

func TestCalendarEventMethodNames(t *testing.T) {
	require.Equal(t, "CalendarEvent/get", (&calendarevent.Get{}).Name())
	require.Equal(t, "CalendarEvent/changes", (&calendarevent.Changes{}).Name())
	require.Equal(t, "CalendarEvent/query", (&calendarevent.Query{}).Name())
	require.Equal(t, "CalendarEvent/queryChanges", (&calendarevent.QueryChanges{}).Name())
	require.Equal(t, "CalendarEvent/set", (&calendarevent.Set{}).Name())
	require.Equal(t, "CalendarEvent/copy", (&calendarevent.Copy{}).Name())
	require.Equal(t, "CalendarEvent/parse", (&calendarevent.Parse{}).Name())
	require.Equal(t, []jmap.URI{calendar.ParseURI}, (&calendarevent.Parse{}).Requires())
}
