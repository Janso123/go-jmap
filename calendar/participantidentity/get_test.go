package participantidentity

import (
	"encoding/json"
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
		IDs:        []jmap.ID{"pi1"},
		Properties: []string{"name", "calendarAddress"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["ParticipantIdentity/get",{"accountId":"u1","ids":["pi1"],"properties":["name","calendarAddress"]},"0"]]}`,
		string(data))
}

func TestGetInvokeWithResultReferences(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account: "u1",
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: "c1",
			Name:     "ParticipantIdentity/changes",
			Path:     "/created",
		},
		ReferenceProperties: &jmap.ResultReference{
			ResultOf: "c2",
			Name:     "Core/echo",
			Path:     "/properties",
		},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["ParticipantIdentity/get",{"accountId":"u1","#ids":{"resultOf":"c1","name":"ParticipantIdentity/changes","path":"/created"},"#properties":{"resultOf":"c2","name":"Core/echo","path":"/properties"}},"0"]]}`,
		string(data))
}

func TestGetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Get{}).Requires())
}
