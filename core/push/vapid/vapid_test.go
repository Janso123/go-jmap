package vapid_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/core/push/vapid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSession(t *testing.T) {
	raw := []byte(`{
		"capabilities": {
			"urn:ietf:params:jmap:webpush-vapid": {
				"applicationServerKey": "BASE64URLKEY"
			}
		},
		"accounts": {},
		"primaryAccounts": {},
		"username": "u",
		"apiUrl": "https://example/jmap",
		"downloadUrl": "https://example/dl",
		"uploadUrl": "https://example/up",
		"eventSourceUrl": "https://example/es",
		"state": "s1"
	}`)
	var sess jmap.Session
	require.NoError(t, json.Unmarshal(raw, &sess))
	c, ok := sess.Capabilities[vapid.URI].(*vapid.Capability)
	require.True(t, ok)
	assert.Equal(t, "BASE64URLKEY", c.ApplicationServerKey)
}
