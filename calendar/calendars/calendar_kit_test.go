package calendars_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/calendars"
	"github.com/stretchr/testify/require"
)

func TestCalendarIsObject(t *testing.T) {
	var _ jmap.Object = calendars.Calendar{}
	require.Equal(t, "Calendar", calendars.Calendar{}.JMAPType())
	require.Equal(t, []jmap.URI{calendar.URI}, calendars.Calendar{}.Requires())
}

func TestCalendarMethodNames(t *testing.T) {
	require.Equal(t, "Calendar/get", (&calendars.Get{}).Name())
	require.Equal(t, "Calendar/changes", (&calendars.Changes{}).Name())
	require.Equal(t, "Calendar/set", (&calendars.Set{}).Name())
}
