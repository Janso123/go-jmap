package participantidentity_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/participantidentity"
	"github.com/stretchr/testify/require"
)

func TestParticipantIdentityIsObject(t *testing.T) {
	var _ jmap.Object = participantidentity.ParticipantIdentity{}
	require.Equal(t, "ParticipantIdentity", participantidentity.ParticipantIdentity{}.JMAPType())
	require.Equal(t, []jmap.URI{calendar.URI}, participantidentity.ParticipantIdentity{}.Requires())
}

func TestParticipantIdentityMethodNames(t *testing.T) {
	require.Equal(t, "ParticipantIdentity/get", (&participantidentity.Get{}).Name())
	require.Equal(t, "ParticipantIdentity/changes", (&participantidentity.Changes{}).Name())
	require.Equal(t, "ParticipantIdentity/set", (&participantidentity.Set{}).Name())
}

func TestParticipantIdentityEventType(t *testing.T) {
	require.Equal(t, jmap.EventType("ParticipantIdentity"), calendar.ParticipantIdentity)
}
