package participantidentity

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetInvokeUsesCalendarsCapabilityAndExtras(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account:               "u1",
		OnSuccessSetIsDefault: jmap.Some(jmap.ID("#create-pi1")),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["ParticipantIdentity/set",{"accountId":"u1","onSuccessSetIsDefault":"#create-pi1"},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"pi1": {
				"name": nil,
			},
		}),
	}

	data, err := jsonv2.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"pi1":{"name":null}}}`,
		string(data))
}

func TestSetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Set{}).Requires())
}
