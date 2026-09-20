package principal

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
			"principal-id": {
				"name": "Jane Doe",
			},
		},
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["Principal/set",{"accountId":"account-id","update":{"principal-id":{"name":"Jane Doe"}}},"0"]]}`
	assert.Equal(t, expected, string(data))

	t.Run("manual", func(t *testing.T) {
		m := &Set{
			Account: "account-id",
			Update: map[jmap.ID]jmap.Patch{
				"principal-id": {
					"description": nil,
				},
			},
		}
		req = &jmap.Request{
			Using: []jmap.URI{sharing.URI},
			Calls: []*jmap.Invocation{
				{
					Name:   "Principal/set",
					Args:   m,
					CallID: "manual",
				},
			},
		}

		data, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["Principal/set",{"accountId":"account-id","update":{"principal-id":{"description":null}}},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
