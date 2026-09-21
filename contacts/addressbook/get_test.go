package addressbook

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account:    "u1",
		IDs:        []jmap.ID{"ab1"},
		Properties: []string{"name", "myRights"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["AddressBook/get",{"accountId":"u1","ids":["ab1"],"properties":["name","myRights"]},"0"]]}`,
		string(data))
}

func TestGetInvokeWithResultReferences(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account: "u1",
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: "c1",
			Name:     "AddressBook/changes",
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
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["AddressBook/get",{"accountId":"u1","#ids":{"resultOf":"c1","name":"AddressBook/changes","path":"/created"},"#properties":{"resultOf":"c2","name":"Core/echo","path":"/properties"}},"0"]]}`,
		string(data))
}

func TestGetRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Get{}).Requires())
}
