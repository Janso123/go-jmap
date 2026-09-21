package sharing_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSessionAndAccount(t *testing.T) {
	raw := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:principals": {
	      "currentUserPrincipalId": "p1"
	    }
	  },
	  "accounts": {
	    "a1": {
	      "name": "user@example.com",
	      "isPersonal": true,
	      "isReadOnly": false,
	      "accountCapabilities": {
	        "urn:ietf:params:jmap:principals:owner": {
	          "accountIdForPrincipal": "a1",
	          "principalId": "p1"
	        }
	      }
	    }
	  },
	  "primaryAccounts": {},
	  "username": "user@example.com",
	  "apiUrl": "https://server.example.com/jmap/api/",
	  "downloadUrl": "https://server.example.com/dl/{accountId}/{blobId}/{name}?accept={type}",
	  "uploadUrl": "https://server.example.com/upload/{accountId}/",
	  "eventSourceUrl": "https://server.example.com/es/",
	  "state": "s1"
	}`

	var session jmap.Session
	require.NoError(t, json.Unmarshal([]byte(raw), &session))

	sessionCap, ok := session.Capabilities[sharing.URI].(*sharing.Capability)
	require.True(t, ok)
	require.NotNil(t, sessionCap.CurrentUserPrincipalID)
	assert.Equal(t, jmap.ID("p1"), *sessionCap.CurrentUserPrincipalID)

	accountCap, ok := session.Accounts["a1"].Capabilities[sharing.OwnerURI].(*sharing.OwnerCapability)
	require.True(t, ok)
	assert.Equal(t, jmap.ID("a1"), accountCap.AccountIDForPrincipal)
	assert.Equal(t, jmap.ID("p1"), accountCap.PrincipalID)
}
