package blob_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSessionAndAccount(t *testing.T) {
	raw := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:core": {"maxSizeUpload": 1},
	    "urn:ietf:params:jmap:blob": {}
	  },
	  "accounts": {
	    "u1": {
	      "name": "user@example.com",
	      "isPersonal": true,
	      "isReadOnly": false,
	      "accountCapabilities": {
	        "urn:ietf:params:jmap:blob": {
	          "maxSizeBlobSet": 1024,
	          "maxDataSources": 3,
	          "supportedTypeNames": ["image/png", "application/pdf"],
	          "supportedDigestAlgorithms": ["sha", "sha-256"]
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

	sessionCap, ok := session.Capabilities[blob.URI].(*blob.AccountCapability)
	require.True(t, ok)
	assert.Nil(t, sessionCap.MaxSizeBlobSet)
	assert.Zero(t, sessionCap.MaxDataSources)
	assert.Nil(t, sessionCap.SupportedTypeNames)
	assert.Nil(t, sessionCap.SupportedDigestAlgorithms)

	accountCap, ok := session.Accounts["u1"].Capabilities[blob.URI].(*blob.AccountCapability)
	require.True(t, ok)
	require.NotNil(t, accountCap.MaxSizeBlobSet)
	assert.Equal(t, uint64(1024), *accountCap.MaxSizeBlobSet)
	assert.Equal(t, uint64(3), accountCap.MaxDataSources)
	assert.Equal(t, []string{"image/png", "application/pdf"}, accountCap.SupportedTypeNames)
	assert.Equal(t, []string{"sha", "sha-256"}, accountCap.SupportedDigestAlgorithms)
}
