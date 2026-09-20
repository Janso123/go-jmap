package addressbook

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetInvokeUsesContactsCapabilityAndExtras(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account:                 "u1",
		OnDestroyRemoveContents: true,
		OnSuccessSetIsDefault:   "#create-ab1",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["AddressBook/set",{"accountId":"u1","onDestroyRemoveContents":true,"onSuccessSetIsDefault":"#create-ab1"},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: map[jmap.ID]jmap.Patch{
			"ab1": {
				"description": nil,
			},
		},
	}

	data, err := json.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"ab1":{"description":null}}}`,
		string(data))
}

func TestSetRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Set{}).Requires())
}
