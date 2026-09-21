package principal_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/sharing"
	"github.com/Janso123/go-jmap/sharing/principal"
	"github.com/stretchr/testify/require"
)

func TestPrincipalIsObject(t *testing.T) {
	var _ jmap.Object = principal.Principal{}
	require.Equal(t, "Principal", principal.Principal{}.JMAPType())
	require.Equal(t, []jmap.URI{sharing.URI}, principal.Principal{}.Requires())
}

func TestPrincipalGetName(t *testing.T) {
	require.Equal(t, "Principal/get", (&principal.Get{}).Name())
}

func TestPrincipalChangesName(t *testing.T) {
	require.Equal(t, "Principal/changes", (&principal.Changes{}).Name())
}

func TestPrincipalSetName(t *testing.T) {
	require.Equal(t, "Principal/set", (&principal.Set{}).Name())
}

func TestPrincipalQueryName(t *testing.T) {
	require.Equal(t, "Principal/query", (&principal.Query{}).Name())
}

func TestPrincipalQueryChangesName(t *testing.T) {
	require.Equal(t, "Principal/queryChanges", (&principal.QueryChanges{}).Name())
}

func TestPrincipalGetAvailabilityRemainsManual(t *testing.T) {
	require.Equal(t, "Principal/getAvailability", (&principal.GetAvailability{}).Name())
	require.Equal(t, []jmap.URI{sharing.URI, calendar.AvailabilityURI}, (&principal.GetAvailability{}).Requires())
}
