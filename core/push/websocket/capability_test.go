package websocket_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSession(t *testing.T) {
	blob := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:core": {"maxSizeUpload": 1},
	    "urn:ietf:params:jmap:websocket": {
	      "url": "wss://server.example.com/jmap/ws/",
	      "supportsPush": true
	    }
	  },
	  "accounts": {},
	  "primaryAccounts": {},
	  "username": "u",
	  "apiUrl": "https://server.example.com/jmap/api/",
	  "downloadUrl": "https://server.example.com/dl/{accountId}/{blobId}/{name}?accept={type}",
	  "uploadUrl": "https://server.example.com/upload/{accountId}/",
	  "eventSourceUrl": "https://server.example.com/es/",
	  "state": "s1"
	}`
	s := &jmap.Session{}
	require.NoError(t, json.Unmarshal([]byte(blob), s))
	cap, ok := s.Capabilities[websocket.URI].(*websocket.WebSocket)
	require.True(t, ok)
	assert.Equal(t, "wss://server.example.com/jmap/ws/", cap.URL)
	assert.True(t, cap.SupportsPush)
}

func TestStateChangePushState(t *testing.T) {
	raw := `{"@type":"StateChange","changed":{"a1":{"Email":"e1"}},"pushState":"bbb"}`
	var sc jmap.StateChange
	require.NoError(t, json.Unmarshal([]byte(raw), &sc))
	assert.Equal(t, "bbb", sc.PushState)
}
