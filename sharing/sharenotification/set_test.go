package sharenotification

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing"
	"github.com/stretchr/testify/assert"
)

func TestSetDestroyOnly(t *testing.T) {
	m := &Set{
		Account: "account-id",
		Destroy: []jmap.ID{"notification-id"},
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/set",{"accountId":"account-id","destroy":["notification-id"]},"0"]]}`
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

		data, err := jsonv2.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/set",{"accountId":"account-id","destroy":["notification-id"]},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
