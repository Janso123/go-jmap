package participantidentity

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetInvokeUsesCalendarsCapabilityAndExtras(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account:               "u1",
		OnSuccessSetIsDefault: "#create-pi1",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["ParticipantIdentity/set",{"accountId":"u1","onSuccessSetIsDefault":"#create-pi1"},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: map[jmap.ID]jmap.Patch{
			"pi1": {
				"name": nil,
			},
		},
	}

	data, err := json.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"pi1":{"name":null}}}`,
		string(data))
}

func TestSetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Set{}).Requires())
}
