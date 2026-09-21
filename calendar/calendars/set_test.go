package calendars

import (
	"encoding/json"
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
		OnDestroyRemoveEvents: true,
		OnSuccessSetIsDefault: "#create-cal1",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["Calendar/set",{"accountId":"u1","onDestroyRemoveEvents":true,"onSuccessSetIsDefault":"#create-cal1"},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: map[jmap.ID]jmap.Patch{
			"cal1": {
				"description": nil,
			},
		},
	}

	data, err := json.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"cal1":{"description":null}}}`,
		string(data))
}

func TestSetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Set{}).Requires())
}
