package sharenotification

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/sharing"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	m := &Set{
		Account: "account-id",
		Update: map[jmap.ID]jmap.Patch{
			"notification-id": {
				"newRights/mayReadItems": true,
			},
		},
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/set",{"accountId":"account-id","update":{"notification-id":{"newRights/mayReadItems":true}}},"0"]]}`
	assert.Equal(t, expected, string(data))

	t.Run("manual", func(t *testing.T) {
		m := &Set{
			Account: "account-id",
			Destroy: []jmap.ID{"notification-id"},
		}
		req = &jmap.Request{
			Using: []jmap.URI{sharing.URI},
			Calls: []*jmap.Invocation{
				{
					Name:   "ShareNotification/set",
					Args:   m,
					CallID: "manual",
				},
			},
		}

		data, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/set",{"accountId":"account-id","destroy":["notification-id"]},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
