package contacts_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSessionAndAccount(t *testing.T) {
	raw := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:contacts": {}
	  },
	  "accounts": {
	    "u1": {
	      "name": "user@example.com",
	      "isPersonal": true,
	      "isReadOnly": false,
	      "accountCapabilities": {
	        "urn:ietf:params:jmap:contacts": {
	          "maxAddressBooksPerCard": 7,
	          "mayCreateAddressBook": true
	        }
	      }
	    }
	  },
	  "primaryAccounts": {},
	  "username": "u",
	  "apiUrl": "https://server.example.com/jmap/api/",
	  "downloadUrl": "https://server.example.com/dl/{accountId}/{blobId}/{name}?accept={type}",
	  "uploadUrl": "https://server.example.com/upload/{accountId}/",
	  "eventSourceUrl": "https://server.example.com/es/",
	  "state": "s1"
	}`

	var session jmap.Session
	require.NoError(t, json.Unmarshal([]byte(raw), &session))

	sessionCap, ok := session.Capabilities[contacts.URI].(*contacts.Capability)
	require.True(t, ok)
	assert.Nil(t, sessionCap.MaxAddressBooksPerCard)
	assert.Nil(t, sessionCap.MayCreateAddressBook)

	accountCap, ok := session.Accounts["u1"].Capabilities[contacts.URI].(*contacts.Capability)
	require.True(t, ok)
	require.NotNil(t, accountCap.MaxAddressBooksPerCard)
	assert.Equal(t, uint64(7), *accountCap.MaxAddressBooksPerCard)
	require.NotNil(t, accountCap.MayCreateAddressBook)
	assert.True(t, *accountCap.MayCreateAddressBook)
}
