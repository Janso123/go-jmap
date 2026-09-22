package blob_test

import (
	jsonv2 "encoding/json/v2"
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
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &session))

	sessionCap, ok := session.Capabilities[blob.URI].(*blob.AccountCapability)
	require.True(t, ok)
	assert.True(t, sessionCap.MaxSizeBlobSet.IsZero())
	assert.Zero(t, sessionCap.MaxDataSources)
	assert.Nil(t, sessionCap.SupportedTypeNames)
	assert.Nil(t, sessionCap.SupportedDigestAlgorithms)

	accountCap, ok := session.Accounts["u1"].Capabilities[blob.URI].(*blob.AccountCapability)
	require.True(t, ok)
	maxSize, okMax := accountCap.MaxSizeBlobSet.Value()
	require.True(t, okMax)
	assert.Equal(t, jmap.UnsignedInt(1024), maxSize)
	assert.Equal(t, jmap.UnsignedInt(3), accountCap.MaxDataSources)
	assert.Equal(t, []string{"image/png", "application/pdf"}, accountCap.SupportedTypeNames)
	assert.Equal(t, []string{"sha", "sha-256"}, accountCap.SupportedDigestAlgorithms)
}

func TestAccountCapabilityNullMaxSizeBlobSet(t *testing.T) {
	t.Parallel()
	var cap blob.AccountCapability
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxSizeBlobSet":null}`), &cap))
	b, err := jsonv2.Marshal(&cap)
	require.NoError(t, err)
	require.Contains(t, string(b), `"maxSizeBlobSet":null`)
}
