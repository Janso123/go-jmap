package addressbook

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetInvokeUsesContactsCapabilityAndExtras(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account:                 "u1",
		OnDestroyRemoveContents: true,
		OnSuccessSetIsDefault:   jmap.Some(jmap.ID("#create-ab1")),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["AddressBook/set",{"accountId":"u1","onDestroyRemoveContents":true,"onSuccessSetIsDefault":"#create-ab1"},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"ab1": {
				"description": nil,
			},
		}),
	}

	data, err := jsonv2.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"ab1":{"description":null}}}`,
		string(data))
}

func TestSetRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Set{}).Requires())
}

func TestOnSuccessSetIsDefaultNull(t *testing.T) {
	t.Parallel()
	var set Set
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"onSuccessSetIsDefault":null}`), &set))
	require.True(t, set.OnSuccessSetIsDefault.IsNull())
	data, err := jsonv2.Marshal(&set)
	require.NoError(t, err)
	require.Contains(t, string(data), `"onSuccessSetIsDefault":null`)
}
