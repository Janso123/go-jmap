package sieve_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/sieve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSessionAndAccount(t *testing.T) {
	raw := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:core": {"maxSizeUpload": 1},
	    "urn:ietf:params:jmap:sieve": {
	      "implementation": "Example Sieve 1.0"
	    }
	  },
	  "accounts": {
	    "u1": {
	      "name": "user@example.com",
	      "isPersonal": true,
	      "isReadOnly": false,
	      "accountCapabilities": {
	        "urn:ietf:params:jmap:sieve": {
	          "maxSizeScriptName": 512,
	          "maxSizeScript": 65536,
	          "maxNumberScripts": 5,
	          "maxNumberRedirects": null,
	          "sieveExtensions": ["fileinto", "vacation", "enotify"],
	          "notificationMethods": ["mailto"],
	          "externalLists": null
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

	sessionCap, ok := session.Capabilities[sieve.URI].(*sieve.Capability)
	require.True(t, ok)
	assert.Equal(t, "Example Sieve 1.0", sessionCap.Implementation)

	accountCap, ok := session.Accounts["u1"].Capabilities[sieve.URI].(*sieve.Capability)
	require.True(t, ok)
	assert.Equal(t, uint64(512), accountCap.MaxSizeScriptName)
	require.NotNil(t, accountCap.MaxSizeScript)
	assert.Equal(t, uint64(65536), *accountCap.MaxSizeScript)
	require.NotNil(t, accountCap.MaxNumberScripts)
	assert.Equal(t, uint64(5), *accountCap.MaxNumberScripts)
	assert.Nil(t, accountCap.MaxNumberRedirects)
	assert.Equal(t, []string{"fileinto", "vacation", "enotify"}, accountCap.SieveExtensions)
	assert.Equal(t, []string{"mailto"}, accountCap.NotificationMethods)
	assert.Nil(t, accountCap.ExternalLists)
}

func TestSieveScriptJSON(t *testing.T) {
	script := sieve.SieveScript{
		ID:       jmap.ID("script1"),
		Name:     "vacation",
		BlobID:   jmap.ID("blob1"),
		IsActive: true,
	}

	data, err := json.Marshal(script)
	require.NoError(t, err)

	assert.JSONEq(t, `{
	  "id": "script1",
	  "name": "vacation",
	  "blobId": "blob1",
	  "isActive": true
	}`, string(data))
}
