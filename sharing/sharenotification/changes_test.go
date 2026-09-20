package sharenotification

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/sharing"
	"github.com/stretchr/testify/assert"
)

func TestChanges(t *testing.T) {
	m := &Changes{
		Account:    "account-id",
		SinceState: "1234",
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/changes",{"accountId":"account-id","sinceState":"1234"},"0"]]}`
	assert.Equal(t, expected, string(data))

	t.Run("manual", func(t *testing.T) {
		m := &Changes{
			Account:    "account-id",
			SinceState: "1234",
		}
		req = &jmap.Request{
			Using: []jmap.URI{sharing.URI},
			Calls: []*jmap.Invocation{
				{
					Name:   "ShareNotification/changes",
					Args:   m,
					CallID: "manual",
				},
			},
		}

		data, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/changes",{"accountId":"account-id","sinceState":"1234"},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
