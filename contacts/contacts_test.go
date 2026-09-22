package contacts_test

import (
	jsonv2 "encoding/json/v2"
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
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &session))

	sessionCap, ok := session.Capabilities[contacts.URI].(*contacts.Capability)
	require.True(t, ok)
	assert.True(t, sessionCap.MaxAddressBooksPerCard.IsZero())
	assert.Nil(t, sessionCap.MayCreateAddressBook)

	accountCap, ok := session.Accounts["u1"].Capabilities[contacts.URI].(*contacts.Capability)
	require.True(t, ok)
	perCard, ok := accountCap.MaxAddressBooksPerCard.Value()
	require.True(t, ok)
	assert.Equal(t, jmap.UnsignedInt(7), perCard)
	require.NotNil(t, accountCap.MayCreateAddressBook)
	assert.True(t, *accountCap.MayCreateAddressBook)
}

func TestMaxAddressBooksPerCardNull(t *testing.T) {
	t.Parallel()
	var cap contacts.Capability
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxAddressBooksPerCard":null}`), &cap))
	require.True(t, cap.MaxAddressBooksPerCard.IsNull())
	data, err := jsonv2.Marshal(&cap)
	require.NoError(t, err)
	require.Contains(t, string(data), `"maxAddressBooksPerCard":null`)
}
