package sieve_test

import (
	jsonv2 "encoding/json/v2"
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
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &session))

	sessionCap, ok := session.Capabilities[sieve.URI].(*sieve.Capability)
	require.True(t, ok)
	assert.Equal(t, "Example Sieve 1.0", sessionCap.Implementation)

	accountCap, ok := session.Accounts["u1"].Capabilities[sieve.URI].(*sieve.Capability)
	require.True(t, ok)
	assert.Equal(t, jmap.UnsignedInt(512), accountCap.MaxSizeScriptName)
	maxSize, ok := accountCap.MaxSizeScript.Value()
	require.True(t, ok)
	assert.Equal(t, jmap.UnsignedInt(65536), maxSize)
	maxScripts, ok := accountCap.MaxNumberScripts.Value()
	require.True(t, ok)
	assert.Equal(t, jmap.UnsignedInt(5), maxScripts)
	assert.True(t, accountCap.MaxNumberRedirects.IsNull())
	assert.Equal(t, []string{"fileinto", "vacation", "enotify"}, accountCap.SieveExtensions)
	methods, ok := accountCap.NotificationMethods.Value()
	require.True(t, ok)
	assert.Equal(t, []string{"mailto"}, methods)
	assert.True(t, accountCap.ExternalLists.IsNull())
}

func TestSieveScriptJSON(t *testing.T) {
	script := sieve.SieveScript{
		ID:       jmap.ID("script1"),
		Name:     jmap.Some("vacation"),
		BlobID:   jmap.ID("blob1"),
		IsActive: new(true),
	}

	data, err := jsonv2.Marshal(script)
	require.NoError(t, err)

	assert.JSONEq(t, `{
	  "id": "script1",
	  "name": "vacation",
	  "blobId": "blob1",
	  "isActive": true
	}`, string(data))
}

func TestSieveScriptNameNull(t *testing.T) {
	t.Parallel()
	var got sieve.SieveScript
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"s1","name":null}`), &got))
	b, err := jsonv2.Marshal(got)
	require.NoError(t, err)
	require.Contains(t, string(b), `"name":null`)
}

func TestSieveScriptIsActiveFalseOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(sieve.SieveScript{ID: "s1", IsActive: new(false)})
	require.NoError(t, err)
	require.Contains(t, string(b), `"isActive":false`)
}
