package calendareventnotification_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/calendareventnotification"
	"github.com/stretchr/testify/require"
)

func TestCalendarEventNotificationIsObject(t *testing.T) {
	var _ jmap.Object = calendareventnotification.CalendarEventNotification{}
	require.Equal(t, "CalendarEventNotification", calendareventnotification.CalendarEventNotification{}.JMAPType())
	require.Equal(t, []jmap.URI{calendar.URI}, calendareventnotification.CalendarEventNotification{}.Requires())
}

func TestCalendarEventNotificationMethodNames(t *testing.T) {
	require.Equal(t, "CalendarEventNotification/get", (&calendareventnotification.Get{}).Name())
	require.Equal(t, "CalendarEventNotification/changes", (&calendareventnotification.Changes{}).Name())
	require.Equal(t, "CalendarEventNotification/query", (&calendareventnotification.Query{}).Name())
	require.Equal(t, "CalendarEventNotification/queryChanges", (&calendareventnotification.QueryChanges{}).Name())
	require.Equal(t, "CalendarEventNotification/set", (&calendareventnotification.Set{}).Name())
}
