package sharenotification_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing"
	"github.com/Janso123/go-jmap/sharing/sharenotification"
	"github.com/stretchr/testify/require"
)

func TestShareNotificationIsObject(t *testing.T) {
	var _ jmap.Object = sharenotification.ShareNotification{}
	require.Equal(t, "ShareNotification", sharenotification.ShareNotification{}.JMAPType())
	require.Equal(t, []jmap.URI{sharing.URI}, sharenotification.ShareNotification{}.Requires())
}

func TestShareNotificationGetName(t *testing.T) {
	require.Equal(t, "ShareNotification/get", (&sharenotification.Get{}).Name())
}

func TestShareNotificationChangesName(t *testing.T) {
	require.Equal(t, "ShareNotification/changes", (&sharenotification.Changes{}).Name())
}

func TestShareNotificationSetName(t *testing.T) {
	require.Equal(t, "ShareNotification/set", (&sharenotification.Set{}).Name())
}

func TestShareNotificationQueryName(t *testing.T) {
	require.Equal(t, "ShareNotification/query", (&sharenotification.Query{}).Name())
}

func TestShareNotificationQueryChangesName(t *testing.T) {
	require.Equal(t, "ShareNotification/queryChanges", (&sharenotification.QueryChanges{}).Name())
}
